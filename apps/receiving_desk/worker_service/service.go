package workerService

import (
	"context"

	commonConstants "github.com/Orange-Health/citadel/common/constants"
	commonUtils "github.com/Orange-Health/citadel/common/utils"
	"github.com/Orange-Health/citadel/worker"
	workerTasks "github.com/Orange-Health/citadel/worker_tasks"
)

type ReceivingDeskWorkerServiceInterface interface {
	UpdateTaskAndTaskMetadataPostReceiving(ctx context.Context, omsOrderId string) error
}

func (ws *ReceivingDeskWorkerService) UpdateTaskAndTaskMetadataPostReceiving(ctx context.Context,
	omsOrderId string) error {
	server, err := worker.StartServer(true)
	if err != nil {
		commonUtils.AddLog(ctx, commonConstants.ERROR_LEVEL, commonUtils.GetCurrentFunctionName(),
			map[string]interface{}{
				"error_message": commonConstants.ERROR_FAILED_TO_START_SERVER,
			}, err)
		return err
	}

	addTaskSignature := workerTasks.UpdateTaskPostReceivingTaskSignature(omsOrderId)
	_, err = server.SendTask(addTaskSignature)
	if err != nil {
		commonUtils.AddLog(ctx, commonConstants.ERROR_LEVEL, commonUtils.GetCurrentFunctionName(),
			map[string]interface{}{
				"error_message": commonConstants.ERROR_FAILED_TO_SEND_TASK,
			}, err)
	}
	return err
}
