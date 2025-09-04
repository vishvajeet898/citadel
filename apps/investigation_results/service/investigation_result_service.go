package service

import (
	"context"
	"fmt"
	"net/http"
	"sort"
	"sync"

	"gorm.io/gorm"

	mapper "github.com/Orange-Health/citadel/apps/investigation_results/mapper"
	"github.com/Orange-Health/citadel/apps/investigation_results/structures"
	commonConstants "github.com/Orange-Health/citadel/common/constants"
	commonStructures "github.com/Orange-Health/citadel/common/structures"
	commonUtils "github.com/Orange-Health/citadel/common/utils"
	commonModels "github.com/Orange-Health/citadel/models"
)

type InvestigationResultServiceInterface interface {
	GetPatientPastRecords(ctx context.Context, patientId string) (
		commonStructures.PatientPastRecordsApiResponse, *commonStructures.CommonError)
	GetInvestigationResultByTaskId(taskId uint) (
		[]structures.InvestigationResultDetails, *commonStructures.CommonError)
	GetInvestigationResultModelsByTaskId(taskId uint) (
		[]commonModels.InvestigationResult, *commonStructures.CommonError)
	GetInvestigationsByInvestigationIds(investigationIds []uint) (
		[]commonModels.InvestigationResult, *commonStructures.CommonError)
	GetInvestigationDataByInvestigationResultsIds(investigationResultsIds []uint) (
		[]commonModels.InvestigationData, *commonStructures.CommonError)
	GetDeltaValuesFromPatientId(ctx context.Context, deltaVal structures.DeltaValuesRequest) (
		map[uint][]commonStructures.DeltaValuesStructResponse, *commonStructures.CommonError)
	GetLastInvValueByPatientId(ctx context.Context, patientId string,
		masterInvestigationIds []uint) (map[uint]commonStructures.DeltaValuesStructResponse,
		*commonStructures.CommonError)
	GetInvestigationResultsByTestDetailsIds(testDetailsIds []uint) (
		[]commonModels.InvestigationResult, *commonStructures.CommonError)
	GetInvestigationAbnormality(ctx context.Context, taskId uint, investigationCode string,
		investigationValue string) (string, *commonStructures.CommonError)
	GetDerivedInvestigationsAndAbnormality(ctx context.Context, taskId uint, masterInvestigationId uint,
		investigationValue string, modifyValueRequest []structures.ModifyValueRequest) (
		map[uint]structures.ModifyValueResponse, *commonStructures.CommonError)
	GetInvestigationResultsByTaskIdAndOmsTestId(taskId uint, omsTestId string) (
		[]commonModels.InvestigationResult, *commonStructures.CommonError)
	GetInvestigationResultsMetadataByInvestigationResultIds(ctx context.Context,
		investigationResultsIds []uint) (map[uint]commonModels.InvestigationResultMetadata, *commonStructures.CommonError)

	CreateInvestigationResultsWithTx(tx *gorm.DB,
		investigationResults []commonModels.InvestigationResult) (
		[]commonModels.InvestigationResult, *commonStructures.CommonError)
	UpdateInvestigationResultsWithTx(tx *gorm.DB,
		investigationResults []commonModels.InvestigationResult) (
		[]commonModels.InvestigationResult, *commonStructures.CommonError)
	DeleteInvestigationResultsByIdsWithTx(tx *gorm.DB,
		investigationIds []uint) *commonStructures.CommonError
	UpdateInvestigationsDataWithTx(tx *gorm.DB,
		investigationData []commonModels.InvestigationData) (
		[]commonModels.InvestigationData, *commonStructures.CommonError)

	CreateInvestigationDataWithTx(tx *gorm.DB,
		investigationData []commonModels.InvestigationData) (
		[]commonModels.InvestigationData, *commonStructures.CommonError)
	UpdateInvestigationDataWithTx(tx *gorm.DB,
		investigationData []commonModels.InvestigationData) (
		[]commonModels.InvestigationData, *commonStructures.CommonError)

	CreateInvestigationResultsMetadataWithTx(ctx context.Context, tx *gorm.DB,
		investigationResultsMetadata []commonModels.InvestigationResultMetadata) (
		[]commonModels.InvestigationResultMetadata, *commonStructures.CommonError)
	UpdateInvestigationResultsMetadataWithTx(ctx context.Context, tx *gorm.DB,
		investigationResultsMetadata []commonModels.InvestigationResultMetadata, existingMetadataMap map[uint]commonModels.InvestigationResultMetadata) (
		[]commonModels.InvestigationResultMetadata, *commonStructures.CommonError)
	DeleteInvestigationResultsMetadataByIdsWithTx(ctx context.Context, tx *gorm.DB,
		investigationResultsMetadataIds []uint) *commonStructures.CommonError
}

