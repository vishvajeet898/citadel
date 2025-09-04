package service

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"strings"
	"time"

	"gorm.io/gorm"

	orderDetailMapper "github.com/Orange-Health/citadel/apps/order_details/mapper"
	"github.com/Orange-Health/citadel/apps/samples/constants"
	mappers "github.com/Orange-Health/citadel/apps/samples/mappers"
	"github.com/Orange-Health/citadel/apps/samples/structures"
	"github.com/Orange-Health/citadel/apps/samples/utils"
	testDetailMapper "github.com/Orange-Health/citadel/apps/test_detail/mapper"
	commonConstants "github.com/Orange-Health/citadel/common/constants"
	commonStructures "github.com/Orange-Health/citadel/common/structures"
	commonUtils "github.com/Orange-Health/citadel/common/utils"
	commonModels "github.com/Orange-Health/citadel/models"
)

func getLogisticsMinutesSpentBySampleCollectedAt(sampleCollectedAt *time.Time) string {
	if sampleCollectedAt == nil {
		return ""
	}
	duration := time.Since(*sampleCollectedAt)
	hours := int(duration.Hours())
	minutes := int(duration.Minutes()) % 60
	if hours > 0 {
		return fmt.Sprintf("%dh %dm", hours, minutes)
	}
	return fmt.Sprintf("%dm", minutes)
}

func validateAddBarcodeDetailsRequest(addBarcodeRequest structures.AddBarcodesRequest) *commonStructures.CommonError {
	for _, accession := range addBarcodeRequest.Accessions {
		if accession.Id == 0 {
			return &commonStructures.CommonError{
				Message:    commonConstants.ERROR_SAMPLE_ID_REQUIRED,
				StatusCode: http.StatusBadRequest,
			}
		} else if accession.Barcode == "" && accession.CollectLaterReason == "" {
			return &commonStructures.CommonError{
				Message:    constants.ERROR_BARCODE_AND_COLLECT_LATER_REASON_EMPTY,
				StatusCode: http.StatusBadRequest,
			}
		} else if accession.Barcode != "" {
			if accession.BarcodeImageURL == "" || accession.BarcodeScannedTime == nil {
				return &commonStructures.CommonError{
					Message:    constants.ERROR_BARCODE_URL_AND_SCANNED_TIME_REQUIRED,
					StatusCode: http.StatusBadRequest,
				}
			}
		}
	}
	return nil
}

func validateSampleForSampleRejection(samples []commonModels.Sample) *commonStructures.CommonError {
	allSamplesRejected := true
	for _, sample := range samples {
		if sample.Id == 0 {
			return &commonStructures.CommonError{
				StatusCode: http.StatusNotFound,
				Message:    commonConstants.ERROR_SAMPLE_NOT_FOUND,
			}
		}

		if sample.Status != commonConstants.SampleRejected {
			allSamplesRejected = false
		}
	}
	if allSamplesRejected {
		return &commonStructures.CommonError{
			StatusCode: http.StatusBadRequest,
			Message:    constants.ERROR_SAMPLE_ALREADY_REJECTED,
		}
	}
	return nil
}

func checkIfTicketShouldBeSentBasedOnReason(reason string) bool {
	reason = strings.ToLower(reason)
	return !commonUtils.SliceContainsString(constants.ReasonsToNotTriggerFreshdeskTicket, reason)
}

func (sampleService *SampleService) CreateFreshDeskTicketForSampleNotReceived(orderDetails commonModels.OrderDetails,
	sampleNumbers []uint) {
	subject := "Sample not received, action needed"

	testDetails, cErr := sampleService.SampleDao.GetAllSampleTestsBySampleNumbers(sampleNumbers, orderDetails.OmsOrderId)
	if cErr != nil {
		return
	}

	testNames := []string{}
	for _, testDetail := range testDetails {
		testNames = append(testNames, testDetail.TestName)
	}
	testNamesString := strings.Join(testNames, ", ")

	patientDetails, cErr := sampleService.PatientDetailsService.GetPatientDetailsById(orderDetails.PatientDetailsId)
	if cErr != nil {
		return
	}

	omsOrderIdUint := commonUtils.GetUintOrderIdWithoutStringPart(orderDetails.OmsOrderId)
	description := fmt.Sprint("Patient: " + patientDetails.Name + "\n" +
		"Tests: " + testNamesString + "\n" +
		"Order Id: " + orderDetails.OmsOrderId + "\n" +
		"Request Id: " + orderDetails.OmsRequestId + "\n")
	requestLink := fmt.Sprintf("%s/request/%s/order/%d", commonUtils.GetOmsBaseDomain(orderDetails.CityCode), orderDetails.OmsRequestId, omsOrderIdUint)
	description += "\n OMS Link: " + requestLink

	sampleService.AccountsApiClient.CreateFreshDeskTicket(context.Background(),
		commonConstants.FreshDeskCentralLogisticsGroupId, commonConstants.FreshDeskCentralLogisticsCreatorName, subject,
		description)
}

func (sampleService *SampleService) CreateFreshDeskTicketForSampleRejection(omsOrderId, subject, reason string,
	testDetails []commonModels.TestDetail) {
	var testNamesString string

	if len(testDetails) == 0 {
		return
	}
	var testNames []string
	for _, test := range testDetails {
		testNames = append(testNames, test.TestName)
	}
	testNamesString = strings.Join(testNames, ", ")

	orderDetails, cErr := sampleService.OrderDetailsService.GetOrderDetailsByOmsOrderId(omsOrderId)
	if cErr != nil {
		return
	}
	patientDetails, cErr := sampleService.PatientDetailsService.GetPatientDetailsById(orderDetails.PatientDetailsId)
	if cErr != nil {
		return
	}

	omsOrderIdUint := commonUtils.GetUintOrderIdWithoutStringPart(orderDetails.OmsOrderId)
	description := fmt.Sprint("Patient: " + patientDetails.Name + "\n" +
		"Tests: " + testNamesString + "\n" +
		"Order Id: " + omsOrderId + "\n" +
		"Request Id: " + orderDetails.OmsRequestId + "\n")
	requestLink := fmt.Sprintf("%s/request/%s/order/%d", commonUtils.GetOmsBaseDomain(orderDetails.CityCode), orderDetails.OmsRequestId, omsOrderIdUint)
	description += "\n OMS Link: " + requestLink
	description += "\n Reason Submitted:" + reason
	groupId, creatorName := utils.GetGroupIdAndNameForFreshdeskTicketBasedOnReason(reason)
	sampleService.AccountsApiClient.CreateFreshDeskTicket(context.Background(), groupId, creatorName, subject, description)
}

func (sampleService *SampleService) GetSampleDetailsByOmsOrderIdAndLabId(omsOrderId string) (
	structures.OrderTestsDetail, *commonStructures.CommonError) {
	orderDetailModel, testDetailModels, cErr := sampleService.GetOrderAndTestDetails(omsOrderId, "", false)
	if cErr != nil {
		return structures.OrderTestsDetail{}, cErr
	}

	orderDetail := orderDetailMapper.MapOrderDetailModelToOrderDetail(orderDetailModel)
	testDetails := testDetailMapper.MapTestDetails(testDetailModels)

	samples, cErr := sampleService.SampleDao.GetSamplesByOmsOrderId(omsOrderId)
	if cErr != nil {
		return structures.OrderTestsDetail{}, cErr
	}

	testSampleMappings, cErr := sampleService.TestSampleMappingService.GetTestSampleMappingsByOrderId(omsOrderId)
	if cErr != nil {
		return structures.OrderTestsDetail{}, cErr
	}

	return structures.OrderTestsDetail{
		OrderDetail:            orderDetail,
		TestDetails:            testDetails,
		SampleInfos:            samples,
		TestSampleMappingInfos: testSampleMappings,
	}, nil
}

func (sampleService *SampleService) GetSampleDetailsByOrderDetails(orderDetails []commonModels.OrderDetails) (
	[]structures.OrderTestsDetail, *commonStructures.CommonError) {

	omsOrderIds, omsOrderIdToOrderDetailsMap := []string{}, map[string]commonModels.OrderDetails{}
	for _, orderDetail := range orderDetails {
		omsOrderIds = append(omsOrderIds, orderDetail.OmsOrderId)
		omsOrderIdToOrderDetailsMap[orderDetail.OmsOrderId] = orderDetail
	}
	testDetails, cErr := sampleService.TestDetailsService.GetTestDetailsByOmsOrderIds(omsOrderIds)
	if cErr != nil {
		return nil, cErr
	}

	samples, cErr := sampleService.SampleDao.GetSamplesByOmsOrderIds(omsOrderIds)
	if cErr != nil {
		return nil, cErr
	}

	testSampleMappings, cErr := sampleService.TestSampleMappingService.GetTestSampleMappingsByOrderIds(omsOrderIds)
	if cErr != nil {
		return nil, cErr
	}

	omsOrderIdToTestDetailsMap, omsOrderIdToSamplesMap, omsOrderIdToTestSampleMappingsMap :=
		map[string][]commonModels.TestDetail{}, map[string][]commonStructures.SampleInfo{},
		map[string][]commonStructures.TestSampleMappingInfo{}

	for _, testDetail := range testDetails {
		omsOrderIdToTestDetailsMap[testDetail.OmsOrderId] = append(omsOrderIdToTestDetailsMap[testDetail.OmsOrderId],
			testDetail)
	}

	for _, sample := range samples {
		omsOrderIdToSamplesMap[sample.OmsOrderId] = append(omsOrderIdToSamplesMap[sample.OmsOrderId], sample)
	}

	for _, testSampleMapping := range testSampleMappings {
		omsOrderIdToTestSampleMappingsMap[testSampleMapping.OmsOrderId] =
			append(omsOrderIdToTestSampleMappingsMap[testSampleMapping.OmsOrderId], testSampleMapping)
	}

	orderTestsDetails := []structures.OrderTestsDetail{}
	for _, omsOrderId := range omsOrderIds {
		orderDetail := omsOrderIdToOrderDetailsMap[omsOrderId]
		testDetails := omsOrderIdToTestDetailsMap[omsOrderId]
		samples := omsOrderIdToSamplesMap[omsOrderId]
		testSampleMappings := omsOrderIdToTestSampleMappingsMap[omsOrderId]

		orderTestsDetails = append(orderTestsDetails, structures.OrderTestsDetail{
			OrderDetail:            orderDetailMapper.MapOrderDetailModelToOrderDetail(orderDetail),
			TestDetails:            testDetailMapper.MapTestDetails(testDetails),
			SampleInfos:            samples,
			TestSampleMappingInfos: testSampleMappings,
		})
	}

	return orderTestsDetails, nil
}

func (sampleService *SampleService) GetSampleDetailsByOmsRequestId(omsRequestId string) (
	[]structures.OrderTestsDetail, *commonStructures.CommonError) {
	orderDetails, cErr := sampleService.OrderDetailsService.GetOrderDetailsByOmsRequestId(omsRequestId)
	if cErr != nil {
		return nil, cErr
	}

	return sampleService.GetSampleDetailsByOrderDetails(orderDetails)
}

func (sampleService *SampleService) GetSampleDetailsByOmsOrderIds(omsOrderIds []string) (
	[]structures.OrderTestsDetail, *commonStructures.CommonError) {
	orderDetails, cErr := sampleService.OrderDetailsService.GetOrderDetailsByOmsOrderIds(omsOrderIds)
	if cErr != nil {
		return nil, cErr
	}

	return sampleService.GetSampleDetailsByOrderDetails(orderDetails)
}

func (sampleService *SampleService) GetCollectedSamples(omsOrderId string, labId uint) ([]commonModels.Sample,
	[]commonModels.SampleMetadata, *commonStructures.CommonError) {
	return sampleService.SampleDao.GetCollectedSamples(omsOrderId, labId)
}

func (sampleService *SampleService) GetSampleByBarcodeForReceiving(barcode string) (commonModels.Sample, *commonStructures.CommonError) {
	return sampleService.SampleDao.GetSampleByBarcodeForReceiving(barcode)
}

func (sampleService *SampleService) GetTestDetailsBySampleIds(sampleIds []uint) ([]commonModels.TestDetail,
	*commonStructures.CommonError) {
	return sampleService.SampleDao.GetTestDetailsBySampleIds(sampleIds)
}

func (sampleService *SampleService) GetAllTestsAndSampleMappingsBySampleNumbers(sampleNumbers []uint, omsOrderId string) (
	[]commonModels.TestDetail, []commonModels.TestSampleMapping, *commonStructures.CommonError) {

	return sampleService.SampleDao.GetAllTestsAndSampleMappingsBySampleNumbers(sampleNumbers, omsOrderId)
}

func (sampleService *SampleService) GetSampleDataBySampleId(sampleId uint) (
	commonModels.Sample, commonModels.SampleMetadata, *commonStructures.CommonError) {
	return sampleService.SampleDao.GetSampleDataBySampleId(sampleId)
}

func (sampleService *SampleService) GetSamplesDataBySampleIds(sampleIds []uint) (
	[]commonModels.Sample, []commonModels.SampleMetadata, *commonStructures.CommonError) {
	return sampleService.SampleDao.GetSamplesDataBySampleIds(sampleIds)
}

func (sampleService *SampleService) GetAllSampleTestsBySampleNumber(sampleNumber uint, omsOrderId string) (
	[]commonModels.TestDetail, *commonStructures.CommonError) {
	return sampleService.SampleDao.GetAllSampleTestsBySampleNumber(sampleNumber, omsOrderId)
}

func (sampleService *SampleService) CreateSamplesWithTx(ctx context.Context, tx *gorm.DB,
	omsOrderId string) *commonStructures.CommonError {
	orderDetails, testDetails, cErr := sampleService.GetOrderAndTestDetailsWithTx(tx, omsOrderId)
	if cErr != nil {
		return cErr
	}

	return sampleService.CreateSamplesWithOrderDetailsAndTestDetailsWithTx(ctx, tx, orderDetails, testDetails, 0)
}

