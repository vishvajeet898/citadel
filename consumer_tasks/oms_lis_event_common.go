package consumerTasks

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"

	"gorm.io/gorm"

	"github.com/Orange-Health/citadel/common/constants"
	"github.com/Orange-Health/citadel/common/structures"
	"github.com/Orange-Health/citadel/common/utils"
	"github.com/Orange-Health/citadel/models"
)

func isInvestigationCritical(abnormality string) bool {
	return utils.SliceContainsString(constants.OhCriticalityStringSlice, abnormality)
}

func getInitialTestStatusForTestDetails(isAutoApproved bool, eventType string, areAttuneStatusesApproved bool) string {
	if isAutoApproved || (eventType == constants.OmsApprovedEvent && areAttuneStatusesApproved) {
		return constants.TEST_STATUS_APPROVE
	}
	return constants.TEST_STATUS_RESULT_SAVED
}

func getTestApprovalSource(eventType string, isTestAutoApproved bool) string {
	if eventType == constants.OmsApprovedEvent {
		return constants.TEST_APPROVAL_SOURCE_ATTUNE
	}

	if isTestAutoApproved {
		return constants.TEST_APPROVAL_SOURCE_OH
	}

	return ""
}

func getStatusForTestDetails(testId string, testIdInitialTestDetailsMap map[string]structures.InitialTestDetails,
	qcFailedTestIds []string) string {
	if utils.SliceContainsString(qcFailedTestIds, testId) {
		return constants.TEST_STATUS_RERUN_REQUESTED
	}
	status := constants.TEST_STATUS_RESULT_PENDING
	if initialTestDetails, ok := testIdInitialTestDetailsMap[testId]; ok {
		status = initialTestDetails.TestStatus
	}
	return status
}

func getDoctorTatForTestDetails(testId string,
	testIdInitialTestDetailsMap map[string]structures.InitialTestDetails) *time.Time {
	var doctorTat *time.Time
	if initialTestDetails, ok := testIdInitialTestDetailsMap[testId]; ok {
		doctorTat = initialTestDetails.DoctorTat
	}
	return doctorTat
}

func getAutoApprovedStatusForTestDetails(testId string,
	testIdInitialTestDetailsMap map[string]structures.InitialTestDetails) bool {
	if initialTestDetails, ok := testIdInitialTestDetailsMap[testId]; ok {
		return initialTestDetails.IsAutoApproved
	}
	return false
}

func getDoctorTatForTask(testDetails []models.TestDetail) *time.Time {
	minDoctorTat := &time.Time{}
	testDetailsStatuses := []string{
		constants.TEST_STATUS_RESULT_SAVED,
		constants.TEST_STATUS_RERUN_RESULT_SAVED,
		constants.TEST_STATUS_WITHHELD,
		constants.TEST_STATUS_CO_AUTHORIZE,
	}
	for _, testDetail := range testDetails {
		if testDetail.DoctorTat != nil && utils.SliceContainsString(testDetailsStatuses, testDetail.Status) {
			if minDoctorTat.IsZero() || testDetail.DoctorTat.Before(*minDoctorTat) {
				minDoctorTat = testDetail.DoctorTat
			}
		}
	}
	return minDoctorTat
}

func getTaskCriticality(testDetailsMetadata []models.TestDetailsMetadata,
	testDetailsMap map[uint]models.TestDetail) bool {
	isCritical := false
	testStatuses := []string{
		constants.TEST_STATUS_RESULT_SAVED,
		constants.TEST_STATUS_RERUN_RESULT_SAVED,
		constants.TEST_STATUS_WITHHELD,
		constants.TEST_STATUS_CO_AUTHORIZE,
	}
	for _, testDetailMetadata := range testDetailsMetadata {
		testDetail := testDetailsMap[testDetailMetadata.TestDetailsId]
		if utils.SliceContainsString(testStatuses, testDetail.Status) && testDetailMetadata.IsCritical {
			isCritical = true
			break
		}
	}

	return isCritical
}

func getIsAbnormalFlagForInvestigation(ohAbnormality string) bool {
	return utils.SliceContainsString(constants.OhAbnormalityStringSlice, ohAbnormality)
}

func getReferenceRangeText(referenceRange structures.MasterReferenceRange, defaultReferenceRange string) string {
	// Currently this is a patch, where the ReferenceRangeText is being fetched from NormalRange Only
	if referenceRange.NormalRange.ReferenceRangeText != "" {
		return referenceRange.NormalRange.ReferenceRangeText
	}

	//	adding this as a fallback, if the reference range text is not found in the normal range
	return defaultReferenceRange
}

func getMethodTypeForInvestigation(deviceId, methodName string) string {
	if strings.ToLower(strings.TrimSpace(methodName)) == constants.METHOD_TYPE_CALCULATED {
		return constants.METHOD_TYPE_CALCULATED
	}
	if deviceId != "" {
		return constants.METHOD_TYPE_DEVICE_MEASURED
	}
	return constants.METHOD_TYPE_MANUAL
}

func fetchInvestigationCodesFromAttunePayload(tests structures.OmsTestDetails) []string {
	investigationCodes := []string{}

	queue, index := []interface{}{}, 0
	for _, lisTestUpdateInfo := range tests.OrderInfo {
		if lisTestUpdateInfo.MetaData.TestType == constants.InvestigationShortHand {
			investigationCodes = append(investigationCodes, lisTestUpdateInfo.TestCode)
		} else {
			for _, orderContentListInfo := range lisTestUpdateInfo.MetaData.OrderContentListInfo {
				queue = append(queue, orderContentListInfo)
			}
		}
	}

	for index < len(queue) {
		currentNode := queue[index]
		index += 1
		marshalledNode, err := json.Marshal(currentNode)
		if err != nil {
			continue
		}
		attuneOrderContentListInfo := structures.AttuneOrderContentListInfoEvent{}
		err = json.Unmarshal(marshalledNode, &attuneOrderContentListInfo)
		if err != nil {
			continue
		}
		if attuneOrderContentListInfo.TestType == constants.InvestigationShortHand {
			investigationCodes = append(investigationCodes, attuneOrderContentListInfo.TestCode)
		} else {
			queue = append(queue, attuneOrderContentListInfo.ParameterListInfo...)
		}
	}

	return investigationCodes
}

func createInitialRemarkModel(remarkType, remarkDescription string) models.Remark {
	remark := models.Remark{
		Description: remarkDescription,
		RemarkType:  remarkType,
		RemarkBy:    constants.LisSystemId,
	}
	remark.CreatedBy = constants.CitadelSystemId
	remark.UpdatedBy = constants.CitadelSystemId

	return remark
}

func createPatientDetailsDto(omsLisEvent structures.OmsLisEvent) models.PatientDetail {
	omsPatientDetails := omsLisEvent.Patient
	patientExpectedDob := utils.GetDobByYearsMonthsDays(int(omsPatientDetails.AgeYears),
		int(omsPatientDetails.AgeMonths), int(omsPatientDetails.AgeDays))

	patientDetails := models.PatientDetail{
		Name:            omsPatientDetails.Name,
		ExpectedDob:     &patientExpectedDob,
		Gender:          omsPatientDetails.Gender,
		Number:          omsPatientDetails.Number,
		SystemPatientId: omsPatientDetails.PatientId,
	}
	patientDetails.CreatedBy = constants.CitadelSystemId
	patientDetails.UpdatedBy = constants.CitadelSystemId

	return patientDetails
}

func createTaskDto(omsLisEvent structures.OmsLisEvent,
	patientDetails models.PatientDetail,
) models.Task {
	omsRequestDetails := omsLisEvent.Request
	omsOrderDetails := omsLisEvent.Order

	task := models.Task{
		OrderId:          omsOrderDetails.Id,
		RequestId:        omsRequestDetails.Id,
		OmsOrderId:       omsOrderDetails.AlnumOrderId,
		OmsRequestId:     omsRequestDetails.AlnumRequestId,
		LabId:            omsOrderDetails.LabId,
		CityCode:         omsOrderDetails.CityCode,
		Status:           constants.TASK_STATUS_PENDING,
		PreviousStatus:   constants.TASK_STATUS_PENDING,
		PatientDetailsId: patientDetails.Id,
		OrderType:        omsOrderDetails.OrderType,
		IsActive:         true,
	}
	task.CreatedBy = constants.CitadelSystemId
	task.UpdatedBy = constants.CitadelSystemId

	return task
}

func createTaskMetadataDto(omsLisEvent structures.OmsLisEvent,
	task models.Task,
) models.TaskMetadata {
	omsOrderDetails := omsLisEvent.Order

	taskMetadata := models.TaskMetadata{
		ContainsPackage: omsOrderDetails.ContainsPackage,
		TaskId:          task.Id,
		ContainsMorphle: false,
		DoctorName:      omsOrderDetails.DoctorName,
		DoctorNumber:    omsOrderDetails.DoctorNumber,
		DoctorNotes:     omsOrderDetails.DoctorNotes,
		PartnerName:     omsOrderDetails.PartnerName,
	}
	taskMetadata.CreatedBy = constants.CitadelSystemId
	taskMetadata.UpdatedBy = constants.CitadelSystemId

	return taskMetadata
}

func createTaskVisitMappingDto(omsLisEvent structures.OmsLisEvent,
	task models.Task,
) []models.TaskVisitMapping {
	omsVisitDetails := omsLisEvent.Visits
	taskVisitMappings := []models.TaskVisitMapping{}
	for _, visit := range omsVisitDetails {
		taskVisitMapping := models.TaskVisitMapping{
			VisitId: visit.Id,
			TaskId:  task.Id,
		}
		taskVisitMapping.CreatedBy = constants.CitadelSystemId
		taskVisitMapping.UpdatedBy = constants.CitadelSystemId
		taskVisitMappings = append(taskVisitMappings, taskVisitMapping)
	}

	return taskVisitMappings
}

func createTestDetailsDto(omsLisEvent structures.OmsLisEvent,
	task models.Task, testIdInitialTestDetailsMap map[string]structures.InitialTestDetails,
	masterTestDepartmentMap map[string]string, qcFailedTestIds []string,
) []models.TestDetail {
	omsTestDetails := omsLisEvent.Tests
	testDetails := []models.TestDetail{}
	for _, test := range omsTestDetails.TestDetails {
		testDetail := models.TestDetail{
			TaskId:           task.Id,
			OmsTestId:        utils.GetUintTestIdWithoutStringPart(test.TestId), // TODO
			CentralOmsTestId: test.TestId,
			TestName:         test.TestName,
			LisCode:          test.TestCode,
			MasterTestId:     test.MasterTestId,
			TestType:         test.TestType,
			Department:       masterTestDepartmentMap[fmt.Sprint(test.MasterTestId)],
			Status:           getStatusForTestDetails(test.TestId, testIdInitialTestDetailsMap, qcFailedTestIds),
			DoctorTat:        getDoctorTatForTestDetails(test.TestId, testIdInitialTestDetailsMap),
			IsAutoApproved:   getAutoApprovedStatusForTestDetails(test.TestId, testIdInitialTestDetailsMap),
			ApprovalSource:   testIdInitialTestDetailsMap[test.TestId].ApprovalSource,
		}
		testDetail.UpdatedBy = constants.CitadelSystemId
		testDetail.CreatedBy = constants.CitadelSystemId
		testDetails = append(testDetails, testDetail)
	}

	return testDetails
}

func createTestDetailsMetadataDto(omsLisEvent structures.OmsLisEvent,
	testDetails []models.TestDetail, testIdInitialTestDetailsMap map[string]structures.InitialTestDetails,
) []models.TestDetailsMetadata {
	omsTestDetails := omsLisEvent.Tests
	testIdBarcodeMap := map[string]string{}
	for _, test := range omsTestDetails.TestDetails {
		testIdBarcodeMap[test.TestId] = strings.Join(test.Barcodes, ",")
	}

	testDetailsMetadatas := []models.TestDetailsMetadata{}
	for _, testDetail := range testDetails {
		if barcodes, ok := testIdBarcodeMap[testDetail.CentralOmsTestId]; ok {
			testDetailsMetadata := models.TestDetailsMetadata{
				TestDetailsId: testDetail.Id,
				Barcodes:      barcodes,
				IsCritical:    testIdInitialTestDetailsMap[testDetail.CentralOmsTestId].IsCritical,
			}
			testDetailsMetadata.CreatedBy = constants.CitadelSystemId
			testDetailsMetadata.UpdatedBy = constants.CitadelSystemId
			testDetailsMetadatas = append(testDetailsMetadatas, testDetailsMetadata)
		}
	}

	return testDetailsMetadatas
}

func createInvestigationResultsDto(testDetails []models.TestDetail,
	testIdInvestigationResultsMap map[string][]models.InvestigationResult,
	qcFailedTestIds []string,
) []models.InvestigationResult {
	finalInvestigationResults := []models.InvestigationResult{}

	omsTestIdTestDetailsMap := map[string]models.TestDetail{}
	for _, testDetail := range testDetails {
		omsTestIdTestDetailsMap[testDetail.CentralOmsTestId] = testDetail
	}

	for testId, investigationResults := range testIdInvestigationResultsMap {
		for _, investigationResult := range investigationResults {
			investigationResult.TestDetailsId = omsTestIdTestDetailsMap[testId].Id
			if utils.SliceContainsString(qcFailedTestIds, testId) {
				investigationResult.InvestigationStatus = constants.INVESTIGATION_STATUS_RERUN
			}
			finalInvestigationResults = append(finalInvestigationResults, investigationResult)
		}
	}

	return finalInvestigationResults
}

