package dao

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"gorm.io/gorm"

	mappers "github.com/Orange-Health/citadel/apps/samples/mappers"
	"github.com/Orange-Health/citadel/apps/samples/structures"
	commonConstants "github.com/Orange-Health/citadel/common/constants"
	commonStructures "github.com/Orange-Health/citadel/common/structures"
	commonUtils "github.com/Orange-Health/citadel/common/utils"
	commonModels "github.com/Orange-Health/citadel/models"
)

func getStringsForReverseLogisiticsSamples() []string {
	return []string{
		"order_details.oms_request_id AS request_id",
		"order_details.oms_order_id AS order_id",
		"order_details.city_code as city_code",
		"patient_details.name AS patient_name",
		"order_details.trf_id AS trf_id",
		"patient_details.expected_dob AS patient_age",
		"patient_details.gender AS patient_gender",
		"samples.barcode",
		"samples.vial_type_id",
		"sample_metadata.collected_at AS sample_collected_time",
	}
}

func (sampleDao *SampleDao) BeginTransaction() *gorm.DB {
	return sampleDao.Db.Begin()
}

func (sampleDao *SampleDao) CreateSampleWithTx(tx *gorm.DB, sampleInfo commonStructures.SampleInfo) (
	commonStructures.SampleInfo, *commonStructures.CommonError) {
	sample := mappers.MapSampleInfoToSample(sampleInfo)
	sampleMetadata := mappers.MapSampleInfoToSampleMetadata(sampleInfo)

	if err := tx.Create(&sample).Error; err != nil {
		return sampleInfo, commonUtils.HandleORMError(err)
	}

	sampleMetadata.SampleId = sample.Id
	if err := tx.Create(&sampleMetadata).Error; err != nil {
		return sampleInfo, commonUtils.HandleORMError(err)
	}

	sampleInfo.Id = sample.Id

	return sampleInfo, nil
}

func (sampleDao *SampleDao) AssignTaskSequenceToSamples(tx *gorm.DB, taskId uint, alnumTestIds []string, isAdditionalTask uint,
	omsRequestId string) *commonStructures.CommonError {
	taskSequence := taskId
	if isAdditionalTask == 0 {
		taskSequence = 0 // For primary collection tasks, task sequence is always 0
	}
	sampleIds := []uint{}
	err := tx.Model(&commonModels.Sample{}).
		Joins("INNER JOIN test_sample_mapping ON samples.oms_order_id = test_sample_mapping.oms_order_id AND samples.sample_number = test_sample_mapping.sample_number").
		Where("samples.status = ?", commonConstants.SampleDefault).
		Where("samples.oms_request_id = ?", omsRequestId).
		Where("test_sample_mapping.oms_test_id IN (?)", alnumTestIds).
		Where("samples.deleted_at IS NULL AND test_sample_mapping.deleted_at IS NULL").
		Pluck("samples.id", &sampleIds).Error
	if err != nil {
		return commonUtils.HandleORMError(err)
	}
	if len(sampleIds) == 0 {
		return nil // No samples found, nothing to update
	}

	if err := tx.Model(&commonModels.SampleMetadata{}).
		Where("sample_id IN (?)", sampleIds).
		Updates(map[string]interface{}{
			"task_sequence": taskSequence,
			"updated_by":    commonConstants.CitadelSystemId,
		}).Error; err != nil {
		return commonUtils.HandleORMError(err)
	}

	return nil
}

func (sampleDao *SampleDao) UpdateBulkSamples(samples []commonStructures.SampleInfo) ([]commonStructures.SampleInfo,
	*commonStructures.CommonError) {
	sample := mappers.MapBulkSampleInfoToSample(samples)
	sampleMetadata := mappers.MapBulkSampleInfoToSampleMetadata(samples)

	if len(sample) == 0 {
		return samples, &commonStructures.CommonError{
			StatusCode: http.StatusNotFound,
			Message:    commonConstants.ERROR_SAMPLE_NOT_FOUND,
		}
	}

	if txErr := sampleDao.Db.Transaction(func(tx *gorm.DB) error {
		for idx := range sample {
			if err := tx.Where("id = ?", sample[idx].Id).Updates(&sample[idx]).Error; err != nil {
				return err
			}
		}

		for idx := range sampleMetadata {
			if err := tx.Where("sample_id = ?", sampleMetadata[idx].SampleId).
				Updates(&sampleMetadata[idx]).Error; err != nil {
				return err
			}
		}

		return nil
	}); txErr != nil {
		return samples, commonUtils.HandleORMError(txErr)
	}

	return mappers.MapBulkSampleSampleMetaToSampleInfo(sample, sampleMetadata), nil
}

func (sampleDao *SampleDao) DeleteSampleByOrderIdAndSampleNumberWithTx(tx *gorm.DB, omsOrderId string,
	sampleNumber uint) *commonStructures.CommonError {
	var sampleId uint
	currentTime := commonUtils.GetCurrentTime()
	if err := tx.Model(&commonModels.Sample{}).
		Where("oms_order_id = ? AND sample_number = ?", omsOrderId, sampleNumber).
		Select("id").
		Scan(&sampleId).Error; err != nil {
		return commonUtils.HandleORMError(err)
	}

	if err := tx.Model(&commonModels.Sample{}).
		Where("id = ?", sampleId).
		Updates(map[string]interface{}{
			"status":     commonConstants.SampleDeleted,
			"deleted_at": currentTime,
			"deleted_by": commonConstants.CitadelSystemId,
		}).Error; err != nil {
		return commonUtils.HandleORMError(err)
	}

	if err := tx.Model(&commonModels.SampleMetadata{}).
		Where("sample_id = ?", sampleId).
		Updates(map[string]interface{}{
			"deleted_at": currentTime,
			"deleted_by": commonConstants.CitadelSystemId,
		}).Error; err != nil {
		return commonUtils.HandleORMError(err)
	}
	return nil
}

func (sampleDao *SampleDao) GetSampleForCollectedB2C(omsRequestId string, taskId uint, statusArray []string,
	appendTaskSequence, useSampleNumber bool) ([]commonStructures.SampleInfo, *commonStructures.CommonError) {
	var samples []commonStructures.SampleInfo
	query := sampleDao.Db.Table(commonConstants.TableSamples).
		Select("sample_metadata.*, samples.*").
		Joins("INNER JOIN sample_metadata ON samples.id = sample_metadata.sample_id")
	if useSampleNumber {
		query = query.Joins("INNER JOIN test_sample_mapping ON samples.oms_order_id = test_sample_mapping.oms_order_id AND samples.sample_number = test_sample_mapping.sample_number")
	}
	query = query.Where("samples.oms_request_id = ?", omsRequestId).
		Where("samples.status IN (?)", statusArray).
		Where("samples.deleted_at IS NULL").
		Where("sample_metadata.deleted_at IS NULL")
	if useSampleNumber {
		query = query.Where("test_sample_mapping.deleted_at IS NULL")
	}

	if appendTaskSequence {
		query = query.Where("sample_metadata.task_sequence = ?", taskId)
	}

	if err := query.Find(&samples).Error; err != nil {
		return nil, commonUtils.HandleORMError(err)
	}
	if len(samples) == 0 {
		return nil, &commonStructures.CommonError{
			StatusCode: http.StatusNotFound,
			Message:    commonConstants.ERROR_SAMPLE_NOT_FOUND,
		}
	}

	return samples, nil
}

