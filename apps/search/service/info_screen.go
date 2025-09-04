package service

import (
	"context"
	"fmt"
	"net/http"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/Orange-Health/citadel/apps/search/constants"
	"github.com/Orange-Health/citadel/apps/search/structures"
	commonConstants "github.com/Orange-Health/citadel/common/constants"
	commonStructures "github.com/Orange-Health/citadel/common/structures"
	commonUtils "github.com/Orange-Health/citadel/common/utils"
)

func getTestStatus(status, omsStatus string) string {
	switch omsStatus {
	case commonConstants.TEST_STATUS_COMPLETED_NOT_SENT:
		return commonConstants.TEST_STATUS_COMPLETED_NOT_SENT
	case commonConstants.TEST_STATUS_COMPLETED_SENT:
		return commonConstants.TEST_STATUS_COMPLETED_SENT
	}
	return status
}

func isTatBreached(labTat string, testStatus string) bool {
	testStatuses := []string{commonConstants.TEST_STATUS_REQUESTED,
		commonConstants.TEST_STATUS_RESULT_PENDING,
		commonConstants.TEST_STATUS_RERUN_REQUESTED}
	if labTat != "" {
		labTatTime, _ := time.Parse(commonConstants.DateTimeUTCLayoutWithoutTZOffset, labTat)
		labTatTime = labTatTime.Add(-30 * time.Minute)
		if time.Now().After(labTatTime) && commonUtils.SliceContainsString(testStatuses, testStatus) {
			return true
		}
	}
	return false
}

func isSampleRejectedForTest(testStatus string, isSampleRejected bool) bool {
	return (testStatus == commonConstants.TEST_STATUS_REQUESTED ||
		testStatus == commonConstants.TEST_STATUS_LAB_RECEIVED ||
		testStatus == commonConstants.TEST_STATUS_RERUN_REQUESTED ||
		testStatus == commonConstants.TEST_STATUS_REJECTED) &&
		isSampleRejected
}

func showReprintBarcode(sampleStatus string) bool {
	return sampleStatus == commonConstants.SampleAccessioned ||
		sampleStatus == commonConstants.SampleSynced ||
		sampleStatus == commonConstants.SampleReceived ||
		sampleStatus == commonConstants.SamplePartiallyRejected ||
		sampleStatus == commonConstants.SampleTransferred ||
		sampleStatus == commonConstants.SampleInTransfer ||
		sampleStatus == commonConstants.SampleOutsourced
}

func showCollectionButton(sampleStatus string) bool {
	return sampleStatus == commonConstants.SampleDefault ||
		sampleStatus == commonConstants.SampleNotCollectedEmedic
}

func updateTestStatusInTestDetails(testDetails []structures.InfoScreenTestDetails) {
	for i := range testDetails {
		if (testDetails[i].Status == commonConstants.TEST_STATUS_REQUESTED) &&
			(testDetails[i].SampleStatus == commonConstants.SampleAccessioned ||
				testDetails[i].SampleStatus == commonConstants.SampleSynced ||
				testDetails[i].SampleStatus == commonConstants.SampleReceived ||
				testDetails[i].SampleStatus == commonConstants.SamplePartiallyRejected ||
				testDetails[i].SampleStatus == commonConstants.SampleOutsourced) {
			testDetails[i].Status = commonConstants.TEST_STATUS_LAB_RECEIVED
		}
	}
}

func (searchService *SearchService) getTestDetailsByOrderId(omsOrderId string, servicingCityCode string) (
	[]structures.InfoScreenTestDetails, *commonStructures.CommonError) {
	testDetails, cErr := searchService.SearchDao.GetTestDetailsByOrderId(omsOrderId, servicingCityCode)
	if cErr != nil {
		return testDetails, cErr
	}
	if len(testDetails) == 0 {
		return testDetails, &commonStructures.CommonError{
			StatusCode: http.StatusNotFound,
			Message:    commonConstants.ERROR_ORDER_ID_NOT_FOUND,
		}
	}
	updateTestStatusInTestDetails(testDetails)
	return testDetails, nil
}