func createInvestigationResultMetadataDto(investigation_results []models.InvestigationResult,
	masterInvestigationIdInvestigationResultMetadataMap map[uint]models.InvestigationResultMetadata,
) []models.InvestigationResultMetadata {
	finalInvestigationResultsMetadata := []models.InvestigationResultMetadata{}
	masterInvestigationIdInvestigationResultMap := map[uint]models.InvestigationResult{}

	for _, investigationResult := range investigation_results {
		masterInvestigationIdInvestigationResultMap[investigationResult.MasterInvestigationId] = investigationResult
	}

	for masterInvestigationId, investigationResultMetadata := range masterInvestigationIdInvestigationResultMetadataMap {
		investigationResultMetadata.InvestigationResultId = masterInvestigationIdInvestigationResultMap[masterInvestigationId].Id
		finalInvestigationResultsMetadata = append(finalInvestigationResultsMetadata, investigationResultMetadata)
	}
	return finalInvestigationResultsMetadata
}

func createInvestigationDataDto(investigationResults []models.InvestigationResult,
	investigationCodeInvestigationDataMap map[string]models.InvestigationData,
) []models.InvestigationData {
	investigationData := []models.InvestigationData{}

	for _, investigationResult := range investigationResults {
		if investigationDataDto, ok := investigationCodeInvestigationDataMap[investigationResult.LisCode]; ok {
			investigationDataDto.InvestigationResultId = investigationResult.Id
			investigationDataDto.CreatedBy = constants.CitadelSystemId
			investigationDataDto.UpdatedBy = constants.CitadelSystemId
			investigationData = append(investigationData, investigationDataDto)
		}
	}

	return investigationData
}

func createRemarksDto(investigationResults []models.InvestigationResult,
	investigationCodeMedicalRemarkMap map[string]models.Remark,
	investigationCodeTechnicianRemarkMap map[string]models.Remark,
) []models.Remark {
	remarks := []models.Remark{}

	for _, investigationResult := range investigationResults {
		if medicalRemark, ok := investigationCodeMedicalRemarkMap[investigationResult.LisCode]; ok {
			medicalRemark.InvestigationResultId = investigationResult.Id
			remarks = append(remarks, medicalRemark)
		}

		if technicianRemark, ok := investigationCodeTechnicianRemarkMap[investigationResult.LisCode]; ok {
			technicianRemark.InvestigationResultId = investigationResult.Id
			remarks = append(remarks, technicianRemark)
		}
	}

	return remarks
}

func updateTaskAndTaskMetadataBasedOnTestDetailsOnCreate(task models.Task, taskMetadata models.TaskMetadata,
	testDetails []models.TestDetail, testDetailsMetadata []models.TestDetailsMetadata, labIdLabMap map[uint]structures.Lab) (
	models.Task, models.TaskMetadata) {
	allInhouseTestsCompleted := true

	testCompleteStatuses := []string{
		constants.TEST_STATUS_APPROVE,
		constants.TEST_STATUS_COMPLETED_NOT_SENT,
		constants.TEST_STATUS_COMPLETED_SENT,
		constants.TEST_STATUS_SAMPLE_NOT_RECEIVED,
	}

	for _, testDetail := range testDetails {
		if !utils.SliceContainsString(testCompleteStatuses, testDetail.Status) &&
			utils.IsTestInhouse(testDetail.ProcessingLabId, task.LabId, labIdLabMap) {
			allInhouseTestsCompleted = false
			break
		}
	}

	for _, testDetailMetadata := range testDetailsMetadata {
		if testDetailMetadata.IsCritical {
			taskMetadata.IsCritical = true
			break
		}
	}

	if allInhouseTestsCompleted {
		task.Status = constants.TASK_STATUS_COMPLETED
		task.CompletedAt = utils.GetCurrentTime()
	}

	doctorTat := getDoctorTatForTask(testDetails)
	if !doctorTat.IsZero() {
		task.DoctorTat = doctorTat
	}

	return task, taskMetadata
}

func createUpdatePatientDetailsDto(omsLisEvent structures.OmsLisEvent,
	patientDetails models.PatientDetail,
) models.PatientDetail {
	omsPatientDetails := omsLisEvent.Patient
	if omsPatientDetails.Dob != nil {
		patientDetails.Dob = omsPatientDetails.Dob
		patientDetails.ExpectedDob = omsPatientDetails.Dob
	} else if omsPatientDetails.ExpectedDob != nil {
		patientDetails.ExpectedDob = omsPatientDetails.ExpectedDob
	} else {
		patientExpectedDob := utils.GetDobByYearsMonthsDays(int(omsPatientDetails.AgeYears),
			int(omsPatientDetails.AgeMonths), int(omsPatientDetails.AgeDays))
		patientDetails.ExpectedDob = &patientExpectedDob
	}

	patientDetails.Name = omsPatientDetails.Name
	patientDetails.Number = omsPatientDetails.Number
	patientDetails.SystemPatientId = omsPatientDetails.PatientId
	patientDetails.UpdatedBy = constants.CitadelSystemId

	return patientDetails
}

func createUpdateTaskMetadataDto(omsLisEvent structures.OmsLisEvent,
	taskMetadata models.TaskMetadata,
) models.TaskMetadata {
	omsOrderDetails := omsLisEvent.Order

	taskMetadata.ContainsPackage = omsOrderDetails.ContainsPackage
	taskMetadata.DoctorName = omsOrderDetails.DoctorName
	taskMetadata.DoctorNumber = omsOrderDetails.DoctorNumber
	taskMetadata.DoctorNotes = omsOrderDetails.DoctorNotes
	taskMetadata.PartnerName = omsOrderDetails.PartnerName
	taskMetadata.UpdatedBy = constants.CitadelSystemId

	return taskMetadata
}

func getCreateDeleteTaskVisits(omsLisEvent structures.OmsLisEvent,
	taskVisitMappings []models.TaskVisitMapping,
) ([]string, []string) {

	omsVisitDetails := omsLisEvent.Visits
	omsVisitIds := []string{}
	for _, visit := range omsVisitDetails {
		omsVisitIds = append(omsVisitIds, visit.Id)
	}

	taskVisitIds := []string{}
	for _, taskVisitMapping := range taskVisitMappings {
		taskVisitIds = append(taskVisitIds, taskVisitMapping.VisitId)
	}

	createVisitIds := utils.GetDifferenceBetweenStringSlices(omsVisitIds, taskVisitIds)
	deleteVisitIds := utils.GetDifferenceBetweenStringSlices(taskVisitIds, omsVisitIds)

	return createVisitIds, deleteVisitIds
}

func getUpdatedTestDetailsBasedOnStatus(newStatus, previousStatus string, testDetail models.TestDetail,
	testIdInitialTestDetailsMap map[string]structures.InitialTestDetails) models.TestDetail {
	if previousStatus == constants.TEST_STATUS_RESULT_PENDING ||
		previousStatus == constants.TEST_STATUS_RERUN_REQUESTED ||
		previousStatus == constants.TEST_STATUS_REJECTED {
		if previousStatus == constants.TEST_STATUS_RERUN_REQUESTED &&
			newStatus == constants.TEST_STATUS_RESULT_SAVED {
			testDetail.Status = constants.TEST_STATUS_RERUN_RESULT_SAVED
		} else if previousStatus == constants.TEST_STATUS_RERUN_REQUESTED &&
			newStatus == constants.TEST_STATUS_RESULT_PENDING {
			testDetail.Status = constants.TEST_STATUS_RERUN_REQUESTED
		} else {
			testDetail.Status = newStatus
			if testIdInitialTestDetailsMap[testDetail.CentralOmsTestId].IsAutoApproved {
				testDetail.IsAutoApproved = true
			}
		}
		testDetail.DoctorTat = getDoctorTatForTestDetails(testDetail.CentralOmsTestId, testIdInitialTestDetailsMap)
	} else if newStatus == constants.TEST_STATUS_APPROVE &&
		(previousStatus == constants.TEST_STATUS_RESULT_SAVED ||
			previousStatus == constants.TEST_STATUS_RERUN_RESULT_SAVED ||
			previousStatus == constants.TEST_STATUS_WITHHELD ||
			previousStatus == constants.TEST_STATUS_CO_AUTHORIZE) {
		if testIdInitialTestDetailsMap[testDetail.CentralOmsTestId].IsAutoApproved {
			testDetail.IsAutoApproved = true
		}
		testDetail.Status = constants.TEST_STATUS_APPROVE
	}
	return testDetail
}

func getCreateUpdateTestDetailsDto(omsLisEvent structures.OmsLisEvent,
	task models.Task, testDetails []models.TestDetail,
	testIdInitialTestDetailsMap map[string]structures.InitialTestDetails, masterTestDepartmentMap map[string]string,
	qcFailedTestIds []string,
) ([]models.TestDetail, []models.TestDetail) {
	toCreateTestDetails, toUpdateTestDetails := []models.TestDetail{}, []models.TestDetail{}
	omsTestDetails := omsLisEvent.Tests

	omsTestIds := []string{}
	for _, test := range omsTestDetails.TestDetails {
		omsTestIds = append(omsTestIds, test.TestId)
	}

	testDetailsOmsIds := []string{}
	for _, testDetail := range testDetails {
		testDetailsOmsIds = append(testDetailsOmsIds, testDetail.CentralOmsTestId)
	}

	toCreateTestIds := utils.GetDifferenceBetweenStringSlices(testDetailsOmsIds, testDetailsOmsIds)
	toUpdateTestIds := utils.GetCommonElementsBetweenStringSlices(testDetailsOmsIds, omsTestIds)

	for _, testDetail := range testDetails {
		if utils.SliceContainsString(toUpdateTestIds, testDetail.CentralOmsTestId) {
			status := getStatusForTestDetails(testDetail.CentralOmsTestId, testIdInitialTestDetailsMap, qcFailedTestIds)
			previousStatus := testDetail.Status
			testDetail = getUpdatedTestDetailsBasedOnStatus(status, previousStatus, testDetail, testIdInitialTestDetailsMap)
			if previousStatus != constants.TEST_STATUS_APPROVE {
				testDetail.ApprovalSource = testIdInitialTestDetailsMap[testDetail.CentralOmsTestId].ApprovalSource
			}
			if utils.SliceContainsString(qcFailedTestIds, testDetail.CentralOmsTestId) {
				testDetail.Status = constants.TEST_STATUS_RERUN_REQUESTED
			}
			toUpdateTestDetails = append(toUpdateTestDetails, testDetail)
		}
	}

	for _, testDetail := range omsTestDetails.TestDetails {
		if utils.SliceContainsString(toCreateTestIds, testDetail.TestId) {
			testDetail := models.TestDetail{
				TaskId:           task.Id,
				OmsTestId:        utils.GetUintTestIdWithoutStringPart(testDetail.TestId), // TODO
				CentralOmsTestId: testDetail.TestId,
				TestName:         testDetail.TestName,
				LisCode:          testDetail.TestCode,
				MasterTestId:     testDetail.MasterTestId,
				TestType:         testDetail.TestType,
				Department:       masterTestDepartmentMap[fmt.Sprint(testDetail.MasterTestId)],
				Status:           getStatusForTestDetails(testDetail.TestId, testIdInitialTestDetailsMap, qcFailedTestIds),
				DoctorTat:        getDoctorTatForTestDetails(testDetail.TestId, testIdInitialTestDetailsMap),
				IsAutoApproved:   getAutoApprovedStatusForTestDetails(testDetail.TestId, testIdInitialTestDetailsMap),
				ApprovalSource:   testIdInitialTestDetailsMap[testDetail.TestId].ApprovalSource,
			}
			testDetail.CreatedBy = constants.CitadelSystemId
			testDetail.UpdatedBy = constants.CitadelSystemId
			toCreateTestDetails = append(toCreateTestDetails, testDetail)
		}
	}

	return toCreateTestDetails, toUpdateTestDetails
}

