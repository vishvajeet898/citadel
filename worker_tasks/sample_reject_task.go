package workerTasks

import (
	"context"
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"

	"github.com/Orange-Health/citadel/common/constants"
	"github.com/Orange-Health/citadel/common/utils"
	"github.com/Orange-Health/citadel/models"
	"github.com/RichardKnop/machinery/v1/log"
	"github.com/RichardKnop/machinery/v1/tasks"
)

func (wt *WorkerTaskService) SampleRejectionTask(omsOrderId string, testIds []string) error {
	log.INFO.Print("SampleRejectionTask :: Started")
	ctx := context.Background()

	if omsOrderId == "" || len(testIds) == 0 {
		log.INFO.Print("SampleRejectionTask :: Invalid input")
		return nil
	}

	task, cErr := wt.TaskService.GetTaskByOmsOrderId(omsOrderId)
	if cErr != nil {
		return nil
	}

	testDetails, cErr := wt.TestDetailService.GetTestDetailsByOmsOrderId(omsOrderId)
	if cErr != nil {
		err := errors.New(cErr.Message)
		utils.AddLog(ctx, constants.ERROR_LEVEL, utils.GetCurrentFunctionName(), nil, err)
	}

	testDetailsToBeUpdated := []models.TestDetail{}
	for index := range testDetails {
		if utils.SliceContainsString(testIds, testDetails[index].CentralOmsTestId) {
			testDetailsToBeUpdated = append(testDetailsToBeUpdated, testDetails[index])
		}
	}

	if len(testDetailsToBeUpdated) == 0 {
		return nil
	}

	for index := range testDetailsToBeUpdated {
		testDetailsToBeUpdated[index].Status = constants.TEST_STATUS_REJECTED
		testDetailsToBeUpdated[index].ProcessingLabId = testDetailsToBeUpdated[index].LabId
	}

	err := wt.Db.Transaction(func(tx *gorm.DB) error {
		if task.Id != 0 {
			cErr = wt.TaskPathMappingService.MarkAllTaskPathMapInactiveByTaskIdWithTx(tx, task.Id)
			if cErr != nil {
				return errors.New(cErr.Message)
			}
		}

		_, cErr = wt.TestDetailService.UpdateTestDetailsWithTx(tx, testDetailsToBeUpdated)
		if cErr != nil {
			return errors.New(cErr.Message)
		}

		return nil
	})

	if err != nil {
		utils.AddLog(ctx, constants.ERROR_LEVEL, utils.GetCurrentFunctionName(), nil, err)
		return err
	}

	log.INFO.Print("SampleRejectionTask :: Ended")
	return nil
}

func SampleRejectionTaskSignature(omsOrderId string, testIds []string) *tasks.Signature {
	groupId := fmt.Sprintf("%s:%v", constants.SampleRejectionTask, time.Now().Unix())
	return &tasks.Signature{
		Name: constants.SampleRejectionTask,
		Args: []tasks.Arg{
			{
				Type:  "string",
				Value: omsOrderId,
				Name:  "omsOrderId",
			},
			{
				Type:  "[]string",
				Value: testIds,
				Name:  "testIds",
			},
		},
		BrokerMessageGroupId: groupId,
	}
}
