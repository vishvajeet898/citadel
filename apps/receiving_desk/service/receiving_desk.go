package service

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/Orange-Health/citadel/apps/receiving_desk/constants"
	"github.com/Orange-Health/citadel/apps/receiving_desk/structures"
	commonConstants "github.com/Orange-Health/citadel/common/constants"
	commonStructures "github.com/Orange-Health/citadel/common/structures"
	commonUtils "github.com/Orange-Health/citadel/common/utils"
	"github.com/Orange-Health/citadel/models"
	commonModels "github.com/Orange-Health/citadel/models"
)

func validateCollectedSamplesRequest(
	collectedSamplesRequest structures.CollectedSamplesRequest) *commonStructures.CommonError {
	if collectedSamplesRequest.LabId == 0 {
		return &commonStructures.CommonError{
			StatusCode: http.StatusBadRequest,
			Message:    constants.ERROR_LAB_ID_REQUIRED,
		}
	}

	if !commonUtils.SliceContainsString(constants.SearchTypes, collectedSamplesRequest.SearchType) {
		return &commonStructures.CommonError{
			StatusCode: http.StatusBadRequest,
			Message:    constants.ERROR_INVALID_SEARCH_TYPE,
		}
	}

	if collectedSamplesRequest.SearchType == constants.SearchTypeOrderId && collectedSamplesRequest.OmsOrderId == "" {
		return &commonStructures.CommonError{
			StatusCode: http.StatusBadRequest,
			Message:    constants.ERROR_INVALID_ORDER_ID,
		}
	}

	if collectedSamplesRequest.SearchType == constants.SearchTypeBarcode && collectedSamplesRequest.Barcode == "" {
		return &commonStructures.CommonError{
			StatusCode: http.StatusBadRequest,
			Message:    commonConstants.ERROR_BARCODE_REQUIRED,
		}
	}

	if collectedSamplesRequest.SearchType == constants.SearchTypeTrfId && collectedSamplesRequest.TrfId == "" {
		return &commonStructures.CommonError{
			StatusCode: http.StatusBadRequest,
			Message:    constants.ERROR_TRF_ID_REQUIRED,
		}
	}

	return nil
}

func validateReceiveAndSyncSamplesRequest(
	receiveSamplesRequest structures.ReceiveSamplesRequest) *commonStructures.CommonError {
	if receiveSamplesRequest.ReceivingLabId == 0 {
		return &commonStructures.CommonError{
			Message:    constants.ERROR_LAB_ID_REQUIRED,
			StatusCode: http.StatusBadRequest,
		}
	}

	for _, sample := range receiveSamplesRequest.Samples {
		if sample.Id == 0 {
			return &commonStructures.CommonError{
				Message:    commonConstants.ERROR_SAMPLE_ID_REQUIRED,
				StatusCode: http.StatusBadRequest,
			}
		}

		if sample.Barcode == "" {
			return &commonStructures.CommonError{
				Message:    commonConstants.ERROR_BARCODE_REQUIRED,
				StatusCode: http.StatusBadRequest,
			}
		}
	}

	return nil
}

func getLabType(metaLabId, sessionLabId uint, labIdLabMap map[uint]commonStructures.Lab) string {
	if metaLabId == sessionLabId {
		return constants.LabTypeInhouse
	}

	if labIdLabMap[metaLabId].Inhouse {
		return constants.LabTypeInterlab
	}

	return constants.LabTypeOutsource
}

func createSampleNumberMaps(testDetails []commonModels.TestDetail, testSampleMappings []commonModels.TestSampleMapping,
	sessionLabId uint) (map[uint]uint, map[uint][]structures.TestDetailsRdResponse) {

	sampleNumberToNewLabIdMap := map[uint]uint{}
	omsTestIdToTestDetailsMap := map[string]commonModels.TestDetail{}

	for _, testDetail := range testDetails {
		omsTestIdToTestDetailsMap[testDetail.CentralOmsTestId] = testDetail
	}

	sampleNumberToTestsResponseMap := map[uint][]structures.TestDetailsRdResponse{}
	for _, testSampleMapping := range testSampleMappings {
		testDetail := omsTestIdToTestDetailsMap[testSampleMapping.OmsTestId]
		if sampleNumberToTestsResponseMap[testSampleMapping.SampleNumber] == nil {
			sampleNumberToTestsResponseMap[testSampleMapping.SampleNumber] = []structures.TestDetailsRdResponse{}
		}
		sampleNumberToTestsResponseMap[testSampleMapping.SampleNumber] = append(
			sampleNumberToTestsResponseMap[testSampleMapping.SampleNumber],
			structures.TestDetailsRdResponse{
				Id:              testSampleMapping.OmsTestId,
				Name:            testDetail.TestName,
				ProcessingLabId: testDetail.ProcessingLabId,
			})
		if _, ok := sampleNumberToNewLabIdMap[testSampleMapping.SampleNumber]; !ok {
			sampleNumberToNewLabIdMap[testSampleMapping.SampleNumber] = testDetail.ProcessingLabId
		} else {
			// This is to keep the new lab id as the inhouse lab id if any of the test is consucted to be inhouse
			// so that if the sample is being received, it should be synced first and then to any other state in future
			newLabId := sampleNumberToNewLabIdMap[testSampleMapping.SampleNumber]
			if newLabId != sessionLabId {
				sampleNumberToNewLabIdMap[testSampleMapping.SampleNumber] = testDetail.ProcessingLabId
			}
		}
	}

	return sampleNumberToNewLabIdMap, sampleNumberToTestsResponseMap
}

