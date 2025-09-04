package service

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/Orange-Health/citadel/apps/report_generation/constants"
	"github.com/Orange-Health/citadel/apps/report_generation/structures"
	commonConstants "github.com/Orange-Health/citadel/common/constants"
	commonStructures "github.com/Orange-Health/citadel/common/structures"
	commonUtils "github.com/Orange-Health/citadel/common/utils"
	commonModels "github.com/Orange-Health/citadel/models"
)

type ReportGenerationInterface interface {
	TriggerReportGenerationEvent(ctx context.Context, omsOrderId string, omsTestIds []string) *commonStructures.CommonError
	TriggerOrderApprovedEvent(ctx context.Context, omsOrderId string, isDummyReport bool) (interface{}, *commonStructures.CommonError)
}

func (s *ReportGenerationService) TriggerReportGenerationEvent(ctx context.Context, omsOrderId string,
	omsTestIds []string) *commonStructures.CommonError {

	orderDetails, cErr := s.OrderDetailsService.GetOrderDetailsByOmsOrderId(omsOrderId)
	if cErr != nil {
		return cErr
	}
	if orderDetails.Id == 0 {
		return &commonStructures.CommonError{
			StatusCode: http.StatusNotFound,
			Message:    commonConstants.ERROR_ORDER_ID_NOT_FOUND,
		}
	}

	reportGenerationEvent, reportGenerationEventAllVisits, cErr :=
		s.fetchReportGenerationEventDetails(ctx, orderDetails, omsTestIds)
	if cErr != nil {
		return cErr
	}
	attachments, err := s.fetchAttachmentDetails(omsOrderId)
	if err != nil {
		return err
	}
	reportGenerationEventAllVisits.ServicingCityCode = orderDetails.CityCode
	go s.GetAndPublishReportGenerationEventForAllVisits(context.Background(), orderDetails, reportGenerationEventAllVisits, attachments, false)

	if len(reportGenerationEvent.Visits) == 0 {
		return &commonStructures.CommonError{
			StatusCode: http.StatusNotFound,
			Message:    commonConstants.ERROR_NO_INHOUSE_VISIT_TESTS_FOUND,
		}
	}

	reportGenerationEvent.ServicingCityCode = orderDetails.CityCode
	reportGenerationEvent.IndexDetails.IsCollectionTypeInClinic =
		orderDetails.CollectionType == commonConstants.CollectionTypePickUpFromPartner

	omsTestIds = s.fetchTestIdsFromReportGenerationEvent(reportGenerationEvent)

	investigations, cErr := s.Dao.GetInvestigationDataForReportGeneration(ctx, omsOrderId, omsTestIds)
	if cErr != nil {
		return cErr
	}
	attachmentsMap := make(map[uint][]string)
	for _, attachment := range attachments {
		if attachment.AttachmentType == commonConstants.AttachmentTypeTestDocument && attachment.InvestigationResultId > 0 {
			attachmentsMap[attachment.InvestigationResultId] = append(attachmentsMap[attachment.InvestigationResultId], attachment.AttachmentUrl)
		}
	}

	approvedByIds := []string{}
	for index := range investigations {
		if attachmentUrls, exists := attachmentsMap[investigations[index].Id]; exists && len(attachmentUrls) > 0 {
			investigations[index].TestDocument = append(investigations[index].TestDocument, attachmentUrls...)
		}
		approvedByIds = append(approvedByIds, investigations[index].ApprovedBy)
	}

	approvedByIds = commonUtils.CreateUniqueSliceString(approvedByIds)

	users, cErr := s.Dao.GetUserDetails(ctx, approvedByIds)
	if cErr != nil {
		return cErr
	}

	approvedByIdToSystemUserIdMap := map[string]string{}
	for _, user := range users {
		approvedByIdToSystemUserIdMap[fmt.Sprint(user.Id)] = user.SystemUserId
	}

	for index := range investigations {
		originalApprovedBy := investigations[index].ApprovedBy
		if systemUserId, exists := approvedByIdToSystemUserIdMap[originalApprovedBy]; exists && systemUserId != "" {
			investigations[index].ApprovedBy = systemUserId
		} else {
			investigations[index].ApprovedBy = constants.DEFAULT_SYSTEM_USER_ID
		}
	}

	// Map investigations to test_id in reportGenerationEvent
	s.mapInvestigationsAndTests(omsOrderId, orderDetails.CityCode, &reportGenerationEvent, investigations)
	reportIndexDetails, cErr := s.getIndexDetails(ctx, orderDetails, false)
	if cErr != nil {
		return cErr
	}
	reportGenerationEvent.IndexDetails = reportIndexDetails
	commonUtils.AddLog(ctx, commonConstants.DEBUG_LEVEL, "ReportGenerationEvent", map[string]interface{}{
		"reportGenerationEvent": reportGenerationEvent,
	}, nil)

	messageBody, messageAttributes := s.PubsubService.GetReportGenerationEvent(ctx, reportGenerationEvent)
	if messageBody == nil {
		commonUtils.AddLog(ctx, commonConstants.DEBUG_LEVEL, commonConstants.ERROR_NO_REPORT_GENERATION_EVENT_FOUND, nil, nil)
		return &commonStructures.CommonError{
			StatusCode: http.StatusNotFound,
			Message:    commonConstants.ERROR_NO_REPORT_GENERATION_EVENT_FOUND,
		}
	}
	cErr = s.SnsClient.PublishTo(ctx, messageBody, messageAttributes, fmt.Sprint(omsOrderId),
		commonConstants.OrderReportUpdateTopicArn, "")
	if cErr != nil {
		return cErr
	}

	return nil
}