func filterDuplicateAndEmptyPatientPastRecords(
	omsPatientPastRecords []commonStructures.PatientPastRecords,
	citadelPatientPastRecords []commonStructures.PatientPastRecords) []commonStructures.PatientPastRecords {
	patientPastRecordsMap := make(map[string]bool)
	patientPastRecords := []commonStructures.PatientPastRecords{}

	for _, pastRecord := range citadelPatientPastRecords {
		if pastRecord.InvestigationValue != "" {
			patientPastRecords = append(patientPastRecords, pastRecord)
			key := fmt.Sprint(pastRecord.OrderId, pastRecord.CityCode, pastRecord.MasterInvestigationId)
			patientPastRecordsMap[key] = true
		}
	}

	for _, pastRecord := range omsPatientPastRecords {
		if pastRecord.InvestigationValue != "" {
			key := fmt.Sprint(pastRecord.OrderId, pastRecord.CityCode, pastRecord.MasterInvestigationId)
			if _, ok := patientPastRecordsMap[key]; !ok {
				patientPastRecords = append(patientPastRecords, pastRecord)
			}
		}
	}

	return patientPastRecords
}

func filterDuplicateDeltaValues(
	omsDeltaValues []commonStructures.DeltaValuesStruct,
	citadelDeltaValues []commonStructures.DeltaValuesStruct) []commonStructures.DeltaValuesStruct {
	citadelDeltaValuesMap := make(map[string]bool)

	deltaValuesStructs := citadelDeltaValues

	for _, deltaValue := range citadelDeltaValues {
		key := fmt.Sprint(deltaValue.OrderId, deltaValue.CityCode, deltaValue.MasterInvestigationId)
		citadelDeltaValuesMap[key] = true
	}

	for _, deltaValue := range omsDeltaValues {
		key := fmt.Sprint(deltaValue.OrderId, deltaValue.CityCode, deltaValue.MasterInvestigationId)
		if _, ok := citadelDeltaValuesMap[key]; !ok {
			deltaValuesStructs = append(deltaValuesStructs, deltaValue)
		}
	}

	return deltaValuesStructs
}

func filterDeltaValuesWithoutApprovedAtAndInvestigationValue(
	deltaValues []commonStructures.DeltaValuesStruct) []commonStructures.DeltaValuesStruct {
	filteredDeltaValues := []commonStructures.DeltaValuesStruct{}
	for _, deltaValue := range deltaValues {
		if deltaValue.ApprovedAt != nil && deltaValue.InvestigationValue != "" {
			filteredDeltaValues = append(filteredDeltaValues, deltaValue)
		}
	}
	return filteredDeltaValues
}

func sortDeltaValuesByApprovedAtAndApplyLimit(
	deltaValues []commonStructures.DeltaValuesStructResponse, limit uint) []commonStructures.DeltaValuesStructResponse {
	sort.Slice(deltaValues, func(i, j int) bool {
		return deltaValues[i].ApprovedAt.After(*deltaValues[j].ApprovedAt)
	})

	if limit > 0 && uint(len(deltaValues)) > limit {
		deltaValues = deltaValues[:limit]
	}

	return deltaValues
}

func (invResService *InvestigationResultService) GetPatientPastRecords(ctx context.Context, patientId string) (
	commonStructures.PatientPastRecordsApiResponse, *commonStructures.CommonError) {
	patientDetails, patientPastRecords := invResService.getCitadelOMSPastRecords(ctx, patientId)
	return mapper.CreateResponseForPatientPastRecords(patientDetails, patientPastRecords), nil
}