func (sampleDao *SampleDao) GetAllSampleTestsBySampleNumber(sampleNumber uint, omsOrderId string) ([]commonModels.TestDetail,
	*commonStructures.CommonError) {
	testDetails := []commonModels.TestDetail{}
	if err := sampleDao.Db.Table("test_details").
		Joins("INNER JOIN test_sample_mapping ON test_details.central_oms_test_id = test_sample_mapping.oms_test_id AND test_details.oms_order_id = test_sample_mapping.oms_order_id").
		Where("test_sample_mapping.sample_number = ?", sampleNumber).
		Where("test_sample_mapping.oms_order_id = ?", omsOrderId).
		Where("test_sample_mapping.is_rejected = false").
		Where("test_sample_mapping.deleted_at is null").
		Find(&testDetails).Error; err != nil {
		return testDetails, commonUtils.HandleORMError(err)
	}

	return testDetails, nil
}

func (sampleDao *SampleDao) GetAllSampleTestsBySampleNumbers(sampleNumbers []uint, omsOrderId string) (
	[]commonModels.TestDetail, *commonStructures.CommonError) {
	testDetails := []commonModels.TestDetail{}
	if err := sampleDao.Db.Table("test_details").
		Joins("INNER JOIN test_sample_mapping ON test_details.central_oms_test_id = test_sample_mapping.oms_test_id AND test_details.oms_order_id = test_sample_mapping.oms_order_id").
		Where("test_sample_mapping.sample_number IN (?)", sampleNumbers).
		Where("test_sample_mapping.oms_order_id = ?", omsOrderId).
		Where("test_sample_mapping.is_rejected = false").
		Where("test_sample_mapping.deleted_at is null").
		Find(&testDetails).Error; err != nil {
		return testDetails, commonUtils.HandleORMError(err)
	}

	return testDetails, nil
}

func (sampleDao *SampleDao) GetSamplesByOmsOrderId(omsOrderId string) ([]commonStructures.SampleInfo,
	*commonStructures.CommonError) {
	samples, samplesMetadata := []commonModels.Sample{}, []commonModels.SampleMetadata{}

	if err := sampleDao.Db.Table("samples").
		Where("oms_order_id = ?", omsOrderId).
		Find(&samples).Error; err != nil {
		return nil, commonUtils.HandleORMError(err)
	}

	if err := sampleDao.Db.Table("sample_metadata").
		Where("oms_order_id = ?", omsOrderId).
		Find(&samplesMetadata).Error; err != nil {
		return nil, commonUtils.HandleORMError(err)
	}

	if len(samples) == 0 {
		return nil, &commonStructures.CommonError{
			StatusCode: http.StatusNotFound,
			Message:    commonConstants.ERROR_SAMPLE_NOT_FOUND,
		}
	}

	return mappers.MapBulkSampleSampleMetaToSampleInfo(samples, samplesMetadata), nil
}

func (sampleDao *SampleDao) GetSamplesByOmsOrderIds(omsOrderIds []string) ([]commonStructures.SampleInfo,
	*commonStructures.CommonError) {
	samples, samplesMetadata := []commonModels.Sample{}, []commonModels.SampleMetadata{}
	if err := sampleDao.Db.Table("samples").
		Where("oms_order_id IN (?)", omsOrderIds).
		Find(&samples).Error; err != nil {
		return nil, commonUtils.HandleORMError(err)
	}

	if len(samples) == 0 {
		return nil, &commonStructures.CommonError{
			StatusCode: http.StatusNotFound,
			Message:    commonConstants.ERROR_SAMPLE_NOT_FOUND,
		}
	}

	if err := sampleDao.Db.Table("sample_metadata").
		Where("oms_order_id IN (?)", omsOrderIds).
		Find(&samplesMetadata).Error; err != nil {
		return nil, commonUtils.HandleORMError(err)
	}

	return mappers.MapBulkSampleSampleMetaToSampleInfo(samples, samplesMetadata), nil
}

func getCaseStatementSampleInfoStatusUint(inField, outField string) string {
	var caseStatements []string
	for key, value := range commonConstants.SampleStatusMap {
		caseStatements = append(caseStatements, fmt.Sprintf("WHEN '%s' THEN %d", key, value))
	}
	return fmt.Sprintf("CASE %s %s END AS %s", inField, strings.Join(caseStatements, " "), outField)
}

func (sampleDao *SampleDao) GetSampleByOrderIdAndSampleNumber(omsOrderId string, sampleNumber uint) (
	commonStructures.SampleInfo, *commonStructures.CommonError) {
	sampleInfo := commonStructures.SampleInfo{}

	if err := sampleDao.Db.Table("samples").
		Select("sample_metadata.*, samples.*, "+getCaseStatementSampleInfoStatusUint("samples.status", "status_uint")).
		Joins("INNER JOIN sample_metadata ON samples.id = sample_metadata.sample_id").
		Where("samples.oms_order_id = ? AND samples.sample_number = ?", omsOrderId, sampleNumber).
		Find(&sampleInfo); err.Error != nil {
		return sampleInfo, commonUtils.HandleORMError(err.Error)
	}
	return sampleInfo, nil
}

func (sampleDao *SampleDao) GetSamplesByOmsOrderIdAndSampleNumbers(omsOrderId string, sampleNumbers []uint) (
	[]commonModels.Sample, []commonModels.SampleMetadata, *commonStructures.CommonError) {
	samples, sampleMetadatas := []commonModels.Sample{}, []commonModels.SampleMetadata{}
	if err := sampleDao.Db.Model(&commonModels.Sample{}).
		Where("samples.oms_order_id = ? AND samples.sample_number IN (?)", omsOrderId, sampleNumbers).
		Find(&samples).Error; err != nil {
		return samples, sampleMetadatas, commonUtils.HandleORMError(err)
	}

	if len(samples) == 0 {
		return samples, sampleMetadatas, &commonStructures.CommonError{
			StatusCode: http.StatusNotFound,
			Message:    commonConstants.ERROR_SAMPLE_NOT_FOUND,
		}
	}

	sampleIds := []uint{}
	for _, sample := range samples {
		sampleIds = append(sampleIds, sample.Id)
	}

	if err := sampleDao.Db.Model(&commonModels.SampleMetadata{}).
		Where("sample_metadata.sample_id IN (?)", sampleIds).
		Find(&sampleMetadatas).Error; err != nil {
		return samples, sampleMetadatas, commonUtils.HandleORMError(err)
	}

	return samples, sampleMetadatas, nil
}

func (sampleDao *SampleDao) GetCovidTestSamples(omsOrderId string) []commonStructures.SampleInfo {
	sampleInfos := []commonStructures.SampleInfo{}
	samples, sampleMetadatas := []commonModels.Sample{}, []commonModels.SampleMetadata{}
	if err := sampleDao.Db.Table(commonConstants.TableSamples).
		Joins("INNER JOIN test_sample_mapping ON samples.oms_order_id = test_sample_mapping.oms_order_id AND samples.sample_number = test_sample_mapping.sample_number").
		Joins("INNER JOIN test_details ON test_sample_mapping.oms_test_id = test_details.central_oms_test_id").
		Where("samples.oms_order_id = ?", omsOrderId).
		Where("test_details.master_test_id IN (?)", commonUtils.GetCovid19MasterTestIds()).
		Where("samples.deleted_at IS NULL AND test_sample_mapping.deleted_at IS NULL AND test_details.deleted_at IS NULL").
		Select("samples.*").Find(&samples).Error; err != nil {
		return sampleInfos
	}
	if len(samples) == 0 {
		return sampleInfos
	}

	sampleIds := []uint{}
	for _, sample := range samples {
		sampleIds = append(sampleIds, sample.Id)
	}
	if err := sampleDao.Db.Table(commonConstants.TableSampleMetadata).
		Where("sample_id IN (?)", sampleIds).
		Where("deleted_at IS NULL").
		Find(&sampleMetadatas).Error; err != nil {
		return sampleInfos
	}
	if len(sampleMetadatas) == 0 {
		return sampleInfos
	}
	return mappers.MapBulkSampleSampleMetaToSampleInfo(samples, sampleMetadatas)
}

