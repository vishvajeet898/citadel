package mapper

import (
	"sort"
	"strings"

	"github.com/Orange-Health/citadel/apps/investigation_results/structures"
	commonConstants "github.com/Orange-Health/citadel/common/constants"
	commonStructures "github.com/Orange-Health/citadel/common/structures"
	commonUtils "github.com/Orange-Health/citadel/common/utils"
	commonModels "github.com/Orange-Health/citadel/models"
)

func MapInvestigationData(invResultData []commonModels.InvestigationData) map[uint]commonModels.InvestigationData {
	invResultDataMap := make(map[uint]commonModels.InvestigationData)
	for _, invData := range invResultData {
		invResultDataMap[invData.InvestigationResultId] = invData
	}
	return invResultDataMap
}

func MapInvResultFromDbStruct(invResult structures.InvestigationResultDbResponse,
	invResultData commonModels.InvestigationData) structures.InvestigationResult {
	invRes := structures.InvestigationResult{
		TestDetailId:              invResult.TestDetailsId,
		TestDetailsStatus:         invResult.TestDetailsStatus,
		MasterInvestigationId:     invResult.MasterInvestigationId,
		InvestigationName:         invResult.InvestigationName,
		InvestigationValue:        strings.TrimSpace(invResult.InvestigationValue),
		DeviceValue:               invResult.DeviceValue,
		ResultRepresentationType:  invResult.ResultRepresentationType,
		Department:                invResult.Department,
		Uom:                       invResult.Uom,
		ReferenceRangeText:        invResult.ReferenceRangeText,
		Method:                    invResult.Method,
		MethodType:                invResult.MethodType,
		LisCode:                   invResult.LisCode,
		Abnormality:               invResult.Abnormality,
		ApprovedBy:                invResult.ApprovedBy,
		ApprovedAt:                commonUtils.GetTimeInString(invResult.ApprovedAt),
		EnteredBy:                 invResult.EnteredBy,
		EnteredAt:                 commonUtils.GetTimeInString(invResult.EnteredAt),
		Status:                    invResult.InvestigationStatus,
		IsCritical:                invResult.IsCritical,
		MasterTestId:              invResult.MasterTestId,
		Barcodes:                  invResult.Barcodes,
		AutoApprovalFailureReason: invResult.AutoApprovalFailureReason,
		ProcessingLabId:           invResult.ProcessingLabId,
		QcFlag:                    invResult.QcFlag,
		QcLotNumber:               invResult.QcLotNumber,
		QcValue:                   invResult.QcValue,
		QcWestGardWarning:         invResult.QcWestGardWarning,
		QcStatus:                  invResult.QcStatus,
	}

	invRes.Id = invResult.Id

	invDataStruct := structures.InvestigationData{}
	invDataStruct.Id = invResultData.Id
	invDataStruct.Data = invResultData.Data
	invDataStruct.DataType = invResultData.DataType

	invRes.InvestigationData = invDataStruct

	return invRes
}

func MapInvestigationResultsFromDbStruct(investigationResults []structures.InvestigationResultDbResponse,
	invResultDataMap map[uint]commonModels.InvestigationData) []structures.InvestigationResult {
	var invRes []structures.InvestigationResult
	for _, invResult := range investigationResults {
		invRes = append(invRes, MapInvResultFromDbStruct(invResult, invResultDataMap[invResult.Id]))
	}
	return invRes
}

func MapInvestigationIdsFromInvestigationResultsDbStructs(
	investigationResults []structures.InvestigationResultDbResponse) []uint {
	var investigationResultIds []uint
	for _, investigation := range investigationResults {
		investigationResultIds = append(investigationResultIds, investigation.Id)
	}
	return investigationResultIds
}

