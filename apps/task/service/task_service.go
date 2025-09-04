package service

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"gorm.io/gorm"

	mapper "github.com/Orange-Health/citadel/apps/task/mapper"
	"github.com/Orange-Health/citadel/apps/task/structures"
	"github.com/Orange-Health/citadel/common/constants"
	commonConstants "github.com/Orange-Health/citadel/common/constants"
	commonStructures "github.com/Orange-Health/citadel/common/structures"
	commonUtils "github.com/Orange-Health/citadel/common/utils"
	commonModels "github.com/Orange-Health/citadel/models"
)

func (taskService *TaskService) GetTaskById(taskID uint) (structures.TaskDetail, *commonStructures.CommonError) {

	task, cErr := taskService.TaskDao.GetTaskWithPatientDetailsById(taskID)

	if cErr != nil {
		return structures.TaskDetail{}, cErr
	}

	taskMetadata, cErr := taskService.TaskDao.GetTaskMetadataByTaskId(task.Id)
	if cErr != nil {
		return structures.TaskDetail{}, cErr
	}

	visits, cErr := taskService.SampleService.GetVisitDetailsForTaskByOmsOrderId(task.OmsOrderId)
	if cErr != nil {
		return structures.TaskDetail{}, cErr
	}

	go taskService.TestDetailService.UpdatePickedAtTimeBasedOnActiveTests(taskID)

	return mapper.GetTaskDetails(task, taskMetadata, visits)
}

func (taskService *TaskService) GetTaskModelById(taskID uint) (commonModels.Task,
	*commonStructures.CommonError) {
	return taskService.TaskDao.GetTaskById(taskID)
}

func (taskService *TaskService) GetTaskByOmsOrderId(omsOrderId string) (commonModels.Task, *commonStructures.CommonError) {
	return taskService.TaskDao.GetTaskByOmsOrderId(omsOrderId)
}

func (taskService *TaskService) GetTasksCount() (uint, *commonStructures.CommonError) {
	return taskService.TaskDao.GetTasksCount()
}

func (taskService *TaskService) GetAmendmentTasksCount(ctx context.Context) (uint, *commonStructures.CommonError) {
	amendmentTaskCountKey := commonConstants.AmendmendTaskCountKey
	var cacheResponse interface{}

	_ = taskService.Cache.Get(ctx, amendmentTaskCountKey, &cacheResponse)
	if cacheResponse != nil {
		countFloat := cacheResponse.(float64)
		return uint(countFloat), nil
	}

	count, cErr := taskService.TaskDao.GetAmendmentTasksCount()
	if cErr == nil {
		_ = taskService.Cache.Set(ctx, amendmentTaskCountKey, count, commonConstants.CacheExpiry5MinutesInt)
	}

	return count, cErr
}

func (taskService *TaskService) UpdateTask(task commonModels.Task) (commonModels.Task, *commonStructures.CommonError) {
	return taskService.TaskDao.UpdateTask(task)
}

func (taskService *TaskService) CreateTaskWithTx(tx *gorm.DB, task commonModels.Task) (commonModels.Task, *commonStructures.CommonError) {
	return taskService.TaskDao.CreateTaskWithTx(tx, task)
}

func (taskService *TaskService) UpdateTaskWithTx(tx *gorm.DB, task commonModels.Task) (commonModels.Task, *commonStructures.CommonError) {
	return taskService.TaskDao.UpdateTaskWithTx(tx, task)
}

func (taskService *TaskService) DeleteTaskWithTx(tx *gorm.DB, taskID uint) *commonStructures.CommonError {
	return taskService.TaskDao.DeleteTaskWithTx(tx, taskID)
}

func (taskService *TaskService) fetchAllTaskData(taskId, userId uint) (
	[]commonModels.TestDetail, []commonModels.InvestigationResult, commonModels.User, *commonStructures.CommonError) {

	user := commonModels.User{}
	task, cErr := taskService.TaskDao.GetTaskById(taskId)
	if cErr != nil {
		return nil, nil, user, cErr
	}

	testDetails, cErr := taskService.TestDetailService.GetTestDetailsByOmsOrderId(task.OmsOrderId)
	if cErr != nil {
		return nil, nil, user, cErr
	}

	testDetailIds := []uint{}
	for _, testDetail := range testDetails {
		testDetailIds = append(testDetailIds, testDetail.Id)
	}

	investigations, cErr := taskService.InvestigationResultsService.GetInvestigationResultsByTestDetailsIds(testDetailIds)
	if cErr != nil {
		return nil, nil, user, cErr
	}

	user, cErr = taskService.UserService.GetUserModel(userId)
	if cErr != nil {
		return nil, nil, user, cErr
	}

	return testDetails, investigations, user, nil
}

func (taskService *TaskService) fetchAllUpdateTaskData(taskId uint, testDetailIds,
	investigationIds, remarkIds []uint) (commonModels.Task, []commonModels.TestDetail,
	[]commonModels.InvestigationResult, []commonModels.InvestigationData, []commonModels.Remark,
	*commonStructures.CommonError) {

	remarks := []commonModels.Remark{}

	task, cErr := taskService.TaskDao.GetTaskById(taskId)
	if cErr != nil {
		return commonModels.Task{}, nil, nil, nil, nil, cErr
	}

	testDetails, cErr := taskService.TestDetailService.GetTestDetailsByTestIds(testDetailIds)
	if cErr != nil {
		return commonModels.Task{}, nil, nil, nil, nil, cErr
	}

	investigations, cErr := taskService.InvestigationResultsService.GetInvestigationsByInvestigationIds(investigationIds)
	if cErr != nil {
		return commonModels.Task{}, nil, nil, nil, nil, cErr
	}

	investigationsData, cErr := taskService.InvestigationResultsService.GetInvestigationDataByInvestigationResultsIds(
		investigationIds)
	if cErr != nil {
		return commonModels.Task{}, nil, nil, nil, nil, cErr
	}

	if len(remarkIds) != 0 {
		remarks, cErr = taskService.RemarkService.GetRemarksByRemarkIds(remarkIds)
		if cErr != nil {
			return commonModels.Task{}, nil, nil, nil, nil, cErr
		}
	}

	return task, testDetails, investigations, investigationsData, remarks, nil
}