func (sampleDao *SampleDao) GetMaxSampleNumberByOmsOrderId(omsOrderId string) uint {
	var sampleNumber uint
	sampleDao.Db.Unscoped().Table("samples").
		Where("oms_order_id = ?", omsOrderId).
		Select("COALESCE(MAX(sample_number), 0)").
		Scan(&sampleNumber)
	return sampleNumber
}

func (sampleDao *SampleDao) GetMaxSampleNumberByOmsOrderIdWithTx(tx *gorm.DB, omsOrderId string) uint {
	var sampleNumber uint
	tx.Unscoped().Table("samples").
		Where("oms_order_id = ?", omsOrderId).
		Select("COALESCE(MAX(sample_number), 0)").
		Scan(&sampleNumber)
	return sampleNumber
}

func (sampleDao *SampleDao) GetSamplesForTests(omsTestIds []string) ([]commonModels.Sample, []commonModels.SampleMetadata,
	*commonStructures.CommonError) {
	samples, sampleMetadatas := []commonModels.Sample{}, []commonModels.SampleMetadata{}
	if err := sampleDao.Db.Table("samples").
		Joins("INNER JOIN test_sample_mapping ON samples.oms_order_id = test_sample_mapping.oms_order_id AND samples.sample_number = test_sample_mapping.sample_number").
		Where("test_sample_mapping.oms_test_id IN (?)", omsTestIds).
		Where("test_sample_mapping.deleted_at IS NULL AND samples.deleted_at IS NULL").
		Find(&samples).Error; err != nil {
		return samples, sampleMetadatas, commonUtils.HandleORMError(err)
	}
	if len(samples) == 0 {
		return samples, sampleMetadatas, &commonStructures.CommonError{
			StatusCode: http.StatusNotFound,
			Message:    commonConstants.ERROR_SAMPLE_NOT_FOUND,
		}
	}

	uniqueSampleIdsMap, uniqueSamples, sampleIds := map[uint]bool{}, []commonModels.Sample{}, []uint{}
	for _, sample := range samples {
		if uniqueSampleIdsMap[sample.Id] {
			continue
		}
		uniqueSampleIdsMap[sample.Id] = true
		uniqueSamples = append(uniqueSamples, sample)
		sampleIds = append(sampleIds, sample.Id)
	}

	if err := sampleDao.Db.Table("sample_metadata").
		Where("sample_id IN (?)", sampleIds).
		Where("deleted_at IS NULL").
		Find(&sampleMetadatas).Error; err != nil {
		return uniqueSamples, sampleMetadatas, commonUtils.HandleORMError(err)
	}
	if len(sampleMetadatas) == 0 {
		return uniqueSamples, sampleMetadatas, &commonStructures.CommonError{
			StatusCode: http.StatusNotFound,
			Message:    commonConstants.ERROR_SAMPLE_NOT_FOUND,
		}
	}

	return uniqueSamples, sampleMetadatas, nil
}

func (sampleDao *SampleDao) UpdateTaskIdByOmsTestIdsWithTx(tx *gorm.DB, omsTestIds []string,
	taskId uint) *commonStructures.CommonError {

	sampleIds := []uint{}
	if err := tx.Model(&commonModels.Sample{}).
		Joins("INNER JOIN test_sample_mapping ON samples.oms_order_id = test_sample_mapping.oms_order_id AND samples.sample_number = test_sample_mapping.sample_number").
		Where("test_sample_mapping.oms_test_id IN (?)", omsTestIds).
		Where("test_sample_mapping.is_rejected = false").
		Where("samples.deleted_at IS NULL AND test_sample_mapping.deleted_at IS NULL").
		Pluck("samples.id", &sampleIds).Error; err != nil {
		return commonUtils.HandleORMError(err)
	}

	if len(sampleIds) == 0 {
		return nil
	}

	sampleMetadataUpdates := map[string]interface{}{
		"task_sequence": taskId,
		"updated_by":    commonConstants.CitadelSystemId,
	}

	if err := tx.Model(&commonModels.SampleMetadata{}).
		Where("sample_id IN (?)", sampleIds).
		Updates(sampleMetadataUpdates).Error; err != nil {
		return commonUtils.HandleORMError(err)
	}

	return nil
}

func (sampleDao *SampleDao) GetSampleDetailsForScheduler(sampleDetailsRequest structures.SampleDetailsRequest) (
	[]structures.VialStructForQuery, *commonStructures.CommonError) {

	vialTypeStructsFromDb := []structures.VialStructForQuery{}
	selectStrings := []string{
		"samples.id as accession_id",
		"samples.oms_order_id as order_id",
		"samples.barcode as barcode",
		"sample_metadata.barcode_image_url as image_url",
		"sample_metadata.barcode_scanned_at as barcode_scanned_at",
		"sample_metadata.collected_volume as collected_volume",
		"sample_metadata.task_sequence as task_sequence",
		"sample_metadata.collect_later_reason as reason_for_skip",
		"samples.vial_type_id as vial_type",
		"JSONB_AGG(DISTINCT jsonb_build_object('id',test_details.central_oms_test_id)) as tests",
	}

	query := sampleDao.Db.Table("samples").
		Select(selectStrings).
		Joins("INNER JOIN sample_metadata ON samples.id = sample_metadata.sample_id").
		Joins("INNER JOIN test_sample_mapping ON samples.sample_number = test_sample_mapping.sample_number AND samples.oms_order_id = test_sample_mapping.oms_order_id").
		Joins("INNER JOIN test_details on test_sample_mapping.oms_test_id = test_details.central_oms_test_id").
		Where("test_details.deleted_at is NULL").
		Where("samples.oms_order_id IN (?)", sampleDetailsRequest.OrderIds).
		Where("samples.deleted_at IS NULL")

	if sampleDetailsRequest.IsAdditionalTask {
		query = query.Where("sample_metadata.task_sequence = ?", sampleDetailsRequest.TaskId).
			Where("samples.status IN (?)", []string{commonConstants.SampleDefault})
	} else {
		if sampleDetailsRequest.IsCampTask {
			sampleStatuses := []string{commonConstants.SampleDefault, commonConstants.SampleCollectionDone,
				commonConstants.SampleAccessioned, commonConstants.SampleNotCollectedEmedic}
			query = query.Where("samples.status IN (?)", sampleStatuses)
		} else {
			query = query.Where("sample_metadata.task_sequence = ?", 0).
				Where("samples.status IN (?)", []string{commonConstants.SampleDefault})
		}
	}

	query = query.Group("samples.id, samples.oms_order_id, samples.barcode, sample_metadata.barcode_image_url, sample_metadata.barcode_scanned_at, sample_metadata.collected_volume, sample_metadata.task_sequence, sample_metadata.collect_later_reason, samples.vial_type_id")
	if err := query.Scan(&vialTypeStructsFromDb).Error; err != nil {
		return nil, commonUtils.HandleORMError(err)
	}

	return vialTypeStructsFromDb, nil
}

