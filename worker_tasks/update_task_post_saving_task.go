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

func (wt *WorkerTaskService) UpdateTaskPostSavingTask(taskId uint) error {
	log.INFO.Print("UpdateTaskPostSavingTask :: Started")

	ctx := context.Background()
	utils.AddLog(ctx, constants.DEBUG_LEVEL, utils.GetCurrentFunctionName(), map[string]interface{}{
		"taskId": taskId,
	}, nil)

	redisKey := fmt.Sprintf(constants.UpdateTaskPostSavingTaskKey, taskId)
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

	completeMasterTestIds, incompleteMasterTestIds, totalMasterTestIds := []uint{}, []uint{}, []uint{}
	dedupMasterTestIds := []uint{}

	task, taskMetadata, testDetails, testMetadaDetails, investigations, cErr :=
		wt.fetchDataForTaskUpdationAfterSavingTask(taskId)
	if cErr != nil {
		return errors.New(cErr.Message)
	}

	testDetailToCriticalMap, isTaskCritical := map[uint]bool{}, false
	for _, investigation := range investigations {
		if utils.SliceContainsString(constants.INVESTIGATION_STATUSES_APPROVE, investigation.InvestigationStatus) &&
			investigation.IsCritical {
			testDetailToCriticalMap[investigation.TestDetailsId] = true
			isTaskCritical = true
		}
	}

	updateTestDetailsMetadata := []models.TestDetailsMetadata{}
	for _, testMetadata := range testMetadaDetails {
		if testDetailToCriticalMap[testMetadata.TestDetailsId] {
			testMetadata.IsCritical = true
			updateTestDetailsMetadata = append(updateTestDetailsMetadata, testMetadata)
		}
	}

	taskMetadata.IsCritical = isTaskCritical

	taskDoctorTat := task.DoctorTat

	testDetailsStatuses := []string{
		constants.TEST_STATUS_RESULT_SAVED,
		constants.TEST_STATUS_RERUN_RESULT_SAVED,
		constants.TEST_STATUS_WITHHELD,
		constants.TEST_STATUS_CO_AUTHORIZE,
	}

	for _, testDetail := range testDetails {
		totalMasterTestIds = append(totalMasterTestIds, testDetail.MasterTestId)
		if testDetail.Status == constants.TEST_STATUS_APPROVE {
			completeMasterTestIds = append(completeMasterTestIds, testDetail.MasterTestId)
		} else if !testDetail.IsDuplicate {
			incompleteMasterTestIds = append(incompleteMasterTestIds, testDetail.MasterTestId)
		}
		if utils.SliceContainsString(testDetailsStatuses, testDetail.Status) {
			testDetailsTat := testDetail.DoctorTat
			if taskDoctorTat.After(*testDetailsTat) {
				taskDoctorTat = testDetail.DoctorTat
			}
		}
	}
	totalMasterTestIds = utils.CreateUniqueSliceUint(totalMasterTestIds)
	completeMasterTestIds = utils.CreateUniqueSliceUint(completeMasterTestIds)
	incompleteMasterTestIds = utils.CreateUniqueSliceUint(incompleteMasterTestIds)

	dedupResponse, err := wt.CdsService.GetDeduplicatedTestsAndPackages(ctx, totalMasterTestIds, []uint{})
	if err != nil {
		utils.AddLog(ctx, constants.ERROR_LEVEL, utils.GetCurrentFunctionName(), nil, err)
	}

	for _, removableTest := range dedupResponse.Recommendation.RemoveTests {
		if len(removableTest.CompleteOverlapDetails) > 0 {
			for _, completeOverlap := range removableTest.CompleteOverlapDetails {
				if utils.SliceContainsUint(completeMasterTestIds, completeOverlap.OverlappedTestId) &&
					utils.SliceContainsUint(incompleteMasterTestIds, removableTest.TestId) {
					dedupMasterTestIds = append(dedupMasterTestIds, removableTest.TestId)
					break
				}
			}
		}
	}
	for _, completeMasterTestId := range completeMasterTestIds {
		if utils.SliceContainsUint(incompleteMasterTestIds, completeMasterTestId) {
			dedupMasterTestIds = append(dedupMasterTestIds, completeMasterTestId)
		}
	}
	dedupMasterTestIds = utils.CreateUniqueSliceUint(dedupMasterTestIds)

	err = wt.Db.Transaction(func(tx *gorm.DB) error {
		if len(updateTestDetailsMetadata) > 0 {
			_, cErr := wt.TestDetailService.UpdateTestDetailsMetadataWithTx(tx, updateTestDetailsMetadata)
			if cErr != nil {
				return errors.New(cErr.Message)
			}
		}

		cErr := wt.TestDetailService.UpdateDuplicateTestDetailsByTaskIdWithTx(tx, taskId, dedupMasterTestIds)
		if cErr != nil {
			return errors.New(cErr.Message)
		}

		_, cErr = wt.TaskService.UpdateTaskMetadataWithTx(tx, taskMetadata)
		if cErr != nil {
			return errors.New(cErr.Message)
		}

		_, cErr = wt.TaskService.UpdateTaskWithTx(tx, task)
		if cErr != nil {
			return errors.New(cErr.Message)
		}

		return nil
	})

	if err != nil {
		utils.AddLog(ctx, constants.ERROR_LEVEL, utils.GetCurrentFunctionName(), nil, err)
		return err
	}

	wt.CommonTaskProcessor.CheckForTaskCompletion(ctx, taskId)

	log.INFO.Print("UpdateTaskPostSavingTask :: Ended")
	return nil
}

func (wt *WorkerTaskService) fetchDataForTaskUpdationAfterSavingTask(taskId uint) (
	task models.Task, taskMetadata models.TaskMetadata, testDetails []models.TestDetail,
	testDetailsMetadata []models.TestDetailsMetadata, investigationResults []models.InvestigationResult,
	cErr *structures.CommonError,
) {
	task, cErr = wt.TaskService.GetTaskModelById(taskId)
	if cErr != nil {
		return
	}

	taskMetadata, cErr = wt.TaskService.GetTaskMetadataByTaskId(taskId)
	if cErr != nil {
		return
	}

	testDetails, cErr = wt.TestDetailService.GetTestDetailsByOmsOrderId(task.OmsOrderId)
	if cErr != nil {
		return
	}

	testDetailsIds := []uint{}
	for _, testDetail := range testDetails {
		testDetailsIds = append(testDetailsIds, testDetail.Id)
	}

	testDetailsMetadata, cErr = wt.TestDetailService.GetTestDetailsMetadataByTestDetailIds(testDetailsIds)
	if cErr != nil {
		return
	}

	investigationResults, cErr = wt.InvestigationResultsService.GetInvestigationResultsByTestDetailsIds(testDetailsIds)
	if cErr != nil {
		return
	}

	return
}

func UpdateTaskPostSavingTaskSignature(taskId uint) *tasks.Signature {
	groupId := fmt.Sprintf("%s:%v", constants.UpdateTaskPostSavingTask, time.Now().Unix())
	return &tasks.Signature{
		Name: constants.UpdateTaskPostSavingTask,
		Args: []tasks.Arg{
			{
				Type:  "uint",
				Value: taskId,
			},
		},
		BrokerMessageGroupId: groupId,
	}
}