func (invResService *InvestigationResultService) getCitadelOMSPastRecords(ctx context.Context, patientId string) (
	commonStructures.PatientDetailsResponse, []commonStructures.PatientPastRecords) {
	similarPatientDetails, err := invResService.PatientServiceClient.GetSimilarPatientDetails(ctx, patientId)
	omsPatientPastRecords, citadelPatientPastRecords :=
		[]commonStructures.PatientPastRecords{}, []commonStructures.PatientPastRecords{}
	patientDetails := commonStructures.PatientDetailsResponse{}
	if err != nil {
		commonUtils.AddLog(ctx, commonConstants.ERROR_LEVEL, commonConstants.ERROR_WHILE_GETTING_SIMILAR_PATIENT_DETAILS, nil, err)
	}

	patientIds := []string{patientId}
	for _, similarPatient := range similarPatientDetails {
		patientIds = append(patientIds, similarPatient.PatientID)
	}

	patientIds = commonUtils.CreateUniqueSliceString(patientIds)

	wg := sync.WaitGroup{}
	wg.Add(3)

	go func() {
		defer wg.Done()
		patientDetails, _ = invResService.PatientServiceClient.GetPatientDetailsByPatientId(ctx, patientId)
	}()

	go func() {
		defer wg.Done()
		omsPatientPastRecords, _ = invResService.OmsClient.GetPatientPastRecordsFromPatientIds(ctx, patientIds)
	}()

	go func() {
		defer wg.Done()
		citadelPatientPastRecords, _ = invResService.InvResDao.GetPatientPastRecords(patientIds)
	}()

	wg.Wait()

	patientPastRecords := filterDuplicateAndEmptyPatientPastRecords(omsPatientPastRecords, citadelPatientPastRecords)
	return patientDetails, patientPastRecords
}

func (invResService *InvestigationResultService) GetInvestigationResultByTaskId(taskId uint) (
	[]structures.InvestigationResultDetails, *commonStructures.CommonError) {

	investigationResults := []structures.InvestigationResultDbResponse{}
	investigationsDataMap := map[uint]commonModels.InvestigationData{}
	rerunDetails, remarks := []commonModels.RerunInvestigationResult{}, []commonModels.Remark{}

	// Fetching Test Details
	testDetails, cErr := invResService.InvResDao.GetTestDetailsByTaskId(taskId)
	if cErr != nil {
		return nil, cErr
	}

	// If no test details found, return empty result
	if len(testDetails) == 0 {
		return []structures.InvestigationResultDetails{}, nil
	}

	testDetailIds := mapper.MapTestDetailsIdsFromTestDetailsStructs(testDetails)
	var errList []*commonStructures.CommonError

	wg, mu := sync.WaitGroup{}, sync.Mutex{}
	wg.Add(2)

	// Fetching Investigation Results and Remarks
	go func(testDetailIds []uint) {
		defer wg.Done()

		var cErr *commonStructures.CommonError

		// Fetching Investigation Results
		investigationResults, cErr = invResService.InvResDao.GetInvestigationResultsDbStructsByTestDetailIds(testDetailIds)
		if cErr != nil {
			mu.Lock()
			errList = append(errList, cErr)
			mu.Unlock()
			return
		}

		if len(investigationResults) == 0 {
			return
		}

		investigationResultIds := mapper.MapInvestigationIdsFromInvestigationResultsDbStructs(investigationResults)

		wg1 := sync.WaitGroup{}
		wg1.Add(2)

		// Fetching Investigation Data
		go func(investigationResultIds []uint) {
			defer wg1.Done()
			investigationsData, cErr := invResService.InvResDao.GetInvestigationDataByInvestigationResultsIds(
				investigationResultIds)
			if cErr != nil {
				mu.Lock()
				errList = append(errList, cErr)
				mu.Unlock()
				return
			}
			investigationsDataMap = mapper.MapInvestigationData(investigationsData)
		}(investigationResultIds)

		// Fetching Remarks
		go func(investigationResultIds []uint) {
			defer wg1.Done()
			remarks, cErr = invResService.RemarkService.GetRemarksByInvestigationResultIds(nil, investigationResultIds)
			if cErr != nil {
				mu.Lock()
				errList = append(errList, cErr)
				mu.Unlock()
				return
			}
		}(investigationResultIds)

		wg1.Wait()
	}(testDetailIds)

	// Fetching Rerun Details
	go func(testDetailIds []uint) {
		defer wg.Done()
		var cErr *commonStructures.CommonError
		rerunDetails, cErr = invResService.RerunService.GetRerunDetailsByTestDetailsIds(testDetailIds)
		if cErr != nil {
			mu.Lock()
			errList = append(errList, cErr)
			mu.Unlock()
		}
	}(testDetailIds)

	wg.Wait()

	if len(errList) != 0 {
		return []structures.InvestigationResultDetails{}, errList[0]
	}

	// Map Investigation Results with Master Tests and filter the tests whose cpEnabled flag is not enabled
	investigationResults = mapper.MapInvestigationResultsDbStructsWithMasterTests(investigationResults, testDetails)

	// Group remarks and rerun details by investigation result ID
	return mapper.CreateInvestigationResultsResponse(
		mapper.MapInvestigationResultsFromDbStruct(investigationResults, investigationsDataMap),
		mapper.MapRemarks(remarks), mapper.MapRerunDetails(rerunDetails)), nil
}

