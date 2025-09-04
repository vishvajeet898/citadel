package dao

import (
	"net/http"

	"gorm.io/gorm"

	"github.com/Orange-Health/citadel/apps/test_sample_mapping/mapper"
	commonConstants "github.com/Orange-Health/citadel/common/constants"
	commonStructures "github.com/Orange-Health/citadel/common/structures"
	commonUtils "github.com/Orange-Health/citadel/common/utils"
	commonModels "github.com/Orange-Health/citadel/models"
)

type DataLayer interface {
	GetTestSampleMappingsByOrderId(omsOrderId string) ([]commonStructures.TestSampleMappingInfo, *commonStructures.CommonError)
	GetTestSampleMappingsByOrderIds(omsOrderIds []string) ([]commonStructures.TestSampleMappingInfo,
		*commonStructures.CommonError)
	GetTestSampleMappingsByTestIds(omsTestIds []string) ([]commonModels.TestSampleMapping, *commonStructures.CommonError)
	GetTestSampleMappingsModelsByOrderId(omsOrderId string) ([]commonModels.TestSampleMapping, *commonStructures.CommonError)
	GetTestSampleMappingByOrderIdWithTx(tx *gorm.DB, omsOrderId string) ([]commonModels.TestSampleMapping,
		*commonStructures.CommonError)
	GetTestSampleMappingsByOrderIdAndTestIds(omsOrderId string, omsTestIds []string) (
		[]commonModels.TestSampleMapping, *commonStructures.CommonError)
	GetTestSampleMappingsByOrderIdAndTestIdsUnscoped(omsOrderId string, omsTestIds []string) (
		[]commonModels.TestSampleMapping, *commonStructures.CommonError)
	GetTsmByOrderIdTestIdAndSampleNumber(omsOrderId, omsTestId string, sampleNumber uint) (commonModels.TestSampleMapping,
		*commonStructures.CommonError)
	OrdersWithRecollectionsPendingPresent(omsOrderIds, omsTestIds []string) (bool, *commonStructures.CommonError)
	OrdersWithRecollectionsPendingPresentByOmsOrderId(omsOrderId string) (bool, *commonStructures.CommonError)
	CheckIfAllTestsRejectedByOrderIdTestIdAndSampleNumber(omsOrderId, omsTestId string, sampleNumber uint) (bool, error)

	CreateTestSampleMappingWithTx(*gorm.DB, commonModels.TestSampleMapping) (
		commonModels.TestSampleMapping, *commonStructures.CommonError)
	CreateBulkTestSampleMappingWithTx(*gorm.DB, []commonModels.TestSampleMapping) (
		[]commonModels.TestSampleMapping, *commonStructures.CommonError)
	UpdateTestSampleMappingWithTx(tx *gorm.DB, tsm commonModels.TestSampleMapping) *commonStructures.CommonError
	UpdateBulkTestSampleMapping([]commonModels.TestSampleMapping) (
		[]commonModels.TestSampleMapping, *commonStructures.CommonError)
	UpdateBulkTestSampleMappingWithTx(tx *gorm.DB, tsmInfos []commonModels.TestSampleMapping) (
		[]commonModels.TestSampleMapping, *commonStructures.CommonError)
	UpdateTsmForTestIdAndRecollectionPendingTrueWithTx(tx *gorm.DB, omsTestIds []string) *commonStructures.CommonError
	RejectTestSampleMappingBySampleIdWithTx(tx *gorm.DB, sampleId, userId uint) *commonStructures.CommonError
	DeleteTestSampleMappingByOrderIdTestIdAndSampleNumberWithTx(tx *gorm.DB, omsOrderId, omsTestId string,
		sampleNumber uint) *commonStructures.CommonError
	DeleteTestSampleMappingByOmsTestIdsWithTx(tx *gorm.DB, omsTestIds []string) *commonStructures.CommonError
	DeleteTestSampleMappingsWithTx(tx *gorm.DB, tsmInfos []commonModels.TestSampleMapping) *commonStructures.CommonError
}

func (tsmDao *TestSampleMappingDao) CreateTestSampleMappingWithTx(tx *gorm.DB, tsm commonModels.TestSampleMapping) (
	commonModels.TestSampleMapping, *commonStructures.CommonError) {
	if err := tx.Create(&tsm).Error; err != nil {
		return commonModels.TestSampleMapping{}, commonUtils.HandleORMError(err)
	}
	return tsm, nil
}

