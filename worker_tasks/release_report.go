package workerTasks

import (
	"context"
	"fmt"
	"time"

	"github.com/RichardKnop/machinery/v1/log"
	"github.com/RichardKnop/machinery/v1/tasks"

	"github.com/Orange-Health/citadel/common/constants"
	"github.com/Orange-Health/citadel/common/utils"
)

func (wt *WorkerTaskService) ReleaseReportTask(taskId uint, reportReleaseKey string, forcefulUpdate bool) error {
	log.INFO.Print("ReleaseReportTask :: Started")

	ctx := context.Background()
	utils.AddLog(ctx, constants.DEBUG_LEVEL, utils.GetCurrentFunctionName(), map[string]interface{}{
		"taskId":           taskId,
		"reportReleaseKey": reportReleaseKey,
	}, nil)

	releaseReportKeyExists, err := wt.Cache.Exists(ctx, reportReleaseKey)
	if err != nil {
		utils.AddLog(ctx, constants.DEBUG_LEVEL, utils.GetCurrentFunctionName(), map[string]interface{}{
			"message": "some error",
		}, err)
	}

	if !releaseReportKeyExists {
		utils.AddLog(ctx, constants.DEBUG_LEVEL, utils.GetCurrentFunctionName(), map[string]interface{}{
			"message": "release report key not found",
		}, nil)
		return nil
	}

	err = wt.CommonTaskProcessor.ReleaseReportTask(ctx, taskId, forcefulUpdate)
	if err != nil {
		utils.AddLog(ctx, constants.DEBUG_LEVEL, utils.GetCurrentFunctionName(), nil, err)
	}

	err = wt.Cache.Delete(ctx, reportReleaseKey)
	if err != nil {
		utils.AddLog(ctx, constants.DEBUG_LEVEL, utils.GetCurrentFunctionName(), nil, err)
	}

	log.INFO.Print("ReleaseReportTask :: Ended")
	return nil
}

func ReleaseReportSignature(taskId uint, reportReleaseKey string, forcefulUpdate bool) *tasks.Signature {
	eta := time.Now().Add(5 * time.Second)
	groupId := fmt.Sprintf("%s:%v", constants.ReleaseReportTask, time.Now().Unix())
	return &tasks.Signature{
		Name: constants.ReleaseReportTask,
		Args: []tasks.Arg{
			{
				Type:  "uint",
				Value: taskId,
			},
			{
				Type:  "string",
				Value: reportReleaseKey,
			},
			{
				Type:  "bool",
				Value: forcefulUpdate,
			},
		},
		ETA:                  &eta,
		BrokerMessageGroupId: groupId,
	}
}