func getCreateUpdateTestDetailsMetadataDtos(omsLisEvent structures.OmsLisEvent,
	testDetailsMetadata []models.TestDetailsMetadata, createTestDetails, updateTestDetails []models.TestDetail,
	testIdInitialTestDetailsMap map[string]structures.InitialTestDetails,
) ([]models.TestDetailsMetadata, []models.TestDetailsMetadata) {
	createTestDetailsMetadata, updateTestDetailsMetadata := []models.TestDetailsMetadata{}, []models.TestDetailsMetadata{}
	omsTestDetails := omsLisEvent.Tests

	testIdTestDetailsMetadataMap := map[uint]models.TestDetailsMetadata{}
	for _, testDetailMetadata := range testDetailsMetadata {
		testIdTestDetailsMetadataMap[testDetailMetadata.TestDetailsId] = testDetailMetadata
	}

	testIdBarcodeMap := map[string]string{}
	for _, testDetail := range omsTestDetails.TestDetails {
		testIdBarcodeMap[testDetail.TestId] = strings.Join(testDetail.Barcodes, ",")
	}

	for _, testDetail := range createTestDetails {
		if barcodes, ok := testIdBarcodeMap[testDetail.CentralOmsTestId]; ok {
			testDetailMetadata := models.TestDetailsMetadata{
				TestDetailsId: testDetail.Id,
				Barcodes:      barcodes,
				IsCritical:    testIdInitialTestDetailsMap[testDetail.CentralOmsTestId].IsCritical,
			}
			testDetailMetadata.CreatedBy = constants.CitadelSystemId
			testDetailMetadata.UpdatedBy = constants.CitadelSystemId
			createTestDetailsMetadata = append(createTestDetailsMetadata, testDetailMetadata)
		}
	}

	for _, testDetail := range updateTestDetails {
		if barcodes, ok := testIdBarcodeMap[testDetail.CentralOmsTestId]; ok {
			if testDetailsMetadata, ok := testIdTestDetailsMetadataMap[testDetail.Id]; ok {
				testDetailsMetadata.Barcodes = barcodes
				testDetailsMetadata.IsCritical = testIdInitialTestDetailsMap[testDetail.CentralOmsTestId].IsCritical
				updateTestDetailsMetadata = append(updateTestDetailsMetadata, testDetailsMetadata)
			}
		}
	}

	return createTestDetailsMetadata, updateTestDetailsMetadata
}

func getCreateUpdateInvestigationResultsDto(createTestDetails, updateTestDetails []models.TestDetail,
	testIdInvestigationResultsMap map[string][]models.InvestigationResult,
	currentInvestigationResults []models.InvestigationResult, qcFailedTestIds []string) (
	[]models.InvestigationResult, []models.InvestigationResult) {
	finalInvestigationResultsForCreation, finalInvestigationResultsForUpdation :=
		[]models.InvestigationResult{}, []models.InvestigationResult{}

	// Creation DTOs for new TestDetails
	omsTestIdTestDetailsMapForCreation := map[string]models.TestDetail{}
	for _, testDetail := range createTestDetails {
		omsTestIdTestDetailsMapForCreation[testDetail.CentralOmsTestId] = testDetail
	}

	for testId, investigationResults := range testIdInvestigationResultsMap {
		if _, ok := omsTestIdTestDetailsMapForCreation[testId]; ok {
			for _, investigationResult := range investigationResults {
				investigationResult.TestDetailsId = omsTestIdTestDetailsMapForCreation[testId].Id
				if utils.SliceContainsString(qcFailedTestIds, testId) {
					investigationResult.InvestigationStatus = constants.INVESTIGATION_STATUS_RERUN
				}
				finalInvestigationResultsForCreation = append(finalInvestigationResultsForCreation, investigationResult)
			}
		}
	}

	omsTestIdTestDetailsMapForUpdation := map[string]models.TestDetail{}
	for _, testDetail := range updateTestDetails {
		omsTestIdTestDetailsMapForUpdation[testDetail.CentralOmsTestId] = testDetail
	}

	testDetailsIdToOmsTestIdMap := map[uint]string{}
	for _, testDetail := range updateTestDetails {
		testDetailsIdToOmsTestIdMap[testDetail.Id] = testDetail.CentralOmsTestId
	}

	testDetailsIdToCurrentInvestigationResultsMap := map[uint][]models.InvestigationResult{}
	for _, investigationResult := range currentInvestigationResults {
		if _, ok := testDetailsIdToCurrentInvestigationResultsMap[investigationResult.TestDetailsId]; !ok {
			testDetailsIdToCurrentInvestigationResultsMap[investigationResult.TestDetailsId] = []models.InvestigationResult{}
		}
		testDetailsIdToCurrentInvestigationResultsMap[investigationResult.TestDetailsId] =
			append(testDetailsIdToCurrentInvestigationResultsMap[investigationResult.TestDetailsId], investigationResult)
	}

	// Creation DTOs for existing TestDetails
	testDetailsIdToNewInvestigationResultsMap := map[uint][]models.InvestigationResult{}
	for _, updateTestDetail := range updateTestDetails {
		testDetailsId := updateTestDetail.Id
		if _, ok := testDetailsIdToCurrentInvestigationResultsMap[testDetailsId]; !ok {
			omsTestId := testDetailsIdToOmsTestIdMap[testDetailsId]
			if _, ok := testIdInvestigationResultsMap[omsTestId]; ok {
				testDetailsIdToNewInvestigationResultsMap[testDetailsId] = testIdInvestigationResultsMap[omsTestId]
				for _, newInvestigationResult := range testIdInvestigationResultsMap[omsTestId] {
					newInvestigationResult.TestDetailsId = testDetailsId
					if utils.SliceContainsString(qcFailedTestIds, omsTestId) {
						newInvestigationResult.InvestigationStatus = constants.INVESTIGATION_STATUS_RERUN
					}
					finalInvestigationResultsForCreation = append(finalInvestigationResultsForCreation,
						newInvestigationResult)
				}
			}
		}
	}

	// Updation DTOs for exiting TestDetails
	for testDetailsId, currentInvestigationResults := range testDetailsIdToCurrentInvestigationResultsMap {
		omsTestId := testDetailsIdToOmsTestIdMap[testDetailsId]
		updatedInvestigationResults := testIdInvestigationResultsMap[omsTestId]
		finalInvestigationResultsForUpdation = append(finalInvestigationResultsForUpdation,
			fetchUpdatedInvestigationResults(currentInvestigationResults, updatedInvestigationResults)...)
	}

	return finalInvestigationResultsForCreation, finalInvestigationResultsForUpdation
}

func getCreateUpdateInvestigationResultsMetadataDto(
	createInvestigationResults, updateInvestigationResults []models.InvestigationResult,
	masterInvestigationIdInvestigationResultMetadataMap map[uint]models.InvestigationResultMetadata,
) ([]models.InvestigationResultMetadata, []models.InvestigationResultMetadata) {
	finalInvestigationResultsMetadataForCreation, finalInvestigationResultsMetadataForUpdation :=
		[]models.InvestigationResultMetadata{}, []models.InvestigationResultMetadata{}

	masterInvestigationIdsForCreation, masterInvestigationIdsForUpdation, masterInvestigationIdsInvestigationResultMap :=
		[]uint{}, []uint{}, map[uint]models.InvestigationResult{}

	for _, investigationResult := range createInvestigationResults {
		masterInvestigationIdsForCreation = append(masterInvestigationIdsForCreation, investigationResult.MasterInvestigationId)
		masterInvestigationIdsInvestigationResultMap[investigationResult.MasterInvestigationId] = investigationResult
	}
	for _, investigationResult := range updateInvestigationResults {
		masterInvestigationIdsForUpdation = append(masterInvestigationIdsForUpdation, investigationResult.MasterInvestigationId)
		masterInvestigationIdsInvestigationResultMap[investigationResult.MasterInvestigationId] = investigationResult
	}

	for masterInvestigationId, investigationResultMetadata := range masterInvestigationIdInvestigationResultMetadataMap {
		investigationResultMetadata.InvestigationResultId = masterInvestigationIdsInvestigationResultMap[masterInvestigationId].Id
		if utils.SliceContainsUint(masterInvestigationIdsForCreation, masterInvestigationId) {
			finalInvestigationResultsMetadataForCreation = append(finalInvestigationResultsMetadataForCreation, investigationResultMetadata)
		} else if utils.SliceContainsUint(masterInvestigationIdsForUpdation, masterInvestigationId) {
			finalInvestigationResultsMetadataForUpdation = append(finalInvestigationResultsMetadataForUpdation, investigationResultMetadata)
		}
	}

	return finalInvestigationResultsMetadataForCreation, finalInvestigationResultsMetadataForUpdation
}

func getCreateUpdateInvestigationDataDto(
	createInvestigationResults, updateInvestigationResults []models.InvestigationResult,
	investigationCodeInvestigationDataMap map[string]models.InvestigationData,
	currentInvestigationDatas []models.InvestigationData,
) ([]models.InvestigationData, []models.InvestigationData) {
	finalInvestigationDataForCreation, finalInvestigationDataForUpdation :=
		[]models.InvestigationData{}, []models.InvestigationData{}

	currentInvestigationDatasMap := map[uint]models.InvestigationData{}
	for _, investigationData := range currentInvestigationDatas {
		currentInvestigationDatasMap[investigationData.InvestigationResultId] = investigationData
	}

	// Creation DTOs for new InvestigationResults
	for _, investigationResult := range createInvestigationResults {
		if investigationData, ok := investigationCodeInvestigationDataMap[investigationResult.LisCode]; ok {
			createInvestigationData := models.InvestigationData{
				InvestigationResultId: investigationResult.Id,
				Data:                  investigationData.Data,
				DataType:              investigationData.DataType,
			}
			createInvestigationData.CreatedBy = constants.CitadelSystemId
			createInvestigationData.UpdatedBy = constants.CitadelSystemId
			finalInvestigationDataForCreation = append(finalInvestigationDataForCreation, createInvestigationData)
		}
	}

	// Updation DTOs for existing InvestigationResults
	for _, investigationResult := range updateInvestigationResults {
		if investigationData, ok := investigationCodeInvestigationDataMap[investigationResult.LisCode]; ok {
			updatedInvestigationData := currentInvestigationDatasMap[investigationResult.Id]
			updatedInvestigationData.Data = investigationData.Data
			updatedInvestigationData.DataType = investigationData.DataType
			updatedInvestigationData.UpdatedBy = constants.CitadelSystemId
			finalInvestigationDataForUpdation = append(finalInvestigationDataForUpdation, updatedInvestigationData)
		}
	}

	return finalInvestigationDataForCreation, finalInvestigationDataForUpdation
}

func getInvestigationIdForTestDocumentMap(investigationResults []models.InvestigationResult,
	testDocumentMap *map[string][]structures.TestDocumentInfoResponse,
) map[uint][]structures.TestDocumentInfoResponse {
	investigationIdTestDocumentMap := map[uint][]structures.TestDocumentInfoResponse{}
	for _, investigationResult := range investigationResults {
		if testDocuments, ok := (*testDocumentMap)[investigationResult.LisCode]; ok {
			for index := range testDocuments {
				testDocuments[index].InvestigationId = investigationResult.Id
			}
			investigationIdTestDocumentMap[investigationResult.Id] = testDocuments
		}
	}
	return investigationIdTestDocumentMap
}

func getCreateUpdateRemarksDto(createInvestigationResults []models.InvestigationResult,
	updateInvestigationResults []models.InvestigationResult,
	currentMedicalRemarks []models.Remark,
	currentTechnicianRemarks []models.Remark,
	investigationCodeMedicalRemarkMap map[string]models.Remark,
	investigationCodeTechnicianRemarkMap map[string]models.Remark,
) ([]models.Remark, []models.Remark) {

	createRemarks, updateRemarks := []models.Remark{}, []models.Remark{}

	currentMedicalRemarksMap := map[uint]models.Remark{}
	for _, medicalRemark := range currentMedicalRemarks {
		currentMedicalRemarksMap[medicalRemark.InvestigationResultId] = medicalRemark
	}

	currentTechnicianRemarksMap := map[uint]models.Remark{}
	for _, technicianRemark := range currentTechnicianRemarks {
		currentTechnicianRemarksMap[technicianRemark.InvestigationResultId] = technicianRemark
	}

	investigationIdLisCodeMap := map[uint]string{}
	for _, investigationResult := range createInvestigationResults {
		investigationIdLisCodeMap[investigationResult.Id] = investigationResult.LisCode
	}

	for _, investigationResult := range updateInvestigationResults {
		investigationIdLisCodeMap[investigationResult.Id] = investigationResult.LisCode
	}

	// Creation DTOs for new Remarks
	for _, investigationResult := range createInvestigationResults {
		if _, ok := investigationCodeMedicalRemarkMap[investigationResult.LisCode]; ok {
			medicalRemark := investigationCodeMedicalRemarkMap[investigationResult.LisCode]
			medicalRemark.InvestigationResultId = investigationResult.Id
			createRemarks = append(createRemarks, medicalRemark)
		}

		if _, ok := investigationCodeTechnicianRemarkMap[investigationResult.LisCode]; ok {
			technicianRemark := investigationCodeTechnicianRemarkMap[investigationResult.LisCode]
			technicianRemark.InvestigationResultId = investigationResult.Id
			createRemarks = append(createRemarks, technicianRemark)
		}
	}

	// Creation/Updation DTOs for existing Remarks
	for _, investigationResult := range updateInvestigationResults {
		if medicalRemark, ok := investigationCodeMedicalRemarkMap[investigationResult.LisCode]; ok {
			if currentMedicalRemark, ok := currentMedicalRemarksMap[investigationResult.Id]; ok {
				currentMedicalRemark.Description = medicalRemark.Description
				updateRemarks = append(updateRemarks, currentMedicalRemark)
			} else {
				medicalRemark.InvestigationResultId = investigationResult.Id
				createRemarks = append(createRemarks, medicalRemark)
			}
		}

		if technicianRemark, ok := investigationCodeTechnicianRemarkMap[investigationResult.LisCode]; ok {
			if currentTechnicianRemark, ok := currentTechnicianRemarksMap[investigationResult.Id]; ok {
				currentTechnicianRemark.Description = technicianRemark.Description
				updateRemarks = append(updateRemarks, currentTechnicianRemark)
			} else {
				technicianRemark.InvestigationResultId = investigationResult.Id
				createRemarks = append(createRemarks, technicianRemark)
			}
		}
	}

	return createRemarks, updateRemarks
}

