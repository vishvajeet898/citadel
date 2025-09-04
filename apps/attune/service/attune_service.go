package service

import (
	"context"
	"fmt"
	"net/http"
	"sort"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/Orange-Health/citadel/apps/attune/structures"
	commonConstants "github.com/Orange-Health/citadel/common/constants"
	commonStructures "github.com/Orange-Health/citadel/common/structures"
	commonUtils "github.com/Orange-Health/citadel/common/utils"
	commonModels "github.com/Orange-Health/citadel/models"
)

func getSampleIdString(samples []commonModels.Sample) string {
	values := []string{}
	for _, sample := range samples {
		values = append(values, fmt.Sprint(sample.Id))
	}
	sort.Slice(values, func(index1 int, index2 int) bool { return values[index1] < values[index2] })
	return strings.Join(values, ",")
}

func getPatientFormattedDob(patientDetails commonModels.PatientDetail) string {
	loc, _ := time.LoadLocation(commonConstants.LocalTimeZoneLocation)
	dob := patientDetails.Dob
	if dob == nil {
		dob = patientDetails.ExpectedDob
	}
	*dob = (*dob).In(loc)
	formattedDob := dob.Format(commonConstants.DateLayoutReverse)
	return formattedDob
}

func ConstructTestSampleListFromLisSyncData(test commonStructures.AttuneOrderInfo, uniqueSampleMap map[string]bool,
	testList []commonStructures.LisTestInfo, sampleList []commonStructures.LisSampleInfo) (
	[]commonStructures.LisTestInfo, []commonStructures.LisSampleInfo) {

	if test.TestType == commonConstants.InvestigationShortHand && test.BarcodeNumber != "" {
		if _, ok := uniqueSampleMap[test.BarcodeNumber]; !ok {
			uniqueSampleMap[test.BarcodeNumber] = true
			syncedSampleData := commonStructures.LisSampleInfo{
				SampleName: test.SampleName,
				Barcode:    test.BarcodeNumber,
			}
			sampleList = append(sampleList, syncedSampleData)
		}
	} else if test.TestType == commonConstants.GroupShortHand {
		for _, groupItem := range test.OrderContentListInfo {
			if groupItem.SampleName != "" && groupItem.BarcodeNumber != "" {
				if _, ok := uniqueSampleMap[groupItem.BarcodeNumber]; !ok {
					uniqueSampleMap[groupItem.BarcodeNumber] = true
					syncedSampleData := commonStructures.LisSampleInfo{
						SampleName: groupItem.SampleName,
						Barcode:    groupItem.BarcodeNumber,
					}
					sampleList = append(sampleList, syncedSampleData)
				}
			}
		}
	}
	syncedTestData := commonStructures.LisTestInfo{
		TestCode: test.TestCode,
		TestName: test.TestName,
		TestType: test.TestType,
	}
	testList = append(testList, syncedTestData)

	return testList, sampleList
}

func ConstructSyncDetailsForVisitID(attuneResponse commonStructures.AttuneOrderResponse,
	visitID string) commonStructures.LisSyncDetails {
	syncDetails := commonStructures.LisSyncDetails{}
	syncDetails.VisitID = visitID
	if attuneResponse.PatientInfo.FirstName != "" {
		syncDetails.PatientName = attuneResponse.PatientInfo.FirstName
	} else {
		syncDetails.PatientName = attuneResponse.PatientInfo.Name
	}

	var testList []commonStructures.LisTestInfo = []commonStructures.LisTestInfo{}
	var sampleList []commonStructures.LisSampleInfo = []commonStructures.LisSampleInfo{}
	var uniqueSampleMap map[string]bool = make(map[string]bool)
	var uniqueTestSampleMap map[string]map[string]string = make(map[string]map[string]string)
	for _, test := range attuneResponse.OrderInfo {
		if test.TestStatus != commonConstants.AttuneTestStatusCancel &&
			test.TestStatus != commonConstants.AttuneTestStatusRetest {
			if _, ok := uniqueTestSampleMap[test.TestID]; !ok {
				uniqueTestSampleMap[test.TestID] = map[string]string{
					test.SampleName: test.BarcodeNumber,
				}
				testList, sampleList = ConstructTestSampleListFromLisSyncData(test, uniqueSampleMap, testList, sampleList)
			} else if ok {
				uniqueTestSampleMap[test.TestID][test.SampleName] = test.BarcodeNumber
				testList, sampleList = ConstructTestSampleListFromLisSyncData(test, uniqueSampleMap, testList, sampleList)
			}
		}
	}
	syncDetails.SyncedTests = testList
	syncDetails.SyncedSamples = sampleList
	syncDetails.LisSyncTime = attuneResponse.PatientVisitInfo.VisitDate
	if len(attuneResponse.OrderInfo) > 0 {
		for _, orderInfo := range attuneResponse.OrderInfo {
			if orderInfo.ResultCapturedAt != "" {
				syncDetails.LisSyncTime = orderInfo.ResultCapturedAt
				break
			}
		}
	}

	return syncDetails
}