func createResponseForGetCollectedSamples(testDetails []commonModels.TestDetail,
	testSampleMappings []commonModels.TestSampleMapping, sessionLabId uint, samples []commonModels.Sample,
	sampleIdToSampleMetadataMap map[uint]commonModels.SampleMetadata, orderDetails commonModels.OrderDetails,
	patientDetails commonModels.PatientDetail,
	labIdLabMap map[uint]commonStructures.Lab) structures.CollectedSamplesResponse {

	response := structures.CollectedSamplesResponse{}

	patientAgeYears := uint(0)
	if patientDetails.Dob != nil {
		patientAgeYears, _, _ = commonUtils.GetAgeYearsMonthsAndDaysFromDob(*patientDetails.Dob)
	} else if patientDetails.ExpectedDob != nil {
		patientAgeYears, _, _ = commonUtils.GetAgeYearsMonthsAndDaysFromDob(*patientDetails.ExpectedDob)
	}

	_, sampleNumberToTestDetailsMap := createSampleNumberMaps(testDetails, testSampleMappings, sessionLabId)

	for _, sample := range samples {
		sampleResponse := structures.CollectedSample{
			Id:                 sample.Id,
			ParentSampleId:     sample.ParentSampleId,
			OrderId:            orderDetails.OmsOrderId,
			RequestId:          orderDetails.OmsRequestId,
			ServicingCityCode:  orderDetails.CityCode,
			Barcode:            sample.Barcode,
			TrfId:              orderDetails.TrfId,
			VialType:           sample.VialTypeId,
			SampleNumber:       sample.SampleNumber,
			Status:             sample.Status,
			CollectLaterReason: sampleIdToSampleMetadataMap[sample.Id].CollectLaterReason,
			NotReceivedReason:  sampleIdToSampleMetadataMap[sample.Id].NotReceivedReason,
			PatientDetails: structures.PatientDetailsRdResponse{
				Name:        patientDetails.Name,
				Dob:         patientDetails.Dob,
				ExpectedDob: patientDetails.ExpectedDob,
				Age:         patientAgeYears,
				Gender:      patientDetails.Gender,
			},
		}
		testDetails := sampleNumberToTestDetailsMap[sample.SampleNumber]
		if len(testDetails) == 0 {
			continue
		}

		processingLabIdToTestDetailsMap := map[uint][]structures.TestDetailsRdResponse{}
		for _, testDetail := range testDetails {
			if processingLabIdToTestDetailsMap[testDetail.ProcessingLabId] == nil {
				processingLabIdToTestDetailsMap[testDetail.ProcessingLabId] = []structures.TestDetailsRdResponse{}
			}
			processingLabIdToTestDetailsMap[testDetail.ProcessingLabId] = append(
				processingLabIdToTestDetailsMap[testDetail.ProcessingLabId], testDetail)
		}
		for processingLabId, testDetails := range processingLabIdToTestDetailsMap {
			sampleResponse.Tests = testDetails
			sampleResponse.LabDetails = structures.LabDetailsRdResponse{
				Id:      processingLabId,
				Name:    labIdLabMap[processingLabId].LabName,
				LabType: getLabType(processingLabId, sessionLabId, labIdLabMap),
			}
			response.CollectedSamples = append(response.CollectedSamples, sampleResponse)
		}
	}

	return response
}

func (rdService *ReceivingDeskService) getOrderDetailsForReceivingDeskResponse(ctx context.Context,
	collectedSamplesRequest structures.CollectedSamplesRequest) (
	commonModels.OrderDetails, *commonStructures.CommonError) {

	orderDetails := commonModels.OrderDetails{}
	var cErr *commonStructures.CommonError

	switch collectedSamplesRequest.SearchType {
	case constants.SearchTypeOrderId:
		orderDetails, cErr = rdService.OrderDetailsService.GetOrderDetailsByOmsOrderId(collectedSamplesRequest.OmsOrderId)
		if cErr != nil {
			return orderDetails, cErr
		}
		if orderDetails.Id == 0 {
			return orderDetails, &commonStructures.CommonError{
				Message:    commonConstants.ERROR_ORDER_ID_NOT_FOUND,
				StatusCode: http.StatusNotFound,
			}
		}
	case constants.SearchTypeBarcode:
		var sample commonModels.Sample
		sample, cErr = rdService.SampleService.GetSampleByBarcodeForReceiving(collectedSamplesRequest.Barcode)
		if cErr != nil {
			return orderDetails, cErr
		}
		if sample.Id == 0 {
			return orderDetails, &commonStructures.CommonError{
				Message:    commonConstants.ERROR_BARCODE_NOT_FOUND,
				StatusCode: http.StatusNotFound,
			}
		}
		inhouseLabIds := rdService.CdsService.GetInhouseLabIds(ctx)
		if sample.DestinationLabId != collectedSamplesRequest.LabId &&
			commonUtils.SliceContainsUint(inhouseLabIds, sample.DestinationLabId) {
			destinationLab, _ := rdService.CdsService.GetLabById(ctx, sample.DestinationLabId)
			return commonModels.OrderDetails{}, &commonStructures.CommonError{
				Message:    fmt.Sprintf(constants.ERROR_SCANNED_AT_INCORRECT_LAB, strings.ToLower(destinationLab.LabName)),
				StatusCode: http.StatusBadRequest,
			}
		}
		orderDetails, cErr = rdService.OrderDetailsService.GetOrderDetailsByOmsOrderId(sample.OmsOrderId)
		if cErr != nil {
			return orderDetails, cErr
		}
	case constants.SearchTypeTrfId:
		orderDetails, cErr = rdService.OrderDetailsService.GetOrderDetailsByTrfId(collectedSamplesRequest.TrfId)
		if cErr != nil {
			return orderDetails, cErr
		}
		if orderDetails.Id == 0 {
			return orderDetails, &commonStructures.CommonError{
				Message:    commonConstants.ERROR_TRF_ID_NOT_FOUND,
				StatusCode: http.StatusNotFound,
			}
		}
	}

	return orderDetails, nil
}