func (searchService *SearchService) getTestDetailsByBarcode(barcode, serviceType string, labId uint) (
	[]structures.InfoScreenTestDetails, *commonStructures.CommonError) {
	testDetails, err := searchService.SearchDao.GetTestDetailsByBarcode(barcode, serviceType, labId)
	if err != nil {
		return testDetails, &commonStructures.CommonError{
			StatusCode: http.StatusInternalServerError,
			Message:    constants.ERROR_WHILE_FETCHING_TEST_DETAILS,
		}
	}
	if len(testDetails) == 0 {
		return testDetails, &commonStructures.CommonError{
			StatusCode: http.StatusNotFound,
			Message:    constants.ERROR_NO_TEST_FOUND_IN_BARCODE,
		}
	}
	filteredTestDetails := []structures.InfoScreenTestDetails{}
	for _, testDetail := range testDetails {
		if testDetail.ProcessingLabId == labId {
			filteredTestDetails = append(filteredTestDetails, testDetail)
		}
	}
	if len(filteredTestDetails) == 0 {
		return filteredTestDetails, &commonStructures.CommonError{
			StatusCode: http.StatusNotFound,
			Message:    constants.ERROR_NO_INHOUSE_TESTS_FOUND,
		}
	}
	updateTestStatusInTestDetails(filteredTestDetails)
	return filteredTestDetails, nil
}

func createBasicSearchResponse(orderDetails structures.InfoScreenOrderDetails,
	labIdMap map[uint]commonStructures.Lab) structures.InfoScreenSearchResponse {

	etsSearchResponse := structures.InfoScreenSearchResponse{
		Patient: structures.InfoScreenPatientDetailsResponse{
			Name:        orderDetails.PatientName,
			ExpectedDob: orderDetails.PatientExpectedDob,
			Dob:         orderDetails.PatientDob,
			Gender:      orderDetails.PatientGender,
		},
		Request: structures.InfoScreenRequestDetailsResponse{
			ID:                fmt.Sprint(orderDetails.RequestID),
			ServicingLabId:    orderDetails.ServicingLabId,
			ServicingCityCode: orderDetails.ServicingCityCode,
			ServicingLabName:  labIdMap[orderDetails.ServicingLabId].LabName,
			MetaData: map[string]interface{}{
				"is_camp_request": orderDetails.CampId != 0,
			},
		},
		Order: structures.InfoScreenOrderDetailsResponse{
			ID:        orderDetails.OrderID,
			Status:    orderDetails.OrderStatus,
			CreatedAt: orderDetails.CreatedAt,
			DeletedAt: orderDetails.DeletedAt,
		},
	}
	return etsSearchResponse
}

