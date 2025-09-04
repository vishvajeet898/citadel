package service

import (
	"context"
	"strconv"
	"time"

	commonConstants "github.com/Orange-Health/citadel/common/constants"
	commonStructures "github.com/Orange-Health/citadel/common/structures"
	commonUtils "github.com/Orange-Health/citadel/common/utils"
	commonModels "github.com/Orange-Health/citadel/models"
)

func (taskService *TaskService) GetVisitIdToAttuneResponseMap(
	ctx context.Context, omsTestIds []string,
) (map[string]commonStructures.AttuneOrderResponse, *commonStructures.CommonError) {
	visitIdToAttuneResponseMap := map[string]commonStructures.AttuneOrderResponse{}
	visitIdLabMap, cErr := taskService.SampleService.GetVisitLabMapByOmsTestIds(omsTestIds)
	if cErr != nil {
		return visitIdToAttuneResponseMap, cErr
	}

	for visitId, labId := range visitIdLabMap {
		if visitId == "" || labId == 0 {
			continue
		}
		attuneResponse, cErr := taskService.AttuneClient.GetPatientVisitDetailsbyVisitNo(ctx, visitId,
			commonConstants.AttuneReportWithStationery, labId)
		if cErr != nil {
			return visitIdToAttuneResponseMap, cErr
		}
		visitIdToAttuneResponseMap[visitId] = attuneResponse
	}

	return visitIdToAttuneResponseMap, nil
}

func GetUpdatedAttuneOrderDetailsForRerunTests(
	response commonStructures.AttuneOrderResponse,
	lisCodeForTestsToBeRerun []string,
	lisCodeToValueMap map[string]string,
	lisCodeToRerunInvestigationMap map[string]bool,
	lisCodeToRerunDetailsMap map[string]commonModels.RerunInvestigationResult,
	user commonModels.User,
) commonStructures.AttuneOrderResponse {
	testStatus := commonConstants.ATTUNE_TEST_STATUS_RECHECK
	currentTime := time.Now()
	rerunTime := currentTime.Format(commonConstants.DateTimeUTCWithFractionSecWithoutZOffset)
	orderInfoStruct, newOrderInfoStruct := response.OrderInfo, []commonStructures.AttuneOrderInfo{}

	attuneUserId, _ := strconv.ParseUint(user.AttuneUserId, 10, 64)

	for _, orderInfo := range orderInfoStruct {
		if commonUtils.SliceContainsString(lisCodeForTestsToBeRerun, orderInfo.TestCode) {
			orderInfo.TestStatus = testStatus
			orderInfo.UserID = int(attuneUserId)
			orderInfo.RerunTime = rerunTime
			orderInfo.RerunReason = commonConstants.DEFAULT_RERUN_REASON
			orderInfo.RerunRemarks = commonConstants.DEFAULT_RERUN_REMARK

			if orderInfo.TestType == commonConstants.InvestigationShortHand {
				if value, ok := lisCodeToValueMap[orderInfo.TestCode]; ok {
					orderInfo.TestValue = value
				}

				if rerunDetails, ok := lisCodeToRerunDetailsMap[orderInfo.TestCode]; ok {
					orderInfo.RerunReason = rerunDetails.RerunReason
					orderInfo.RerunRemarks = rerunDetails.RerunRemarks
				}
			}

			for j := range orderInfo.OrderContentListInfo {
				if value, ok := lisCodeToValueMap[orderInfo.OrderContentListInfo[j].TestCode]; ok {
					orderInfo.OrderContentListInfo[j].TestValue = value
				}
				orderInfo.OrderContentListInfo[j].UserID = int(attuneUserId)
				if rerunDetails, ok := lisCodeToRerunDetailsMap[orderInfo.OrderContentListInfo[j].TestCode]; ok {
					orderInfo.OrderContentListInfo[j].TestStatus = testStatus
					orderInfo.OrderContentListInfo[j].RerunReason = rerunDetails.RerunReason
					orderInfo.OrderContentListInfo[j].RerunRemarks = rerunDetails.RerunRemarks
					orderInfo.OrderContentListInfo[j].RerunTime = rerunTime
				}

				params := orderInfo.OrderContentListInfo[j].ParameterListInfo
				for k := range params {
					params[k].UserID = int(attuneUserId)
					if rerunDetails, ok := lisCodeToRerunDetailsMap[params[k].TestCode]; ok {
						params[k].TestStatus = testStatus
						params[k].RerunReason = rerunDetails.RerunReason
						params[k].RerunRemarks = rerunDetails.RerunRemarks
						params[k].RerunTime = rerunTime
					}
					if value, ok := lisCodeToValueMap[params[k].TestCode]; ok {
						params[k].TestValue = value
					}
				}
				orderInfo.OrderContentListInfo[j].ParameterListInfo = params
			}

			newOrderInfoStruct = append(newOrderInfoStruct, orderInfo)
		}
	}

	return commonStructures.AttuneOrderResponse{
		OrderId:   response.OrderId,
		OrgCode:   response.OrgCode,
		OrderInfo: newOrderInfoStruct,
	}
}

