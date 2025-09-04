package dao

import (
	"database/sql"
	"time"

	commonConstants "github.com/Orange-Health/citadel/common/constants"
	commonStructures "github.com/Orange-Health/citadel/common/structures"
	commonUtils "github.com/Orange-Health/citadel/common/utils"
	commonModels "github.com/Orange-Health/citadel/models"
)

func (attuneDao *AttuneDao) GetOrderDetailsByVisitId(visitId string) (commonModels.OrderDetails,
	*commonStructures.CommonError) {

	orderDetails := commonModels.OrderDetails{}
	if err := attuneDao.Db.Table(commonConstants.TableOrderDetails).
		Joins("INNER JOIN samples ON order_details.oms_order_id = samples.oms_order_id").
		Where("samples.visit_id = ?", visitId).
		Select("order_details.*").
		Find(&orderDetails).Error; err != nil {
		return orderDetails, commonUtils.HandleORMError(err)
	}

	return orderDetails, nil
}

func (attune *AttuneDao) GetSampleByVisitId(visitId string) (commonModels.Sample, *commonStructures.CommonError) {
	sample := commonModels.Sample{}
	if err := attune.Db.Table(commonConstants.TableSamples).
		Where("visit_id = ?", visitId).
		Select("samples.*").
		Find(&sample).Error; err != nil {
		return sample, commonUtils.HandleORMError(err)
	}

	return sample, nil
}

func (attuneDao *AttuneDao) GetOrderDetailsByOmsOrderId(omsOrderId string) (commonModels.OrderDetails,
	*commonStructures.CommonError) {
	orderDetails := commonModels.OrderDetails{}
	if err := attuneDao.Db.Table(commonConstants.TableOrderDetails).
		Find(&orderDetails, "oms_order_id = ?", omsOrderId).Error; err != nil {
		return orderDetails, commonUtils.HandleORMError(err)
	}

	return orderDetails, nil
}

func (attuneDao *AttuneDao) GetPatientDetailsById(patientDetailsId uint) (commonModels.PatientDetail,
	*commonStructures.CommonError) {

	patientDetails := commonModels.PatientDetail{}
	if err := attuneDao.Db.Table(commonConstants.TablePatientDetails).
		Find(&patientDetails, "id = ?", patientDetailsId).Error; err != nil {
		return patientDetails, commonUtils.HandleORMError(err)
	}

	return patientDetails, nil
}

func (attuneDao *AttuneDao) GetOrderDetailsAndPatientDetailsByVisitId(visitId string) (commonModels.OrderDetails,
	commonModels.PatientDetail, *commonStructures.CommonError) {

	orderDetails, cErr := attuneDao.GetOrderDetailsByVisitId(visitId)
	if cErr != nil {
		return orderDetails, commonModels.PatientDetail{}, cErr
	}

	patientDetails, cErr := attuneDao.GetPatientDetailsById(orderDetails.PatientDetailsId)
	if cErr != nil {
		return orderDetails, patientDetails, cErr
	}

	return orderDetails, patientDetails, nil
}

func (attuneDao *AttuneDao) GetSampleCollectedAtByVisitId(visitId string) (*time.Time, *commonStructures.CommonError) {
	sampleMetadata := commonModels.SampleMetadata{}
	if err := attuneDao.Db.Table(commonConstants.TableSampleMetadata).
		Joins("INNER JOIN samples ON samples.id = sample_metadata.sample_id").
		Select("sample_metadata.*").
		Where("samples.visit_id = ?", visitId).
		Find(&sampleMetadata).Error; err != nil {
		return nil, commonUtils.HandleORMError(err)
	}

	sampleCollectedAt := sampleMetadata.CollectedAt
	if sampleCollectedAt == nil {
		sampleCollectedAt = sampleMetadata.ReceivedAt
	}

	return sampleCollectedAt, nil
}

func (attuneDao *AttuneDao) GetSampleCollectedAtBySampleIds(sampleIds []uint) (*time.Time, *commonStructures.CommonError) {
	var sampleCollectedAt sql.NullTime

	if err := attuneDao.Db.Table(commonConstants.TableSampleMetadata).
		Select("collected_at").
		Where("sample_id IN (?)", sampleIds).
		Order("collected_at ASC").
		Limit(1).
		Scan(&sampleCollectedAt).Error; err != nil {
		return nil, commonUtils.HandleORMError(err)
	}

	if !sampleCollectedAt.Valid {
		return nil, nil // Return nil if collected_at is NULL
	}

	return &sampleCollectedAt.Time, nil // Return the actual time value
}

func (attuneDao *AttuneDao) GetAttuneTestSampleMapByOmsOrderId(omsOrderId string, sampleIds []uint) (
	[]commonStructures.AttuneTestSampleMapSnakeCase, *commonStructures.CommonError) {
	testSampleMap := []commonStructures.AttuneTestSampleMapSnakeCase{}
	if err := attuneDao.Db.Table(commonConstants.TableTestSampleMapping).
		Select("test_sample_mapping.oms_test_id as test_id, test_sample_mapping.sample_id, samples.barcode as barcode, test_sample_mapping.vial_type_id as vial_type_id").
		Joins("INNER JOIN samples ON samples.sample_number = test_sample_mapping.sample_number AND samples.oms_order_id = test_sample_mapping.oms_order_id").
		Where("samples.id IN (?)", sampleIds).
		Where("samples.oms_order_id = ?", omsOrderId).
		Scan(&testSampleMap).Error; err != nil {
		return testSampleMap, commonUtils.HandleORMError(err)
	}

	return testSampleMap, nil
}

func (attuneDao *AttuneDao) GetOrderDetailsAndPatientDetailsByOmsOrderId(omsOrderId string) (commonModels.OrderDetails,
	commonModels.PatientDetail, *commonStructures.CommonError) {

	orderDetails, cErr := attuneDao.GetOrderDetailsByOmsOrderId(omsOrderId)
	if cErr != nil {
		return orderDetails, commonModels.PatientDetail{}, cErr
	}

	patientDetails, cErr := attuneDao.GetPatientDetailsById(orderDetails.PatientDetailsId)
	if cErr != nil {
		return orderDetails, patientDetails, cErr
	}

	return orderDetails, patientDetails, nil
}

func (attuneDao *AttuneDao) GetTestDetailsForSyncingToAttune(omsOrderId string, sampleIds []uint, labId uint) (
	[]commonModels.TestDetail, *commonStructures.CommonError) {
	testDetails := []commonModels.TestDetail{}

	if err := attuneDao.Db.Table(commonConstants.TableTestDetails).
		Select("test_details.*").
		Joins("INNER JOIN test_sample_mapping ON test_details.central_oms_test_id = test_sample_mapping.oms_test_id AND test_details.oms_order_id = test_sample_mapping.oms_order_id").
		Joins("INNER JOIN samples ON samples.sample_number = test_sample_mapping.sample_number").
		Where("samples.oms_order_id = ?", omsOrderId).
		Where("test_sample_mapping.oms_order_id = ?", omsOrderId).
		Where("samples.id IN (?)", sampleIds).
		Where("test_details.processing_lab_id = ?", labId).
		Where("samples.deleted_at IS NULL").
		Where("test_sample_mapping.deleted_at IS NULL").
		Where("test_sample_mapping.is_rejected = false").
		Scan(&testDetails).Error; err != nil {
		return testDetails, commonUtils.HandleORMError(err)
	}

	return testDetails, nil
}