func createBarcodeVisitMapAndVisitSampleMap(basicVisitDetails []structures.InfoScreenBasicVisitDetails,
	servicingLabId uint, labIdMap map[uint]commonStructures.Lab) (
	map[string][]uint, map[uint][]structures.InfoScreenSampleDetailsResponse) {

	visitSampleMap := map[string]map[uint]structures.InfoScreenSampleDetailsResponse{}
	visitIdToSampleNumbersMap := map[string][]uint{}
	sampleNumberToSampleDetailsMap := map[uint][]structures.InfoScreenSampleDetailsResponse{}
	for _, basicVisitDetail := range basicVisitDetails {
		if _, ok := visitSampleMap[basicVisitDetail.VisitID]; !ok {
			visitSampleMap[basicVisitDetail.VisitID] = map[uint]structures.InfoScreenSampleDetailsResponse{}
		}
		if _, ok := visitIdToSampleNumbersMap[basicVisitDetail.VisitID]; !ok {
			visitIdToSampleNumbersMap[basicVisitDetail.VisitID] = []uint{}
		}
		visitIdToSampleNumbersMap[basicVisitDetail.VisitID] = append(visitIdToSampleNumbersMap[basicVisitDetail.VisitID],
			basicVisitDetail.SampleNumber)

		if _, ok := sampleNumberToSampleDetailsMap[basicVisitDetail.SampleNumber]; !ok {
			sampleNumberToSampleDetailsMap[basicVisitDetail.SampleNumber] = []structures.InfoScreenSampleDetailsResponse{}
		}
		labId := basicVisitDetail.LabId
		if labId == 0 {
			labId = basicVisitDetail.DestinationLabId
		}
		if labId == 0 {
			labId = servicingLabId
		}
		currentLabId := labId
		if basicVisitDetail.CurrentStatus == commonConstants.SampleOutsourced {
			currentLabId = basicVisitDetail.DestinationLabId
		}
		labName := ""
		if labId != 0 {
			labName = labIdMap[labId].LabName
		}
		sampleNumberToSampleDetailsMap[basicVisitDetail.SampleNumber] = append(
			sampleNumberToSampleDetailsMap[basicVisitDetail.SampleNumber], structures.InfoScreenSampleDetailsResponse{
				SampleId:         basicVisitDetail.SampleId,
				ParentSampleId:   basicVisitDetail.ParentSampleId,
				LabId:            labId,
				CurrentLabId:     currentLabId,
				LabName:          labName,
				CurrentLabName:   labIdMap[currentLabId].LabName,
				SampleNumber:     basicVisitDetail.SampleNumber,
				Barcode:          basicVisitDetail.Barcode,
				VialTypeID:       basicVisitDetail.VialTypeID,
				CurrentStatus:    basicVisitDetail.CurrentStatus,
				CollectedAt:      basicVisitDetail.CollectedAt,
				ReceivedAt:       basicVisitDetail.ReceivedAt,
				TransferredAt:    basicVisitDetail.TransferredAt,
				OutsourcedAt:     basicVisitDetail.OutsourcedAt,
				RejectedAt:       basicVisitDetail.RejectedAt,
				NotReceivedAt:    basicVisitDetail.NotReceivedAt,
				LisSyncAt:        basicVisitDetail.LisSyncAt,
				BarcodeScannedAt: basicVisitDetail.BarcodeScannedAt,
				CreatedAt:        basicVisitDetail.CreatedAt,
				DeletedAt:        basicVisitDetail.DeletedAt,
				Metadata: map[string]interface{}{
					"is_sample_rejected":           basicVisitDetail.CurrentStatus == commonConstants.SampleRejected,
					"is_sample_partially_rejected": basicVisitDetail.CurrentStatus == commonConstants.SamplePartiallyRejected,
					"show_reprint_barcode":         showReprintBarcode(basicVisitDetail.CurrentStatus),
					"show_collection_button":       showCollectionButton(basicVisitDetail.CurrentStatus),
				},
			})
	}
	return visitIdToSampleNumbersMap, sampleNumberToSampleDetailsMap
}

func createSampleTestMap(testDetails []structures.InfoScreenTestDetails,
	labIdMap map[uint]commonStructures.Lab) map[uint]map[uint]structures.InfoScreenTestDetailsResponse {
	sampleTestMap := map[uint]map[uint]structures.InfoScreenTestDetailsResponse{}
	for _, testDetail := range testDetails {
		if _, ok := sampleTestMap[testDetail.SampleNumber]; !ok {
			sampleTestMap[testDetail.SampleNumber] = map[uint]structures.InfoScreenTestDetailsResponse{}
		}
		status := getTestStatus(testDetail.Status, testDetail.OmsStatus)
		sampleTestMap[testDetail.SampleNumber][testDetail.ID] = structures.InfoScreenTestDetailsResponse{
			ID:                testDetail.ID,
			AlnumTestId:       testDetail.AlnumTestId,
			Name:              testDetail.TestName,
			Department:        testDetail.Department,
			Status:            status,
			StatusLabel:       commonConstants.TEST_STATUSES_LABEL_MAP[status],
			ProcessingLabId:   testDetail.ProcessingLabId,
			ProcessingLabName: labIdMap[testDetail.ProcessingLabId].LabName,
			CreatedAt:         testDetail.CreatedAt,
			DeletedAt:         testDetail.DeletedAt,
			Metadata: map[string]interface{}{
				"is_tat_breached": isTatBreached(testDetail.LabTat, testDetail.Status),
			},
		}
	}
	return sampleTestMap
}

func getSortedTestDetailsFromTestMap(testMap map[uint]structures.InfoScreenTestDetailsResponse) []structures.InfoScreenTestDetailsResponse {
	testIDs := []uint{}
	for testID := range testMap {
		testIDs = append(testIDs, testID)
	}
	commonUtils.SortUintSlice(testIDs)

	testDetails := []structures.InfoScreenTestDetailsResponse{}
	for _, testID := range testIDs {
		testDetails = append(testDetails, testMap[testID])
	}
	return testDetails
}