func getAttunePayloadForSyncingData(attunePayloadMeta structures.AttunePayloadMeta) (
	commonStructures.AttuneSyncDataToLisRequest, *commonStructures.CommonError) {
	location, err := commonUtils.GetLocalLocation()
	if err != nil {
		return commonStructures.AttuneSyncDataToLisRequest{}, &commonStructures.CommonError{
			Message:    commonConstants.ERROR_FAILED_TO_GET_LOCAL_LOCATION,
			StatusCode: http.StatusInternalServerError,
		}
	}
	return commonStructures.AttuneSyncDataToLisRequest{
		MessageType: commonConstants.AttuneMessageTypeNew,
		OrderID:     attunePayloadMeta.VisitId,
		OrgCode:     commonUtils.GetAttuneOrgCodeByLabId(attunePayloadMeta.LabId),
		PatientInfo: commonStructures.AttunePatientInfo{
			Salutation:    attunePayloadMeta.PatientSalutation,
			PatientId:     attunePayloadMeta.PatientId,
			PatientNumber: attunePayloadMeta.PatientId,
			FirstName:     attunePayloadMeta.PatientName,
			Gender:        attunePayloadMeta.PatientGender,
			Dob:           attunePayloadMeta.PatientDob,
			MobileNumber:  attunePayloadMeta.PatientNumber,
			Name:          attunePayloadMeta.PatientName,
		},
		PatientVisitInfo: commonStructures.AttunePatientVisitInfo{
			VisitDate:      time.Now().In(location).Format(commonConstants.DateTimeInSecLayout),
			CollectedDate:  attunePayloadMeta.SampleCollectedAt,
			ClientName:     commonConstants.AttuneClientName,
			SrfId:          attunePayloadMeta.SrfId,
			ReportLanguage: nil,
		},
		TestDetailsList: attunePayloadMeta.AttuneTestDetailsList,
	}, nil
}

func (attuneService *AttuneService) GetAttuneSampleDetailsToCancel(ctx context.Context,
	sampleInfo commonStructures.SampleInfo) (commonStructures.AttuneSampleDetails, *commonStructures.CommonError) {
	masterVialType, cErr := attuneService.CdsService.GetMasterVialType(ctx, sampleInfo.VialTypeId)
	if cErr != nil {
		return commonStructures.AttuneSampleDetails{}, cErr
	}

	sampleDetails := commonStructures.AttuneSampleDetails{
		SID:         sampleInfo.Barcode,
		SampleID:    fmt.Sprint(masterVialType.SampleId),
		ContainerID: fmt.Sprint(masterVialType.ContainerId),
	}

	return sampleDetails, nil
}

func (attuneService *AttuneService) GetTestBarcodeMap(ctx context.Context,
	testSampleMap []commonStructures.AttuneTestSampleMap) map[string][]commonStructures.AttuneSampleDetails {
	testBarcodeMap := map[string][]commonStructures.AttuneSampleDetails{}
	vialTypeMap := attuneService.CdsService.GetMasterVialTypeAsMap(ctx)
	for _, item := range testSampleMap {
		vialType := vialTypeMap[item.VialTypeId]
		barcode := commonStructures.AttuneSampleDetails{
			SID:         item.Barcode,
			SampleID:    fmt.Sprint(vialType.SampleId),
			ContainerID: fmt.Sprint(vialType.ContainerId),
		}
		testBarcodeMap[item.TestId] = append(testBarcodeMap[item.TestId], barcode)
	}

	return testBarcodeMap
}

