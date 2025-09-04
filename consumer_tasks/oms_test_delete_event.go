package consumerTasks

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"gorm.io/gorm"

	"github.com/Orange-Health/citadel/common/constants"
	"github.com/Orange-Health/citadel/common/structures"
	"github.com/Orange-Health/citadel/common/utils"
	"github.com/Orange-Health/citadel/models"
)

func (eventProcessor *EventProcessor) OmsTestDeleteEventTask(ctx context.Context, eventPayload string) error {
	omsTestDeleteEvent := structures.OmsTestDeleteEvent{}
	err := json.Unmarshal(json.RawMessage(eventPayload), &omsTestDeleteEvent)
	if err != nil {
		eventProcessor.Sentry.LogError(ctx, constants.ERROR_FAILED_TO_UNMARSHAL_JSON, err, nil)
		return err
	}

	redisKey := fmt.Sprintf(constants.OmsCreateUpdateOrderEventKey, omsTestDeleteEvent.AlnumOrderId)
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

	toDeleteOmsTestId, toDeleteTestDetailsId := omsTestDeleteEvent.AlnumTestId, uint(0)
	toDeleteTestDetail := models.TestDetail{}
	investigations := []models.InvestigationResult{}

	omsOrderId := omsTestDeleteEvent.AlnumOrderId

	task, cErr := eventProcessor.TaskService.GetTaskByOmsOrderId(omsOrderId)
	if cErr != nil {
		return nil
	}

	testDetails, cErr := eventProcessor.TestDetailService.GetTestDetailsByOmsOrderId(omsOrderId)
	if cErr != nil {
		err := errors.New(cErr.Message)
		utils.AddLog(ctx, constants.ERROR_LEVEL, utils.GetCurrentFunctionName(), nil, err)
		return err
	}

	if task.Id != 0 {
		investigations, cErr = eventProcessor.InvestigationResultsService.
			GetInvestigationResultsByTaskIdAndOmsTestId(task.Id, omsTestDeleteEvent.AlnumTestId)
		if cErr != nil {
			err := errors.New(cErr.Message)
			utils.AddLog(ctx, constants.ERROR_LEVEL, utils.GetCurrentFunctionName(), nil, err)
			return err
		}
	}
	investigationIds := []uint{}
	for _, investigation := range investigations {
		investigationIds = append(investigationIds, investigation.Id)
	}

	updatedTestDetails := []models.TestDetail{}
	for _, testDetail := range testDetails {
		if testDetail.CentralOmsTestId == omsTestDeleteEvent.AlnumTestId {
			toDeleteTestDetail = testDetail
			toDeleteTestDetailsId = testDetail.Id
			continue
		}
		updatedTestDetails = append(updatedTestDetails, testDetail)
	}

	err = eventProcessor.Db.Transaction(func(tx *gorm.DB) error {
		if task.Id != 0 {
			cErr = eventProcessor.TaskPathMappingService.MarkAllTaskPathMapInactiveByTaskIdWithTx(tx, task.Id)
			if cErr != nil {
				return errors.New(cErr.Message)
			}
		}

		if len(investigationIds) > 0 {
			cErr = eventProcessor.InvestigationResultsService.DeleteInvestigationResultsByIdsWithTx(
				tx, investigationIds)
			if cErr != nil {
				return errors.New(cErr.Message)
			}
		}

		cErr = eventProcessor.TestDetailService.DeleteTestDetailsByOmsTestIdsWithTx(tx, []string{toDeleteOmsTestId})
		if cErr != nil {
			return errors.New(cErr.Message)
		}

		if toDeleteTestDetail.Id != 0 {
			cErr = eventProcessor.TestDetailService.DeleteTestDetailsMetadataByTestDetailsIdsWithTx(tx,
				[]uint{toDeleteTestDetailsId})
			if cErr != nil {
				return errors.New(cErr.Message)
			}

			cErr = eventProcessor.SampleService.DeleteSamplesAndTestDetailsWithTx(ctx, tx, omsOrderId,
				[]models.TestDetail{toDeleteTestDetail})
			if cErr != nil {
				return errors.New(cErr.Message)
			}
		}

		if task.Id != 0 {
			labIdLabMap := eventProcessor.CdsService.GetLabIdLabMap(ctx)
			task, deleteTask := eventProcessor.updateTaskBasedOnTestDetailsDeletion(task, updatedTestDetails, labIdLabMap)
			if deleteTask {
				cErr = eventProcessor.TaskService.DeleteTaskWithTx(tx, task.Id)
				if cErr != nil {
					return errors.New(cErr.Message)
				}
			} else {
				_, cErr = eventProcessor.TaskService.UpdateTaskWithTx(tx, task)
				if cErr != nil {
					return errors.New(cErr.Message)
				}
			}
		}

		return nil
	})

	if err != nil {
		utils.AddLog(ctx, constants.ERROR_LEVEL, utils.GetCurrentFunctionName(), nil, err)
		return err
	}

	go func() {
		if task.Id != 0 {
			cErr = eventProcessor.ReleaseReport(ctx, task.Id)
			if cErr != nil {
				utils.AddLog(ctx, constants.ERROR_LEVEL, utils.GetCurrentFunctionName(), nil, errors.New(cErr.Message))
			}
		}

		eventProcessor.EtsService.GetAndPublishEtsTestBasicEvent(ctx, []string{omsTestDeleteEvent.AlnumTestId})
	}()

	return nil
}

func (eventProcessor *EventProcessor) updateTaskBasedOnTestDetailsDeletion(task models.Task,
	allTestDetails []models.TestDetail, labIdLabMap map[uint]structures.Lab) (models.Task, bool) {
	if len(allTestDetails) == 0 {
		return task, true
	}

	allInhouseTestsCompleted := true
	testCompleteStatuses := []string{
		constants.TEST_STATUS_APPROVE,
		constants.TEST_STATUS_COMPLETED_NOT_SENT,
		constants.TEST_STATUS_COMPLETED_SENT,
		constants.TEST_STATUS_SAMPLE_NOT_RECEIVED,
	}
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

	doctorTat := getDoctorTatForTask(allTestDetails)
	if !doctorTat.IsZero() {
		task.DoctorTat = doctorTat
	}

	if task.Status == constants.TASK_STATUS_IN_PROGRESS {
		task.Status = task.PreviousStatus
	}

	return task, false
}
