package consumerTasks

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"gorm.io/gorm"

	"github.com/Orange-Health/citadel/common/constants"
	"github.com/Orange-Health/citadel/common/structures"
	"github.com/Orange-Health/citadel/common/utils"
	"github.com/Orange-Health/citadel/models"
)

func (eventProcessor *EventProcessor) LisEventTask(ctx context.Context, eventPayload string) error {
	lisEvent := structures.LisEvent{}
	err := json.Unmarshal(json.RawMessage(eventPayload), &lisEvent)
	if err != nil {
		eventProcessor.Sentry.LogError(ctx, constants.ERROR_FAILED_TO_UNMARSHAL_JSON, err, nil)
		return err
	}

	utils.AddLog(ctx, constants.INFO_LEVEL, utils.GetCurrentFunctionName(), map[string]interface{}{
		"entity_id":  lisEvent.EntityID,
		"event_type": constants.LisEvent,
	}, nil)

	redisKey := fmt.Sprintf(constants.LisEventKey, lisEvent.EntityID)
	keyExists, err := eventProcessor.Cache.Exists(ctx, redisKey)
	if err != nil || keyExists {
		utils.AddLog(ctx, constants.DEBUG_LEVEL, utils.GetCurrentFunctionName(), nil, err)
		return errors.New(constants.ERROR_LIS_EVENT_TASK_IN_PROGRESS)
	}

	err = eventProcessor.Cache.Set(ctx, redisKey, true, constants.CacheExpiry10MinutesInt)
	if err != nil {
		utils.AddLog(ctx, constants.DEBUG_LEVEL, utils.GetCurrentFunctionName(), nil, err)
		return err
	}

	defer func() {
		err := eventProcessor.Cache.Delete(ctx, redisKey)
		if err != nil {
			utils.AddLog(ctx, constants.DEBUG_LEVEL, utils.GetCurrentFunctionName(), nil, err)
		}
	}()

	task := models.Task{}
	lisOrderUpdateDetailsEvent, cErr := GetLisVisitDataEventFromBase64Data(lisEvent.EntityID, lisEvent.WebhookData)
	if cErr != nil {
		return errors.New(cErr.Message)
	}

	if lisOrderUpdateDetailsEvent.ReportPdfFormat == constants.AttuneReportWithoutStationery {
		utils.AddLog(ctx, constants.INFO_LEVEL, utils.GetCurrentFunctionName(), nil,
			errors.New(constants.ERROR_REPORT_PDF_FORMAT_IGNORED))
		return nil
	}

	utils.AddLog(ctx, constants.INFO_LEVEL, utils.GetCurrentFunctionName(), map[string]interface{}{
		"entity_id":         lisEvent.EntityID,
		"event_type":        constants.LisEvent,
		"lis_order_details": lisOrderUpdateDetailsEvent,
	}, nil)

	sample, cErr := eventProcessor.AttuneService.GetSampleByVisitId(lisEvent.EntityID)
	if cErr != nil {
		return errors.New(cErr.Message)
	}
	if sample.Id == 0 {
		return nil
	}

	orderDetails, cErr := eventProcessor.OrderDetailsService.GetOrderDetailsByOmsOrderId(sample.OmsOrderId)
	if cErr != nil {
		return errors.New(cErr.Message)
	}
	if orderDetails.Id == 0 {
		utils.AddLog(ctx, constants.ERROR_LEVEL, constants.ERROR_ORDER_ID_NOT_FOUND, nil, nil)
		return nil
	}

	patientDetails, cErr := eventProcessor.PatientDetailService.GetPatientDetailsById(orderDetails.PatientDetailsId)
	if cErr != nil {
		return errors.New(cErr.Message)
	}

	pathologistsList, cErr := eventProcessor.UserService.GetAllPathologistsModels()
	if cErr != nil {
		utils.AddLog(ctx, constants.ERROR_LEVEL, utils.GetCurrentFunctionName(), nil, errors.New(cErr.Message))
		return errors.New(cErr.Message)
	}

	masterTestDepartmentMap, err := eventProcessor.fetchMasterTestDepartmentMap(ctx)
	if err != nil {
		return err
	}

	testDetails, cErr := eventProcessor.SampleService.GetTestDetailsForLisEventByVisitId(lisEvent.EntityID)
	if cErr != nil {
		return errors.New(cErr.Message)
	}
	if len(testDetails) == 0 {
		return errors.New(constants.ERROR_NO_TEST_DETAILS_FOUND)
	}

	sampleLab, cErr := eventProcessor.CdsService.GetLabById(ctx, sample.LabId)
	if cErr != nil {
		return errors.New(cErr.Message)
	}

	lisOrderUpdateDetailsEvent = eventProcessor.ResyncLisOrderDataInCaseOfAnomaly(ctx, lisOrderUpdateDetailsEvent,
		testDetails, orderDetails, patientDetails, pathologistsList, sample.LabId)

	// Trigger Attune Report Generation
	if _, ok := lisOrderUpdateDetailsEvent.OrderInfo[constants.ATTUNE_TEST_STATUS_APPROVED]; ok {
		go eventProcessor.SaveReportFromLis(orderDetails, patientDetails, lisOrderUpdateDetailsEvent, testDetails,
			sample.LabId)
	}
	testDocumentMap := map[string][]structures.TestDocumentInfoResponse{}
	labIdLabMap := eventProcessor.CdsService.GetLabIdLabMap(ctx)
	if _, ok := lisOrderUpdateDetailsEvent.OrderInfo[constants.ATTUNE_TEST_STATUS_APPROVED]; ok {
		omsApprovedEvent, cErr := eventProcessor.GetLisEventDetails(ctx, orderDetails, patientDetails,
			constants.ATTUNE_TEST_STATUS_APPROVED, lisOrderUpdateDetailsEvent, testDetails)
		if cErr != nil {
			return errors.New(cErr.Message)
		}
		task, testDocumentMap, cErr = eventProcessor.processLisApprovedEvent(ctx, omsApprovedEvent, masterTestDepartmentMap, pathologistsList,
			patientDetails.SystemPatientId, sampleLab, labIdLabMap)
		if cErr != nil {
			return errors.New(cErr.Message)
		}
	}
	if _, ok := lisOrderUpdateDetailsEvent.OrderInfo[constants.ATTUNE_TEST_STATUS_COMPLETED]; ok {
		omsCompletedEvent, cErr := eventProcessor.GetLisEventDetails(ctx, orderDetails, patientDetails,
			constants.ATTUNE_TEST_STATUS_COMPLETED, lisOrderUpdateDetailsEvent, testDetails)
		if cErr != nil {
			return errors.New(cErr.Message)
		}
		task, testDocumentMap, cErr = eventProcessor.processLisCompletedEvent(ctx, omsCompletedEvent, masterTestDepartmentMap,
			pathologistsList, patientDetails.SystemPatientId, sampleLab, labIdLabMap)
		if cErr != nil {
			return errors.New(cErr.Message)
		}
	}
	if _, ok := lisOrderUpdateDetailsEvent.OrderInfo[constants.ATTUNE_TEST_STATUS_RERUN]; ok {
		omsRerunEvent, cErr := eventProcessor.GetLisRerunEventDetails(ctx, orderDetails,
			lisOrderUpdateDetailsEvent, testDetails)
		if cErr != nil {
			return errors.New(cErr.Message)
		}
		err = eventProcessor.processLisRerunEventTask(ctx, omsRerunEvent)
		if err != nil {
			return err
		}
	}

	visitDocuments := []string{}
	for _, visitMap := range lisOrderUpdateDetailsEvent.VisitDocumentInfo {
		visitDocuments = append(visitDocuments, visitMap.VisitDocument)
	}
	if len(visitDocuments) > 0 && task.Id > 0 {
		err = eventProcessor.CommonTaskProcessor.AddVisitDocumentTaskByLisEvent(ctx, lisEvent.EntityID, task.Id,
			visitDocuments)
		if err != nil {
			utils.AddLog(ctx, constants.ERROR_LEVEL, utils.GetCurrentFunctionName(), nil, err)
		}
	}

	if len(testDocumentMap) > 0 {
		err = eventProcessor.CommonTaskProcessor.AddTestDocumentTaskByLisEvent(ctx, task.Id, testDocumentMap)
		if err != nil {
			utils.AddLog(ctx, constants.ERROR_LEVEL, utils.GetCurrentFunctionName(), nil, err)
		}
	}

	if task.Id != 0 {
		cErr := eventProcessor.ReleaseReport(ctx, task.Id)
		if cErr != nil {
			utils.AddLog(ctx, constants.ERROR_LEVEL, utils.GetCurrentFunctionName(), nil, errors.New(cErr.Message))
		}
	}

	basicTestDetails, cErr := eventProcessor.TestDetailService.GetOmsTestIdAndStatusByOmsOrderId(orderDetails.OmsOrderId)
	if cErr != nil {
		return nil
	}

	approvedTestIds, resultSavedTestIds := []string{}, []string{}
	for _, testDetail := range basicTestDetails {
		switch testDetail.Status {
		case constants.TEST_STATUS_APPROVE:
			approvedTestIds = append(approvedTestIds, testDetail.OmsTestId)
		case constants.TEST_STATUS_RESULT_SAVED, constants.TEST_STATUS_RERUN_RESULT_SAVED:
			resultSavedTestIds = append(resultSavedTestIds, testDetail.OmsTestId)
		}
	}

	if len(approvedTestIds) > 0 {
		eventProcessor.EtsService.GetAndPublishEtsTestEventForLisWebhook(ctx, approvedTestIds, constants.AttuneTestStatusApprove)
	}
	if len(resultSavedTestIds) > 0 {
		eventProcessor.EtsService.GetAndPublishEtsTestEventForLisWebhook(ctx, resultSavedTestIds,
			constants.AttuneTestStatusCompleted)
	}

	return nil
}