func getSortedSampleDetailsResponse(sampleDetailsResponse []structures.InfoScreenSampleDetailsResponse) []structures.InfoScreenSampleDetailsResponse {
	sort.Slice(sampleDetailsResponse, func(i, j int) bool {
		return sampleDetailsResponse[i].SampleNumber < sampleDetailsResponse[j].SampleNumber
	})
	return sampleDetailsResponse
}

func getSortedVisitDetailsResponse(visitDetailsResponse []structures.InfoScreenVisitDetailsResponse) []structures.InfoScreenVisitDetailsResponse {
	sort.Slice(visitDetailsResponse, func(i, j int) bool {
		return visitDetailsResponse[i].ID < visitDetailsResponse[j].ID
	})
	return visitDetailsResponse
}

func modifyVisitDetailResponseForRejectedSamples(visitDetailsResponse []structures.InfoScreenVisitDetailsResponse) []structures.InfoScreenVisitDetailsResponse {
	rejectedSamples, nonRejectedSamples := []structures.InfoScreenSampleDetailsResponse{}, []structures.InfoScreenSampleDetailsResponse{}
	emptyVisitIndex := -1
	rejectedParentSampleIds := []uint{}
	for visitIndex := range visitDetailsResponse {
		if visitDetailsResponse[visitIndex].ID == "" {
			emptyVisitIndex = visitIndex
			for _, sample := range visitDetailsResponse[visitIndex].Samples {
				if sample.Metadata["is_sample_rejected"].(bool) {
					rejectedSamples = append(rejectedSamples, sample)
					if sample.ParentSampleId != 0 {
						rejectedParentSampleIds = append(rejectedParentSampleIds, sample.ParentSampleId)
					}
				} else {
					nonRejectedSamples = append(nonRejectedSamples, sample)
				}
			}
			break
		}
	}

	finalNonRejectedSamples, extraRejectedSamples :=
		[]structures.InfoScreenSampleDetailsResponse{}, []structures.InfoScreenSampleDetailsResponse{}

	// Remove the empty visit from the response
	if emptyVisitIndex != -1 {
		visitDetailsResponse = append(visitDetailsResponse[:emptyVisitIndex], visitDetailsResponse[emptyVisitIndex+1:]...)
	}

	for _, nonRejectedSample := range nonRejectedSamples {
		if commonUtils.SliceContainsUint(rejectedParentSampleIds, nonRejectedSample.SampleId) {
			// If the parent sample is rejected, we don't want to show this sample under normal samples
			extraRejectedSamples = append(extraRejectedSamples, nonRejectedSample)
		} else {
			finalNonRejectedSamples = append(finalNonRejectedSamples, nonRejectedSample)
		}
	}

	if len(finalNonRejectedSamples) > 0 {
		visitDetailsResponse = append(visitDetailsResponse, structures.InfoScreenVisitDetailsResponse{
			Samples:  finalNonRejectedSamples,
			Metadata: map[string]interface{}{},
		})
	}

	rejectedSamples = append(rejectedSamples, extraRejectedSamples...)
	if len(rejectedSamples) > 0 {
		visitDetailsResponse = append(visitDetailsResponse, structures.InfoScreenVisitDetailsResponse{
			Samples: rejectedSamples,
			Metadata: map[string]interface{}{
				"are_all_samples_rejected": true,
			},
		})
	}
	return visitDetailsResponse
}