func (sampleDao *SampleDao) GetTestDetailsForLisEventByVisitId(visitId string) (
	[]commonStructures.TestDetailsForLisEvent, *commonStructures.CommonError) {

	testDetailsForLisEvent := []commonStructures.TestDetailsForLisEvent{}

	selectStrings := []string{
		"test_details.central_oms_test_id as test_id",
		"test_details.lis_code as test_code",
		"test_details.master_test_id as master_test_id",
		"test_details.master_package_id as master_package_id",
		"test_details.test_name as test_name",
		"test_details.test_type as test_type",
		"STRING_AGG(DISTINCT samples.barcode, ', ') as barcodes",
	}

	sampleStatuses := []string{commonConstants.SamplePartiallyRejected, commonConstants.SampleSynced,
		commonConstants.SampleReceived, commonConstants.SampleAccessioned}
	if err := sampleDao.Db.Table(commonConstants.TableTestDetails).
		Select(selectStrings).
		Joins("INNER JOIN test_sample_mapping ON test_details.oms_order_id = test_sample_mapping.oms_order_id and test_sample_mapping.oms_test_id = test_details.central_oms_test_id").
		Joins("INNER JOIN samples ON samples.oms_order_id = test_sample_mapping.oms_order_id and samples.sample_number = test_sample_mapping.sample_number").
		Where("samples.visit_id = ?", visitId).
		Where("samples.status IN (?)", sampleStatuses).
		Where("test_sample_mapping.is_rejected = false").
		Where("test_details.deleted_at IS NULL").
		Where("test_sample_mapping.deleted_at IS NULL").
		Where("samples.deleted_at IS NULL").
		Group("test_details.central_oms_test_id, test_details.lis_code, test_details.master_test_id, test_details.master_package_id, test_details.test_name, test_details.test_type").
		Scan(&testDetailsForLisEvent).Error; err != nil {
		return testDetailsForLisEvent, commonUtils.HandleORMError(err)
	}

	return testDetailsForLisEvent, nil
}

func (sampleDao *SampleDao) GetOmsTestDetailsByVisitId(visitId string) (
	[]commonStructures.OmsTestDetailsForLis, *commonStructures.CommonError) {

	omsTestDetailsForLis := []commonStructures.OmsTestDetailsForLis{}

	selectStrings := []string{
		"test_details.oms_order_id as order_id",
		"test_details.central_oms_test_id as test_id",
		"test_details.lis_code as test_code",
		"test_details.master_test_id as master_test_id",
		"test_details.master_package_id as master_package_id",
	}

	sampleStatuses := []string{commonConstants.SamplePartiallyRejected, commonConstants.SampleSynced,
		commonConstants.SampleReceived, commonConstants.SampleAccessioned}
	if err := sampleDao.Db.Table(commonConstants.TableTestDetails).
		Select(selectStrings).
		Joins("INNER JOIN test_sample_mapping ON test_details.oms_order_id = test_sample_mapping.oms_order_id and test_sample_mapping.oms_test_id = test_details.central_oms_test_id").
		Joins("INNER JOIN samples ON samples.oms_order_id = test_sample_mapping.oms_order_id and samples.sample_number = test_sample_mapping.sample_number").
		Where("samples.visit_id = ?", visitId).
		Where("samples.status IN (?)", sampleStatuses).
		Where("test_sample_mapping.is_rejected = false").
		Where("test_details.deleted_at IS NULL").
		Where("test_sample_mapping.deleted_at IS NULL").
		Where("samples.deleted_at IS NULL").
		// Group("test_details.oms_order_id, test_details.central_oms_test_id, test_details.lis_code, test_details.master_test_id, test_details.master_package_id").
		Scan(&omsTestDetailsForLis).Error; err != nil {
		return omsTestDetailsForLis, commonUtils.HandleORMError(err)
	}

	return omsTestDetailsForLis, nil
}

func (sampleDao *SampleDao) GetVisitIdsByOmsOrderId(omsOrderId string) ([]string, *commonStructures.CommonError) {
	var visitIds []string
	if err := sampleDao.Db.Table(commonConstants.TableSamples).
		Select("visit_id").
		Where("oms_order_id = ?", omsOrderId).
		Pluck("visit_id", &visitIds).Error; err != nil {
		return visitIds, commonUtils.HandleORMError(err)
	}

	return commonUtils.RemoveEmptyStringsFromStringSlice(commonUtils.CreateUniqueSliceString(visitIds)), nil
}

func (sampleDao *SampleDao) GetVisitLabMapByOmsTestIds(omsTestIds []string) (
	map[string]uint, *commonStructures.CommonError) {
	samples := []commonModels.Sample{}
	if err := sampleDao.Db.Model(&commonModels.Sample{}).
		Joins("INNER JOIN test_sample_mapping ON samples.oms_order_id = test_sample_mapping.oms_order_id AND samples.sample_number = test_sample_mapping.sample_number").
		Where("test_sample_mapping.oms_test_id IN (?)", omsTestIds).
		Where("test_sample_mapping.is_rejected = false").
		Find(&samples).Error; err != nil {
		return nil, commonUtils.HandleORMError(err)
	}
	visitIdLabMap := make(map[string]uint)
	for _, sample := range samples {
		if sample.VisitId != "" {
			visitIdLabMap[sample.VisitId] = sample.LabId
		}
	}

	return visitIdLabMap, nil
}

func (sampleDao *SampleDao) GetSampleDataBySampleNumberAndTestId(sampleNumber uint, testId string) ([]commonModels.Sample,
	[]commonModels.SampleMetadata, *commonStructures.CommonError) {
	samples, sampleMetadatas := []commonModels.Sample{}, []commonModels.SampleMetadata{}
	if err := sampleDao.Db.Table(commonConstants.TableSamples).
		Joins("INNER JOIN test_sample_mapping ON samples.oms_order_id = test_sample_mapping.oms_order_id AND samples.sample_number = test_sample_mapping.sample_number").
		Where("test_sample_mapping.sample_number = ? AND test_sample_mapping.oms_test_id = ?", sampleNumber, testId).
		Find(&samples).Error; err != nil {
		return samples, sampleMetadatas, commonUtils.HandleORMError(err)
	}

	sampleIds := []uint{}
	for _, sample := range samples {
		sampleIds = append(sampleIds, sample.Id)
	}

	if err := sampleDao.Db.Table(commonConstants.TableSampleMetadata).
		Where("sample_id IN (?)", sampleIds).
		Find(&sampleMetadatas).Error; err != nil {
		return samples, sampleMetadatas, commonUtils.HandleORMError(err)
	}

	return samples, sampleMetadatas, nil
}

func (sampleDao *SampleDao) GetSampleByVisitId(visitId string) (commonModels.Sample, *commonStructures.CommonError) {
	sample := commonModels.Sample{}
	if err := sampleDao.Db.Table("samples").
		Where("visit_id = ?", visitId).
		Where("deleted_at IS NULL").
		First(&sample).Error; err != nil {
		return sample, commonUtils.HandleORMError(err)
	}

	return sample, nil
}