func (taskService *TaskService) validateTestStatusesAndGetTaskStatus(newTestDetails,
	oldTestDetails []commonModels.TestDetail) (
	string, *commonStructures.CommonError) {
	newTestDetailIdToStatusMap := make(map[uint]string)
	testDetailStatuses, uniqueTestStatuses := []string{}, []string{}

	for _, testDetail := range newTestDetails {
		testDetailStatuses = append(testDetailStatuses, testDetail.Status)
		if !commonUtils.SliceContainsString(commonConstants.TEST_STATUSES, testDetail.Status) {
			return "", &commonStructures.CommonError{
				Message:    commonConstants.ERROR_INVALID_TEST_STATUS,
				StatusCode: http.StatusBadRequest,
			}
		}
		newTestDetailIdToStatusMap[testDetail.Id] = testDetail.Status
	}

	for _, testDetail := range oldTestDetails {
		if _, ok := newTestDetailIdToStatusMap[testDetail.Id]; !ok {
			testDetailStatuses = append(testDetailStatuses, testDetail.Status)
		}
	}

	testDetailStatusMap := make(map[string]bool)
	for _, testDetailStatus := range testDetailStatuses {
		testDetailStatusMap[testDetailStatus] = true
	}

	for testDetailStatus := range testDetailStatusMap {
		uniqueTestStatuses = append(uniqueTestStatuses, testDetailStatus)
	}

	if testDetailStatusMap[commonConstants.TEST_STATUS_WITHHELD] &&
		testDetailStatusMap[commonConstants.TEST_STATUS_CO_AUTHORIZE] {
		return "", &commonStructures.CommonError{
			Message:    commonConstants.ERROR_TEST_STATUS_WITHHELD_CO_AUTHORIZE_NON_COUPLING,
			StatusCode: http.StatusBadRequest,
		}
	}

	switch {
	case testDetailStatusMap[commonConstants.TEST_STATUS_CO_AUTHORIZE]:
		return commonConstants.TASK_STATUS_CO_AUTHORIZE, nil
	case testDetailStatusMap[commonConstants.TEST_STATUS_WITHHELD]:
		return commonConstants.TASK_STATUS_WITHHELD_APPROVAL, nil
	case testDetailStatusMap[commonConstants.TEST_STATUS_RERUN_REQUESTED]:
		return commonConstants.TASK_STATUS_PENDING, nil
	case len(uniqueTestStatuses) == 2 && testDetailStatusMap[commonConstants.TEST_STATUS_RESULT_SAVED] &&
		testDetailStatusMap[commonConstants.TEST_STATUS_RERUN_RESULT_SAVED]:
		return commonConstants.TASK_STATUS_PENDING, nil
	case len(uniqueTestStatuses) == 1 && testDetailStatusMap[commonConstants.TEST_STATUS_APPROVE]:
		return commonConstants.TASK_STATUS_COMPLETED, nil
	}

	return commonConstants.TASK_STATUS_PENDING, nil
}