func (attuneService *AttuneService) GetAttuneTestSampleMapByOmsOrderId(omsOrderId string, sampleIds []uint,
	sampleIdBarcodeMap map[uint]string) ([]commonStructures.AttuneTestSampleMap, *commonStructures.CommonError) {

	attuneTestSampleMaps := []commonStructures.AttuneTestSampleMap{}

	attuneTestSampleMapSnakeCase, cErr := attuneService.AttuneDao.GetAttuneTestSampleMapByOmsOrderId(omsOrderId, sampleIds)
	if cErr != nil {
		return nil, cErr
	}

	for index, testSampleMap := range attuneTestSampleMapSnakeCase {
		if barcode, ok := sampleIdBarcodeMap[testSampleMap.SampleId]; ok {
			attuneTestSampleMapSnakeCase[index].Barcode = barcode
		}
		attuneTestSampleMaps = append(attuneTestSampleMaps, commonStructures.AttuneTestSampleMap{
			TestId:     attuneTestSampleMapSnakeCase[index].TestId,
			VialTypeId: attuneTestSampleMapSnakeCase[index].VialTypeId,
			Barcode:    attuneTestSampleMapSnakeCase[index].Barcode,
		})
	}

	return attuneTestSampleMaps, nil
}

func (attuneService *AttuneService) GetTestDetailsListForSyncingToAttune(omsOrderId string, sampleIds []uint, labId uint,
	testBarcodeMap map[string][]commonStructures.AttuneSampleDetails) ([]commonStructures.AttuneTestDetails,
	*commonStructures.CommonError) {

	testDetailsList := []commonStructures.AttuneTestDetails{}
	testDetails, cErr := attuneService.AttuneDao.GetTestDetailsForSyncingToAttune(omsOrderId, sampleIds, labId)
	if cErr != nil {
		return testDetailsList, cErr
	}

	for _, test := range testDetails {
		// Sync only the tests who have barcode entered and are not rejected
		if _, ok := testBarcodeMap[test.CentralOmsTestId]; ok {
			testDetails := commonStructures.AttuneTestDetails{
				Status:    commonConstants.AttuneTestStatusOrdered,
				TestCode:  test.LisCode,
				TestName:  test.TestName,
				TestType:  test.TestType,
				TestID:    int(test.Id),
				BarcodeNo: testBarcodeMap[test.CentralOmsTestId],
				Price:     "100", // TODO @shrish: Ask Priyanka: Remove this
			}

			testDetailsList = append(testDetailsList, testDetails)
		}
	}

	if len(testDetailsList) == 0 {
		return nil, &commonStructures.CommonError{
			Message:    commonConstants.ERROR_NO_TESTS_TO_BE_SYNCED,
			StatusCode: http.StatusBadRequest,
		}
	}

	return testDetailsList, nil
}

func (attuneService *AttuneService) GetVisitId(ctx context.Context, omsOrderId string) string {
	// Generate visitID
	// FORMAT: 12345V01 : <oms_order_id>V<visit_count>
	// Store visit count in redis for 10 days -> Get count, increment and store
	visitCount, visitIdRedisKey := int(1), fmt.Sprintf(commonConstants.CacheKeySyncToLisVisitCount, omsOrderId)
	err := attuneService.Cache.Get(ctx, visitIdRedisKey, &visitCount)
	if err != nil {
		if err == redis.Nil {
			visitCount = 1
		} else {
			commonUtils.AddLog(ctx, commonConstants.ERROR_LEVEL, commonConstants.ERROR_FAILED_TO_GET_VISIT_ID_FROM_REDIS,
				map[string]interface{}{"oms_order_id": omsOrderId}, err)
		}
	}
	paddedvisitCount := fmt.Sprintf("%02d", visitCount)

	visitId := fmt.Sprintf("%sV%s", omsOrderId, paddedvisitCount)

	err = attuneService.Cache.Set(ctx, visitIdRedisKey, visitCount+1, commonConstants.CacheExpiry10DaysInt)
	if err != nil {
		commonUtils.AddLog(ctx, commonConstants.ERROR_LEVEL, commonConstants.ERROR_FAILED_TO_SET_VISIT_ID_IN_REDIS,
			map[string]interface{}{"oms_order_id": omsOrderId}, err)
	}

	return visitId
}

