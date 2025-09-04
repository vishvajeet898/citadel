package workerTasks

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/RichardKnop/machinery/v1/log"
	"github.com/RichardKnop/machinery/v1/tasks"
	"gorm.io/gorm"

	"github.com/Orange-Health/citadel/common/constants"
	"github.com/Orange-Health/citadel/common/structures"
	"github.com/Orange-Health/citadel/common/utils"
	"github.com/Orange-Health/citadel/models"
)

func (wt *WorkerTaskService) CreateUpdateTaskByOmsOrderIdTask(omsOrderId string) error {
	log.INFO.Print("CreateUpdateTaskByOmsOrderIdTask :: Started")

	ctx := context.Background()
	utils.AddLog(ctx, constants.DEBUG_LEVEL, utils.GetCurrentFunctionName(), map[string]interface{}{
		"omsOrderId": omsOrderId,
	}, nil)

	redisKey := fmt.Sprintf(constants.CreateUpdateTaskTaskKey, omsOrderId)
	keyExists, err := wt.Cache.Exists(ctx, redisKey)
	if err != nil || keyExists {
		utils.AddLog(ctx, constants.DEBUG_LEVEL, utils.GetCurrentFunctionName(), nil, err)
		return errors.New("worker task in progress")
	}

	err = wt.Cache.Set(ctx, redisKey, true, constants.CacheExpiry1HourInt)
	if err != nil {
		utils.AddLog(ctx, constants.DEBUG_LEVEL, utils.GetCurrentFunctionName(), nil, err)
		return err
	}

	defer func() {
		err := wt.Cache.Delete(ctx, redisKey)
		if err != nil {
			utils.AddLog(ctx, constants.DEBUG_LEVEL, utils.GetCurrentFunctionName(), nil, err)
		}
	}()

	// Fetch task by OMS order ID
	orderDetails, task, taskMetadata, cErr := wt.fetchDataForTaskUpdationPostReceiving(omsOrderId)
	if cErr != nil {
		return errors.New(cErr.Message)
	}

	orderDetailsResponse, err := wt.OmsClient.GetOrderById(ctx, omsOrderId, orderDetails.CityCode)
	if err != nil {
		utils.AddLog(ctx, constants.ERROR_LEVEL, utils.GetCurrentFunctionName(), nil, err)
		return nil
	}
	if orderDetailsResponse.Id == 0 {
		utils.AddLog(ctx, constants.DEBUG_LEVEL, utils.GetCurrentFunctionName(), nil, nil)
		return nil
	}

	// Fetch all test_details for the given OMS order ID
	testDetails, cErr := wt.TestDetailService.GetTestDetailsByOmsOrderIdWithSampleStatus(omsOrderId, constants.SampleSynced)
	if cErr != nil {
		err := errors.New(cErr.Message)
		utils.AddLog(ctx, constants.ERROR_LEVEL, utils.GetCurrentFunctionName(), nil, err)
		return err
	}
	if len(testDetails) == 0 {
		utils.AddLog(ctx, constants.DEBUG_LEVEL, utils.GetCurrentFunctionName(), nil, nil)
		return nil
	}

	orderDetails.PartnerId = orderDetailsResponse.PartnerId
	orderDetails.DoctorId = orderDetailsResponse.SystemDoctorId
	orderDetails.ReferredBy = orderDetailsResponse.ReferredBy

	if task.Id == 0 {
		task = models.Task{
			OrderId:          utils.GetUintOrderIdWithoutStringPart(orderDetails.OmsOrderId),
			RequestId:        utils.GetUintRequestIdWithoutStringPart(orderDetails.OmsRequestId),
			OmsOrderId:       orderDetails.OmsOrderId,
			OmsRequestId:     orderDetails.OmsRequestId,
			LabId:            orderDetails.ServicingLabId,
			CityCode:         orderDetails.CityCode,
			Status:           constants.TASK_STATUS_PENDING,
			PreviousStatus:   constants.TASK_STATUS_PENDING,
			OrderType:        constants.CollecTypeToOrderTypeMap[orderDetails.CollectionType],
			PatientDetailsId: orderDetails.PatientDetailsId,
			IsActive:         true,
		}

		doctor, partner := structures.Doctor{}, structures.Partner{}
		if orderDetailsResponse.SystemDoctorId != 0 {
			doctor, _ = wt.HealthApiClient.GetDoctorById(ctx, orderDetailsResponse.SystemDoctorId)
		}
		if orderDetailsResponse.PartnerId != 0 {
			partner, _ = wt.PartnerApiClient.GetPartnerById(ctx, orderDetailsResponse.PartnerId)
		}

		taskMetadata.DoctorName = doctor.Name
		taskMetadata.DoctorNumber = doctor.Number
		taskMetadata.PartnerName = partner.PartnerName
		taskMetadata.DoctorNotes = orderDetailsResponse.Notes
		taskMetadata.ContainsPackage = wt.TestDetailService.ContainsPackageTests(omsOrderId)
	}

	labIdLabMap := wt.CdsService.GetLabIdLabMap(ctx)
	toBeUpdatedTestDetails := []models.TestDetail{}
	for _, testDetail := range testDetails {
		if labIdLabMap[testDetail.ProcessingLabId].Inhouse {
			toBeUpdatedTestDetails = append(toBeUpdatedTestDetails, testDetail)
		}
	}

	txErr := wt.Db.Transaction(func(tx *gorm.DB) error {
		_, cErr = wt.OrderDetailsService.UpdateOrderDetailsWithTx(tx, orderDetails)
		if cErr != nil {
			return errors.New(cErr.Message)
		}

		if task.Id == 0 {
			_, cErr = wt.TaskService.CreateTaskWithTx(tx, task)
			if cErr != nil {
				return errors.New(cErr.Message)
			}

			taskMetadata.TaskId = task.Id
			_, cErr = wt.TaskService.CreateTaskMetadataWithTx(tx, taskMetadata)
			if cErr != nil {
				return errors.New(cErr.Message)
			}
		} else {
			_, cErr = wt.TaskService.UpdateTaskWithTx(tx, task)
			if cErr != nil {
				return errors.New(cErr.Message)
			}
			taskMetadata.TaskId = task.Id
			_, cErr = wt.TaskService.UpdateTaskMetadataWithTx(tx, taskMetadata)
			if cErr != nil {
				return errors.New(cErr.Message)
			}
		}

		if len(toBeUpdatedTestDetails) > 0 {
			for index := range toBeUpdatedTestDetails {
				toBeUpdatedTestDetails[index].TaskId = task.Id
			}
			_, cErr = wt.TestDetailService.UpdateTestDetailsWithTx(tx, toBeUpdatedTestDetails)
			if cErr != nil {
				return errors.New(cErr.Message)
			}
		}

		return nil
	})

	if txErr != nil {
		utils.AddLog(ctx, constants.ERROR_LEVEL, utils.GetCurrentFunctionName(), nil, txErr)
		return txErr
	}

	log.INFO.Print("CreateUpdateTaskByOmsOrderIdTask :: Ended")
	return nil
}

func CreateUpdateTaskByOmsOrderIdTaskSignature(omsOrderId string) *tasks.Signature {
	groupId := fmt.Sprintf("%s:%v", constants.CreateUpdateTaskByOmsOrderIdTask, time.Now().Unix())
	return &tasks.Signature{
		Name: constants.CreateUpdateTaskByOmsOrderIdTask,
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