func (s *ReportGenerationService) fetchTestIdsFromReportGenerationEvent(
	reportGenerationEvent structures.ReportGenerationEvent) []string {
	testIds := []string{}
	for _, visit := range reportGenerationEvent.Visits {
		for _, test := range visit.Tests {
			testIds = append(testIds, test.Id)
		}
	}
	return testIds
}

func (s *ReportGenerationService) fetchReportGenerationEventDetails(ctx context.Context,
	orderDetails commonModels.OrderDetails, omsTestIds []string) (
	structures.ReportGenerationEvent, structures.ReportGenerationEvent, *commonStructures.CommonError) {
	visits, allVisits, orderEvent := []structures.VisitEvent{}, []structures.VisitEvent{}, structures.OrderEvent{}
	reportGenerationEvent, reportGenerationEventAllVisits :=
		structures.ReportGenerationEvent{}, structures.ReportGenerationEvent{}
	var cErr *commonStructures.CommonError
	var cErrList []*commonStructures.CommonError
	wg := sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()
		visits, allVisits, cErr = s.fetchVisitEventDetails(ctx, orderDetails.OmsOrderId, omsTestIds)
		if cErr != nil {
			cErrList = append(cErrList, cErr)
		}
	}()

	go func() {
		defer wg.Done()
		orderEvent, cErr = s.fetchOrderEventDetails(ctx, orderDetails)
		if cErr != nil {
			cErrList = append(cErrList, cErr)
		}
	}()

	wg.Wait()
	if len(cErrList) > 0 {
		return reportGenerationEvent, reportGenerationEventAllVisits, cErrList[0]
	}

	reportGenerationEvent = structures.ReportGenerationEvent{
		Order:  orderEvent,
		Visits: visits,
	}

	reportGenerationEventAllVisits = structures.ReportGenerationEvent{
		Order:  orderEvent,
		Visits: allVisits,
	}

	return reportGenerationEvent, reportGenerationEventAllVisits, nil
}

func (s *ReportGenerationService) fetchAttachmentDetails(omsOrderId string) ([]commonModels.Attachment, *commonStructures.CommonError) {
	task, cErr := s.TaskService.GetTaskByOmsOrderId(omsOrderId)
	if cErr != nil {
		return []commonModels.Attachment{}, cErr
	}
	if task.Id == 0 {
		return []commonModels.Attachment{}, &commonStructures.CommonError{
			StatusCode: http.StatusNotFound,
			Message:    commonConstants.ERROR_TASK_NOT_FOUND,
		}
	}
	attachments, cErr := s.AttachmentService.GetAttachmentDtosByTaskId(task.Id, []string{commonConstants.AttachmentTypeTestDocument})
	if cErr != nil {
		return []commonModels.Attachment{}, cErr
	}
	return attachments, nil
}