func GetLisVisitDataEventFromBase64Data(lisOrderId string, lisPayload string) (structures.LisOrderUpdateDetails,
	*structures.CommonError) {
	lisOrderInfoEvent := structures.LisOrderInfo{}

	lisOrderInfoEventByte, err := base64.StdEncoding.DecodeString(lisPayload)
	if err != nil {
		return structures.LisOrderUpdateDetails{}, &structures.CommonError{
			Message:    constants.ERROR_FAILED_TO_DECODE_BASE64_STRING,
			StatusCode: http.StatusInternalServerError,
		}
	}

	err = json.Unmarshal(lisOrderInfoEventByte, &lisOrderInfoEvent)
	if err != nil {
		return structures.LisOrderUpdateDetails{}, &structures.CommonError{
			Message:    constants.ERROR_FAILED_TO_UNMARSHAL_JSON,
			StatusCode: http.StatusInternalServerError,
		}
	}

	lisOrderUpdateDetails := structures.LisOrderUpdateDetails{}
	lisOrderUpdateDetails.LisVisitId = lisOrderId
	allOrderInfo := lisOrderInfoEvent.OrderInfo
	lisOrderUpdateDetails.OrderInfo = map[string]map[string]structures.LisTestUpdateInfo{}

	for _, orderInfo := range allOrderInfo {
		testStatus := orderInfo.TestStatus
		testCode := orderInfo.TestCode
		_, ok := lisOrderUpdateDetails.OrderInfo[testStatus]
		if !ok {
			lisOrderUpdateDetails.OrderInfo[testStatus] = make(map[string]structures.LisTestUpdateInfo)
		}
		lisOrderUpdateDetails.OrderInfo[testStatus][testCode] = structures.LisTestUpdateInfo{
			TestCode: testCode,
			TestName: orderInfo.TestName,
			MetaData: orderInfo,
		}
	}

	lisOrderUpdateDetails.PdfResult = lisOrderInfoEvent.ResultAsPdf
	lisOrderUpdateDetails.ReportPdfFormat = lisOrderInfoEvent.ReportPdfFormat
	lisOrderUpdateDetails.VisitDocumentInfo = lisOrderInfoEvent.VisitDocumentinfo

	return lisOrderUpdateDetails, nil
}

func GetLisVisitDataFromAttuneOrderResponse(orderResponse structures.AttuneOrderResponse,
	reportPdfFormat string) structures.LisOrderUpdateDetails {
	lisOrderUpdateDetails := structures.LisOrderUpdateDetails{}
	lisOrderUpdateDetails.LisVisitId = orderResponse.OrderId
	allOrderInfo := orderResponse.OrderInfo
	lisOrderUpdateDetails.OrderInfo = map[string]map[string]structures.LisTestUpdateInfo{}

	for _, orderInfo := range allOrderInfo {
		testStatus := orderInfo.TestStatus
		testCode := orderInfo.TestCode
		_, ok := lisOrderUpdateDetails.OrderInfo[testStatus]
		if !ok {
			lisOrderUpdateDetails.OrderInfo[testStatus] = make(map[string]structures.LisTestUpdateInfo)
		}
		lisOrderUpdateDetails.OrderInfo[testStatus][testCode] = structures.LisTestUpdateInfo{
			TestCode: testCode,
			TestName: orderInfo.TestName,
			MetaData: orderInfo,
		}
	}

	lisOrderUpdateDetails.PdfResult = orderResponse.ResultAsPdf
	lisOrderUpdateDetails.ReportPdfFormat = reportPdfFormat
	lisOrderUpdateDetails.VisitDocumentInfo = orderResponse.VisitDocumentinfo

	return lisOrderUpdateDetails
}

func fetchAllTestCodesFromMetadata(orderInfo map[string]map[string]structures.LisTestUpdateInfo) []string {
	testCodes := []string{}
	for attuneStatus, testResults := range orderInfo {
		if attuneStatus == constants.AttuneTestStatusCompleted ||
			attuneStatus == constants.AttuneTestStatusApprove {
			for testCode := range testResults {
				if utils.SliceContainsString(constants.TestCodesToBeSkippedInRetry, testCode) {
					continue
				}
				testCodes = append(testCodes, testCode)
			}
		}
	}
	return testCodes
}

func flattenCdsResponse(cdsTestDetails structures.InvestigationDetails) (
	map[string]uint, map[uint][]string, map[string]structures.Investigation) {
	panels, investigations := cdsTestDetails.Panels, cdsTestDetails.Investigations
	testCodeMimmMap, mimmTestCodeMap := map[string]uint{}, map[uint][]string{}
	testCodeInvestigationMap := map[string]structures.Investigation{}
	index, queue := 0, []structures.Panel{}
	queue = append(queue, panels...)

	for _, investigation := range cdsTestDetails.Investigations {
		investigations = append(investigations, investigation)
		testCodeMimmMap[investigation.LisCode] = investigation.InvestigationMethodMappingId
		if _, keyExists := mimmTestCodeMap[investigation.InvestigationMethodMappingId]; !keyExists {
			mimmTestCodeMap[investigation.InvestigationMethodMappingId] = []string{}
		}
		mimmTestCodeMap[investigation.InvestigationMethodMappingId] = append(
			mimmTestCodeMap[investigation.InvestigationMethodMappingId], investigation.LisCode)
		testCodeInvestigationMap[investigation.LisCode] = investigation
	}

	for index < len(queue) {
		currentNode := queue[index]
		index += 1

		for _, investigation := range currentNode.Investigations {
			investigations = append(investigations, investigation)
			testCodeMimmMap[investigation.LisCode] = investigation.InvestigationMethodMappingId
			if _, keyExists := mimmTestCodeMap[investigation.InvestigationMethodMappingId]; !keyExists {
				mimmTestCodeMap[investigation.InvestigationMethodMappingId] = []string{}
			}
			mimmTestCodeMap[investigation.InvestigationMethodMappingId] = append(
				mimmTestCodeMap[investigation.InvestigationMethodMappingId], investigation.LisCode)
			testCodeInvestigationMap[investigation.LisCode] = investigation
		}

		queue = append(queue, currentNode.Panels...)
	}

	return testCodeMimmMap, mimmTestCodeMap, testCodeInvestigationMap
}

