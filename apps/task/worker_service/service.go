package workerService

import (
	"context"
	"fmt"

	commonConstants "github.com/Orange-Health/citadel/common/constants"
	commonUtils "github.com/Orange-Health/citadel/common/utils"
	"github.com/Orange-Health/citadel/worker"
	workerTasks "github.com/Orange-Health/citadel/worker_tasks"
)

type TaskWorkerServiceInterface interface {
	ReleaseReportTask(ctx context.Context, taskId uint, forcefulUpdate bool) error
	UpdateTaskPostSavingTask(ctx context.Context, taskId uint) error
	CreateUpdateTaskByOmsOrderIdTask(ctx context.Context, omsOrderId string) error
}

func (ws *TaskWorkerService) ReleaseReportTask(ctx context.Context, taskId uint, forcefulUpdate bool) error {
	server, err := worker.StartServer(true)
	if err != nil {
		commonUtils.AddLog(ctx, commonConstants.ERROR_LEVEL, commonUtils.GetCurrentFunctionName(),
			map[string]interface{}{
				"error_message": commonConstants.ERROR_FAILED_TO_START_SERVER,
			}, err)
		return err
	}

	reportReleaseKey := fmt.Sprintf(commonConstants.ReportReleaseKey, taskId)
	err = ws.Cache.Set(ctx, reportReleaseKey, true, commonConstants.CacheExpiry1DayInt)
	if err != nil {
		commonUtils.AddLog(ctx, commonConstants.ERROR_LEVEL, commonUtils.GetCurrentFunctionName(),
			map[string]interface{}{
				"error_message": commonConstants.ERROR_FAILED_TO_SET_REPORT_RELEASE_KEY,
			}, err)
		return err
	}

	addTaskSignature := workerTasks.ReleaseReportSignature(taskId, reportReleaseKey, forcefulUpdate)
	_, err = server.SendTask(addTaskSignature)
	if err != nil {
		commonUtils.AddLog(ctx, commonConstants.ERROR_LEVEL, commonUtils.GetCurrentFunctionName(),
			map[string]interface{}{
				"error_message": commonConstants.ERROR_FAILED_TO_SEND_TASK,
			}, err)
	}
	return err
}

func (ws *TaskWorkerService) UpdateTaskPostSavingTask(ctx context.Context, taskId uint) error {
	server, err := worker.StartServer(true)
	if err != nil {
		commonUtils.AddLog(ctx, commonConstants.ERROR_LEVEL, commonUtils.GetCurrentFunctionName(),
			map[string]interface{}{
				"error_message": commonConstants.ERROR_FAILED_TO_START_SERVER,
			}, err)
		return err
	}

	addTaskSignature := workerTasks.UpdateTaskPostSavingTaskSignature(taskId)
	_, err = server.SendTask(addTaskSignature)
	if err != nil {
		commonUtils.AddLog(ctx, commonConstants.ERROR_LEVEL, commonUtils.GetCurrentFunctionName(),
			map[string]interface{}{
				"error_message": commonConstants.ERROR_FAILED_TO_SEND_TASK,
			}, err)
	}
	return err
}

func (ws *TaskWorkerService) CreateUpdateTaskByOmsOrderIdTask(ctx context.Context, omsOrderId string) error {
	server, err := worker.StartServer(true)
	if err != nil {
		commonUtils.AddLog(ctx, commonConstants.ERROR_LEVEL, commonUtils.GetCurrentFunctionName(),
			map[string]interface{}{
				"error_message": commonConstants.ERROR_FAILED_TO_START_SERVER,
			}, err)
		return err
	}

	addTaskSignature := workerTasks.CreateUpdateTaskByOmsOrderIdTaskSignature(omsOrderId)
	_, err = server.SendTask(addTaskSignature)
	if err != nil {
		commonUtils.AddLog(ctx, commonConstants.ERROR_LEVEL, commonUtils.GetCurrentFunctionName(),
			map[string]interface{}{
				"error_message": commonConstants.ERROR_FAILED_TO_SEND_TASK,
			}, err)
	}
	return err
}
