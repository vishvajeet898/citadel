package service

import (
	"context"
	"fmt"
	"net/http"

	"gorm.io/gorm"

	"github.com/Orange-Health/citadel/apps/test_detail/constants"
	"github.com/Orange-Health/citadel/apps/test_detail/structures"
	commonConstants "github.com/Orange-Health/citadel/common/constants"
	commonStructures "github.com/Orange-Health/citadel/common/structures"
	commonUtils "github.com/Orange-Health/citadel/common/utils"
	commonModels "github.com/Orange-Health/citadel/models"
)

type TestDetailServiceInterface interface {
	GetTestDetailByIds(testDetailIDs []uint) (
		[]structures.TestDetail, *commonStructures.CommonError)
	GetTestDetailByIdsWithTx(tx *gorm.DB, testDetailIDs []uint) (
		[]structures.TestDetail, *commonStructures.CommonError)
	GetTestDetailModelById(testDetailID uint) (
		commonModels.TestDetail, *commonStructures.CommonError)
	GetTestDetailModelByOmsTestId(omsTestId string) (
		commonModels.TestDetail, *commonStructures.CommonError)
	GetTestDetailModelByOmsTestIds(omsTestIds []string) (
		[]commonModels.TestDetail, *commonStructures.CommonError)
	GetTestDetailByIdWithTx(tx *gorm.DB, testDetailID uint) (
		commonModels.TestDetail, *commonStructures.CommonError)
	GetTestDetailsByTaskId(taskID uint) (
		[]commonModels.TestDetail, *commonStructures.CommonError)
	GetTestDetailsByOmsOrderId(omsOrderId string) (
		[]commonModels.TestDetail, *commonStructures.CommonError)
	GetTestDetailsByOmsOrderIds(omsOrderIds []string) (
		[]commonModels.TestDetail, *commonStructures.CommonError)
	GetTestDetailsByOmsOrderIdWithTx(tx *gorm.DB, omsOrderId string) (
		[]commonModels.TestDetail, *commonStructures.CommonError)
	GetTestDetailsByTestIds(testIds []uint) (
		[]commonModels.TestDetail, *commonStructures.CommonError)
	GetTestBasicDetailsForSearchScreenByTaskIds(taskIds []uint) (
		[]structures.TestBasicDetails, *commonStructures.CommonError)
	GetOmsTestIdAndStatusByOmsOrderId(omsOrderId string) (
		[]structures.TestBasicDetails, *commonStructures.CommonError)
	GetAllBasicTestDetailsByOmsOrderId(ctx context.Context, omsOrderId string) (
		structures.TestBasicDetailsByOmsOrderId, *commonStructures.CommonError)
	GetTestDetailsByOmsOrderIdWithSampleStatus(omsOrderId, sampleStatus string) (
		[]commonModels.TestDetail, *commonStructures.CommonError)
	ContainsPackageTests(omsOrderId string) bool

	UpdateTestDetails(testDetails []commonModels.TestDetail) (
		[]commonModels.TestDetail, *commonStructures.CommonError)
	CreateTestDetailsWithTx(tx *gorm.DB, testDetails []commonModels.TestDetail) (
		[]commonModels.TestDetail, *commonStructures.CommonError)
	UpdateTestDetailsWithTx(tx *gorm.DB, testDetails []commonModels.TestDetail) (
		[]commonModels.TestDetail, *commonStructures.CommonError)
	UpdateTaskIdInTestDetailsWithOmsTestIdWithTx(tx *gorm.DB, centralOmsTestIds []string, taskId uint) *commonStructures.CommonError
	UpdateTaskIdInTestDetailsWithOmsTestIds(centralOmsTestIds []string, taskId uint) *commonStructures.CommonError
	UpdateTestDetailsById(testDetails commonModels.TestDetail) *commonStructures.CommonError
	UpdateTestStatusesByOmsTestIdsWithTx(tx *gorm.DB, testIdStatusMap map[string]string, userId uint) *commonStructures.CommonError
	UpdateProcessingLabForTestDetails(ctx context.Context, request structures.UpdateProcessingLabRequest) *commonStructures.CommonError
	DeleteTestDetailsByOmsTestIdsWithTx(tx *gorm.DB, omsTestIds []string) *commonStructures.CommonError
	DeleteTestDetailsByTaskIdAndOmsTestIdWithTx(tx *gorm.DB,
		taskID uint, omsTestId string) *commonStructures.CommonError
	DeleteTestDetailsByOmsOrderIdWithTx(tx *gorm.DB, omsOrderId string) *commonStructures.CommonError
	UpdateReportStatusByOmsTestIds(omsTestIds []string, oldStatus, newStatus string) *commonStructures.CommonError

	// Test Details Metadata
	GetTestDetailsMetadataByTestDetailIds(testDetailsIds []uint) (
		[]commonModels.TestDetailsMetadata, *commonStructures.CommonError)
	CreateTestDetailsMetadataWithTx(tx *gorm.DB,
		testDetailsMetadata []commonModels.TestDetailsMetadata) (
		[]commonModels.TestDetailsMetadata, *commonStructures.CommonError)
	UpdateTestDetailsMetadataWithTx(tx *gorm.DB,
		testDetailsMetadata []commonModels.TestDetailsMetadata) (
		[]commonModels.TestDetailsMetadata, *commonStructures.CommonError)
	UpdateDuplicateTestDetailsByTaskIdWithTx(tx *gorm.DB, taskId uint, masterTestIds []uint) *commonStructures.CommonError
	UpdatePickedAtTimeBasedOnActiveTests(taskID uint) *commonStructures.CommonError
	DeleteTestDetailsMetadataByTestDetailsIdsWithTx(tx *gorm.DB, testDetailsIds []uint) *commonStructures.CommonError
}