func (taskService *TaskService) validateInvestigations(oldInvestigations,
	newInvestigations []commonModels.InvestigationResult,
	newTestDetails []commonModels.TestDetail) *commonStructures.CommonError {

	newTestDetailIdToInvestigationMap, oldTestDetailsIdToInvestigationMap :=
		make(map[uint][]commonModels.InvestigationResult), make(map[uint][]commonModels.InvestigationResult)

	for _, investigation := range newInvestigations {
		if !commonUtils.SliceContainsString(commonConstants.INVESTIGATION_RESULT_STATUSES,
			investigation.InvestigationStatus) {
			return &commonStructures.CommonError{
				Message:    commonConstants.ERROR_INVALID_INVESTIGATION_STATUS,
				StatusCode: http.StatusBadRequest,
			}
		}
	}

	for _, investigation := range newInvestigations {
		if _, ok := newTestDetailIdToInvestigationMap[investigation.TestDetailsId]; !ok {
			newTestDetailIdToInvestigationMap[investigation.TestDetailsId] = []commonModels.InvestigationResult{}
		}
		newTestDetailIdToInvestigationMap[investigation.TestDetailsId] =
			append(newTestDetailIdToInvestigationMap[investigation.TestDetailsId], investigation)
	}

	for _, investigation := range oldInvestigations {
		if _, ok := oldTestDetailsIdToInvestigationMap[investigation.TestDetailsId]; !ok {
			oldTestDetailsIdToInvestigationMap[investigation.TestDetailsId] = []commonModels.InvestigationResult{}
		}
		oldTestDetailsIdToInvestigationMap[investigation.TestDetailsId] =
			append(oldTestDetailsIdToInvestigationMap[investigation.TestDetailsId], investigation)
	}

	for newTestId, newInvestigations := range newTestDetailIdToInvestigationMap {
		oldInvestigations := oldTestDetailsIdToInvestigationMap[newTestId]
		oldInvestigationIds, newInvestigationIds := []uint{}, []uint{}
		for _, oldInvestigation := range oldInvestigations {
			oldInvestigationIds = append(oldInvestigationIds, oldInvestigation.Id)
		}

		for _, newInvestigation := range newInvestigations {
			newInvestigationIds = append(newInvestigationIds, newInvestigation.Id)
		}

		if !commonUtils.AreEqualUintSlices(oldInvestigationIds, newInvestigationIds) {
			return &commonStructures.CommonError{
				Message:    commonConstants.ERROR_INVESTIGATION_IDS_MISMATCH,
				StatusCode: http.StatusBadRequest,
			}
		}
	}

	for index := range newTestDetails {
		investigations := newTestDetailIdToInvestigationMap[newTestDetails[index].Id]
		investigationStatus := investigations[0].InvestigationStatus
		ifInvestigationRerun := false
		for _, investigation := range investigations {
			if commonUtils.SliceContainsString(commonConstants.INVESTIGATION_STATUSES_RERUN, investigation.InvestigationStatus) {
				ifInvestigationRerun = true
				break
			}
		}
		if ifInvestigationRerun {
			newTestDetails[index].Status = commonConstants.TEST_STATUS_RERUN_REQUESTED
			for _, investigation := range investigations {
				if !commonUtils.SliceContainsString(commonConstants.INVESTIGATION_STATUSES_RERUN, investigation.InvestigationStatus) &&
					!commonUtils.SliceContainsString(commonConstants.INVESTIGATION_STATUSES_APPROVE, investigation.InvestigationStatus) {
					return &commonStructures.CommonError{
						Message:    commonConstants.ERROR_RERUN_STATUS_COUPLING,
						StatusCode: http.StatusBadRequest,
					}
				}
				if commonUtils.SliceContainsString(commonConstants.INVESTIGATION_STATUSES_APPROVE, investigation.InvestigationStatus) {
					investigation.InvestigationStatus = commonConstants.INVESTIGATION_STATUS_PENDING
					investigation.ApprovedAt = nil
					investigation.ApprovedBy = 0
					for index := range newInvestigations {
						if newInvestigations[index].Id == investigation.Id {
							newInvestigations[index] = investigation
							break
						}
					}
				}
			}
		} else {
			for _, investigation := range investigations {
				if investigation.InvestigationStatus != investigationStatus {
					return &commonStructures.CommonError{
						Message:    commonConstants.ERROR_INVESTIGATION_STATUS_MISMATCH,
						StatusCode: http.StatusBadRequest,
					}
				}
			}
			newTestDetails[index].Status = investigationStatus
			if newTestDetails[index].Status == commonConstants.TEST_STATUS_APPROVE {
				newTestDetails[index].ApprovalSource = commonConstants.TEST_APPROVAL_SOURCE_OH
			}
		}
	}

	return nil
}

func (taskService *TaskService) validateUpdateTaskDetailsAndUpdateData(task *commonModels.Task,
	oldTestDetails []commonModels.TestDetail,
	oldInvestigations []commonModels.InvestigationResult,
	newTestDetails []commonModels.TestDetail,
	newInvestigations []commonModels.InvestigationResult) *commonStructures.CommonError {

	cErr := taskService.validateInvestigations(oldInvestigations, newInvestigations, newTestDetails)
	if cErr != nil {
		return cErr
	}

	taskStatus, cErr := taskService.validateTestStatusesAndGetTaskStatus(newTestDetails, oldTestDetails)
	if cErr != nil {
		return cErr
	}

	task.Status = taskStatus
	if taskStatus == commonConstants.TASK_STATUS_COMPLETED {
		task.CompletedAt = commonUtils.GetCurrentTime()
	}

	return nil
}

func (taskService *TaskService) getTestDetailsToBeRerun(testDetails []commonModels.TestDetail) ([]uint, []string) {
	testDetailsIdsToBeRerun, omsTestIdsToRerun := []uint{}, []string{}

	for _, testDetail := range testDetails {
		if testDetail.Status == commonConstants.TEST_STATUS_RERUN_REQUESTED {
			testDetailsIdsToBeRerun = append(testDetailsIdsToBeRerun, testDetail.Id)
			omsTestIdsToRerun = append(omsTestIdsToRerun, testDetail.CentralOmsTestId)
		}
	}

	return testDetailsIdsToBeRerun, omsTestIdsToRerun
}