func (eventProcessor *EventProcessor) SendAlertsInCaseOfAnomaly(ctx context.Context,
	orderDetails models.OrderDetails, patientName, visitId string,
	missingTestsInCds, missingTestsInLis, blankTests, incorrectResultTypes, doctorDetailsMissingTests []string,
	testCodeInvestigationMap map[string]structures.Investigation, testCodeNameMapLis map[string]string) {

	if len(missingTestsInCds) > 0 {
		slackMessage := utils.GetSlackMessageForMissingParametersInCds(patientName, visitId, orderDetails.CityCode,
			orderDetails.OmsOrderId, orderDetails.OmsRequestId, missingTestsInCds, testCodeNameMapLis)
		eventProcessor.HealthApiClient.SendGenericSlackMessage(ctx, constants.SlackMissingParametersChannel, slackMessage)
	}
	if len(missingTestsInLis) > 0 {
		slackMessage := utils.GetSlackMessageForMissingParametersInLis(patientName, visitId, orderDetails.CityCode,
			orderDetails.OmsOrderId, orderDetails.OmsRequestId, missingTestsInLis, testCodeInvestigationMap)
		eventProcessor.HealthApiClient.SendGenericSlackMessage(ctx, constants.SlackMissingParametersChannel, slackMessage)
	}
	if len(blankTests) > 0 {
		slackMessage := utils.GetSlackMessageForBlankValues(patientName, visitId, orderDetails.CityCode,
			orderDetails.OmsOrderId, orderDetails.OmsRequestId, blankTests, testCodeInvestigationMap)
		eventProcessor.HealthApiClient.SendGenericSlackMessage(ctx, constants.SlackMissingParametersChannel, slackMessage)
	}
	if len(incorrectResultTypes) > 0 {
		slackMessage := utils.GetSlackMessageForIncorrectResultTypes(patientName, visitId, orderDetails.CityCode,
			orderDetails.OmsOrderId, orderDetails.OmsRequestId, incorrectResultTypes, testCodeInvestigationMap)
		eventProcessor.HealthApiClient.SendGenericSlackMessage(ctx, constants.SlackMissingParametersChannel, slackMessage)
	}
	if len(doctorDetailsMissingTests) > 0 {
		slackMessage := utils.GetSlackMessageForDoctorDetailsMissing(patientName, visitId, orderDetails.CityCode,
			orderDetails.OmsOrderId, orderDetails.OmsRequestId, doctorDetailsMissingTests)
		eventProcessor.HealthApiClient.SendGenericSlackMessage(ctx, constants.SlackMissingParametersChannel, slackMessage)
	}
}

func (eventProcessor *EventProcessor) ResyncLisOrderDataInCaseOfAnomaly(ctx context.Context,
	lisOrderUpdateDetails structures.LisOrderUpdateDetails, omsTestDetails []structures.TestDetailsForLisEvent,
	orderDetails models.OrderDetails, patientDetails models.PatientDetail,
	pathologistsList []models.User, sampleLabId uint) structures.LisOrderUpdateDetails {

	attunePathologistsUserIds := []string{}
	for _, pathologist := range pathologistsList {
		attunePathologistsUserIds = append(attunePathologistsUserIds, pathologist.AttuneUserId)
	}

	// Fetch test-codes from lisOrderUpdateDetails
	testCodes := fetchAllTestCodesFromMetadata(lisOrderUpdateDetails.OrderInfo)
	if len(testCodes) == 0 {
		return lisOrderUpdateDetails
	}

	var err error
	maxRetries := constants.LisMissingValuesMaxRetries
	newLisOrderUpdateDetails := lisOrderUpdateDetails
	retryDelay := constants.LisMissingValuesBackoffTime

	patientDob := patientDetails.Dob
	if patientDob == nil {
		patientDob = patientDetails.ExpectedDob
	}

	if patientDob == nil {
		return lisOrderUpdateDetails
	}
	patientDobString := patientDob.Format(constants.DateLayout)
	patientGender := utils.GetGenderConstant(patientDetails.Gender)

	cdsTestDetails, err := eventProcessor.CdsClient.GetPanelDetails(ctx, testCodes, nil, orderDetails.CityCode,
		orderDetails.ServicingLabId, true, false, patientDobString, patientGender)
	if err != nil {
		utils.AddLog(ctx, constants.ERROR_LEVEL, utils.GetCurrentFunctionName(), nil, err)
		return lisOrderUpdateDetails
	}

	testCodeMimmMap, mimmTestCodeMap, testCodeInvestigationMap := flattenCdsResponse(cdsTestDetails.InvestigationDetails)

	missingTestsInCds, missingTestsInLis, blankTests, incorrectResultTypes, testCodeNameMapLis, doctorDetailsMissingTests :=
		utils.ValidatePayload(lisOrderUpdateDetails, omsTestDetails, testCodeInvestigationMap, testCodeMimmMap,
			mimmTestCodeMap, attunePathologistsUserIds)
	if len(missingTestsInCds) == 0 && len(missingTestsInLis) == 0 && len(blankTests) == 0 &&
		len(incorrectResultTypes) == 0 && len(doctorDetailsMissingTests) == 0 {
		return lisOrderUpdateDetails
	}

	utils.AddLog(ctx, constants.INFO_LEVEL, utils.GetCurrentFunctionName(), map[string]interface{}{
		"resyncing":                    true,
		"visit_id":                     lisOrderUpdateDetails.LisVisitId,
		"missing_tests_in_cds":         missingTestsInCds,
		"missing_tests_in_lis":         missingTestsInLis,
		"blank_tests":                  blankTests,
		"incorrect_result_types":       incorrectResultTypes,
		"doctor_details_missing_tests": doctorDetailsMissingTests,
	}, nil)
	for attempt := range maxRetries {
		attuneOrderResponse, cErr := eventProcessor.AttuneClient.GetPatientVisitDetailsbyVisitNo(ctx,
			newLisOrderUpdateDetails.LisVisitId, constants.AttuneReportWithStationery, sampleLabId)
		if cErr != nil {
			return lisOrderUpdateDetails
		}

		newLisOrderUpdateDetails = GetLisVisitDataFromAttuneOrderResponse(attuneOrderResponse,
			constants.AttuneReportWithStationery)
		lisOrderUpdateDetails = newLisOrderUpdateDetails
		testCodes := fetchAllTestCodesFromMetadata(lisOrderUpdateDetails.OrderInfo)
		cdsTestDetails, err := eventProcessor.CdsClient.GetPanelDetails(ctx, testCodes, nil, orderDetails.CityCode,
			orderDetails.ServicingLabId, true, false, patientDobString, patientGender)
		if err != nil {
			return lisOrderUpdateDetails
		}

		testCodeMimmMap, mimmTestCodeMap, testCodeInvestigationMap :=
			flattenCdsResponse(cdsTestDetails.InvestigationDetails)
		missingTestsInCds, missingTestsInLis, blankTests, incorrectResultTypes,
			testCodeNameMapLis, doctorDetailsMissingTests = utils.ValidatePayload(lisOrderUpdateDetails, omsTestDetails,
			testCodeInvestigationMap, testCodeMimmMap, mimmTestCodeMap, attunePathologistsUserIds)
		if len(missingTestsInCds) == 0 && len(missingTestsInLis) == 0 && len(blankTests) == 0 &&
			len(incorrectResultTypes) == 0 && len(doctorDetailsMissingTests) == 0 {
			return lisOrderUpdateDetails
		}

		time.Sleep(utils.GetExponentialBackoff(attempt, retryDelay))
	}

	eventProcessor.SendAlertsInCaseOfAnomaly(ctx, orderDetails, patientDetails.Name, lisOrderUpdateDetails.LisVisitId,
		missingTestsInCds, missingTestsInLis, blankTests, incorrectResultTypes, doctorDetailsMissingTests,
		testCodeInvestigationMap, testCodeNameMapLis)

	return lisOrderUpdateDetails
}