func (testDetailService *TestDetailService) GetTestDetailByIds(testDetailIds []uint) (
	[]structures.TestDetail, *commonStructures.CommonError) {
	testDetails, err := testDetailService.TestDetailDao.GetTestDetailByIds(testDetailIds)
	return testDetails, err
}

func (testDetailService *TestDetailService) GetTestDetailByIdsWithTx(tx *gorm.DB, testDetailIds []uint) (
	[]structures.TestDetail, *commonStructures.CommonError) {
	testDetails, err := testDetailService.TestDetailDao.GetTestDetailByIdsWithTx(tx, testDetailIds)
	return testDetails, err
}

func (testDetailService *TestDetailService) GetTestDetailModelById(testDetailId uint) (
	commonModels.TestDetail, *commonStructures.CommonError) {
	testDetail, err := testDetailService.TestDetailDao.GetTestDetailById(testDetailId)
	if testDetail.Id == 0 {
		return commonModels.TestDetail{}, &commonStructures.CommonError{
			Message:    commonConstants.ERROR_NO_TEST_DETAILS_FOUND,
			StatusCode: http.StatusNotFound,
		}
	}
	return testDetail, err
}

func (testDetailService *TestDetailService) GetTestDetailModelByOmsTestId(omsTestId string) (
	commonModels.TestDetail, *commonStructures.CommonError) {
	testDetail, err := testDetailService.TestDetailDao.GetTestDetailByOmsTestId(omsTestId)
	if testDetail.Id == 0 {
		return commonModels.TestDetail{}, &commonStructures.CommonError{
			Message:    commonConstants.ERROR_NO_TEST_DETAILS_FOUND,
			StatusCode: http.StatusNotFound,
		}
	}
	return testDetail, err
}

func (testDetailService *TestDetailService) GetTestDetailModelByOmsTestIds(omsTestIds []string) (
	[]commonModels.TestDetail, *commonStructures.CommonError) {
	testDetails, err := testDetailService.TestDetailDao.GetTestDetailByOmsTestIds(omsTestIds)
	if len(testDetails) == 0 {
		return []commonModels.TestDetail{}, &commonStructures.CommonError{
			Message:    commonConstants.ERROR_NO_TEST_DETAILS_FOUND,
			StatusCode: http.StatusNotFound,
		}
	}
	return testDetails, err
}