func (s *ReportGenerationService) fetchVisitEventDetails(ctx context.Context, omsOrderId string, omsTestIds []string) (
	[]structures.VisitEvent, []structures.VisitEvent, *commonStructures.CommonError) {
	visits, allVisits, uniqueMasterTestIdsMap := []structures.VisitEvent{}, []structures.VisitEvent{}, map[uint]bool{}
	visitIdToTestEventMap, visitIdToTestEventMapForAllVisits :=
		map[string][]structures.TestEvent{}, map[string][]structures.TestEvent{}

	masterTestsMap := s.CdsService.GetMasterTestAsMap(ctx)
	labIdLabMap := s.CdsService.GetLabIdLabMap(ctx)

	visitIdTests, cErr := s.Dao.FetchVisitToTestMapping(omsOrderId, omsTestIds)
	if cErr != nil {
		return visits, allVisits, cErr
	}
	if len(visitIdTests) == 0 {
		return visits, allVisits, &commonStructures.CommonError{
			StatusCode: http.StatusBadRequest,
			Message:    commonConstants.ERROR_NO_VISIT_TESTS_FOUND,
		}
	}
	filteredVisitIdTests, allVisitIdTests := []structures.VisitTestsStruct{}, []structures.VisitTestsStruct{}
	for _, visitIdTest := range visitIdTests {
		if !labIdLabMap[visitIdTest.ProcessingLabId].Inhouse {
			continue
		}
		inhouseReportEnabled := masterTestsMap[visitIdTest.MasterTestId].TestLabMeta[visitIdTest.LabId].InhouseReportEnabled
		visitIdTest.DepartmentName = masterTestsMap[visitIdTest.MasterTestId].Department
		allVisitIdTests = append(allVisitIdTests, visitIdTest)
		if inhouseReportEnabled {
			filteredVisitIdTests = append(filteredVisitIdTests, visitIdTest)
		}
	}

	// if len(filteredVisitIdTests) == 0 {
	// 	return visits, allVisits, &commonStructures.CommonError{
	// 		StatusCode: http.StatusBadRequest,
	// 		Message:    commonConstants.ERROR_NO_INHOUSE_VISIT_TESTS_FOUND,
	// 	}
	// }

	vialTypeMap := s.CdsService.GetMasterVialTypeAsMap(ctx)
	for index := range filteredVisitIdTests {
		filteredVisitIdTests[index].SampleName = vialTypeMap[filteredVisitIdTests[index].VialTypeId].SampleName
	}

	// Map tests to visit_id
	for _, visitIdTest := range filteredVisitIdTests {
		// make sure if the same master_test_id is present always send any one
		if uniqueMasterTestIdsMap[visitIdTest.MasterTestId] {
			continue
		}
		uniqueMasterTestIdsMap[visitIdTest.MasterTestId] = true

		visitId := visitIdTest.VisitId
		if _, ok := visitIdToTestEventMap[visitId]; !ok {
			visitIdToTestEventMap[visitId] = []structures.TestEvent{}
		}
		visitIdToTestEventMap[visitId] = append(visitIdToTestEventMap[visitId], mapTestEvent(visitIdTest))
	}

	uniqueMasterTestIdsMap = map[uint]bool{}
	for _, visitIdTest := range allVisitIdTests {
		if uniqueMasterTestIdsMap[visitIdTest.MasterTestId] {
			continue
		}
		uniqueMasterTestIdsMap[visitIdTest.MasterTestId] = true

		visitId := visitIdTest.VisitId
		if _, ok := visitIdToTestEventMapForAllVisits[visitId]; !ok {
			visitIdToTestEventMapForAllVisits[visitId] = []structures.TestEvent{}
		}
		visitIdToTestEventMapForAllVisits[visitId] = append(visitIdToTestEventMapForAllVisits[visitId],
			mapTestEvent(visitIdTest))
	}

	for key, value := range visitIdToTestEventMap {
		if len(value) == 0 {
			continue
		}
		visit := structures.VisitEvent{
			Id:    key,
			Date:  value[0].SampleReceivedAt,
			LabId: value[0].LabId,
			Tests: value,
		}
		visits = append(visits, visit)
	}

	for key, value := range visitIdToTestEventMapForAllVisits {
		if len(value) == 0 {
			continue
		}
		visit := structures.VisitEvent{
			Id:    key,
			Date:  value[0].SampleReceivedAt,
			LabId: value[0].LabId,
			Tests: value,
		}
		allVisits = append(allVisits, visit)
	}

	return visits, allVisits, nil
}

