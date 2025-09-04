package service

import (
	"net/http"

	"gorm.io/gorm"

	commonConstants "github.com/Orange-Health/citadel/common/constants"
	commonStructures "github.com/Orange-Health/citadel/common/structures"
	commonUtils "github.com/Orange-Health/citadel/common/utils"
	commonModels "github.com/Orange-Health/citadel/models"
)

type TestSampleMappingServiceInterface interface {
	GetTestSampleMappingsByOrderId(omsOrderId string) ([]commonStructures.TestSampleMappingInfo, *commonStructures.CommonError)
	GetTestSampleMappingsByOrderIds(omsOrderIds []string) ([]commonStructures.TestSampleMappingInfo,
		*commonStructures.CommonError)
	GetTestSampleMappingByTestIds(omsTestIds []string) ([]commonModels.TestSampleMapping, *commonStructures.CommonError)
	GetTestSampleMappingsModelsByOrderId(omsOrderId string) ([]commonModels.TestSampleMapping, *commonStructures.CommonError)
	GetTestSampleMappingByOrderIdWithTx(tx *gorm.DB, omsOrderId string) ([]commonModels.TestSampleMapping,
		*commonStructures.CommonError)
	GetTestSampleMappingByOrderIdAndTestIds(omsOrderId string, omsTestIds []string) ([]commonModels.TestSampleMapping,
		*commonStructures.CommonError)
	GetTestSampleMappingByOrderIdAndTestIdsUnscoped(omsOrderId string, omsTestIds []string) ([]commonModels.TestSampleMapping,
		*commonStructures.CommonError)
	OrdersWithRecollectionsPendingPresent(omsOrderIds, omsTestIds []string) (bool, *commonStructures.CommonError)
	OrdersWithRecollectionsPendingPresentByOmsOrderId(omsOrderId string) (bool, *commonStructures.CommonError)
	CheckIfAllTestsRejectedByOrderIdTestIdAndSampleNumber(omsOrderId, omsTestId string, sampleNumber uint) (bool, error)

	CreateTestSampleMappingWithTx(tx *gorm.DB, tsm commonModels.TestSampleMapping) (commonModels.TestSampleMapping,
		*commonStructures.CommonError)
	CreateBulkTestSampleMappingWithTx(tx *gorm.DB, tsmInfos []commonModels.TestSampleMapping) (
		[]commonModels.TestSampleMapping, *commonStructures.CommonError)
	UpdateBulkTestSampleMapping(tsmInfos []commonModels.TestSampleMapping) (
		[]commonModels.TestSampleMapping, *commonStructures.CommonError)
	UpdateBulkTestSampleMappingWithTx(tx *gorm.DB, tsmInfos []commonModels.TestSampleMapping) (
		[]commonModels.TestSampleMapping, *commonStructures.CommonError)
	UpdateTsmForTestIdAndRecollectionPendingTrueWithTx(tx *gorm.DB, omsTestIds []string) *commonStructures.CommonError
	DeleteTestSampleMappingByOrderIdTestIdAndSampleNumberWithTx(tx *gorm.DB,
		omsOrderId, omsTestId string, sampleNumber uint) *commonStructures.CommonError
	DeleteTestSampleMappingByOmsTestIdsWithTx(tx *gorm.DB, omsTestIds []string) *commonStructures.CommonError
	DeleteTestSampleMappingsWithTx(tx *gorm.DB, tsmInfos []commonModels.TestSampleMapping) *commonStructures.CommonError
	RejectTestSampleMappingBySampleIdWithTx(tx *gorm.DB, sampleId, userId uint) *commonStructures.CommonError
	RejectTestSampleMappingsByOmsOrderIdSampleNumberAndTestIdWithTx(tx *gorm.DB, omsOrderId string, sampleNumber uint,
		omsTestId string, rejectionReason string) *commonStructures.CommonError
}

func (tsms *TestSampleMappingService) CreateTestSampleMappingWithTx(tx *gorm.DB, tsm commonModels.TestSampleMapping) (
	commonModels.TestSampleMapping, *commonStructures.CommonError) {
	return tsms.TestSampleMappingDao.CreateTestSampleMappingWithTx(tx, tsm)
}

func (tsms *TestSampleMappingService) CreateBulkTestSampleMappingWithTx(tx *gorm.DB,
	tsmInfos []commonModels.TestSampleMapping) ([]commonModels.TestSampleMapping, *commonStructures.CommonError) {
	return tsms.TestSampleMappingDao.CreateBulkTestSampleMappingWithTx(tx, tsmInfos)
}