func (testDetailService *TestDetailService) GetTestDetailByIdWithTx(tx *gorm.DB, testDetailId uint) (
	commonModels.TestDetail, *commonStructures.CommonError) {
	testDetail, err := testDetailService.TestDetailDao.GetTestDetailByIdWithTx(tx, testDetailId)
	if testDetail.Id == 0 {
		return commonModels.TestDetail{}, &commonStructures.CommonError{
			Message:    commonConstants.ERROR_NO_TEST_DETAILS_FOUND,
			StatusCode: http.StatusNotFound,
		}
	}
	return testDetail, err
}

func (testDetailService *TestDetailService) GetTestDetailsByTaskId(taskID uint) (
	[]commonModels.TestDetail, *commonStructures.CommonError) {
	return testDetailService.TestDetailDao.GetTestDetailsByTaskId(taskID)
}

func (testDetailService *TestDetailService) GetTestDetailsByOmsOrderId(omsOrderId string) (
	[]commonModels.TestDetail, *commonStructures.CommonError) {
	testDetails, err := testDetailService.TestDetailDao.GetTestDetailsByOmsOrderId(omsOrderId)
	if err != nil {
		return []commonModels.TestDetail{}, err
	}

	if len(testDetails) == 0 {
		return []commonModels.TestDetail{}, &commonStructures.CommonError{
			Message:    commonConstants.ERROR_NO_TEST_DETAILS_FOUND,
			StatusCode: http.StatusNotFound,
		}
	}

	return testDetails, nil
}

func (testDetailService *TestDetailService) GetTestDetailsByOmsOrderIds(omsOrderIds []string) (
	[]commonModels.TestDetail, *commonStructures.CommonError) {
	testDetails, err := testDetailService.TestDetailDao.GetTestDetailsByOmsOrderIds(omsOrderIds)
	if err != nil {
		return []commonModels.TestDetail{}, err
	}
	if len(testDetails) == 0 {
		return []commonModels.TestDetail{}, &commonStructures.CommonError{
			Message:    commonConstants.ERROR_NO_TEST_DETAILS_FOUND,
			StatusCode: http.StatusNotFound,
		}
	}
	return testDetails, nil
}

func (testDetailService *TestDetailService) GetTestDetailsByOmsOrderIdWithTx(tx *gorm.DB, omsOrderId string) (
	[]commonModels.TestDetail, *commonStructures.CommonError) {
	testDetails, err := testDetailService.TestDetailDao.GetTestDetailsByOmsOrderIdWithTx(tx, omsOrderId)
	if err != nil {
		return []commonModels.TestDetail{}, err
	}

	if len(testDetails) == 0 {
		return []commonModels.TestDetail{}, &commonStructures.CommonError{
			Message:    commonConstants.ERROR_NO_TEST_DETAILS_FOUND,
			StatusCode: http.StatusNotFound,
		}
	}

	return testDetails, nil
}

func (testDetailService *TestDetailService) GetTestDetailsIdsByOmsOrderIdWithTx(tx *gorm.DB, omsOrderId string) (
	[]uint, *commonStructures.CommonError) {
	return testDetailService.TestDetailDao.GetTestDetailsIdsByOmsOrderIdWithTx(tx, omsOrderId)
}

func (testDetailService *TestDetailService) GetTestDetailsByTestIds(testIds []uint) (
	[]commonModels.TestDetail, *commonStructures.CommonError) {
	return testDetailService.TestDetailDao.GetTestDetailsByTestIds(testIds)
}

func (testDetailService *TestDetailService) GetTestBasicDetailsForSearchScreenByTaskIds(taskIds []uint) (
	[]structures.TestBasicDetails, *commonStructures.CommonError) {
	return testDetailService.TestDetailDao.GetTestBasicDetailsForSearchScreenByTaskIds(taskIds)
}

func (testDetailService *TestDetailService) GetOmsTestIdAndStatusByOmsOrderId(omsOrderId string) (
	[]structures.TestBasicDetails, *commonStructures.CommonError) {
	return testDetailService.TestDetailDao.GetOmsTestIdAndStatusByOmsOrderId(omsOrderId)
}