func createSearchResponse(orderDetails structures.InfoScreenOrderDetails,
	basicVisitDetails []structures.InfoScreenBasicVisitDetails, testDetails []structures.InfoScreenTestDetails,
	labIdMap map[uint]commonStructures.Lab) structures.InfoScreenSearchResponse {
	etsSearchResponse := createBasicSearchResponse(orderDetails, labIdMap)

	visitIdToSampleNumbersMap, sampleNumberToSampleDetailsMap := createBarcodeVisitMapAndVisitSampleMap(basicVisitDetails,
		orderDetails.ServicingLabId, labIdMap)

	sampleTestMap := createSampleTestMap(testDetails, labIdMap)
	sampleIdUsedMap := map[uint]bool{}

	visitDetailsResponse := []structures.InfoScreenVisitDetailsResponse{}
	for visitId, sampleNumbers := range visitIdToSampleNumbersMap {
		if visitId == "" {
			continue
		}
		sampleDetailsResponse := []structures.InfoScreenSampleDetailsResponse{}
		for _, sampleNumber := range sampleNumbers {
			tests := getSortedTestDetailsFromTestMap(sampleTestMap[sampleNumber])
			sampleDetails := sampleNumberToSampleDetailsMap[sampleNumber]
			for index := range sampleDetails {
				sampleDetails[index].Tests = tests
			}
			for _, sampleDetail := range sampleDetails {
				if !sampleIdUsedMap[sampleDetail.SampleId] {
					sampleDetailsResponse = append(sampleDetailsResponse, sampleDetail)
					sampleIdUsedMap[sampleDetail.SampleId] = true
				}
			}
		}
		visitDetailsResponse = append(visitDetailsResponse, structures.InfoScreenVisitDetailsResponse{
			ID:      visitId,
			Samples: getSortedSampleDetailsResponse(sampleDetailsResponse),
			Metadata: map[string]interface{}{
				"is_resync_enabled": true,
			},
		})
	}

	emptyVisitIdSampleNumbers := visitIdToSampleNumbersMap[""]
	if len(emptyVisitIdSampleNumbers) > 0 {
		sampleDetailsResponse := []structures.InfoScreenSampleDetailsResponse{}
		for _, sampleNumber := range emptyVisitIdSampleNumbers {
			tests := getSortedTestDetailsFromTestMap(sampleTestMap[sampleNumber])
			sampleDetails := sampleNumberToSampleDetailsMap[sampleNumber]
			for index := range sampleDetails {
				sampleDetails[index].Tests = tests
			}
			for _, sampleDetail := range sampleDetails {
				if !sampleIdUsedMap[sampleDetail.SampleId] {
					sampleDetailsResponse = append(sampleDetailsResponse, sampleDetail)
					sampleIdUsedMap[sampleDetail.SampleId] = true
				}
			}
		}
		visitDetailsResponse = append(visitDetailsResponse, structures.InfoScreenVisitDetailsResponse{
			ID:       "",
			Samples:  getSortedSampleDetailsResponse(sampleDetailsResponse),
			Metadata: map[string]interface{}{},
		})

	}

	visitDetailsResponse = getSortedVisitDetailsResponse(visitDetailsResponse)
	visitDetailsResponse = modifyVisitDetailResponseForRejectedSamples(visitDetailsResponse)

	etsSearchResponse.Visits = visitDetailsResponse
	return etsSearchResponse
}

func createResponseForBarcodeDetails(barcode string, orderDetails structures.InfoScreenOrderDetails, visitDetails structures.InfoScreenBasicVisitDetails, testDetails []structures.InfoScreenTestDetails) structures.BarcodeDetailsResponse {
	barcodeDetailsResponse := structures.BarcodeDetailsResponse{
		Barcode:    barcode,
		VialTypeID: visitDetails.VialTypeID,
		Patient: structures.InfoScreenPatientDetailsResponse{
			Name:        orderDetails.PatientName,
			ExpectedDob: orderDetails.PatientExpectedDob,
			Dob:         orderDetails.PatientDob,
			Gender:      orderDetails.PatientGender,
		},
		CreatedAt:        visitDetails.CreatedAt,
		DeletedAt:        visitDetails.DeletedAt,
		VisitID:          visitDetails.VisitID,
		OrderID:          orderDetails.OrderID,
		SampleReceivedAt: visitDetails.ReceivedAt,
	}

	isTatBreachedForOrder := false
	testDetailsResponse := []structures.InfoScreenTestDetailsResponse{}
	for _, testDetail := range testDetails {
		isTatBreachedForTest := isTatBreached(testDetail.LabTat, testDetail.Status)
		if isTatBreachedForTest {
			isTatBreachedForOrder = true
		}
		testDetailsResponse = append(testDetailsResponse, structures.InfoScreenTestDetailsResponse{
			ID:          testDetail.ID,
			AlnumTestId: testDetail.AlnumTestId,
			Name:        testDetail.TestName,
			Status:      testDetail.Status,
			Department:  testDetail.Department,
			StatusLabel: commonConstants.TEST_STATUSES_LABEL_MAP[testDetail.Status],
			CreatedAt:   testDetail.CreatedAt,
			DeletedAt:   testDetail.DeletedAt,
			Metadata: map[string]interface{}{
				"is_tat_breached":    isTatBreachedForTest,
				"is_sample_rejected": isSampleRejectedForTest(testDetail.Status, testDetail.IsSampleRejected),
			},
		})
	}
	barcodeDetailsResponse.Tests = testDetailsResponse
	barcodeDetailsResponse.Metadata = map[string]interface{}{
		"is_tat_breached":              isTatBreachedForOrder,
		"is_camp_order":                orderDetails.CampId != 0,
		"is_sample_rejected":           visitDetails.CurrentStatus == commonConstants.SampleRejected,
		"is_sample_partially_rejected": visitDetails.CurrentStatus == commonConstants.SamplePartiallyRejected,
		"show_reprint_barcode":         showReprintBarcode(visitDetails.CurrentStatus),
	}
	return barcodeDetailsResponse
}