func (sampleDao *SampleDao) RemoveSamplesNotLinkedToAnyTests(omsOrderId string) *commonStructures.CommonError {
	sampleIdsFromSamples := []uint{}
	if err := sampleDao.Db.Table("samples").
		Where("parent_sample_id = 0").
		Where("oms_order_id = ?", omsOrderId).
		Where("deleted_at IS NULL").
		Pluck("id", &sampleIdsFromSamples).Error; err != nil {
		return commonUtils.HandleORMError(err)
	}

	sampleIdsFromTsms := []uint{}
	if err := sampleDao.Db.Table("test_sample_mapping").
		Where("oms_order_id = ?", omsOrderId).
		Where("deleted_at IS NULL").
		Pluck("sample_id", &sampleIdsFromTsms).Error; err != nil {
		return commonUtils.HandleORMError(err)
	}
	sampleIdsFromSamples = commonUtils.CreateUniqueSliceUint(sampleIdsFromSamples)
	sampleIdsFromTsms = commonUtils.CreateUniqueSliceUint(sampleIdsFromTsms)

	sampleIds := commonUtils.GetDifferenceBetweenUintSlices(sampleIdsFromSamples, sampleIdsFromTsms)

	commonUtils.AddLog(context.Background(), commonConstants.DEBUG_LEVEL, commonUtils.GetCurrentFunctionName(), map[string]interface{}{
		"sample_ids":              sampleIds,
		"oms_order_id":            omsOrderId,
		"sample_ids_from_samples": sampleIdsFromSamples,
		"sample_ids_from_tsms":    sampleIdsFromTsms,
	}, nil)

	if len(sampleIds) == 0 {
		return nil
	}

	txErr := sampleDao.Db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&commonModels.Sample{}).
			Where("id IN (?)", sampleIds).
			Updates(map[string]interface{}{
				"status":     commonConstants.SampleDeleted,
				"deleted_at": commonUtils.GetCurrentTime(),
				"deleted_by": commonConstants.CitadelSystemId,
				"updated_by": commonConstants.CitadelSystemId,
			}).Error; err != nil {
			return err
		}
		if err := tx.Model(&commonModels.SampleMetadata{}).
			Where("sample_id IN (?)", sampleIds).
			Updates(map[string]interface{}{
				"deleted_at": commonUtils.GetCurrentTime(),
				"deleted_by": commonConstants.CitadelSystemId,
				"updated_by": commonConstants.CitadelSystemId,
			}).Error; err != nil {
			return err
		}
		return nil
	})
	if txErr != nil {
		return commonUtils.HandleORMError(txErr)
	}

	return nil
}

func (sampleDao *SampleDao) BarcodesExistsInSystem(barcodes []string) (bool, *commonStructures.CommonError) {
	if commonUtils.StringSliceHasDuplicates(barcodes) {
		return true, nil
	}
	samples := []commonModels.Sample{}
	if err := sampleDao.Db.Where("barcode IN (?)", barcodes).Find(&samples).Error; err != nil {
		return false, commonUtils.HandleORMError(err)
	}

	return len(samples) != 0, nil
}

func (sampleDao *SampleDao) GetVisitDetailsForTaskByOmsOrderId(omsOrderId string) (
	[]commonStructures.VisitDetailsForTask, *commonStructures.CommonError) {
	visitDetailsForTask := []commonStructures.VisitDetailsForTask{}
	selectStrings := []string{
		"samples.visit_id as visit_id",
		"sample_metadata.collected_at as sample_collected_at",
		"sample_metadata.received_at as sample_received_at",
	}
	if err := sampleDao.Db.Table("samples").
		Joins("INNER JOIN sample_metadata ON sample_metadata.sample_id = samples.id").
		Select(selectStrings).
		Where("samples.oms_order_id = ?", omsOrderId).
		Where("sample_metadata.deleted_at IS NULL").
		Find(&visitDetailsForTask).Error; err != nil {
		return nil, commonUtils.HandleORMError(err)
	}

	return visitDetailsForTask, nil
}

func (sampleDao *SampleDao) IsSampleCollected(omsOrderId string) (bool, *commonStructures.CommonError) {
	sampleStatuses := []string{
		commonConstants.SampleRejected,
		commonConstants.SamplePartiallyRejected,
		commonConstants.SampleSynced,
		commonConstants.SampleReceived,
		commonConstants.SampleCollectionDone,
		commonConstants.SampleAccessioned,
		commonConstants.SampleTransferred,
		commonConstants.SampleTransferFailed,
		commonConstants.SampleInTransfer,
		commonConstants.SampleOutsourced,
	}

	var count int64
	if err := sampleDao.Db.Table("samples").
		Where("oms_order_id = ? AND status IN (?)", omsOrderId, sampleStatuses).
		Count(&count).Error; err != nil {
		return false, commonUtils.HandleORMError(err)
	}
	return count != 0, nil
}

func (sampleDao *SampleDao) UpdateSamplesAndSamplesMetadata(samples []commonModels.Sample,
	samplesMetadata []commonModels.SampleMetadata) ([]commonModels.Sample, []commonModels.SampleMetadata,
	*commonStructures.CommonError) {

	txErr := sampleDao.Db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(&samples).Error; err != nil {
			return err
		}

		if err := tx.Save(&samplesMetadata).Error; err != nil {
			return err
		}

		return nil
	})

	if txErr != nil {
		return samples, samplesMetadata, commonUtils.HandleORMError(txErr)
	}

	return samples, samplesMetadata, nil
}

func (sampleDao *SampleDao) UpdateSampleAndSampleMetadataWithTx(tx *gorm.DB, sample commonModels.Sample, sampleMetadata commonModels.SampleMetadata) (
	commonModels.Sample, commonModels.SampleMetadata, *commonStructures.CommonError) {
	err := tx.Save(&sample).Error
	if err != nil {
		return sample, sampleMetadata, commonUtils.HandleORMError(err)
	}

	err = tx.Save(&sampleMetadata).Error
	if err != nil {
		return sample, sampleMetadata, commonUtils.HandleORMError(err)
	}

	return sample, sampleMetadata, nil
}

func (sampleDao *SampleDao) UpdateSamplesAndSamplesMetadataWithTx(tx *gorm.DB, samples []commonModels.Sample,
	samplesMetadata []commonModels.SampleMetadata) ([]commonModels.Sample, []commonModels.SampleMetadata,
	*commonStructures.CommonError) {

	if err := tx.Save(&samples).Error; err != nil {
		return samples, samplesMetadata, commonUtils.HandleORMError(err)
	}

	if err := tx.Save(&samplesMetadata).Error; err != nil {
		return samples, samplesMetadata, commonUtils.HandleORMError(err)
	}

	return samples, samplesMetadata, nil
}

func (sampleDao *SampleDao) GetCollectedSamples(omsOrderId string, labId uint) ([]commonModels.Sample,
	[]commonModels.SampleMetadata, *commonStructures.CommonError) {
	collectedSamples, collectedSampleMetadatas := []commonModels.Sample{}, []commonModels.SampleMetadata{}
	collectedSampleStatuses := []string{
		commonConstants.SampleCollectionDone,
		commonConstants.SampleNotCollectedEmedic,
		commonConstants.SampleInTransfer,
		commonConstants.SampleReceived,
	}

	if err := sampleDao.Db.Table(commonConstants.TableSamples).
		Where("status IN (?)", collectedSampleStatuses).
		Where("oms_order_id = ? AND destination_lab_id = ?", omsOrderId, labId).
		Find(&collectedSamples).Error; err != nil {
		return collectedSamples, collectedSampleMetadatas, commonUtils.HandleORMError(err)
	}

	if len(collectedSamples) == 0 {
		return collectedSamples, collectedSampleMetadatas, nil
	}

	sampleIds := []uint{}
	for _, sample := range collectedSamples {
		sampleIds = append(sampleIds, sample.Id)
	}

	if err := sampleDao.Db.Table(commonConstants.TableSampleMetadata).
		Where("sample_id IN (?)", sampleIds).
		Where("deleted_at IS NULL").
		Find(&collectedSampleMetadatas).Error; err != nil {
		return collectedSamples, collectedSampleMetadatas, commonUtils.HandleORMError(err)
	}

	return collectedSamples, collectedSampleMetadatas, nil
}

func (sampleDao *SampleDao) GetSampleByBarcodeForReceiving(barcode string) (commonModels.Sample, *commonStructures.CommonError) {
	sample := commonModels.Sample{}
	collectedSampleStatuses := []string{
		commonConstants.SampleCollectionDone,
		commonConstants.SampleNotCollectedEmedic,
		commonConstants.SampleInTransfer,
		commonConstants.SampleReceived,
		commonConstants.SampleSynced,
	}
	if err := sampleDao.Db.Where("barcode = ?", barcode).
		Where("status IN (?)", collectedSampleStatuses).Find(&sample).Error; err != nil {
		return sample, commonUtils.HandleORMError(err)
	}

	return sample, nil
}

