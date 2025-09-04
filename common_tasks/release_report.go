package commonTasks

import (
	"context"
	"errors"

	"github.com/Orange-Health/citadel/common/constants"
	"github.com/Orange-Health/citadel/common/structures"
	"github.com/Orange-Health/citadel/common/utils"
	"github.com/Orange-Health/citadel/models"
)

func (ctp *CommonTaskProcessor) ReleaseReportTask(ctx context.Context, taskId uint, forcefulUpdate bool) error {
	reportReleaseTestIds, attuneTestIds := []string{}, []uint{}
	attuneTestDetails := []models.TestDetail{}

	task, cErr := ctp.TaskService.GetTaskModelById(taskId)
	if cErr != nil {
		return errors.New(cErr.Message)
	}

	testDetails, cErr := ctp.TestDetailService.GetTestDetailsByOmsOrderId(task.OmsOrderId)
	if cErr != nil {
		return errors.New(cErr.Message)
	}

	for _, testDetail := range testDetails {
		if testDetail.Status == constants.TEST_STATUS_APPROVE {
			reportReleaseTestIds = append(reportReleaseTestIds, testDetail.CentralOmsTestId)
			if !testDetail.CpEnabled {
				continue
			}
			if forcefulUpdate {
				attuneTestIds = append(attuneTestIds, testDetail.Id)
				attuneTestDetails = append(attuneTestDetails, testDetail)
			} else if testDetail.ReportSentAt == nil {
				attuneTestIds = append(attuneTestIds, testDetail.Id)
				attuneTestDetails = append(attuneTestDetails, testDetail)
			}
		}
	}

	if len(reportReleaseTestIds) == 0 {
		return nil
	}

	investigationResults, cErr := ctp.InvestigationResultsService.GetInvestigationResultsByTestDetailsIds(attuneTestIds)
	if cErr != nil {
		return errors.New(cErr.Message)
	}

	investigationResultIds := []uint{}
	for _, investigationResult := range investigationResults {
		investigationResultIds = append(investigationResultIds, investigationResult.Id)
	}

	investigationsData, cErr := ctp.InvestigationResultsService.GetInvestigationDataByInvestigationResultsIds(
		investigationResultIds)
	if cErr != nil {
		return errors.New(cErr.Message)
	}

	investigationIdToInvestigationDataMap := map[uint]models.InvestigationData{}
	for _, investigationData := range investigationsData {
		investigationIdToInvestigationDataMap[investigationData.InvestigationResultId] = investigationData
	}

	for index := range investigationResults {
		if investigationData, ok := investigationIdToInvestigationDataMap[investigationResults[index].Id]; ok {
			investigationResults[index].InvestigationValue = investigationData.Data
		}
	}

	medicalRemarks, cErr := ctp.RemarkService.GetRemarksByInvestigationResultIds(
		[]string{constants.REMARK_TYPE_MEDICAL_REMARK}, investigationResultIds)
	if cErr != nil {
		return errors.New(cErr.Message)
	}

	approvedByIds := []uint{}
	for _, investigationResult := range investigationResults {
		approvedByIds = append(approvedByIds, investigationResult.ApprovedBy)
	}

	approvedByIds = utils.CreateUniqueSliceUint(approvedByIds)

	approvedByUsers, cErr := ctp.UserService.GetUsersByIds(approvedByIds)
	if cErr != nil {
		return errors.New(cErr.Message)
	}

	visitIdToAttuneResponseMap, cErr := ctp.TaskService.FetchAttuneDataForApprovingTests(ctx,
		task, attuneTestDetails, investigationResults, medicalRemarks, approvedByUsers)
	if cErr != nil {
		return errors.New(cErr.Message)
	}

	for _, attuneResponse := range visitIdToAttuneResponseMap {
		for index := range attuneResponse.OrderInfo {
			newAttuneResponse := structures.AttuneOrderResponse{
				OrderId: attuneResponse.OrderId,
				OrgCode: attuneResponse.OrgCode,
				OrderInfo: []structures.AttuneOrderInfo{
					attuneResponse.OrderInfo[index],
				},
			}
			cErr = ctp.AttuneClient.InsertTestDataToAttune(ctx, newAttuneResponse)
			if cErr != nil {
				utils.AddLog(ctx, constants.ERROR_LEVEL, utils.GetCurrentFunctionName(), nil, errors.New(cErr.Message))
			}
		}
	}

	cErr = ctp.ReportGenerationService.TriggerReportGenerationEvent(ctx, task.OmsOrderId, reportReleaseTestIds)
	if cErr != nil {
		return errors.New(cErr.Message)
	}

	currentTime := utils.GetCurrentTime()
	for index := range attuneTestDetails {
		attuneTestDetails[index].ReportSentAt = currentTime
		attuneTestDetails[index].UpdatedBy = constants.CitadelSystemId
		attuneTestDetails[index].ReportStatus = constants.TEST_REPORT_STATUS_QUEUED
	}
	_, cErr = ctp.TestDetailService.UpdateTestDetails(attuneTestDetails)
	if cErr != nil {
		return errors.New(cErr.Message)
	}

	return nil
}

func (ctp *CommonTaskProcessor) ReleaseReportByOmsOrderIdPostManualUploadReportTask(ctx context.Context, omsOrderId string,
	omsTestIds []string) error {

	cErr := ctp.ReportGenerationService.TriggerReportGenerationEvent(ctx, omsOrderId, []string{})
	if cErr != nil {
		return errors.New(cErr.Message)
	}

	cErr = ctp.TestDetailService.UpdateReportStatusByOmsTestIds(omsTestIds, constants.TEST_REPORT_STATUS_NOT_READY,
		constants.TEST_REPORT_STATUS_QUEUED)
	if cErr != nil {
		utils.AddLog(ctx, constants.ERROR_LEVEL, utils.GetCurrentFunctionName(), map[string]interface{}{
			"omsTestIds": omsTestIds,
			"omsOrderId": omsOrderId,
		}, errors.New(cErr.Message))
		return errors.New(cErr.Message)
	}

	return nil
}