func (tsmDao *TestSampleMappingDao) CreateBulkTestSampleMappingWithTx(tx *gorm.DB,
	tsms []commonModels.TestSampleMapping) ([]commonModels.TestSampleMapping, *commonStructures.CommonError) {
	if err := tx.Create(&tsms).Error; err != nil {
		return []commonModels.TestSampleMapping{}, commonUtils.HandleORMError(err)
	}
	return tsms, nil
}

func (tsmDao *TestSampleMappingDao) UpdateTestSampleMappingWithTx(tx *gorm.DB,
	tsm commonModels.TestSampleMapping) *commonStructures.CommonError {
	if err := tx.Save(&tsm).Error; err != nil {
		return commonUtils.HandleORMError(err)
	}
	return nil
}

func (tsmDao *TestSampleMappingDao) CheckIfAllTestsRejectedByOrderIdTestIdAndSampleNumber(omsOrderId, omsTestId string,
	sampleNumber uint) (bool, error) {
	testCount := int64(0)
	query := tsmDao.Db.Model(&commonModels.TestSampleMapping{}).
		Where("sample_number = ? AND oms_test_id != ? AND oms_order_id = ? AND deleted_at IS NULL AND is_rejected = ?", sampleNumber, omsTestId, omsOrderId, false).
		Count(&testCount)
	return testCount == 0, query.Error
}

func (tsmDao *TestSampleMappingDao) UpdateBulkTestSampleMapping(tsms []commonModels.TestSampleMapping) (
	[]commonModels.TestSampleMapping, *commonStructures.CommonError) {
	if txErr := tsmDao.Db.Transaction(func(tx *gorm.DB) error {
		for _, tsm := range tsms {
			if err := tx.Model(&tsm).Updates(&tsm).Error; err != nil {
				return err
			}
		}
		return nil
	}); txErr != nil {
		return []commonModels.TestSampleMapping{}, commonUtils.HandleORMError(txErr)
	}
	return tsms, nil
}

func (tsmDao *TestSampleMappingDao) UpdateBulkTestSampleMappingWithTx(tx *gorm.DB, tsms []commonModels.TestSampleMapping) (
	[]commonModels.TestSampleMapping, *commonStructures.CommonError) {
	for _, tsm := range tsms {
		if err := tx.Model(&tsm).Save(&tsm).Error; err != nil {
			return []commonModels.TestSampleMapping{}, commonUtils.HandleORMError(err)
		}
	}
	return tsms, nil
}

func (tsmDao *TestSampleMappingDao) UpdateTsmForTestIdAndRecollectionPendingTrueWithTx(tx *gorm.DB,
	omsTestIds []string) *commonStructures.CommonError {
	tsmUpdates := map[string]interface{}{
		"recollection_pending": false,
		"updated_by":           commonConstants.CitadelSystemId,
	}

	if err := tx.Model(&commonModels.TestSampleMapping{}).
		Where("oms_test_id IN (?) AND recollection_pending = true", omsTestIds).
		Updates(tsmUpdates).Error; err != nil {
		return commonUtils.HandleORMError(err)
	}

	return nil
}

func (tsmDao *TestSampleMappingDao) RejectTestSampleMappingBySampleIdWithTx(tx *gorm.DB,
	sampleId, userId uint) *commonStructures.CommonError {
	updatesMap := map[string]interface{}{
		"is_rejected":          true,
		"recollection_pending": true,
		"updated_by":           userId,
		"updated_at":           commonUtils.GetCurrentTime(),
	}

	if err := tx.Model(&commonModels.TestSampleMapping{}).
		Where("sample_id = ?", sampleId).
		Updates(updatesMap).Error; err != nil {
		return commonUtils.HandleORMError(err)
	}

	return nil
}

func (tsmDao *TestSampleMappingDao) DeleteTestSampleMappingByOrderIdTestIdAndSampleNumberWithTx(tx *gorm.DB,
	omsOrderId, omsTestId string, sampleNumber uint) *commonStructures.CommonError {
	updates := map[string]interface{}{
		"deleted_at": commonUtils.GetCurrentTime(),
		"deleted_by": commonConstants.CitadelSystemId,
		"updated_by": commonConstants.CitadelSystemId,
	}
	if err := tx.Table(commonConstants.TableTestSampleMapping).Where(`oms_order_id = ?`, omsOrderId).
		Where(`oms_test_id = ?`, omsTestId).
		Where(`sample_number = ?`, sampleNumber).
		Updates(updates).Error; err != nil {
		return commonUtils.HandleORMError(err)
	}
	return nil
}