func (taskService *TaskService) updateTaskDetailsAfterValidation(ctx context.Context,
	task commonModels.Task,
	testDetails []commonModels.TestDetail,
	investigations []commonModels.InvestigationResult,
	investigationsData []commonModels.InvestigationData,
	coAuthorizePathologist commonModels.CoAuthorizedPathologists,
	createRemarks []commonModels.Remark,
	updateRemarks []commonModels.Remark,
	deleteRemarkIds []uint,
	rerunInvestigationResults []commonModels.RerunInvestigationResult,
	testDetailsIdsToBeRerun []uint,
	visitIdToAttuneOrderMap map[string]commonStructures.AttuneOrderResponse,
	userId uint) *commonStructures.CommonError {

	tx := taskService.TaskDao.GetDbTransactionObject()
	defer tx.Rollback()

	_, cErr := taskService.TaskDao.UpdateTaskWithTx(tx, task)
	if cErr != nil {
		return cErr
	}

	_, cErr = taskService.TestDetailService.UpdateTestDetailsWithTx(tx, testDetails)
	if cErr != nil {
		return cErr
	}

	_, cErr = taskService.InvestigationResultsService.UpdateInvestigationResultsWithTx(tx, investigations)
	if cErr != nil {
		return cErr
	}

	_, cErr = taskService.InvestigationResultsService.UpdateInvestigationsDataWithTx(tx, investigationsData)
	if cErr != nil {
		return cErr
	}

	if len(rerunInvestigationResults) > 0 {
		_, cErr = taskService.RerunService.CreateRerunInvestigationResultsWithTx(tx, rerunInvestigationResults)
		if cErr != nil {
			return cErr
		}
	}

	cErr = taskService.CoAuthorizePathologistService.DeleteCurrentCoAuthorizePathologistWithTx(tx, task.Id, userId)
	if cErr != nil {
		return cErr
	}

	if coAuthorizePathologist.CoAuthorizedTo != 0 {
		_, cErr = taskService.CoAuthorizePathologistService.CreateCoAuthrizePathologistWithTx(tx, coAuthorizePathologist)
		if cErr != nil {
			return cErr
		}
	}

	if len(createRemarks) > 0 {
		_, cErr = taskService.RemarkService.CreateRemarksWithTx(tx, createRemarks)
		if cErr != nil {
			return cErr
		}
	}

	if len(updateRemarks) > 0 {
		_, cErr = taskService.RemarkService.UpdateRemarksWithTx(tx, updateRemarks)
		if cErr != nil {
			return cErr
		}
	}

	if len(deleteRemarkIds) > 0 {
		cErr = taskService.RemarkService.DeleteRemarksWithTx(tx, deleteRemarkIds, userId)
		if cErr != nil {
			return cErr
		}
	}

	if len(testDetailsIdsToBeRerun) > 0 {
		for _, attuneResponse := range visitIdToAttuneOrderMap {
			cErr = taskService.AttuneClient.InsertTestDataToAttune(ctx, attuneResponse)
			if cErr != nil {
				return cErr
			}
		}
	}

	tx.Commit()

	return nil
}

func (taskService *TaskService) GetUpdatedRemarkDetails(investigationStruct structures.UpdateInvestigationStruct,
	remarkIdToRemarkMap map[uint]commonModels.Remark, investigationId uint,
	createRemarks, updateRemarks []commonModels.Remark, deleteRemarksIds []uint, userId uint) (
	[]commonModels.Remark, []commonModels.Remark, []uint) {

	if investigationStruct.Remarks.MedicalRemark.Description != "" {
		if investigationStruct.Remarks.MedicalRemark.Id != 0 {
			if investigationStruct.Remarks.MedicalRemark.ToDelete {
				deleteRemarksIds = append(deleteRemarksIds, investigationStruct.Remarks.MedicalRemark.Id)
			} else {
				updateRemark := remarkIdToRemarkMap[investigationStruct.Remarks.MedicalRemark.Id]
				updateRemark.Description = investigationStruct.Remarks.MedicalRemark.Description
				updateRemark.RemarkType = commonConstants.REMARK_TYPE_MEDICAL_REMARK
				updateRemark.RemarkBy = userId
				updateRemark.UpdatedBy = userId
				updateRemarks = append(updateRemarks, updateRemark)
			}
		} else {
			createRemark := commonModels.Remark{
				InvestigationResultId: investigationId,
				Description:           investigationStruct.Remarks.MedicalRemark.Description,
				RemarkType:            commonConstants.REMARK_TYPE_MEDICAL_REMARK,
			}
			createRemark.RemarkBy = userId
			createRemark.CreatedBy = userId
			createRemark.UpdatedBy = userId
			createRemarks = append(createRemarks, createRemark)
		}
	} else {
		if investigationStruct.Remarks.MedicalRemark.Id != 0 && investigationStruct.Remarks.MedicalRemark.ToDelete {
			deleteRemarksIds = append(deleteRemarksIds, investigationStruct.Remarks.MedicalRemark.Id)
		}
	}

	if investigationStruct.Status == commonConstants.INVESTIGATION_STATUS_WITHHELD &&
		investigationStruct.Remarks.WithheldReason.Description != "" {
		createRemark := commonModels.Remark{
			InvestigationResultId: investigationId,
			Description:           investigationStruct.Remarks.WithheldReason.Description,
			RemarkType:            commonConstants.REMARK_TYPE_WITHHELD_REASON,
		}
		createRemark.RemarkBy = userId
		createRemark.CreatedBy = userId
		createRemark.UpdatedBy = userId
		createRemarks = append(createRemarks, createRemark)
	}

	return createRemarks, updateRemarks, deleteRemarksIds
}

func (taskService *TaskService) getRerunInvestigationStruct(testDetail commonModels.TestDetail,
	investigation commonModels.InvestigationResult,
	investigationStruct structures.UpdateInvestigationStruct,
	currentTime *time.Time, userId uint) commonModels.RerunInvestigationResult {
	rerunInvestigationResult := commonModels.RerunInvestigationResult{
		TestDetailsId:            testDetail.Id,
		MasterInvestigationId:    investigation.MasterInvestigationId,
		InvestigationName:        investigation.InvestigationName,
		InvestigationValue:       investigation.InvestigationValue,
		DeviceValue:              investigation.DeviceValue,
		ResultRepresentationType: investigation.ResultRepresentationType,
		LisCode:                  investigation.LisCode,
		RerunTriggeredBy:         userId,
		RerunTriggeredAt:         currentTime,
		RerunReason:              investigationStruct.RerunDetails.RerunReason,
		RerunRemarks:             investigationStruct.RerunDetails.RerunReason,
		EnteredBy:                investigation.EnteredBy,
		EnteredAt:                investigation.EnteredAt,
	}
	if rerunInvestigationResult.RerunReason == "" {
		rerunInvestigationResult.RerunReason = commonConstants.DEFAULT_RERUN_REASON
		rerunInvestigationResult.RerunRemarks = rerunInvestigationResult.RerunReason
	}

	rerunInvestigationResult.CreatedBy = userId
	rerunInvestigationResult.UpdatedBy = userId

	return rerunInvestigationResult
}