func (s *ReportGenerationService) fetchOrderEventDetails(ctx context.Context, orderDetails commonModels.OrderDetails) (
	structures.OrderEvent, *commonStructures.CommonError) {
	orderEvent := structures.OrderEvent{}

	patientDetails, cErr := s.Dao.GetPatientDetailsById(orderDetails.PatientDetailsId)
	if cErr != nil {
		return orderEvent, cErr
	}

	patientAge := patientDetails.Dob
	if patientAge == nil {
		patientAge = patientDetails.ExpectedDob
	}

	patientAgeYears, patientAgeMonths, patientAgeDays := commonUtils.GetAgeYearsMonthsAndDaysFromDob(*patientAge)
	patientDetailsEvent := structures.PatientDetailsEvent{
		Id:        patientDetails.SystemPatientId,
		Name:      patientDetails.Name,
		Number:    patientDetails.Number,
		Gender:    patientDetails.Gender,
		AgeYears:  float32(patientAgeYears),
		AgeMonths: patientAgeMonths,
		AgeDays:   patientAgeDays,
	}
	if patientDetailsEvent.Id == "" {
		patientDetailsEvent.Id = "OH" + fmt.Sprint(orderDetails.OmsOrderId)
	}

	orderEvent.Id = orderDetails.OmsOrderId
	orderEvent.AlnumOrderId = orderDetails.OmsOrderId
	orderEvent.Token = orderDetails.Uuid
	orderEvent.PartnerId = orderDetails.PartnerId
	orderEvent.RequestId = orderDetails.OmsRequestId
	orderEvent.CreatedOn = orderDetails.CreatedAt.Format(commonConstants.DateTimeInSecLayout)
	orderEvent.PatientDetails = patientDetailsEvent
	orderEvent.ReferredBy = orderDetails.ReferredBy
	orderEvent.ServicingLabId = orderDetails.ServicingLabId

	if orderDetails.DoctorId != 0 {
		doctor, err := s.HealthApiClient.GetDoctorById(ctx, orderDetails.DoctorId)
		if err != nil {
			commonUtils.AddLog(ctx, commonConstants.ERROR_LEVEL, commonConstants.ERROR_WHILE_GETTING_DOCTORS, nil, err)
		}
		orderEvent.ReferredBy = doctor.Name
	}

	servicingLab, cErr := s.CdsService.GetLabById(ctx, orderDetails.ServicingLabId)
	if cErr != nil {
		return orderEvent, cErr
	}
	orderEvent.CityCode = servicingLab.City

	return orderEvent, nil
}

func mapTestEvent(visitIdTest structures.VisitTestsStruct) structures.TestEvent {
	if visitIdTest.CollectedAt == "" {
		visitIdTest.CollectedAt = visitIdTest.ReceivedAt
	}
	return structures.TestEvent{
		Id:                visitIdTest.CentralOmsTestId,
		MasterTestId:      visitIdTest.MasterTestId,
		Name:              visitIdTest.TestName,
		IsPanel:           visitIdTest.TestType == commonConstants.GroupShortHand,
		LabId:             visitIdTest.LabId,
		SampleType:        visitIdTest.SampleName,
		SampleId:          visitIdTest.Barcode,
		SampleCollectedAt: visitIdTest.CollectedAt,
		SampleReceivedAt:  visitIdTest.ReceivedAt,
		DepartmerntName:   visitIdTest.DepartmentName,
	}
}

func (s *ReportGenerationService) GetTestApprovedAt(currentTestApprovedAt string,
	investigationApprovedAt *time.Time) string {
	if investigationApprovedAt == nil {
		return currentTestApprovedAt
	}

	if currentTestApprovedAt == "" ||
		investigationApprovedAt.After(commonUtils.StringToTime(currentTestApprovedAt,
			commonConstants.DateTimeUTCLayoutWithoutTZOffset)) {
		return investigationApprovedAt.Format(commonConstants.DateTimeUTCLayoutWithoutTZOffset)
	}

	return currentTestApprovedAt
}