func (tsmDao *TestSampleMappingDao) DeleteTestSampleMappingByOmsTestIdsWithTx(tx *gorm.DB,
	omsTestIds []string) *commonStructures.CommonError {
	updates := map[string]interface{}{
		"deleted_at": commonUtils.GetCurrentTime(),
		"deleted_by": commonConstants.CitadelSystemId,
		"updated_by": commonConstants.CitadelSystemId,
	}

	if err := tx.Table(commonConstants.TableTestSampleMapping).
		Where("oms_test_id IN (?)", omsTestIds).
		Updates(updates).Error; err != nil {
		return commonUtils.HandleORMError(err)
	}

	return nil
}

func (tsmDao *TestSampleMappingDao) DeleteTestSampleMappingsWithTx(tx *gorm.DB,
	tsmInfos []commonModels.TestSampleMapping) *commonStructures.CommonError {
	updates := map[string]interface{}{
		"deleted_at": commonUtils.GetCurrentTime(),
		"deleted_by": commonConstants.CitadelSystemId,
		"updated_by": commonConstants.CitadelSystemId,
	}
	tsmIds := []uint{}
	for _, tsm := range tsmInfos {
		tsmIds = append(tsmIds, tsm.Id)
	}

	if err := tx.Table(commonConstants.TableTestSampleMapping).
		Where("id IN (?)", tsmIds).
		Updates(updates).Error; err != nil {
		return commonUtils.HandleORMError(err)
	}

	return nil

}

func (tsmDao *TestSampleMappingDao) GetTestSampleMappingsByTestIds(omsTestIds []string) (
	[]commonModels.TestSampleMapping, *commonStructures.CommonError) {
	var testSampleMaps []commonModels.TestSampleMapping
	if err := tsmDao.Db.Where("oms_test_id IN ?", omsTestIds).Find(&testSampleMaps).Error; err != nil {
		return nil, commonUtils.HandleORMError(err)
	}
	if len(testSampleMaps) == 0 {
		return nil, &commonStructures.CommonError{
			Message:    commonConstants.ERROR_TSM_NOT_FOUND,
			StatusCode: http.StatusNotFound,
		}
	}
	return testSampleMaps, nil
}

func (tsmDao *TestSampleMappingDao) GetTestSampleMappingsModelsByOrderId(omsOrderId string) (
	[]commonModels.TestSampleMapping, *commonStructures.CommonError) {
	var tsms []commonModels.TestSampleMapping
	if err := tsmDao.Db.Where("oms_order_id = ?", omsOrderId).Find(&tsms).Error; err != nil {
		return []commonModels.TestSampleMapping{}, commonUtils.HandleORMError(err)
	}
	if len(tsms) == 0 {
		return []commonModels.TestSampleMapping{}, &commonStructures.CommonError{
			Message:    commonConstants.ERROR_TSM_NOT_FOUND,
			StatusCode: http.StatusNotFound,
		}
	}
	return tsms, nil
}

func (tsmDao *TestSampleMappingDao) GetTestSampleMappingByOrderIdWithTx(tx *gorm.DB, omsOrderId string) (
	[]commonModels.TestSampleMapping, *commonStructures.CommonError) {
	var tsms []commonModels.TestSampleMapping
	if err := tx.Where("oms_order_id = ?", omsOrderId).Find(&tsms).Error; err != nil {
		return []commonModels.TestSampleMapping{}, commonUtils.HandleORMError(err)
	}
	return tsms, nil
}

func (tsmDao *TestSampleMappingDao) GetTestSampleMappingsByOrderIdAndTestIds(omsOrderId string, omsTestIds []string) (
	[]commonModels.TestSampleMapping, *commonStructures.CommonError) {
	var tsms []commonModels.TestSampleMapping
	if err := tsmDao.Db.Where("oms_order_id = ? AND oms_test_id in (?)", omsOrderId, omsTestIds).Find(&tsms).Error; err != nil {
		return []commonModels.TestSampleMapping{}, commonUtils.HandleORMError(err)
	}
	if len(tsms) == 0 {
		return []commonModels.TestSampleMapping{}, &commonStructures.CommonError{
			Message:    commonConstants.ERROR_TSM_NOT_FOUND,
			StatusCode: http.StatusNotFound,
		}
	}
	return tsms, nil
}