func fetchUpdatedInvestigationResults(currentInvestigationResults,
	updatedInvestigationResults []models.InvestigationResult) []models.InvestigationResult {
	finalInvestigationResults := []models.InvestigationResult{}
	currentInvestigationResultsMap := map[string]models.InvestigationResult{}
	for _, investigationResult := range currentInvestigationResults {
		currentInvestigationResultsMap[investigationResult.LisCode] = investigationResult
	}

	updatedInvestigationResultsMap := map[string]models.InvestigationResult{}
	for _, investigationResult := range updatedInvestigationResults {
		updatedInvestigationResultsMap[investigationResult.LisCode] = investigationResult
	}

	for lisCode, currentInvestigationResult := range currentInvestigationResultsMap {
		if _, ok := updatedInvestigationResultsMap[lisCode]; ok {
			if currentInvestigationResult.InvestigationStatus == constants.INVESTIGATION_STATUS_APPROVE &&
				updatedInvestigationResultsMap[lisCode].InvestigationStatus == constants.INVESTIGATION_STATUS_PENDING {
				// If the current investigation result is approved and the updated one is pending, we skip updating it.
				continue
			}

			if currentInvestigationResult.InvestigationStatus != constants.INVESTIGATION_STATUS_APPROVE &&
				updatedInvestigationResultsMap[lisCode].InvestigationStatus == constants.INVESTIGATION_STATUS_APPROVE {
				// If the current investigation result is approved and the updated one is also approve,
				// we skip updating some fields.
				currentInvestigationResult.ApprovalSource = updatedInvestigationResultsMap[lisCode].ApprovalSource
				currentInvestigationResult.IsAutoApproved = updatedInvestigationResultsMap[lisCode].IsAutoApproved
				currentInvestigationResult.ApprovedBy = updatedInvestigationResultsMap[lisCode].ApprovedBy
				currentInvestigationResult.ApprovedAt = updatedInvestigationResultsMap[lisCode].ApprovedAt
			}
			currentInvestigationResult.InvestigationValue = updatedInvestigationResultsMap[lisCode].InvestigationValue
			currentInvestigationResult.InvestigationStatus = updatedInvestigationResultsMap[lisCode].InvestigationStatus
			currentInvestigationResult.MasterInvestigationId = updatedInvestigationResultsMap[lisCode].MasterInvestigationId
			currentInvestigationResult.MasterInvestigationMethodMappingId =
				updatedInvestigationResultsMap[lisCode].MasterInvestigationMethodMappingId
			currentInvestigationResult.Abnormality = updatedInvestigationResultsMap[lisCode].Abnormality
			currentInvestigationResult.IsAbnormal = updatedInvestigationResultsMap[lisCode].IsAbnormal
			currentInvestigationResult.EnteredBy = updatedInvestigationResultsMap[lisCode].EnteredBy
			currentInvestigationResult.EnteredAt = updatedInvestigationResultsMap[lisCode].EnteredAt
			currentInvestigationResult.IsCritical = updatedInvestigationResultsMap[lisCode].IsCritical
			currentInvestigationResult.DeviceValue = updatedInvestigationResultsMap[lisCode].DeviceValue
			currentInvestigationResult.MethodType = updatedInvestigationResultsMap[lisCode].MethodType
			currentInvestigationResult.AutoApprovalFailureReason = updatedInvestigationResultsMap[lisCode].AutoApprovalFailureReason

			finalInvestigationResults = append(finalInvestigationResults, currentInvestigationResult)
		}
	}

	return finalInvestigationResults
}

func updateTaskAndTaskMetadataBasedOnTestDetailsOnUpdate(task models.Task, taskMetadata models.TaskMetadata,
	createUpdateTestDetails []models.TestDetail, allTestDetails []models.TestDetail,
	createUpdateTestMetadataDetails []models.TestDetailsMetadata,
	allTestDetailsMetadata []models.TestDetailsMetadata, labIdLabMap map[uint]structures.Lab,
) (models.Task, models.TaskMetadata) {
	testDetailsMap := map[uint]models.TestDetail{}
	for _, testDetail := range allTestDetails {
		testDetailsMap[testDetail.Id] = testDetail
	}
	for _, testDetail := range createUpdateTestDetails {
		testDetailsMap[testDetail.Id] = testDetail
	}

	testDetails := []models.TestDetail{}
	for _, testDetail := range testDetailsMap {
		testDetails = append(testDetails, testDetail)
	}

	if len(testDetails) == 0 {
		return task, taskMetadata
	}

	testDetailsMetadataMap := map[uint]models.TestDetailsMetadata{}
	for _, testDetailsMetadata := range allTestDetailsMetadata {
		testDetailsMetadataMap[testDetailsMetadata.TestDetailsId] = testDetailsMetadata
	}
	for _, testDetailsMetadata := range createUpdateTestMetadataDetails {
		testDetailsMetadataMap[testDetailsMetadata.TestDetailsId] = testDetailsMetadata
	}

	testDetailsMetadata := []models.TestDetailsMetadata{}
	for _, testDetailMetadata := range testDetailsMetadataMap {
		testDetailsMetadata = append(testDetailsMetadata, testDetailMetadata)
	}

	testCompleteStatuses := []string{
		constants.TEST_STATUS_APPROVE,
		constants.TEST_STATUS_COMPLETED_NOT_SENT,
		constants.TEST_STATUS_COMPLETED_SENT,
		constants.TEST_STATUS_SAMPLE_NOT_RECEIVED,
	}

	allInhouseTestsCompleted := true
	for _, testDetail := range testDetails {
		if !utils.SliceContainsString(testCompleteStatuses, testDetail.Status) &&
			utils.IsTestInhouse(testDetail.ProcessingLabId, task.LabId, labIdLabMap) {
			allInhouseTestsCompleted = false
			break
		}
	}

	taskMetadata.IsCritical = getTaskCriticality(testDetailsMetadata, testDetailsMap)

	// Set Task Status as completed if all tests are approved
	if allInhouseTestsCompleted {
		task.PreviousStatus = task.Status
		task.Status = constants.TASK_STATUS_COMPLETED
		task.CompletedAt = utils.GetCurrentTime()
	}

	// Reopen task if any of the tests is not approved
	if !allInhouseTestsCompleted && task.Status == constants.TASK_STATUS_COMPLETED {
		task.Status = constants.TASK_STATUS_PENDING
	}

	doctorTat := getDoctorTatForTask(testDetails)
	if !doctorTat.IsZero() {
		task.DoctorTat = doctorTat
	}

	return task, taskMetadata
}

func (eventProcessor *EventProcessor) getInvestigationInhouseAbnormality(ctx context.Context,
	investigationValue string, investigation structures.Investigation) string {
	if investigationValue == "" {
		return ""
	}
	abnormality := eventProcessor.AbnormalityService.GetInvestigationAbnormality(ctx, investigationValue, investigation)
	return abnormality
}

func (eventProcessor *EventProcessor) getInvestigationInhouseAutoApprovalStatus(ctx context.Context,
	investigationValue string, investigation structures.Investigation) bool {
	autoApprovalStatus := eventProcessor.AbnormalityService.GetInvestigationAutoApprovalStatus(ctx,
		investigationValue, investigation)
	return autoApprovalStatus
}

func (eventProcessor *EventProcessor) getInvestigationAutoApprovalStatus(ctx context.Context,
	investigationValue, entryTime, methodType string, imDevice, imDeviceFlag bool,
	investigation structures.Investigation, latestPastValueMap map[uint]structures.DeltaValuesStructResponse,
	qcFailedTestCodes []string) (bool, string, string) {

	// Failing auto approval for QC failed tests
	if utils.SliceContainsString(qcFailedTestCodes, investigation.LisCode) {
		return false, constants.APPROVAL_SOURCE_NA, constants.AUTO_APPROVAL_FAIL_REASON_QC_FAILED
	}
	if investigationValue == "" || utils.ConvertStringToFloat32ForAbnormality(investigationValue) < 0 || investigationValue[0] == '-' {
		return false, constants.APPROVAL_SOURCE_NA, constants.AUTO_APPROVAL_FAIL_REASON_INVALID_INVESTIGATION_VALUE
	}

	// Approving default auto approval codes
	if utils.SliceContainsString(constants.DefaultAutoApprovalCodes, investigation.LisCode) {
		return true, constants.APPROVAL_SOURCE_OH, constants.AUTO_APPROVAL_FAIL_REASON_NA
	}

	// Preventing approval based on delta check
	pastValueApprovalStatus := eventProcessor.getApprovalStatusBasedOnPastInvestigationResults(ctx, investigationValue,
		investigation, entryTime, latestPastValueMap)
	if !pastValueApprovalStatus {
		return false, constants.APPROVAL_SOURCE_NA, constants.AUTO_APPROVAL_FAIL_REASON_PAST_RECORD
	}

	// Checking auto approval based on imDevice flag
	if imDevice {
		if imDeviceFlag {
			return true, constants.APPROVAL_SOURCE_IM, constants.AUTO_APPROVAL_FAIL_REASON_NA
		}
		return false, constants.APPROVAL_SOURCE_NA, constants.AUTO_APPROVAL_FAIL_REASON_IM_DEVICE
	}

	// Preventing auto approval for manual input method type
	if methodType == constants.METHOD_TYPE_MANUAL {
		return false, constants.APPROVAL_SOURCE_NA, constants.AUTO_APPROVAL_FAIL_REASON_MANUAL_INPUT
	}

	inHouseAutoApproved := eventProcessor.getInvestigationInhouseAutoApprovalStatus(ctx, investigationValue, investigation)

	// Prevent auto approval if value lies outside the reference range
	if !inHouseAutoApproved {
		return false, constants.APPROVAL_SOURCE_NA, constants.AUTO_APPROVAL_FAIL_REASON_REF_RANGE
	}

	return true, constants.APPROVAL_SOURCE_OH, constants.AUTO_APPROVAL_FAIL_REASON_NA
}

func (eventProcessor *EventProcessor) getApprovalStatusBasedOnPastInvestigationResults(ctx context.Context,
	investigationValue string, investigation structures.Investigation,
	enteredTime string, pastResultsMap map[uint]structures.DeltaValuesStructResponse) bool {
	if investigationValue == "" {
		return false
	}

	if utils.SliceContainsInt(constants.DeltaCheckWhitelistedMasterInvIds, int(investigation.InvestigationId)) {
		utils.AddLog(ctx, constants.DEBUG_LEVEL, fmt.Sprintf("Skipping delta check for whitelisted investigation id: %d",
			investigation.InvestigationId), nil, nil)
		return true
	}

	if enteredTime == "" {
		utils.AddLog(ctx, constants.ERROR_LEVEL, constants.INVALID_INVESTIGATION_ENTERED_TIME, nil,
			fmt.Errorf(constants.EMPTY_ENTERED_TIME, investigation.LisCode))
		return false
	}
	enteredTimeUtc := utils.GetTimeFromString(enteredTime)
	if lastInvestigation, ok := pastResultsMap[investigation.InvestigationId]; !ok {
		// Approve the investigation if there is no last result for the investigation value.
		return true
	} else {
		utils.AddLog(ctx, constants.DEBUG_LEVEL, fmt.Sprintf(constants.LAST_INVESTIGATION_RESULT,
			investigation.InvestigationName, lastInvestigation), nil, nil)

		// Approve the investigation if the last result is older than the past value threshold days.
		diffLastInvestigationTime := utils.GetDaysBetweenTimes(*lastInvestigation.ApprovedAt, *enteredTimeUtc)
		utils.AddLog(ctx, constants.DEBUG_LEVEL, fmt.Sprintf(constants.DAYS_DIFF_LAST_INV_CURRENT_INV,
			diffLastInvestigationTime), nil, nil)
		if investigation.PastValueThresholdDays != 0 && diffLastInvestigationTime > int(investigation.PastValueThresholdDays) {
			utils.AddLog(ctx, constants.DEBUG_LEVEL, fmt.Sprintf(constants.SKIP_DELTA_CHECK_OLD_PAST_VALUE,
				investigation.InvestigationId), nil, nil)
			return true
		}

		floatInvValue, err := strconv.ParseFloat(investigationValue, 64)
		if err != nil {
			return true
		}

		// Approve if the investigation doesn't have RCV values.
		if investigation.RcvPositive == 0 && investigation.RcvNegative == 0 {
			return true
		}

		floatPastInvValue, err := strconv.ParseFloat(lastInvestigation.InvestigationValue, 64)
		if err != nil {
			utils.AddLog(ctx, constants.ERROR_LEVEL, constants.ERROR_WHILE_PARSING_FLOAT_VALUE, nil, err)
			return false
		}
		// Rounding off to 5 decimal places to avoid floating number precision issues.
		percentDiff := (floatInvValue - floatPastInvValue) / floatPastInvValue * 100
		percentDiff = math.Round(percentDiff*1e5) / 1e5

		if percentDiff > 0 && percentDiff <= investigation.RcvPositive {
			return true
		} else if percentDiff < 0 && percentDiff >= investigation.RcvNegative {
			return true
		} else if percentDiff == 0 {
			return true
		}
		return false
	}
}

