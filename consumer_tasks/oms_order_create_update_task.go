package consumerTasks

import (
	"context"
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

func (eventProcessor *EventProcessor) OmsOrderCreateUpdateEventTask(ctx context.Context,
	eventPayload string) error {

	omsOrderCreateUpdateEvent := structures.OmsOrderCreateUpdateEvent{}
	err := json.Unmarshal(json.RawMessage(eventPayload), &omsOrderCreateUpdateEvent)
	if err != nil {
		eventProcessor.Sentry.LogError(ctx, constants.ERROR_FAILED_TO_UNMARSHAL_JSON, err, nil)
		return err
	}
	utils.AddLog(ctx, constants.INFO_LEVEL, utils.GetCurrentFunctionName(),
		map[string]interface{}{
			"event_payload": eventPayload,
			"event_type":    constants.OmsOrderCreatedEvent,
		}, nil)

	if validateOmsOrderCreateUpdateEvent(omsOrderCreateUpdateEvent) != nil {
		utils.AddLog(ctx, constants.DEBUG_LEVEL, utils.GetCurrentFunctionName(),
			map[string]interface{}{"orderEvent": omsOrderCreateUpdateEvent},
			errors.New(constants.ERROR_INVALID_OMS_ORDER_CREATE_UPDATE_EVENT),
		)
		return nil
	}

	redisKey := fmt.Sprintf(constants.OmsCreateUpdateOrderEventKey, omsOrderCreateUpdateEvent.Order.AlnumOrderId)
	keyExists, err := eventProcessor.Cache.Exists(ctx, redisKey)
	if err != nil || keyExists {
		utils.AddLog(ctx, constants.DEBUG_LEVEL, utils.GetCurrentFunctionName(), nil, err)
		return errors.New(constants.ERROR_CREATE_UPDATE_ORDER_TASK_IN_PROGRESS)
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

	cErr := eventProcessor.ProcessOmsOrderCreateUpdateEvent(ctx, omsOrderCreateUpdateEvent)
	if cErr != nil {
		return errors.New(cErr.Message)
	}

	return nil
}

func validateOmsOrderCreateUpdateEvent(omsOrderEvent structures.OmsOrderCreateUpdateEvent) *structures.CommonError {
	if omsOrderEvent.CityCode == "" || omsOrderEvent.Order.AlnumOrderId == "" ||
		omsOrderEvent.Request.AlnumRequestId == "" {
		return &structures.CommonError{
			Message:    constants.ERROR_INVALID_OMS_ORDER_CREATE_UPDATE_EVENT,
			StatusCode: http.StatusBadRequest,
		}
	}

	for _, test := range omsOrderEvent.Tests {
		if test.AlnumTestId == "" || test.MasterTestId == 0 {
			return &structures.CommonError{
				Message:    constants.ERROR_INVALID_OMS_ORDER_CREATE_UPDATE_EVENT,
				StatusCode: http.StatusBadRequest,
			}
		}
	}

	if omsOrderEvent.Order.OriginalOrderId != 0 {
		return &structures.CommonError{
			Message:    constants.ERROR_INVALID_OMS_ORDER_CREATE_UPDATE_EVENT,
			StatusCode: http.StatusBadRequest,
		}
	}

	return nil
}

func (eventProcessor *EventProcessor) ProcessOmsOrderCreateUpdateEvent(ctx context.Context,
	omsOrderEvent structures.OmsOrderCreateUpdateEvent) *structures.CommonError {

	omsOrderId := omsOrderEvent.Order.AlnumOrderId
	orderDetails, cErr := eventProcessor.OrderDetailsService.GetOrderDetailsByOmsOrderId(omsOrderId)
	if cErr != nil && cErr.Message != constants.ERROR_ORDER_ID_NOT_FOUND {
		return cErr
	}

	masterTestsMap := eventProcessor.CdsService.GetMasterTestAsMap(ctx)
	outsourceLabIds := eventProcessor.CdsService.GetOutsourceLabIds(ctx)
	inhouseLabIds := eventProcessor.CdsService.GetInhouseLabIds(ctx)
	nrlCpEnabledMasterTestIds := eventProcessor.CdsService.GetNrlCpEnabledMasterTestIds(ctx)

	if orderDetails.Id == 0 {
		return eventProcessor.createOrderDetailsData(ctx, omsOrderEvent, masterTestsMap,
			outsourceLabIds, nrlCpEnabledMasterTestIds)
	} else {
		if orderDetails.ServicingLabId == omsOrderEvent.Request.ServicingLabId {
			return eventProcessor.createUpdateOrderDetailsData(ctx, omsOrderEvent, orderDetails,
				masterTestsMap, inhouseLabIds, outsourceLabIds, nrlCpEnabledMasterTestIds)
		} else {
			return eventProcessor.recreateOrderDetailsDataForServicingLabChange(ctx, omsOrderEvent, orderDetails,
				masterTestsMap, outsourceLabIds, nrlCpEnabledMasterTestIds)
		}
	}
}

func getCpEnabledTestFlagForTestDetails(testFinalProcessingLabId, masterTestId uint,
	outsourceLabIds, nrlCpEnabledMasterTestIds []uint) bool {
	if utils.SliceContainsUint(outsourceLabIds, testFinalProcessingLabId) {
		return false
	}
	if testFinalProcessingLabId != constants.NrlLabId {
		return true
	}

	if utils.SliceContainsUint(nrlCpEnabledMasterTestIds, masterTestId) {
		return true
	}

	return false
}

func getOmsTestStatus(omsTestStatus uint) string {
	if constants.TestStatusUintToStringMap[omsTestStatus] == "" {
		return constants.TEST_STATUS_REQUESTED
	}
	return constants.TestStatusUintToStringMap[omsTestStatus]
}

func getLabEta(labEta string) *time.Time {
	if labEta != "" {
		labEtaTime, err := time.Parse(constants.DateTimeUTCLayoutWithoutTZOffset, labEta)
		if err == nil {
			return &labEtaTime
		}
	}
	return nil
}

func getReportEta(reportEta string) *time.Time {
	if reportEta != "" {
		reportEtaTime, err := time.Parse(constants.DateTimeUTCLayoutWithoutTZOffset, reportEta)
		if err == nil {
			return &reportEtaTime
		}
	}
	return nil
}

func getFinalTestCodeAndLabIdByMasterTest(masterTest structures.CdsTestMaster, labId uint, testCode string,
	outsourceLabIds []uint) (string, uint) {
	if masterTest.Id == 0 {
		return testCode, labId
	}
	testCode = masterTest.TestLabMeta[labId].TestCode
	finalProcessingLabId := masterTest.TestLabMeta[labId].LabId
	for {
		if labMeta, exists := masterTest.TestLabMeta[labId]; exists {
			labMetaLabId := labMeta.LabId
			labId = labMetaLabId
			testCode = labMeta.TestCode
			finalProcessingLabId = labMeta.LabId
			if labMetaLabId == labId || utils.SliceContainsUint(outsourceLabIds, labMetaLabId) {
				break
			}
		} else {
			break
		}
	}

	return testCode, finalProcessingLabId
}

func createPatientDetailDtoForOmsOrderEvent(omsOrderEvent structures.OmsOrderCreateUpdateEvent) models.PatientDetail {
	patientExpectedDob := utils.GetDobByYearsMonthsDays(int(omsOrderEvent.Order.PatientAge),
		int(omsOrderEvent.Order.PatientAgeMonths), int(omsOrderEvent.Order.PatientAgeDays))
	patientDetails := models.PatientDetail{
		Name:            omsOrderEvent.Order.PatientName,
		ExpectedDob:     &patientExpectedDob,
		Gender:          strings.ToLower(omsOrderEvent.Order.PatientGender),
		Number:          omsOrderEvent.Order.PatientNumber,
		SystemPatientId: omsOrderEvent.Order.PatientId,
	}
	patientDetails.CreatedBy = constants.CitadelSystemId
	patientDetails.UpdatedBy = constants.CitadelSystemId

	return patientDetails
}

func createOrderDetailsDtoForOmsOrderEvent(omsOrderEvent structures.OmsOrderCreateUpdateEvent,
	patientDetails models.PatientDetail) models.OrderDetails {
	orderDetails := models.OrderDetails{
		OmsOrderId:       omsOrderEvent.Order.AlnumOrderId,
		OmsRequestId:     omsOrderEvent.Request.AlnumRequestId,
		Uuid:             omsOrderEvent.Request.Uuid,
		CityCode:         omsOrderEvent.Request.ServicingCityCode,
		PatientDetailsId: patientDetails.Id,
		OrderStatus:      constants.OrderStatusMapUint[omsOrderEvent.Order.Status],
		PartnerId:        omsOrderEvent.Order.PartnerId,
		DoctorId:         omsOrderEvent.Order.SystemDoctorId,
		TrfId:            omsOrderEvent.Order.TrfId,
		ServicingLabId:   omsOrderEvent.Request.ServicingLabId,
		CollectionType:   omsOrderEvent.Request.CollectionType,
		RequestSource:    omsOrderEvent.Request.Source,
		BulkOrderId:      omsOrderEvent.Request.BulkOrderDetailId,
		CampId:           omsOrderEvent.Request.CampId,
		ReferredBy:       omsOrderEvent.Order.ReferredBy,
		SrfId:            omsOrderEvent.Order.SrfId,
	}
	orderDetails.CollectedOn = nil
	if omsOrderEvent.Order.CollectedOn != "" {
		collectedOnTime, err := time.Parse(constants.DateTimeUTCLayoutWithoutTZOffset, omsOrderEvent.Order.CollectedOn)
		if err == nil {
			orderDetails.CollectedOn = &collectedOnTime
		}
	}
	orderDetails.CreatedBy = constants.CitadelSystemId
	orderDetails.UpdatedBy = constants.CitadelSystemId

	return orderDetails
}

func createTestDetailsDtoForOmsOrderEvent(omsOrderEvent structures.OmsOrderCreateUpdateEvent,
	masterTestsMap map[uint]structures.CdsTestMaster,
	outsourceLabIds, nrlCpEnabledMasterTestIds []uint) []models.TestDetail {
	testDetails := make([]models.TestDetail, 0)
	for _, test := range omsOrderEvent.Tests {
		lisCode, finalProcessingLabId := getFinalTestCodeAndLabIdByMasterTest(masterTestsMap[test.MasterTestId], test.LabId,
			test.TestCode, outsourceLabIds)
		cpEnabledFlag := getCpEnabledTestFlagForTestDetails(finalProcessingLabId, test.MasterTestId, outsourceLabIds,
			nrlCpEnabledMasterTestIds)
		testDetail := models.TestDetail{
			OmsOrderId:       omsOrderEvent.Order.AlnumOrderId,
			OmsTestId:        test.Id,
			CentralOmsTestId: test.AlnumTestId,
			CityCode:         omsOrderEvent.CityCode,
			TestName:         test.TestName,
			LabId:            test.LabId,
			ProcessingLabId:  test.LabId,
			LisCode:          lisCode,
			MasterTestId:     test.MasterTestId,
			MasterPackageId:  test.MasterPackageId,
			TestType:         test.TestType,
			Department:       masterTestsMap[test.MasterTestId].Department,
			Status:           constants.TEST_STATUS_REQUESTED,
			OmsStatus:        getOmsTestStatus(test.Status),
			LabTat:           test.LabTat,
			LabEta:           getLabEta(test.LabEta),
			ReportEta:        getReportEta(test.ReportTat),
			ReportStatus:     constants.TEST_REPORT_STATUS_NOT_READY,
			CpEnabled:        cpEnabledFlag,
		}
		testDetail.CreatedBy = constants.CitadelSystemId
		testDetail.UpdatedBy = constants.CitadelSystemId

		testDetails = append(testDetails, testDetail)
	}

	return testDetails
}

func createTestDetailsMetadataDtoForOmsOrderEvent(testDetails []models.TestDetail) []models.TestDetailsMetadata {
	testDetailsMetadata := make([]models.TestDetailsMetadata, 0)
	for _, test := range testDetails {
		testDetailMetadata := models.TestDetailsMetadata{
			TestDetailsId: test.Id,
		}
		testDetailMetadata.CreatedBy = constants.CitadelSystemId
		testDetailMetadata.UpdatedBy = constants.CitadelSystemId
		testDetailsMetadata = append(testDetailsMetadata, testDetailMetadata)
	}

	return testDetailsMetadata
}

func createUpdatePatientDetailsDtoForOmsOrderEvent(omsOrderEvent structures.OmsOrderCreateUpdateEvent,
	patientDetails models.PatientDetail) models.PatientDetail {
	patientExpectedDob := utils.GetDobByYearsMonthsDays(int(omsOrderEvent.Order.PatientAge),
		int(omsOrderEvent.Order.PatientAgeMonths), int(omsOrderEvent.Order.PatientAgeDays))

	patientDetails.Gender = strings.ToLower(omsOrderEvent.Order.PatientGender)
	patientDetails.Name = omsOrderEvent.Order.PatientName
	patientDetails.ExpectedDob = &patientExpectedDob
	patientDetails.Number = omsOrderEvent.Order.PatientNumber
	patientDetails.SystemPatientId = omsOrderEvent.Order.PatientId
	patientDetails.UpdatedBy = constants.CitadelSystemId

	return patientDetails
}

func createUpdateOrderDetailsDtoForOmsOrderEvent(omsOrderEvent structures.OmsOrderCreateUpdateEvent,
	orderDetails models.OrderDetails) models.OrderDetails {

	orderDetails.OrderStatus = constants.OrderStatusMapUint[omsOrderEvent.Order.Status]
	orderDetails.Uuid = omsOrderEvent.Request.Uuid
	orderDetails.CityCode = omsOrderEvent.Request.ServicingCityCode
	orderDetails.TrfId = omsOrderEvent.Order.TrfId
	orderDetails.PartnerId = omsOrderEvent.Order.PartnerId
	orderDetails.DoctorId = omsOrderEvent.Order.SystemDoctorId
	orderDetails.ServicingLabId = omsOrderEvent.Request.ServicingLabId
	orderDetails.CollectionType = omsOrderEvent.Request.CollectionType
	orderDetails.RequestSource = omsOrderEvent.Request.Source
	orderDetails.BulkOrderId = omsOrderEvent.Request.BulkOrderDetailId
	orderDetails.UpdatedBy = constants.CitadelSystemId
	orderDetails.CampId = omsOrderEvent.Request.CampId
	orderDetails.ReferredBy = omsOrderEvent.Order.ReferredBy
	orderDetails.SrfId = omsOrderEvent.Order.SrfId
	orderDetails.CollectedOn = nil
	if omsOrderEvent.Order.CollectedOn != "" {
		collectedOnTime, err := time.Parse(constants.DateTimeUTCLayoutWithoutTZOffset, omsOrderEvent.Order.CollectedOn)
		if err == nil {
			orderDetails.CollectedOn = &collectedOnTime
		}
	}

	return orderDetails

}

func getCreateUpdateTestDetailsDtoForOmsOrderEvent(omsOrderEvent structures.OmsOrderCreateUpdateEvent,
	testDetails []models.TestDetail, task models.Task, masterTestsMap map[uint]structures.CdsTestMaster,
	inhouseLabIds, outsourceLabIds, nrlCpEnabledMasterTestIds []uint) (
	toCreateTestDetails, toUpdateTestDetails, labIdChangeTestDetails []models.TestDetail, toDeleteTestIds []string) {
	toCreateTestDetails = []models.TestDetail{}
	omsTestIds := []string{}
	for _, test := range omsOrderEvent.Tests {
		omsTestIds = append(omsTestIds, test.AlnumTestId)
	}

	testDetailsOmsIds := []string{}
	for _, testDetail := range testDetails {
		testDetailsOmsIds = append(testDetailsOmsIds, testDetail.CentralOmsTestId)
	}

	toCreateTestIds := utils.GetDifferenceBetweenStringSlices(omsTestIds, testDetailsOmsIds)
	toUpdateTestIds := utils.GetCommonElementsBetweenStringSlices(omsTestIds, testDetailsOmsIds)
	toDeleteTestIds = utils.GetDifferenceBetweenStringSlices(testDetailsOmsIds, omsTestIds)

	omsTestIdToUpdateTestsMap := map[string]models.TestDetail{}
	for _, testDetail := range testDetails {
		if utils.SliceContainsString(toUpdateTestIds, testDetail.CentralOmsTestId) {
			omsTestIdToUpdateTestsMap[testDetail.CentralOmsTestId] = testDetail
		}
	}

	for _, test := range omsOrderEvent.Tests {
		taskId := uint(0)
		if utils.SliceContainsUint(inhouseLabIds, test.LabId) {
			taskId = task.Id
		}
		if utils.SliceContainsString(toCreateTestIds, test.AlnumTestId) {
			lisCode, finalProcessingLabId := getFinalTestCodeAndLabIdByMasterTest(masterTestsMap[test.MasterTestId], test.LabId, test.TestCode,
				outsourceLabIds)
			cpEnabledFlag := getCpEnabledTestFlagForTestDetails(finalProcessingLabId, test.MasterTestId, outsourceLabIds,
				nrlCpEnabledMasterTestIds)
			testDetail := models.TestDetail{
				OmsOrderId:       omsOrderEvent.Order.AlnumOrderId,
				TaskId:           taskId,
				OmsTestId:        test.Id,
				CentralOmsTestId: test.AlnumTestId,
				CityCode:         omsOrderEvent.CityCode,
				TestName:         test.TestName,
				Status:           constants.TEST_STATUS_REQUESTED,
				OmsStatus:        getOmsTestStatus(test.Status),
				LabId:            test.LabId,
				ProcessingLabId:  test.LabId,
				LisCode:          lisCode,
				MasterTestId:     test.MasterTestId,
				MasterPackageId:  test.MasterPackageId,
				TestType:         test.TestType,
				Department:       masterTestsMap[test.MasterTestId].Department,
				LabTat:           test.LabTat,
				LabEta:           getLabEta(test.LabEta),
				ReportEta:        getReportEta(test.ReportTat),
				CpEnabled:        cpEnabledFlag,
			}
			testDetail.CreatedBy = constants.CitadelSystemId
			testDetail.UpdatedBy = constants.CitadelSystemId
			toCreateTestDetails = append(toCreateTestDetails, testDetail)
		} else if utils.SliceContainsString(toUpdateTestIds, test.AlnumTestId) {
			labIdChange := false
			lisCode, finalProcessingLabId := getFinalTestCodeAndLabIdByMasterTest(masterTestsMap[test.MasterTestId],
				test.LabId, test.TestCode, outsourceLabIds)
			cpEnabledFlag := getCpEnabledTestFlagForTestDetails(finalProcessingLabId, test.MasterTestId,
				outsourceLabIds, nrlCpEnabledMasterTestIds)
			testDetail := omsTestIdToUpdateTestsMap[test.AlnumTestId]
			testDetail.TaskId = taskId
			testDetail.TestName = test.TestName
			testDetail.OmsStatus = getOmsTestStatus(test.Status)
			if testDetail.LabId != test.LabId {
				testDetail.LabId = test.LabId
				testDetail.ProcessingLabId = test.LabId
				labIdChange = true
			}
			testDetail.LisCode = lisCode
			testDetail.MasterTestId = test.MasterTestId
			testDetail.MasterPackageId = test.MasterPackageId
			testDetail.TestType = test.TestType
			testDetail.Department = masterTestsMap[test.MasterTestId].Department
			testDetail.LabTat = test.LabTat
			testDetail.LabEta = getLabEta(test.LabEta)
			testDetail.CpEnabled = cpEnabledFlag
			testDetail.ReportEta = getReportEta(test.ReportTat)
			testDetail.UpdatedBy = constants.CitadelSystemId
			toUpdateTestDetails = append(toUpdateTestDetails, testDetail)
			if labIdChange {
				labIdChangeTestDetails = append(labIdChangeTestDetails, testDetail)
			}
		}
	}

	return
}

func (eventProcessor *EventProcessor) fetchCurrentDataForOmsOrderEvent(
	omsOrderEvent structures.OmsOrderCreateUpdateEvent, orderDetails models.OrderDetails) (
	models.Task, models.TaskMetadata, models.PatientDetail, []models.TestDetail, bool, *structures.CommonError) {
	taskMetadata, omsOrderId := models.TaskMetadata{}, omsOrderEvent.Order.AlnumOrderId
	task, cErr := eventProcessor.TaskService.GetTaskByOmsOrderId(omsOrderId)
	if cErr != nil {
		return task, models.TaskMetadata{}, models.PatientDetail{}, []models.TestDetail{}, false, cErr
	}

	if task.Id != 0 {
		taskMetadata, cErr = eventProcessor.TaskService.GetTaskMetadataByTaskId(task.Id)
		if cErr != nil {
			return task, taskMetadata, models.PatientDetail{}, []models.TestDetail{}, false, cErr
		}
	}

	patientDetails, cErr := eventProcessor.PatientDetailService.GetPatientDetailsById(orderDetails.PatientDetailsId)
	if cErr != nil {
		return task, taskMetadata, patientDetails, []models.TestDetail{}, false, cErr
	}

	testDetails, cErr := eventProcessor.TestDetailService.GetTestDetailsByOmsOrderId(omsOrderId)
	if cErr != nil && cErr.Message != constants.ERROR_NO_TEST_DETAILS_FOUND {
		return task, taskMetadata, patientDetails, testDetails, false, cErr
	}

	isSampleCollected, cErr := eventProcessor.SampleService.IsSampleCollected(omsOrderId)
	if cErr != nil {
		return task, taskMetadata, patientDetails, testDetails, false, cErr
	}

	return task, taskMetadata, patientDetails, testDetails, isSampleCollected, nil
}

func (eventProcessor *EventProcessor) createOrderDetailsData(ctx context.Context,
	omsOrderEvent structures.OmsOrderCreateUpdateEvent, masterTestsMap map[uint]structures.CdsTestMaster,
	outsourceLabIds, nrlCpEnabledMasterTestIds []uint) *structures.CommonError {

	// Start a Transaction and create all the details
	err := eventProcessor.Db.Transaction(func(tx *gorm.DB) error {
		// Create Patient Details
		patientDetails, cErr := eventProcessor.PatientDetailService.CreatePatientDetailsWithTx(tx,
			createPatientDetailDtoForOmsOrderEvent(omsOrderEvent))
		if cErr != nil {
			return errors.New(cErr.Message)
		}

		// Create Order Details
		orderDetails, cErr := eventProcessor.OrderDetailsService.CreateOrderDetailsWithTx(tx,
			createOrderDetailsDtoForOmsOrderEvent(omsOrderEvent, patientDetails))
		if cErr != nil {
			return errors.New(cErr.Message)
		}

		// Create Test Details
		testDetails, cErr := eventProcessor.TestDetailService.CreateTestDetailsWithTx(tx,
			createTestDetailsDtoForOmsOrderEvent(omsOrderEvent, masterTestsMap, outsourceLabIds, nrlCpEnabledMasterTestIds))
		if cErr != nil {
			return errors.New(cErr.Message)
		}

		// Create Test Details Metadata
		_, cErr = eventProcessor.TestDetailService.CreateTestDetailsMetadataWithTx(tx,
			createTestDetailsMetadataDtoForOmsOrderEvent(testDetails))
		if cErr != nil {
			return errors.New(cErr.Message)
		}

		// Create Samples
		if len(testDetails) > 0 {
			cErr = eventProcessor.SampleService.CreateSamplesWithTx(ctx, tx, omsOrderEvent.Order.AlnumOrderId)
			if cErr != nil {
				return errors.New(cErr.Message)
			}
		}

		if omsOrderEvent.SynchronizeTasks {
			cErr = eventProcessor.SampleService.SynchronizeTasksWithSamplesWithTx(tx, orderDetails.OmsRequestId,
				omsOrderEvent.Tasks, omsOrderEvent.TaskTestsMapping)
			if cErr != nil {
				utils.AddLog(ctx, constants.ERROR_LEVEL, utils.GetCurrentFunctionName(), nil, errors.New(cErr.Message))
			}
		}

		return nil
	})

	if err != nil {
		return &structures.CommonError{
			Message:    err.Error(),
			StatusCode: http.StatusInternalServerError,
		}
	}

	return nil
}

func (eventProcessor *EventProcessor) createUpdateOrderDetailsData(ctx context.Context,
	omsOrderEvent structures.OmsOrderCreateUpdateEvent, orderDetails models.OrderDetails,
	masterTestsMap map[uint]structures.CdsTestMaster,
	inhouseLabIds, outsourceLabIds, nrlCpEnabledMasterTestIds []uint) *structures.CommonError {

	toDeleteTestDetailsIds := []uint{}
	testIdLisSyncAtMap := map[string]*time.Time{}

	omsOrderId := omsOrderEvent.Order.AlnumOrderId
	toDeleteOmsTestIds := []string{}

	task, taskMetadata, patientDetails, testDetails, isSampleCollected, cErr :=
		eventProcessor.fetchCurrentDataForOmsOrderEvent(omsOrderEvent, orderDetails)
	if cErr != nil {
		utils.AddLog(ctx, constants.ERROR_LEVEL, utils.GetCurrentFunctionName(), nil, errors.New(cErr.Message))
		return cErr
	}

	if task.Id != 0 {
		taskMetadata.DoctorNotes = omsOrderEvent.Order.Notes
	}

	omsTestIdToTestDetailMap := map[string]models.TestDetail{}
	for _, testDetail := range testDetails {
		omsTestIdToTestDetailMap[testDetail.CentralOmsTestId] = testDetail
	}

	err := eventProcessor.Db.Transaction(func(tx *gorm.DB) error {
		// Update Patient Details
		patientDetails, cErr = eventProcessor.PatientDetailService.UpdatePatientDetailsWithTx(tx,
			createUpdatePatientDetailsDtoForOmsOrderEvent(omsOrderEvent, patientDetails))
		if cErr != nil {
			return errors.New(cErr.Message)
		}

		// Update Order Details
		orderDetails, cErr = eventProcessor.OrderDetailsService.UpdateOrderDetailsWithTx(tx,
			createUpdateOrderDetailsDtoForOmsOrderEvent(omsOrderEvent, orderDetails))
		if cErr != nil {
			return errors.New(cErr.Message)
		}

		if task.Id != 0 {
			// Update Task Metadata
			_, cErr = eventProcessor.TaskService.UpdateTaskMetadataWithTx(tx, taskMetadata)
			if cErr != nil {
				return errors.New(cErr.Message)
			}
		}

		// Update Test Details
		createTestDetails, updateTestDetails, labIdChangeTestDetails, toDeleteOmsTestIds :=
			getCreateUpdateTestDetailsDtoForOmsOrderEvent(omsOrderEvent, testDetails, task,
				masterTestsMap, inhouseLabIds, outsourceLabIds, nrlCpEnabledMasterTestIds)
		_, cErr = eventProcessor.TestDetailService.CreateTestDetailsWithTx(tx, createTestDetails)
		if cErr != nil {
			return errors.New(cErr.Message)
		}
		_, cErr = eventProcessor.TestDetailService.CreateTestDetailsMetadataWithTx(tx,
			createTestDetailsMetadataDtoForOmsOrderEvent(createTestDetails))
		if cErr != nil {
			return errors.New(cErr.Message)
		}

		_, cErr := eventProcessor.TestDetailService.UpdateTestDetailsWithTx(tx, updateTestDetails)
		if cErr != nil {
			return errors.New(cErr.Message)
		}

		if !isSampleCollected {
			if len(createTestDetails) > 0 || len(labIdChangeTestDetails) > 0 {
				cErr = eventProcessor.SampleService.CreateSamplesAndTestDetailsWithTx(ctx, tx, omsOrderId,
					createTestDetails, labIdChangeTestDetails)
				if cErr != nil {
					return errors.New(cErr.Message)
				}
			}
		} else {
			testIdLisSyncAtMap, cErr = eventProcessor.SampleService.CreateUpdateSamplesPostCollectionWithTx(ctx, tx,
				omsOrderId, createTestDetails, updateTestDetails, omsOrderEvent.Tests)
			if cErr != nil {
				return errors.New(cErr.Message)
			}
		}

		if !isSampleCollected || omsOrderEvent.IsPartialCancellationFlow {
			cErr = eventProcessor.TestDetailService.DeleteTestDetailsByOmsTestIdsWithTx(tx, toDeleteOmsTestIds)
			if cErr != nil {
				return errors.New(cErr.Message)
			}

			toDeleteTestDetails := []models.TestDetail{}
			for _, omsTestId := range toDeleteOmsTestIds {
				toDeleteTestDetails = append(toDeleteTestDetails, omsTestIdToTestDetailMap[omsTestId])
				toDeleteTestDetailsIds = append(toDeleteTestDetailsIds, omsTestIdToTestDetailMap[omsTestId].Id)
			}

			if len(toDeleteTestDetails) > 0 {
				cErr = eventProcessor.TestDetailService.DeleteTestDetailsMetadataByTestDetailsIdsWithTx(tx,
					toDeleteTestDetailsIds)
				if cErr != nil {
					return errors.New(cErr.Message)
				}

				cErr = eventProcessor.SampleService.DeleteTestSampleMappingForDeletedTestIds(ctx, tx, omsOrderId,
					toDeleteTestDetails)
				if cErr != nil {
					return errors.New(cErr.Message)
				}
			}
		}

		if omsOrderEvent.SynchronizeTasks {
			cErr = eventProcessor.SampleService.SynchronizeTasksWithSamplesWithTx(tx, orderDetails.OmsRequestId,
				omsOrderEvent.Tasks, omsOrderEvent.TaskTestsMapping)
			if cErr != nil {
				utils.AddLog(ctx, constants.ERROR_LEVEL, utils.GetCurrentFunctionName(), nil, errors.New(cErr.Message))
			}
		}

		return nil
	})

	if err != nil {
		return &structures.CommonError{
			Message:    err.Error(),
			StatusCode: http.StatusInternalServerError,
		}
	}

	if isSampleCollected {
		if len(testIdLisSyncAtMap) > 0 {
			centralOmsTestIds := []string{}
			for omsTestId, lisSyncAtTime := range testIdLisSyncAtMap {
				centralOmsTestIds = append(centralOmsTestIds, omsTestId)
				messageBody, messageAttributes := eventProcessor.PubsubService.GetLabEtaUpdateEvent(omsOrderId,
					[]string{omsTestId}, lisSyncAtTime, orderDetails.CityCode)
				eventProcessor.SnsClient.PublishTo(ctx, messageBody, messageAttributes, "", constants.OmsUpdatesTopicArn, "")
			}

			cErr = eventProcessor.TestDetailService.UpdateTaskIdInTestDetailsWithOmsTestIds(centralOmsTestIds, task.Id)
			if cErr != nil {
				utils.AddLog(ctx, constants.ERROR_LEVEL, utils.GetCurrentFunctionName(), nil, errors.New(cErr.Message))
			}
		}

		recollectionPendingPresent, cErr := eventProcessor.TestSampleMappingService.
			OrdersWithRecollectionsPendingPresentByOmsOrderId(omsOrderId)
		if cErr != nil {
			utils.AddLog(ctx, constants.ERROR_LEVEL, utils.GetCurrentFunctionName(), nil, errors.New(cErr.Message))
		}

		if !recollectionPendingPresent {
			go eventProcessor.SampleService.PublishRemoveSampleRejectedTagEvent(orderDetails.OmsRequestId,
				[]string{orderDetails.OmsOrderId}, !recollectionPendingPresent, orderDetails.CityCode)
		}
	}

	eventProcessor.SampleService.RemoveSamplesNotLinkedToAnyTests(omsOrderId)

	if omsOrderEvent.IsPartialCancellationFlow {
		if task.Id != 0 {
			cErr = eventProcessor.ReleaseReport(ctx, task.Id)
			if cErr != nil {
				utils.AddLog(ctx, constants.ERROR_LEVEL, utils.GetCurrentFunctionName(), nil, errors.New(cErr.Message))
			}
		}

		if len(toDeleteOmsTestIds) > 0 {
			eventProcessor.EtsService.GetAndPublishEtsTestBasicEvent(ctx, toDeleteOmsTestIds)
		}
	}

	return nil
}

func (eventProcessor *EventProcessor) recreateOrderDetailsDataForServicingLabChange(ctx context.Context,
	omsOrderEvent structures.OmsOrderCreateUpdateEvent, orderDetails models.OrderDetails,
	masterTestsMap map[uint]structures.CdsTestMaster,
	outsourceLabIds, nrlCpEnabledMasterTestIds []uint) *structures.CommonError {

	omsOrderId := omsOrderEvent.Order.AlnumOrderId

	patientDetails, cErr := eventProcessor.PatientDetailService.GetPatientDetailsById(orderDetails.PatientDetailsId)
	if cErr != nil {
		return cErr
	}

	err := eventProcessor.Db.Transaction(func(tx *gorm.DB) error {
		// Delete Test Details by Oms Order Id
		cErr := eventProcessor.TestDetailService.DeleteTestDetailsByOmsOrderIdWithTx(tx, omsOrderId)
		if cErr != nil {
			return errors.New(cErr.Message)
		}

		// Delete All Samples Data by Oms Order Id
		cErr = eventProcessor.SampleService.DeleteAllSamplesDataByOmsOrderIdWithTx(tx, omsOrderId)
		if cErr != nil {
			return errors.New(cErr.Message)
		}

		// Update Patient Details
		patientDetails, cErr = eventProcessor.PatientDetailService.UpdatePatientDetailsWithTx(tx,
			createUpdatePatientDetailsDtoForOmsOrderEvent(omsOrderEvent, patientDetails))
		if cErr != nil {
			return errors.New(cErr.Message)
		}

		// Update Order Details
		orderDetails, cErr = eventProcessor.OrderDetailsService.UpdateOrderDetailsWithTx(tx,
			createUpdateOrderDetailsDtoForOmsOrderEvent(omsOrderEvent, orderDetails))
		if cErr != nil {
			return errors.New(cErr.Message)
		}

		// Create Test Details
		testDetails, cErr := eventProcessor.TestDetailService.CreateTestDetailsWithTx(tx,
			createTestDetailsDtoForOmsOrderEvent(omsOrderEvent, masterTestsMap,
				outsourceLabIds, nrlCpEnabledMasterTestIds))
		if cErr != nil {
			return errors.New(cErr.Message)
		}

		// Create Test Details Metadata
		_, cErr = eventProcessor.TestDetailService.CreateTestDetailsMetadataWithTx(tx,
			createTestDetailsMetadataDtoForOmsOrderEvent(testDetails))
		if cErr != nil {
			return errors.New(cErr.Message)
		}

		// Create Samples
		if len(testDetails) > 0 {
			cErr = eventProcessor.SampleService.CreateSamplesWithTx(ctx, tx, omsOrderId)
			if cErr != nil {
				return errors.New(cErr.Message)
			}
		}

		return nil
	})

	if err != nil {
		return &structures.CommonError{
			Message:    err.Error(),
			StatusCode: http.StatusInternalServerError,
		}
	}

	return nil

}