func (eventProcessor *EventProcessor) GetLisEventDetails(ctx context.Context, orderDetails models.OrderDetails,
	patientDetails models.PatientDetail, lisStatus string, lisOrderUpdateDetails structures.LisOrderUpdateDetails,
	testDetails []structures.TestDetailsForLisEvent) (structures.OmsLisEvent, *structures.CommonError) {

	omsLisEvent := structures.OmsLisEvent{}
	orderInfo := lisOrderUpdateDetails.OrderInfo[lisStatus]
	servicingLabId := orderDetails.ServicingLabId
	omsOrderId := orderDetails.OmsOrderId
	partner, doctor := structures.Partner{}, structures.Doctor{}
	var err error

	if orderDetails.PartnerId != 0 {
		partner, err = eventProcessor.PartnerApiClient.GetPartnerById(ctx, orderDetails.PartnerId)
		if err != nil {
			utils.AddLog(ctx, constants.ERROR_LEVEL, utils.GetCurrentFunctionName(), nil, err)
		}
	}

	if orderDetails.DoctorId != 0 {
		doctor, err = eventProcessor.HealthApiClient.GetDoctorById(ctx, orderDetails.DoctorId)
		if err != nil {
			utils.AddLog(ctx, constants.ERROR_LEVEL, utils.GetCurrentFunctionName(), nil, err)
		}
	}

	visitIds, cErr := eventProcessor.SampleService.GetVisitIdsByOmsOrderId(omsOrderId)
	if cErr != nil {
		return omsLisEvent, cErr
	}

	servicingLab, cErr := eventProcessor.CdsService.GetLabById(ctx, servicingLabId)
	if cErr != nil {
		return omsLisEvent, cErr
	}

	omsLisEvent = CreateLisEvent(servicingLab, orderDetails, patientDetails, doctor, partner, visitIds, testDetails,
		orderInfo)
	return omsLisEvent, nil
}

func (eventProcessor *EventProcessor) GetLisRerunEventDetails(ctx context.Context, orderDetails models.OrderDetails,
	lisOrderUpdateDetails structures.LisOrderUpdateDetails, testDetails []structures.TestDetailsForLisEvent) (
	structures.OmsRerunEvent, *structures.CommonError) {

	omsRerunEvent := structures.OmsRerunEvent{}
	omsRerunEvent.AlnumOrderId = orderDetails.OmsOrderId
	servicingLabId := orderDetails.ServicingLabId

	servicingLab, cErr := eventProcessor.CdsService.GetLabById(ctx, servicingLabId)
	if cErr != nil {
		return omsRerunEvent, cErr
	}

	omsRerunEvent.CityCode = servicingLab.City

	omsRerunEvent.Tests.OrderInfo = lisOrderUpdateDetails.OrderInfo[constants.ATTUNE_TEST_STATUS_RERUN]

	testDetailList := []structures.OmsTestStruct{}
	for _, testDetail := range testDetails {
		if _, keyExists := omsRerunEvent.Tests.OrderInfo[testDetail.TestCode]; !keyExists {
			continue
		}
		testDetailList = append(testDetailList, structures.OmsTestStruct{
			TestId:       testDetail.TestId,
			TestCode:     testDetail.TestCode,
			MasterTestId: testDetail.MasterTestId,
			TestName:     testDetail.TestName,
			TestType:     testDetail.TestType,
			Barcodes:     strings.Split(testDetail.Barcodes, ","),
		})
	}
	omsRerunEvent.Tests.TestDetails = testDetailList

	return omsRerunEvent, nil
}

func CreateLisEvent(lab structures.Lab, orderDetails models.OrderDetails, patientDetails models.PatientDetail,
	doctor structures.Doctor, partner structures.Partner, visitIds []string, testDetails []structures.TestDetailsForLisEvent,
	orderInfo map[string]structures.LisTestUpdateInfo,
) structures.OmsLisEvent {

	containsPackage := false
	testDetailList := []structures.OmsTestStruct{}
	for _, testDetail := range testDetails {
		testDetailList = append(testDetailList, structures.OmsTestStruct{
			TestId:       testDetail.TestId,
			TestCode:     testDetail.TestCode,
			MasterTestId: testDetail.MasterTestId,
			TestName:     testDetail.TestName,
			TestType:     testDetail.TestType,
			Barcodes:     strings.Split(testDetail.Barcodes, ","),
		})
		if testDetail.MasterPackageId != 0 {
			containsPackage = true
		}
	}

	updatedVisitDetails := []structures.OmsVisitDetails{}
	for _, visitId := range visitIds {
		if visitId == "" {
			continue
		}
		updatedVisitDetails = append(updatedVisitDetails, structures.OmsVisitDetails{
			Id: visitId,
		})
	}

	citadelEvent := structures.OmsLisEvent{
		Request: structures.OmsRequestDetails{
			AlnumRequestId: orderDetails.OmsRequestId,
		},
		Order: structures.OmsOrderDetails{
			AlnumOrderId:    orderDetails.OmsOrderId,
			LabId:           lab.Id,
			CityCode:        lab.City,
			OrderType:       constants.CollecTypeToOrderTypeMap[orderDetails.CollectionType],
			ContainsPackage: containsPackage,
			DoctorName:      doctor.Name,
			DoctorNumber:    doctor.Number,
			PartnerName:     partner.PartnerName,
		},
		Visits: updatedVisitDetails,
		Patient: structures.OmsPatientDetails{
			Name:        patientDetails.Name,
			ExpectedDob: patientDetails.ExpectedDob,
			Dob:         patientDetails.Dob,
			Gender:      patientDetails.Gender,
			Number:      patientDetails.Number,
			PatientId:   patientDetails.SystemPatientId,
		},
		Tests: structures.OmsTestDetails{
			OrderInfo:   orderInfo,
			TestDetails: testDetailList,
		},
	}
	return citadelEvent
}

func (eventProcessor *EventProcessor) processLisApprovedEvent(ctx context.Context, omsApprovedEvent structures.OmsLisEvent,
	masterTestDepartmentMap map[string]string, usersList []models.User, systemPatientId string, sampleLab structures.Lab,
	labIdLabMap map[uint]structures.Lab) (
	models.Task, map[string][]structures.TestDocumentInfoResponse, *structures.CommonError) {

	var cErr *structures.CommonError
	task := models.Task{}
	masterInvestigationIdInvestigationResultMetadataMap := map[uint]models.InvestigationResultMetadata{}
	masterInvestigationDetailsMap, err := eventProcessor.fetchInvestigationDetailsFromCds(ctx,
		omsApprovedEvent.Tests, sampleLab.City, sampleLab.Id, omsApprovedEvent.Patient)
	if err != nil {
		return task, nil, &structures.CommonError{
			Message:    err.Error(),
			StatusCode: http.StatusInternalServerError,
		}
	}
	userIdToAttuneUserIdMap := map[uint]int{}
	for _, user := range usersList {
		attuneUserId := utils.ConvertStringToInt(user.AttuneUserId)
		userIdToAttuneUserIdMap[user.Id] = attuneUserId
	}
	testIdInvestigationResultsMap, _, testIdInitialTestDetailsMap, investigationCodeInvestigationDataMap,
		investigationCodeMedicalRemarkMap, investigationCodeTechnicianRemarkMap, _, testDocumentMap :=
		eventProcessor.createFlattenedInvestigationResultsAndInitialTestDetailsMap(ctx, constants.OmsApprovedEvent,
			omsApprovedEvent.Order.CityCode, omsApprovedEvent.Tests, masterInvestigationDetailsMap, usersList,
			systemPatientId, sampleLab.Id)
	task, cErr = eventProcessor.TaskService.GetTaskByOmsOrderId(omsApprovedEvent.Order.AlnumOrderId)
	if task.Id != 0 {
		task, cErr = eventProcessor.createUpdateData(ctx, omsApprovedEvent, task,
			testIdInvestigationResultsMap, masterInvestigationIdInvestigationResultMetadataMap,
			investigationCodeInvestigationDataMap, testIdInitialTestDetailsMap, masterTestDepartmentMap,
			investigationCodeMedicalRemarkMap, investigationCodeTechnicianRemarkMap, []string{}, userIdToAttuneUserIdMap,
			labIdLabMap, &testDocumentMap)
	} else if cErr != nil && cErr.StatusCode == http.StatusNotFound || task.Id == 0 {
		task, cErr = eventProcessor.createData(ctx, omsApprovedEvent, testIdInvestigationResultsMap,
			masterInvestigationIdInvestigationResultMetadataMap, investigationCodeInvestigationDataMap,
			testIdInitialTestDetailsMap, masterTestDepartmentMap, investigationCodeMedicalRemarkMap,
			investigationCodeTechnicianRemarkMap, []string{}, userIdToAttuneUserIdMap, labIdLabMap, &testDocumentMap)
	} else if cErr != nil {
		utils.AddLog(ctx, constants.ERROR_LEVEL, utils.GetCurrentFunctionName(), nil, errors.New(cErr.Message))
	}

	if cErr != nil {
		return task, testDocumentMap, cErr
	}

	return task, testDocumentMap, nil
}