func prependBarcodeInSearchResponse(etsSearchResponse structures.InfoScreenSearchResponse, barcode string) structures.InfoScreenSearchResponse {
	visits := etsSearchResponse.Visits
	swapSuccess := false
	for visitIndex := range visits {
		samples := visits[visitIndex].Samples
		for sampleIndex := range samples {
			if samples[sampleIndex].Barcode == barcode {
				samples[0], samples[sampleIndex] = samples[sampleIndex], samples[0]
				swapSuccess = true
				break
			}
		}
		if swapSuccess {
			visits[0], visits[visitIndex] = visits[visitIndex], visits[0]
			break
		}
	}
	etsSearchResponse.Visits = visits
	return etsSearchResponse
}

func prependVisitInSearchResponse(etsSearchResponse structures.InfoScreenSearchResponse, visitId string) structures.InfoScreenSearchResponse {
	visits := etsSearchResponse.Visits
	for visitIndex := range visits {
		if visits[visitIndex].ID == visitId {
			visits[0], visits[visitIndex] = visits[visitIndex], visits[0]
			break
		}
	}
	etsSearchResponse.Visits = visits
	return etsSearchResponse
}

func prependSessionLabIdSamplesInSearchSearchResponse(
	etsSearchResponse structures.InfoScreenSearchResponse, sessionLabId uint, searchType string) structures.InfoScreenSearchResponse {
	if sessionLabId == 0 {
		return etsSearchResponse
	}

	visits := etsSearchResponse.Visits
	if len(visits) <= 1 {
		return etsSearchResponse
	}

	startingIndex := 0
	if searchType != constants.INFO_SCREEN_SEARCH_TYPE_ORDER_ID {
		startingIndex = 1 // Exclude first index for non-order-id searches
	}

	visitLabIdMap := make(map[string]uint)
	for i := startingIndex; i < len(visits); i++ { // Exclude first index
		visit := visits[i]
		if visit.ID == "" && visit.Metadata["are_all_samples_rejected"] != nil && visit.Metadata["are_all_samples_rejected"].(bool) {
			visitLabIdMap[visit.ID] = 0 // This is to ensure that we don't prepend this visit
			continue
		}
		parentLabId, childLabId := uint(0), uint(0)
		if visit.ID != "" {
			for _, sample := range visit.Samples {
				if sample.ParentSampleId != 0 {
					childLabId = sample.CurrentLabId
				} else {
					parentLabId = sample.CurrentLabId
				}
			}
			visitLabIdMap[visit.ID] = parentLabId
			if childLabId != 0 {
				visitLabIdMap[visit.ID] = childLabId
			}
		} else {
			for _, sample := range visit.Samples {
				if sample.ParentSampleId != 0 {
					childLabId = sample.LabId
				} else {
					parentLabId = sample.LabId
				}
				visitLabIdMap[visit.ID] = parentLabId
				if childLabId != 0 {
					visitLabIdMap[visit.ID] = childLabId
				}
			}
		}
	}

	// Prepend visits with labId == sessionLabId
	prependedVisits := []structures.InfoScreenVisitDetailsResponse{}
	appendedVisits := []structures.InfoScreenVisitDetailsResponse{}
	for i := startingIndex; i < len(visits); i++ {
		visit := visits[i]
		if visitLabIdMap[visit.ID] == sessionLabId {
			prependedVisits = append(prependedVisits, visit)
		} else {
			appendedVisits = append(appendedVisits, visit)
		}
	}

	finalVisits := []structures.InfoScreenVisitDetailsResponse{}
	if startingIndex != 0 {
		finalVisits = append(finalVisits, visits[0]) // Always keep the first visit
	}
	finalVisits = append(finalVisits, append(prependedVisits, appendedVisits...)...)

	etsSearchResponse.Visits = finalVisits
	return etsSearchResponse
}

