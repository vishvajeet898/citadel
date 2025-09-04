package commonTasks

import (
	"context"
	"errors"

	"github.com/Orange-Health/citadel/common/constants"
	"github.com/Orange-Health/citadel/common/utils"
)

func (ctp *CommonTaskProcessor) CheckForTaskCompletion(ctx context.Context, taskId uint) {
	labIdLabMap := ctp.CdsService.GetLabIdLabMap(ctx)

	task, cErr := ctp.TaskService.GetTaskModelById(taskId)
	if cErr != nil {
		return
	}

	testDetails, cErr := ctp.TestDetailService.GetTestDetailsByOmsOrderId(task.OmsOrderId)
	if cErr != nil {
		return
	}

	allInhouseTestsCompleted := true
	testCompleteStatuses := []string{
		constants.TEST_STATUS_APPROVE,
		constants.TEST_STATUS_COMPLETED_NOT_SENT,
		constants.TEST_STATUS_COMPLETED_SENT,
		constants.TEST_STATUS_SAMPLE_NOT_RECEIVED,
	}
	for _, testDetail := range testDetails {
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

		// Update task in the database
		_, cErr = ctp.TaskService.UpdateTask(task)
		if cErr != nil {
			utils.AddLog(ctx, constants.ERROR_LEVEL, utils.GetCurrentFunctionName(), nil, errors.New(cErr.Message))
		}
	}
}