func (eventProcessor *EventProcessor) fetchInvestigationDetailsFromCds(ctx context.Context,
	tests structures.OmsTestDetails, cityCode string, labId uint,
	patientDetails structures.OmsPatientDetails) (
	map[string]structures.Investigation, error) {

	var patientDob string
	if patientDetails.Dob != nil {
		patientDob = patientDetails.Dob.Format(constants.DateLayout)
	} else if patientDetails.ExpectedDob != nil {
		patientDob = patientDetails.ExpectedDob.Format(constants.DateLayout)
	} else {
		patientDob = utils.CalculatePatientDob(int(patientDetails.AgeYears), int(patientDetails.AgeMonths),
			int(patientDetails.AgeDays))
	}

	investigationCodes := fetchInvestigationCodesFromAttunePayload(tests)

	masterInvestigations, err := eventProcessor.CdsClient.GetInvestigationDetails(ctx, investigationCodes,
		cityCode, labId, patientDob, utils.GetGenderConstant(patientDetails.Gender))
	if err != nil {
		utils.AddLog(ctx, constants.ERROR_LEVEL, constants.ERROR_IN_GETTING_INVESTIGATION_DETAILS, nil, err)
		return map[string]structures.Investigation{}, err
	}

	investigationDetailsMap := map[string]structures.Investigation{}
	for _, investigation := range masterInvestigations {
		investigationDetailsMap[investigation.LisCode] = investigation
	}

	return investigationDetailsMap, nil
}

func (eventProcessor *EventProcessor) fetchMasterTestDepartmentMap(ctx context.Context) (map[string]string, error) {
	departmentMapping, err := eventProcessor.CdsClient.GetDepartmentMapping(ctx)
	if err != nil {
		return map[string]string{}, err
	}
	return departmentMapping, nil
}

func (eventProcessor *EventProcessor) createInvestigationResultDetailsForInvestigation(
	ctx context.Context, orderInfo structures.AttuneOrderInfo,
	masterInvestigationDetailsMap map[string]structures.Investigation,
	latestPastValueMap map[uint]structures.DeltaValuesStructResponse,
	qcFailedTestCodes []string) (models.InvestigationResult, models.InvestigationResultMetadata) {

	masterInvestigationId := masterInvestigationDetailsMap[orderInfo.TestCode].InvestigationId

	masterInvestigationMethodMappingId := masterInvestigationDetailsMap[orderInfo.TestCode].InvestigationMethodMappingId

	investigationName := utils.GetNonEmptyString(masterInvestigationDetailsMap[orderInfo.TestCode].InvestigationName,
		orderInfo.TestName)

	resultRepresentationType := masterInvestigationDetailsMap[orderInfo.TestCode].ResultRepresentationType

	departmentName := utils.GetNonEmptyString(
		utils.ConvertStringToCamelCase(masterInvestigationDetailsMap[orderInfo.TestCode].DepartmentName),
		utils.ConvertStringToCamelCase(orderInfo.DepartmentName))

	uom := utils.GetNonEmptyString(masterInvestigationDetailsMap[orderInfo.TestCode].Unit, orderInfo.UOMCode)

	methodName := utils.GetNonEmptyString(
		strings.TrimSpace(masterInvestigationDetailsMap[orderInfo.TestCode].Method),
		strings.TrimSpace(orderInfo.MethodName))

	methodType := getMethodTypeForInvestigation(orderInfo.DeviceID, orderInfo.MethodName)

	referenceRangeText := getReferenceRangeText(masterInvestigationDetailsMap[orderInfo.TestCode].ReferenceRange,
		orderInfo.ReferenceRange)

	imDevice := utils.StringsEqualIgnoreCase(orderInfo.IMDevice, constants.SmallYes)

	imDeviceFlag := utils.StringsEqualIgnoreCase(orderInfo.IMDeviceFlag, constants.AutoVerified)

	ohAbnormality := eventProcessor.getInvestigationInhouseAbnormality(ctx, orderInfo.TestValue,
		masterInvestigationDetailsMap[orderInfo.TestCode])

	isInvestigationAutoApproved, approvalSource, approvalFailureReason :=
		eventProcessor.getInvestigationAutoApprovalStatus(ctx, orderInfo.TestValue, orderInfo.ResultCapturedAt,
			methodType, imDevice, imDeviceFlag, masterInvestigationDetailsMap[orderInfo.TestCode],
			latestPastValueMap, qcFailedTestCodes)

	isInvestigationCritical := isInvestigationCritical(ohAbnormality)

	resultCapturedBy := orderInfo.ResultCapturedBy
	if resultCapturedBy < 0 {
		resultCapturedBy = 0
	}

	investigationResult := models.InvestigationResult{
		MasterInvestigationId:              masterInvestigationId,
		MasterInvestigationMethodMappingId: masterInvestigationMethodMappingId,
		InvestigationName:                  investigationName,
		ResultRepresentationType:           resultRepresentationType,
		Department:                         departmentName,
		Uom:                                uom,
		Method:                             methodName,
		MethodType:                         methodType,
		ReferenceRangeText:                 referenceRangeText,
		LisCode:                            orderInfo.TestCode,
		Abnormality:                        ohAbnormality,
		IsAbnormal:                         getIsAbnormalFlagForInvestigation(ohAbnormality),
		EnteredBy:                          uint(resultCapturedBy),
		EnteredAt:                          utils.GetEnteredAtTime(orderInfo.ResultCapturedAt),
		InvestigationStatus:                constants.INVESTIGATION_STATUS_PENDING,
		IsAutoApproved:                     isInvestigationAutoApproved,
		IsNablApproved:                     masterInvestigationDetailsMap[orderInfo.TestCode].IsNablAccredited,
		Source:                             constants.SourceAttune,
		IsCritical:                         isInvestigationCritical,
		ApprovalSource:                     approvalSource,
		AutoApprovalFailureReason:          approvalFailureReason,
	}

	investigationResult.CreatedBy = constants.CitadelSystemId
	investigationResult.UpdatedBy = constants.CitadelSystemId

	investigationResultMetadata := models.InvestigationResultMetadata{
		QcFlag:            orderInfo.QcFlag,
		QcLotNumber:       orderInfo.QcLotNumber,
		QcValue:           orderInfo.QcValue,
		QcWestGardWarning: orderInfo.QcWestGardWarning,
		QcStatus:          orderInfo.QcStatus,
	}
	investigationResultMetadata.CreatedBy = constants.CitadelSystemId
	investigationResultMetadata.UpdatedBy = constants.CitadelSystemId

	return investigationResult, investigationResultMetadata
}

func (eventProcessor *EventProcessor) createInvestigationResultsDetailsForPanel(ctx context.Context,
	masterInvestigationDetailsMap map[string]structures.Investigation,
	orderInfo structures.AttuneOrderContentListInfoEvent) (models.InvestigationResult,
	models.InvestigationResultMetadata) {

	masterInvestigationId := masterInvestigationDetailsMap[orderInfo.TestCode].InvestigationId

	masterInvestigationMethodMappingId := masterInvestigationDetailsMap[orderInfo.TestCode].InvestigationMethodMappingId

	investigationName := utils.GetNonEmptyString(
		masterInvestigationDetailsMap[orderInfo.TestCode].InvestigationName, orderInfo.TestName)

	resultRepresentationType := masterInvestigationDetailsMap[orderInfo.TestCode].ResultRepresentationType

	departmentName := utils.GetNonEmptyString(
		utils.ConvertStringToCamelCase(masterInvestigationDetailsMap[orderInfo.TestCode].DepartmentName),
		utils.ConvertStringToCamelCase(orderInfo.DepartmentName))

	uom := utils.GetNonEmptyString(masterInvestigationDetailsMap[orderInfo.TestCode].Unit, orderInfo.UOMCode)

	methodName := utils.GetNonEmptyString(
		strings.TrimSpace(masterInvestigationDetailsMap[orderInfo.TestCode].Method),
		strings.TrimSpace(orderInfo.MethodName))

	methodType := getMethodTypeForInvestigation(orderInfo.DeviceID, methodName)

	referenceRangeText := getReferenceRangeText(
		masterInvestigationDetailsMap[orderInfo.TestCode].ReferenceRange, orderInfo.ReferenceRange)

	ohAbnormality := eventProcessor.getInvestigationInhouseAbnormality(ctx, orderInfo.TestValue,
		masterInvestigationDetailsMap[orderInfo.TestCode])

	isInvestigationCritical := isInvestigationCritical(ohAbnormality)

	investigationResult := models.InvestigationResult{
		MasterInvestigationId:              masterInvestigationId,
		MasterInvestigationMethodMappingId: masterInvestigationMethodMappingId,
		InvestigationName:                  investigationName,
		ResultRepresentationType:           resultRepresentationType,
		Department:                         departmentName,
		Uom:                                uom,
		Method:                             methodName,
		MethodType:                         methodType,
		ReferenceRangeText:                 referenceRangeText,
		LisCode:                            orderInfo.TestCode,
		Abnormality:                        ohAbnormality,
		IsAbnormal:                         getIsAbnormalFlagForInvestigation(ohAbnormality),
		EnteredBy:                          uint(orderInfo.ResultCapturedBy),
		EnteredAt:                          utils.GetEnteredAtTime(orderInfo.ResultCapturedAt),
		InvestigationStatus:                constants.INVESTIGATION_STATUS_PENDING,
		IsNablApproved:                     masterInvestigationDetailsMap[orderInfo.TestCode].IsNablAccredited,
		Source:                             constants.SourceAttune,
		IsCritical:                         isInvestigationCritical,
	}

	investigationResult.CreatedBy = constants.CitadelSystemId
	investigationResult.UpdatedBy = constants.CitadelSystemId

	investigationResultMetadata := models.InvestigationResultMetadata{
		QcFlag:            orderInfo.QcFlag,
		QcLotNumber:       orderInfo.QcLotNumber,
		QcValue:           orderInfo.QcValue,
		QcWestGardWarning: orderInfo.QcWestGardWarning,
		QcStatus:          orderInfo.QcStatus,
	}
	investigationResultMetadata.CreatedBy = constants.CitadelSystemId
	investigationResultMetadata.UpdatedBy = constants.CitadelSystemId

	return investigationResult, investigationResultMetadata
}