func (attuneService *AttuneService) SendSlackMessageIfLisSyncFailed(ctx context.Context,
	attuneTestSampleMap []commonStructures.AttuneTestSampleMap,
	testDetailsListForAttune []commonStructures.AttuneTestDetails, patientName string,
	orderDetails commonModels.OrderDetails) {

	slackMessage := commonUtils.GetSlackMessageForLisSyncFailed(patientName, orderDetails.CityCode, orderDetails.OmsOrderId,
		orderDetails.OmsRequestId, attuneTestSampleMap, testDetailsListForAttune)
	attuneService.HealthApiClient.SendGenericSlackMessage(ctx, commonConstants.SlackReportSyncFailureChannel, slackMessage)
}

func (attuneService *AttuneService) GetSampleByVisitId(visitId string) (commonModels.Sample, *commonStructures.CommonError) {
	sample, cErr := attuneService.AttuneDao.GetSampleByVisitId(visitId)
	if cErr != nil {
		return commonModels.Sample{}, cErr
	}
	return sample, nil
}

func (attuneService *AttuneService) GetOrderDetailsByVisitId(visitId string) (commonModels.OrderDetails,
	*commonStructures.CommonError) {
	return attuneService.AttuneDao.GetOrderDetailsByVisitId(visitId)
}

func (attuneService *AttuneService) GetLisSyncData(ctx context.Context, visitId, reportPdfFormat string) (
	commonStructures.LisSyncDetails, commonStructures.AttuneOrderResponse,
	*commonStructures.CommonError) {

	lisSyncDetails := commonStructures.LisSyncDetails{VisitID: visitId}
	attuneOrderResponse := commonStructures.AttuneOrderResponse{}

	sample, cErr := attuneService.AttuneDao.GetSampleByVisitId(visitId)
	if cErr != nil {
		return lisSyncDetails, attuneOrderResponse, cErr
	}

	attuneOrderResponse, cErr = attuneService.AttuneClient.GetPatientVisitDetailsbyVisitNo(ctx, visitId, reportPdfFormat,
		sample.LabId)
	if cErr != nil {
		return lisSyncDetails, attuneOrderResponse, cErr
	}

	if attuneOrderResponse.OrderId == "" {
		commonUtils.AddLog(ctx, commonConstants.ERROR_LEVEL, commonUtils.GetCurrentFunctionName(), map[string]interface{}{
			"visit_id":            visitId,
			"labId":               sample.LabId,
			"attuneOrderResponse": attuneOrderResponse,
			"error":               commonConstants.ERROR_LAB_VISIT_ID_NOT_FOUND_IN_ATTUNE,
		}, nil)
		return lisSyncDetails, attuneOrderResponse, &commonStructures.CommonError{
			Message:    commonConstants.ERROR_LAB_VISIT_ID_NOT_FOUND_IN_ATTUNE,
			StatusCode: http.StatusNotFound,
		}
	}

	lisSyncDetails = ConstructSyncDetailsForVisitID(attuneOrderResponse, visitId)

	return lisSyncDetails, attuneOrderResponse, nil
}

func (attuneService *AttuneService) SyncDataToLisByOmsOrderId(ctx context.Context, omsOrderId string, labId uint,
	samples []commonModels.Sample, samplesMetadata []commonModels.SampleMetadata, sampleIdBarcodeMap map[uint]string) (
	string, *commonStructures.CommonError) {
	sampleIds := []uint{}
	for _, sample := range samples {
		sampleIds = append(sampleIds, sample.Id)
	}

	sampleIdToSampleMetadataMap := make(map[uint]commonModels.SampleMetadata)
	for _, sampleMetadata := range samplesMetadata {
		sampleIdToSampleMetadataMap[sampleMetadata.SampleId] = sampleMetadata
	}

	sampleIdString := getSampleIdString(samples)
	redisKey := fmt.Sprintf(commonConstants.CacheKeyCurrentInProgressSample, omsOrderId)
	value := ""
	err := attuneService.Cache.Get(ctx, redisKey, &value)
	if err == redis.Nil {
		value = sampleIdString
		_ = attuneService.Cache.Set(ctx, redisKey, value, commonConstants.CacheExpiry10DaysInt)

		defer func() {
			_ = attuneService.Cache.Delete(ctx, redisKey)
		}()
	} else if value == sampleIdString {
		return "", &commonStructures.CommonError{
			Message:    commonConstants.ERROR_LIS_SYNC_ALREADY_IN_PROGRESS,
			StatusCode: http.StatusConflict,
		}
	}

	visitId, cErr := attuneService.SyncDataToLis(ctx, omsOrderId, labId, sampleIds, sampleIdBarcodeMap)
	if cErr != nil {
		return "", cErr
	}

	return visitId, nil
}

