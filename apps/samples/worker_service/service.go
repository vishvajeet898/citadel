package workerService

import (
	"context"

	commonConstants "github.com/Orange-Health/citadel/common/constants"
	commonUtils "github.com/Orange-Health/citadel/common/utils"
	"github.com/Orange-Health/citadel/worker"
	workerTasks "github.com/Orange-Health/citadel/worker_tasks"
)

type SampleWorkerServiceInterface interface {
	SampleRejectionTask(ctx context.Context, omsOrderId string, testIds []string) error
}

func (ws *SampleWorkerService) SampleRejectionTask(ctx context.Context, omsOrderId string, testIds []string) error {
	server, err := worker.StartServer(true)
	if err != nil {
		commonUtils.AddLog(ctx, commonConstants.ERROR_LEVEL, commonUtils.GetCurrentFunctionName(),
			map[string]interface{}{
				"error_message": commonConstants.ERROR_FAILED_TO_START_SERVER,
			}, err)
		return err
	}

	addTaskSignature := workerTasks.SampleRejectionTaskSignature(omsOrderId, testIds)
	_, err = server.SendTask(addTaskSignature)
	if err != nil {
		commonUtils.AddLog(ctx, commonConstants.ERROR_LEVEL, commonUtils.GetCurrentFunctionName(),
			map[string]interface{}{
				"error_message": commonConstants.ERROR_FAILED_TO_SEND_TASK,
			}, err)
	}
	return err
}