func (sampleDao *SampleDao) GetTestDetailsBySampleIds(sampleIds []uint) ([]commonModels.TestDetail,
	*commonStructures.CommonError) {
	testDetails := []commonModels.TestDetail{}
	if err := sampleDao.Db.Table(commonConstants.TableTestDetails).
		Joins("INNER JOIN test_sample_mapping ON test_details.central_oms_test_id = test_sample_mapping.oms_test_id AND test_details.oms_order_id = test_sample_mapping.oms_order_id").
		Where("test_sample_mapping.sample_id IN (?)", sampleIds).
		Where("test_sample_mapping.deleted_at IS NULL").
		Where("test_details.deleted_at IS NULL").
		Find(&testDetails).Error; err != nil {
		return testDetails, commonUtils.HandleORMError(err)
	}

	// Create unique test details
	uniqueTestDetailsMap := map[uint]commonModels.TestDetail{}
	for _, testDetail := range testDetails {
		uniqueTestDetailsMap[testDetail.Id] = testDetail
	}
	testDetails = []commonModels.TestDetail{}
	for _, testDetail := range uniqueTestDetailsMap {
		testDetails = append(testDetails, testDetail)
	}

	return testDetails, nil
}

func (sampleDao *SampleDao) GetAllTestsAndSampleMappingsBySampleNumbers(sampleNumbers []uint, omsOrderId string) (
	[]commonModels.TestDetail, []commonModels.TestSampleMapping, *commonStructures.CommonError) {
	testDetails, testSampleMappings := []commonModels.TestDetail{}, []commonModels.TestSampleMapping{}

	if err := sampleDao.Db.Table(commonConstants.TableTestSampleMapping).
		Where("sample_number IN (?)", sampleNumbers).
		Where("oms_order_id = ?", omsOrderId).
		Where("is_rejected = false").
		Where("deleted_at is null").
		Find(&testSampleMappings).Error; err != nil {
		return testDetails, testSampleMappings, commonUtils.HandleORMError(err)
	}

	tsmIds := []uint{}
	for _, tsm := range testSampleMappings {
		tsmIds = append(tsmIds, tsm.Id)
	}

	if err := sampleDao.Db.Table(commonConstants.TableTestDetails).
		Joins("INNER JOIN test_sample_mapping ON test_details.central_oms_test_id = test_sample_mapping.oms_test_id").
		Where("test_sample_mapping.sample_number IN (?)", sampleNumbers).
		Where("test_sample_mapping.id IN (?)", tsmIds).
		Find(&testDetails).Error; err != nil {
		return testDetails, testSampleMappings, commonUtils.HandleORMError(err)
	}

	return testDetails, testSampleMappings, nil
}

func (sampleDao *SampleDao) GetSampleDataBySampleId(sampleId uint) (commonModels.Sample, commonModels.SampleMetadata,
	*commonStructures.CommonError) {
	sample, sampleMetadata := commonModels.Sample{}, commonModels.SampleMetadata{}

	if err := sampleDao.Db.Where("id = ?", sampleId).Find(&sample).Error; err != nil {
		return sample, sampleMetadata, commonUtils.HandleORMError(err)
	}

	if err := sampleDao.Db.Where("sample_id = ?", sampleId).Find(&sampleMetadata).Error; err != nil {
		return sample, sampleMetadata, commonUtils.HandleORMError(err)
	}

	return sample, sampleMetadata, nil
}

func (sampleDao *SampleDao) GetSampleDataByBarcodeForRejection(barcode string) ([]commonModels.Sample, []commonModels.SampleMetadata,
	*commonStructures.CommonError) {
	samples, sampleMetadatas := []commonModels.Sample{}, []commonModels.SampleMetadata{}

	if err := sampleDao.Db.Where("barcode = ?", barcode).Find(&samples).Error; err != nil {
		return samples, sampleMetadatas, commonUtils.HandleORMError(err)
	}

	sampleIds := []uint{}
	for _, sample := range samples {
		sampleIds = append(sampleIds, sample.Id)
	}

	if len(sampleIds) == 0 {
		return samples, sampleMetadatas, &commonStructures.CommonError{
			StatusCode: http.StatusNotFound,
			Message:    commonConstants.ERROR_SAMPLE_NOT_FOUND,
		}
	}

	if err := sampleDao.Db.Where("sample_id IN (?)", sampleIds).Find(&sampleMetadatas).Error; err != nil {
		return samples, sampleMetadatas, commonUtils.HandleORMError(err)
	}

	return samples, sampleMetadatas, nil
}

func (sampleDao *SampleDao) GetSamplesDataBySampleIds(sampleIds []uint) ([]commonModels.Sample, []commonModels.SampleMetadata,
	*commonStructures.CommonError) {
	samples, sampleMetadata := []commonModels.Sample{}, []commonModels.SampleMetadata{}

	if err := sampleDao.Db.Where("id IN (?)", sampleIds).Find(&samples).Error; err != nil {
		return samples, sampleMetadata, commonUtils.HandleORMError(err)
	}

	if err := sampleDao.Db.Where("sample_id IN (?)", sampleIds).Find(&sampleMetadata).Error; err != nil {
		return samples, sampleMetadata, commonUtils.HandleORMError(err)
	}

	return samples, sampleMetadata, nil
}

func (sampleDao *SampleDao) DeleteAllSamplesDataByOmsOrderIdWithTx(tx *gorm.DB, omsOrderId string) *commonStructures.CommonError {
	currentTime := commonUtils.GetCurrentTime()
	updates := map[string]interface{}{
		"deleted_at": currentTime,
		"updated_at": currentTime,
		"deleted_by": commonConstants.CitadelSystemId,
		"updated_by": commonConstants.CitadelSystemId,
	}

	if err := tx.Model(&commonModels.SampleMetadata{}).
		Where("oms_order_id = ?", omsOrderId).
		Updates(updates).Error; err != nil {
		return commonUtils.HandleORMError(err)
	}

	if err := tx.Model(&commonModels.TestSampleMapping{}).
		Where("oms_order_id = ?", omsOrderId).
		Updates(updates).Error; err != nil {
		return commonUtils.HandleORMError(err)
	}

	updates["status"] = commonConstants.SampleDeleted
	if err := tx.Model(&commonModels.Sample{}).
		Where("oms_order_id = ?", omsOrderId).
		Updates(updates).Error; err != nil {
		return commonUtils.HandleORMError(err)
	}

	return nil
}

func (sampleDao *SampleDao) MarkSampleAsEmedicNotCollected(omsRequestId string) *commonStructures.CommonError {
	samples := []commonModels.Sample{}

	err := sampleDao.Db.Table("samples").
		Joins("INNER JOIN sample_metadata ON samples.id = sample_metadata.sample_id").
		Where("samples.oms_request_id = ?", omsRequestId).
		Where("sample_metadata.collect_later_reason != ''").
		Where("samples.status = ?", commonConstants.SampleDefault).Find(&samples).Error
	if err != nil {
		return commonUtils.HandleORMError(err)
	}
	if len(samples) == 0 {
		return &commonStructures.CommonError{
			StatusCode: http.StatusNotFound,
			Message:    commonConstants.ERROR_SAMPLE_NOT_FOUND,
		}
	}

	for index := range samples {
		samples[index].Status = commonConstants.SampleNotCollectedEmedic
		samples[index].UpdatedBy = commonConstants.CitadelSystemId
	}

	err = sampleDao.Db.Save(&samples).Error
	if err != nil {
		return commonUtils.HandleORMError(err)
	}

	return nil
}