func getCoAuthorizePathologistStrict(coAuthorizePathologist commonModels.CoAuthorizedPathologists,
	taskId, userId, coAuthorizeToId uint, currentTime *time.Time) commonModels.CoAuthorizedPathologists {
	coAuthorizePathologist.TaskId = taskId
	coAuthorizePathologist.CoAuthorizedBy = userId
	coAuthorizePathologist.CoAuthorizedTo = coAuthorizeToId
	coAuthorizePathologist.CoAuthorizedAt = currentTime
	coAuthorizePathologist.CreatedBy = userId
	coAuthorizePathologist.UpdatedBy = userId

	return coAuthorizePathologist
}

func (taskService *TaskService) updateNewTaskDetails(testDetails []commonModels.TestDetail,
	investigations []commonModels.InvestigationResult,
	investigationsData []commonModels.InvestigationData,
	remarks []commonModels.Remark,
	taskDetails structures.UpdateTaskStruct,
	userId uint) (
	[]commonModels.TestDetail,
	[]commonModels.InvestigationResult,
	[]commonModels.InvestigationData,
	[]commonModels.Remark,
	[]commonModels.Remark,
	[]uint,
	[]commonModels.RerunInvestigationResult,
	commonModels.CoAuthorizedPathologists,
	*commonStructures.CommonError,
) {
	testDetailIdToTestDetailMap := make(map[uint]commonModels.TestDetail)
	testDetailIdToInvestigationMap := make(map[uint][]commonModels.InvestigationResult)
	investigationIdToInvestigationMap := make(map[uint]commonModels.InvestigationResult)
	investigationIdToInvestigationDataMap := make(map[uint]commonModels.InvestigationData)
	remarkIdToRemarkMap := make(map[uint]commonModels.Remark)

	rerunInvestigationResults := make([]commonModels.RerunInvestigationResult, 0)
	updateTestDetails, updateInvestigations, updateInvestigationsData :=
		[]commonModels.TestDetail{}, []commonModels.InvestigationResult{}, []commonModels.InvestigationData{}
	createRemarks, updateRemarks := []commonModels.Remark{}, []commonModels.Remark{}
	deleteRemarksIds, coAuthorizeToId := []uint{}, uint(0)
	coAuthorizePathologist := commonModels.CoAuthorizedPathologists{}

	for _, testDetail := range testDetails {
		testDetailIdToTestDetailMap[testDetail.Id] = testDetail
	}

	for _, investigation := range investigations {
		investigationIdToInvestigationMap[investigation.Id] = investigation
		if _, ok := testDetailIdToInvestigationMap[investigation.TestDetailsId]; !ok {
			testDetailIdToInvestigationMap[investigation.TestDetailsId] = []commonModels.InvestigationResult{}
		}
		testDetailIdToInvestigationMap[investigation.TestDetailsId] =
			append(testDetailIdToInvestigationMap[investigation.TestDetailsId], investigation)
	}

	for _, investigationData := range investigationsData {
		investigationIdToInvestigationDataMap[investigationData.InvestigationResultId] = investigationData
	}

	for _, remark := range remarks {
		remarkIdToRemarkMap[remark.Id] = remark
	}

	currentTime := commonUtils.GetCurrentTime()

	for _, testDetailStruct := range taskDetails.TestDetails {
		testDetail := testDetailIdToTestDetailMap[testDetailStruct.Id]
		testDetailIdToTestDetailMap[testDetailStruct.Id] = testDetail
		updateTestDetails = append(updateTestDetails, testDetail)

		// Update Investigations
		for _, investigationStruct := range testDetailStruct.Investigations {
			investigation := investigationIdToInvestigationMap[investigationStruct.Id]
			investigation.InvestigationStatus = investigationStruct.Status
			investigation.Abnormality = investigationStruct.Abnormality
			investigation.IsAbnormal = getIsAbnormalFlagForInvestigation(investigationStruct.Abnormality)
			if investigationStruct.Status == commonConstants.INVESTIGATION_STATUS_APPROVE {
				investigation.ApprovedBy = userId
				investigation.ApprovedAt = currentTime
			}

			// Update Investigation Data
			if investigationStruct.InvestigationData != "" {
				investigationData := investigationIdToInvestigationDataMap[investigation.Id]
				investigationData.Data = investigationStruct.InvestigationData
				investigationData.UpdatedBy = userId
				updateInvestigationsData = append(updateInvestigationsData, investigationData)
			}

			// Update Rerun Details
			if investigationStruct.Status == commonConstants.INVESTIGATION_STATUS_RERUN {
				rerunInvestigationResults = append(rerunInvestigationResults,
					taskService.getRerunInvestigationStruct(testDetail, investigation, investigationStruct,
						currentTime, userId))
			}

			if investigationStruct.Status == commonConstants.INVESTIGATION_STATUS_CO_AUTHORIZE {
				if coAuthorizeToId == 0 {
					if investigationStruct.CoAuthorizeTo == 0 {
						return nil, nil, nil, nil, nil, nil, nil, coAuthorizePathologist, &commonStructures.CommonError{
							Message:    commonConstants.ERROR_CO_AUTHORIZE_TO_MISSING,
							StatusCode: http.StatusBadRequest,
						}
					}
					if investigationStruct.CoAuthorizeTo == userId {
						return nil, nil, nil, nil, nil, nil, nil, coAuthorizePathologist, &commonStructures.CommonError{
							Message:    commonConstants.ERROR_CO_AUTHORIZE_TO_SELF,
							StatusCode: http.StatusBadRequest,
						}
					}
					coAuthorizeToId = investigationStruct.CoAuthorizeTo
					coAuthorizePathologist = getCoAuthorizePathologistStrict(coAuthorizePathologist,
						taskDetails.Id, userId, coAuthorizeToId, currentTime)
				} else if coAuthorizeToId != investigationStruct.CoAuthorizeTo {
					return nil, nil, nil, nil, nil, nil, nil, coAuthorizePathologist, &commonStructures.CommonError{
						Message:    commonConstants.ERROR_CO_AUTHORIZE_TO_MISMATCH,
						StatusCode: http.StatusBadRequest,
					}
				}
			}

			// Update Remarks
			createRemarks, updateRemarks, deleteRemarksIds = taskService.GetUpdatedRemarkDetails(investigationStruct,
				remarkIdToRemarkMap, investigation.Id, createRemarks, updateRemarks, deleteRemarksIds, userId)

			investigation.InvestigationValue = investigationStruct.InvestigationValue
			investigation.UpdatedBy = userId
			investigationIdToInvestigationMap[investigationStruct.Id] = investigation
			updateInvestigations = append(updateInvestigations, investigation)
		}
	}

	deleteRemarksIds = commonUtils.CreateUniqueSliceUint(deleteRemarksIds)
	return updateTestDetails, updateInvestigations, updateInvestigationsData, createRemarks,
		updateRemarks, deleteRemarksIds, rerunInvestigationResults, coAuthorizePathologist, nil
}