func (rdService *ReceivingDeskService) getTaskAndTaskMetadataDtos(orderDetails commonModels.OrderDetails) (
	commonModels.Task, commonModels.TaskMetadata) {
	task := commonModels.Task{
		OrderId:          commonUtils.GetUintOrderIdWithoutStringPart(orderDetails.OmsOrderId),
		RequestId:        commonUtils.GetUintRequestIdWithoutStringPart(orderDetails.OmsRequestId),
		OmsOrderId:       orderDetails.OmsOrderId,
		OmsRequestId:     orderDetails.OmsRequestId,
		LabId:            orderDetails.ServicingLabId,
		CityCode:         orderDetails.CityCode,
		Status:           commonConstants.TASK_STATUS_PENDING,
		PreviousStatus:   commonConstants.TASK_STATUS_PENDING,
		OrderType:        commonConstants.CollecTypeToOrderTypeMap[orderDetails.CollectionType],
		PatientDetailsId: orderDetails.PatientDetailsId,
		IsActive:         true,
	}

	taskMetadata := commonModels.TaskMetadata{}

	return task, taskMetadata
}

func (rdService *ReceivingDeskService) GetCollectedSamples(ctx context.Context,
	collectedSamplesRequest structures.CollectedSamplesRequest) (structures.CollectedSamplesResponse,
	*commonStructures.CommonError) {

	response := structures.CollectedSamplesResponse{}
	cErr := validateCollectedSamplesRequest(collectedSamplesRequest)
	if cErr != nil {
		return response, cErr
	}

	orderDetails, cErr := rdService.getOrderDetailsForReceivingDeskResponse(ctx, collectedSamplesRequest)
	if cErr != nil {
		return response, cErr
	}

	samples, sampleMetadatas, cErr := rdService.SampleService.GetCollectedSamples(orderDetails.OmsOrderId,
		collectedSamplesRequest.LabId)
	if cErr != nil {
		return response, cErr
	}

	if len(samples) == 0 && orderDetails.TrfId != "" {
		return response, &commonStructures.CommonError{
			Message:    constants.ERROR_DIGITISATION_ISSUE,
			StatusCode: http.StatusNotFound,
		}
	}

	if len(samples) == 0 {
		return response, nil
	}

	sampleIdToSampleMetadataMap := map[uint]commonModels.SampleMetadata{}
	for _, sampleMetadata := range sampleMetadatas {
		sampleIdToSampleMetadataMap[sampleMetadata.SampleId] = sampleMetadata
	}

	patientDetails, testDetails := commonModels.PatientDetail{}, []commonModels.TestDetail{}
	testSampleMappings, labIdLabMap := []commonModels.TestSampleMapping{}, map[uint]commonStructures.Lab{}
	errList, wg, mu := []*commonStructures.CommonError{}, sync.WaitGroup{}, sync.Mutex{}

	sampleNumbers := []uint{}
	for _, sample := range samples {
		sampleNumbers = append(sampleNumbers, sample.SampleNumber)
	}

	wg.Add(3)

	// Fetch patientDetails in parallel
	go func() {
		defer wg.Done()
		details, err := rdService.PatientDetailService.GetPatientDetailsById(orderDetails.PatientDetailsId)
		mu.Lock()
		defer mu.Unlock()
		if err != nil {
			errList = append(errList, err)
		} else {
			patientDetails = details
		}
	}()

	// Fetch testDetails and testSampleMappings in parallel
	go func() {
		defer wg.Done()
		td, tsm, err := rdService.SampleService.GetAllTestsAndSampleMappingsBySampleNumbers(sampleNumbers,
			orderDetails.OmsOrderId)
		mu.Lock()
		defer mu.Unlock()
		if err != nil {
			errList = append(errList, err)
		} else {
			testDetails = td
			testSampleMappings = tsm
		}
	}()

	// Fetch labIdLabMap in parallel
	go func() {
		defer wg.Done()
		labMap := rdService.CdsService.GetLabIdLabMap(ctx)
		mu.Lock()
		labIdLabMap = labMap
		mu.Unlock()
	}()

	wg.Wait()

	if len(errList) > 0 {
		// Return the first error encountered
		return response, errList[0]
	}

	return createResponseForGetCollectedSamples(testDetails, testSampleMappings, collectedSamplesRequest.LabId,
		samples, sampleIdToSampleMetadataMap, orderDetails, patientDetails, labIdLabMap), nil
}