func (testDetailService *TestDetailService) GetAllBasicTestDetailsByOmsOrderId(ctx context.Context, omsOrderId string) (
	structures.TestBasicDetailsByOmsOrderId, *commonStructures.CommonError) {

	orderDetails, cErr := testDetailService.OrderDetailService.GetOrderDetailsByOmsOrderId(omsOrderId)
	if cErr != nil {
		return structures.TestBasicDetailsByOmsOrderId{}, cErr
	}

	omsTestBasicDetailDbStructs, cErr := testDetailService.TestDetailDao.GetAllBasicTestDetailsByOmsOrderId(omsOrderId)
	if cErr != nil {
		return structures.TestBasicDetailsByOmsOrderId{}, cErr
	}

	labIdMap := testDetailService.CdsService.GetLabIdLabMap(ctx)

	inhouseOmsTestBasicDetails, outsourceOmsTestBasicDetails :=
		[]structures.OmsTestBasicDetails{}, []structures.OmsTestBasicDetails{}
	for _, dbStruct := range omsTestBasicDetailDbStructs {
		inhouse := dbStruct.LabId == orderDetails.ServicingLabId ||
			dbStruct.LabId != orderDetails.ServicingLabId && labIdMap[dbStruct.LabId].Inhouse
		omsBasicTestDetail := structures.OmsTestBasicDetails{
			OrderID:      orderDetails.OmsOrderId,
			RequestID:    orderDetails.OmsRequestId,
			TestId:       dbStruct.TestId,
			TestName:     dbStruct.TestName,
			TestStatus:   fmt.Sprint(commonConstants.OmsTestStatusToUintMap[dbStruct.TestStatus]),
			MasterTestId: dbStruct.MasterTestId,
			LabId:        dbStruct.LabId,
			Inhouse:      inhouse,
		}
		if inhouse {
			inhouseOmsTestBasicDetails = append(inhouseOmsTestBasicDetails, omsBasicTestDetail)
		} else {
			outsourceOmsTestBasicDetails = append(outsourceOmsTestBasicDetails, omsBasicTestDetail)
		}
	}

	return structures.TestBasicDetailsByOmsOrderId{
		InhouseTests:    inhouseOmsTestBasicDetails,
		OutsourcedTests: outsourceOmsTestBasicDetails,
		CityCode:        orderDetails.CityCode,
	}, nil
}

func (testDetailService *TestDetailService) GetTestDetailsByOmsOrderIdWithSampleStatus(omsOrderId, sampleStatus string) (
	[]commonModels.TestDetail, *commonStructures.CommonError) {
	return testDetailService.TestDetailDao.GetTestDetailsByOmsOrderIdWithSampleStatus(omsOrderId, sampleStatus)
}

func (testDetailService *TestDetailService) ContainsPackageTests(omsOrderId string) bool {
	return testDetailService.TestDetailDao.ContainsPackageTests(omsOrderId)
}

func (testDetailService *TestDetailService) UpdateTestDetails(testDetails []commonModels.TestDetail) (
	[]commonModels.TestDetail, *commonStructures.CommonError) {
	if len(testDetails) == 0 {
		return testDetails, nil
	}
	return testDetailService.TestDetailDao.UpdateTestDetails(testDetails)
}

func (testDetailService *TestDetailService) CreateTestDetailsWithTx(tx *gorm.DB,
	testDetails []commonModels.TestDetail) (
	[]commonModels.TestDetail, *commonStructures.CommonError) {
	if len(testDetails) == 0 {
		return testDetails, nil
	}
	return testDetailService.TestDetailDao.CreateTestDetailsWithTx(tx, testDetails)
}

func (testDetailService *TestDetailService) UpdateTestDetailsById(
	testDetails commonModels.TestDetail) *commonStructures.CommonError {
	return testDetailService.TestDetailDao.UpdateTestDetailsById(testDetails)
}

func (testDetailService *TestDetailService) UpdateTestDetailsWithTx(tx *gorm.DB,
	testDetails []commonModels.TestDetail) (
	[]commonModels.TestDetail, *commonStructures.CommonError) {
	if len(testDetails) == 0 {
		return testDetails, nil
	}
	return testDetailService.TestDetailDao.UpdateTestDetailsWithTx(tx, testDetails)
}