func (s *ReportGenerationService) mapInvestigationsAndTests(omsOrderId string, cityCode string,
	reportGenerationEvent *structures.ReportGenerationEvent,
	investigations []structures.InvestigationEvent) {

	// Map investigations to test_id
	investigationsMap := map[string][]structures.InvestigationEvent{}
	for _, investigation := range investigations {
		investigation.OrderId = omsOrderId
		if investigation.Data != "" {
			investigation.Value = investigation.Data
			investigation.Data = ""
		}
		investigationsMap[investigation.TestId] = append(investigationsMap[investigation.TestId], investigation)
	}

	// Map investigations to test_id in reportGenerationEvent
	for visitIndex, visit := range reportGenerationEvent.Visits {
		for testIndex, test := range visit.Tests {
			if investigations, ok := investigationsMap[test.Id]; ok {
				reportGenerationEvent.Visits[visitIndex].Tests[testIndex].Investigations = investigations
				testApprovedAt := s.GetTestApprovedAt(test.ApprovedAt, investigations[0].ApprovedAt)
				testApprovedBy := investigations[0].ApprovedBy
				reportReleasedBy := getDutyDoctorUserId(investigations[0].ApprovedAt, cityCode)
				reportGenerationEvent.Visits[visitIndex].Tests[testIndex].ApprovedAt = testApprovedAt
				reportGenerationEvent.Visits[visitIndex].Tests[testIndex].ReportReleasedAt = testApprovedAt
				reportGenerationEvent.Visits[visitIndex].Tests[testIndex].ApprovedBy = testApprovedBy
				reportGenerationEvent.Visits[visitIndex].Tests[testIndex].ReportReleasedBy = reportReleasedBy
			}
		}
	}
}

func getDutyDoctorUserId(approvedAt *time.Time, cityCode string) int {
	dutyDoctorMap := commonConstants.DutyDoctorMap

	if approvedAt == nil {
		return 0
	}

	istLocation, _ := time.LoadLocation(commonConstants.LocalTimeZoneLocation)
	approvedTime := approvedAt.In(istLocation)
	approvedDayStr := strings.ToLower(approvedTime.Weekday().String())

	dutyDoctor := dutyDoctorMap[approvedDayStr].(map[string]interface{})
	dutyDoctorPreviousDay :=
		dutyDoctorMap[strings.ToLower(approvedTime.AddDate(0, 0, -1).Weekday().String())].(map[string]interface{})

	approvedDayStartTime := dutyDoctor["shift_1"].(map[string]interface{})["start_time"].(string)
	startTime, _ := time.Parse(commonConstants.TimeLayout, approvedDayStartTime)

	if approvedTime.Hour() > 00 && approvedTime.Hour() < startTime.Hour() {
		dutyDoctor = dutyDoctorPreviousDay
	}

	for _, shift := range dutyDoctor {
		doctorUserId := getDoctorUserId(shift, approvedTime, cityCode)
		if doctorUserId != "" {
			intDoctorUserId, _ := strconv.Atoi(doctorUserId)
			return intDoctorUserId
		}
	}

	return 0
}

func getDoctorUserId(shift interface{}, approvedTime time.Time, cityCode string) string {
	shiftMap := shift.(map[string]interface{})
	startTimeStr := shiftMap["start_time"].(string)
	endTimeStr := shiftMap["end_time"].(string)
	doctorNameMap := shiftMap["doctor"].(map[string]interface{})
	startTime, _ := time.Parse(commonConstants.TimeLayout, startTimeStr)
	endTime, _ := time.Parse(commonConstants.TimeLayout, endTimeStr)

	startTime = time.Date(approvedTime.Year(), approvedTime.Month(), approvedTime.Day(), startTime.Hour(),
		startTime.Minute(), startTime.Second(), startTime.Nanosecond(), approvedTime.Location())
	endTime = time.Date(approvedTime.Year(), approvedTime.Month(), approvedTime.Day(), endTime.Hour(), endTime.Minute(),
		endTime.Second(), endTime.Nanosecond(), approvedTime.Location())
	approvedTime = time.Date(approvedTime.Year(), approvedTime.Month(), approvedTime.Day(), approvedTime.Hour(),
		approvedTime.Minute(), approvedTime.Second(), 0, approvedTime.Location())

	if endTime.Before(startTime) {
		endTime = endTime.AddDate(0, 0, +1)
	}

	commonUtils.AddLog(context.Background(), commonConstants.DEBUG_LEVEL, commonUtils.GetCurrentFunctionName(),
		map[string]interface{}{
			"start_time":    startTime.String(),
			"end_time":      endTime.String(),
			"approved_time": approvedTime.String(),
		}, nil)

	if (approvedTime.After(startTime) || approvedTime.Equal(startTime)) && approvedTime.Before(endTime) {
		doctorUserId := doctorNameMap[strings.ToLower(cityCode)]
		if doctorUserId != nil {
			return fmt.Sprint(doctorUserId)
		}
	}
	return ""
}