func (rdService *ReceivingDeskService) updateLabEtaAndLabTat(omsOrderId string, omsTestIds []string, lisSyncAt *time.Time,
	labId uint, cityCode string) {

	ctx := context.Background()
	if len(omsTestIds) == 0 || lisSyncAt == nil {
		messageBody, messageAttributes := rdService.PubsubService.GetLabEtaUpdateEvent(omsOrderId, omsTestIds, lisSyncAt,
			cityCode)
		rdService.SnsClient.PublishTo(ctx, messageBody, messageAttributes, "", commonConstants.OmsUpdatesTopicArn, "")
		return
	}

	testDetails, cErr := rdService.TestDetailsService.GetTestDetailModelByOmsTestIds(omsTestIds)
	if cErr != nil {
		commonUtils.AddLog(ctx, commonConstants.ERROR_LEVEL, commonUtils.GetCurrentFunctionName(), nil,
			errors.New(cErr.Message))
	}

	masterTestIds, masterTestIdToTestDetailsMap := []uint{}, map[uint][]commonModels.TestDetail{}
	for _, testDetail := range testDetails {
		masterTestIds = append(masterTestIds, testDetail.MasterTestId)
		if _, ok := masterTestIdToTestDetailsMap[testDetail.MasterTestId]; !ok {
			masterTestIdToTestDetailsMap[testDetail.MasterTestId] = []commonModels.TestDetail{}
		}
		masterTestIdToTestDetailsMap[testDetail.MasterTestId] = append(masterTestIdToTestDetailsMap[testDetail.MasterTestId],
			testDetail)
	}

	masterTests := rdService.CdsService.GetMasterTestsByIds(ctx, masterTestIds)
	masterTestIdToMasterTestMap := map[uint]commonStructures.CdsTestMaster{}
	for _, masterTest := range masterTests {
		masterTestIdToMasterTestMap[masterTest.Id] = masterTest
	}

	for _, masterTest := range masterTests {
		if _, ok := masterTest.TestLabMeta[labId]; !ok {
			continue
		}
		labTat := masterTest.TestLabMeta[labId].LabTat
		labEta := commonUtils.CalculateLabEta(lisSyncAt, labTat)
		testDetails := masterTestIdToTestDetailsMap[masterTest.Id]
		for index := range testDetails {
			testDetails[index].LabEta = labEta
			testDetails[index].LabTat = labTat
			testDetails[index].UpdatedAt = commonUtils.GetCurrentTime()
			testDetails[index].UpdatedBy = commonConstants.CitadelSystemId
		}
		masterTestIdToTestDetailsMap[masterTest.Id] = testDetails
	}

	finalTestDetails := []commonModels.TestDetail{}
	for _, testDetails := range masterTestIdToTestDetailsMap {
		finalTestDetails = append(finalTestDetails, testDetails...)
	}

	_, cErr = rdService.TestDetailsService.UpdateTestDetails(finalTestDetails)
	if cErr != nil {
		commonUtils.AddLog(ctx, commonConstants.ERROR_LEVEL, commonUtils.GetCurrentFunctionName(), nil,
			errors.New(cErr.Message))
	}

	messageBody, messageAttributes := rdService.PubsubService.GetLabEtaUpdateEvent(omsOrderId, omsTestIds, lisSyncAt,
		cityCode)
	rdService.SnsClient.PublishTo(ctx, messageBody, messageAttributes, "", commonConstants.OmsUpdatesTopicArn, "")
}