func (invResService *InvestigationResultService) GetInvestigationResultModelsByTaskId(taskId uint) (
	[]commonModels.InvestigationResult, *commonStructures.CommonError) {

	return invResService.InvResDao.GetInvestigationResultModelsByTaskId(taskId)
}

func (invResService *InvestigationResultService) GetInvestigationsByInvestigationIds(
	investigationIds []uint) ([]commonModels.InvestigationResult, *commonStructures.CommonError) {

	return invResService.InvResDao.GetInvestigationsByInvestigationIds(investigationIds)
}

func (invResService *InvestigationResultService) GetDeltaValuesFromPatientId(ctx context.Context, deltaVal structures.DeltaValuesRequest) (
	map[uint][]commonStructures.DeltaValuesStructResponse, *commonStructures.CommonError) {
	omsDeltaValues, citadelDeltaValues := []commonStructures.DeltaValuesStruct{}, []commonStructures.DeltaValuesStruct{}

	similarPatientDetails, err := invResService.PatientServiceClient.GetSimilarPatientDetails(ctx, deltaVal.PatientId)
	if err != nil {
		commonUtils.AddLog(ctx, commonConstants.ERROR_LEVEL, commonConstants.ERROR_WHILE_GETTING_SIMILAR_PATIENT_DETAILS, nil, err)
	}

	patientIds := []string{deltaVal.PatientId}
	for _, similarPatient := range similarPatientDetails {
		patientIds = append(patientIds, similarPatient.PatientID)
	}

	patientIds = commonUtils.CreateUniqueSliceString(patientIds)
	deltaVal.MasterInvestigationIds = commonUtils.CreateUniqueSliceUint(deltaVal.MasterInvestigationIds)

	wg := sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()
		omsDeltaValues, _ = invResService.OmsClient.GetDeltaValuesFromPatientIds(ctx, patientIds,
			deltaVal.MasterInvestigationIds)
	}()

	go func() {
		defer wg.Done()
		citadelDeltaValues, _ = invResService.InvResDao.GetDeltaValuesFromPatientIds(patientIds,
			deltaVal.MasterInvestigationIds)
	}()

	wg.Wait()

	deltaValuesStructs := filterDuplicateDeltaValues(omsDeltaValues, citadelDeltaValues)

	deltaValuesStructs = filterDeltaValuesWithoutApprovedAtAndInvestigationValue(deltaValuesStructs)

	masterInvestigationIdToDeltaValues := make(map[uint][]commonStructures.DeltaValuesStructResponse)
	for _, deltaValuesStruct := range deltaValuesStructs {
		masterInvestigationIdToDeltaValues[deltaValuesStruct.MasterInvestigationId] = append(
			masterInvestigationIdToDeltaValues[deltaValuesStruct.MasterInvestigationId], commonStructures.DeltaValuesStructResponse{
				InvestigationValue:    deltaValuesStruct.InvestigationValue,
				MasterInvestigationId: deltaValuesStruct.MasterInvestigationId,
				ApprovedAt:            deltaValuesStruct.ApprovedAt,
			})
	}

	for index, deltaValues := range masterInvestigationIdToDeltaValues {
		deltaValues = sortDeltaValuesByApprovedAtAndApplyLimit(deltaValues, deltaVal.Limit)
		masterInvestigationIdToDeltaValues[index] = deltaValues
	}

	return masterInvestigationIdToDeltaValues, nil
}