func (taskService *TaskService) FetchAttuneDataForReruningTests(ctx context.Context, cityCode string,
	testDetailsIdsToBeRerun []uint,
	testDetailIdToTestDetailMap map[uint]commonModels.TestDetail,
	investigations []commonModels.InvestigationResult,
	rerunDetails []commonModels.RerunInvestigationResult,
	user commonModels.User) (
	map[string]commonStructures.AttuneOrderResponse, *commonStructures.CommonError) {

	updatedVisitIdToAttuneResponseMap := map[string]commonStructures.AttuneOrderResponse{}
	lisCodeToValueMap := map[string]string{}
	lisCodeToRerunDetailsMap := map[string]commonModels.RerunInvestigationResult{}
	lisCodeToRerunInvestigationMap := map[string]bool{}
	lisCodesForTestToBeRerun := []string{}
	omsTestIds := []string{}

	for _, testDetailId := range testDetailsIdsToBeRerun {
		omsTestIds = append(omsTestIds, testDetailIdToTestDetailMap[testDetailId].CentralOmsTestId)
		lisCodesForTestToBeRerun = append(lisCodesForTestToBeRerun, testDetailIdToTestDetailMap[testDetailId].LisCode)
	}

	if len(omsTestIds) == 0 {
		return updatedVisitIdToAttuneResponseMap, nil
	}

	visitIdToAttuneResponseMap, cErr := taskService.GetVisitIdToAttuneResponseMap(ctx, omsTestIds)
	if cErr != nil {
		return updatedVisitIdToAttuneResponseMap, cErr
	}

	for _, investigation := range investigations {
		lisCodeToValueMap[investigation.LisCode] = investigation.InvestigationValue
		if commonUtils.SliceContainsString(commonConstants.INVESTIGATION_STATUSES_RERUN, investigation.InvestigationStatus) {
			lisCodeToRerunInvestigationMap[investigation.LisCode] = true
		}
	}

	for _, rerunDetail := range rerunDetails {
		lisCodeToRerunDetailsMap[rerunDetail.LisCode] = rerunDetail
	}

	for visitId, attuneResponse := range visitIdToAttuneResponseMap {
		updatedAttuneResponse := GetUpdatedAttuneOrderDetailsForRerunTests(attuneResponse,
			lisCodesForTestToBeRerun, lisCodeToValueMap, lisCodeToRerunInvestigationMap, lisCodeToRerunDetailsMap, user)
		updatedVisitIdToAttuneResponseMap[visitId] = updatedAttuneResponse
	}

	commonUtils.AddLog(ctx, commonConstants.DEBUG_LEVEL, commonUtils.GetCurrentFunctionName(), map[string]interface{}{
		"updatedVisitIdToAttuneResponseMap": updatedVisitIdToAttuneResponseMap,
	}, nil)

	return updatedVisitIdToAttuneResponseMap, nil
}

func GetUpdatedAttuneOrderDetailsForApprovedTests(
	response commonStructures.AttuneOrderResponse,
	lisCodeForTestsToBeApproved []string,
	lisCodeToValueMap map[string]string,
	lisCodeToMedicalRemarkMap map[string]string,
	testLisCodeToApprovedByMap map[string]uint,
	approvedByIdToAttuneIdMap map[uint]uint,
) commonStructures.AttuneOrderResponse {
	testStatus := commonConstants.ATTUNE_TEST_STATUS_APPROVED
	currentTime := time.Now()
	approvedAt := currentTime.Format(commonConstants.DateTimeUTCWithFractionSecWithoutZOffset)
	orderInfoStruct, newOrderInfoStruct := response.OrderInfo, []commonStructures.AttuneOrderInfo{}

	for _, orderInfo := range orderInfoStruct {
		if commonUtils.SliceContainsString(lisCodeForTestsToBeApproved, orderInfo.TestCode) {
			userId := int(approvedByIdToAttuneIdMap[testLisCodeToApprovedByMap[orderInfo.TestCode]])
			orderInfo.TestStatus = testStatus
			orderInfo.UserID = userId
			orderInfo.ResultApprovedAt = &approvedAt
			orderInfo.RerunReason = ""
			orderInfo.RerunRemarks = ""

			if orderInfo.TestType == commonConstants.InvestigationShortHand {
				if value, ok := lisCodeToValueMap[orderInfo.TestCode]; ok {
					orderInfo.TestValue = value
				}

				if medicalRemark, ok := lisCodeToMedicalRemarkMap[orderInfo.TestCode]; ok {
					orderInfo.MedicalRemarks = medicalRemark
				} else {
					orderInfo.MedicalRemarks = ""
				}
			}

			for j := range orderInfo.OrderContentListInfo {
				orderInfo.OrderContentListInfo[j].TestStatus = testStatus
				if value, ok := lisCodeToValueMap[orderInfo.OrderContentListInfo[j].TestCode]; ok {
					orderInfo.OrderContentListInfo[j].TestValue = value
				}

				if medicalRemark, ok := lisCodeToMedicalRemarkMap[orderInfo.OrderContentListInfo[j].TestCode]; ok {
					orderInfo.OrderContentListInfo[j].MedicalRemarks = medicalRemark
				} else {
					orderInfo.OrderContentListInfo[j].MedicalRemarks = ""
				}

				orderInfo.OrderContentListInfo[j].UserID = userId
				orderInfo.OrderContentListInfo[j].ResultApprovedAt = &approvedAt
				orderInfo.OrderContentListInfo[j].RerunReason = ""
				orderInfo.OrderContentListInfo[j].RerunRemarks = ""

				params := orderInfo.OrderContentListInfo[j].ParameterListInfo
				for k := range params {
					params[k].TestStatus = testStatus
					params[k].UserID = userId
					params[k].ResultApprovedAt = &approvedAt
					params[k].RerunReason = ""
					params[k].RerunRemarks = ""
					if value, ok := lisCodeToValueMap[params[k].TestCode]; ok {
						params[k].TestValue = value
					}

					if medicalRemark, ok := lisCodeToMedicalRemarkMap[params[k].TestCode]; ok {
						params[k].MedicalRemarks = medicalRemark
					} else {
						params[k].MedicalRemarks = ""
					}
				}
				orderInfo.OrderContentListInfo[j].ParameterListInfo = params
			}

			newOrderInfoStruct = append(newOrderInfoStruct, orderInfo)
		}
	}

	return commonStructures.AttuneOrderResponse{
		OrderId:   response.OrderId,
		OrgCode:   response.OrgCode,
		OrderInfo: newOrderInfoStruct,
	}
}