func (eventProcessor *EventProcessor) createFlattenedInvestigationResultsAndInitialTestDetailsMap(ctx context.Context,
	eventType, cityCode string, tests structures.OmsTestDetails,
	masterInvestigationDetailsMap map[string]structures.Investigation, usersList []models.User, patientId string,
	labId uint) (
	map[string][]models.InvestigationResult, map[uint]models.InvestigationResultMetadata, map[string]structures.InitialTestDetails,
	map[string]models.InvestigationData, map[string]models.Remark, map[string]models.Remark, []string, map[string][]structures.TestDocumentInfoResponse) {

	testIdInvestigationResultsMap := map[string][]models.InvestigationResult{}
	masterInvestigationIdInvestigationResultMetadataMap := map[uint]models.InvestigationResultMetadata{}
	investigationCodeInvestigationDataMap := map[string]models.InvestigationData{}
	testIdInitialTestDetailsMap := map[string]structures.InitialTestDetails{}
	investigationCodeMedicalRemarkMap := map[string]models.Remark{}
	investigationCodeTechnicianRemarkMap := map[string]models.Remark{}
	testCodeTestIdMap := map[string]string{}
	testDocumentMap := map[string][]structures.TestDocumentInfoResponse{}

	for _, testDetails := range tests.TestDetails {
		testCodeTestIdMap[testDetails.TestCode] = testDetails.TestId
	}
	qcFailedOmsTestIds, qcFailedTestCodes := getQcFailedOmsTestIds(tests.OrderInfo, testCodeTestIdMap, labId)

	masterInvestigationIds := []uint{}
	for _, invResults := range masterInvestigationDetailsMap {
		masterInvestigationIds = append(masterInvestigationIds, invResults.InvestigationId)
	}
	pastResultsMap := map[uint]structures.DeltaValuesStructResponse{}
	var cErr *structures.CommonError
	if patientId != "" {
		pastResultsMap, cErr = eventProcessor.InvestigationResultsService.GetLastInvValueByPatientId(ctx, patientId,
			masterInvestigationIds)
		if cErr != nil {
			utils.AddLog(ctx, constants.ERROR_LEVEL, constants.ERROR_IN_GETTING_LAST_INVESTIGATION_VALUE, nil,
				errors.New(cErr.Message))
		}
	}

	attuneUserIdToUserIdMap := map[int]uint{}
	userIdToAttuneUserIdMap := map[uint]int{}
	for _, user := range usersList {
		attuneUserId := utils.ConvertStringToInt(user.AttuneUserId)
		attuneUserIdToUserIdMap[attuneUserId] = user.Id
		userIdToAttuneUserIdMap[user.Id] = attuneUserId
	}

	for _, lisTestInfo := range tests.OrderInfo {
		lisOrderInfo := lisTestInfo.MetaData

		investigationResults, initialTestDetails := []models.InvestigationResult{}, structures.InitialTestDetails{}
		switch lisOrderInfo.TestType {
		case constants.InvestigationShortHand:
			investigationResults, initialTestDetails, investigationCodeInvestigationDataMap,
				investigationCodeMedicalRemarkMap, investigationCodeTechnicianRemarkMap, testDocumentMap =
				eventProcessor.createInvestigationResultsAndInitialTestDetailsForInvestigation(ctx, eventType, cityCode,
					lisOrderInfo, attuneUserIdToUserIdMap, masterInvestigationIdInvestigationResultMetadataMap, masterInvestigationDetailsMap,
					investigationCodeInvestigationDataMap, investigationCodeMedicalRemarkMap,
					investigationCodeTechnicianRemarkMap, pastResultsMap, qcFailedTestCodes, testDocumentMap)
		case constants.GroupShortHand:
			investigationResults, initialTestDetails, investigationCodeInvestigationDataMap,
				investigationCodeMedicalRemarkMap, investigationCodeTechnicianRemarkMap, testDocumentMap =
				eventProcessor.createInvestigationResultsAndInitialTestDetailsForPanel(ctx, eventType, cityCode,
					lisOrderInfo, attuneUserIdToUserIdMap, masterInvestigationIdInvestigationResultMetadataMap, masterInvestigationDetailsMap,
					investigationCodeInvestigationDataMap, investigationCodeMedicalRemarkMap,
					investigationCodeTechnicianRemarkMap, pastResultsMap, qcFailedTestCodes, testDocumentMap)
		}
		testIdInitialTestDetailsMap[testCodeTestIdMap[lisOrderInfo.TestCode]] = initialTestDetails
		testIdInvestigationResultsMap[testCodeTestIdMap[lisOrderInfo.TestCode]] = investigationResults
	}

	return testIdInvestigationResultsMap, masterInvestigationIdInvestigationResultMetadataMap, testIdInitialTestDetailsMap,
		investigationCodeInvestigationDataMap, investigationCodeMedicalRemarkMap, investigationCodeTechnicianRemarkMap,
		qcFailedOmsTestIds, testDocumentMap
}

func (eventProcessor *EventProcessor) createInvestigationResultsAndInitialTestDetailsForInvestigation(
	ctx context.Context, eventType, cityCode string, lisOrderInfo structures.AttuneOrderInfo,
	attuneUserIdToUserIdMap map[int]uint,
	masterInvestigationIdInvestigationResultMetadataMap map[uint]models.InvestigationResultMetadata,
	masterInvestigationDetailsMap map[string]structures.Investigation,
	investigationCodeInvestigationDataMap map[string]models.InvestigationData,
	investigationCodeMedicalRemarkMap map[string]models.Remark,
	investigationCodeTechnicianRemarkMap map[string]models.Remark,
	latestPastValueMap map[uint]structures.DeltaValuesStructResponse,
	qcFailedTestCodes []string,
	testDocumentMap map[string][]structures.TestDocumentInfoResponse,
) ([]models.InvestigationResult, structures.InitialTestDetails,
	map[string]models.InvestigationData, map[string]models.Remark, map[string]models.Remark, map[string][]structures.TestDocumentInfoResponse) {
	orderInfo, isTestAutoApproved, isTestCritical := lisOrderInfo, true, false
	areAttuneStatusesApproved := eventType == constants.OmsApprovedEvent
	investigationResults := []models.InvestigationResult{}

	autoApproveId, _ := strconv.ParseUint(constants.AutoApprovalIdsMap[strings.ToLower(cityCode)], 10, 64)
	currentTime := utils.GetCurrentTime()

	investigationResult, investigationResultMetadata := eventProcessor.createInvestigationResultDetailsForInvestigation(ctx,
		orderInfo, masterInvestigationDetailsMap, latestPastValueMap, qcFailedTestCodes)

	if investigationResult.IsCritical {
		isTestCritical = true
	}

	if eventType == constants.OmsApprovedEvent {
		investigationResult.InvestigationStatus = constants.INVESTIGATION_STATUS_APPROVE
		approvedBy := attuneUserIdToUserIdMap[orderInfo.ResultApprovedBy]
		if approvedBy == 0 {
			approvedBy = uint(autoApproveId)
		}
		investigationResult.ApprovedBy = approvedBy
		investigationResult.ApprovedAt = utils.GetApprovedAtTime(*orderInfo.ResultApprovedAt)
	} else if investigationResult.IsAutoApproved {
		investigationResult.InvestigationStatus = constants.INVESTIGATION_STATUS_APPROVE
		investigationResult.ApprovedBy = uint(autoApproveId)
		investigationResult.ApprovedAt = currentTime
	}

	if masterInvestigationDetailsMap[orderInfo.TestCode].ResultRepresentationType == constants.BiopatternRepresentationType {
		investigationResult.InvestigationValue = orderInfo.TestValue
		investigationResult.DeviceValue = orderInfo.DeviceActualValue
	} else {
		data, _ := utils.ParseXMLToJSON(ctx, orderInfo.TestValue)
		investigationData := models.InvestigationData{
			Data:     data,
			DataType: constants.InvestigationValue,
		}
		investigationData.CreatedBy = constants.CitadelSystemId
		investigationData.UpdatedBy = constants.CitadelSystemId
		investigationCodeInvestigationDataMap[orderInfo.TestCode] = investigationData
	}

	if orderInfo.MedicalRemarks != "" {
		investigationCodeMedicalRemarkMap[orderInfo.TestCode] = createInitialRemarkModel(
			constants.REMARK_TYPE_MEDICAL_REMARK, orderInfo.MedicalRemarks)
	}
	if orderInfo.TechnicalRemarks != "" {
		investigationCodeTechnicianRemarkMap[orderInfo.TestCode] = createInitialRemarkModel(
			constants.REMARK_TYPE_TECHINICIAN_REMARK, orderInfo.TechnicalRemarks)
	}

	investigationResults = append(investigationResults, investigationResult)
	masterInvestigationIdInvestigationResultMetadataMap[investigationResult.MasterInvestigationId] = investigationResultMetadata

	if !investigationResult.IsAutoApproved {
		isTestAutoApproved = false
	}

	doctorTat := utils.AddMinutesToTime(utils.GetEnteredAtTime(orderInfo.ResultCapturedAt), constants.DoctorTatDuration)
	initialTestDetails := structures.InitialTestDetails{
		DoctorTat:      doctorTat,
		IsAutoApproved: isTestAutoApproved,
		TestStatus:     getInitialTestStatusForTestDetails(isTestAutoApproved, eventType, areAttuneStatusesApproved),
		ApprovalSource: getTestApprovalSource(eventType, isTestAutoApproved),
		IsCritical:     isTestCritical,
	}
	return investigationResults, initialTestDetails, investigationCodeInvestigationDataMap,
		investigationCodeMedicalRemarkMap, investigationCodeTechnicianRemarkMap, testDocumentMap
}

func (eventProcessor *EventProcessor) createInvestigationResultsAndInitialTestDetailsForPanel(ctx context.Context,
	eventType, cityCode string, lisOrderInfo structures.AttuneOrderInfo, attuneUserIdToUserIdMap map[int]uint,
	masterInvestigationIdInvestigationResultMetadataMap map[uint]models.InvestigationResultMetadata,
	masterInvestigationDetailsMap map[string]structures.Investigation,
	investigationCodeInvestigationDataMap map[string]models.InvestigationData,
	investigationCodeMedicalRemarkMap map[string]models.Remark,
	investigationCodeTechnicianRemarkMap map[string]models.Remark,
	pastResultsMap map[uint]structures.DeltaValuesStructResponse,
	qcFailedTestCodes []string,
	testDocumentMap map[string][]structures.TestDocumentInfoResponse,
) ([]models.InvestigationResult, structures.InitialTestDetails,
	map[string]models.InvestigationData, map[string]models.Remark, map[string]models.Remark,
	map[string][]structures.TestDocumentInfoResponse) {

	investigationResults, queue, index := []models.InvestigationResult{}, []interface{}{}, 0
	isTestAutoApproved, isTestCritical, areAttuneStatusesApproved := true, false, true
	imDevice, imDeviceFlag := false, true

	autoApproveId, _ := strconv.ParseUint(constants.AutoApprovalIdsMap[strings.ToLower(cityCode)], 10, 64)
	currentTime := utils.GetCurrentTime()

	queue = append(queue, lisOrderInfo.OrderContentListInfo)

	for index < len(queue) {
		currentNode := queue[index]
		index += 1
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

		for _, orderInfo := range attuneOrderContentListInfo {
			switch orderInfo.TestType {
			case constants.InvestigationShortHand:
				investigationResult, investigationResultMetadata := eventProcessor.createInvestigationResultsDetailsForPanel(ctx,
					masterInvestigationDetailsMap, orderInfo)
				if len(orderInfo.TestDocumentInfo) > 0 {
					for _, testDocument := range orderInfo.TestDocumentInfo {
						testDocumentMap[orderInfo.TestCode] = append(testDocumentMap[orderInfo.TestCode], structures.TestDocumentInfoResponse{
							TestCode:     orderInfo.TestCode,
							TestDocument: testDocument.TestDocument,
						})
					}
				}

				if investigationResult.IsCritical {
					isTestCritical = true
				}

				if eventType == constants.OmsApprovedEvent && orderInfo.TestStatus != constants.AttuneTestStatusApprove {
					// If any of the investigation is not approved, then the test is not approved
					areAttuneStatusesApproved = false
				}

				if masterInvestigationDetailsMap[orderInfo.TestCode].ResultRepresentationType ==
					constants.BiopatternRepresentationType {
					investigationResult.InvestigationValue = orderInfo.TestValue
					investigationResult.DeviceValue = orderInfo.DeviceActualValue
				} else {
					data, _ := utils.ParseXMLToJSON(ctx, orderInfo.TestValue)
					investigationData := models.InvestigationData{
						Data:     data,
						DataType: constants.InvestigationValue,
					}
					investigationData.CreatedBy = constants.CitadelSystemId
					investigationData.UpdatedBy = constants.CitadelSystemId
					investigationCodeInvestigationDataMap[orderInfo.TestCode] = investigationData
				}

				if orderInfo.MedicalRemarks != "" {
					investigationCodeMedicalRemarkMap[orderInfo.TestCode] = createInitialRemarkModel(
						constants.REMARK_TYPE_MEDICAL_REMARK, orderInfo.MedicalRemarks)
				}
				if orderInfo.TechnicalRemarks != "" {
					investigationCodeTechnicianRemarkMap[orderInfo.TestCode] = createInitialRemarkModel(
						constants.REMARK_TYPE_TECHINICIAN_REMARK, orderInfo.TechnicalRemarks)
				}

				if utils.StringsEqualIgnoreCase(orderInfo.IMDevice, constants.SmallYes) {
					imDevice = true
				}

				investigationResults = append(investigationResults, investigationResult)
				masterInvestigationIdInvestigationResultMetadataMap[investigationResult.MasterInvestigationId] = investigationResultMetadata
				if !utils.StringsEqualIgnoreCase(orderInfo.IMDeviceFlag, constants.AutoVerified) {
					imDeviceFlag = false
				}
			case constants.GroupShortHand:
				queue = append(queue, orderInfo.ParameterListInfo)
			}
		}
	}

	for index, investigationResult := range investigationResults {
		isInvestigationAutoApproved, approvalSource, approvalFailureReason :=
			eventProcessor.getInvestigationAutoApprovalStatus(ctx, investigationResult.InvestigationValue,
				utils.GetTimeInString(investigationResult.EnteredAt), investigationResult.MethodType,
				imDevice, imDeviceFlag, masterInvestigationDetailsMap[investigationResult.LisCode],
				pastResultsMap, qcFailedTestCodes)
		investigationResult.IsAutoApproved = isInvestigationAutoApproved
		investigationResult.ApprovalSource = approvalSource
		investigationResult.AutoApprovalFailureReason = approvalFailureReason
		investigationResults[index] = investigationResult
	}

	for _, investigationResult := range investigationResults {
		if !investigationResult.IsAutoApproved {
			isTestAutoApproved = false
			break
		}
	}

	if eventType == constants.OmsApprovedEvent && areAttuneStatusesApproved {
		for index := range investigationResults {
			investigationResults[index].InvestigationStatus = constants.INVESTIGATION_STATUS_APPROVE
			approvedBy := attuneUserIdToUserIdMap[lisOrderInfo.ResultApprovedBy]
			if approvedBy == 0 {
				approvedBy = uint(autoApproveId)
			}
			investigationResults[index].ApprovedBy = approvedBy
			investigationResults[index].ApprovedAt = utils.GetApprovedAtTime(*lisOrderInfo.ResultApprovedAt)
		}
	} else if isTestAutoApproved {
		for index := range investigationResults {
			investigationResults[index].InvestigationStatus = constants.INVESTIGATION_STATUS_APPROVE
			investigationResults[index].ApprovedAt = currentTime
			investigationResults[index].ApprovedBy = uint(autoApproveId)
		}
	}

	doctorTat := utils.AddMinutesToTime(utils.GetEnteredAtTime(lisOrderInfo.ResultCapturedAt), constants.DoctorTatDuration)
	initialTestDetails := structures.InitialTestDetails{
		DoctorTat:      doctorTat,
		IsAutoApproved: isTestAutoApproved,
		TestStatus:     getInitialTestStatusForTestDetails(isTestAutoApproved, eventType, areAttuneStatusesApproved),
		ApprovalSource: getTestApprovalSource(eventType, isTestAutoApproved),
		IsCritical:     isTestCritical,
	}

	return investigationResults, initialTestDetails, investigationCodeInvestigationDataMap,
		investigationCodeMedicalRemarkMap, investigationCodeTechnicianRemarkMap, testDocumentMap
}