// TODO: Optimize this function further by leveraging goroutines for fetching data concurrently post stability
func (attuneService *AttuneService) SyncDataToLis(ctx context.Context, omsOrderId string, labId uint, sampleIds []uint,
	sampleIdBarcodeMap map[uint]string) (string, *commonStructures.CommonError) {

	attunePayloadMeta := structures.AttunePayloadMeta{}
	var cErr *commonStructures.CommonError
	attunePayloadMeta.AttuneTestSampleMap, cErr = attuneService.GetAttuneTestSampleMapByOmsOrderId(omsOrderId,
		sampleIds, sampleIdBarcodeMap)
	if cErr != nil {
		return "", cErr
	}
	orderDetails, patientDetails, cErr := attuneService.AttuneDao.GetOrderDetailsAndPatientDetailsByOmsOrderId(omsOrderId)
	if cErr != nil {
		return "", cErr
	}

	attunePayloadMeta.SrfId = orderDetails.SrfId
	attunePayloadMeta.LabId = labId

	// Map test and barcodes in the required structure
	attunePayloadMeta.TestBarcodeMap = attuneService.GetTestBarcodeMap(ctx, attunePayloadMeta.AttuneTestSampleMap)

	// Get Test Details List for Attune
	attunePayloadMeta.AttuneTestDetailsList, cErr = attuneService.GetTestDetailsListForSyncingToAttune(omsOrderId,
		sampleIds, labId, attunePayloadMeta.TestBarcodeMap)
	if cErr != nil {
		return "", cErr
	}

	attunePayloadMeta.PatientSalutation = commonUtils.GetSalutationByGender(patientDetails.Gender)
	attunePayloadMeta.PatientName = patientDetails.Name
	attunePayloadMeta.PatientGender = commonUtils.ToTitleCase(patientDetails.Gender)
	attunePayloadMeta.PatientNumber = commonConstants.AttunePatientDefaultContact
	attunePayloadMeta.PatientId = commonUtils.GetPatientIdByOrderIdAndPatientId(orderDetails.OmsOrderId,
		patientDetails.SystemPatientId)

	loc, _ := time.LoadLocation(commonConstants.LocalTimeZoneLocation)
	attunePayloadMeta.PatientDob = getPatientFormattedDob(patientDetails)

	// Fetch Visit Id
	attunePayloadMeta.VisitId = attuneService.GetVisitId(ctx, omsOrderId)

	// Get Collection time
	fomattedCollectedAt := ""
	collectedAt, cErr := attuneService.AttuneDao.GetSampleCollectedAtBySampleIds(sampleIds)
	if cErr != nil {
		return "", cErr
	}
	if collectedAt != nil {
		fomattedCollectedAt = (*collectedAt).In(loc).Format(commonConstants.DateTimeInSecLayout)
	} else {
		fomattedCollectedAt = time.Now().In(loc).Format(commonConstants.DateTimeInSecLayout)
	}
	attunePayloadMeta.SampleCollectedAt = fomattedCollectedAt

	attunePayload, cErr := getAttunePayloadForSyncingData(attunePayloadMeta)
	if cErr != nil {
		return "", cErr
	}

	cErr = attuneService.AttuneClient.SyncDataToAttune(ctx, attunePayloadMeta.LabId, attunePayload)
	if cErr != nil {
		attuneService.SendSlackMessageIfLisSyncFailed(ctx, attunePayloadMeta.AttuneTestSampleMap,
			attunePayload.TestDetailsList, patientDetails.Name, orderDetails)
		return "", cErr
	}

	return attunePayloadMeta.VisitId, nil
}