func (testDetailService *TestDetailService) UpdateTaskIdInTestDetailsWithOmsTestIdWithTx(tx *gorm.DB, omsTestIds []string,
	taskId uint) *commonStructures.CommonError {
	if len(omsTestIds) == 0 {
		return nil
	}
	return testDetailService.TestDetailDao.UpdateTaskIdInTestDetailsWithOmsTestIdWithTx(tx, omsTestIds, taskId)
}

func (testDetailService *TestDetailService) UpdateTaskIdInTestDetailsWithOmsTestIds(centralOmsTestIds []string,
	taskId uint) *commonStructures.CommonError {
	if len(centralOmsTestIds) == 0 {
		return nil
	}
	return testDetailService.TestDetailDao.UpdateTaskIdInTestDetailsWithOmsTestIds(centralOmsTestIds, taskId)
}

func (testDetailService *TestDetailService) UpdateTestStatusesByOmsTestIdsWithTx(tx *gorm.DB,
	testIdStatusMap map[string]string, userId uint) *commonStructures.CommonError {
	return testDetailService.TestDetailDao.UpdateTestStatusesByOmsTestIdsWithTx(tx, testIdStatusMap, userId)
}

func (testDetailService *TestDetailService) UpdateProcessingLabForTestDetails(ctx context.Context,
	request structures.UpdateProcessingLabRequest) *commonStructures.CommonError {
	if len(request.TestDetails) == 0 {
		return &commonStructures.CommonError{
			Message:    commonConstants.ERROR_NO_TEST_DETAILS_FOUND,
			StatusCode: http.StatusBadRequest,
		}
	}

	loggingAttributes := map[string]interface{}{}
	loggingAttributes[constants.USER_ID] = request.UserId

	omsTestIds, omsTestIdToProcessingLabIdMap := []string{}, map[string]uint{}
	for _, testDetail := range request.TestDetails {
		if testDetail.OmsTestId == "" {
			return &commonStructures.CommonError{
				Message:    commonConstants.ERROR_INVALID_PROCESSING_LAB_REQUEST,
				StatusCode: http.StatusBadRequest,
			}
		}
		omsTestIds = append(omsTestIds, testDetail.OmsTestId)
		omsTestIdToProcessingLabIdMap[testDetail.OmsTestId] = testDetail.ProcessingLabId
	}

	labIdLabMap := testDetailService.CdsService.GetLabIdLabMap(ctx)

	testDetails, cErr := testDetailService.TestDetailDao.GetTestDetailByOmsTestIds(omsTestIds)
	if cErr != nil {
		return cErr
	}

	for index := range testDetails {
		if labIdLabMap[testDetails[index].ProcessingLabId].Inhouse {
			return &commonStructures.CommonError{
				Message:    commonConstants.ERROR_PROCESSING_LAB_IS_INHOUSE,
				StatusCode: http.StatusBadRequest,
			}
		}
		loggingAttributes[testDetails[index].CentralOmsTestId] = fmt.Sprintf(constants.PROCESSING_LAB_ID_MODIFIED_LOG,
			testDetails[index].ProcessingLabId, omsTestIdToProcessingLabIdMap[testDetails[index].CentralOmsTestId],
			testDetails[index].CentralOmsTestId)
		testDetails[index].ProcessingLabId = omsTestIdToProcessingLabIdMap[testDetails[index].CentralOmsTestId]
		testDetails[index].UpdatedBy = request.UserId
	}

	_, cErr = testDetailService.TestDetailDao.UpdateTestDetails(testDetails)
	if cErr != nil {
		return cErr
	}

	commonUtils.AddLog(ctx, commonConstants.INFO_LEVEL, constants.PROCESSING_LAB_ID_MODIFIED, loggingAttributes, nil)

	return nil
}

func (testDetailService *TestDetailService) DeleteTestDetailsByOmsTestIdsWithTx(tx *gorm.DB,
	omsTestIds []string) *commonStructures.CommonError {
	if len(omsTestIds) == 0 {
		return nil
	}
	return testDetailService.TestDetailDao.DeleteTestDetailsByOmsTestIdsWithTx(tx, omsTestIds)
}