func (taskService *TaskService) validatePathologist(taskId, userId uint) bool {
	taskPathMap, cErr := taskService.TaskPathService.GetActiveTaskPathMapByTaskId(taskId)
	if cErr != nil {
		return false
	}

	if taskPathMap.PathologistID != userId || !taskPathMap.IsActive {
		return false
	}

	pathologist, cErr := taskService.UserService.GetUserModel(userId)
	if cErr != nil {
		return false
	}

	userTypes := []string{
		commonConstants.USER_TYPE_PATHOLOGIST,
		commonConstants.USER_TYPE_SUPER_ADMIN,
	}

	return commonUtils.SliceContainsString(userTypes, pathologist.UserType) && pathologist.AttuneUserId != "" &&
		pathologist.SystemUserId != ""
}

func (taskService *TaskService) UpdateAllTaskDetails(ctx context.Context,
	allTaskDetails structures.UpdateAllTaskDetailsStruct) *commonStructures.CommonError {

	taskDetails := allTaskDetails.Task
	userId := allTaskDetails.UserId
	taskId := taskDetails.Id
	testDetailIds, investigationIds, remarkIds := []uint{}, []uint{}, []uint{}
	testDetailIdToTestDetailMap := make(map[uint]commonModels.TestDetail)
	remarkIdToRemarkMap := make(map[uint]commonModels.Remark)

	task, user := commonModels.Task{}, commonModels.User{}
	newTestDetails, oldTestDetails := []commonModels.TestDetail{}, []commonModels.TestDetail{}
	newInvestigations, oldInvestigations := []commonModels.InvestigationResult{}, []commonModels.InvestigationResult{}
	newInvestigationsData, newRemarks := []commonModels.InvestigationData{}, []commonModels.Remark{}
	var cErr *commonStructures.CommonError
	var errList []*commonStructures.CommonError
	var mu sync.Mutex
	visitIdToAttuneOrderMap := make(map[string]commonStructures.AttuneOrderResponse)

	isPathologistValid := taskService.validatePathologist(taskId, userId)
	if !isPathologistValid {
		return &commonStructures.CommonError{
			Message:    commonConstants.ERROR_ACTION_NOT_ALLOWED_FOR_PATHOLOGIST,
			StatusCode: http.StatusConflict,
		}
	}

	if !validateApprovedInvestigationsByValues(allTaskDetails.Task.TestDetails) {
		return &commonStructures.CommonError{
			Message:    commonConstants.ERROR_NEGATIVE_INVESTIGATION_VALUE,
			StatusCode: http.StatusBadRequest,
		}
	}

	for _, testDetail := range taskDetails.TestDetails {
		if testDetail.Id == 0 {
			continue
		}
		testDetailIds = append(testDetailIds, testDetail.Id)
		for _, investigation := range testDetail.Investigations {
			if investigation.Id == 0 {
				continue
			}
			investigationIds = append(investigationIds, investigation.Id)
			if investigation.Remarks.MedicalRemark.Id != 0 {
				remarkIds = append(remarkIds, investigation.Remarks.MedicalRemark.Id)
			}
		}
	}

	wg := sync.WaitGroup{}
	wg.Add(2)
	go func() {
		defer wg.Done()

		task, newTestDetails, newInvestigations, newInvestigationsData, newRemarks, cErr =
			taskService.fetchAllUpdateTaskData(taskId, testDetailIds, investigationIds, remarkIds)
		if cErr != nil {
			mu.Lock()
			errList = append(errList, cErr)
			mu.Unlock()
		}
	}()

	go func() {
		defer wg.Done()

		oldTestDetails, oldInvestigations, user, cErr = taskService.fetchAllTaskData(taskId, userId)
		if cErr != nil {
			mu.Lock()
			errList = append(errList, cErr)
			mu.Unlock()
		}
	}()

	wg.Wait()

	if len(errList) > 0 {
		return errList[0]
	}

	newTestDetails, newInvestigations, newInvestigationsData, createRemarks,
		updateRemarks, deleteRemarkIds, rerunDetails, coAuthorizePathologist, cErr :=
		taskService.updateNewTaskDetails(newTestDetails, newInvestigations,
			newInvestigationsData, newRemarks, taskDetails, userId)
	if cErr != nil {
		return cErr
	}

	cErr = taskService.validateUpdateTaskDetailsAndUpdateData(&task, oldTestDetails,
		oldInvestigations, newTestDetails, newInvestigations)
	if cErr != nil {
		return cErr
	}

	testDetailsIdsToBeRerun, omsTestIdsToRerun := taskService.getTestDetailsToBeRerun(newTestDetails)

	for _, testDetail := range newTestDetails {
		testDetailIdToTestDetailMap[testDetail.Id] = testDetail
	}

	if len(testDetailsIdsToBeRerun) != 0 {
		visitIdToAttuneOrderMap, cErr = taskService.FetchAttuneDataForReruningTests(ctx, task.CityCode,
			testDetailsIdsToBeRerun, testDetailIdToTestDetailMap, newInvestigations, rerunDetails, user)
		if cErr != nil {
			return cErr
		}
	}

	for _, remark := range newRemarks {
		remarkIdToRemarkMap[remark.Id] = remark
	}

	cErr = taskService.updateTaskDetailsAfterValidation(ctx, task, newTestDetails, newInvestigations,
		newInvestigationsData, coAuthorizePathologist,
		createRemarks, updateRemarks, deleteRemarkIds,
		rerunDetails, testDetailsIdsToBeRerun, visitIdToAttuneOrderMap, userId)
	if cErr != nil {
		return cErr
	}
	if len(omsTestIdsToRerun) != 0 {
		go taskService.EtsService.GetAndPublishEtsTestEventForLisWebhook(context.Background(), omsTestIdsToRerun,
			commonConstants.AttuneTestStatusRerun)
	}

	return nil
}

