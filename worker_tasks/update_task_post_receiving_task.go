package workerTasks

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Orange-Health/citadel/common/constants"
	"github.com/Orange-Health/citadel/common/structures"
	"github.com/Orange-Health/citadel/common/utils"
	"github.com/Orange-Health/citadel/models"
	"github.com/RichardKnop/machinery/v1/log"
	"github.com/RichardKnop/machinery/v1/tasks"
	"gorm.io/gorm"
)

func (wt *WorkerTaskService) UpdateTaskPostReceivingTask(omsOrderId string) error {
	log.INFO.Print("UpdateTaskPostReceivingTask :: Started")

	ctx := context.Background()
	utils.AddLog(ctx, constants.DEBUG_LEVEL, utils.GetCurrentFunctionName(), map[string]interface{}{
		"omsOrderId": omsOrderId,
	}, nil)

	orderDetails, task, taskMetadata, cErr := wt.fetchDataForTaskUpdationPostReceiving(omsOrderId)
	if cErr != nil {
		return errors.New(cErr.Message)
	}

	orderDetailsResponse, err := wt.OmsClient.GetOrderById(ctx, omsOrderId, orderDetails.CityCode)
	if err != nil {
		utils.AddLog(ctx, constants.ERROR_LEVEL, utils.GetCurrentFunctionName(), nil, err)
	}

	if orderDetailsResponse.Id != 0 {
		orderDetails.PartnerId = orderDetailsResponse.PartnerId
		orderDetails.DoctorId = orderDetailsResponse.SystemDoctorId
		orderDetails.ReferredBy = orderDetailsResponse.ReferredBy
		if task.Id != 0 {
			taskMetadata.DoctorNotes = orderDetailsResponse.Notes
		}
	}

	doctor, partner := structures.Doctor{}, structures.Partner{}
	if orderDetails.DoctorId != 0 && task.Id != 0 {
		doctor, _ = wt.HealthApiClient.GetDoctorById(ctx, orderDetails.DoctorId)
		taskMetadata.DoctorName = doctor.Name
		taskMetadata.DoctorNumber = doctor.Number
	}
	if orderDetails.PartnerId != 0 && task.Id != 0 {
		partner, _ = wt.PartnerApiClient.GetPartnerById(ctx, orderDetails.PartnerId)
		taskMetadata.PartnerName = partner.PartnerName
	}

	if task.Id != 0 {
		taskMetadata.ContainsPackage = wt.TestDetailService.ContainsPackageTests(omsOrderId)
	}

	txErr := wt.Db.Transaction(func(tx *gorm.DB) error {
		_, cErr = wt.OrderDetailsService.UpdateOrderDetailsWithTx(tx, orderDetails)
		if cErr != nil {
			return errors.New(cErr.Message)
		}

		if task.Id != 0 {
			_, cErr = wt.TaskService.UpdateTaskMetadataWithTx(tx, taskMetadata)
			if cErr != nil {
				return errors.New(cErr.Message)
			}
		}

		return nil
	})

	if txErr != nil {
		return txErr
	}

	log.INFO.Print("UpdateTaskPostReceivingTask :: Ended")
	return nil
}

func (wt *WorkerTaskService) fetchDataForTaskUpdationPostReceiving(omsOrderId string) (
	orderDetails models.OrderDetails, task models.Task, taskMetadata models.TaskMetadata, cErr *structures.CommonError) {

	orderDetails, cErr = wt.OrderDetailsService.GetOrderDetailsByOmsOrderId(omsOrderId)
	if cErr != nil {
		return
	}

	task, cErr = wt.TaskService.GetTaskByOmsOrderId(omsOrderId)
	if cErr != nil {
		return
	}

	if task.Id != 0 {
		taskMetadata, cErr = wt.TaskService.GetTaskMetadataByTaskId(task.Id)
		if cErr != nil {
			return
		}
	}

	return
}

func UpdateTaskPostReceivingTaskSignature(omsOrderId string) *tasks.Signature {
	groupId := fmt.Sprintf("%s:%v", constants.UpdateTaskPostReceivingTask, time.Now().Unix())
	return &tasks.Signature{
		Name: constants.UpdateTaskPostReceivingTask,
		Args: []tasks.Arg{
			{
				Type:  "string",
				Value: omsOrderId,
				Name:  "omsOrderId",
			},
		},
		BrokerMessageGroupId: groupId,
	}
}