func (attuneService *AttuneService) CancelLisSyncData(ctx context.Context, testDetails []commonModels.TestDetail,
	visitId string, sampleInfo commonStructures.SampleInfo) *commonStructures.CommonError {

	orderDetails, patientDetails, cErr := attuneService.AttuneDao.GetOrderDetailsAndPatientDetailsByVisitId(visitId)
	if cErr != nil {
		return cErr
	}

	gender := commonUtils.ToTitleCase(patientDetails.Gender)
	salutation := commonUtils.GetSalutationByGender(patientDetails.Gender)
	patientNumber := commonConstants.AttunePatientDefaultContact
	patientId := commonUtils.GetPatientIdByOrderIdAndPatientId(orderDetails.OmsOrderId, patientDetails.SystemPatientId)

	// Calculating date of birth
	loc, _ := time.LoadLocation(commonConstants.LocalTimeZoneLocation)
	formattedDob := getPatientFormattedDob(patientDetails)

	// Get sampleDetails of the accession which has to be rejected
	sampleDetails, cErr := attuneService.GetAttuneSampleDetailsToCancel(ctx, sampleInfo)
	if cErr != nil {
		return cErr
	}

	// Get collection time for the the sample
	fomattedCollectedAt := ""
	collectedAt, cErr := attuneService.AttuneDao.GetSampleCollectedAtByVisitId(visitId)
	if cErr != nil {
		return cErr
	}
	if collectedAt != nil {
		fomattedCollectedAt = (*collectedAt).In(loc).Format(commonConstants.DateTimeInSecLayout)
	} else {
		fomattedCollectedAt = time.Now().In(loc).Format(commonConstants.DateTimeInSecLayout)
	}

	attuneTestDetailsList := []commonStructures.AttuneTestDetails{}
	for _, test := range testDetails {
		attuneTestDetail := commonStructures.AttuneTestDetails{
			Status:   commonConstants.AttuneTestStatusCancel,
			TestCode: test.LisCode,
			TestName: test.TestName,
			TestType: test.TestType,
			TestID:   int(test.Id),
			BarcodeNo: []commonStructures.AttuneSampleDetails{
				sampleDetails,
			},
			Price: commonConstants.DefaultPriceAttuneTest,
		}

		attuneTestDetailsList = append(attuneTestDetailsList, attuneTestDetail)
	}

	attuneOrgCode := commonUtils.GetAttuneOrgCodeByLabId(sampleInfo.LabId)

	payload := commonStructures.AttuneSyncDataToLisRequest{
		MessageType: commonConstants.AttuneMessageTypeCancel,
		OrderID:     visitId,
		OrgCode:     attuneOrgCode,
		PatientInfo: commonStructures.AttunePatientInfo{
			Salutation:    salutation,
			PatientId:     patientId,
			PatientNumber: patientId,
			Name:          patientDetails.Name, // PII data. Ideally send anonymous data: OrangePatient_Id
			Gender:        gender,
			Dob:           formattedDob,
			MobileNumber:  patientNumber,
		},
		PatientVisitInfo: commonStructures.AttunePatientVisitInfo{
			VisitDate:      time.Now().In(loc).Format(commonConstants.DateTimeInSecLayout),
			CollectedDate:  fomattedCollectedAt,
			ClientName:     commonConstants.AttuneClientName,
			ReportLanguage: nil,
		},
		TestDetailsList: attuneTestDetailsList,
	}

	cErr = attuneService.AttuneClient.SyncDataToAttune(ctx, sampleInfo.LabId, payload)
	if cErr != nil {
		return cErr
	}

	return nil
}