func (eventProcessor *EventProcessor) processLisCompletedEvent(ctx context.Context, omsCompletedEvent structures.OmsLisEvent,
	masterTestDepartmentMap map[string]string, usersList []models.User, systemPatientId string, sampleLab structures.Lab,
	labIdLabMap map[uint]structures.Lab) (
	models.Task, map[string][]structures.TestDocumentInfoResponse, *structures.CommonError) {
	var cErr *structures.CommonError
	task := models.Task{}
	masterInvestigationDetailsMap, err := eventProcessor.fetchInvestigationDetailsFromCds(ctx, omsCompletedEvent.Tests,
		sampleLab.City, sampleLab.Id, omsCompletedEvent.Patient)
	if err != nil {
		return task, nil, &structures.CommonError{
			Message:    err.Error(),
			StatusCode: http.StatusInternalServerError,
		}
	}

	masterInvestigationIdToCdsInvestigationMap := make(map[uint]structures.Investigation)
	for _, inv := range masterInvestigationDetailsMap {
		masterInvestigationIdToCdsInvestigationMap[inv.InvestigationId] = inv
	}
	userIdToAttuneUserIdMap := map[uint]int{}
	for _, user := range usersList {
		attuneUserId := utils.ConvertStringToInt(user.AttuneUserId)
		userIdToAttuneUserIdMap[user.Id] = attuneUserId
	}
	testIdInvestigationResultsMap, masterInvestigationIdInvestigationResultMetadataMap, testIdInitialTestDetailsMap,
		investigationCodeInvestigationDataMap, investigationCodeMedicalRemarkMap, investigationCodeTechnicianRemarkMap,
		qcFailedOmsTestIds, testDocumentMap := eventProcessor.createFlattenedInvestigationResultsAndInitialTestDetailsMap(ctx,
		constants.OmsCompletedEvent, omsCompletedEvent.Order.CityCode, omsCompletedEvent.Tests,
		masterInvestigationDetailsMap, usersList, systemPatientId, sampleLab.Id)
	task, cErr = eventProcessor.TaskService.GetTaskByOmsOrderId(omsCompletedEvent.Order.AlnumOrderId)
	if task.Id != 0 {
		task, cErr = eventProcessor.createUpdateData(ctx, omsCompletedEvent, task,
			testIdInvestigationResultsMap, masterInvestigationIdInvestigationResultMetadataMap,
			investigationCodeInvestigationDataMap, testIdInitialTestDetailsMap, masterTestDepartmentMap,
			investigationCodeMedicalRemarkMap, investigationCodeTechnicianRemarkMap, qcFailedOmsTestIds,
			userIdToAttuneUserIdMap, labIdLabMap, &testDocumentMap)
	} else if cErr != nil && cErr.StatusCode == http.StatusNotFound || task.Id == 0 {
		task, cErr = eventProcessor.createData(ctx, omsCompletedEvent, testIdInvestigationResultsMap,
			masterInvestigationIdInvestigationResultMetadataMap, investigationCodeInvestigationDataMap,
			testIdInitialTestDetailsMap, masterTestDepartmentMap, investigationCodeMedicalRemarkMap,
			investigationCodeTechnicianRemarkMap, qcFailedOmsTestIds, userIdToAttuneUserIdMap, labIdLabMap, &testDocumentMap)
	} else if cErr != nil {
		utils.AddLog(ctx, constants.ERROR_LEVEL, utils.GetCurrentFunctionName(), nil, errors.New(cErr.Message))
	}

	if cErr != nil {
		return task, testDocumentMap, cErr
	}

	if constants.EnableAutoApprovalSlackAlerts {
		uniqueMasterInvIds := []uint{}
		invResults := []models.InvestigationResult{}
		testIdTestNameMap := make(map[string]string)
		masterInvIdOmsTestIdsMap := make(map[uint][]string)
		for testId, investigations := range testIdInvestigationResultsMap {
			for _, investigation := range investigations {
				if !utils.SliceContainsUint(uniqueMasterInvIds, investigation.MasterInvestigationId) {
					invResults = append(invResults, investigation)
					uniqueMasterInvIds = append(uniqueMasterInvIds, investigation.MasterInvestigationId)
				}
				if _, exists := masterInvIdOmsTestIdsMap[investigation.MasterInvestigationId]; !exists {
					masterInvIdOmsTestIdsMap[investigation.MasterInvestigationId] = []string{}
				}
				masterInvIdOmsTestIdsMap[investigation.MasterInvestigationId] =
					append(masterInvIdOmsTestIdsMap[investigation.MasterInvestigationId], testId)
			}
		}

		for _, test := range omsCompletedEvent.Tests.TestDetails {
			testIdTestNameMap[test.TestId] = test.TestName
		}

		eventProcessor.findAndNotifyAutoApprovalFailedTests(ctx, invResults, omsCompletedEvent,
			masterInvIdOmsTestIdsMap, testIdTestNameMap)
	}
	return task, testDocumentMap, nil
}