func (sampleService *SampleService) CreateSamplesForRecollectionWithTx(ctx context.Context, tx *gorm.DB,
	orderDetails commonModels.OrderDetails, testDetails []commonModels.TestDetail,
	omsTaskId uint) *commonStructures.CommonError {

	if len(testDetails) == 0 {
		return nil
	}

	cErr := sampleService.CreateSamplesWithOrderDetailsAndTestDetailsWithTx(ctx, tx, orderDetails, testDetails, omsTaskId)
	if cErr != nil {
		return cErr
	}

	omsTestIds := []string{}
	for _, testDetail := range testDetails {
		omsTestIds = append(omsTestIds, testDetail.CentralOmsTestId)
	}

	return sampleService.TestSampleMappingService.UpdateTsmForTestIdAndRecollectionPendingTrueWithTx(tx, omsTestIds)
}

func (sampleService *SampleService) CreateSamplesAndTestDetailsWithTx(ctx context.Context, tx *gorm.DB, omsOrderId string,
	toBeCreatedTestDetails, labIdChangeTestDetails []commonModels.TestDetail) *commonStructures.CommonError {

	samples, cErr := sampleService.SampleDao.GetSamplesByOmsOrderId(omsOrderId)
	if cErr != nil && cErr.Message != commonConstants.ERROR_SAMPLE_NOT_FOUND {
		return cErr
	}

	if len(samples) == 0 {
		return sampleService.CreateSamplesWithTx(ctx, tx, omsOrderId)
	}

	orderDetails, testDetails, cErr := sampleService.GetOrderAndTestDetails(omsOrderId, "", false)
	if cErr != nil && cErr.Message != commonConstants.ERROR_NO_TEST_DETAILS_FOUND {
		return cErr
	}

	return sampleService.CreateUpdateTestsWithSampleWithTx(ctx, tx, orderDetails, testDetails, toBeCreatedTestDetails,
		labIdChangeTestDetails, samples)
}

func (sampleService *SampleService) CreateUpdateSamplesPostCollectionWithTx(ctx context.Context, tx *gorm.DB,
	omsOrderId string, createTestDetails, updateTestDetails []commonModels.TestDetail,
	omsTests []commonStructures.OmsTestModelDetails) (map[string]*time.Time, *commonStructures.CommonError) {

	orderDetails, testDetails, cErr := sampleService.GetOrderAndTestDetails(omsOrderId, "", false)
	if cErr != nil && cErr.Message != commonConstants.ERROR_NO_TEST_DETAILS_FOUND {
		return nil, cErr
	}

	samples, cErr := sampleService.SampleDao.GetSamplesByOmsOrderId(omsOrderId)
	if cErr != nil && cErr.Message != commonConstants.ERROR_SAMPLE_NOT_FOUND {
		return nil, cErr
	}

	return sampleService.CreateUpdateTestsWithSamplePostCollectionWithTx(ctx, tx, orderDetails, testDetails,
		createTestDetails, updateTestDetails, samples, omsTests)
}

func (sampleService *SampleService) SynchronizeTasksWithSamplesWithTx(tx *gorm.DB, omsRequestId string,
	tasks []commonStructures.OmsTaskModelDetails,
	taskTestsMapping [][]commonStructures.TestsJsonStruct) *commonStructures.CommonError {

	if len(tasks) == 0 {
		return nil
	}
	filteredTasks := []commonStructures.OmsTaskModelDetails{}
	for _, task := range tasks {
		if task.Id != 0 && task.DeletedOn == nil && task.TaskMarkedDescheduledTime == nil && task.TaskMarkedRnrTime == nil {
			filteredTasks = append(filteredTasks, task)
		}
	}
	sort.Slice(filteredTasks, func(i, j int) bool {
		return filteredTasks[i].Id < filteredTasks[j].Id
	})

	for idx, task := range filteredTasks {
		if idx >= len(taskTestsMapping) {
			return &commonStructures.CommonError{
				Message:    commonConstants.ERROR_SYNCHRONIZING_TASKS,
				StatusCode: http.StatusInternalServerError,
			}
		}
		selectedTests := taskTestsMapping[idx]
		var alnumTestIds []string

		for _, s := range selectedTests {
			if s.Tests != nil {
				for _, test := range s.Tests {
					alnumTestIds = append(alnumTestIds, test.AlnumTestId)
				}
			}
		}

		if len(alnumTestIds) != 0 {
			if cErr := sampleService.SampleDao.AssignTaskSequenceToSamples(tx, task.Id, alnumTestIds, task.IsAdditionalTask, omsRequestId); cErr != nil {
				return cErr
			}
		}
	}
	return nil
}