func (eventProcessor *EventProcessor) createData(ctx context.Context,
	omsLisEvent structures.OmsLisEvent,
	testIdInvestigationResultsMap map[string][]models.InvestigationResult,
	masterInvestigationIdInvestigationResultMetadataMap map[uint]models.InvestigationResultMetadata,
	investigationCodeInvestigationDataMap map[string]models.InvestigationData,
	testIdInitialTestDetailsMap map[string]structures.InitialTestDetails,
	masterTestDepartmentMap map[string]string,
	investigationCodeMedicalRemarkMap map[string]models.Remark,
	investigationCodeTechnicianRemarkMap map[string]models.Remark,
	qcFailedOmsTestIds []string, userIdToAttuneUserIdMap map[uint]int, labIdLabMap map[uint]structures.Lab,
	testDocumentMap *map[string][]structures.TestDocumentInfoResponse,
) (models.Task, *structures.CommonError) {

	task := models.Task{}

	err := eventProcessor.Db.Transaction(func(tx *gorm.DB) error {
		// Create PatientDetails
		patientDetails, cErr := eventProcessor.PatientDetailService.CreatePatientDetailsWithTx(tx,
			createPatientDetailsDto(omsLisEvent))
		if cErr != nil {
			return errors.New(cErr.Message)
		}

		// Create Task
		task, cErr = eventProcessor.TaskService.CreateTaskWithTx(tx,
			createTaskDto(omsLisEvent, patientDetails))
		if cErr != nil {
			return errors.New(cErr.Message)
		}

		// Create TaskMetadata
		taskMetadata, cErr := eventProcessor.TaskService.CreateTaskMetadataWithTx(tx,
			createTaskMetadataDto(omsLisEvent, task))
		if cErr != nil {
			return errors.New(cErr.Message)
		}

		// Create TaskVisitMapping
		_, cErr = eventProcessor.TaskService.CreateTaskVisitMappingWithTx(tx,
			createTaskVisitMappingDto(omsLisEvent, task))
		if cErr != nil {
			return errors.New(cErr.Message)
		}

		// Create TestDetails
		testDetails, cErr := eventProcessor.TestDetailService.CreateTestDetailsWithTx(tx,
			createTestDetailsDto(omsLisEvent, task, testIdInitialTestDetailsMap, masterTestDepartmentMap,
				qcFailedOmsTestIds))
		if cErr != nil {
			return errors.New(cErr.Message)
		}

		omsTestIdTestDetailMap := map[string]models.TestDetail{}
		for _, testDetail := range testDetails {
			omsTestIdTestDetailMap[testDetail.CentralOmsTestId] = testDetail
		}

		// Create TestDetailsMetadata
		testDetailsMetadata, cErr := eventProcessor.TestDetailService.CreateTestDetailsMetadataWithTx(tx,
			createTestDetailsMetadataDto(omsLisEvent, testDetails, testIdInitialTestDetailsMap))
		if cErr != nil {
			return errors.New(cErr.Message)
		}

		// Create InvestigationResults
		investigationResults, cErr := eventProcessor.InvestigationResultsService.CreateInvestigationResultsWithTx(tx,
			createInvestigationResultsDto(testDetails, testIdInvestigationResultsMap, qcFailedOmsTestIds))
		if cErr != nil {
			return errors.New(cErr.Message)
		}

		// Create InvestigationResultMetadata
		_, cErr = eventProcessor.InvestigationResultsService.CreateInvestigationResultsMetadataWithTx(ctx, tx,
			createInvestigationResultMetadataDto(investigationResults, masterInvestigationIdInvestigationResultMetadataMap))
		if cErr != nil {
			return errors.New(cErr.Message)
		}

		// Create InvestigationData
		_, cErr = eventProcessor.InvestigationResultsService.CreateInvestigationDataWithTx(tx,
			createInvestigationDataDto(investigationResults, investigationCodeInvestigationDataMap))
		if cErr != nil {
			return errors.New(cErr.Message)
		}

		// Map test documents to investigation IDs
		getInvestigationIdForTestDocumentMap(investigationResults, testDocumentMap)

		// Create Remarks
		_, cErr = eventProcessor.RemarkService.CreateRemarksWithTx(tx,
			createRemarksDto(investigationResults, investigationCodeMedicalRemarkMap, investigationCodeTechnicianRemarkMap))
		if cErr != nil {
			return errors.New(cErr.Message)
		}

		// Update Task and Task Metadata based on TestDetails
		task, taskMetadata = updateTaskAndTaskMetadataBasedOnTestDetailsOnCreate(task, taskMetadata, testDetails,
			testDetailsMetadata, labIdLabMap)
		_, cErr = eventProcessor.TaskService.UpdateTaskWithTx(tx, task)
		if cErr != nil {
			return errors.New(cErr.Message)
		}
		_, cErr = eventProcessor.TaskService.UpdateTaskMetadataWithTx(tx, taskMetadata)
		if cErr != nil {
			return errors.New(cErr.Message)
		}

		// Handle rerun data and sync to attune steps
		rerunDetails, visitIdToAttuneOrderMap, cErr := eventProcessor.getQcFailedRerunData(ctx, qcFailedOmsTestIds,
			omsTestIdTestDetailMap, omsLisEvent.Order.CityCode, investigationResults, []models.RerunInvestigationResult{}, userIdToAttuneUserIdMap)
		if cErr != nil {
			return errors.New(cErr.Message)
		}

		_, cErr = eventProcessor.RerunService.UpdateRerunInvestigationResultsWithTx(tx, rerunDetails)
		if cErr != nil {
			return errors.New(cErr.Message)
		}

		for _, attuneResponse := range visitIdToAttuneOrderMap {
			cErr = eventProcessor.AttuneClient.InsertTestDataToAttune(ctx, attuneResponse)
			if cErr != nil {
				return errors.New(cErr.Message)
			}
		}

		return nil
	})

	if err != nil {
		utils.AddLog(ctx, constants.ERROR_LEVEL, constants.ERROR_IN_TRANSACTION, nil, err)
		return task, &structures.CommonError{
			Message:    err.Error(),
			StatusCode: http.StatusInternalServerError,
		}
	}

	return task, nil
}

func (eventProcessor *EventProcessor) createUpdateData(ctx context.Context,
	omsLisEvent structures.OmsLisEvent, task models.Task,
	testIdInvestigationResultsMap map[string][]models.InvestigationResult,
	masterInvestigationIdInvestigationResultMetadataMap map[uint]models.InvestigationResultMetadata,
	investigationCodeInvestigationDataMap map[string]models.InvestigationData,
	testIdInitialTestDetailsMap map[string]structures.InitialTestDetails,
	masterTestDepartmentMap map[string]string,
	investigationCodeMedicalRemarkMap map[string]models.Remark,
	investigationCodeTechnicianRemarkMap map[string]models.Remark,
	qcFailedOmsTestIds []string, userIdToAttuneUserIdMap map[uint]int, labIdLabMap map[uint]structures.Lab,
	testDocumentMap *map[string][]structures.TestDocumentInfoResponse,
) (models.Task, *structures.CommonError) {
	patientDetails, taskMetadata, taskVisitMappings, testDetails, testDetailsMetadata,
		investigationResults, investigationResultsMetadataMap, investigationDatas, rerunInvResults, medicalRemarks,
		technicianRemarks, cErr := eventProcessor.fetchCurrentData(ctx, task)
	if cErr != nil {
		utils.AddLog(ctx, constants.ERROR_LEVEL, utils.GetCurrentFunctionName(), nil, errors.New(cErr.Message))
		return task, cErr
	}

	err := eventProcessor.Db.Transaction(func(tx *gorm.DB) error {
		// Update Patient Details
		patientDetails, cErr = eventProcessor.PatientDetailService.UpdatePatientDetailsWithTx(tx,
			createUpdatePatientDetailsDto(omsLisEvent, patientDetails))
		if cErr != nil {
			return errors.New(cErr.Message)
		}

		// Update TaskMetadata
		taskMetadata, cErr = eventProcessor.TaskService.UpdateTaskMetadataWithTx(tx,
			createUpdateTaskMetadataDto(omsLisEvent, taskMetadata))
		if cErr != nil {
			return errors.New(cErr.Message)
		}

		// Create/Update TaskVisitMapping
		createVisitIds, deleteVisitIds := getCreateDeleteTaskVisits(omsLisEvent, taskVisitMappings)
		cErr = eventProcessor.TaskService.CreateDeleteTaskVisitMappingsWithTx(tx, task.Id,
			createVisitIds, deleteVisitIds)
		if cErr != nil {
			return errors.New(cErr.Message)
		}

		// Create/Update TestDetails
		createTestDetails, updateTestDetails := getCreateUpdateTestDetailsDto(omsLisEvent, task,
			testDetails, testIdInitialTestDetailsMap, masterTestDepartmentMap, qcFailedOmsTestIds)
		createTestDetails, cErr = eventProcessor.TestDetailService.CreateTestDetailsWithTx(tx, createTestDetails)
		if cErr != nil {
			return errors.New(cErr.Message)
		}
		updateTestDetails, cErr = eventProcessor.TestDetailService.UpdateTestDetailsWithTx(tx, updateTestDetails)
		if cErr != nil {
			return errors.New(cErr.Message)
		}

		omsTestIdTestDetailMap := map[string]models.TestDetail{}
		for _, testDetail := range createTestDetails {
			omsTestIdTestDetailMap[testDetail.CentralOmsTestId] = testDetail
		}
		for _, testDetail := range updateTestDetails {
			omsTestIdTestDetailMap[testDetail.CentralOmsTestId] = testDetail
		}

		// Create/Update TestDetailsMetadata
		createTestDetailsMetadata, updateTestDetailsMetadata := getCreateUpdateTestDetailsMetadataDtos(
			omsLisEvent, testDetailsMetadata, createTestDetails, updateTestDetails, testIdInitialTestDetailsMap)
		_, cErr = eventProcessor.TestDetailService.CreateTestDetailsMetadataWithTx(tx, createTestDetailsMetadata)
		if cErr != nil {
			return errors.New(cErr.Message)
		}
		_, cErr = eventProcessor.TestDetailService.UpdateTestDetailsMetadataWithTx(tx, updateTestDetailsMetadata)
		if cErr != nil {
			return errors.New(cErr.Message)
		}

		// Create/Update InvestigationResults
		createInvestigationResults, updateInvestigationResults := getCreateUpdateInvestigationResultsDto(
			createTestDetails, updateTestDetails, testIdInvestigationResultsMap, investigationResults,
			qcFailedOmsTestIds)
		createInvestigationResults, cErr = eventProcessor.InvestigationResultsService.CreateInvestigationResultsWithTx(
			tx, createInvestigationResults)
		if cErr != nil {
			return errors.New(cErr.Message)
		}
		updateInvestigationResults, cErr = eventProcessor.InvestigationResultsService.UpdateInvestigationResultsWithTx(
			tx, updateInvestigationResults)
		if cErr != nil {
			return errors.New(cErr.Message)
		}

		// Create/Update InvestigationResultMetadata
		createInvestigationResultsMetadata, updateInvestigationResultsMetadata := getCreateUpdateInvestigationResultsMetadataDto(
			createInvestigationResults, updateInvestigationResults, masterInvestigationIdInvestigationResultMetadataMap)
		_, cErr = eventProcessor.InvestigationResultsService.CreateInvestigationResultsMetadataWithTx(ctx,
			tx, createInvestigationResultsMetadata)
		if cErr != nil {
			return errors.New(cErr.Message)
		}
		_, cErr = eventProcessor.InvestigationResultsService.UpdateInvestigationResultsMetadataWithTx(ctx,
			tx, updateInvestigationResultsMetadata, investigationResultsMetadataMap)
		if cErr != nil {
			return errors.New(cErr.Message)
		}

		// Create/Update InvestigationData
		createInvestigationData, updateInvestigationData := getCreateUpdateInvestigationDataDto(
			createInvestigationResults, updateInvestigationResults,
			investigationCodeInvestigationDataMap, investigationDatas)
		_, cErr = eventProcessor.InvestigationResultsService.CreateInvestigationDataWithTx(tx, createInvestigationData)
		if cErr != nil {
			return errors.New(cErr.Message)
		}
		_, cErr = eventProcessor.InvestigationResultsService.UpdateInvestigationDataWithTx(tx, updateInvestigationData)
		if cErr != nil {
			return errors.New(cErr.Message)
		}

		// Map test documents to investigation IDs
		combinedInvestigationResults := append(createInvestigationResults, updateInvestigationResults...)
		getInvestigationIdForTestDocumentMap(combinedInvestigationResults, testDocumentMap)

		// Create/Update Remarks
		createRemarks, updateRemarks := getCreateUpdateRemarksDto(createInvestigationResults, updateInvestigationResults,
			medicalRemarks, technicianRemarks, investigationCodeMedicalRemarkMap, investigationCodeTechnicianRemarkMap)
		_, cErr = eventProcessor.RemarkService.CreateRemarksWithTx(tx, createRemarks)
		if cErr != nil {
			return errors.New(cErr.Message)
		}
		_, cErr = eventProcessor.RemarkService.UpdateRemarksWithTx(tx, updateRemarks)
		if cErr != nil {
			return errors.New(cErr.Message)
		}

		// Update Task and Task Metadata based on TestDetails
		createUpdateTestDetails := append(createTestDetails, updateTestDetails...)
		createUpdateTestMetadataDetails := append(createTestDetailsMetadata, updateTestDetailsMetadata...)
		task, taskMetadata = updateTaskAndTaskMetadataBasedOnTestDetailsOnUpdate(task, taskMetadata,
			createUpdateTestDetails, testDetails, createUpdateTestMetadataDetails, testDetailsMetadata, labIdLabMap)
		_, cErr = eventProcessor.TaskService.UpdateTaskWithTx(tx, task)
		if cErr != nil {
			return errors.New(cErr.Message)
		}
		_, cErr = eventProcessor.TaskService.UpdateTaskMetadataWithTx(tx, taskMetadata)
		if cErr != nil {
			return errors.New(cErr.Message)
		}

		// Handle rerun data and sync to attune steps
		allInvestigationResults := createInvestigationResults
		allInvestigationResults = append(allInvestigationResults, updateInvestigationResults...)
		rerunDetails, visitIdToAttuneOrderMap, cErr := eventProcessor.getQcFailedRerunData(ctx,
			qcFailedOmsTestIds, omsTestIdTestDetailMap, omsLisEvent.Order.CityCode, allInvestigationResults,
			rerunInvResults, userIdToAttuneUserIdMap)
		if cErr != nil {
			return errors.New(cErr.Message)
		}

		_, cErr = eventProcessor.RerunService.UpdateRerunInvestigationResultsWithTx(tx, rerunDetails)
		if cErr != nil {
			return errors.New(cErr.Message)
		}

		for _, attuneResponse := range visitIdToAttuneOrderMap {
			cErr = eventProcessor.AttuneClient.InsertTestDataToAttune(ctx, attuneResponse)
			if cErr != nil {
				return errors.New(cErr.Message)
			}
		}

		return nil
	})

	if err != nil {
		utils.AddLog(ctx, constants.ERROR_LEVEL, constants.ERROR_IN_TRANSACTION, nil, err)
		return task, &structures.CommonError{
			Message:    err.Error(),
			StatusCode: http.StatusInternalServerError,
		}
	}

	return task, nil
}