func (attuneService *AttuneService) ModifyLisDataPostSyncByOrderId(ctx context.Context, omsOrderId string, labId uint,
	testDetail commonModels.TestDetail, sample commonStructures.SampleInfo) *commonStructures.CommonError {

	sampleModel := commonModels.Sample{}
	sampleModel.Id = sample.Id
	sampleIdString := getSampleIdString([]commonModels.Sample{sampleModel})
	redisKey := fmt.Sprintf(commonConstants.CacheKeyCurrentInProgressSample, omsOrderId)
	value := ""
	err := attuneService.Cache.Get(ctx, redisKey, &value)
	if err == redis.Nil {
		value = sampleIdString
		_ = attuneService.Cache.Set(ctx, redisKey, value, commonConstants.CacheExpiry10DaysInt)

		defer func() {
			_ = attuneService.Cache.Delete(ctx, redisKey)
		}()
	} else if value == sampleIdString {
		return &commonStructures.CommonError{
			Message:    commonConstants.ERROR_LIS_SYNC_ALREADY_IN_PROGRESS,
			StatusCode: http.StatusConflict,
		}
	}

	visitId := sample.VisitId

	testSampleMap := []commonStructures.AttuneTestSampleMap{{
		TestId:     testDetail.CentralOmsTestId,
		VialTypeId: sample.VialTypeId,
		Barcode:    sample.Barcode,
	}}

	orderDetails, patientDetails, cErr := attuneService.AttuneDao.GetOrderDetailsAndPatientDetailsByOmsOrderId(omsOrderId)
	if cErr != nil {
		return cErr
	}

	testBarcodeMap := attuneService.GetTestBarcodeMap(ctx, testSampleMap)

	testDetailsList := []commonStructures.AttuneTestDetails{{
		Status:    commonConstants.AttuneTestStatusOrdered,
		TestCode:  testDetail.LisCode,
		TestName:  testDetail.TestName,
		TestType:  testDetail.TestType,
		TestID:    int(testDetail.Id),
		BarcodeNo: testBarcodeMap[testDetail.CentralOmsTestId],
		Price:     commonConstants.DefaultPriceAttuneTest,
	}}
	orgCode := commonUtils.GetAttuneOrgCodeByLabId(labId)

	attunePayload := commonStructures.AttuneSyncDataToLisAfterSyncRequest{
		MessageType:     commonConstants.AttuneMessageTypeModified,
		OrderID:         visitId,
		OrgCode:         orgCode,
		TestDetailsList: testDetailsList,
	}

	cErr = attuneService.AttuneClient.SyncDataToAttuneAfterSync(ctx, labId, attunePayload)
	if cErr != nil {
		attuneService.SendSlackMessageIfLisSyncFailed(ctx, testSampleMap, testDetailsList, patientDetails.Name, orderDetails)
		return cErr
	}

	return nil
}

func (attuneService *AttuneService) UpdateSrfIdToAttune(ctx context.Context, sample commonStructures.SampleInfo,
	srfId string) *commonStructures.CommonError {
	if sample.VisitId == "" {
		return &commonStructures.CommonError{
			Message:    commonConstants.ERROR_VISIT_ID_NOT_FOUND,
			StatusCode: http.StatusBadRequest,
		}
	}
	_, attuneResponse, cErr := attuneService.GetLisSyncData(ctx, sample.VisitId, commonConstants.AttuneReportWithStationery)
	if cErr != nil {
		return cErr
	}
	payload := commonStructures.AttuneSyncDataToLisRequest{
		MessageType:      commonConstants.AttuneMessageTypeModified,
		OrderID:          attuneResponse.OrderId,
		OrgCode:          attuneResponse.OrgCode,
		PatientInfo:      attuneResponse.PatientInfo,
		PatientVisitInfo: attuneResponse.PatientVisitInfo,
	}

	// if attune doesn't send the collected date
	if payload.PatientVisitInfo.CollectedDate == "" {
		loc, _ := time.LoadLocation(commonConstants.LocalTimeZoneLocation)
		fomattedCollectedAt := ""
		collectedAt, cErr := attuneService.AttuneDao.GetSampleCollectedAtByVisitId(sample.VisitId)
		if cErr != nil {
			return cErr
		}
		if collectedAt != nil {
			fomattedCollectedAt = (*collectedAt).In(loc).Format(commonConstants.DateTimeInSecLayout)
		} else {
			fomattedCollectedAt = time.Now().In(loc).Format(commonConstants.DateTimeInSecLayout)
		}
		payload.PatientVisitInfo.CollectedDate = fomattedCollectedAt
	}

	payload.PatientVisitInfo.SrfId = srfId
	cErr = attuneService.AttuneClient.SyncDataToAttune(ctx, sample.LabId, payload)
	if cErr != nil {
		return cErr
	}

	return nil
}