func (sampleService *SampleService) CreateUpdateTestsWithSampleWithTx(ctx context.Context, tx *gorm.DB,
	orderDetails commonModels.OrderDetails,
	testDetails, toBeCreatedTestDetails, labIdChangeTestDetails []commonModels.TestDetail,
	samples []commonStructures.SampleInfo) *commonStructures.CommonError {

	allTestDetails, omsTestIdToTestDetailsMap := []commonModels.TestDetail{}, map[string]commonModels.TestDetail{}
	labIdChangeOmsTestIds := []string{}
	for _, testDetail := range testDetails {
		omsTestIdToTestDetailsMap[testDetail.CentralOmsTestId] = testDetail
	}
	for _, testDetail := range toBeCreatedTestDetails {
		omsTestIdToTestDetailsMap[testDetail.CentralOmsTestId] = testDetail
	}
	for _, testDetail := range labIdChangeTestDetails {
		omsTestIdToTestDetailsMap[testDetail.CentralOmsTestId] = testDetail
		labIdChangeOmsTestIds = append(labIdChangeOmsTestIds, testDetail.CentralOmsTestId)
	}
	for _, testDetail := range omsTestIdToTestDetailsMap {
		allTestDetails = append(allTestDetails, testDetail)
	}
	recomputedSampleDetails, cErr := sampleService.GetSamples(ctx, orderDetails, allTestDetails)
	if cErr != nil {
		return cErr
	}
	commonUtils.AddLog(ctx, commonConstants.DEBUG_LEVEL, commonUtils.GetCurrentFunctionName(), map[string]interface{}{
		"recomputedSampleDetails": recomputedSampleDetails,
	}, nil)
	commonUtils.AddLog(ctx, commonConstants.DEBUG_LEVEL, commonUtils.GetCurrentFunctionName(), map[string]interface{}{
		"orderDetails":    orderDetails,
		"testDetails":     testDetails,
		"current_samples": samples,
	}, nil)

	testSampleMappings, cErr :=
		sampleService.TestSampleMappingService.GetTestSampleMappingsModelsByOrderId(orderDetails.OmsOrderId)
	if cErr != nil {
		return cErr
	}

	omsTestIdToTestSampleMappingsMap := map[string][]commonModels.TestSampleMapping{}
	for _, testSampleMapping := range testSampleMappings {
		omsTestIdToTestSampleMappingsMap[testSampleMapping.OmsTestId] =
			append(omsTestIdToTestSampleMappingsMap[testSampleMapping.OmsTestId], testSampleMapping)
	}
	sampleNumberToTestDetailsMap := utils.CreateSampleNumberMaps(testDetails, testSampleMappings)
	sampleIdToProcessingLabIdMap := map[uint]uint{}
	for _, sample := range samples {
		sampleId := sample.Id
		sampleNumber := sample.SampleNumber
		testDetailsForSample := sampleNumberToTestDetailsMap[sampleNumber]
		if len(testDetailsForSample) == 0 {
			continue
		}
		processingLabId := testDetailsForSample[0].ProcessingLabId
		sampleIdToProcessingLabIdMap[sampleId] = processingLabId
	}

	existingTestIds := make(map[string]bool)
	for _, testDetail := range testDetails {
		if commonUtils.SliceContainsString(labIdChangeOmsTestIds, testDetail.CentralOmsTestId) {
			continue
		}
		existingTestIds[testDetail.CentralOmsTestId] = true
	}

	oldSamplesMap, oldSampleNumberMap := map[string]commonStructures.SampleInfo{}, map[string]uint{}
	for _, sample := range samples {
		finalProcessingLabId := sampleIdToProcessingLabIdMap[sample.Id]
		key := fmt.Sprint(sample.CollectionSequenceNumber, ":", sample.VialTypeId, ":", finalProcessingLabId)
		oldSamplesMap[key] = sample
		oldSampleNumberMap[key] = sample.SampleNumber
	}

	latestSampleNumber := sampleService.SampleDao.GetMaxSampleNumberByOmsOrderId(orderDetails.OmsOrderId)

	for _, sampleTestMappingDetails := range recomputedSampleDetails {
		sampleFound, oldSample := false, commonStructures.SampleInfo{}
		sampleInfo := sampleTestMappingDetails.SampleInfo
		omsTestIdOfTestInfoDetail := sampleTestMappingDetails.TestSampleMappingInfos[0].OmsTestId
		processingLabId := omsTestIdToTestDetailsMap[omsTestIdOfTestInfoDetail].ProcessingLabId
		key := fmt.Sprint(sampleInfo.CollectionSequenceNumber, ":", sampleInfo.VialTypeId, ":", processingLabId)
		if _, ok := oldSamplesMap[key]; ok {
			sampleFound = true
			oldSample = oldSamplesMap[key]
		}

		if sampleFound {
			for _, testSampleMappingInfo := range sampleTestMappingDetails.TestSampleMappingInfos {
				if existingTestIds[testSampleMappingInfo.OmsTestId] {
					continue
				}
				if commonUtils.SliceContainsString(labIdChangeOmsTestIds, testSampleMappingInfo.OmsTestId) {
					sampleService.TestSampleMappingService.DeleteTestSampleMappingByOmsTestIdsWithTx(tx,
						[]string{testSampleMappingInfo.OmsTestId})
				}
				testSampleMappingInfo.SampleNumber = oldSample.SampleNumber
				testSampleMappingInfo.SampleId = oldSample.SampleId
				_, err := sampleService.TestSampleMappingService.CreateTestSampleMappingWithTx(tx, testSampleMappingInfo)
				if err != nil {
					return err
				}
			}
		} else {
			latestSampleNumber++
			sampleInfo.SampleNumber = latestSampleNumber
			newSample, err := sampleService.SampleDao.CreateSampleWithTx(tx, sampleInfo)
			if err != nil {
				return err
			}

			for _, testSampleMappingInfo := range sampleTestMappingDetails.TestSampleMappingInfos {
				previousSampleMappings := omsTestIdToTestSampleMappingsMap[testSampleMappingInfo.OmsTestId]
				sampleService.TestSampleMappingService.DeleteTestSampleMappingsWithTx(tx, previousSampleMappings)
				testSampleMappingInfo.SampleNumber = newSample.SampleNumber
				testSampleMappingInfo.SampleId = newSample.Id
				_, err := sampleService.TestSampleMappingService.CreateTestSampleMappingWithTx(tx, testSampleMappingInfo)
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func (sampleService *SampleService) CreateUpdateTestsWithSamplePostCollectionWithTx(ctx context.Context, tx *gorm.DB,
	orderDetails commonModels.OrderDetails, testDetails, createTestDetails, updateTestDetails []commonModels.TestDetail,
	samples []commonStructures.SampleInfo, omsTests []commonStructures.OmsTestModelDetails) (
	map[string]*time.Time, *commonStructures.CommonError) {

	sampleNumberToSampleMapping := make(map[uint]commonStructures.SampleInfo)
	for _, sample := range samples {
		sampleNumberToSampleMapping[sample.SampleNumber] = sample
	}
	createOmsTestIds, updateOmsTestIds := []string{}, []string{}
	testIdToLisSyncAtTime := map[string](*time.Time){}

	updateTestSampleMappings := []commonModels.TestSampleMapping{}
	testSampleMappings, cErr := sampleService.TestSampleMappingService.GetTestSampleMappingsModelsByOrderId(
		orderDetails.OmsOrderId)
	if cErr != nil {
		return nil, cErr
	}
	sampleNumberToTestSampleMapping, omsTestIdToTestSampleMapping :=
		map[uint][]commonModels.TestSampleMapping{}, map[string]commonModels.TestSampleMapping{}
	for _, testSampleMapping := range testSampleMappings {
		sampleNumberToTestSampleMapping[testSampleMapping.SampleNumber] = append(
			sampleNumberToTestSampleMapping[testSampleMapping.SampleNumber], testSampleMapping)
		omsTestIdToTestSampleMapping[testSampleMapping.OmsTestId] = testSampleMapping
	}

	allTestDetailsMap := make(map[string]commonModels.TestDetail)
	for _, testDetail := range testDetails {
		allTestDetailsMap[testDetail.CentralOmsTestId] = testDetail
	}
	for _, testDetail := range createTestDetails {
		allTestDetailsMap[testDetail.CentralOmsTestId] = testDetail
		createOmsTestIds = append(createOmsTestIds, testDetail.CentralOmsTestId)
	}
	for _, testDetail := range updateTestDetails {
		allTestDetailsMap[testDetail.CentralOmsTestId] = testDetail
		updateOmsTestIds = append(updateOmsTestIds, testDetail.CentralOmsTestId)
	}

	sampleStatusesForModifyingAttuneVisit := []string{commonConstants.SampleAccessioned, commonConstants.SampleSynced,
		commonConstants.SamplePartiallyRejected}

	omsTestsMap, newSampleNumberOmsTestIds := make(map[uint]commonStructures.OmsTestModelDetails), []string{}
	newSampleNumerDefaultTests, newSampleNumberCollectedTests := []commonModels.TestDetail{}, []commonModels.TestDetail{}
	for _, omsTest := range omsTests {
		omsTestsMap[omsTest.Id] = omsTest
		if omsTest.TestAddUpdateType == commonConstants.AddNewSampleDefaultType {
			newSampleNumerDefaultTests = append(newSampleNumerDefaultTests, allTestDetailsMap[omsTest.AlnumTestId])
			newSampleNumberOmsTestIds = append(newSampleNumberOmsTestIds, omsTest.AlnumTestId)
		} else if omsTest.TestAddUpdateType == commonConstants.AddNewSampleCollectedType {
			newSampleNumberCollectedTests = append(newSampleNumberCollectedTests, allTestDetailsMap[omsTest.AlnumTestId])
			newSampleNumberOmsTestIds = append(newSampleNumberOmsTestIds, omsTest.AlnumTestId)
		} else if omsTest.TestAddUpdateType == commonConstants.MapExistingSampleType {
			if commonUtils.SliceContainsString(updateOmsTestIds, omsTest.AlnumTestId) {
				currentTestSampleMapping := omsTestIdToTestSampleMapping[omsTest.AlnumTestId]
				if currentTestSampleMapping.Id == 0 {
					return nil, &commonStructures.CommonError{
						Message:    fmt.Sprintf(commonConstants.ERROR_NO_SAMPLE_TEST_MAPPING_DETAILS),
						StatusCode: http.StatusBadRequest,
					}
				}
				newSample := sampleNumberToSampleMapping[omsTest.MappingSampleNumber]
				newTestSampleMapping := currentTestSampleMapping
				newTestSampleMapping.SampleId = newSample.Id
				newTestSampleMapping.SampleNumber = newSample.SampleNumber
				updateTestSampleMappings = append(updateTestSampleMappings, newTestSampleMapping)

				// TODO @shrish: Adds check if the visit is completed on Attune
				if commonUtils.SliceContainsString(sampleStatusesForModifyingAttuneVisit,
					sampleNumberToSampleMapping[currentTestSampleMapping.SampleNumber].Status) {
					testIdToLisSyncAtTime[newTestSampleMapping.OmsTestId] = newSample.LisSyncAt
					cErr = sampleService.AttuneService.ModifyLisDataPostSyncByOrderId(ctx, orderDetails.OmsOrderId,
						orderDetails.ServicingLabId, allTestDetailsMap[omsTest.AlnumTestId], newSample)
					if cErr != nil {
						return nil, cErr
					}
				}
			} else if commonUtils.SliceContainsString(createOmsTestIds, omsTest.AlnumTestId) {
				sample := sampleNumberToSampleMapping[omsTest.MappingSampleNumber]
				newTestSampleMapping := commonModels.TestSampleMapping{
					OmsCityCode:  orderDetails.CityCode,
					OmsTestId:    omsTest.AlnumTestId,
					OmsOrderId:   orderDetails.OmsOrderId,
					SampleId:     sample.Id,
					SampleNumber: sample.SampleNumber,
					VialTypeId:   sample.VialTypeId,
				}
				updateTestSampleMappings = append(updateTestSampleMappings, newTestSampleMapping)
				if commonUtils.SliceContainsString(sampleStatusesForModifyingAttuneVisit,
					sampleNumberToSampleMapping[testSampleMappings[0].SampleNumber].Status) {
					testIdToLisSyncAtTime[newTestSampleMapping.OmsTestId] = sample.LisSyncAt
					cErr = sampleService.AttuneService.ModifyLisDataPostSyncByOrderId(ctx, orderDetails.OmsOrderId,
						orderDetails.ServicingLabId, allTestDetailsMap[omsTest.AlnumTestId], sample)
					if cErr != nil {
						return nil, cErr
					}
				}
			}
		} else if omsTest.TestAddUpdateType == commonConstants.LabIdModificationType {
			currentTestSampleMapping := omsTestIdToTestSampleMapping[omsTest.AlnumTestId]
			if currentTestSampleMapping.Id == 0 {
				return nil, &commonStructures.CommonError{
					Message:    fmt.Sprintf(commonConstants.ERROR_NO_SAMPLE_TEST_MAPPING_DETAILS),
					StatusCode: http.StatusBadRequest,
				}
			}
			sample := sampleNumberToSampleMapping[currentTestSampleMapping.SampleNumber]
			sample.VisitId = ""
			if sample.Status == commonConstants.SampleSynced {
				sample.Status = commonConstants.SampleReceived
				sample.LisSyncAt = nil
				sample.ReceivedAt = commonUtils.GetCurrentTime()
				_, _, cErr = sampleService.SampleDao.UpdateSampleAndSampleMetadataWithTx(tx,
					mappers.MapSampleInfoToSample(sample), mappers.MapSampleInfoToSampleMetadata(sample))
				if cErr != nil {
					return nil, cErr
				}
			} else if sample.Status == commonConstants.SampleReceived {
				if orderDetails.ServicingLabId == omsTest.LabId {
					sample.Status = commonConstants.SampleSynced
					sampleModel := mappers.MapSampleInfoToSample(sample)
					sampleMetadataModel := mappers.MapSampleInfoToSampleMetadata(sample)
					sampleIdBarcodeMap := map[uint]string{}
					sampleIdBarcodeMap[sample.Id] = sample.Barcode
					visitId, cErr := sampleService.AttuneService.SyncDataToLisByOmsOrderId(ctx, orderDetails.OmsOrderId,
						orderDetails.ServicingLabId, []commonModels.Sample{sampleModel},
						[]commonModels.SampleMetadata{sampleMetadataModel}, sampleIdBarcodeMap)
					if cErr != nil {
						continue
					}
					sample.VisitId = visitId
					sampleModel.VisitId = visitId
					sampleMetadataModel.LisSyncAt = commonUtils.GetCurrentTime()
					_, _, cErr = sampleService.SampleDao.UpdateSampleAndSampleMetadataWithTx(tx, sampleModel,
						sampleMetadataModel)
					if cErr != nil {
						continue
					}
				}
			}
		}
	}

	if len(updateTestSampleMappings) > 0 {
		_, cErr = sampleService.TestSampleMappingService.UpdateBulkTestSampleMappingWithTx(tx, updateTestSampleMappings)
		if cErr != nil {
			return nil, cErr
		}
	}

	// Delete Test Sample Mappings for these tests
	if len(newSampleNumberOmsTestIds) > 0 {
		cErr = sampleService.TestSampleMappingService.DeleteTestSampleMappingByOmsTestIdsWithTx(tx, newSampleNumberOmsTestIds)
		if cErr != nil {
			return nil, cErr
		}
	}

	latestSampleNumber := sampleService.SampleDao.GetMaxSampleNumberByOmsOrderId(orderDetails.OmsOrderId)

	if len(newSampleNumerDefaultTests) > 0 {
		newDefaultSampleDetails, cErr := sampleService.GetSamples(ctx, orderDetails, newSampleNumerDefaultTests)
		if cErr != nil {
			return nil, cErr
		}
		for _, sampleTestMappingDetails := range newDefaultSampleDetails {
			sampleInfo := sampleTestMappingDetails.SampleInfo

			latestSampleNumber++
			sampleInfo.SampleNumber = latestSampleNumber
			newSample, cErr := sampleService.SampleDao.CreateSampleWithTx(tx, sampleInfo)
			if cErr != nil {
				return nil, cErr
			}
			for _, testSampleMappingInfo := range sampleTestMappingDetails.TestSampleMappingInfos {
				testSampleMappingInfo.SampleNumber = newSample.SampleNumber
				testSampleMappingInfo.SampleId = newSample.Id
				testSampleMappingInfo.RecollectionPending = true
				_, cErr := sampleService.TestSampleMappingService.CreateTestSampleMappingWithTx(tx, testSampleMappingInfo)
				if cErr != nil {
					return nil, cErr
				}
			}
		}
	}

	if len(newSampleNumberCollectedTests) > 0 {
		newCollectedSampleDetails, cErr := sampleService.GetSamples(ctx, orderDetails, newSampleNumberCollectedTests)
		if cErr != nil {
			return nil, cErr
		}
		for _, sampleTestMappingDetails := range newCollectedSampleDetails {
			sampleInfo := sampleTestMappingDetails.SampleInfo

			latestSampleNumber++
			sampleInfo.Status = commonConstants.SampleCollectionDone
			sampleInfo.CollectedAt = commonUtils.GetCurrentTime()
			sampleInfo.SampleNumber = latestSampleNumber
			newSample, cErr := sampleService.SampleDao.CreateSampleWithTx(tx, sampleInfo)
			if cErr != nil {
				return nil, cErr
			}
			for _, testSampleMappingInfo := range sampleTestMappingDetails.TestSampleMappingInfos {
				testSampleMappingInfo.SampleNumber = newSample.SampleNumber
				testSampleMappingInfo.SampleId = newSample.Id
				_, cErr := sampleService.TestSampleMappingService.CreateTestSampleMappingWithTx(tx, testSampleMappingInfo)
				if cErr != nil {
					return nil, cErr
				}
			}
		}
	}

	return testIdToLisSyncAtTime, nil
}

func (sampleService *SampleService) CreateSampleWithOrderDetailsAndTestDetails(ctx context.Context,
	orderDetails commonModels.OrderDetails, testDetails []commonModels.TestDetail,
	omsTaskId uint) *commonStructures.CommonError {

	sampleTestMappingDetails, cErr := sampleService.GetSamples(ctx, orderDetails, testDetails)
	if cErr != nil {
		return cErr
	}

	tx := sampleService.SampleDao.BeginTransaction()
	defer tx.Rollback()

	for _, sampleTestMappingDetail := range sampleTestMappingDetails {
		sampleTestMappingDetail.SampleInfo.TaskSequence = omsTaskId
		sample, err := sampleService.SampleDao.CreateSampleWithTx(tx, sampleTestMappingDetail.SampleInfo)
		if err != nil {
			return err
		}
		for idx := range sampleTestMappingDetail.TestSampleMappingInfos {
			sampleTestMappingDetail.TestSampleMappingInfos[idx].SampleId = sample.Id
		}
		_, err = sampleService.TestSampleMappingService.CreateBulkTestSampleMappingWithTx(
			tx, sampleTestMappingDetail.TestSampleMappingInfos)
		if err != nil {
			return err
		}
	}

	tx.Commit()
	return nil
}

func (sampleService *SampleService) RemoveSamplesNotLinkedToAnyTests(omsOrderId string) *commonStructures.CommonError {
	return sampleService.SampleDao.RemoveSamplesNotLinkedToAnyTests(omsOrderId)
}

func (sampleService *SampleService) CreateSamplesWithOrderDetailsAndTestDetailsWithTx(ctx context.Context, tx *gorm.DB,
	orderDetails commonModels.OrderDetails, testDetails []commonModels.TestDetail,
	omsTaskId uint) *commonStructures.CommonError {
	sampleTestMappingDetails, cErr := sampleService.GetSamplesWithTx(ctx, tx, orderDetails, testDetails)
	if cErr != nil {
		return cErr
	}

	for _, sampleTestMappingDetail := range sampleTestMappingDetails {
		sampleTestMappingDetail.SampleInfo.TaskSequence = omsTaskId
		if orderDetails.TrfId != "" {
			sampleTestMappingDetail.SampleInfo.Status = commonConstants.SampleCollectionDone
			sampleTestMappingDetail.SampleInfo.CollectedAt = commonUtils.GetCurrentTime()
		}
		if orderDetails.BulkOrderId != 0 {
			sampleTestMappingDetail.SampleInfo.Status = commonConstants.SampleCollectionDone
			sampleTestMappingDetail.SampleInfo.CollectedAt = orderDetails.CollectedOn
		}
		sample, err := sampleService.SampleDao.CreateSampleWithTx(tx, sampleTestMappingDetail.SampleInfo)
		if err != nil {
			return err
		}
		for idx := range sampleTestMappingDetail.TestSampleMappingInfos {
			sampleTestMappingDetail.TestSampleMappingInfos[idx].SampleId = sample.Id
		}
		_, err = sampleService.TestSampleMappingService.CreateBulkTestSampleMappingWithTx(tx, sampleTestMappingDetail.TestSampleMappingInfos)
		if err != nil {
			return err
		}
	}
	return nil
}

func (sampleService *SampleService) UpdateTaskIdByOmsTestIdsWithTx(tx *gorm.DB, omsTestIds []string,
	taskId uint) *commonStructures.CommonError {
	if len(omsTestIds) == 0 {
		return nil
	}
	return sampleService.SampleDao.UpdateTaskIdByOmsTestIdsWithTx(tx, omsTestIds, taskId)
}

func (sampleService *SampleService) UpdateSampleCollected(
	sampleCollectionRequest commonStructures.SampleCollectedRequest) *commonStructures.CommonError {
	statusArray := []string{commonConstants.SampleDefault}
	var sampleInfos = []commonStructures.SampleInfo{}
	var cErr *commonStructures.CommonError
	if sampleCollectionRequest.CollectionType != "" && sampleCollectionRequest.CollectionType == "B2C" {
		statusArray = append(statusArray, commonConstants.SampleCollectionDone)
	}

	switch sampleCollectionRequest.OmsTaskType {
	case commonConstants.OmsTaskTypePrimaryCollection:
		sampleInfos, cErr = sampleService.SampleDao.GetSampleForCollectedB2C(sampleCollectionRequest.RequestId, 0,
			statusArray, sampleCollectionRequest.IsB2c, false)
		if cErr != nil {
			return cErr
		}
	case commonConstants.OmsTaskTypeRecollection:
		sampleInfos, cErr = sampleService.SampleDao.GetSampleForCollectedB2C(sampleCollectionRequest.RequestId,
			sampleCollectionRequest.OmsTaskId, []string{commonConstants.SampleDefault}, sampleCollectionRequest.IsB2c, true)
		if cErr != nil {
			return cErr
		}
	}

	if len(sampleInfos) == 0 {
		return &commonStructures.CommonError{
			Message:    commonConstants.ERROR_NO_SAMPLES_FOUND_TO_MARK_COLLECTED,
			StatusCode: http.StatusNotFound,
		}
	}

	for idx := range sampleInfos {
		sampleInfos[idx].Status = commonConstants.SampleCollectionDone
		sampleInfos[idx].CollectedAt = commonUtils.GetTimeFromString(sampleCollectionRequest.CollectedAt)
	}

	if _, cErr := sampleService.SampleDao.UpdateBulkSamples(sampleInfos); cErr != nil {
		return cErr
	}
	return nil
}

func (sampleService *SampleService) ForcefullyMarkSampleAsCollected(
	request structures.ForcefullyMarkSampleAsCollectedRequest) *commonStructures.CommonError {
	if len(request.SampleNumbers) == 0 || request.OmsOrderId == "" {
		return &commonStructures.CommonError{
			Message:    constants.ERROR_INVALID_PARAMETERS,
			StatusCode: http.StatusBadRequest,
		}
	}
	omsTestIds := []string{}
	samples, sampleMetadatas, cErr := sampleService.SampleDao.GetSamplesByOmsOrderIdAndSampleNumbers(request.OmsOrderId,
		request.SampleNumbers)
	if cErr != nil {
		return cErr
	}

	for index := range samples {
		samples[index].Status = commonConstants.SampleCollectionDone
		samples[index].UpdatedBy = request.UserId
	}

	collectedAt := commonUtils.GetTimeFromString(request.CollectedAt)
	if collectedAt == nil {
		collectedAt = commonUtils.GetCurrentTime()
	}

	for index := range sampleMetadatas {
		sampleMetadatas[index].CollectedAt = collectedAt
		sampleMetadatas[index].LastUpdatedAt = commonUtils.GetCurrentTime()
		sampleMetadatas[index].CollectLaterReason = ""
		sampleMetadatas[index].UpdatedBy = request.UserId
	}

	testDetails, cErr := sampleService.SampleDao.GetAllSampleTestsBySampleNumbers(request.SampleNumbers, request.OmsOrderId)
	if cErr != nil {
		return cErr
	}

	for index := range testDetails {
		testDetails[index].Status = commonConstants.TEST_STATUS_REQUESTED
		omsTestIds = append(omsTestIds, testDetails[index].CentralOmsTestId)
	}

	tx := sampleService.SampleDao.BeginTransaction()
	defer tx.Rollback()

	_, _, cErr = sampleService.SampleDao.UpdateSamplesAndSamplesMetadataWithTx(tx, samples, sampleMetadatas)
	if cErr != nil {
		return cErr
	}

	_, cErr = sampleService.TestDetailsService.UpdateTestDetailsWithTx(tx, testDetails)
	if cErr != nil {
		return cErr
	}

	tx.Commit()

	sampleService.PublishSampleCollectedEvent(omsTestIds, collectedAt, samples[0].OmsCityCode)
	return nil
}

func (sampleService *SampleService) AddBarcodeDetails(addBarcodeRequest structures.AddBarcodesRequest) (
	map[string]string, *commonStructures.CommonError) {
	cErr := validateAddBarcodeDetailsRequest(addBarcodeRequest)
	if cErr != nil {
		return nil, cErr
	}

	barcodes, sampleIds := []string{}, []uint{}
	for _, accession := range addBarcodeRequest.Accessions {
		if accession.Barcode != "" {
			barcodes = append(barcodes, accession.Barcode)
		}
		sampleIds = append(sampleIds, accession.Id)
	}

	if len(barcodes) > 0 {
		barcodesExists, cErr := sampleService.SampleDao.BarcodesExistsInSystem(barcodes)
		if cErr != nil {
			return nil, cErr
		}
		if barcodesExists {
			return nil, &commonStructures.CommonError{
				Message:    commonConstants.ERROR_DUPLICATE_BARCODE,
				StatusCode: http.StatusBadRequest,
			}
		}
	}

	sampleIdToTestDetailsMap, testIdStatusMap := map[uint][]commonModels.TestDetail{}, map[string]string{}
	updatedSamples, updatedSamplesMetadata := []commonModels.Sample{}, []commonModels.SampleMetadata{}
	sampleIdToSampleMap, SampleIdToSampleMetadataMap :=
		map[uint]commonModels.Sample{}, map[uint]commonModels.SampleMetadata{}

	samples, samplesMetadata, cErr := sampleService.GetSamplesDataBySampleIds(sampleIds)
	if cErr != nil {
		return nil, cErr
	}

	for _, sample := range samples {
		sampleIdToSampleMap[sample.Id] = sample
	}
	for _, sampleMetadata := range samplesMetadata {
		SampleIdToSampleMetadataMap[sampleMetadata.SampleId] = sampleMetadata
	}

	// Optimize this call
	for _, sample := range sampleIdToSampleMap {
		testDetails, cErr := sampleService.SampleDao.GetAllSampleTestsBySampleNumber(sample.SampleNumber, sample.OmsOrderId)
		if cErr != nil {
			return nil, cErr
		}
		sampleIdToTestDetailsMap[sample.Id] = testDetails
	}

	for _, accession := range addBarcodeRequest.Accessions {
		sample := sampleIdToSampleMap[accession.Id]
		sampleMetadata := SampleIdToSampleMetadataMap[accession.Id]
		sample.Barcode = accession.Barcode
		sampleMetadata.BarcodeImageUrl = accession.BarcodeImageURL
		sampleMetadata.BarcodeScannedAt = accession.BarcodeScannedTime
		sampleMetadata.CollectLaterReason = accession.CollectLaterReason
		sampleMetadata.LastUpdatedAt = commonUtils.GetCurrentTime()
		sample.UpdatedBy = commonConstants.CitadelSystemId
		updatedSamples = append(updatedSamples, sample)
		updatedSamplesMetadata = append(updatedSamplesMetadata, sampleMetadata)
		for _, testDetail := range sampleIdToTestDetailsMap[sample.Id] {
			if accession.CollectLaterReason == "" {
				testIdStatusMap[testDetail.CentralOmsTestId] = commonConstants.TEST_STATUS_REQUESTED
			}
			if accession.CollectLaterReason != "" && accession.CollectLaterReason != commonConstants.NotCollectedReasonCollectLater {
				testIdStatusMap[testDetail.CentralOmsTestId] = commonConstants.TEST_STATUS_COLLECT_SAMPLE_LATER
			}
			if accession.CollectLaterReason != "" && accession.CollectLaterReason == commonConstants.NotCollectedReasonCollectLater {
				testIdStatusMap[testDetail.CentralOmsTestId] = commonConstants.TEST_STATUS_REQUESTED
			}
		}
	}

	tx := sampleService.SampleDao.BeginTransaction()
	defer tx.Rollback()

	_, _, cErr = sampleService.SampleDao.UpdateSamplesAndSamplesMetadataWithTx(tx, updatedSamples, updatedSamplesMetadata)
	if cErr != nil {
		return nil, cErr
	}
	if cErr := sampleService.TestDetailsService.UpdateTestStatusesByOmsTestIdsWithTx(tx, testIdStatusMap,
		commonConstants.CitadelSystemId); cErr != nil {
		return nil, cErr
	}

	tx.Commit()

	return testIdStatusMap, nil
}

func (sampleService *SampleService) AddBarcodeDetailsForOrangers(accessionBody structures.UpdateAccessionBody) (
	map[string]string, *commonStructures.CommonError) {
	if accessionBody.AccessionId == 0 {
		return nil, &commonStructures.CommonError{
			Message:    commonConstants.ERROR_SAMPLE_ID_REQUIRED,
			StatusCode: http.StatusBadRequest,
		}
	}
	testIdStatusMap := map[string]string{}

	sample, sampleMetadata, cErr := sampleService.GetSampleDataBySampleId(accessionBody.AccessionId)
	if cErr != nil {
		return nil, cErr
	}

	testDetails, cErr := sampleService.SampleDao.GetAllSampleTestsBySampleNumber(sample.SampleNumber, sample.OmsOrderId)
	if cErr != nil {
		return nil, cErr
	}

	currentTime := commonUtils.GetCurrentTime()

	switch {
	case accessionBody.SkipReason == "" && accessionBody.Barcode != "" && accessionBody.BarcodeImageURL != "":
		barcodesExists, cErr := sampleService.SampleDao.BarcodesExistsInSystem([]string{accessionBody.Barcode})
		if cErr != nil {
			return nil, cErr
		}
		if barcodesExists {
			return nil, &commonStructures.CommonError{
				Message:    commonConstants.ERROR_DUPLICATE_BARCODE,
				StatusCode: http.StatusBadRequest,
			}
		}
		sample.Barcode = accessionBody.Barcode
		sampleMetadata.BarcodeImageUrl = accessionBody.BarcodeImageURL
		sampleMetadata.BarcodeScannedAt = currentTime
		sampleMetadata.CollectLaterReason = ""
	case accessionBody.Barcode == "" && accessionBody.SkipReason != "":
		sampleMetadata.CollectLaterReason = accessionBody.SkipReason
		sample.Barcode = ""
		sampleMetadata.BarcodeImageUrl = ""
		sampleMetadata.BarcodeScannedAt = nil
	}

	sample.UpdatedBy = commonConstants.CitadelSystemId
	sampleMetadata.LastUpdatedAt = currentTime
	sample.UpdatedAt = currentTime
	sampleMetadata.UpdatedBy = commonConstants.CitadelSystemId
	sampleMetadata.UpdatedAt = currentTime

	for _, testDetail := range testDetails {
		if accessionBody.Barcode != "" && accessionBody.SkipReason == "" {
			testIdStatusMap[testDetail.CentralOmsTestId] = commonConstants.TEST_STATUS_REQUESTED
		} else if accessionBody.SkipReason != "" {
			if accessionBody.SkipReason == commonConstants.NotCollectedReasonCollectLater {
				testIdStatusMap[testDetail.CentralOmsTestId] = commonConstants.TEST_STATUS_REQUESTED
			} else {
				testIdStatusMap[testDetail.CentralOmsTestId] = commonConstants.TEST_STATUS_COLLECT_SAMPLE_LATER
			}
		}
	}

	tx := sampleService.SampleDao.BeginTransaction()
	defer tx.Rollback()

	_, _, cErr = sampleService.SampleDao.UpdateSampleAndSampleMetadataWithTx(tx, sample, sampleMetadata)
	if cErr != nil {
		return nil, cErr
	}
	if cErr := sampleService.TestDetailsService.UpdateTestStatusesByOmsTestIdsWithTx(tx, testIdStatusMap,
		commonConstants.CitadelSystemId); cErr != nil {
		return nil, cErr
	}

	tx.Commit()

	return testIdStatusMap, nil
}

func (sampleService *SampleService) GetSampleDetailsForScheduler(sampleDetailsRequest structures.SampleDetailsRequest) (
	map[string][]structures.SampleDetailsResponse, *commonStructures.CommonError) {

	response := map[string][]structures.SampleDetailsResponse{}
	vialsResponseFromDb, cErr := sampleService.SampleDao.GetSampleDetailsForScheduler(sampleDetailsRequest)
	if cErr != nil {
		return nil, cErr
	}

	for _, vial := range vialsResponseFromDb {
		if _, ok := response[vial.OrderId]; !ok {
			response[vial.OrderId] = []structures.SampleDetailsResponse{}
		}
		vialObject := structures.SampleDetailsResponse{
			AccessionId:      vial.AccessionId,
			TaskSequence:     vial.TaskSequence,
			VialType:         vial.VialType,
			Barcode:          vial.Barcode,
			ImageUrl:         vial.ImageUrl,
			ReasonForSkip:    vial.ReasonForSkip,
			BarcodeScannedAt: vial.BarcodeScannedAt,
			CollectedVolume:  vial.CollectedVolume,
		}

		err := json.Unmarshal(vial.Tests, &vialObject.Tests)
		if err != nil {
			return nil, &commonStructures.CommonError{
				Message:    fmt.Sprintf(commonConstants.ERROR_FAILED_TO_UNMARSHAL_JSON),
				StatusCode: http.StatusInternalServerError,
			}
		}

		response[vial.OrderId] = append(response[vial.OrderId], vialObject)
	}
	return response, nil
}

func (sampleService *SampleService) RejectSampleByBarcode(ctx context.Context, barcode string,
	requestBody structures.RejectSampleRequest) (string, []string, *commonStructures.CommonError) {

	barcode = strings.TrimSpace(barcode)
	if barcode == "" {
		return "", nil, &commonStructures.CommonError{
			StatusCode: http.StatusBadRequest,
			Message:    commonConstants.ERROR_BARCODE_REQUIRED,
		}
	}

	samples, sampleMetadatas, cErr := sampleService.SampleDao.GetSampleDataByBarcodeForRejection(barcode)
	if cErr != nil {
		return "", nil, cErr
	}
	sampleIdToSampleMetadataMap := map[uint]commonModels.SampleMetadata{}
	for _, sampleMetadata := range sampleMetadatas {
		sampleIdToSampleMetadataMap[sampleMetadata.SampleId] = sampleMetadata
	}

	cErr = validateSampleForSampleRejection(samples)
	if cErr != nil {
		return "", nil, cErr
	}

	omsTestIds, cErr := sampleService.RejectSample(ctx, requestBody.UserId, requestBody.LabId, requestBody.RejectionReason,
		samples, sampleMetadatas)
	if cErr != nil {
		return "", nil, cErr
	}

	return samples[0].OmsOrderId, omsTestIds, nil
}

func (sampleService *SampleService) RejectSamplePartiallyBySampleNumberAndTestId(ctx context.Context,
	requestBody structures.RejectSamplePartiallyRequest) (string, []string, *commonStructures.CommonError) {

	samples, sampleMetadatas, cErr := sampleService.SampleDao.GetSampleDataBySampleNumberAndTestId(
		requestBody.SampleNumber, requestBody.TestId)
	if cErr != nil {
		return "", nil, cErr
	}
	sampleIdToSampleMetadataMap := map[uint]commonModels.SampleMetadata{}
	for _, sampleMetadata := range sampleMetadatas {
		sampleIdToSampleMetadataMap[sampleMetadata.SampleId] = sampleMetadata
	}

	finalOmsTestIds := []string{}

	for _, sample := range samples {
		omsTestIds, cErr := sampleService.PartiallyRejectSample(ctx, requestBody.TestId, requestBody.UserId,
			requestBody.LabId, samples, sampleMetadatas, requestBody.RejectionReason)
		if cErr != nil {
			return "", nil, cErr
		}

		if sample.ParentSampleId == 0 {
			finalOmsTestIds = append(finalOmsTestIds, omsTestIds...)
		}
	}

	return samples[0].OmsOrderId, finalOmsTestIds, nil
}

func (sampleService *SampleService) GetTestDetailsForLisEventByVisitId(visitId string) (
	[]commonStructures.TestDetailsForLisEvent, *commonStructures.CommonError) {
	return sampleService.SampleDao.GetTestDetailsForLisEventByVisitId(visitId)
}

func (sampleService *SampleService) GetOmsTestDetailsByVisitId(visitId string) ([]commonStructures.OmsTestDetailsForLis,
	*commonStructures.CommonError) {
	return sampleService.SampleDao.GetOmsTestDetailsByVisitId(visitId)
}

func (sampleService *SampleService) GetVisitIdsByOmsOrderId(omsOrderId string) ([]string, *commonStructures.CommonError) {
	return sampleService.SampleDao.GetVisitIdsByOmsOrderId(omsOrderId)
}

func (sampleService *SampleService) GetVisitLabMapByOmsTestIds(omsTestIds []string) (
	map[string]uint, *commonStructures.CommonError) {
	return sampleService.SampleDao.GetVisitLabMapByOmsTestIds(omsTestIds)
}

func (sampleService *SampleService) BarcodesExistsInSystem(barcodes []string) (bool, *commonStructures.CommonError) {
	return sampleService.SampleDao.BarcodesExistsInSystem(barcodes)
}

func (sampleService *SampleService) GetVisitDetailsForTaskByOmsOrderId(omsOrderId string) (
	[]commonStructures.VisitDetailsForTask, *commonStructures.CommonError) {
	visits, cErr := sampleService.SampleDao.GetVisitDetailsForTaskByOmsOrderId(omsOrderId)
	if cErr != nil {
		return nil, cErr
	}

	filteredVisits, visitMap := []commonStructures.VisitDetailsForTask{}, map[string]bool{}
	for _, visit := range visits {
		if visit.VisitId != "" && !visitMap[visit.VisitId] {
			visitMap[visit.VisitId] = true
			filteredVisits = append(filteredVisits, visit)
		}
	}

	return filteredVisits, nil
}

func (sampleService *SampleService) IsSampleCollected(omsOrderId string) (bool, *commonStructures.CommonError) {
	return sampleService.SampleDao.IsSampleCollected(omsOrderId)
}

func (sampleService *SampleService) DeleteSamplesAndTestDetailsWithTx(ctx context.Context, tx *gorm.DB, omsOrderId string,
	toDeleteTestDetails []commonModels.TestDetail) *commonStructures.CommonError {
	if len(toDeleteTestDetails) == 0 {
		return nil
	}

	toDeleteTestIds, omsTestIdToTestDetailMap := make([]string, 0), map[string]commonModels.TestDetail{}
	for _, testDetail := range toDeleteTestDetails {
		toDeleteTestIds = append(toDeleteTestIds, testDetail.CentralOmsTestId)
		omsTestIdToTestDetailMap[testDetail.CentralOmsTestId] = testDetail
	}

	tsm, cErr := sampleService.TestSampleMappingService.GetTestSampleMappingByOrderIdAndTestIds(omsOrderId, toDeleteTestIds)
	if cErr != nil && cErr.Message != commonConstants.ERROR_NO_TEST_SAMPLE_MAPPING_FOUND {
		return cErr
	}

	allTsms, cErr := sampleService.TestSampleMappingService.GetTestSampleMappingsModelsByOrderId(omsOrderId)
	if cErr != nil {
		return cErr
	}

	// Map to keep track of sample numbers associated with each test
	testIdToSampleNumbers := make(map[string][]uint)

	// Map to keep track of tests associated with each sample number
	sampleNumberToTestIds := make(map[uint][]string)

	// Populate the maps
	for _, mapping := range tsm {
		if _, ok := testIdToSampleNumbers[mapping.OmsTestId]; !ok {
			testIdToSampleNumbers[mapping.OmsTestId] = []uint{}
		}
		testIdToSampleNumbers[mapping.OmsTestId] = append(testIdToSampleNumbers[mapping.OmsTestId], mapping.SampleNumber)
	}

	for _, mapping := range allTsms {
		if _, ok := sampleNumberToTestIds[mapping.SampleNumber]; !ok {
			sampleNumberToTestIds[mapping.SampleNumber] = []string{}
		}
		sampleNumberToTestIds[mapping.SampleNumber] = append(sampleNumberToTestIds[mapping.SampleNumber], mapping.OmsTestId)
	}

	// Process each test to delete
	for _, omsTestId := range toDeleteTestIds {
		// Get all sample numbers mapped to this test
		sampleNumbers := testIdToSampleNumbers[omsTestId]

		for _, sampleNumber := range sampleNumbers {
			// Get the sample to check its visit ID for LIS sync cancellation
			sample, cErr := sampleService.SampleDao.GetSampleByOrderIdAndSampleNumber(omsOrderId, sampleNumber)
			if cErr != nil {
				return cErr
			}

			// If the sample has a visit ID, we need to trigger LIS sync cancellation for this test
			if sample.VisitId != "" {
				testDetail := omsTestIdToTestDetailMap[omsTestId]

				// Trigger LIS sync cancellation for just this test
				cErr = sampleService.AttuneService.CancelLisSyncData(ctx, []commonModels.TestDetail{testDetail},
					sample.VisitId, sample)
				if cErr != nil {
					return cErr
				}
			}

			// If the sample has only this test mapped to it, delete the sample
			if len(sampleNumberToTestIds[sampleNumber]) == 1 {
				if cErr := sampleService.SampleDao.DeleteSampleByOrderIdAndSampleNumberWithTx(tx, omsOrderId, sampleNumber); cErr != nil {
					return cErr
				}
			}

			// Always delete the test-sample mapping
			if cErr := sampleService.TestSampleMappingService.DeleteTestSampleMappingByOrderIdTestIdAndSampleNumberWithTx(
				tx, omsOrderId, omsTestId, sampleNumber); cErr != nil {
				return cErr
			}
		}
	}

	return nil
}

func (sampleService *SampleService) DeleteTestSampleMappingForDeletedTestIds(ctx context.Context, tx *gorm.DB, omsOrderId string,
	toDeleteTestDetails []commonModels.TestDetail) *commonStructures.CommonError {
	toDeleteTestIds, omsTestIdToTestDetailMap := make([]string, 0), map[string]commonModels.TestDetail{}
	for _, testDetail := range toDeleteTestDetails {
		toDeleteTestIds = append(toDeleteTestIds, testDetail.CentralOmsTestId)
		omsTestIdToTestDetailMap[testDetail.CentralOmsTestId] = testDetail
	}

	tsm, cErr := sampleService.TestSampleMappingService.GetTestSampleMappingByOrderIdAndTestIds(omsOrderId, toDeleteTestIds)
	if cErr != nil && cErr.Message != commonConstants.ERROR_NO_TEST_SAMPLE_MAPPING_FOUND {
		return cErr
	}

	allTsms, cErr := sampleService.TestSampleMappingService.GetTestSampleMappingsModelsByOrderId(omsOrderId)
	if cErr != nil {
		return cErr
	}

	// Map to keep track of sample numbers associated with each test
	testIdToSampleNumbers := make(map[string][]uint)

	// Map to keep track of tests associated with each sample number
	sampleNumberToTestIds := make(map[uint][]string)

	// Populate the maps
	for _, mapping := range tsm {
		if _, ok := testIdToSampleNumbers[mapping.OmsTestId]; !ok {
			testIdToSampleNumbers[mapping.OmsTestId] = []uint{}
		}
		testIdToSampleNumbers[mapping.OmsTestId] = append(testIdToSampleNumbers[mapping.OmsTestId], mapping.SampleNumber)
	}

	for _, mapping := range allTsms {
		if _, ok := sampleNumberToTestIds[mapping.SampleNumber]; !ok {
			sampleNumberToTestIds[mapping.SampleNumber] = []string{}
		}
		sampleNumberToTestIds[mapping.SampleNumber] = append(sampleNumberToTestIds[mapping.SampleNumber], mapping.OmsTestId)
	}

	// Process each test to delete
	for _, omsTestId := range toDeleteTestIds {
		// Get all sample numbers mapped to this test
		sampleNumbers := testIdToSampleNumbers[omsTestId]

		for _, sampleNumber := range sampleNumbers {
			// Get the sample to check its visit ID for LIS sync cancellation
			sample, cErr := sampleService.SampleDao.GetSampleByOrderIdAndSampleNumber(omsOrderId, sampleNumber)
			if cErr != nil {
				return cErr
			}

			// If the sample has a visit ID, we need to trigger LIS sync cancellation for this test
			if sample.VisitId != "" {
				testDetail := omsTestIdToTestDetailMap[omsTestId]

				// Trigger LIS sync cancellation for just this test
				cErr = sampleService.AttuneService.CancelLisSyncData(ctx, []commonModels.TestDetail{testDetail},
					sample.VisitId, sample)
				if cErr != nil {
					return cErr
				}
			}

			// Always delete the test-sample mapping
			if cErr := sampleService.TestSampleMappingService.DeleteTestSampleMappingByOrderIdTestIdAndSampleNumberWithTx(
				tx, omsOrderId, omsTestId, sampleNumber); cErr != nil {
				return cErr
			}
		}
	}
	return nil
}

func (sampleService *SampleService) DeleteAllSamplesDataByOmsOrderIdWithTx(tx *gorm.DB, omsOrderId string) *commonStructures.CommonError {
	return sampleService.SampleDao.DeleteAllSamplesDataByOmsOrderIdWithTx(tx, omsOrderId)

}

func (sampleService *SampleService) RejectSample(ctx context.Context, userId, labId uint, rejectionReason string,
	samples []commonModels.Sample, sampleMetadatas []commonModels.SampleMetadata) ([]string, *commonStructures.CommonError) {

	for index := range samples {
		samples[index].Status = commonConstants.SampleRejected
		samples[index].UpdatedBy = userId
		samples[index].RejectionReason = rejectionReason
	}

	for index := range sampleMetadatas {
		sampleMetadatas[index].RejectedAt = commonUtils.GetCurrentTime()
		sampleMetadatas[index].RejectingLab = labId
		sampleMetadatas[index].UpdatedBy = userId
	}

	centralOmsTestIds := []string{}

	testDetails, _, cErr := sampleService.SampleDao.GetAllTestsAndSampleMappingsBySampleNumbers(
		[]uint{samples[0].SampleNumber}, samples[0].OmsOrderId)
	if cErr != nil {
		return centralOmsTestIds, cErr
	}

	for _, test := range testDetails {
		centralOmsTestIds = append(centralOmsTestIds, test.CentralOmsTestId)
	}

	for index := range samples {
		if samples[index].VisitId != "" {
			cErr = sampleService.AttuneService.CancelLisSyncData(ctx, testDetails, samples[index].VisitId,
				mappers.MapSampleSampleMetaToSampleInfo(samples[index], sampleMetadatas[index]))
			if cErr != nil {
				return centralOmsTestIds, cErr
			}
		}
		samples[index].VisitId = ""
		sampleMetadatas[index].LisSyncAt = nil
	}

	tx := sampleService.SampleDao.BeginTransaction()
	defer tx.Rollback()

	_, _, cErr = sampleService.SampleDao.UpdateSamplesAndSamplesMetadataWithTx(tx, samples, sampleMetadatas)
	if cErr != nil {
		return centralOmsTestIds, cErr
	}

	for _, sample := range samples {
		cErr = sampleService.TestSampleMappingService.RejectTestSampleMappingBySampleIdWithTx(tx, sample.Id, userId)
		if cErr != nil {
			return centralOmsTestIds, cErr
		}
	}

	tx.Commit()

	for _, sample := range samples {
		if sample.ParentSampleId == 0 {
			if checkIfTicketShouldBeSentBasedOnReason(rejectionReason) {
				go sampleService.CreateFreshDeskTicketForSampleRejection(sample.OmsOrderId,
					constants.SampleRejectedFreshDeskSubject, rejectionReason, testDetails)
			}

			omsTestStatusMap := map[string]string{}
			for _, testId := range centralOmsTestIds {
				omsTestStatusMap[testId] = commonConstants.TEST_STATUS_REQUESTED
			}
			go sampleService.PublishUpdateTestStatusEvent(omsTestStatusMap, sample.OmsOrderId, true, sample.OmsCityCode)

			go sampleService.PublishResetTatsEvent(centralOmsTestIds, sample.OmsCityCode)

			go sampleService.PublishAddSampleRejectedTagEvent(sample.OmsOrderId, sample.OmsCityCode)

			go sampleService.EtsService.GetAndPublishEtsTestEventForSampleRejection(context.Background(),
				sample.OmsOrderId, sample.SampleNumber)
		}
	}

	return centralOmsTestIds, nil
}

func (sampleService *SampleService) PartiallyRejectSample(ctx context.Context, omsTestId string, userId, labId uint,
	samples []commonModels.Sample, sampleMetadatas []commonModels.SampleMetadata,
	rejectionReason string) ([]string, *commonStructures.CommonError) {

	if omsTestId == "" {
		return nil, &commonStructures.CommonError{
			Message:    commonConstants.ERROR_TEST_ID_CANNOT_BE_ZERO,
			StatusCode: http.StatusBadRequest,
		}
	}

	testDetails, cErr := sampleService.TestDetailsService.GetTestDetailModelByOmsTestId(omsTestId)
	if cErr != nil {
		return nil, cErr
	}
	omsTestIds := []string{testDetails.CentralOmsTestId}

	sampleIdToSampleMetadataMap := map[uint]commonModels.SampleMetadata{}
	for _, sampleMetadata := range sampleMetadatas {
		sampleIdToSampleMetadataMap[sampleMetadata.SampleId] = sampleMetadata
	}

	currentTime := commonUtils.GetCurrentTime()
	for index := range samples {
		samples[index].Status = commonConstants.SamplePartiallyRejected
		samples[index].RejectionReason = rejectionReason
		samples[index].UpdatedBy = userId
		samples[index].UpdatedAt = currentTime
		sampleMetadata := sampleIdToSampleMetadataMap[samples[index].Id]
		sampleMetadata.RejectedAt = currentTime
		sampleMetadata.RejectingLab = labId
		sampleMetadata.UpdatedBy = userId
		sampleMetadata.UpdatedAt = currentTime
		sampleIdToSampleMetadataMap[samples[index].Id] = sampleMetadata
	}

	for _, sample := range samples {
		if sample.VisitId != "" {
			cErr = sampleService.AttuneService.CancelLisSyncData(ctx, []commonModels.TestDetail{testDetails}, sample.VisitId,
				mappers.MapSampleSampleMetaToSampleInfo(sample, sampleIdToSampleMetadataMap[sample.Id]))
			if cErr != nil {
				return omsTestIds, cErr
			}
			break
		}
	}

	sampleMetadatas = []commonModels.SampleMetadata{}
	for _, sampleMetadata := range sampleIdToSampleMetadataMap {
		sampleMetadatas = append(sampleMetadatas, sampleMetadata)
	}

	tx := sampleService.SampleDao.BeginTransaction()
	defer tx.Rollback()

	_, _, cErr = sampleService.SampleDao.UpdateSamplesAndSamplesMetadataWithTx(tx, samples, sampleMetadatas)
	if cErr != nil {
		return omsTestIds, cErr
	}

	cErr = sampleService.TestSampleMappingService.RejectTestSampleMappingsByOmsOrderIdSampleNumberAndTestIdWithTx(tx,
		samples[0].OmsOrderId, samples[0].SampleNumber, omsTestId, rejectionReason)
	if cErr != nil {
		return omsTestIds, cErr
	}

	tx.Commit()

	go sampleService.PublishResetTatsEvent(omsTestIds, samples[0].OmsCityCode)

	go sampleService.EtsService.GetAndPublishEtsTestEventForPartialRejection(context.Background(), samples[0].OmsOrderId, samples[0].SampleNumber,
		omsTestId)

	if checkIfTicketShouldBeSentBasedOnReason(rejectionReason) {
		go sampleService.CreateFreshDeskTicketForSampleRejection(samples[0].OmsOrderId,
			constants.SamplePartiallyRejectedFreshDeskSubject, rejectionReason, []commonModels.TestDetail{testDetails})
	}

	allSampleRejected, err := sampleService.TestSampleMappingService.CheckIfAllTestsRejectedByOrderIdTestIdAndSampleNumber(
		samples[0].OmsOrderId, omsTestId, samples[0].SampleNumber)
	if err != nil {
		return omsTestIds, &commonStructures.CommonError{
			Message:    err.Error(),
			StatusCode: http.StatusInternalServerError,
		}
	}

	if allSampleRejected {
		for index := range samples {
			samples[index].Status = commonConstants.SampleRejected
			samples[index].VisitId = ""
			samples[index].UpdatedBy = commonConstants.CitadelSystemId
			samples[index].UpdatedAt = currentTime
			sampleMetadata := sampleIdToSampleMetadataMap[samples[index].Id]
			sampleMetadata.LisSyncAt = nil
			sampleMetadata.UpdatedBy = commonConstants.CitadelSystemId
			sampleMetadata.UpdatedAt = currentTime
			sampleIdToSampleMetadataMap[samples[index].Id] = sampleMetadata
		}

		_, _, cErr = sampleService.SampleDao.UpdateSamplesAndSamplesMetadata(samples, sampleMetadatas)
		if cErr != nil {
			return omsTestIds, cErr
		}
	}

	go sampleService.PublishAddSampleRejectedTagEvent(samples[0].OmsOrderId, samples[0].OmsCityCode)

	return omsTestIds, nil
}

func (sampleService *SampleService) GetCollectionSequenceForSample(ctx context.Context,
	collectionSequence commonStructures.CollectionSequenceResponse, masterTestIdToProcessingLabIdMap map[uint]uint) []structures.CollectionSequence {
	// return sequence -> [labId -> [sampleId -> []testIds]]
	sequences := []structures.CollectionSequence{}
	for _, sequence := range collectionSequence.Collections {
		labToSampleTestMapping := getLabToSampleTestMapping(sequence, masterTestIdToProcessingLabIdMap)
		commonUtils.AddLog(ctx, commonConstants.DEBUG_LEVEL, commonUtils.GetCurrentFunctionName(), map[string]interface{}{
			"labToSampleTestMapping": labToSampleTestMapping,
		}, nil)
		sequences = append(sequences, structures.CollectionSequence{
			Sequence:        sequence.Sequence,
			SequenceDetails: labToSampleTestMapping,
		})
	}
	return sequences
}

func (sampleService *SampleService) GetSamples(ctx context.Context, orderDetails commonModels.OrderDetails,
	testDetails []commonModels.TestDetail) ([]structures.SampleTestMappingDetails, *commonStructures.CommonError) {

	sampleNumber := sampleService.SampleDao.GetMaxSampleNumberByOmsOrderId(orderDetails.OmsOrderId)
	masterTestIds := getOrderedMasterTestIds(testDetails)

	collectionSequence, cErr := sampleService.CdsService.GetCollectionSequence(ctx, getCollectionSequenceRequest(
		orderDetails, masterTestIds))
	if cErr != nil {
		return []structures.SampleTestMappingDetails{}, cErr
	}

	masterTestsMap := sampleService.CdsService.GetMasterTestAsMap(ctx)

	masterTestIdToProcessingLabIdMap := utils.GetMasterTestToProcessingLabMap(testDetails)
	collectionSequenceForSample := sampleService.GetCollectionSequenceForSample(ctx, collectionSequence,
		masterTestIdToProcessingLabIdMap)
	stmds := sampleService.GetSampleTestMappingDetails(collectionSequenceForSample, orderDetails, testDetails, sampleNumber,
		masterTestsMap)
	return stmds, nil
}

func (sampleService *SampleService) GetSamplesWithTx(ctx context.Context, tx *gorm.DB, orderDetails commonModels.OrderDetails,
	testDetails []commonModels.TestDetail) ([]structures.SampleTestMappingDetails, *commonStructures.CommonError) {
	sampleNumber := sampleService.SampleDao.GetMaxSampleNumberByOmsOrderIdWithTx(tx, orderDetails.OmsOrderId)
	masterTestIds := getOrderedMasterTestIds(testDetails)

	collectionSequence, cErr := sampleService.CdsService.GetCollectionSequence(ctx, getCollectionSequenceRequest(
		orderDetails, masterTestIds))
	if cErr != nil {
		return []structures.SampleTestMappingDetails{}, cErr
	}

	masterTestsMap := sampleService.CdsService.GetMasterTestAsMap(ctx)

	masterTestIdToProcessingLabIdMap := utils.GetMasterTestToProcessingLabMap(testDetails)
	commonUtils.AddLog(ctx, commonConstants.DEBUG_LEVEL, commonUtils.GetCurrentFunctionName(), map[string]interface{}{
		"masterTestIdToProcessingLabIdMap": masterTestIdToProcessingLabIdMap,
	}, nil)
	collectionSequenceForSample := sampleService.GetCollectionSequenceForSample(ctx, collectionSequence,
		masterTestIdToProcessingLabIdMap)
	commonUtils.AddLog(ctx, commonConstants.DEBUG_LEVEL, commonUtils.GetCurrentFunctionName(), map[string]interface{}{
		"collectionSequenceForSample": collectionSequenceForSample,
	}, nil)
	stmds := sampleService.GetSampleTestMappingDetails(collectionSequenceForSample, orderDetails, testDetails, sampleNumber,
		masterTestsMap)
	return stmds, nil
}

func (sampleService *SampleService) GetSampleInfoForCreate(sampleNumber, vialId, sequenceNumber uint,
	orderDetails commonModels.OrderDetails) commonStructures.SampleInfo {
	sampleInfo := commonStructures.SampleInfo{
		SampleNumber:             sampleNumber,
		OmsCityCode:              orderDetails.CityCode,
		OmsOrderId:               orderDetails.OmsOrderId,
		OmsRequestId:             orderDetails.OmsRequestId,
		VialTypeId:               vialId,
		DestinationLabId:         orderDetails.ServicingLabId,
		CollectionSequenceNumber: sequenceNumber,
		Status:                   commonConstants.SampleDefault,
		CreatedAt:                commonUtils.GetCurrentTime(),
		CreatedBy:                commonConstants.CitadelSystemId,
		UpdatedAt:                commonUtils.GetCurrentTime(),
		UpdatedBy:                commonConstants.CitadelSystemId,
	}

	return sampleInfo
}

func (sampleService *SampleService) GetTestSampleMappingsForCreateWithCollections(baseSampleNumber, vialId uint,
	masterTestIds []uint, orderDetails commonModels.OrderDetails, testDetails []commonModels.TestDetail,
	masterTestsMap map[uint]commonStructures.CdsTestMaster) []commonModels.TestSampleMapping {

	masterTestIdToOmsTestIdMap := make(map[uint][]string)
	tsms := []commonModels.TestSampleMapping{}

	// Build mapping from master test ID to OMS test IDs
	for _, testDetail := range testDetails {
		if _, ok := masterTestIdToOmsTestIdMap[testDetail.MasterTestId]; !ok {
			masterTestIdToOmsTestIdMap[testDetail.MasterTestId] = []string{}
		}
		masterTestIdToOmsTestIdMap[testDetail.MasterTestId] = append(masterTestIdToOmsTestIdMap[testDetail.MasterTestId],
			testDetail.CentralOmsTestId)
	}

	// For each master test, create mappings based on its CollectionsCount
	for _, masterTestId := range masterTestIds {
		omsTestIds := masterTestIdToOmsTestIdMap[masterTestId]

		// Get CollectionsCount for this master test, default to 1
		collectionsCount := uint(1)
		if masterTest, exists := masterTestsMap[masterTestId]; exists && masterTest.CollectionsCount > 0 {
			collectionsCount = masterTest.CollectionsCount
		}

		// Create test sample mappings for each collection of this master test
		for index := uint(0); index < collectionsCount; index++ {
			currentSampleNumber := baseSampleNumber + index + 1
			for _, omsTestId := range omsTestIds {
				tsm := commonModels.TestSampleMapping{
					OmsCityCode:  orderDetails.CityCode,
					OmsTestId:    omsTestId,
					SampleNumber: currentSampleNumber,
					VialTypeId:   vialId,
					OmsOrderId:   orderDetails.OmsOrderId,
				}
				tsm.CreatedBy = commonConstants.CitadelSystemId
				tsm.UpdatedBy = commonConstants.CitadelSystemId
				tsms = append(tsms, tsm)
			}
		}
	}
	return tsms
}

func (sampleService *SampleService) GetSampleTestMappingDetails(collectionSequence []structures.CollectionSequence,
	orderDetails commonModels.OrderDetails, testDetails []commonModels.TestDetail, sampleNumber uint,
	masterTestsMap map[uint]commonStructures.CdsTestMaster) []structures.SampleTestMappingDetails {
	stmds := []structures.SampleTestMappingDetails{}
	for _, sequence := range collectionSequence {
		for _, vialTestMapping := range sequence.SequenceDetails {
			for vialId, masterTestIds := range vialTestMapping {
				// Determine the number of samples to create based on maximum CollectionsCount
				maxCollectionsCount := uint(1) // Default to 1 sample if no CollectionsCount

				// Find the maximum CollectionsCount among all master tests in this vial
				for _, masterTestId := range masterTestIds {
					if masterTest, exists := masterTestsMap[masterTestId]; exists && masterTest.CollectionsCount > 0 {
						if masterTest.CollectionsCount > maxCollectionsCount {
							maxCollectionsCount = masterTest.CollectionsCount
						}
					}
				}

				// Generate test sample mappings based on each master test's CollectionsCount
				baseSampleNumber := sampleNumber
				allTestSampleMappings := sampleService.GetTestSampleMappingsForCreateWithCollections(
					baseSampleNumber, vialId, masterTestIds, orderDetails, testDetails, masterTestsMap)

				// Create samples based on the maximum CollectionsCount
				for index := uint(0); index < maxCollectionsCount; index++ {
					currentSampleNumber := baseSampleNumber + index + 1

					// Filter test sample mappings for this specific sample number
					var testSampleMappingsForCurrentSample []commonModels.TestSampleMapping
					for _, tsm := range allTestSampleMappings {
						if tsm.SampleNumber == currentSampleNumber {
							testSampleMappingsForCurrentSample = append(testSampleMappingsForCurrentSample, tsm)
						}
					}

					stmds = append(stmds, structures.SampleTestMappingDetails{
						SampleInfo: sampleService.GetSampleInfoForCreate(currentSampleNumber, vialId, uint(sequence.Sequence),
							orderDetails),
						TestSampleMappingInfos: testSampleMappingsForCurrentSample,
					})
				}
				sampleNumber += maxCollectionsCount
			}
		}
	}
	return stmds
}

func (sampleService *SampleService) GetOrderAndTestDetails(omsOrderId string, testId string, isOmsTestId bool) (
	commonModels.OrderDetails, []commonModels.TestDetail, *commonStructures.CommonError) {

	orderDetails, testDetails := commonModels.OrderDetails{}, []commonModels.TestDetail{}
	if omsOrderId == "" {
		return orderDetails, testDetails, &commonStructures.CommonError{
			Message:    commonConstants.ERROR_ORDER_ID_CANNOT_BE_ZERO,
			StatusCode: http.StatusBadRequest,
		}
	}
	var cErr *commonStructures.CommonError
	orderDetails, cErr = sampleService.OrderDetailsService.GetOrderDetailsByOmsOrderId(omsOrderId)
	if cErr != nil {
		return orderDetails, testDetails, cErr
	}

	if testId != "" {
		if isOmsTestId {
			testDetail, cErr := sampleService.TestDetailsService.GetTestDetailModelByOmsTestId(testId)
			if cErr != nil {
				return orderDetails, testDetails, cErr
			}
			testDetails = append(testDetails, testDetail)
		} else {
			testDetailsId := commonUtils.ConvertStringToUint(testId)
			if testDetailsId == 0 {
				return orderDetails, testDetails, &commonStructures.CommonError{
					Message:    commonConstants.ERROR_TEST_ID_CANNOT_BE_ZERO,
					StatusCode: http.StatusBadRequest,
				}
			}
			testDetail, cErr := sampleService.TestDetailsService.GetTestDetailModelById(testDetailsId)
			if cErr != nil {
				return orderDetails, testDetails, cErr
			}
			testDetails = append(testDetails, testDetail)
		}
		return orderDetails, testDetails, nil
	} else {
		testDetails, cErr = sampleService.TestDetailsService.GetTestDetailsByOmsOrderId(omsOrderId)
		if cErr != nil {
			return orderDetails, testDetails, cErr
		}
	}

	return orderDetails, testDetails, nil
}

func (sampleService *SampleService) GetOrderAndTestDetailsWithTx(tx *gorm.DB, omsOrderId string) (
	commonModels.OrderDetails, []commonModels.TestDetail, *commonStructures.CommonError) {

	orderDetails, testDetails := commonModels.OrderDetails{}, []commonModels.TestDetail{}
	if omsOrderId == "" {
		return orderDetails, testDetails, &commonStructures.CommonError{
			Message:    commonConstants.ERROR_ORDER_ID_CANNOT_BE_ZERO,
			StatusCode: http.StatusBadRequest,
		}
	}

	orderDetails, cErr := sampleService.OrderDetailsService.GetOrderDetailsByOmsOrderIdWithTx(tx, omsOrderId)
	if cErr != nil {
		return orderDetails, testDetails, cErr
	}

	testDetails, cErr = sampleService.TestDetailsService.GetTestDetailsByOmsOrderIdWithTx(tx, omsOrderId)
	if cErr != nil {
		return orderDetails, testDetails, cErr
	}

	tsm, cErr := sampleService.TestSampleMappingService.GetTestSampleMappingByOrderIdWithTx(tx, omsOrderId)
	if cErr != nil && cErr.StatusCode != http.StatusNotFound {
		return orderDetails, testDetails, cErr
	}

	omsTestIds := make([]string, 0)
	for _, tsm := range tsm {
		omsTestIds = append(omsTestIds, tsm.OmsTestId)
	}

	filteredTestDetails := []commonModels.TestDetail{}
	for _, testDetail := range testDetails {
		if !commonUtils.SliceContainsString(omsTestIds, testDetail.CentralOmsTestId) {
			filteredTestDetails = append(filteredTestDetails, testDetail)
		}
	}

	return orderDetails, filteredTestDetails, nil
}

func (sampleService *SampleService) GetSamplesForTests(testIds []string) (
	[]commonModels.Sample, []commonModels.SampleMetadata, *commonStructures.CommonError) {
	return sampleService.SampleDao.GetSamplesForTests(testIds)
}

func getLabToSampleTestMapping(sequence commonStructures.CollectionDetails,
	masterTestIdToProcessingLabIdMap map[uint]uint) map[uint]map[uint][]uint {

	labToSampleTestMapping := make(map[uint]map[uint][]uint)
	testVialsMap := sequence.TestVialMapping
	masterTestIds := sequence.MasterTestIds

	for _, masterTestId := range masterTestIds {
		if _, ok := labToSampleTestMapping[masterTestIdToProcessingLabIdMap[masterTestId]]; !ok {
			vialTestsMap := make(map[uint][]uint)
			for _, vialId := range testVialsMap[masterTestId] {
				vialTestsMap[vialId] = []uint{masterTestId}
			}
			labToSampleTestMapping[masterTestIdToProcessingLabIdMap[masterTestId]] = vialTestsMap
		} else {
			for _, vialId := range testVialsMap[masterTestId] {
				if _, ok := labToSampleTestMapping[masterTestIdToProcessingLabIdMap[masterTestId]][vialId]; !ok {
					labToSampleTestMapping[masterTestIdToProcessingLabIdMap[masterTestId]][vialId] = []uint{masterTestId}
				} else {
					labToSampleTestMapping[masterTestIdToProcessingLabIdMap[masterTestId]][vialId] =
						append(labToSampleTestMapping[masterTestIdToProcessingLabIdMap[masterTestId]][vialId], masterTestId)
				}
			}
		}
	}
	return labToSampleTestMapping
}

func getOrderedMasterTestIds(testDetails []commonModels.TestDetail) []uint {
	masterTestIds := make([]uint, 0)
	for _, testDetail := range testDetails {
		masterTestIds = append(masterTestIds, testDetail.MasterTestId)
	}
	masterTestIds = commonUtils.CreateUniqueSliceUint(masterTestIds)
	sort.Slice(masterTestIds, func(i, j int) bool {
		return masterTestIds[i] < masterTestIds[j]
	})
	return masterTestIds
}

func getCollectionSequenceRequest(orderDetails commonModels.OrderDetails,
	masterTestIds []uint) commonStructures.CollectionSequenceRequest {
	return commonStructures.CollectionSequenceRequest{
		CityCode: orderDetails.CityCode,
		OrderDetails: []commonStructures.OrderDetails{{
			Indentifier:   fmt.Sprint(orderDetails.OmsOrderId),
			MasterTestIds: masterTestIds,
		}},
		CreateSamples: true,
	}

}

func (sampleService *SampleService) UpdateSamplesAndSamplesMetadataWithTx(tx *gorm.DB,
	samples []commonModels.Sample, samplesMetadata []commonModels.SampleMetadata) (
	[]commonModels.Sample, []commonModels.SampleMetadata, *commonStructures.CommonError) {
	return sampleService.SampleDao.UpdateSamplesAndSamplesMetadataWithTx(tx, samples, samplesMetadata)
}

func (sampleService *SampleService) GetLisSyncDataByVisitId(ctx context.Context, visitId string) (
	commonStructures.LisSyncDetails, *commonStructures.CommonError) {

	sample, cErr := sampleService.SampleDao.GetSampleByVisitId(visitId)
	if cErr != nil {
		return commonStructures.LisSyncDetails{}, cErr
	}
	if sample.VisitId == "" {
		return commonStructures.LisSyncDetails{}, &commonStructures.CommonError{
			Message:    commonConstants.ERROR_VISIT_ID_NOT_FOUND,
			StatusCode: http.StatusBadRequest,
		}
	}

	lisSyncDetails, attuneOrderResponse, cErr := sampleService.AttuneService.GetLisSyncData(ctx, visitId,
		commonConstants.AttuneReportWithStationery)
	if cErr != nil {
		return lisSyncDetails, cErr
	}

	sampleService.PublishLisDataEvent(attuneOrderResponse)

	return lisSyncDetails, nil
}

func (sampleService *SampleService) MarkSampleAsEmedicNotCollected(omsRequestId string) *commonStructures.CommonError {
	return sampleService.SampleDao.MarkSampleAsEmedicNotCollected(omsRequestId)
}

func (sampleService *SampleService) ReMarkSampleDefaultEmedicNotCollected(omsRequestId string, taskSequence uint) *commonStructures.CommonError {
	return sampleService.SampleDao.ReMarkSampleDefaultEmedicNotCollected(omsRequestId, taskSequence)
}

func (sampleService *SampleService) SendAlertForSamplesToBeCollectedLater(omsRequestId string) error {
	orderDetails, _ := sampleService.OrderDetailsService.GetOrderDetailsByOmsRequestId(omsRequestId)
	for _, orderDetail := range orderDetails {
		sampleInfos, cErr := sampleService.SampleDao.GetSamplesByOmsOrderId(orderDetail.OmsOrderId)
		if cErr != nil {
			continue
		}
		filteredSampleInfos := []commonStructures.SampleInfo{}
		for _, sampleInfo := range sampleInfos {
			if sampleInfo.CollectLaterReason != "" && sampleInfo.Status == commonConstants.SampleNotCollectedEmedic {
				filteredSampleInfos = append(filteredSampleInfos, sampleInfo)
			}
		}
		if len(filteredSampleInfos) <= 0 {
			continue
		}
		sampleNumbers := []uint{}
		for _, a := range filteredSampleInfos {
			exists := false
			for _, item := range sampleNumbers {
				if item == a.SampleNumber {
					exists = true
					break
				}
			}
			if !exists {
				sampleNumbers = append(sampleNumbers, a.SampleNumber)
			}
		}
		go sampleService.CreateFreshDeskTicketForSampleNotReceived(orderDetail, sampleNumbers)
	}
	return nil
}

func (sampleService *SampleService) UpdateSampleDetailsPostTaskCompletion(ctx context.Context,
	requestBody commonStructures.UpdateSampleDetailsPostTaskCompletionRequest) *commonStructures.CommonError {
	cErr := sampleService.MarkSampleAsEmedicNotCollected(requestBody.RequestId)
	if cErr != nil && cErr.Message != commonConstants.ERROR_SAMPLE_NOT_FOUND {
		commonUtils.AddLog(ctx, commonConstants.ERROR_LEVEL, commonUtils.GetCurrentFunctionName(), map[string]interface{}{
			"error": cErr.Message}, nil)
	}

	updateSampleCollectedReq := commonStructures.SampleCollectedRequest{
		RequestId:      requestBody.RequestId,
		OmsTaskType:    requestBody.OmsTaskType,
		OmsTaskId:      requestBody.OmsTaskId,
		CollectedAt:    requestBody.CollectedAt,
		CollectionType: requestBody.CollectionType,
		IsB2c:          requestBody.IsB2c,
	}
	cErr = sampleService.UpdateSampleCollected(updateSampleCollectedReq)
	if cErr != nil && cErr.Message != commonConstants.ERROR_SAMPLE_NOT_FOUND {
		commonUtils.AddLog(ctx, commonConstants.ERROR_LEVEL, commonUtils.GetCurrentFunctionName(), map[string]interface{}{
			"error": cErr.Message}, nil)
	}

	err := sampleService.SendAlertForSamplesToBeCollectedLater(requestBody.RequestId)
	if err != nil {
		commonUtils.AddLog(ctx, commonConstants.ERROR_LEVEL, commonUtils.GetCurrentFunctionName(), map[string]interface{}{
			"error": err.Error(),
		}, nil)
	}

	cErr = sampleService.ReMarkSampleDefaultEmedicNotCollected(requestBody.RequestId, requestBody.TaskSequence)
	if cErr != nil && cErr.Message != commonConstants.ERROR_SAMPLE_NOT_FOUND {
		commonUtils.AddLog(ctx, commonConstants.ERROR_LEVEL, commonUtils.GetCurrentFunctionName(), map[string]interface{}{
			"error": cErr.Message}, nil)
	}

	cErr = sampleService.SampleDao.CollectionPortalMarkAccessionAsAccessioned(true, 0, requestBody.RequestId)
	if cErr != nil {
		commonUtils.AddLog(ctx, commonConstants.ERROR_LEVEL, commonUtils.GetCurrentFunctionName(), map[string]interface{}{
			"error": cErr.Message,
		}, nil)
	}

	return nil
}

func (sampleService *SampleService) CollectionPortalMarkAccessionAsAccessioned(
	requestBody structures.MarkAccessionAsAccessionedRequest) *commonStructures.CommonError {
	return sampleService.SampleDao.CollectionPortalMarkAccessionAsAccessioned(requestBody.IsWebhook, requestBody.SampleId,
		requestBody.RequestId)
}

func (sampleService *SampleService) UpdateTaskSequenceForSample(
	reqBody commonStructures.UpdateTaskSequenceRequest) *commonStructures.CommonError {
	if reqBody.RequestId == "" {
		return &commonStructures.CommonError{
			Message:    constants.ERROR_REQUEST_ID_REQUIRED,
			StatusCode: http.StatusBadRequest,
		}
	}
	if reqBody.TaskId == 0 {
		return &commonStructures.CommonError{
			Message:    constants.ERROR_TASK_ID_REQUIRED,
			StatusCode: http.StatusBadRequest,
		}
	}
	if len(reqBody.TestIds) == 0 {
		return &commonStructures.CommonError{
			Message:    constants.ERROR_TEST_IDS_REQUIRED,
			StatusCode: http.StatusBadRequest,
		}
	}

	return sampleService.SampleDao.UpdateTaskSequenceForSample(reqBody.RequestId, reqBody.TaskId, reqBody.TestIds)
}

func (sampleService *SampleService) RemapSamplesToNewTaskSequence(
	requestBody structures.RemapSamplesRequest) *commonStructures.CommonError {
	if len(requestBody.TestIds) == 0 {
		return &commonStructures.CommonError{
			Message:    constants.ERROR_TEST_IDS_REQUIRED,
			StatusCode: http.StatusBadRequest,
		}
	}

	testSampleMaps, cErr := sampleService.TestSampleMappingService.GetTestSampleMappingByTestIds(requestBody.TestIds)
	if cErr != nil {
		return cErr
	}

	orderIdToSampleUpdateMap := map[string][]uint{}
	for _, testSampleMap := range testSampleMaps {
		if !commonUtils.SliceContainsUint(orderIdToSampleUpdateMap[testSampleMap.OmsOrderId], testSampleMap.SampleNumber) {
			orderIdToSampleUpdateMap[testSampleMap.OmsOrderId] = append(orderIdToSampleUpdateMap[testSampleMap.OmsOrderId],
				testSampleMap.SampleNumber)
		}
	}

	return sampleService.SampleDao.RemapSamplesToNewTaskSequence(orderIdToSampleUpdateMap, requestBody.NewSequence)
}

func (sampleService *SampleService) UpdateSampleDetailsForReschedule(
	reqBody structures.UpdateSampleDetailsForRescheduleRequest) *commonStructures.CommonError {
	if reqBody.RequestId == "" {
		return &commonStructures.CommonError{
			Message:    constants.ERROR_REQUEST_ID_REQUIRED,
			StatusCode: http.StatusBadRequest,
		}
	}

	return sampleService.SampleDao.UpdateSampleDetailsForReschedule(reqBody)
}

func (sampleService *SampleService) AddCollectedVolumeToSample(requestBody structures.AddCollectedVolumneRequest) (
	structures.AddVolumeResponse, *commonStructures.CommonError) {
	response := structures.AddVolumeResponse{}
	if requestBody.SampleId == 0 {
		return response, &commonStructures.CommonError{
			Message:    commonConstants.ERROR_SAMPLE_ID_CANNOT_BE_ZERO,
			StatusCode: http.StatusBadRequest,
		}
	}

	cErr := sampleService.SampleDao.AddCollectedVolumeToSample(requestBody.SampleId, requestBody.Volume)
	if cErr != nil {
		return response, cErr
	}

	sample, _, cErr := sampleService.SampleDao.GetSampleDataBySampleId(requestBody.SampleId)
	if cErr != nil {
		return response, cErr
	}
	response.RequestId = sample.OmsRequestId
	response.OrderId = sample.OmsOrderId

	return response, nil
}

func (sampleService *SampleService) UpdateSrfIdToLis(ctx context.Context,
	orderDetails commonModels.OrderDetails) *commonStructures.CommonError {
	samples := sampleService.SampleDao.GetCovidTestSamples(orderDetails.OmsOrderId)
	if len(samples) == 0 {
		return &commonStructures.CommonError{
			Message:    commonConstants.ERROR_NO_COVID_TEST_VISIT_IDS_FOUND,
			StatusCode: http.StatusNotFound,
		}
	}
	for _, sample := range samples {
		if sample.VisitId == "" {
			continue
		}
		cErr := sampleService.AttuneService.UpdateSrfIdToAttune(ctx, sample, orderDetails.SrfId)
		if cErr != nil {
			commonUtils.AddLog(ctx, commonConstants.ERROR_LEVEL, commonUtils.GetCurrentFunctionName(), map[string]interface{}{
				"error": cErr.Message,
			}, nil)
			return cErr
		}
	}
	return nil
}

func (sampleService *SampleService) GetSamplesForDelayedReverseLogisticsDashboard(ctx context.Context,
	cityCode string) []structures.DelayedReverseLogisticsSamplesResponse {
	response := []structures.DelayedReverseLogisticsSamplesResponse{}
	delayedSamples, cErr := sampleService.SampleDao.GetSamplesForDelayedReverseLogistics(
		constants.ReverseLogisticsNormalTatForDashboard,
		constants.ReverseLogisticsCampTatForDashboard,
		constants.ReverseLogisticsInclinicTatForDashboard,
		constants.ReverseLogisticsDelayDaysForDashboard,
	)
	if cErr != nil {
		return response
	}

	masterVialTypeMap := sampleService.CdsService.GetMasterVialTypeAsMap(ctx)

	for _, sample := range delayedSamples {
		if (sample.SampleCollectedTime == nil) || (cityCode != "" && sample.CityCode != cityCode) {
			continue
		}
		response = append(response, structures.DelayedReverseLogisticsSamplesResponse{
			RequestId:            sample.RequestId,
			OrderId:              sample.OrderId,
			TrfId:                sample.TrfId,
			PatientName:          sample.PatientName,
			PatientAge:           commonUtils.GetAgeYearsFromDob(*sample.PatientAge),
			PatientGender:        sample.PatientGender,
			Barcode:              sample.Barcode,
			VialTypeId:           sample.VialTypeId,
			SampleCollectedTime:  sample.SampleCollectedTime,
			LogisticsMinuteSpent: getLogisticsMinutesSpentBySampleCollectedAt(sample.SampleCollectedTime),
			VialTypeName:         masterVialTypeMap[sample.VialTypeId].VialType,
			VialColor:            masterVialTypeMap[sample.VialTypeId].VialColor,
		})
	}

	return response
}

func (sampleService *SampleService) SendSlackAlertForDelayedReverseLogisticsTat(ctx context.Context) {
	delayedSamples, cErr := sampleService.SampleDao.GetSamplesForDelayedReverseLogistics(
		constants.ReverseLogisticsNormalTatForSlack,
		constants.ReverseLogisticsCampTatForSlack,
		constants.ReverseLogisticsInclinicTatForSlack,
		constants.ReverseLogisticsDelayDaysForSlack,
	)
	if cErr != nil {
		return
	}

	if len(delayedSamples) == 0 {
		return
	}

	vialTypeMap := sampleService.CdsService.GetMasterVialTypeAsMap(ctx)

	// Group samples by city code
	samplesByCity := make(map[string][]structures.DelayedReverseLogisticsSamplesDbStruct)
	for _, sample := range delayedSamples {
		cityCode := sample.CityCode
		if cityCode == "" {
			cityCode = "UNKNOWN"
		}
		samplesByCity[cityCode] = append(samplesByCity[cityCode], sample)
	}

	// Send separate alerts for each city
	for cityCode, citySamples := range samplesByCity {
		// Send samples in batches to avoid hitting Slack's 50 block limit
		maxSamplesPerMessage := commonConstants.MaxBlocksForSlackMessage

		for i := 0; i < len(citySamples); i += maxSamplesPerMessage {
			end := i + maxSamplesPerMessage
			if end > len(citySamples) {
				end = len(citySamples)
			}
			batch := citySamples[i:end]

			messageNumber := (i / maxSamplesPerMessage) + 1
			totalMessages := (len(citySamples) + maxSamplesPerMessage - 1) / maxSamplesPerMessage

			headerText := fmt.Sprintf(":alert: Sample Loss Risk Alert - %s", cityCode)
			if totalMessages > 1 {
				headerText = fmt.Sprintf(":alert: Sample Loss Risk Alert - %s (%d/%d)", cityCode, messageNumber, totalMessages)
			}

			blocks := []map[string]interface{}{
				{
					"type": "header",
					"text": map[string]interface{}{
						"type":  "plain_text",
						"text":  headerText,
						"emoji": true,
					},
				},
				{
					"type": "section",
					"fields": []map[string]interface{}{
						{"type": "mrkdwn", "text": "*Vial Type / Order ID / Sample Barcode / City*"},
						{"type": "mrkdwn", "text": "*Time since Collection (hh:mm)*"},
					},
				},
			}

			for _, sample := range batch {
				vialType := vialTypeMap[sample.VialTypeId].VialType
				sampleDetail := fmt.Sprintf("%s / %s / %s / %s", vialType, sample.OrderId, sample.Barcode, sample.CityCode)
				omsOrderIdString := commonUtils.GetStringOrderIdWithoutStringPart(sample.OrderId)
				omsOrderUrl := fmt.Sprintf("%s/request/%s/order/%s",
					commonUtils.GetOmsBaseDomain(sample.CityCode), sample.RequestId, omsOrderIdString)
				blocks = append(blocks, map[string]interface{}{
					"type": "section",
					"fields": []map[string]interface{}{
						{"type": "mrkdwn", "text": fmt.Sprintf("<%s|%s>", omsOrderUrl, sampleDetail)},
						{"type": "mrkdwn", "text": getLogisticsMinutesSpentBySampleCollectedAt(sample.SampleCollectedTime)},
					},
				})
			}

			err := sampleService.SlackClient.SendToSlackDirectly(ctx, commonConstants.SlackSampleLossRiskAlertChannel, blocks)
			if err != nil {
				commonUtils.AddLog(ctx, commonConstants.ERROR_LEVEL, commonUtils.GetCurrentFunctionName(), map[string]interface{}{
					"error":    err.Error(),
					"cityCode": cityCode,
					"batch":    messageNumber,
				}, nil)
			}
		}
	}
}

func (sampleService *SampleService) SendSlackAlertForDelayedInterlabLogisticsTat(ctx context.Context) {
	sampleInfos, cErr := sampleService.SampleDao.GetSamplesForDelayedInterlabLogistics()
	if cErr != nil {
		commonUtils.AddLog(ctx, commonConstants.ERROR_LEVEL, commonUtils.GetCurrentFunctionName(), map[string]interface{}{
			"error": cErr.Message,
		}, nil)
		return
	}

	if len(sampleInfos) == 0 {
		return
	}

	vialTypeMap := sampleService.CdsService.GetMasterVialTypeAsMap(ctx)

	// Group samples by city code
	samplesByCity := make(map[string][]commonStructures.SampleInfo)
	for _, sample := range sampleInfos {
		cityCode := sample.OmsCityCode
		if cityCode == "" {
			cityCode = "UNKNOWN"
		}
		samplesByCity[cityCode] = append(samplesByCity[cityCode], sample)
	}

	// Send separate alerts for each city
	for cityCode, citySamples := range samplesByCity {
		// Send samples in batches to avoid hitting Slack's 50 block limit
		maxSamplesPerMessage := commonConstants.MaxBlocksForSlackMessage

		for i := 0; i < len(citySamples); i += maxSamplesPerMessage {
			end := i + maxSamplesPerMessage
			if end > len(citySamples) {
				end = len(citySamples)
			}
			batch := citySamples[i:end]

			messageNumber := (i / maxSamplesPerMessage) + 1
			totalMessages := (len(citySamples) + maxSamplesPerMessage - 1) / maxSamplesPerMessage

			headerText := fmt.Sprintf(":alert: Samples Under Interlab Logistics Delay - %s", cityCode)
			if totalMessages > 1 {
				headerText = fmt.Sprintf(":alert: Samples Under Interlab Logistics Delay - %s (%d/%d)", cityCode, messageNumber, totalMessages)
			}

			blocks := []map[string]interface{}{
				{
					"type": "header",
					"text": map[string]interface{}{
						"type":  "plain_text",
						"text":  headerText,
						"emoji": true,
					},
				},
				{
					"type": "section",
					"fields": []map[string]interface{}{
						{"type": "mrkdwn", "text": "*Vial Type / Barcode / RL / RL Order ID*"},
						{"type": "mrkdwn", "text": "*Transferred from RL At*"},
					},
				},
			}

			for _, sample := range batch {
				vialType := vialTypeMap[sample.VialTypeId].VialType
				omsOrderIdString := commonUtils.GetStringOrderIdWithoutStringPart(sample.OmsOrderId)
				omsOrderUrl := fmt.Sprintf("%s/request/%s/order/%s",
					commonUtils.GetOmsBaseDomain(sample.OmsCityCode), sample.OmsRequestId, omsOrderIdString)
				sampleDetail := fmt.Sprintf("%s / %s / %s / %s", vialType, sample.Barcode, sample.OmsCityCode, sample.OmsOrderId)
				blocks = append(blocks, map[string]interface{}{
					"type": "section",
					"fields": []map[string]interface{}{
						{"type": "mrkdwn", "text": fmt.Sprintf("<%s|%s>", omsOrderUrl, sampleDetail)},
						{"type": "mrkdwn", "text": commonUtils.UtcToIst(*sample.TransferredAt).Format(commonConstants.PrettyDateTimeLayout)},
					},
				})
			}

			err := sampleService.SlackClient.SendToSlackDirectly(ctx, commonConstants.SlackSampleLossRiskAlertChannel, blocks)
			if err != nil {
				commonUtils.AddLog(ctx, commonConstants.ERROR_LEVEL, commonUtils.GetCurrentFunctionName(), map[string]interface{}{
					"error":    err.Error(),
					"cityCode": cityCode,
					"batch":    messageNumber,
				}, nil)
			}
		}
	}
}

func (sampleService *SampleService) GetSrfOrderIds(ctx context.Context, cityCode string) []string {
	omsOrderIds := []string{}
	cacheKey := commonConstants.CacheKeySrfOrderIdsAll
	if cityCode != "" {
		cacheKey = fmt.Sprintf(commonConstants.CacheKeySrfOrderIds, cityCode)
	}
	keyExists, _ := sampleService.Cache.Exists(ctx, cacheKey)
	if keyExists {
		sampleService.Cache.Get(ctx, cacheKey, &omsOrderIds)
		return omsOrderIds
	}
	omsOrderIds = sampleService.SampleDao.GetSrfOrderIds(cityCode)
	sampleService.Cache.Set(ctx, cacheKey, omsOrderIds, commonConstants.CacheExpiry5MinutesInt)
	return omsOrderIds
}