func (invResService *InvestigationResultService) GetLastInvValueByPatientId(ctx context.Context, patientId string,
	masterInvestigationIds []uint) (map[uint]commonStructures.DeltaValuesStructResponse, *commonStructures.CommonError) {
	similarPatientDetails, err := invResService.PatientServiceClient.GetSimilarPatientDetails(ctx, patientId)
	if err != nil {
		commonUtils.AddLog(ctx, commonConstants.ERROR_LEVEL, commonConstants.ERROR_WHILE_GETTING_SIMILAR_PATIENT_DETAILS, nil, err)
	}
	patientIds := []string{patientId}
	for _, similarPatient := range similarPatientDetails {
		patientIds = append(patientIds, similarPatient.PatientID)
	}
	patientIds = commonUtils.CreateUniqueSliceString(patientIds)

	deltaValues, cErr := invResService.InvResDao.GetDeltaValuesFromPatientIds(patientIds, masterInvestigationIds)
	if cErr != nil {
		return nil, cErr
	}
	deltaValues = filterDeltaValuesWithoutApprovedAtAndInvestigationValue(deltaValues)

	masterInvestigationIdToDeltaValues := make(map[uint][]commonStructures.DeltaValuesStructResponse)
	for _, deltaValuesStruct := range deltaValues {
		masterInvestigationIdToDeltaValues[deltaValuesStruct.MasterInvestigationId] = append(
			masterInvestigationIdToDeltaValues[deltaValuesStruct.MasterInvestigationId], commonStructures.DeltaValuesStructResponse{
				InvestigationValue:    deltaValuesStruct.InvestigationValue,
				MasterInvestigationId: deltaValuesStruct.MasterInvestigationId,
				ApprovedAt:            deltaValuesStruct.ApprovedAt,
			})
	}
	investigationValueMap := make(map[uint]commonStructures.DeltaValuesStructResponse)
	for id, deltaValues := range masterInvestigationIdToDeltaValues {
		deltaValues = sortDeltaValuesByApprovedAtAndApplyLimit(deltaValues, 1)
		investigationValueMap[id] = deltaValues[0]
	}

	return investigationValueMap, nil
}

func (invResService *InvestigationResultService) GetInvestigationResultsByTestDetailsIds(testDetailsIds []uint) (
	[]commonModels.InvestigationResult, *commonStructures.CommonError) {

	return invResService.InvResDao.GetInvestigationResultsByTestDetailsIds(testDetailsIds)
}

func (invResService *InvestigationResultService) GetInvestigationResultsMetadataByInvestigationResultIds(ctx context.Context,
	investigationResultsIds []uint) (map[uint]commonModels.InvestigationResultMetadata, *commonStructures.CommonError) {

	irms, cErr := invResService.InvResDao.GetInvestigationResultsMetadataByInvestigationResultsIds(investigationResultsIds)
	if cErr != nil {
		return nil, cErr
	}
	irmsMap := make(map[uint]commonModels.InvestigationResultMetadata)
	for _, irm := range irms {
		irmsMap[irm.InvestigationResultId] = irm
	}
	return irmsMap, nil
}

func (invResService *InvestigationResultService) CreateInvestigationResultsWithTx(tx *gorm.DB,
	investigationResults []commonModels.InvestigationResult) (
	[]commonModels.InvestigationResult, *commonStructures.CommonError) {

	if len(investigationResults) == 0 {
		return investigationResults, nil
	}

	return invResService.InvResDao.CreateInvestigationResultsWithTx(tx, investigationResults)
}

func (invResService *InvestigationResultService) UpdateInvestigationResultsWithTx(tx *gorm.DB,
	investigationResults []commonModels.InvestigationResult) (
	[]commonModels.InvestigationResult, *commonStructures.CommonError) {

	if len(investigationResults) == 0 {
		return investigationResults, nil
	}

	return invResService.InvResDao.UpdateInvestigationResultsWithTx(tx, investigationResults)
}

func (invResService *InvestigationResultService) DeleteInvestigationResultsByIdsWithTx(tx *gorm.DB,
	investigationIds []uint) *commonStructures.CommonError {

	if len(investigationIds) == 0 {
		return nil
	}

	return invResService.InvResDao.DeleteInvestigationResultsByIdsWithTx(tx, investigationIds)
}

func (invResService *InvestigationResultService) UpdateInvestigationsDataWithTx(tx *gorm.DB,
	investigationData []commonModels.InvestigationData) (
	[]commonModels.InvestigationData, *commonStructures.CommonError) {

	if len(investigationData) == 0 {
		return investigationData, nil
	}

	return invResService.InvResDao.UpdateInvestigationsDataWithTx(tx, investigationData)
}