func (eventProcessor *EventProcessor) fetchCurrentData(ctx context.Context, task models.Task) (
	patientDetails models.PatientDetail,
	taskMetadata models.TaskMetadata,
	taskVisitMappings []models.TaskVisitMapping,
	testDetails []models.TestDetail,
	testDetailsMetadata []models.TestDetailsMetadata,
	investigationResults []models.InvestigationResult,
	investigationResultsMetadata map[uint]models.InvestigationResultMetadata,
	investigationDatas []models.InvestigationData,
	rerunInvDetails []models.RerunInvestigationResult,
	medicalRemarks []models.Remark,
	technicianRemarks []models.Remark,
	cErr *structures.CommonError,
) {
	patientDetails, cErr = eventProcessor.PatientDetailService.GetPatientDetailsById(task.PatientDetailsId)
	if cErr != nil {
		utils.AddLog(ctx, constants.ERROR_LEVEL, utils.GetCurrentFunctionName(), nil, errors.New(cErr.Message))
		return
	}

	taskMetadata, cErr = eventProcessor.TaskService.GetTaskMetadataByTaskId(task.Id)
	if cErr != nil {
		utils.AddLog(ctx, constants.ERROR_LEVEL, utils.GetCurrentFunctionName(), nil, errors.New(cErr.Message))
		return
	}

	taskVisitMappings, cErr = eventProcessor.TaskService.GetTaskVisitMappingsByTaskId(task.Id)
	if cErr != nil {
		utils.AddLog(ctx, constants.ERROR_LEVEL, utils.GetCurrentFunctionName(), nil, errors.New(cErr.Message))
		return
	}

	testDetails, cErr = eventProcessor.TestDetailService.GetTestDetailsByOmsOrderId(task.OmsOrderId)
	if cErr != nil {
		utils.AddLog(ctx, constants.ERROR_LEVEL, utils.GetCurrentFunctionName(), nil, errors.New(cErr.Message))
		return
	}
	testDetailsIds := []uint{}
	for _, testDetail := range testDetails {
		testDetailsIds = append(testDetailsIds, testDetail.Id)
	}

	testDetailsMetadata, cErr = eventProcessor.TestDetailService.GetTestDetailsMetadataByTestDetailIds(testDetailsIds)
	if cErr != nil {
		utils.AddLog(ctx, constants.ERROR_LEVEL, utils.GetCurrentFunctionName(), nil, errors.New(cErr.Message))
		return
	}

	investigationResults, cErr = eventProcessor.InvestigationResultsService.
		GetInvestigationResultsByTestDetailsIds(testDetailsIds)
	if cErr != nil {
		utils.AddLog(ctx, constants.ERROR_LEVEL, utils.GetCurrentFunctionName(), nil, errors.New(cErr.Message))
		return
	}

	investigationResultIds := []uint{}
	for _, investigationResult := range investigationResults {
		investigationResultIds = append(investigationResultIds, investigationResult.Id)
	}

	if len(investigationResultIds) > 0 {
		investigationResultsMetadata, cErr = eventProcessor.InvestigationResultsService.
			GetInvestigationResultsMetadataByInvestigationResultIds(ctx, investigationResultIds)
		if cErr != nil {
			utils.AddLog(ctx, constants.ERROR_LEVEL, utils.GetCurrentFunctionName(), nil, errors.New(cErr.Message))
			return
		}
	}

	if len(investigationResultIds) > 0 {
		investigationDatas, cErr = eventProcessor.InvestigationResultsService.
			GetInvestigationDataByInvestigationResultsIds(investigationResultIds)
		if cErr != nil {
			utils.AddLog(ctx, constants.ERROR_LEVEL, utils.GetCurrentFunctionName(), nil, errors.New(cErr.Message))
			return
		}
	}

	remarks := []models.Remark{}
	if len(investigationResultIds) > 0 {
		remarks, cErr = eventProcessor.RemarkService.
			GetRemarksByInvestigationResultIds([]string{constants.REMARK_TYPE_TECHINICIAN_REMARK,
				constants.REMARK_TYPE_MEDICAL_REMARK}, investigationResultIds)
		if cErr != nil {
			utils.AddLog(ctx, constants.ERROR_LEVEL, utils.GetCurrentFunctionName(), nil, errors.New(cErr.Message))
			return
		}
	}

	medicalRemarks, technicianRemarks = []models.Remark{}, []models.Remark{}
	for _, remark := range remarks {
		switch remark.RemarkType {
		case constants.REMARK_TYPE_MEDICAL_REMARK:
			medicalRemarks = append(medicalRemarks, remark)
		case constants.REMARK_TYPE_TECHINICIAN_REMARK:
			technicianRemarks = append(technicianRemarks, remark)
		}
	}

	rerunInvDetails, cErr = eventProcessor.RerunService.GetRerunDetailsByTaskId(task.Id)
	if cErr != nil {
		utils.AddLog(ctx, constants.ERROR_LEVEL, utils.GetCurrentFunctionName(), nil, errors.New(cErr.Message))
		return
	}
	return
}

func (eventProcessor *EventProcessor) ReleaseReport(ctx context.Context, taskId uint) *structures.CommonError {

	err := eventProcessor.CommonTaskProcessor.ReleaseReportTask(ctx, taskId, false)
	if err != nil {
		utils.AddLog(ctx, constants.DEBUG_LEVEL, utils.GetCurrentFunctionName(), nil, err)
	}
	return nil
}

func getQcFailedOmsTestIds(testUpdates map[string]structures.LisTestUpdateInfo, testCodeTestIdMap map[string]string,
	labId uint) ([]string, []string) {
	qcFailedTestIds, qcFailedTestCodes := []string{}, []string{}
	if !utils.IsQcEnabledForLab(labId) {
		return qcFailedTestIds, qcFailedTestCodes
	}
	for _, testDetails := range testUpdates {
		if testDetails.MetaData.QcStatus == constants.QCStatusFail {
			qcFailedTestCodes = append(qcFailedTestCodes, testDetails.TestCode)
			testId := testCodeTestIdMap[testDetails.TestCode]
			qcFailedTestIds = append(qcFailedTestIds, testId)
		}
		if len(testDetails.MetaData.OrderContentListInfo) > 0 {
			for _, orderContent := range testDetails.MetaData.OrderContentListInfo {
				if orderContent.QcStatus == constants.QCStatusFail {
					qcFailedTestCodes = append(qcFailedTestCodes, orderContent.TestCode)
					qcFailedTestCodes = append(qcFailedTestCodes, testDetails.TestCode)
					testId := testCodeTestIdMap[testDetails.TestCode]
					qcFailedTestIds = append(qcFailedTestIds, testId)
				}
			}
		}
	}
	return utils.CreateUniqueSliceString(qcFailedTestIds), utils.CreateUniqueSliceString(qcFailedTestCodes)
}

func (eventProcessor *EventProcessor) getQcFailedRerunData(ctx context.Context,
	qcFailedOmsTestIds []string, testIdTestDetailsMap map[string]models.TestDetail, cityCode string,
	invResults []models.InvestigationResult, rerunInvResult []models.RerunInvestigationResult,
	userIdToAttuneUserId map[uint]int) (
	[]models.RerunInvestigationResult, map[string]structures.AttuneOrderResponse, *structures.CommonError) {
	qcFailedTestDetailIds := []uint{}
	qcFailedTestDetails := []models.TestDetail{}
	qcFailedInvResults := []models.InvestigationResult{}

	for _, testId := range qcFailedOmsTestIds {
		testDetail := testIdTestDetailsMap[testId]
		qcFailedTestDetailIds = append(qcFailedTestDetailIds, testDetail.Id)
		qcFailedTestDetails = append(qcFailedTestDetails, testDetail)
	}

	for _, invResult := range invResults {
		if utils.SliceContainsUint(qcFailedTestDetailIds, invResult.TestDetailsId) {
			qcFailedInvResults = append(qcFailedInvResults, invResult)
		}
	}
	// This is a map for TestDetialsId to MasterInvId to RerunInvResult
	existingRerunDetailsMap := map[uint]map[uint]models.RerunInvestigationResult{}
	for _, rerunDetail := range rerunInvResult {
		if _, exists := existingRerunDetailsMap[rerunDetail.TestDetailsId]; !exists {
			existingRerunDetailsMap[rerunDetail.TestDetailsId] = map[uint]models.RerunInvestigationResult{}
			existingRerunDetailsMap[rerunDetail.TestDetailsId][rerunDetail.MasterInvestigationId] = rerunDetail
		}
	}
	if len(qcFailedTestDetailIds) > 0 {
		rerunDetails, visitIdToAttuneOrder, cErr := eventProcessor.TaskService.GetQcFailedTestDataToRerun(ctx,
			constants.CitadelSystemId, qcFailedTestDetailIds, cityCode, qcFailedTestDetails, qcFailedInvResults,
			userIdToAttuneUserId)
		if cErr != nil {
			return nil, nil, cErr
		}

		// Updating the rerun details with existing rerun details if any
		for i, rerunDetail := range rerunDetails {
			if existingRerunDetails, exists := existingRerunDetailsMap[rerunDetail.TestDetailsId]; exists {
				if existingRerunDetail, exists := existingRerunDetails[rerunDetail.MasterInvestigationId]; exists {
					rerunDetails[i].Id = existingRerunDetail.Id
				}
			}
		}
		return rerunDetails, visitIdToAttuneOrder, nil
	}
	return []models.RerunInvestigationResult{}, map[string]structures.AttuneOrderResponse{}, nil
}