func (taskService *TaskService) UndoReportRelease(ctx context.Context, taskID uint) *commonStructures.CommonError {
	reportReleaseKey := fmt.Sprintf(commonConstants.ReportReleaseKey, taskID)
	err := taskService.Cache.Delete(ctx, reportReleaseKey)
	if err != nil {
		return &commonStructures.CommonError{
			Message:    commonConstants.ERROR_IN_UNDOING_REPORT_RELEASE,
			StatusCode: http.StatusInternalServerError,
		}
	}

	return nil
}

func (taskService *TaskService) GetTaskMetadataByTaskId(taskID uint) (
	commonModels.TaskMetadata, *commonStructures.CommonError) {
	return taskService.TaskDao.GetTaskMetadataByTaskId(taskID)
}

func (taskService *TaskService) UpdateTaskMetadata(taskMetadata commonModels.TaskMetadata) (
	commonModels.TaskMetadata, *commonStructures.CommonError) {
	return taskService.TaskDao.UpdateTaskMetadata(taskMetadata)
}

func (taskService *TaskService) CreateTaskMetadataWithTx(tx *gorm.DB, taskMetadata commonModels.TaskMetadata) (
	commonModels.TaskMetadata, *commonStructures.CommonError) {
	return taskService.TaskDao.CreateTaskMetadataWithTx(tx, taskMetadata)
}

func (taskService *TaskService) UpdateTaskMetadataWithTx(tx *gorm.DB, taskMetadata commonModels.TaskMetadata) (
	commonModels.TaskMetadata, *commonStructures.CommonError) {
	return taskService.TaskDao.UpdateTaskMetadataWithTx(tx, taskMetadata)
}

func (taskService *TaskService) GetTaskVisitMappingsByTaskId(taskID uint) (
	[]commonModels.TaskVisitMapping, *commonStructures.CommonError) {
	return taskService.TaskDao.GetTaskVisitMappingsByTaskId(taskID)
}

func (taskService *TaskService) CreateTaskVisitMappingWithTx(tx *gorm.DB,
	taskVisitMappings []commonModels.TaskVisitMapping) ([]commonModels.TaskVisitMapping, *commonStructures.CommonError) {
	if len(taskVisitMappings) == 0 {
		return taskVisitMappings, nil
	}
	return taskService.TaskDao.CreateTaskVisitMappingWithTransaction(tx, taskVisitMappings)
}

func (taskService *TaskService) CreateDeleteTaskVisitMappingsWithTx(tx *gorm.DB, taskId uint,
	toCreateVisitIds, toDeleteVisitIds []string) *commonStructures.CommonError {
	if len(toCreateVisitIds) > 0 {
		taskVisitMappings := make([]commonModels.TaskVisitMapping, 0, len(toCreateVisitIds))
		for _, visitId := range toCreateVisitIds {
			taskVisitMappings = append(taskVisitMappings, commonModels.TaskVisitMapping{
				TaskId:  taskId,
				VisitId: visitId,
			})
		}
		_, cErr := taskService.TaskDao.CreateTaskVisitMappingWithTransaction(tx, taskVisitMappings)
		if cErr != nil {
			return cErr
		}
	}

	if len(toDeleteVisitIds) > 0 {
		cErr := taskService.TaskDao.DeleteTaskVisitMappingsByTaskIdAndVisitIds(tx, taskId, toDeleteVisitIds)
		if cErr != nil {
			return cErr
		}
	}

	return nil
}