func MapRemarks(remarks []commonModels.Remark) map[uint][]structures.Remark {
	remarksMap := make(map[uint][]structures.Remark)
	for _, r := range remarks {
		remarksMap[r.InvestigationResultId] = append(remarksMap[r.InvestigationResultId], structures.Remark{
			Id:                    r.Id,
			InvestigationResultId: r.InvestigationResultId,
			Description:           r.Description,
			RemarkType:            r.RemarkType,
			RemarkBy:              r.RemarkBy,
		})
	}
	return remarksMap
}

func MapRerunDetails(rrds []commonModels.RerunInvestigationResult) map[uint][]structures.RerunDetails {
	rerunMap := make(map[uint][]structures.RerunDetails)
	for _, r := range rrds {
		rerunMap[r.MasterInvestigationId] = append(rerunMap[r.MasterInvestigationId], structures.RerunDetails{
			Id:                    r.Id,
			TestDetailId:          r.TestDetailsId,
			MasterInvestigationId: r.MasterInvestigationId,
			InvestigationValue:    r.InvestigationValue,
			DeviceValue:           r.DeviceValue,
			RerunTriggeredBy:      r.RerunTriggeredBy,
			RerunTriggeredAt:      r.RerunTriggeredAt,
			RerunReason:           r.RerunReason,
			RerunRemarks:          r.RerunRemarks,
		})
	}
	return rerunMap
}

func CreateInvestigationResultsResponse(invResults []structures.InvestigationResult, remarksMap map[uint][]structures.Remark,
	rerunMap map[uint][]structures.RerunDetails) []structures.InvestigationResultDetails {
	invResDetails := []structures.InvestigationResultDetails{}
	for _, invRes := range invResults {
		invResDetails = append(invResDetails, structures.InvestigationResultDetails{
			InvestigationResult: invRes,
			Remarks:             remarksMap[invRes.Id],
			RerunDetails:        rerunMap[invRes.MasterInvestigationId],
		})
	}

	return invResDetails
}

func CreateResponseForPatientPastRecords(
	patientDetails commonStructures.PatientDetailsResponse,
	patientPastRecords []commonStructures.PatientPastRecords,
) commonStructures.PatientPastRecordsApiResponse {

	datesSlice := []string{}
	investigationIdToPastRecordsDateMap := make(map[uint]map[string][]commonStructures.PatientPastRecordsValue)
	investigationIdToNameMap := make(map[uint]string)
	departmentToInvestigationIdMap := make(map[string][]uint)

	for _, record := range patientPastRecords {
		approvedAt := record.ApprovedAt
		if approvedAt == nil {
			continue
		}
		approvedAtDate := approvedAt.Format(commonConstants.DateLayout)
		datesSlice = append(datesSlice, approvedAtDate)

		if record.MasterInvestigationId != 0 {
			if _, ok := investigationIdToPastRecordsDateMap[record.MasterInvestigationId]; !ok {
				investigationIdToPastRecordsDateMap[record.MasterInvestigationId] = make(map[string][]commonStructures.PatientPastRecordsValue)
			}
			investigationIdToPastRecordsDateMap[record.MasterInvestigationId][approvedAtDate] = append(
				investigationIdToPastRecordsDateMap[record.MasterInvestigationId][approvedAtDate],
				commonStructures.PatientPastRecordsValue{
					InvestigationValue: record.InvestigationValue,
					Uom:                record.Uom,
				},
			)
			investigationIdToNameMap[record.MasterInvestigationId] = record.InvestigationName
			departmentName := commonConstants.DEPARTMENT_NAME_MAP[strings.ToLower(record.Department)]
			departmentName = strings.TrimSpace(departmentName)
			if departmentName == "" {
				departmentName = record.Department
			}
			departmentToInvestigationIdMap[departmentName] = append(
				departmentToInvestigationIdMap[departmentName], record.MasterInvestigationId)
		}
	}

	datesSlice = commonUtils.CreateUniqueSliceString(datesSlice)
	sort.Slice(datesSlice, func(i, j int) bool {
		return datesSlice[i] > datesSlice[j]
	})

	investigationIdToPastRecordsMap := make(map[uint]commonStructures.PatientPastRecordsInvestigations)
	for investigationId, pastRecordsDateMap := range investigationIdToPastRecordsDateMap {
		investigationIdToPastRecordsMap[investigationId] = commonStructures.PatientPastRecordsInvestigations{
			InvestigationId:   investigationId,
			InvestigationName: investigationIdToNameMap[investigationId],
			PastRecords:       pastRecordsDateMap,
		}
	}

	for department, investigationIds := range departmentToInvestigationIdMap {
		investigationIds = commonUtils.CreateUniqueSliceUint(investigationIds)
		departmentToInvestigationIdMap[department] = investigationIds
	}

	departmentToPastRecordsInvestigationsMap := make(map[string][]commonStructures.PatientPastRecordsInvestigations)
	for department, investigationIds := range departmentToInvestigationIdMap {
		for _, investigationId := range investigationIds {
			departmentToPastRecordsInvestigationsMap[department] = append(
				departmentToPastRecordsInvestigationsMap[department], investigationIdToPastRecordsMap[investigationId])
		}
	}

	patientPastRecordsData := []commonStructures.PatientPastRecordsData{}
	for department, pastRecordsInvestigations := range departmentToPastRecordsInvestigationsMap {
		patientPastRecordsData = append(patientPastRecordsData, commonStructures.PatientPastRecordsData{
			Department:     department,
			Investigations: pastRecordsInvestigations,
		})
	}

	return commonStructures.PatientPastRecordsApiResponse{
		PatientDetails: commonStructures.PatientPastRecordsPatientDetails{
			Name:      patientDetails.FullName,
			AgeYears:  patientDetails.AgeYears,
			AgeMonths: patientDetails.AgeMonths,
			AgeDays:   patientDetails.AgeDays,
			Gender:    commonUtils.GetGenderConstantForPatient(patientDetails.Gender),
		},
		Dates: datesSlice,
		Data:  patientPastRecordsData,
	}
}