func (tsms *TestSampleMappingService) UpdateBulkTestSampleMapping(tsmInfos []commonModels.TestSampleMapping) (
	[]commonModels.TestSampleMapping, *commonStructures.CommonError) {
	return tsms.TestSampleMappingDao.UpdateBulkTestSampleMapping(tsmInfos)
}

func (tsms *TestSampleMappingService) UpdateBulkTestSampleMappingWithTx(tx *gorm.DB,
	tsmInfos []commonModels.TestSampleMapping) ([]commonModels.TestSampleMapping, *commonStructures.CommonError) {
	return tsms.TestSampleMappingDao.UpdateBulkTestSampleMappingWithTx(tx, tsmInfos)
}

func (tsms *TestSampleMappingService) UpdateRecollectionPendingForTestSampleMapping(
	omsOrderId, omsTestId string) *commonStructures.CommonError {
	tsmInfos, err := tsms.GetTestSampleMappingByOrderIdAndTestIdsUnscoped(omsOrderId, []string{omsTestId})
	if err != nil {
		return err
	}

	for idx := range tsmInfos {
		if !tsmInfos[idx].IsRejected {
			continue
		}
		tsmInfos[idx].RecollectionPending = false
		tsmInfos[idx].UpdatedAt = commonUtils.GetCurrentTime()
		tsmInfos[idx].UpdatedBy = commonConstants.CitadelSystemId
	}

	if _, err := tsms.UpdateBulkTestSampleMapping(tsmInfos); err != nil {
		return err
	}

	return nil
}

func (tsms *TestSampleMappingService) UpdateTsmForTestIdAndRecollectionPendingTrueWithTx(tx *gorm.DB,
	omsTestIds []string) *commonStructures.CommonError {
	return tsms.TestSampleMappingDao.UpdateTsmForTestIdAndRecollectionPendingTrueWithTx(tx, omsTestIds)
}

func (tsms *TestSampleMappingService) DeleteTestSampleMappingByOrderIdTestIdAndSampleNumberWithTx(tx *gorm.DB,
	omsOrderId, omsTestId string, sampleNumber uint) *commonStructures.CommonError {
	return tsms.TestSampleMappingDao.DeleteTestSampleMappingByOrderIdTestIdAndSampleNumberWithTx(tx, omsOrderId, omsTestId,
		sampleNumber)
}

func (tsms *TestSampleMappingService) DeleteTestSampleMappingByOmsTestIdsWithTx(tx *gorm.DB,
	omsTestIds []string) *commonStructures.CommonError {
	if len(omsTestIds) == 0 {
		return nil
	}
	return tsms.TestSampleMappingDao.DeleteTestSampleMappingByOmsTestIdsWithTx(tx, omsTestIds)
}

func (tsms *TestSampleMappingService) DeleteTestSampleMappingsWithTx(tx *gorm.DB,
	tsmInfos []commonModels.TestSampleMapping) *commonStructures.CommonError {
	if len(tsmInfos) == 0 {
		return nil
	}
	return tsms.TestSampleMappingDao.DeleteTestSampleMappingsWithTx(tx, tsmInfos)
}

func (tsms *TestSampleMappingService) RejectTestSampleMappingBySampleIdWithTx(tx *gorm.DB,
	sampleId, userId uint) *commonStructures.CommonError {
	return tsms.TestSampleMappingDao.RejectTestSampleMappingBySampleIdWithTx(tx, sampleId, userId)
}

func (tsms *TestSampleMappingService) RejectTestSampleMappingsByOmsOrderIdSampleNumberAndTestIdWithTx(tx *gorm.DB,
	omsOrderId string, sampleNumber uint, omsTestId string, rejectionReason string) *commonStructures.CommonError {
	tsmInfo, cErr := tsms.TestSampleMappingDao.GetTsmByOrderIdTestIdAndSampleNumber(omsOrderId, omsTestId, sampleNumber)
	if cErr != nil {
		return cErr
	}

	tsmInfo.IsRejected = true
	tsmInfo.RecollectionPending = true
	tsmInfo.RejectionReason = rejectionReason
	tsmInfo.UpdatedAt = commonUtils.GetCurrentTime()
	tsmInfo.UpdatedBy = commonConstants.CitadelSystemId

	cErr = tsms.TestSampleMappingDao.UpdateTestSampleMappingWithTx(tx, tsmInfo)
	if cErr != nil {
		return cErr
	}

	return nil
}