func rearrangeTestDetailsForIndexPage(testDetails []commonModels.TestDetail) []commonModels.TestDetail {
	if len(testDetails) == 0 {
		return testDetails
	}

	completedTestDetails, notCompletedTestDetails := []commonModels.TestDetail{}, []commonModels.TestDetail{}
	for _, testDetail := range testDetails {
		if commonUtils.SliceContainsString(commonConstants.OmsCompletedTestStatuses, testDetail.OmsStatus) ||
			testDetail.Status == commonConstants.TEST_STATUS_APPROVE {
			completedTestDetails = append(completedTestDetails, testDetail)
		} else {
			notCompletedTestDetails = append(notCompletedTestDetails, testDetail)
		}
	}

	newTestDetails := []commonModels.TestDetail{}
	if len(completedTestDetails) > 0 {
		newTestDetails = append(newTestDetails, completedTestDetails...)
	}
	if len(notCompletedTestDetails) > 0 {
		newTestDetails = append(newTestDetails, notCompletedTestDetails...)
	}
	return newTestDetails
}

func (s *ReportGenerationService) getIndexDetails(ctx context.Context, orderDetails commonModels.OrderDetails, isDummyReport bool) (
	structures.ReportIndexDetails, *commonStructures.CommonError) {
	hasIndividualTests := false
	testDetails, cErr := s.TestDetailService.GetTestDetailsByOmsOrderId(orderDetails.OmsOrderId)
	if cErr != nil {
		return structures.ReportIndexDetails{}, cErr
	}

	masterTestIds := []uint{}
	masterPackageIds := []uint{}

	for _, testDetail := range testDetails {
		masterTestIds = append(masterTestIds, testDetail.MasterTestId)
	}
	dedupResponse, err := s.CdsService.GetDeduplicatedTestsAndPackages(ctx, masterTestIds, []uint{})
	if err != nil {
		commonUtils.AddLog(ctx, commonConstants.ERROR_LEVEL, commonUtils.GetCurrentFunctionName(), nil, err)
	}

	dedupTestMap := map[uint]uint{}
	for _, removableTest := range dedupResponse.Recommendation.RemoveTests {
		if len(removableTest.CompleteOverlapDetails) > 0 {
			dedupTestMap[removableTest.TestId] = removableTest.CompleteOverlapDetails[0].OverlappedTestId
		}
	}

	masterTestIdToReportIndexStatusMap := map[uint]uint{}
	masterTestIdToReportIndexPresentMap := map[uint]bool{}

	reportIndexTestDetails := []structures.ReportTestDetails{}
	testDetails = rearrangeTestDetailsForIndexPage(testDetails)
	for _, testDetail := range testDetails {
		if masterTestIdToReportIndexPresentMap[testDetail.MasterTestId] {
			continue
		}
		holdReason := ""
		if testDetail.Status == commonConstants.TEST_STATUS_REJECTED {
			holdReason = constants.RECOLLECTION_REQUIRED
		} else if testDetail.Status == commonConstants.TEST_STATUS_SAMPLE_NOT_RECEIVED {
			holdReason = constants.SAMPLE_NOT_RECEIVED
		} else if testDetail.Status == commonConstants.TEST_STATUS_REQUESTED && testDetail.ReportEta == nil {
			holdReason = constants.SAMPLE_NOT_COLLECTED
		}
		mappingId := uint(0)
		if !testDetail.IsManualReportUpload && (isDummyReport || testDetail.LabId == orderDetails.ServicingLabId) {
			if _, ok := dedupTestMap[testDetail.MasterTestId]; ok {
				continue
			} else {
				mappingId = testDetail.MasterTestId
			}
		}

		reportEta := ""
		if testDetail.ReportEta != nil {
			reportEta = testDetail.ReportEta.Format(commonConstants.DateTimeInSecLayout)
		}
		testStatus := commonConstants.OmsTestStatusToUintMap[testDetail.OmsStatus]
		if !commonUtils.SliceContainsUint(commonConstants.OmsCompletedTestStatusesUint, testStatus) &&
			testDetail.Status == commonConstants.TEST_STATUS_APPROVE {
			testStatus = commonConstants.TestStatusCompletedNotSentUint
		}
		masterTestIdToReportIndexStatusMap[testDetail.MasterTestId] = testStatus
		reportIndexTestDetail := structures.ReportTestDetails{
			MasterTestId:   testDetail.MasterTestId,
			TestStatus:     testStatus,
			TestEta:        reportEta,
			TestHoldReason: holdReason,
			MappingId:      mappingId,
		}
		reportIndexTestDetail.IsNew = testDetail.Status == commonConstants.TEST_STATUS_APPROVE &&
			testDetail.ReportStatus != commonConstants.TEST_REPORT_STATUS_DELIVERED
		reportIndexTestDetails = append(reportIndexTestDetails, reportIndexTestDetail)
		masterTestIdToReportIndexPresentMap[testDetail.MasterTestId] = true
		hasIndividualTests = hasIndividualTests || testDetail.MasterPackageId == 0
		if testDetail.MasterPackageId != 0 && !commonUtils.SliceContainsUint(masterPackageIds, testDetail.MasterPackageId) {
			masterPackageIds = append(masterPackageIds, testDetail.MasterPackageId)
		}
	}

	for _, testDetail := range testDetails {
		if masterTestIdToReportIndexPresentMap[testDetail.MasterTestId] {
			continue
		}
		holdReason := ""
		if testDetail.Status == commonConstants.TEST_STATUS_REJECTED {
			holdReason = constants.RECOLLECTION_REQUIRED
		} else if testDetail.Status == commonConstants.TEST_STATUS_SAMPLE_NOT_RECEIVED {
			holdReason = constants.SAMPLE_NOT_RECEIVED
		} else if testDetail.Status == commonConstants.TEST_STATUS_REQUESTED && testDetail.ReportEta == nil {
			holdReason = constants.SAMPLE_NOT_COLLECTED
		}
		mappingId := uint(0)
		if !testDetail.IsManualReportUpload && (isDummyReport || testDetail.LabId == orderDetails.ServicingLabId) {
			if val, ok := dedupTestMap[testDetail.MasterTestId]; ok {
				mappingId = val
			}
		}

		reportEta := ""
		if testDetail.ReportEta != nil {
			reportEta = testDetail.ReportEta.Format(commonConstants.DateTimeInSecLayout)
		}
		reportIndexTestDetail := structures.ReportTestDetails{
			MasterTestId:   testDetail.MasterTestId,
			TestStatus:     masterTestIdToReportIndexStatusMap[mappingId],
			TestEta:        reportEta,
			TestHoldReason: holdReason,
			MappingId:      mappingId,
		}
		reportIndexTestDetail.IsNew = testDetail.Status == commonConstants.TEST_STATUS_APPROVE &&
			testDetail.ReportStatus != commonConstants.TEST_REPORT_STATUS_DELIVERED
		reportIndexTestDetails = append(reportIndexTestDetails, reportIndexTestDetail)
		masterTestIdToReportIndexPresentMap[testDetail.MasterTestId] = true
		hasIndividualTests = hasIndividualTests || testDetail.MasterPackageId == 0
		if testDetail.MasterPackageId != 0 && !commonUtils.SliceContainsUint(masterPackageIds, testDetail.MasterPackageId) {
			masterPackageIds = append(masterPackageIds, testDetail.MasterPackageId)
		}
	}

	return structures.ReportIndexDetails{
		IsPackageOrder: !hasIndividualTests,
		PackageIds:     masterPackageIds,
		TestDetails:    reportIndexTestDetails,
	}, nil
}