func (invResService *InvestigationResultService) GetInvestigationAbnormality(ctx context.Context, taskId uint,
	investigationCode string, investigationValue string) (string, *commonStructures.CommonError) {

	abnormalityDetails, cErr := invResService.InvResDao.GetDetailsForAbnormality(taskId)
	if cErr != nil {
		return "", cErr
	}

	patientDob := ""
	if abnormalityDetails.PatientExpectedDob != nil {
		patientDob = abnormalityDetails.PatientExpectedDob.Format(commonConstants.DateLayout)
	}
	if abnormalityDetails.PatientDob != nil {
		patientDob = abnormalityDetails.PatientDob.Format(commonConstants.DateLayout)
	}

	patientGender := commonUtils.GetGenderConstant(abnormalityDetails.PatientGender)

	investigations, err := invResService.CdsClient.GetInvestigationDetails(ctx, []string{investigationCode},
		abnormalityDetails.CityCode, abnormalityDetails.LabId, patientDob, patientGender)
	if err != nil {
		return "", &commonStructures.CommonError{
			Message:    err.Error(),
			StatusCode: http.StatusInternalServerError,
		}
	}
	investigation := investigations[0]

	return invResService.AbnormalityService.GetInvestigationAbnormality(ctx, investigationValue, investigation), nil
}

func (invResService *InvestigationResultService) GetDerivedInvestigationsAndAbnormality(ctx context.Context, taskId uint,
	masterInvestigationId uint, investigationValue string, modifyValueRequest []structures.ModifyValueRequest) (
	map[uint]structures.ModifyValueResponse, *commonStructures.CommonError) {

	abnormalityDetails, cErr := invResService.InvResDao.GetDetailsForAbnormality(taskId)
	if cErr != nil {
		return nil, cErr
	}

	patientDob := ""
	if abnormalityDetails.PatientExpectedDob != nil {
		patientDob = abnormalityDetails.PatientExpectedDob.Format(commonConstants.DateLayout)
	}
	if abnormalityDetails.PatientDob != nil {
		patientDob = abnormalityDetails.PatientDob.Format(commonConstants.DateLayout)
	}

	patientGender := commonUtils.GetGenderConstant(abnormalityDetails.PatientGender)

	dependentInvestigationsApiResponse, err := invResService.CdsClient.GetDependentInvestigations(ctx,
		masterInvestigationId, abnormalityDetails.LabId, abnormalityDetails.CityCode, patientDob, patientGender)
	if err != nil {
		return nil, &commonStructures.CommonError{
			Message:    commonConstants.TECHNICAL_ERROR,
			StatusCode: http.StatusInternalServerError,
		}
	}

	investigationResults, cErr := invResService.InvResDao.GetInvestigationResultModelsByTaskId(taskId)
	if cErr != nil {
		return nil, cErr
	}

	masterInvestigationIdToNewInvestigationValueMap := make(map[uint]string)
	for _, modifyValueRequest := range modifyValueRequest {
		masterInvestigationIdToNewInvestigationValueMap[modifyValueRequest.MasterInvestigationId] =
			modifyValueRequest.InvestigationValue
	}

	masterInvestigationIdToInvestigationResultMap :=
		mapper.MapInvestigationResultsByMasterInvestigationId(investigationResults)

	for masterInvestigationId, newInvestigationValue := range masterInvestigationIdToNewInvestigationValueMap {
		if investigationResult, ok := masterInvestigationIdToInvestigationResultMap[masterInvestigationId]; ok {
			investigationResult.InvestigationValue = newInvestigationValue
			masterInvestigationIdToInvestigationResultMap[masterInvestigationId] = investigationResult
		}
	}

	invRes := masterInvestigationIdToInvestigationResultMap[masterInvestigationId]
	invRes.InvestigationValue = investigationValue
	masterInvestigationIdToInvestigationResultMap[masterInvestigationId] = invRes

	response := map[uint]structures.ModifyValueResponse{}
	response[dependentInvestigationsApiResponse.Investigation.InvestigationId] = structures.ModifyValueResponse{
		Value: investigationValue,
		Abnormality: invResService.AbnormalityService.GetInvestigationAbnormalityByMasterReferenceRange(ctx,
			investigationValue, dependentInvestigationsApiResponse.Investigation.ReferenceRange),
	}

	for _, derivedInvestigation := range dependentInvestigationsApiResponse.DerivedInvestigations {
		if investigationResult, ok :=
			masterInvestigationIdToInvestigationResultMap[derivedInvestigation.InvestigationId]; ok {

			modifiedValue, cErr := invResService.CalculationsService.GetModifiedValue(ctx, derivedInvestigation.Formula,
				derivedInvestigation.Decimals, derivedInvestigation.DependentInvestigationIds,
				masterInvestigationIdToInvestigationResultMap)
			if cErr != nil {
				return nil, cErr
			}
			response[derivedInvestigation.InvestigationId] = structures.ModifyValueResponse{
				Value: modifiedValue,
				Abnormality: invResService.AbnormalityService.GetInvestigationAbnormalityByMasterReferenceRange(ctx,
					investigationResult.InvestigationValue, derivedInvestigation.ReferenceRange),
			}
		}
	}

	return response, nil
}

