package consumerTasks

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"gorm.io/gorm"

	"github.com/Orange-Health/citadel/common/constants"
	"github.com/Orange-Health/citadel/common/structures"
	"github.com/Orange-Health/citadel/common/utils"
	"github.com/Orange-Health/citadel/models"
)

func (eventProcessor *EventProcessor) OmsManualReportUploadEventTask(ctx context.Context, eventPayload string) error {
	omsManualReportUploadEvent := structures.OmsManualReportUploadEvent{}
	err := json.Unmarshal(json.RawMessage(eventPayload), &omsManualReportUploadEvent)
	if err != nil {
		eventProcessor.Sentry.LogError(ctx, constants.ERROR_FAILED_TO_UNMARSHAL_JSON, err, nil)
		return err
	}

	redisKey := fmt.Sprintf(constants.OmsManualUploadKey, omsManualReportUploadEvent.AlnumOrderId,
		omsManualReportUploadEvent.CityCode)
	keyExists, err := eventProcessor.Cache.Exists(ctx, redisKey)
	if err != nil || keyExists {
		utils.AddLog(ctx, constants.DEBUG_LEVEL, utils.GetCurrentFunctionName(), nil, err)
		return errors.New(constants.ERROR_MANUAL_REPORT_UPLOAD_TASK_IN_PROGRESS)
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

	taskMetadata := models.TaskMetadata{}
	omsOrderId := omsManualReportUploadEvent.AlnumOrderId
	task, cErr := eventProcessor.TaskService.GetTaskByOmsOrderId(omsOrderId)
	if cErr != nil {
		return errors.New(constants.ERROR_TASK_NOT_FOUND)
	}

	if task.Id != 0 {
		taskMetadata, cErr = eventProcessor.TaskService.GetTaskMetadataByTaskId(task.Id)
		if cErr != nil {
			err := errors.New(cErr.Message)
			utils.AddLog(ctx, constants.ERROR_LEVEL, utils.GetCurrentFunctionName(), nil, err)
			return err
		}
	}

	testDetails, cErr := eventProcessor.TestDetailService.GetTestDetailsByOmsOrderId(omsOrderId)
	if cErr != nil {
		err := errors.New(cErr.Message)
		utils.AddLog(ctx, constants.ERROR_LEVEL, utils.GetCurrentFunctionName(), nil, err)
		return err
	}

	testDetailsToBeUpdated, testDetailsIdsToBeUpdated := []models.TestDetail{}, []uint{}
	testDetailsIds := []uint{}
	for index := range testDetails {
		testDetailsIds = append(testDetailsIds, testDetails[index].Id)
		if testDetails[index].Status != constants.TEST_STATUS_APPROVE &&
			utils.SliceContainsString(omsManualReportUploadEvent.AlnumTestIds, testDetails[index].CentralOmsTestId) {
			testDetails[index].Status = constants.TEST_STATUS_APPROVE
			testDetails[index].IsManualReportUpload = true
			testDetailsToBeUpdated = append(testDetailsToBeUpdated, testDetails[index])
			testDetailsIdsToBeUpdated = append(testDetailsIdsToBeUpdated, testDetails[index].Id)
		}
	}

	if len(testDetailsToBeUpdated) == 0 {
		err = eventProcessor.CommonTaskProcessor.ReleaseReportByOmsOrderIdPostManualUploadReportTask(ctx, omsOrderId,
			omsManualReportUploadEvent.AlnumTestIds)
		if err != nil {
			return err
		}
		return nil
	}

	testDetailsMetadata, cErr := eventProcessor.TestDetailService.GetTestDetailsMetadataByTestDetailIds(testDetailsIds)
	if cErr != nil {
		err := errors.New(cErr.Message)
		utils.AddLog(ctx, constants.ERROR_LEVEL, utils.GetCurrentFunctionName(), nil, err)
		return err
	}

	investigations, cErr := eventProcessor.InvestigationResultsService.
		GetInvestigationResultsByTestDetailsIds(testDetailsIdsToBeUpdated)
	if cErr != nil {
		err := errors.New(cErr.Message)
		utils.AddLog(ctx, constants.ERROR_LEVEL, utils.GetCurrentFunctionName(), nil, err)
		return err
	}

	approveId, _ := strconv.ParseUint(
		constants.AutoApprovalIdsMap[strings.ToLower(omsManualReportUploadEvent.CityCode)], 10, 64)
	for index := range investigations {
		investigations[index].InvestigationStatus = constants.INVESTIGATION_STATUS_APPROVE
		investigations[index].ApprovedAt = utils.GetCurrentTime()
		investigations[index].ApprovedBy = uint(approveId)
	}

	if task.Id != 0 {
		labIdLabMap := eventProcessor.CdsService.GetLabIdLabMap(ctx)
		task, taskMetadata = updateTaskBasedOnTestDetailsCompletion(task, taskMetadata, testDetails, testDetailsMetadata,
			labIdLabMap)
	}

	err = eventProcessor.updateDbPostManualReportUploadEvent(task, taskMetadata, testDetailsToBeUpdated, investigations)
	if err != nil {
		utils.AddLog(ctx, constants.ERROR_LEVEL, utils.GetCurrentFunctionName(), nil, err)
		return err
	}

	err = eventProcessor.CommonTaskProcessor.ReleaseReportByOmsOrderIdPostManualUploadReportTask(ctx, omsOrderId,
		omsManualReportUploadEvent.AlnumTestIds)
	if err != nil {
		return err
	}

	eventProcessor.EtsService.GetAndPublishEtsTestBasicEvent(ctx, omsManualReportUploadEvent.AlnumTestIds)

	return nil
}

func updateTaskBasedOnTestDetailsCompletion(task models.Task, taskMetadata models.TaskMetadata,
	allTestDetails []models.TestDetail, testDetailsMetadata []models.TestDetailsMetadata,
	labIdLabMap map[uint]structures.Lab) (models.Task, models.TaskMetadata) {

	if len(allTestDetails) == 0 {
		return task, taskMetadata
	}

	testCompleteStatuses := []string{
		constants.TEST_STATUS_APPROVE,
		constants.TEST_STATUS_COMPLETED_NOT_SENT,
		constants.TEST_STATUS_COMPLETED_SENT,
		constants.TEST_STATUS_SAMPLE_NOT_RECEIVED,
	}

	testDetailsMap := make(map[uint]models.TestDetail)
	for _, testDetail := range allTestDetails {
		testDetailsMap[testDetail.Id] = testDetail
	}

	allInhouseTestsCompleted := true
	for _, testDetail := range allTestDetails {
		if !utils.SliceContainsString(testCompleteStatuses, testDetail.Status) &&
			utils.IsTestInhouse(testDetail.ProcessingLabId, task.LabId, labIdLabMap) {
			allInhouseTestsCompleted = false
			break
		}
	}

	// Set Task Status as completed if all tests are approved
	if allInhouseTestsCompleted {
		task.PreviousStatus = task.Status
		task.Status = constants.TASK_STATUS_COMPLETED
		task.CompletedAt = utils.GetCurrentTime()
	}

	taskMetadata.IsCritical = getTaskCriticality(testDetailsMetadata, testDetailsMap)

	doctorTat := getDoctorTatForTask(allTestDetails)
	if !doctorTat.IsZero() {
		task.DoctorTat = doctorTat
	}

	if task.Status == constants.TASK_STATUS_IN_PROGRESS {
		task.Status = task.PreviousStatus
	}

	return task, taskMetadata
}

func (eventProcessor *EventProcessor) updateDbPostManualReportUploadEvent(task models.Task,
	taskMetadata models.TaskMetadata, testDetailsToBeUpdated []models.TestDetail,
	investigations []models.InvestigationResult) error {
	err := eventProcessor.Db.Transaction(func(tx *gorm.DB) error {
		if task.Id > 0 {
			cErr := eventProcessor.TaskPathMappingService.MarkAllTaskPathMapInactiveByTaskIdWithTx(tx, task.Id)
			if cErr != nil {
				return errors.New(cErr.Message)
			}
		}

		if len(investigations) > 0 {
			_, cErr := eventProcessor.InvestigationResultsService.UpdateInvestigationResultsWithTx(tx, investigations)
			if cErr != nil {
				return errors.New(cErr.Message)
			}
		}

		if len(testDetailsToBeUpdated) > 0 {
			_, cErr := eventProcessor.TestDetailService.UpdateTestDetailsWithTx(tx, testDetailsToBeUpdated)
			if cErr != nil {
				return errors.New(cErr.Message)
			}
		}

		if task.Id > 0 {
			_, cErr := eventProcessor.TaskService.UpdateTaskWithTx(tx, task)
			if cErr != nil {
				return errors.New(cErr.Message)
			}
		}

		if task.Id > 0 {
			_, cErr := eventProcessor.TaskService.UpdateTaskMetadataWithTx(tx, taskMetadata)
			if cErr != nil {
				return errors.New(cErr.Message)
			}
		}

		return nil
	})

	return err
}