func (taskService *TaskService) FetchAttuneDataForApprovingTests(ctx context.Context, task commonModels.Task,
	testDetailsToBeAproved []commonModels.TestDetail,
	investigations []commonModels.InvestigationResult,
	medicalRemarks []commonModels.Remark,
	approvedByUsers []commonModels.User,
) (map[string]commonStructures.AttuneOrderResponse, *commonStructures.CommonError) {

	updatedVisitIdToAttuneResponseMap := map[string]commonStructures.AttuneOrderResponse{}
	lisCodeToValueMap, investigationIdToLisCodeMap := map[string]string{}, map[uint]string{}
	lisCodeToMedicalRemarkMap, approvedByIdToAttuneIdMap := map[string]string{}, map[uint]uint{}
	testDetailIdToTestLisCodeMap, testLisCodeToApprovedByMap := map[uint]string{}, map[string]uint{}
	omsTestIds, lisCodesForTestToBeApproved := []string{}, []string{}

	if len(testDetailsToBeAproved) == 0 {
		return updatedVisitIdToAttuneResponseMap, nil
	}

	for _, testDetail := range testDetailsToBeAproved {
		omsTestIds = append(omsTestIds, testDetail.CentralOmsTestId)
		lisCodesForTestToBeApproved = append(lisCodesForTestToBeApproved, testDetail.LisCode)
		testDetailIdToTestLisCodeMap[testDetail.Id] = testDetail.LisCode
	}

	visitIdToAttuneResponseMap, cErr := taskService.GetVisitIdToAttuneResponseMap(ctx, omsTestIds)
	if cErr != nil {
		return updatedVisitIdToAttuneResponseMap, cErr
	}

	for _, investigation := range investigations {
		lisCodeToValueMap[investigation.LisCode] = investigation.InvestigationValue
		investigationIdToLisCodeMap[investigation.Id] = investigation.LisCode
		testLisCodeToApprovedByMap[testDetailIdToTestLisCodeMap[investigation.TestDetailsId]] = investigation.ApprovedBy
	}

	for _, medicalRemark := range medicalRemarks {
		lisCodeToMedicalRemarkMap[investigationIdToLisCodeMap[medicalRemark.InvestigationResultId]] = medicalRemark.Description
	}

	for _, user := range approvedByUsers {
		attuneUserId, _ := strconv.ParseUint(user.AttuneUserId, 10, 64)
		approvedByIdToAttuneIdMap[user.Id] = uint(attuneUserId)
	}

	for visitId, attuneResponse := range visitIdToAttuneResponseMap {
		updatedOrderDetails := GetUpdatedAttuneOrderDetailsForApprovedTests(attuneResponse,
			lisCodesForTestToBeApproved, lisCodeToValueMap, lisCodeToMedicalRemarkMap, testLisCodeToApprovedByMap, approvedByIdToAttuneIdMap)
		if len(updatedOrderDetails.OrderInfo) > 0 {
			updatedVisitIdToAttuneResponseMap[visitId] = updatedOrderDetails
		}
	}

	commonUtils.AddLog(ctx, commonConstants.DEBUG_LEVEL, commonUtils.GetCurrentFunctionName(), map[string]interface{}{
		"updatedVisitIdToAttuneResponseMap": updatedVisitIdToAttuneResponseMap,
	}, nil)

	return updatedVisitIdToAttuneResponseMap, nil
}