func MapInvestigationResultsByMasterInvestigationId(investigationResults []commonModels.InvestigationResult) map[uint]commonModels.InvestigationResult {
	investigationResultsByMasterInvestigationId := make(map[uint]commonModels.InvestigationResult)
	for _, invRes := range investigationResults {
		investigationResultsByMasterInvestigationId[invRes.MasterInvestigationId] = invRes
	}
	return investigationResultsByMasterInvestigationId
}

func MapTestDetailsIdsFromTestDetailsStructs(testDetails []structures.TestDetailsDbResponse) []uint {
	testDetailsIds := []uint{}
	for _, testDetail := range testDetails {
		testDetailsIds = append(testDetailsIds, testDetail.Id)
	}
	return commonUtils.CreateUniqueSliceUint(testDetailsIds)
}

func MapInvestigationResultsDbStructsWithMasterTests(invResults []structures.InvestigationResultDbResponse,
	testDetails []structures.TestDetailsDbResponse) []structures.InvestigationResultDbResponse {

	testDetailsMap := map[uint]structures.TestDetailsDbResponse{}

	finalInvestigationResults := []structures.InvestigationResultDbResponse{}

	for _, testDetail := range testDetails {
		testDetailsMap[testDetail.Id] = testDetail
	}

	for _, invResult := range invResults {
		testDetail := testDetailsMap[invResult.TestDetailsId]
		if !testDetail.CpEnabled {
			continue
		}

		if invResult.TestDetailsId != 0 {
			invResult.TestDetailsStatus = testDetailsMap[invResult.TestDetailsId].Status
			invResult.MasterTestId = testDetailsMap[invResult.TestDetailsId].MasterTestId
			invResult.Barcodes = testDetailsMap[invResult.TestDetailsId].Barcodes

			finalInvestigationResults = append(finalInvestigationResults, invResult)
		}
	}

	return finalInvestigationResults
}