func (rdService *ReceivingDeskService) ReceiveAndSyncSamplesByOmsOrderId(ctx context.Context,
	omsOrderId string, sessionLabId, userId uint, samples []models.Sample,
	samplesMetadata []models.SampleMetadata) *commonStructures.CommonError {

	taskId, inhouseSampleIds, inhouseSampleNumbers, sampleNumbers := uint(0), []uint{}, []uint{}, []uint{}
	inhouseSamples, outsourceSamples, interlabSamples :=
		[]commonModels.Sample{}, []commonModels.Sample{}, []commonModels.Sample{}
	inhouseSamplesMetadata, outsourceSamplesMetadata, interlabSamplesMetadata :=
		[]commonModels.SampleMetadata{}, []commonModels.SampleMetadata{}, []commonModels.SampleMetadata{}
	task, taskMetadata, taskExists := commonModels.Task{}, commonModels.TaskMetadata{}, true
	sampleIdBarcodeMap := map[uint]string{}
	interlabTestDetails, omsTestIdsToTestDetailsMap, masterTestIdToMasterTestMap, interlabTestDetailsIds :=
		[]commonModels.TestDetail{}, map[string]commonModels.TestDetail{}, map[uint]commonStructures.CdsTestMaster{}, []uint{}

	sampleIdToSampleMetadataMap := map[uint]commonModels.SampleMetadata{}
	for _, sampleMetadata := range samplesMetadata {
		sampleIdToSampleMetadataMap[sampleMetadata.SampleId] = sampleMetadata
	}

	for _, sample := range samples {
		sampleNumbers = append(sampleNumbers, sample.SampleNumber)
		sampleIdBarcodeMap[sample.Id] = sample.Barcode
	}

	redisKey := fmt.Sprintf(commonConstants.CacheKeyReceiveAndSync, omsOrderId)
	keyExists, err := rdService.Cache.Exists(ctx, redisKey)
	if err != nil || keyExists {
		commonUtils.AddLog(ctx, commonConstants.DEBUG_LEVEL, commonUtils.GetCurrentFunctionName(), nil, err)
		return &commonStructures.CommonError{
			Message:    commonConstants.ERROR_RECEIVE_AND_SYNC_IN_PROGRESS,
			StatusCode: http.StatusBadRequest,
		}
	}

	err = rdService.Cache.Set(ctx, redisKey, true, commonConstants.CacheExpiry1HourInt)
	if err != nil {
		commonUtils.AddLog(ctx, commonConstants.DEBUG_LEVEL, commonUtils.GetCurrentFunctionName(), nil, err)
		return &commonStructures.CommonError{
			Message:    err.Error(),
			StatusCode: http.StatusInternalServerError,
		}
	}

	defer func() {
		err := rdService.Cache.Delete(ctx, redisKey)
		if err != nil {
			commonUtils.AddLog(ctx, commonConstants.DEBUG_LEVEL, commonUtils.GetCurrentFunctionName(), nil, err)
		}
	}()

	orderDetails, cErr := rdService.OrderDetailsService.GetOrderDetailsByOmsOrderId(omsOrderId)
	if cErr != nil {
		return cErr
	}

	testDetails, testSampleMappings, cErr := rdService.SampleService.GetAllTestsAndSampleMappingsBySampleNumbers(
		sampleNumbers, omsOrderId)
	if cErr != nil {
		return cErr
	}

	masterTestIds := []uint{}
	for _, testDetail := range testDetails {
		masterTestIds = append(masterTestIds, testDetail.MasterTestId)
		omsTestIdsToTestDetailsMap[testDetail.CentralOmsTestId] = testDetail
	}

	masterTests := rdService.CdsService.GetMasterTestsByIds(ctx, masterTestIds)
	for _, masterTest := range masterTests {
		masterTestIdToMasterTestMap[masterTest.Id] = masterTest
	}

	labIdLabMap := rdService.CdsService.GetLabIdLabMap(ctx)

	sampleNumberToNewLabIdMap, sampleNumberToTestsResponseMap := createSampleNumberMaps(testDetails, testSampleMappings,
		sessionLabId)

	for _, sample := range samples {
		switch getLabType(sampleNumberToNewLabIdMap[sample.SampleNumber], sessionLabId, labIdLabMap) {
		case constants.LabTypeInhouse:
			inhouseSampleIds = append(inhouseSampleIds, sample.Id)
			inhouseSamples = append(inhouseSamples, sample)
			inhouseSampleNumbers = append(inhouseSampleNumbers, sample.SampleNumber)
		case constants.LabTypeInterlab:
			interlabSamples = append(interlabSamples, sample)
			for _, testDetailResponse := range sampleNumberToTestsResponseMap[sample.SampleNumber] {
				interlabTestDetail := omsTestIdsToTestDetailsMap[testDetailResponse.Id]
				if commonUtils.SliceContainsUint(interlabTestDetailsIds, interlabTestDetail.Id) {
					continue
				}
				interlabTestDetailsIds = append(interlabTestDetailsIds, interlabTestDetail.Id)
				interlabTestDetails = append(interlabTestDetails, interlabTestDetail)
			}
		case constants.LabTypeOutsource:
			outsourceSamples = append(outsourceSamples, sample)
		}
	}

	currentTime := commonUtils.GetCurrentTime()
	// Step 1: Update Sample Statuses to SampleReceived
	if len(inhouseSampleIds) > 0 {
		for index := range inhouseSamples {
			inhouseSamples[index].UpdatedBy = userId
			inhouseSamples[index].Barcode = sampleIdBarcodeMap[inhouseSamples[index].Id]
			sampleMetadata := sampleIdToSampleMetadataMap[inhouseSamples[index].Id]
			if inhouseSamples[index].Status != commonConstants.SampleReceived {
				sampleMetadata.ReceivedAt = currentTime
			}
			sampleMetadata.UpdatedBy = userId
			inhouseSamples[index].Status = commonConstants.SampleReceived
			inhouseSamplesMetadata = append(inhouseSamplesMetadata, sampleMetadata)
		}
	}

	for index := range outsourceSamples {
		outsourceSamples[index].UpdatedBy = userId
		outsourceSamples[index].Barcode = sampleIdBarcodeMap[outsourceSamples[index].Id]
		sampleMetadata := sampleIdToSampleMetadataMap[outsourceSamples[index].Id]
		if outsourceSamples[index].Status != commonConstants.SampleReceived {
			sampleMetadata.ReceivedAt = currentTime
		}
		outsourceSamples[index].Status = commonConstants.SampleReceived
		sampleMetadata.UpdatedBy = userId
		outsourceSamplesMetadata = append(outsourceSamplesMetadata, sampleMetadata)
	}

	for index := range interlabSamples {
		interlabSamples[index].UpdatedBy = userId
		interlabSamples[index].Barcode = sampleIdBarcodeMap[interlabSamples[index].Id]
		sampleMetadata := sampleIdToSampleMetadataMap[interlabSamples[index].Id]
		if interlabSamples[index].Status != commonConstants.SampleReceived {
			sampleMetadata.ReceivedAt = currentTime
		}
		interlabSamples[index].Status = commonConstants.SampleReceived
		sampleMetadata.UpdatedBy = userId
		interlabSamplesMetadata = append(interlabSamplesMetadata, sampleMetadata)
	}

	updatedSamples, updatedSamplesMetadata := []commonModels.Sample{}, []commonModels.SampleMetadata{}
	if len(inhouseSamples) > 0 {
		updatedSamples = append(updatedSamples, inhouseSamples...)
		updatedSamplesMetadata = append(updatedSamplesMetadata, inhouseSamplesMetadata...)
	}

	if len(outsourceSamples) > 0 {
		updatedSamples = append(updatedSamples, outsourceSamples...)
		updatedSamplesMetadata = append(updatedSamplesMetadata, outsourceSamplesMetadata...)
	}

	if len(interlabSamples) > 0 {
		updatedSamples = append(updatedSamples, interlabSamples...)
		updatedSamplesMetadata = append(updatedSamplesMetadata, interlabSamplesMetadata...)
	}

	tx := rdService.ReceivingDeskDao.BeginTransaction()
	defer tx.Rollback()

	_, _, cErr = rdService.SampleService.UpdateSamplesAndSamplesMetadataWithTx(tx, updatedSamples, updatedSamplesMetadata)
	if cErr != nil {
		return cErr
	}

	tx.Commit()

	sampleIdToSampleMetadataMap = map[uint]commonModels.SampleMetadata{}
	for _, sampleMetadata := range updatedSamplesMetadata {
		sampleIdToSampleMetadataMap[sampleMetadata.SampleId] = sampleMetadata
	}

	outsourceSamplesMetadata, interlabSamplesMetadata = []commonModels.SampleMetadata{}, []commonModels.SampleMetadata{}

	// Step 2: Update the samples statuses to final statuses
	currentTime = commonUtils.GetCurrentTime()
	if len(inhouseSampleIds) > 0 {
		visitId, cErr := rdService.AttuneService.SyncDataToLisByOmsOrderId(ctx, omsOrderId, sessionLabId, inhouseSamples,
			inhouseSamplesMetadata, sampleIdBarcodeMap)
		if cErr != nil {
			return cErr
		}
		inhouseSamplesMetadata = []commonModels.SampleMetadata{}

		for index := range inhouseSamples {
			inhouseSamples[index].VisitId = visitId
			inhouseSamples[index].UpdatedBy = userId
			inhouseSamples[index].Barcode = sampleIdBarcodeMap[inhouseSamples[index].Id]
			inhouseSamples[index].LabId = inhouseSamples[index].DestinationLabId
			inhouseSamples[index].Status = commonConstants.SampleSynced
			sampleMetadata := sampleIdToSampleMetadataMap[inhouseSamples[index].Id]
			sampleMetadata.UpdatedBy = userId
			sampleMetadata.LisSyncAt = currentTime
			inhouseSamplesMetadata = append(inhouseSamplesMetadata, sampleMetadata)
		}

		taskFromDb, cErr := rdService.TaskService.GetTaskByOmsOrderId(omsOrderId)
		if cErr != nil {
			return cErr
		}
		if taskFromDb.Id != 0 {
			taskId = taskFromDb.Id
		}
		if taskFromDb.Id == 0 {
			taskExists = false
			task, taskMetadata = rdService.getTaskAndTaskMetadataDtos(orderDetails)
		}
	}

	for index := range outsourceSamples {
		outsourceSamples[index].UpdatedBy = userId
		outsourceSamples[index].Barcode = sampleIdBarcodeMap[outsourceSamples[index].Id]
		outsourceSamples[index].LabId = outsourceSamples[index].DestinationLabId
		outsourceSamples[index].DestinationLabId = sampleNumberToNewLabIdMap[outsourceSamples[index].SampleNumber]
		outsourceSamples[index].Status = commonConstants.SampleOutsourced
		sampleMetadata := sampleIdToSampleMetadataMap[outsourceSamples[index].Id]
		sampleMetadata.OutsourcedAt = currentTime
		sampleMetadata.UpdatedBy = userId
		outsourceSamplesMetadata = append(outsourceSamplesMetadata, sampleMetadata)
	}

	for index := range interlabSamples {
		interlabSamples[index].UpdatedBy = userId
		interlabSamples[index].Barcode = sampleIdBarcodeMap[interlabSamples[index].Id]
		interlabSamples[index].LabId = interlabSamples[index].DestinationLabId
		interlabSamples[index].Status = commonConstants.SampleTransferred
		sampleMetadata := sampleIdToSampleMetadataMap[interlabSamples[index].Id]
		sampleMetadata.TransferredAt = currentTime
		sampleMetadata.UpdatedBy = userId
		interlabSamplesMetadata = append(interlabSamplesMetadata, sampleMetadata)
	}

	for index := range interlabTestDetails {
		currentProcessingLab := interlabTestDetails[index].ProcessingLabId
		masterTestMeta := masterTestIdToMasterTestMap[interlabTestDetails[index].MasterTestId]
		nextProcessingLabId := masterTestMeta.TestLabMeta[currentProcessingLab].LabId
		interlabTestDetails[index].ProcessingLabId = nextProcessingLabId
		interlabTestDetails[index].UpdatedBy = commonConstants.CitadelSystemId
	}

	updatedSamples, updatedSamplesMetadata = []commonModels.Sample{}, []commonModels.SampleMetadata{}
	if len(inhouseSamples) > 0 {
		updatedSamples = append(updatedSamples, inhouseSamples...)
		updatedSamplesMetadata = append(updatedSamplesMetadata, inhouseSamplesMetadata...)
	}

	if len(outsourceSamples) > 0 {
		updatedSamples = append(updatedSamples, outsourceSamples...)
		updatedSamplesMetadata = append(updatedSamplesMetadata, outsourceSamplesMetadata...)
	}

	if len(interlabSamples) > 0 {
		updatedSamples = append(updatedSamples, interlabSamples...)
		updatedSamplesMetadata = append(updatedSamplesMetadata, interlabSamplesMetadata...)
	}

	tx = rdService.ReceivingDeskDao.BeginTransaction()
	defer tx.Rollback()

	_, _, cErr = rdService.SampleService.UpdateSamplesAndSamplesMetadataWithTx(tx, updatedSamples, updatedSamplesMetadata)
	if cErr != nil {
		return cErr
	}

	if len(interlabSamples) > 0 {
		if cErr := rdService.SampleService.CreateInterlabSamplesWithTx(ctx, tx, interlabSamples,
			interlabSamplesMetadata, sampleNumberToNewLabIdMap); cErr != nil {
			return cErr
		}
	}

	if len(inhouseSampleIds) > 0 && !taskExists {
		task, cErr = rdService.TaskService.CreateTaskWithTx(tx, task)
		if cErr != nil {
			return cErr
		}

		taskMetadata.TaskId = task.Id
		taskId = task.Id
		_, cErr = rdService.TaskService.CreateTaskMetadataWithTx(tx, taskMetadata)
		if cErr != nil {
			return cErr
		}
	}

	inhouseOmsTestIds := []string{}
	for index := range inhouseSampleNumbers {
		for _, testDetail := range sampleNumberToTestsResponseMap[inhouseSampleNumbers[index]] {
			if testDetail.ProcessingLabId == sessionLabId {
				inhouseOmsTestIds = append(inhouseOmsTestIds, testDetail.Id)
			}
		}
	}

	cErr = rdService.TestDetailsService.UpdateTaskIdInTestDetailsWithOmsTestIdWithTx(tx, inhouseOmsTestIds, taskId)
	if cErr != nil {
		return cErr
	}

	_, cErr = rdService.TestDetailsService.UpdateTestDetailsWithTx(tx, interlabTestDetails)
	if cErr != nil {
		return cErr
	}

	tx.Commit()

	lisSyncAt := &time.Time{}
	if len(inhouseOmsTestIds) > 0 {
		lisSyncAt = inhouseSamplesMetadata[0].LisSyncAt
	}
	go rdService.updateLabEtaAndLabTat(omsOrderId, inhouseOmsTestIds, lisSyncAt, sessionLabId, orderDetails.CityCode)

	return nil
}