func (searchService *SearchService) GetBarcodeDetails(barcode, serviceType string, labId uint) (
	structures.BarcodeDetailsResponse, *commonStructures.CommonError) {
	barcode = strings.TrimSpace(barcode)
	if serviceType == "" {
		serviceType = constants.ServiceTypeScan
	}

	if barcode == "" {
		return structures.BarcodeDetailsResponse{}, &commonStructures.CommonError{
			StatusCode: http.StatusBadRequest,
			Message:    commonConstants.ERROR_BARCODE_REQUIRED,
		}
	}
	if labId == 0 {
		return structures.BarcodeDetailsResponse{}, &commonStructures.CommonError{
			StatusCode: http.StatusBadRequest,
			Message:    commonConstants.ERROR_LAB_NOT_FOUND,
		}
	}

	orderDetails, cErr := searchService.SearchDao.GetOrderDetailsByBarcode(barcode, serviceType)
	if cErr != nil {
		return structures.BarcodeDetailsResponse{}, cErr
	}

	testDetails, cErr := searchService.getTestDetailsByBarcode(barcode, serviceType, labId)
	if cErr != nil {
		return structures.BarcodeDetailsResponse{}, cErr
	}

	visitDetails, cErr := searchService.SearchDao.GetVisitBasicDetailsByBarcode(barcode, serviceType, labId)
	if cErr != nil {
		return structures.BarcodeDetailsResponse{}, cErr
	}

	return createResponseForBarcodeDetails(barcode, orderDetails, visitDetails, testDetails), nil
}

func (searchService *SearchService) GetInfoScreenDataByBarcode(ctx context.Context, barcode string, labId uint) (
	structures.InfoScreenSearchResponse, *commonStructures.CommonError) {
	barcode = strings.TrimSpace(barcode)
	if barcode == "" {
		return structures.InfoScreenSearchResponse{}, &commonStructures.CommonError{
			StatusCode: http.StatusBadRequest,
			Message:    commonConstants.ERROR_BARCODE_REQUIRED,
		}
	}
	infoScreenSearchResponse := structures.InfoScreenSearchResponse{}

	orderDetails, cErr := searchService.SearchDao.GetOrderDetailsByBarcode(barcode, constants.ServiceTypeDetails)
	if cErr != nil {
		return infoScreenSearchResponse, cErr
	}

	var (
		basicVisitDetails []structures.InfoScreenBasicVisitDetails
		testDetails       []structures.InfoScreenTestDetails
		labIdMap          map[uint]commonStructures.Lab
	)

	wg := &sync.WaitGroup{}
	errChan := make(chan *commonStructures.CommonError, 2)
	wg.Add(3)

	// Fetch basicVisitDetails
	go func() {
		defer wg.Done()
		var err *commonStructures.CommonError
		basicVisitDetails, err = searchService.SearchDao.GetVisitBasicDetailsByOrderId(orderDetails.OrderID)
		if err != nil {
			errChan <- err
		}
	}()

	// Fetch testDetails
	go func() {
		defer wg.Done()
		var err *commonStructures.CommonError
		testDetails, err = searchService.getTestDetailsByOrderId(orderDetails.OrderID, orderDetails.ServicingCityCode)
		if err != nil {
			errChan <- err
		}
	}()

	// Fetch labIdMap
	go func() {
		defer wg.Done()
		labIdMap = searchService.CdsService.GetLabIdLabMap(ctx)
	}()

	wg.Wait()
	close(errChan)
	for err := range errChan {
		if err != nil {
			return infoScreenSearchResponse, err
		}
	}

	infoScreenSearchResponse = createSearchResponse(orderDetails, basicVisitDetails, testDetails, labIdMap)
	infoScreenSearchResponse = prependBarcodeInSearchResponse(infoScreenSearchResponse, barcode)
	return prependSessionLabIdSamplesInSearchSearchResponse(infoScreenSearchResponse, labId,
		constants.INFO_SCREEN_SEARCH_TYPE_BARCODE), nil
}