func (tsmDao *TestSampleMappingDao) GetTestSampleMappingsByOrderIdAndTestIdsUnscoped(omsOrderId string, omsTestIds []string) (
	[]commonModels.TestSampleMapping, *commonStructures.CommonError) {
	var tsms []commonModels.TestSampleMapping
	if err := tsmDao.Db.Unscoped().Where("oms_order_id = ? AND oms_test_id in (?)", omsOrderId, omsTestIds).Find(&tsms).Error; err != nil {
		return []commonModels.TestSampleMapping{}, commonUtils.HandleORMError(err)
	}
	if len(tsms) == 0 {
		return []commonModels.TestSampleMapping{}, &commonStructures.CommonError{
			Message:    commonConstants.ERROR_TSM_NOT_FOUND,
			StatusCode: http.StatusNotFound,
		}
	}
	return tsms, nil
}

func (tsmDao *TestSampleMappingDao) GetTsmByOrderIdTestIdAndSampleNumber(omsOrderId, omsTestId string, sampleNumber uint) (
	commonModels.TestSampleMapping, *commonStructures.CommonError) {
	var tsm commonModels.TestSampleMapping
	if err := tsmDao.Db.Where("oms_order_id = ?", omsOrderId).
		Where("oms_test_id = ?", omsTestId).
		Where("sample_number = ?", sampleNumber).
		First(&tsm).Error; err != nil {
		return commonModels.TestSampleMapping{}, commonUtils.HandleORMError(err)
	}
	return tsm, nil
}

func (tsmDao *TestSampleMappingDao) GetTestSampleMappingsByOrderId(omsOrderId string) (
	[]commonStructures.TestSampleMappingInfo, *commonStructures.CommonError) {
	var tsms []commonModels.TestSampleMapping
	if err := tsmDao.Db.Where("oms_order_id = ?", omsOrderId).Find(&tsms).Error; err != nil {
		return nil, commonUtils.HandleORMError(err)
	}
	if len(tsms) == 0 {
		return nil, &commonStructures.CommonError{
			Message:    commonConstants.ERROR_TSM_NOT_FOUND,
			StatusCode: http.StatusNotFound,
		}
	}
	return mapper.MapBulkTsmgModelToTsmInfo(tsms), nil
}

func (tsmDao *TestSampleMappingDao) GetTestSampleMappingsByOrderIds(omsOrderIds []string) (
	[]commonStructures.TestSampleMappingInfo, *commonStructures.CommonError) {
	var tsms []commonModels.TestSampleMapping
	if err := tsmDao.Db.Where("oms_order_id IN ?", omsOrderIds).Find(&tsms).Error; err != nil {
		return nil, commonUtils.HandleORMError(err)
	}
	if len(tsms) == 0 {
		return nil, &commonStructures.CommonError{
			Message:    commonConstants.ERROR_TSM_NOT_FOUND,
			StatusCode: http.StatusNotFound,
		}
	}
	return mapper.MapBulkTsmgModelToTsmInfo(tsms), nil
}

func (tsmDao *TestSampleMappingDao) OrdersWithRecollectionsPendingPresent(omsOrderIds, omsTestIds []string) (bool,
	*commonStructures.CommonError) {
	var count int64
	if err := tsmDao.Db.Model(&commonModels.TestSampleMapping{}).
		Where("oms_order_id IN (?) AND oms_test_id IN (?) AND recollection_pending = true", omsOrderIds, omsTestIds).
		Count(&count).Error; err != nil {
		return false, commonUtils.HandleORMError(err)
	}
	if count > 0 {
		return true, nil
	}
	return false, nil
}

func (tsmDao *TestSampleMappingDao) OrdersWithRecollectionsPendingPresentByOmsOrderId(omsOrderId string) (bool,
	*commonStructures.CommonError) {
	var count int64
	if err := tsmDao.Db.Model(&commonModels.TestSampleMapping{}).
		Where("oms_order_id = ? AND recollection_pending = true", omsOrderId).
		Count(&count).Error; err != nil {
		return false, commonUtils.HandleORMError(err)
	}
	if count > 0 {
		return true, nil
	}
	return false, nil
}