func (rdService *ReceivingDeskService) ReceiveAndSyncSamples(ctx context.Context,
	receiveSamplesRequest structures.ReceiveSamplesRequest) ([]string, *commonStructures.CommonError) {

	cErr := validateReceiveAndSyncSamplesRequest(receiveSamplesRequest)
	if cErr != nil {
		return nil, cErr
	}

	sessionLabId, userId := receiveSamplesRequest.ReceivingLabId, receiveSamplesRequest.UserId
	newBarcodes, sampleIdBarcodeMap := []string{}, map[uint]string{}
	errorStrings := []string{}

	sampleIds := []uint{}
	for _, sample := range receiveSamplesRequest.Samples {
		sampleIds = append(sampleIds, sample.Id)
		sampleIdBarcodeMap[sample.Id] = sample.Barcode
	}

	samples, samplesMetadata, cErr := rdService.SampleService.GetSamplesDataBySampleIds(sampleIds)
	if cErr != nil {
		return nil, cErr
	}

	sampleIdToSampleMetadataMap := map[uint]commonModels.SampleMetadata{}
	for _, sampleMetadata := range samplesMetadata {
		sampleIdToSampleMetadataMap[sampleMetadata.SampleId] = sampleMetadata
	}

	if len(samples) == 0 {
		return nil, &commonStructures.CommonError{
			Message:    commonConstants.ERROR_NO_SAMPLES_FOUND,
			StatusCode: http.StatusNotFound,
		}
	}

	for _, sample := range samples {
		if sample.Barcode == "" {
			newBarcodes = append(newBarcodes, sampleIdBarcodeMap[sample.Id])
		} else if sample.Barcode != sampleIdBarcodeMap[sample.Id] {
			newBarcodes = append(newBarcodes, sampleIdBarcodeMap[sample.Id])
		}
	}

	barcodesExists, cErr := rdService.SampleService.BarcodesExistsInSystem(newBarcodes)
	if cErr != nil {
		return nil, cErr
	}
	if barcodesExists {
		return nil, &commonStructures.CommonError{
			Message:    commonConstants.ERROR_DUPLICATE_BARCODE,
			StatusCode: http.StatusBadRequest,
		}
	}

	omsOrderIds, successfulOmsOrderIds := []string{}, []string{}
	omsOrderIdSampleMap, omsOrderIdSampleMetadataMap := map[string][]models.Sample{}, map[string][]models.SampleMetadata{}
	for _, sample := range samples {
		if sample.OmsOrderId == "" {
			continue
		}
		sample.Barcode = sampleIdBarcodeMap[sample.Id]
		omsOrderIds = append(omsOrderIds, sample.OmsOrderId)
		if _, ok := omsOrderIdSampleMap[sample.OmsOrderId]; !ok {
			omsOrderIdSampleMap[sample.OmsOrderId] = []models.Sample{}
		}
		omsOrderIdSampleMap[sample.OmsOrderId] =
			append(omsOrderIdSampleMap[sample.OmsOrderId], sample)
	}
	for _, sampleMetadata := range samplesMetadata {
		if _, ok := omsOrderIdSampleMetadataMap[sampleMetadata.OmsOrderId]; !ok {
			omsOrderIdSampleMetadataMap[sampleMetadata.OmsOrderId] = []models.SampleMetadata{}
		}
		omsOrderIdSampleMetadataMap[sampleMetadata.OmsOrderId] =
			append(omsOrderIdSampleMetadataMap[sampleMetadata.OmsOrderId], sampleMetadata)
	}

	omsOrderIds = commonUtils.CreateUniqueSliceString(omsOrderIds)
	if len(omsOrderIds) == 0 {
		return nil, &commonStructures.CommonError{
			Message:    commonConstants.ERROR_NO_SAMPLES_FOUND,
			StatusCode: http.StatusNotFound,
		}
	}

	// Any Vial Any Order Disabled based on config
	if len(omsOrderIds) > 1 && !commonConstants.MultipleOrdersReceivingEnabled {
		return nil, &commonStructures.CommonError{
			Message:    commonConstants.ERROR_MULTIPLE_ORDERS_NOT_ALLOWED,
			StatusCode: http.StatusBadRequest,
		}
	}

	wg, mu := sync.WaitGroup{}, sync.Mutex{}
	wg.Add(len(omsOrderIds))

	for _, omsOrderId := range omsOrderIds {
		go func(samples []models.Sample, sampleMetadatas []models.SampleMetadata) {
			defer wg.Done()
			goroutineErr := rdService.ReceiveAndSyncSamplesByOmsOrderId(ctx, omsOrderId, sessionLabId, userId,
				samples, sampleMetadatas)
			if goroutineErr != nil {
				errString := fmt.Sprintf("%s %s, err: %v\n", commonConstants.ERROR_IN_RECEIVING_FOR_THIS_ORDER_ID,
					omsOrderId, goroutineErr.Message)
				mu.Lock()
				errorStrings = append(errorStrings, errString)
				mu.Unlock()
			} else {
				mu.Lock()
				successfulOmsOrderIds = append(successfulOmsOrderIds, omsOrderId)
				mu.Unlock()
			}
		}(omsOrderIdSampleMap[omsOrderId], omsOrderIdSampleMetadataMap[omsOrderId])
	}

	wg.Wait()

	if len(errorStrings) > 0 {
		return successfulOmsOrderIds, &commonStructures.CommonError{
			Message:    strings.Join(errorStrings, "\n"),
			StatusCode: http.StatusInternalServerError,
		}
	}

	return successfulOmsOrderIds, nil
}