func (eventProcessor *EventProcessor) processLisRerunEventTask(ctx context.Context,
	omsRerunEvent structures.OmsRerunEvent) error {
	var err error

	task, taskMetadata, testDetails, testDetailsMetadata, investigationResults, cErr :=
		eventProcessor.fetchCurrentDataForRerunEvent(omsRerunEvent)
	if cErr != nil {
		err = errors.New(cErr.Message)
		utils.AddLog(ctx, constants.ERROR_LEVEL, utils.GetCurrentFunctionName(), nil, err)
		return err
	}

	rerunDetails, investigationIds := createFlattenedRerunInvestigationResults(ctx, investigationResults,
		omsRerunEvent.Tests)

	omsTestIdToTestDetailMap := make(map[string]models.TestDetail)
	testDetailsToBeUpdated, testDetailsIdsToBeUpdated, omsRerunTestIds := []models.TestDetail{}, []uint{}, []string{}
	for _, testDetail := range testDetails {
		omsTestIdToTestDetailMap[testDetail.CentralOmsTestId] = testDetail
	}
	for _, test := range omsRerunEvent.Tests.TestDetails {
		if testDetail, ok := omsTestIdToTestDetailMap[test.TestId]; ok {
			if testDetail.Status == constants.TEST_STATUS_RERUN_REQUESTED {
				continue
			}
			testDetail.Status = constants.TEST_STATUS_RERUN_REQUESTED
			testDetailsToBeUpdated = append(testDetailsToBeUpdated, testDetail)
			testDetailsIdsToBeUpdated = append(testDetailsIdsToBeUpdated, testDetail.Id)
			omsRerunTestIds = append(omsRerunTestIds, testDetail.CentralOmsTestId)
			omsTestIdToTestDetailMap[test.TestId] = testDetail
		}
	}

	filteredRerunDetails := []models.RerunInvestigationResult{}

	for _, rerunDetail := range rerunDetails {
		if utils.SliceContainsUint(testDetailsIdsToBeUpdated, rerunDetail.TestDetailsId) {
			filteredRerunDetails = append(filteredRerunDetails, rerunDetail)
		}
	}

	investigationResultsToBeUpdated := []models.InvestigationResult{}
	for _, investigation := range investigationResults {
		if utils.SliceContainsUint(investigationIds, investigation.Id) {
			investigation.InvestigationStatus = constants.INVESTIGATION_STATUS_RERUN
			investigationResultsToBeUpdated = append(investigationResultsToBeUpdated, investigation)
		}
	}

	testDetails = []models.TestDetail{}
	for _, testDetail := range omsTestIdToTestDetailMap {
		testDetails = append(testDetails, testDetail)
	}

	task, taskMetadata = updateTaskBasedOnRerunEventChanges(task, taskMetadata, testDetails, testDetailsMetadata)

	err = eventProcessor.Db.Transaction(func(tx *gorm.DB) error {
		_, cErr = eventProcessor.TaskService.UpdateTaskWithTx(tx, task)
		if cErr != nil {
			return errors.New(cErr.Message)
		}

		_, cErr = eventProcessor.TaskService.UpdateTaskMetadataWithTx(tx, taskMetadata)
		if cErr != nil {
			return errors.New(cErr.Message)
		}

		_, cErr = eventProcessor.TestDetailService.UpdateTestDetailsWithTx(tx, testDetailsToBeUpdated)
		if cErr != nil {
			return errors.New(cErr.Message)
		}

		_, cErr = eventProcessor.RerunService.CreateRerunInvestigationResultsWithTx(tx, filteredRerunDetails)
		if cErr != nil {
			return errors.New(cErr.Message)
		}

		_, cErr = eventProcessor.InvestigationResultsService.UpdateInvestigationResultsWithTx(
			tx, investigationResultsToBeUpdated)
		if cErr != nil {
			return errors.New(cErr.Message)
		}

		return nil
	})
	if err != nil {
		utils.AddLog(ctx, constants.ERROR_LEVEL, utils.GetCurrentFunctionName(), nil, err)
		return err
	}

	if len(omsRerunTestIds) != 0 {
		go eventProcessor.EtsService.GetAndPublishEtsTestEventForLisWebhook(context.Background(), omsRerunTestIds, constants.AttuneTestStatusRerun)
	}
	return nil
}

func (eventProcessor *EventProcessor) fetchCurrentDataForRerunEvent(
	omsRerunEvent structures.OmsRerunEvent) (task models.Task, taskMetadata models.TaskMetadata,
	testDetails []models.TestDetail, testDetailsMetadata []models.TestDetailsMetadata,
	investigations []models.InvestigationResult, cErr *structures.CommonError) {

	taskMetadata = models.TaskMetadata{}
	testDetails = []models.TestDetail{}
	testDetailsMetadata = []models.TestDetailsMetadata{}
	investigations = []models.InvestigationResult{}
	task, cErr = eventProcessor.TaskService.GetTaskByOmsOrderId(omsRerunEvent.AlnumOrderId)
	if cErr != nil || task.Id == 0 {
		return
	}

	taskMetadata, cErr = eventProcessor.TaskService.GetTaskMetadataByTaskId(task.Id)
	if cErr != nil {
		return
	}

	testDetails, cErr = eventProcessor.TestDetailService.GetTestDetailsByOmsOrderId(task.OmsOrderId)
	if cErr != nil {
		return
	}

	testDetailsIds := []uint{}
	for index := range testDetails {
		testDetailsIds = append(testDetailsIds, testDetails[index].Id)
	}

	testDetailsMetadata, cErr = eventProcessor.TestDetailService.GetTestDetailsMetadataByTestDetailIds(testDetailsIds)
	if cErr != nil {
		return
	}

	investigations, cErr = eventProcessor.InvestigationResultsService.GetInvestigationResultsByTestDetailsIds(testDetailsIds)
	if cErr != nil {
		return
	}

	return
}

func createFlattenedRerunInvestigationResults(ctx context.Context, investigations []models.InvestigationResult,
	omsTestDetails structures.OmsTestDetails) (rerunDetails []models.RerunInvestigationResult, investigationIds []uint) {

	investigationCodeToInvestigationMap := make(map[string]models.InvestigationResult)
	investigationIdToInvestigationCodeMap := make(map[uint]string)
	for _, investigation := range investigations {
		investigationCodeToInvestigationMap[investigation.LisCode] = investigation
		investigationIdToInvestigationCodeMap[investigation.Id] = investigation.LisCode
	}

	rerunDetails = []models.RerunInvestigationResult{}

	for _, lisTestInfo := range omsTestDetails.OrderInfo {
		lisOrderInfo := lisTestInfo.MetaData

		switch lisOrderInfo.TestType {
		case constants.InvestigationShortHand:
			rerunDetail, investigationId := createRerunInvestigationsForInvestigation(lisOrderInfo,
				investigationCodeToInvestigationMap)
			if investigationId != 0 {
				investigationIds = append(investigationIds, investigationId)
			}
			rerunDetails = append(rerunDetails, rerunDetail)
		case constants.GroupShortHand:
			panelRerunDetails, panleInvestigationIds := createRerunInvestigationsForPanel(ctx, lisOrderInfo,
				investigationCodeToInvestigationMap)
			rerunDetails = append(rerunDetails, panelRerunDetails...)
			investigationIds = append(investigationIds, panleInvestigationIds...)
		}
	}

	return
}

func createRerunInvestigationsForInvestigation(lisOrderInfo structures.AttuneOrderInfo,
	investigationCodeToInvestigationMap map[string]models.InvestigationResult) (
	rerunDetail models.RerunInvestigationResult, investigationId uint) {
	investigation := investigationCodeToInvestigationMap[lisOrderInfo.TestCode]
	if !utils.SliceContainsString(constants.INVESTIGATION_STATUSES_RERUN, investigation.InvestigationStatus) {
		investigationId = investigation.Id
	}
	rerunTriggerTime := utils.GetTimeFromString(lisOrderInfo.RerunTime)
	if rerunTriggerTime == nil {
		rerunTriggerTime = utils.GetCurrentTime()
	}

	rerunDetail = models.RerunInvestigationResult{
		TestDetailsId:            investigation.TestDetailsId,
		MasterInvestigationId:    investigation.MasterInvestigationId,
		InvestigationName:        investigation.InvestigationName,
		InvestigationValue:       investigation.InvestigationValue,
		DeviceValue:              investigation.DeviceValue,
		ResultRepresentationType: investigation.ResultRepresentationType,
		LisCode:                  investigation.LisCode,
		RerunReason:              lisOrderInfo.RerunReason,
		RerunRemarks:             lisOrderInfo.RerunRemarks,
		RerunTriggeredAt:         rerunTriggerTime,
		RerunTriggeredBy:         constants.LisSystemId,
	}
	rerunDetail.CreatedBy = constants.CitadelSystemId
	rerunDetail.UpdatedBy = constants.CitadelSystemId

	return
}