func (tsms *TestSampleMappingService) CheckIfAllTestsRejectedByOrderIdTestIdAndSampleNumber(omsOrderId, omsTestId string,
	sampleNumber uint) (bool, error) {
	return tsms.TestSampleMappingDao.CheckIfAllTestsRejectedByOrderIdTestIdAndSampleNumber(omsOrderId, omsTestId,
		sampleNumber)
}

func (tsms *TestSampleMappingService) GetTestSampleMappingByTestIds(omsTestIds []string) (
	[]commonModels.TestSampleMapping, *commonStructures.CommonError) {
	return tsms.TestSampleMappingDao.GetTestSampleMappingsByTestIds(omsTestIds)
}

func (tsms *TestSampleMappingService) GetTestSampleMappingByOrderIdAndTestIds(omsOrderId string, omsTestIds []string) (
	[]commonModels.TestSampleMapping, *commonStructures.CommonError) {
	tsmgs, cErr := tsms.TestSampleMappingDao.GetTestSampleMappingsByOrderIdAndTestIds(omsOrderId, omsTestIds)
	if cErr != nil {
		return []commonModels.TestSampleMapping{}, cErr
	}

	if len(tsmgs) == 0 {
		return []commonModels.TestSampleMapping{}, &commonStructures.CommonError{
			Message:    commonConstants.ERROR_NO_TEST_SAMPLE_MAPPING_FOUND,
			StatusCode: http.StatusNotFound,
		}
	}

	return tsmgs, nil
}

func (tsms *TestSampleMappingService) GetTestSampleMappingByOrderIdAndTestIdsUnscoped(omsOrderId string,
	omsTestIds []string) ([]commonModels.TestSampleMapping, *commonStructures.CommonError) {
	tsmgs, cErr := tsms.TestSampleMappingDao.GetTestSampleMappingsByOrderIdAndTestIdsUnscoped(omsOrderId, omsTestIds)
	if cErr != nil {
		return []commonModels.TestSampleMapping{}, cErr
	}

	if len(tsmgs) == 0 {
		return []commonModels.TestSampleMapping{}, &commonStructures.CommonError{
			Message:    commonConstants.ERROR_NO_TEST_SAMPLE_MAPPING_FOUND,
			StatusCode: http.StatusNotFound,
		}
	}

	return tsmgs, nil
}

func (tsms *TestSampleMappingService) GetTestSampleMappingsModelsByOrderId(omsOrderId string) (
	[]commonModels.TestSampleMapping, *commonStructures.CommonError) {
	return tsms.TestSampleMappingDao.GetTestSampleMappingsModelsByOrderId(omsOrderId)
}

func (tsms *TestSampleMappingService) GetTestSampleMappingByOrderIdWithTx(tx *gorm.DB, omsOrderId string) (
	[]commonModels.TestSampleMapping, *commonStructures.CommonError) {
	return tsms.TestSampleMappingDao.GetTestSampleMappingByOrderIdWithTx(tx, omsOrderId)
}

func (tsms *TestSampleMappingService) GetTestSampleMappingsByOrderId(omsOrderId string) (
	[]commonStructures.TestSampleMappingInfo, *commonStructures.CommonError) {
	return tsms.TestSampleMappingDao.GetTestSampleMappingsByOrderId(omsOrderId)
}

func (tsms *TestSampleMappingService) GetTestSampleMappingsByOrderIds(omsOrderIds []string) (
	[]commonStructures.TestSampleMappingInfo, *commonStructures.CommonError) {
	return tsms.TestSampleMappingDao.GetTestSampleMappingsByOrderIds(omsOrderIds)
}

func (tsms *TestSampleMappingService) OrdersWithRecollectionsPendingPresent(omsOrderIds, omsTestIds []string) (
	bool, *commonStructures.CommonError) {
	return tsms.TestSampleMappingDao.OrdersWithRecollectionsPendingPresent(omsOrderIds, omsTestIds)
}

func (tsms *TestSampleMappingService) OrdersWithRecollectionsPendingPresentByOmsOrderId(omsOrderId string) (
	bool, *commonStructures.CommonError) {
	return tsms.TestSampleMappingDao.OrdersWithRecollectionsPendingPresentByOmsOrderId(omsOrderId)
}