func (taskService *TaskService) GetTaskCallingDetails(taskId, userId uint, callingType string) (
	structures.TaskCallingDetailsResponse, *commonStructures.CommonError) {

	if !commonUtils.SliceContainsString(commonConstants.CALLING_TYPES, callingType) {
		return structures.TaskCallingDetailsResponse{}, &commonStructures.CommonError{
			Message:    commonConstants.ERROR_INVALID_CALLING_TYPE,
			StatusCode: http.StatusBadRequest,
		}
	}

	mobileNumber := ""
	user := commonModels.User{}
	var errList []*commonStructures.CommonError

	wg, mu := sync.WaitGroup{}, sync.Mutex{}
	wg.Add(2)

	go func() {
		defer wg.Done()
		switch callingType {
		case commonConstants.CALLING_TYPE_CUSTOMER:
			task, cErr := taskService.TaskDao.GetTaskWithPatientDetailsById(taskId)
			if cErr != nil {
				mu.Lock()
				errList = append(errList, cErr)
				mu.Unlock()
				return
			}
			mobileNumber = task.PatientDetail.Number
		case commonConstants.CALLING_TYPE_DOCTOR:
			taskMetadata, cErr := taskService.TaskDao.GetTaskMetadataByTaskId(taskId)
			if cErr != nil {
				mu.Lock()
				errList = append(errList, cErr)
				mu.Unlock()
				return
			}
			mobileNumber = taskMetadata.DoctorNumber
		}
	}()

	go func() {
		defer wg.Done()
		var cErr *commonStructures.CommonError
		user, cErr = taskService.UserService.GetUserModel(userId)
		if cErr != nil {
			mu.Lock()
			errList = append(errList, cErr)
			mu.Unlock()
		}
	}()

	wg.Wait()

	if len(errList) > 0 {
		return structures.TaskCallingDetailsResponse{}, errList[0]
	}

	return structures.TaskCallingDetailsResponse{
		MobileNumber: mobileNumber,
		AgentId:      user.AgentId,
	}, nil
}

func (taskService *TaskService) GetTaskIdByVisitId(visitId string) (uint, *commonStructures.CommonError) {
	taskVisitMapping, cErr := taskService.TaskDao.GetTaskVisitMappingsByVisitId(visitId)
	if cErr != nil {
		return 0, cErr
	}

	return taskVisitMapping.TaskId, nil
}

func (taskService *TaskService) GetQcFailedTestDataToRerun(ctx context.Context, userId uint, testDetailsIdsToRerun []uint,
	cityCode string, testDetails []commonModels.TestDetail, qcFailedInvestigations []commonModels.InvestigationResult,
	userIdToAttuneUserId map[uint]int) (
	[]commonModels.RerunInvestigationResult, map[string]commonStructures.AttuneOrderResponse, *commonStructures.CommonError) {
	testDetailIdToTestDetailMap := map[uint]commonModels.TestDetail{}
	for _, testDetail := range testDetails {
		testDetailIdToTestDetailMap[testDetail.Id] = testDetail
	}
	user, cErr := taskService.UserService.GetUserModel(userId)
	if cErr != nil {
		return nil, nil, cErr
	}

	if len(testDetailsIdsToRerun) > 0 {
		currentTime := commonUtils.GetCurrentTime()
		rerunDetails := []commonModels.RerunInvestigationResult{}
		for _, investigation := range qcFailedInvestigations {
			investigationStruct := getUpdateInvestigationStruct(investigation)
			rerunDetail := taskService.getRerunInvestigationStruct(testDetailIdToTestDetailMap[investigation.TestDetailsId],
				investigation, investigationStruct, currentTime, userId)
			rerunDetails = append(rerunDetails, rerunDetail)
		}

		newVisitIdToAttuneOrderMap, cErr := taskService.FetchAttuneDataForReruningTests(ctx, cityCode,
			testDetailsIdsToRerun, testDetailIdToTestDetailMap, qcFailedInvestigations, rerunDetails, user)

		// Adding userId from auto approval of the city code here.
		userIdForRerun, _ := (strconv.ParseUint(constants.AutoApprovalIdsMap[strings.ToLower(cityCode)], 10, 64))
		userIdForRerunUint := uint(userIdForRerun)
		attuneUserId := userIdToAttuneUserId[userIdForRerunUint]
		for visitId, attuneOrder := range newVisitIdToAttuneOrderMap {
			for idx, order := range attuneOrder.OrderInfo {
				order.UserID = attuneUserId
				attuneOrder.OrderInfo[idx] = order
			}
			newVisitIdToAttuneOrderMap[visitId] = attuneOrder
		}

		if cErr != nil {
			return []commonModels.RerunInvestigationResult{}, map[string]commonStructures.AttuneOrderResponse{}, cErr
		}
		return rerunDetails, newVisitIdToAttuneOrderMap, nil
	}

	return []commonModels.RerunInvestigationResult{}, map[string]commonStructures.AttuneOrderResponse{}, nil
}

func getUpdateInvestigationStruct(investigation commonModels.InvestigationResult) structures.UpdateInvestigationStruct {
	return structures.UpdateInvestigationStruct{
		Id:     investigation.Id,
		Status: commonConstants.INVESTIGATION_STATUS_RERUN,
		RerunDetails: structures.UpdateRerunDetailsStruct{
			RerunReason: commonConstants.RERUN_REASON_QC_FAIL,
		},
	}
}

func validateApprovedInvestigationsByValues(testDetails []structures.UpdateTestDetailsStruct) bool {
	for _, testDetail := range testDetails {
		for _, investigation := range testDetail.Investigations {
			if investigation.Status == commonConstants.INVESTIGATION_STATUS_APPROVE {
				floatVal := commonUtils.ConvertStringToFloat32ForAbnormality(investigation.InvestigationValue)
				if floatVal < 0 {
					return false
				}
				if investigation.InvestigationValue != "" && investigation.InvestigationValue[0] == '-' {
					return false
				}
			}

		}
	}
	return true
}

func getIsAbnormalFlagForInvestigation(ohAbnormality string) bool {
	return commonUtils.SliceContainsString(constants.OhAbnormalityStringSlice, ohAbnormality)
}