func (s *ReportGenerationService) TriggerOrderApprovedEvent(ctx context.Context,
	omsOrderId string, isDummyReport bool) (interface{}, *commonStructures.CommonError) {
	orderDetails, cErr := s.OrderDetailsService.GetOrderDetailsByOmsOrderId(omsOrderId)
	if cErr != nil {
		return nil, cErr
	}
	if orderDetails.Id == 0 {
		return nil, &commonStructures.CommonError{
			StatusCode: http.StatusNotFound,
			Message:    commonConstants.ERROR_ORDER_ID_NOT_FOUND,
		}
	}

	_, reportGenerationEventAllVisits, cErr := s.fetchReportGenerationEventDetails(ctx, orderDetails, []string{})
	if cErr != nil {
		return nil, cErr
	}
	var payload interface{}
	reportGenerationEventAllVisits.IsDummyReport = isDummyReport
	reportGenerationEventAllVisits.ServicingCityCode = orderDetails.CityCode
	attachments, cErr := s.fetchAttachmentDetails(omsOrderId)
	if cErr != nil {
		return nil, cErr
	}
	payload, cErr = s.GetAndPublishReportGenerationEventForAllVisits(ctx, orderDetails, reportGenerationEventAllVisits, attachments, isDummyReport)
	if cErr != nil {
		return nil, cErr
	}

	return payload, nil
}