func (sampleDao *SampleDao) ReMarkSampleDefaultEmedicNotCollected(omsRequestId string,
	taskSequence uint) *commonStructures.CommonError {
	samples, sampleMetadatas := []commonModels.Sample{}, []commonModels.SampleMetadata{}

	if err := sampleDao.Db.Table("samples").
		Select("sample_metadata.*, samples.*").
		Joins("INNER JOIN sample_metadata ON samples.id = sample_metadata.sample_id").
		Where("samples.oms_request_id = ?", omsRequestId).
		Where("sample_metadata.collect_later_reason = ?", commonConstants.NotCollectedReasonCollectLater).
		Where("samples.status = ?", commonConstants.SampleNotCollectedEmedic).
		Find(&samples).Error; err != nil {
		return commonUtils.HandleORMError(err)
	}

	if len(samples) == 0 {
		return &commonStructures.CommonError{
			StatusCode: http.StatusNotFound,
			Message:    commonConstants.ERROR_SAMPLE_NOT_FOUND,
		}
	}

	sampleIds := []uint{}
	for _, sample := range samples {
		sampleIds = append(sampleIds, sample.Id)
	}

	if er := sampleDao.Db.Table("sample_metadata").
		Where("sample_id IN (?)", sampleIds).
		Find(&sampleMetadatas).Error; er != nil {
		return commonUtils.HandleORMError(er)
	}

	sampleUpdates := map[string]interface{}{
		"status":     commonConstants.SampleDefault,
		"updated_by": commonConstants.CitadelSystemId,
	}

	sampleMetadataUpdates := map[string]interface{}{
		"task_sequence":        taskSequence,
		"collect_later_reason": "",
	}
	txErr := sampleDao.Db.Transaction(func(tx *gorm.DB) error {
		err := tx.Table("samples").Where("id IN (?)", sampleIds).Updates(sampleUpdates).Error
		if err != nil {
			return err
		}

		err = tx.Table("sample_metadata").Where("sample_id IN (?)", sampleIds).Updates(sampleMetadataUpdates).Error
		if err != nil {
			return err
		}

		return nil
	})

	if txErr != nil {
		return commonUtils.HandleORMError(txErr)
	}

	return nil
}

func (sampleDao *SampleDao) CollectionPortalMarkAccessionAsAccessioned(isWebhook bool, sampleId uint,
	omsRequestId string) *commonStructures.CommonError {
	if isWebhook {
		sample := commonModels.Sample{}
		if err := sampleDao.Db.Model(&sample).
			Joins("INNER JOIN sample_metadata ON samples.id = sample_metadata.sample_id").
			Where("samples.oms_request_id = ?", omsRequestId).
			Where("sample_metadata.accessioned_at IS NOT NULL").
			Where("sample_metadata.rejected_at IS NULL").
			Where("samples.deleted_at IS NULL").
			Where("sample_metadata.deleted_at IS NULL").
			Find(&sample).Error; err != nil {
			return commonUtils.HandleORMError(err)
		}
		if sample.Id == 0 {
			return &commonStructures.CommonError{
				StatusCode: http.StatusNotFound,
				Message:    commonConstants.ERROR_SAMPLE_NOT_FOUND,
			}
		}
		sample.Status = commonConstants.SampleAccessioned
		sample.UpdatedBy = commonConstants.CitadelSystemId
		err := sampleDao.Db.Save(&sample).Error
		if err != nil {
			return commonUtils.HandleORMError(err)
		}
	} else {
		currentTime := commonUtils.GetCurrentTime()
		fieldsToBeUpdated := map[string]interface{}{
			"barcode_scanned_at": currentTime,
			"collected_at":       currentTime,
			"received_at":        currentTime,
			"accessioned_at":     currentTime,
		}
		err := sampleDao.Db.Model(&commonModels.SampleMetadata{}).Where("sample_id = ?", sampleId).
			Updates(fieldsToBeUpdated).Error
		if err != nil {
			return commonUtils.HandleORMError(err)
		}
	}

	return nil
}

func (sampleDao *SampleDao) AddCollectedVolumeToSample(sampleId, volume uint) *commonStructures.CommonError {
	updates := map[string]interface{}{
		"collected_volume": volume,
		"updated_by":       commonConstants.CitadelSystemId,
	}

	if err := sampleDao.Db.Model(&commonModels.SampleMetadata{}).Where("sample_id = ?", sampleId).
		Updates(updates).Error; err != nil {
		return commonUtils.HandleORMError(err)
	}

	return nil
}

func (sampleDao *SampleDao) UpdateTaskSequenceForSample(omsRequestId string, taskId uint,
	testIds []string) *commonStructures.CommonError {

	sampleIds := []uint{}
	if err := sampleDao.Db.Model(&commonModels.Sample{}).
		Joins("INNER JOIN test_sample_mapping ON samples.oms_order_id = test_sample_mapping.oms_order_id AND samples.sample_number = test_sample_mapping.sample_number").
		Where("samples.oms_request_id = ?", omsRequestId).
		Where("test_sample_mapping.oms_test_id IN (?)", testIds).
		Where("samples.deleted_at IS NULL AND test_sample_mapping.deleted_at IS NULL").
		Pluck("samples.id", &sampleIds).Error; err != nil {
		return commonUtils.HandleORMError(err)
	}

	if len(sampleIds) == 0 {
		return &commonStructures.CommonError{
			StatusCode: http.StatusNotFound,
			Message:    commonConstants.ERROR_SAMPLE_NOT_FOUND,
		}
	}

	updates := map[string]interface{}{
		"task_sequence": taskId,
		"updated_by":    commonConstants.CitadelSystemId,
	}

	if err := sampleDao.Db.Model(&commonModels.SampleMetadata{}).
		Where("sample_id IN (?)", sampleIds).
		Updates(updates).Error; err != nil {
		return commonUtils.HandleORMError(err)
	}

	return nil
}

func (sampleDao *SampleDao) RemapSamplesToNewTaskSequence(orderIdToSampleUpdateMap map[string][]uint,
	newSequence uint) *commonStructures.CommonError {

	for omsOrderId, sampleNumbers := range orderIdToSampleUpdateMap {
		t := time.Now()
		samples := []commonModels.Sample{}
		if err := sampleDao.Db.Where("oms_order_id = ? AND sample_number IN (?)", omsOrderId, sampleNumbers).
			Find(&samples).Error; err != nil {
			return commonUtils.HandleORMError(err)
		}
		if len(samples) == 0 {
			continue
		}
		sampleIds := []uint{}
		for _, sample := range samples {
			sampleIds = append(sampleIds, sample.Id)
		}
		if err := sampleDao.Db.Model(&commonModels.SampleMetadata{}).
			Where("sample_id IN (?)", sampleIds).
			Updates(map[string]interface{}{
				"task_sequence": newSequence,
				"updated_at":    t,
			}).Error; err != nil {
			return commonUtils.HandleORMError(err)
		}
	}

	return nil
}