func createRerunInvestigationsForPanel(ctx context.Context, lisOrderInfo structures.AttuneOrderInfo,
	investigationCodeToInvestigationMap map[string]models.InvestigationResult) (
	rerunDetails []models.RerunInvestigationResult, investigationIds []uint) {

	rerunDetails, investigationIds, queue, index :=
		[]models.RerunInvestigationResult{}, []uint{}, []interface{}{}, 0

	queue = append(queue, lisOrderInfo.OrderContentListInfo)

	for index < len(queue) {
		currentNode := queue[index]
		index++
		marshalledNode, err := json.Marshal(currentNode)
		if err != nil {
			utils.AddLog(ctx, constants.DEBUG_LEVEL, utils.GetCurrentFunctionName(), nil, err)
			continue
		}

		attuneOrderContentListInfo := []structures.AttuneOrderContentListInfoEvent{}
		err = json.Unmarshal(marshalledNode, &attuneOrderContentListInfo)
		if err != nil {
			utils.AddLog(ctx, constants.DEBUG_LEVEL, utils.GetCurrentFunctionName(), nil, err)
			continue
		}

		for _, attuneOrderContentListInfo := range attuneOrderContentListInfo {
			switch attuneOrderContentListInfo.TestType {
			case constants.InvestigationShortHand:
				if attuneOrderContentListInfo.TestStatus != constants.ATTUNE_TEST_STATUS_RERUN {
					continue
				}
				investigation := investigationCodeToInvestigationMap[attuneOrderContentListInfo.TestCode]
				rerunTriggerTime := utils.GetTimeFromString(attuneOrderContentListInfo.RerunTime)
				if rerunTriggerTime == nil {
					rerunTriggerTime = utils.GetCurrentTime()
				}
				rerunDetail := models.RerunInvestigationResult{
					TestDetailsId:            investigation.TestDetailsId,
					MasterInvestigationId:    investigation.MasterInvestigationId,
					InvestigationName:        investigation.InvestigationName,
					InvestigationValue:       investigation.InvestigationValue,
					DeviceValue:              investigation.DeviceValue,
					ResultRepresentationType: investigation.ResultRepresentationType,
					LisCode:                  investigation.LisCode,
					RerunReason:              attuneOrderContentListInfo.RerunReason,
					RerunRemarks:             attuneOrderContentListInfo.RerunRemarks,
					RerunTriggeredAt:         rerunTriggerTime,
					RerunTriggeredBy:         constants.LisSystemId,
				}
				rerunDetail.CreatedBy = constants.CitadelSystemId
				rerunDetail.UpdatedBy = constants.CitadelSystemId
				rerunDetails = append(rerunDetails, rerunDetail)
				if !utils.SliceContainsString(constants.INVESTIGATION_STATUSES_RERUN, investigation.InvestigationStatus) {
					investigationIds = append(investigationIds, investigation.Id)
				}
			case constants.GroupShortHand:
				queue = append(queue, attuneOrderContentListInfo.ParameterListInfo)
			}
		}
	}

	return
}

func updateTaskBasedOnRerunEventChanges(task models.Task, taskMetadata models.TaskMetadata,
	allTestDetails []models.TestDetail, testDetailsMetadata []models.TestDetailsMetadata) (
	models.Task, models.TaskMetadata) {
	if len(allTestDetails) == 0 {
		return task, taskMetadata
	}

	testDetailsMap := make(map[uint]models.TestDetail)
	for _, testDetail := range allTestDetails {
		testDetailsMap[testDetail.Id] = testDetail
	}

	doctorTat := getDoctorTatForTask(allTestDetails)
	if !doctorTat.IsZero() {
		task.DoctorTat = doctorTat
	}

	taskMetadata.IsCritical = getTaskCriticality(testDetailsMetadata, testDetailsMap)

	return task, taskMetadata
}

func (eventProcessor *EventProcessor) fetchAttuneEventWithRetry(ctx context.Context, reportPdfFormat string,
	lisOrderUpdateEvent structures.LisOrderUpdateDetails, sampleLabId uint) structures.LisOrderUpdateDetails {
	lisOrderUpdateEvent.ReportPdfFormat = reportPdfFormat
	for attempt := range constants.LisMissingPdfMaxRetries {
		attuneResponse, _ := eventProcessor.AttuneClient.GetPatientVisitDetailsbyVisitNo(ctx, lisOrderUpdateEvent.LisVisitId,
			reportPdfFormat, sampleLabId)
		if attuneResponse.ResultAsPdf == "" {
			sleepTime := utils.GetExponentialBackoff(attempt, constants.LisMissingPdfBackoffTime)
			time.Sleep(sleepTime)
			utils.AddLog(ctx, constants.DEBUG_LEVEL, utils.GetCurrentFunctionName(),
				map[string]interface{}{
					"attempt":           attempt + 1,
					"visit_id":          lisOrderUpdateEvent.LisVisitId,
					"report_pdf_format": lisOrderUpdateEvent.ReportPdfFormat,
				},
				errors.New(constants.ERROR_REPORT_PDF_NOT_FOUND))
			continue
		}
		return GetLisVisitDataFromAttuneOrderResponse(attuneResponse, reportPdfFormat)
	}
	return lisOrderUpdateEvent
}

func (eventProcessor *EventProcessor) SaveReportFromLis(orderDetails models.OrderDetails,
	patientDetails models.PatientDetail, stationaryLisOrderUpdateDetailsEvent structures.LisOrderUpdateDetails,
	testDetails []structures.TestDetailsForLisEvent, sampleLabId uint) {
	ctx := context.Background()
	cobrandedFilePath := ""
	omsTestIds := []string{}
	for _, testDetail := range testDetails {
		if _, keyExists := stationaryLisOrderUpdateDetailsEvent.OrderInfo[constants.ATTUNE_TEST_STATUS_APPROVED][testDetail.TestCode]; keyExists {
			omsTestIds = append(omsTestIds, testDetail.TestId)
		}
	}

	attuneResponse, cErr := eventProcessor.AttuneClient.GetPatientVisitDetailsbyVisitNo(ctx,
		stationaryLisOrderUpdateDetailsEvent.LisVisitId, constants.AttuneReportWithoutStationery, sampleLabId)
	if cErr != nil {
		utils.AddLog(ctx, constants.ERROR_LEVEL, utils.GetCurrentFunctionName(), nil, errors.New(cErr.Message))
		return
	}
	nonStationaryLisOrderUpdateDetailsEvent := GetLisVisitDataFromAttuneOrderResponse(attuneResponse,
		constants.AttuneReportWithoutStationery)

	if stationaryLisOrderUpdateDetailsEvent.PdfResult == "" {
		stationaryLisOrderUpdateDetailsEvent = eventProcessor.fetchAttuneEventWithRetry(ctx,
			constants.AttuneReportWithStationery, stationaryLisOrderUpdateDetailsEvent, sampleLabId)
	}
	if nonStationaryLisOrderUpdateDetailsEvent.PdfResult == "" {
		nonStationaryLisOrderUpdateDetailsEvent = eventProcessor.fetchAttuneEventWithRetry(ctx,
			constants.AttuneReportWithoutStationery, nonStationaryLisOrderUpdateDetailsEvent, sampleLabId)
	}

	if stationaryLisOrderUpdateDetailsEvent.PdfResult == "" || nonStationaryLisOrderUpdateDetailsEvent.PdfResult == "" {
		utils.AddLog(ctx, constants.ERROR_LEVEL, constants.ERROR_REPORT_PDF_NOT_FOUND, map[string]interface{}{
			"order_id":                 orderDetails.OmsOrderId,
			"visit_id":                 stationaryLisOrderUpdateDetailsEvent.LisVisitId,
			"stationary_pdf_found":     stationaryLisOrderUpdateDetailsEvent.PdfResult != "",
			"non_stationary_pdf_found": nonStationaryLisOrderUpdateDetailsEvent.PdfResult != "",
		}, nil)
		return
	}

	brandedFilePath, _ := eventProcessor.UpsertReportFromBase64String(ctx, orderDetails.OmsOrderId, patientDetails.Name,
		stationaryLisOrderUpdateDetailsEvent.PdfResult, constants.AttuneReportWithStationery)
	nonbrandedFilePath, _ := eventProcessor.UpsertReportFromBase64String(ctx, orderDetails.OmsOrderId, patientDetails.Name,
		nonStationaryLisOrderUpdateDetailsEvent.PdfResult, constants.AttuneReportWithoutStationery)

	if orderDetails.PartnerId != 0 {
		partner, err := eventProcessor.PartnerApiClient.GetPartnerById(ctx, orderDetails.PartnerId)
		if err != nil {
			utils.AddLog(ctx, constants.ERROR_LEVEL, utils.GetCurrentFunctionName(), nil, err)
		}
		if partner.ReportFormat == constants.ReportFormatCoBranded {
			cobrandedReportUrl, err := eventProcessor.UploadCobrandedReport(ctx, nonbrandedFilePath, partner.CobrandedImageUrl)
			if err != nil {
				utils.AddLog(ctx, constants.ERROR_LEVEL, utils.GetCurrentFunctionName(), nil, err)
			}

			cobrandedFilePath, _ = eventProcessor.UpsertReportFromBase64String(ctx,
				orderDetails.OmsOrderId, patientDetails.Name, cobrandedReportUrl, constants.AttuneReportWithCobrandStationery)
		}
	}

	reportPdfEvent := structures.ReportPdfEvent{
		OrderID:     orderDetails.OmsOrderId,
		CityCode:    orderDetails.CityCode,
		TestIds:     omsTestIds,
		VisitId:     stationaryLisOrderUpdateDetailsEvent.LisVisitId,
		AttuneFiles: true,
	}

	reportPdfEvent.ReportPdfBrandedURL = brandedFilePath
	reportPdfEvent.ReportPdfHeaderlessURL = nonbrandedFilePath

	if cobrandedFilePath != "" {
		reportPdfEvent.ReportPdfCobrandedURL = cobrandedFilePath
	}

	utils.AddLog(ctx, constants.DEBUG_LEVEL, utils.GetCurrentFunctionName(), map[string]interface{}{
		"orderId":          orderDetails.OmsOrderId,
		"testIds":          reportPdfEvent.TestIds,
		"report_pdf_event": reportPdfEvent,
	}, nil)

	reportReadyEvent := structures.ReportReadyEvent{
		ReportPdfEvent:    reportPdfEvent,
		ServicingCityCode: orderDetails.CityCode,
	}

	messageBody, messageAttributes := eventProcessor.PubsubService.GetReportReadyEvent(ctx, reportReadyEvent)
	cErr = eventProcessor.SnsClient.PublishTo(ctx, messageBody, messageAttributes, fmt.Sprint(orderDetails.OmsOrderId),
		constants.ReportReadyTopicArn, fmt.Sprintf("%v_%v", orderDetails.OmsOrderId, utils.GetCurrentTimeInMilliseconds()))
	if cErr != nil {
		utils.AddLog(ctx, constants.ERROR_LEVEL, utils.GetCurrentFunctionName(), nil, errors.New(cErr.Message))
	}
}