func (searchService *SearchService) GetInfoScreenDataByOrderId(ctx context.Context, omsOrderId string, labId uint) (
	structures.InfoScreenSearchResponse, *commonStructures.CommonError) {

	etsSearchResponse := structures.InfoScreenSearchResponse{}

	orderDetails, cErr := searchService.SearchDao.GetOrderDetailsByOrderId(omsOrderId)
	if cErr != nil {
		return etsSearchResponse, cErr
	}

	var (
		basicVisitDetails []structures.InfoScreenBasicVisitDetails
		testDetails       []structures.InfoScreenTestDetails
		labIdMap          map[uint]commonStructures.Lab
	)

	wg := &sync.WaitGroup{}
	errChan := make(chan *commonStructures.CommonError, 2)
	wg.Add(3)

	// Fetch basicVisitDetails
	go func() {
		defer wg.Done()
		var err *commonStructures.CommonError
		basicVisitDetails, err = searchService.SearchDao.GetVisitBasicDetailsByOrderId(omsOrderId)
		if err != nil {
			errChan <- err
		}
	}()

	// Fetch testDetails
	go func() {
		defer wg.Done()
		var err *commonStructures.CommonError
		testDetails, err = searchService.getTestDetailsByOrderId(omsOrderId, orderDetails.ServicingCityCode)
		if err != nil {
			errChan <- err
		}
	}()

	// Fetch labIdMap
	go func() {
		defer wg.Done()
		labIdMap = searchService.CdsService.GetLabIdLabMap(ctx)
	}()

	wg.Wait()
	close(errChan)
	for err := range errChan {
		if err != nil {
			return etsSearchResponse, err
		}
	}

	etsSearchResponse = createSearchResponse(orderDetails, basicVisitDetails, testDetails, labIdMap)
	return prependSessionLabIdSamplesInSearchSearchResponse(etsSearchResponse, labId,
		constants.INFO_SCREEN_SEARCH_TYPE_ORDER_ID), nil
}

func (searchService *SearchService) GetInfoScreenDataByVisitId(ctx context.Context, visitId string, labId uint) (
	structures.InfoScreenSearchResponse, *commonStructures.CommonError) {
	visitId = strings.TrimSpace(visitId)
	if visitId == "" {
		return structures.InfoScreenSearchResponse{}, &commonStructures.CommonError{
			StatusCode: http.StatusBadRequest,
			Message:    commonConstants.ERROR_VISIT_ID_REQUIRED,
		}
	}
	etsSearchResponse := structures.InfoScreenSearchResponse{}

	orderDetails, cErr := searchService.SearchDao.GetOrderDetailsByVisitId(visitId)
	if cErr != nil {
		return etsSearchResponse, cErr
	}

	var (
		basicVisitDetails []structures.InfoScreenBasicVisitDetails
		testDetails       []structures.InfoScreenTestDetails
		labIdMap          map[uint]commonStructures.Lab
	)

	wg := &sync.WaitGroup{}
	errChan := make(chan *commonStructures.CommonError, 2)
	wg.Add(3)

	// Fetch basicVisitDetails
	go func() {
		defer wg.Done()
		var err *commonStructures.CommonError
		basicVisitDetails, err = searchService.SearchDao.GetVisitBasicDetailsByOrderId(orderDetails.OrderID)
		if err != nil {
			errChan <- err
		}
	}()

	// Fetch testDetails
	go func() {
		defer wg.Done()
		var err *commonStructures.CommonError
		testDetails, err = searchService.getTestDetailsByOrderId(orderDetails.OrderID, orderDetails.ServicingCityCode)
		if err != nil {
			errChan <- err
		}
	}()

	// Fetch labIdMap
	go func() {
		defer wg.Done()
		labIdMap = searchService.CdsService.GetLabIdLabMap(ctx)
	}()

	wg.Wait()
	close(errChan)
	for err := range errChan {
		if err != nil {
			return etsSearchResponse, err
		}
	}

	etsSearchResponse = createSearchResponse(orderDetails, basicVisitDetails, testDetails, labIdMap)
	etsSearchResponse = prependVisitInSearchResponse(etsSearchResponse, visitId)
	return prependSessionLabIdSamplesInSearchSearchResponse(etsSearchResponse, labId,
		constants.INFO_SCREEN_SEARCH_TYPE_VISIT_ID), nil
}