func (sampleDao *SampleDao) UpdateSampleDetailsForReschedule(
	requestBody structures.UpdateSampleDetailsForRescheduleRequest) *commonStructures.CommonError {
	sampleIds := []uint{}
	sampleQuery := sampleDao.Db.Model(&commonModels.Sample{}).
		Joins("INNER JOIN sample_metadata ON samples.id = sample_metadata.sample_id").
		Where("samples.barcode = ''")
	if requestBody.OldTaskId == 0 {
		sampleQuery = sampleQuery.Where("samples.oms_request_id = ?", requestBody.RequestId).
			Where("samples.status = ?", commonConstants.SampleCollectionDone)
		if requestBody.TaskType == commonConstants.OmsTaskTypePrimaryCollection {
			sampleQuery = sampleQuery.Where("sample_metadata.task_sequence = ?", 0)
		} else {
			sampleQuery = sampleQuery.Where("sample_metadata.task_sequence = ?", requestBody.TaskId)
		}
	} else {
		sampleQuery = sampleQuery.Where("samples.oms_request_id = ?", requestBody.RequestId)
		if requestBody.TaskType == commonConstants.OmsTaskTypePrimaryCollection {
			sampleQuery = sampleQuery.Where("sample_metadata.task_sequence = 0")
		} else {
			sampleQuery = sampleQuery.Where("sample_metadata.task_sequence = ?", requestBody.OldTaskId)
		}
	}
	sampleQuery = sampleQuery.Where("samples.deleted_at IS NULL").
		Pluck("samples.id", &sampleIds)
	if err := sampleQuery.Error; err != nil {
		return commonUtils.HandleORMError(err)
	}

	if len(sampleIds) == 0 {
		return nil
	}

	txErr := sampleDao.Db.Transaction(func(tx *gorm.DB) error {
		sampleUpdates := map[string]interface{}{
			"status":     commonConstants.SampleDefault,
			"updated_by": commonConstants.CitadelSystemId,
		}
		if err := tx.Model(&commonModels.Sample{}).
			Where("id IN (?)", sampleIds).
			Updates(sampleUpdates).Error; err != nil {
			return err
		}
		if requestBody.OldTaskId != 0 {
			sampleMetadataUpdates := map[string]interface{}{
				"updated_by": commonConstants.CitadelSystemId,
			}
			if requestBody.TaskType == commonConstants.OmsTaskTypePrimaryCollection {
				sampleMetadataUpdates["task_sequence"] = 0
			} else {
				sampleMetadataUpdates["task_sequence"] = requestBody.NewTaskId
			}
			if err := tx.Model(&commonModels.SampleMetadata{}).
				Where("sample_id IN (?)", sampleIds).
				Updates(sampleMetadataUpdates).Error; err != nil {
				return err
			}
		}
		return nil
	})
	if txErr != nil {
		return commonUtils.HandleORMError(txErr)
	}

	return nil
}

func (sampleDao *SampleDao) GetSamplesForDelayedReverseLogistics(normalTat, campTat, inclinicTat, days uint) (
	[]structures.DelayedReverseLogisticsSamplesDbStruct, *commonStructures.CommonError) {
	results := []structures.DelayedReverseLogisticsSamplesDbStruct{}

	query := sampleDao.Db.Table(commonConstants.TableSamples).
		Select(getStringsForReverseLogisiticsSamples()).
		Joins("INNER JOIN sample_metadata ON sample_metadata.sample_id = samples.id").
		Joins("INNER JOIN order_details ON order_details.oms_order_id = samples.oms_order_id").
		Joins("INNER JOIN patient_details ON patient_details.id = order_details.patient_details_id").
		Where("samples.deleted_at IS NULL").
		Where("samples.status = ?", commonConstants.SampleCollectionDone).
		Where("order_details.order_status = ?", commonConstants.OrderRequested).
		Where("samples.parent_sample_id = 0").
		Where("samples.vial_type_id NOT IN (?)", commonConstants.VialTypesToBeSkippedForReverseLogistics)

	queryStr := `(
		(
			order_details.collection_type <> 2
			AND COALESCE(order_details.camp_id, 0) = 0
			AND (sample_metadata.collected_at + interval '210 minute') < now ()
			AND (sample_metadata.collected_at + INTERVAL '7 day') > now ()
		)
		OR (
			COALESCE(order_details.camp_id, 0) > 0
			AND (sample_metadata.collected_at + INTERVAL '320 minute') < now ()
			AND (sample_metadata.collected_at + INTERVAL '7 day') > now ()
		)
		OR (
			order_details.collection_type = 2
			AND (sample_metadata.collected_at + INTERVAL '360 minute') < now ()
			AND (sample_metadata.collected_at + INTERVAL '7 day') > now ()
		)
	)`
	query = query.Where(queryStr)

	if err := query.Scan(&results).Error; err != nil {
		return results, commonUtils.HandleORMError(err)
	}

	return results, nil
}

func (sampleDao *SampleDao) GetSamplesForDelayedInterlabLogistics() (
	[]commonStructures.SampleInfo, *commonStructures.CommonError) {

	sampleInfos := []commonStructures.SampleInfo{}

	query := sampleDao.Db.Table(commonConstants.TableSamples).
		Select([]string{
			"samples.barcode",
			"samples.oms_order_id",
			"samples.oms_request_id",
			"sample_metadata.collected_at",
			"samples.status",
			"samples.oms_city_code",
			"samples.vial_type_id",
			"sample_metadata.transferred_at",
		}).
		Joins("INNER JOIN sample_metadata ON samples.id = sample_metadata.sample_id").
		Joins("INNER JOIN order_details ON order_details.oms_order_id = samples.oms_order_id").
		Where("samples.deleted_at IS NULL").
		Where("samples.parent_sample_id > 0").
		Where("samples.status = ?", commonConstants.SampleInTransfer).
		Where("order_details.order_status = ?", commonConstants.OrderRequested).
		Where("samples.vial_type_id NOT IN (?)", commonConstants.VialTypesToBeSkippedForReverseLogistics).
		Where("(sample_metadata.transferred_at + INTERVAL '1 day' + INTERVAL '750 minutes') < now()").
		Where("sample_metadata.transferred_at > (now() - INTERVAL '10 days')")

	if err := query.Scan(&sampleInfos).Error; err != nil {
		return sampleInfos, commonUtils.HandleORMError(err)
	}

	return sampleInfos, nil
}

func (sampleDao *SampleDao) GetSrfOrderIds(cityCode string) []string {
	orderDetails := []structures.SrfOrderIdsOrderDetails{}
	sampleDao.Db.Table(commonConstants.TableOrderDetails).
		Select("order_details.oms_order_id, order_details.city_code").
		Joins("INNER JOIN test_details ON test_details.oms_order_id = order_details.oms_order_id").
		Where("order_details.created_at >= NOW() - INTERVAL '15 days'").
		Where("order_details.order_status IN (?)", []string{commonConstants.OrderRequested, commonConstants.OrderCompleted}).
		Where("order_details.deleted_at IS NULL").
		Where("test_details.master_test_id IN (?)", commonUtils.GetCovid19MasterTestIds()).
		Where("test_details.deleted_at IS NULL").
		Find(&orderDetails)

	omsOrderIds := []string{}
	for _, orderDetail := range orderDetails {
		if cityCode == "" {
			omsOrderIds = append(omsOrderIds, orderDetail.OmsOrderId)
		} else if orderDetail.CityCode == cityCode {
			omsOrderIds = append(omsOrderIds, orderDetail.OmsOrderId)
		}
	}

	if len(omsOrderIds) == 0 {
		return []string{}
	}

	omsOrderIds = commonUtils.CreateUniqueSliceString(omsOrderIds)
	finalOrderIds := []string{}
	sampleDao.Db.Table(commonConstants.TableSamples).
		Joins("INNER JOIN sample_metadata ON samples.id = sample_metadata.sample_id").
		Select("DISTINCT samples.oms_order_id").
		Where("samples.oms_order_id IN (?)", omsOrderIds).
		Where("sample_metadata.collected_at < NOW() - INTERVAL '1 hour'").
		Pluck("samples.oms_order_id", &finalOrderIds)

	return commonUtils.CreateUniqueSliceString(finalOrderIds)
}