func (eventProcessor *EventProcessor) GetOrderFileName(orderId string, patientName string, reportPdfFormat string) string {
	t := time.Now()
	timestamp := t.Format(constants.TimeStampLayout)
	extension := ".pdf"
	patientFirstName := strings.Split(strings.TrimSpace(patientName), " ")[0]
	fileName := fmt.Sprintf("%s_%s_%s_%s", patientFirstName, orderId, timestamp, reportPdfFormat)
	fileName = fmt.Sprintf("%s%s", fileName, extension)
	return fileName
}

func (eventProcessor *EventProcessor) UpsertReportFromBase64String(ctx context.Context, omsOrderId string,
	patientName, reportBase64, reportPdfFormat string) (string, string) {
	fileName := eventProcessor.GetOrderFileName(fmt.Sprint(omsOrderId), patientName, reportPdfFormat)
	uploadedFilePath, uploadFileReference := eventProcessor.S3Client.UploadFileFromString(reportBase64, fileName)
	if uploadedFilePath == "" || uploadFileReference == "" {
		eventProcessor.Sentry.LogError(ctx, "Uploading report to s3 failed", nil, map[string]interface{}{
			"OrderID": omsOrderId,
		})
		return "", ""
	}

	return uploadedFilePath, uploadFileReference
}

func (eventProcessor *EventProcessor) UploadCobrandedReport(ctx context.Context, reportUrl, cobrandedImageUrl string) (
	string, error) {
	tokenizePublicReportFileURL, err := eventProcessor.S3wrapperClient.GetTokenizeOrderFilePublicUrl(ctx, reportUrl)
	if err != nil {
		return "", err
	}
	tokenizePublicCobrandedURL, err := eventProcessor.S3wrapperClient.GetTokenizeOrderFilePublicUrl(ctx, cobrandedImageUrl)
	if err != nil {
		return "", err
	}
	reqBody := structures.CobrandingRequest{
		ReportUrl:         tokenizePublicReportFileURL,
		CobrandedImageUrl: tokenizePublicCobrandedURL,
		HeaderLength:      100, // hardcoded header height for now
	}
	cobrandingResponse, err := eventProcessor.ReportRebrandingClient.AttachCobrandedImage(ctx, reqBody)
	if err != nil {
		return "", err
	}
	return cobrandingResponse.ReportData, nil
}

func (eventProcessor *EventProcessor) findAndNotifyAutoApprovalFailedTests(ctx context.Context,
	invResults []models.InvestigationResult,
	omsCompletedEvent structures.OmsLisEvent,
	invIdToTestIdsMap map[uint][]string,
	testIdToTestNameMap map[string]string,
) {
	deltaCheckFailedInvs := []models.InvestigationResult{}
	refRangeFailedInvs := []models.InvestigationResult{}
	for _, invResult := range invResults {
		switch invResult.AutoApprovalFailureReason {
		case constants.AUTO_APPROVAL_FAIL_REASON_PAST_RECORD:
			deltaCheckFailedInvs = append(deltaCheckFailedInvs, invResult)
		case constants.AUTO_APPROVAL_FAIL_REASON_REF_RANGE:
			refRangeFailedInvs = append(refRangeFailedInvs, invResult)
		}
	}
	eventProcessor.notifyAutoApprovalFailure(ctx, omsCompletedEvent, deltaCheckFailedInvs,
		constants.AUTO_APPROVAL_FAIL_REASON_PAST_RECORD, invIdToTestIdsMap, testIdToTestNameMap)
	eventProcessor.notifyAutoApprovalFailure(ctx, omsCompletedEvent, refRangeFailedInvs,
		constants.AUTO_APPROVAL_FAIL_REASON_REF_RANGE, invIdToTestIdsMap, testIdToTestNameMap)
}

func (eventProcessor *EventProcessor) notifyAutoApprovalFailure(ctx context.Context,
	omsCompletedEvent structures.OmsLisEvent,
	failedInvs []models.InvestigationResult,
	failureReason string,
	invIdToTestIdsMap map[uint][]string,
	testIdToTestNameMap map[string]string,
) {
	if len(failedInvs) > 0 {
		investigationNames := []string{}
		testNames := []string{}
		uniqueTestIds := []string{}
		visitIds := []string{}
		for _, visit := range omsCompletedEvent.Visits {
			if visit.Id != "" && !utils.SliceContainsString(visitIds, visit.Id) {
				visitIds = append(visitIds, visit.Id)
			}
		}
		for _, inv := range failedInvs {
			investigationNames = append(investigationNames, inv.InvestigationName)
			for _, testId := range invIdToTestIdsMap[inv.MasterInvestigationId] {
				if !utils.SliceContainsString(uniqueTestIds, testId) {
					testNames = append(testNames, testIdToTestNameMap[testId])
					uniqueTestIds = append(uniqueTestIds, testId)
				}
			}
		}
		investigationNamesSeparated := strings.Join(investigationNames, ", ")
		testNamesSeparated := strings.Join(testNames, ", ")
		visitIdsSeparated := strings.Join(visitIds, ", ")
		autoApprovalSlackMessageAttributes := structures.AutoApprovalSlackMessageAttributes{
			PatientName:        omsCompletedEvent.Patient.Name,
			InvestigationsName: investigationNamesSeparated,
			OrderId:            omsCompletedEvent.Order.AlnumOrderId,
			RequestID:          omsCompletedEvent.Request.AlnumRequestId,
			VisitID:            visitIdsSeparated,
			CityCode:           omsCompletedEvent.Order.CityCode,
			TestNames:          testNamesSeparated,
		}

		msgBlocks := utils.GetSlackMessageBlocksForAutoApprovalFailure(autoApprovalSlackMessageAttributes, failureReason)
		err := eventProcessor.SlackClient.SendToSlackDirectly(ctx, constants.SlackAutoApprovalFailureChannel, msgBlocks)
		if err != nil {
			utils.AddLog(ctx, constants.ERROR_LEVEL, utils.GetCurrentFunctionName(), nil, err)
		}
	}
}