func (rdService *ReceivingDeskService) MarkAsNotReceived(ctx context.Context,
	markAsNotReceivedRequest commonStructures.MarkAsNotReceivedRequest, sendEvents bool) *commonStructures.CommonError {

	if markAsNotReceivedRequest.SampleId == 0 {
		return &commonStructures.CommonError{
			Message:    commonConstants.ERROR_SAMPLE_ID_REQUIRED,
			StatusCode: http.StatusBadRequest,
		}
	}

	if markAsNotReceivedRequest.NotReceivedReason == "" {
		return &commonStructures.CommonError{
			Message:    constants.ERROR_NOT_RECEIVED_REASON_REQUIRED,
			StatusCode: http.StatusBadRequest,
		}
	}

	omsTestIds, omsTestStatusMap := []string{}, map[string]string{}

	// Fetch Sample Data
	sample, sampleMetadata, cErr := rdService.SampleService.GetSampleDataBySampleId(markAsNotReceivedRequest.SampleId)
	if cErr != nil {
		return cErr
	}
	// Update Sample Status
	sample.Status = commonConstants.SampleNotReceived
	sample.UpdatedBy = markAsNotReceivedRequest.UserId
	sampleMetadata.UpdatedBy = markAsNotReceivedRequest.UserId
	sampleMetadata.NotReceivedAt = commonUtils.GetCurrentTime()
	sampleMetadata.NotReceivedReason = markAsNotReceivedRequest.NotReceivedReason

	// Update Test Status
	testDetails, cErr := rdService.SampleService.GetAllSampleTestsBySampleNumber(sample.SampleNumber, sample.OmsOrderId)
	if cErr != nil {
		return cErr
	}
	for index := range testDetails {
		testDetails[index].Status = commonConstants.TEST_STATUS_SAMPLE_NOT_RECEIVED
		testDetails[index].UpdatedBy = markAsNotReceivedRequest.UserId
		omsTestIds = append(omsTestIds, testDetails[index].CentralOmsTestId)
		omsTestStatusMap[testDetails[index].CentralOmsTestId] = commonConstants.TEST_STATUS_SAMPLE_NOT_RECEIVED
	}

	parentSample, parentSampleMetadata := commonModels.Sample{}, commonModels.SampleMetadata{}
	if sample.ParentSampleId != 0 {
		parentSample, parentSampleMetadata, cErr = rdService.SampleService.GetSampleDataBySampleId(sample.ParentSampleId)
		if cErr != nil {
			return cErr
		}
		parentSample.Status = commonConstants.SampleNotReceived
		parentSample.UpdatedBy = markAsNotReceivedRequest.UserId
		parentSampleMetadata.UpdatedBy = markAsNotReceivedRequest.UserId
		parentSampleMetadata.NotReceivedAt = commonUtils.GetCurrentTime()
		parentSampleMetadata.NotReceivedReason = markAsNotReceivedRequest.NotReceivedReason
	}

	samples, samplesMetadata := []commonModels.Sample{sample}, []commonModels.SampleMetadata{sampleMetadata}
	if parentSample.Id != 0 {
		samples = append(samples, parentSample)
		samplesMetadata = append(samplesMetadata, parentSampleMetadata)
	}

	tx := rdService.ReceivingDeskDao.BeginTransaction()
	defer tx.Rollback()

	_, _, cErr = rdService.SampleService.UpdateSamplesAndSamplesMetadataWithTx(tx, samples, samplesMetadata)
	if cErr != nil {
		return cErr
	}

	_, cErr = rdService.TestDetailsService.UpdateTestDetailsWithTx(tx, testDetails)
	if cErr != nil {
		return cErr
	}

	tx.Commit()

	if len(omsTestIds) > 0 && sendEvents {
		go rdService.SampleService.PublishResetTatsEvent(omsTestIds, sample.OmsCityCode)
		go rdService.SampleService.PublishUpdateTestStatusEvent(omsTestStatusMap, sample.OmsOrderId, true, sample.OmsCityCode)
	}

	return nil
}