func (testDetailService *TestDetailService) DeleteTestDetailsByTaskIdAndOmsTestIdWithTx(tx *gorm.DB, taskID uint,
	omsTestId string) *commonStructures.CommonError {
	return testDetailService.TestDetailDao.DeleteTestDetailsByTaskIdAndOmsTestIdWithTx(tx, taskID, omsTestId)
}

func (testDetailService *TestDetailService) DeleteTestDetailsByOmsOrderIdWithTx(tx *gorm.DB,
	omsOrderId string) *commonStructures.CommonError {
	testDetailsIds, cErr := testDetailService.GetTestDetailsIdsByOmsOrderIdWithTx(tx, omsOrderId)
	if cErr != nil {
		return cErr
	}

	if cErr := testDetailService.TestDetailDao.DeleteTestDetailsByOmsOrderIdWithTx(tx, omsOrderId); cErr != nil {
		return cErr
	}

	return testDetailService.DeleteTestDetailsMetadataByTestDetailsIdsWithTx(tx, testDetailsIds)
}

func (testDetailService *TestDetailService) UpdateReportStatusByOmsTestIds(omsTestIds []string,
	oldStatus, newStatus string) *commonStructures.CommonError {
	if len(omsTestIds) == 0 {
		return nil
	}

	return testDetailService.TestDetailDao.UpdateReportStatusByOmsTestIds(omsTestIds, oldStatus, newStatus)
}

func (testDetailService *TestDetailService) GetTestDetailsMetadataByTestDetailIds(testDetailsIds []uint) (
	[]commonModels.TestDetailsMetadata, *commonStructures.CommonError) {
	return testDetailService.TestDetailDao.GetTestDetailsMetadataByTestDetailIds(testDetailsIds)
}

func (testDetailService *TestDetailService) CreateTestDetailsMetadataWithTx(tx *gorm.DB,
	testDetailsMetadata []commonModels.TestDetailsMetadata) (
	[]commonModels.TestDetailsMetadata, *commonStructures.CommonError) {
	if len(testDetailsMetadata) == 0 {
		return testDetailsMetadata, nil
	}
	return testDetailService.TestDetailDao.CreateTestDetailsMetadataWithTransaction(tx, testDetailsMetadata)
}

func (testDetailService *TestDetailService) UpdateTestDetailsMetadataWithTx(tx *gorm.DB,
	testDetailsMetadata []commonModels.TestDetailsMetadata) (
	[]commonModels.TestDetailsMetadata, *commonStructures.CommonError) {
	if len(testDetailsMetadata) == 0 {
		return testDetailsMetadata, nil
	}
	return testDetailService.TestDetailDao.UpdateTestDetailsMetadataWithTransaction(tx, testDetailsMetadata)
}

func (testDetailService *TestDetailService) UpdateDuplicateTestDetailsByTaskIdWithTx(tx *gorm.DB, taskId uint,
	masterTestIds []uint) *commonStructures.CommonError {
	if len(masterTestIds) == 0 {
		return nil
	}
	return testDetailService.TestDetailDao.UpdateDuplicateTestDetailsByTaskIdWithTx(tx, taskId, masterTestIds)
}

func (testDetailService *TestDetailService) UpdatePickedAtTimeBasedOnActiveTests(taskId uint) *commonStructures.CommonError {
	activeTestDetails, cErr := testDetailService.TestDetailDao.GetActiveTestDetailsByTaskId(taskId)
	if cErr != nil {
		return cErr
	}

	if len(activeTestDetails) == 0 {
		return nil
	}

	activeTestDetailsIds := []uint{}
	for _, testDetail := range activeTestDetails {
		activeTestDetailsIds = append(activeTestDetailsIds, testDetail.Id)
	}
	return testDetailService.TestDetailDao.UpdatePickedAtTimeBasedOnActiveTests(activeTestDetailsIds)
}

func (testDetailService *TestDetailService) DeleteTestDetailsMetadataByTestDetailsIdsWithTx(tx *gorm.DB,
	testDetailsIds []uint) *commonStructures.CommonError {
	return testDetailService.TestDetailDao.DeleteTestDetailsMetadataByTestDetailsIdsWithTx(tx, testDetailsIds)
}