func (invResService *InvestigationResultService) GetInvestigationResultsByTaskIdAndOmsTestId(
	taskId uint, omsTestId string) ([]commonModels.InvestigationResult, *commonStructures.CommonError) {

	return invResService.InvResDao.GetInvestigationResultsByTaskIdAndOmsTestId(taskId, omsTestId)
}

func (invResService *InvestigationResultService) GetInvestigationDataByInvestigationResultsIds(
	investigationResultsIds []uint) ([]commonModels.InvestigationData, *commonStructures.CommonError) {
	if len(investigationResultsIds) == 0 {
		return []commonModels.InvestigationData{}, nil
	}

	return invResService.InvResDao.GetInvestigationDataByInvestigationResultsIds(investigationResultsIds)
}

func (invResService *InvestigationResultService) CreateInvestigationDataWithTx(tx *gorm.DB,
	investigationData []commonModels.InvestigationData) (
	[]commonModels.InvestigationData, *commonStructures.CommonError) {

	if len(investigationData) == 0 {
		return investigationData, nil
	}

	return invResService.InvResDao.CreateInvestigationDataWithTx(tx, investigationData)
}

func (invResService *InvestigationResultService) UpdateInvestigationDataWithTx(tx *gorm.DB,
	investigationData []commonModels.InvestigationData) (
	[]commonModels.InvestigationData, *commonStructures.CommonError) {

	if len(investigationData) == 0 {
		return investigationData, nil
	}

	return invResService.InvResDao.UpdateInvestigationsDataWithTx(tx, investigationData)
}

func (invResService *InvestigationResultService) CreateInvestigationResultsMetadataWithTx(ctx context.Context, tx *gorm.DB,
	investigationResultsMetadata []commonModels.InvestigationResultMetadata) (
	[]commonModels.InvestigationResultMetadata, *commonStructures.CommonError) {

	if len(investigationResultsMetadata) == 0 {
		return investigationResultsMetadata, nil
	}

	return invResService.InvResDao.CreateInvestigationResultsMetadataWithTx(tx, investigationResultsMetadata)
}

func (invResService *InvestigationResultService) UpdateInvestigationResultsMetadataWithTx(ctx context.Context, tx *gorm.DB,
	investigationResultsMetadata []commonModels.InvestigationResultMetadata, existingMetadataMap map[uint]commonModels.InvestigationResultMetadata) (
	[]commonModels.InvestigationResultMetadata, *commonStructures.CommonError) {

	if len(investigationResultsMetadata) == 0 {
		return investigationResultsMetadata, nil
	}
	for idx, irm := range investigationResultsMetadata {
		if existingMetadata, ok := existingMetadataMap[irm.InvestigationResultId]; ok {
			irm.Id = existingMetadata.Id
		}
		investigationResultsMetadata[idx] = irm
	}
	return invResService.InvResDao.UpdateInvestigationResultsMetadataWithTx(tx, investigationResultsMetadata)
}

func (invResService *InvestigationResultService) DeleteInvestigationResultsMetadataByIdsWithTx(ctx context.Context, tx *gorm.DB,
	investigationResultsMetadataIds []uint) *commonStructures.CommonError {

	if len(investigationResultsMetadataIds) == 0 {
		return nil
	}

	return invResService.InvResDao.DeleteInvestigationResultsMetadataByIdsWithTx(tx, investigationResultsMetadataIds)
}