func (s *ReportGenerationService) GetAndPublishReportGenerationEventForAllVisits(ctx context.Context,
	orderDetails commonModels.OrderDetails, reportGenerationEvent structures.ReportGenerationEvent, attachments []commonModels.Attachment, isDummyReport bool) (interface{}, *commonStructures.CommonError) {
	omsTestIds := s.fetchTestIdsFromReportGenerationEvent(reportGenerationEvent)

	investigations, cErr := s.Dao.GetInvestigationDataForReportGeneration(ctx, orderDetails.OmsOrderId, omsTestIds)
	if cErr != nil {
		return nil, cErr
	}

	attachmentsMap := make(map[uint][]string)
	for _, attachment := range attachments {
		if attachment.AttachmentType == commonConstants.AttachmentTypeTestDocument && attachment.InvestigationResultId > 0 {
			attachmentsMap[attachment.InvestigationResultId] = append(attachmentsMap[attachment.InvestigationResultId], attachment.AttachmentUrl)
		}
	}

	approvedByIds := []string{}
	for index := range investigations {
		if attachmentUrls, exists := attachmentsMap[investigations[index].Id]; exists && len(attachmentUrls) > 0 {
			investigations[index].TestDocument = append(investigations[index].TestDocument, attachmentUrls...)
		}
		approvedByIds = append(approvedByIds, investigations[index].ApprovedBy)
	}

	approvedByIds = commonUtils.CreateUniqueSliceString(approvedByIds)

	users, cErr := s.Dao.GetUserDetails(ctx, approvedByIds)
	if cErr != nil {
		return nil, cErr
	}

	approvedByIdToSystemUserIdMap := map[string]string{}
	for _, user := range users {
		approvedByIdToSystemUserIdMap[fmt.Sprint(user.Id)] = user.SystemUserId
	}

	for index := range investigations {
		originalApprovedBy := investigations[index].ApprovedBy
		if systemUserId, exists := approvedByIdToSystemUserIdMap[originalApprovedBy]; exists && systemUserId != "" {
			investigations[index].ApprovedBy = systemUserId
		} else {
			investigations[index].ApprovedBy = constants.DEFAULT_SYSTEM_USER_ID
		}
	}
	// Map investigations to test_id in reportGenerationEvent
	s.mapInvestigationsAndTests(orderDetails.OmsOrderId, orderDetails.CityCode, &reportGenerationEvent, investigations)
	reportIndexDetails, cErr := s.getIndexDetails(ctx, orderDetails, isDummyReport)
	if cErr != nil {
		return nil, cErr
	}
	reportGenerationEvent.IndexDetails = reportIndexDetails

	messageBody, messageAttributes := s.PubsubService.GetReportGenerationEvent(ctx, reportGenerationEvent)
	if messageBody == nil {
		return nil, &commonStructures.CommonError{
			StatusCode: http.StatusNotFound,
			Message:    commonConstants.ERROR_NO_REPORT_GENERATION_EVENT_FOUND,
		}
	}

	if isDummyReport {
		cErr = s.SnsClient.PublishTo(ctx, messageBody, messageAttributes, fmt.Sprint(orderDetails.OmsOrderId),
			commonConstants.OrderReportUpdateTopicArn, "")
		if cErr != nil {
			return nil, cErr
		}
		return messageBody, nil
	}
	cErr = s.SnsClient.PublishTo(ctx, messageBody, messageAttributes, fmt.Sprint(orderDetails.OmsOrderId),
		commonConstants.OrderResultsApprovedTopicArn, "")
	if cErr != nil {
		return nil, cErr
	}

	return nil, nil
}
