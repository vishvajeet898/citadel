package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	mapper "github.com/Orange-Health/citadel/apps/task_pathologist_mapping/mapper"
	"github.com/Orange-Health/citadel/apps/task_pathologist_mapping/structures"
	commonConstants "github.com/Orange-Health/citadel/common/constants"
	commonStructures "github.com/Orange-Health/citadel/common/structures"
	"github.com/Orange-Health/citadel/models"
	"gorm.io/gorm"
)

type TaskPathologistMappingServiceInterface interface {
	GetActiveTaskPathMapByTaskId(taskID uint) (
		structures.TaskPathologistMapping, *commonStructures.CommonError)
	CreateTaskPathMap(ctx context.Context, tpm structures.TaskPathologistMapping) (
		structures.TaskPathologistMapping, *commonStructures.CommonError)
	MarkAllTaskPathMapInactiveByTaskIdWithTx(tx *gorm.DB,
		taskID uint) *commonStructures.CommonError
}

func (tpmService *TaskPathologistMappingService) GetActiveTaskPathMapByTaskId(taskID uint) (
	structures.TaskPathologistMapping, *commonStructures.CommonError) {
	tpm, err := tpmService.TaskPathDao.GetActiveTaskPathMapByTaskId(taskID)
	if err != nil {
		return structures.TaskPathologistMapping{}, err
	}

	return mapper.MapTpm(tpm), nil
}

func (tpmService *TaskPathologistMappingService) CreateTaskPathMap(ctx context.Context,
	tpm structures.TaskPathologistMapping) (
	structures.TaskPathologistMapping, *commonStructures.CommonError) {

	redisKey := fmt.Sprintf(commonConstants.CacheKeyTaskPathlogistMapping, tpm.TaskID)
	if tpm.IsActive {
		// Check if the redis key exists for this task id and run the code accordingly
		var tpmRedisInterface interface{}
		cacheErr := tpmService.Cache.Get(ctx, redisKey, &tpmRedisInterface)
		if tpmRedisInterface != nil && cacheErr == nil {
			tpmRedisBytes, err := json.Marshal(tpmRedisInterface)
			if err == nil {
				tpmRedis := structures.TaskPathologistMapping{}
				err = json.Unmarshal(tpmRedisBytes, &tpmRedis)
				if err == nil {
					if tpmRedis.PathologistID == tpm.PathologistID {
						return tpmRedis, nil
					} else {
						return structures.TaskPathologistMapping{}, &commonStructures.CommonError{
							StatusCode: http.StatusBadRequest,
							Message:    fmt.Sprintf(commonConstants.ERROR_TPM_IN_SAME_STATE, fmt.Sprint(tpm.IsActive)),
						}
					}
				}
			}
		}

		// Create the redis key and db mapping
		err := tpmService.Cache.Set(ctx, redisKey, tpm, commonConstants.TaskPathologistMappingKeyExpiry)
		if err != nil {
			return structures.TaskPathologistMapping{}, &commonStructures.CommonError{
				StatusCode: http.StatusInternalServerError,
				Message:    commonConstants.ERROR_FAILED_TO_SET_TASK_PATH_MAPPING,
			}
		}
		createdTpm, cErr := tpmService.createTaskPathMapAfterRedisUpdation(ctx, tpm)
		if cErr != nil {
			_ = tpmService.Cache.Delete(ctx, redisKey)
			return structures.TaskPathologistMapping{}, cErr
		}
		return createdTpm, nil
	}

	_ = tpmService.Cache.Delete(ctx, redisKey)
	return tpmService.createTaskPathMapAfterRedisUpdation(ctx, tpm)
}

func (tpmService *TaskPathologistMappingService) createTaskPathMapAfterRedisUpdation(ctx context.Context,
	tpm structures.TaskPathologistMapping) (
	structures.TaskPathologistMapping, *commonStructures.CommonError) {

	tpmModel := models.TaskPathologistMapping{
		TaskId:        tpm.TaskID,
		PathologistId: tpm.PathologistID,
		IsActive:      tpm.IsActive,
	}
	curTpm, _ := tpmService.TaskPathDao.GetTaskPathMapByTaskId(tpmModel.TaskId)

	isUpdated, updatedTpm, cErr := tpmService.updateIfTpmExists(&tpmModel, &curTpm)
	if cErr != nil {
		return structures.TaskPathologistMapping{}, cErr
	}

	if isUpdated {
		return mapper.MapTpm(updatedTpm), nil
	}

	if curTpm.IsActive && curTpm.TaskId == tpmModel.TaskId {
		return structures.TaskPathologistMapping{}, &commonStructures.CommonError{
			StatusCode: http.StatusBadRequest,
			Message:    commonConstants.ERROR_TPM_EXISTS,
		}
	}

	tpmModel.CreatedBy = tpm.PathologistID
	tpmModel.UpdatedBy = tpm.PathologistID

	createdTpm, cErr := tpmService.TaskPathDao.CreateTaskPathMap(tpmModel)
	if cErr != nil {
		if cErr.StatusCode == http.StatusInternalServerError {
			tpmService.Sentry.LogError(ctx, cErr.Message, errors.New(cErr.Message), map[string]interface{}{
				"task_id":        tpmModel.TaskId,
				"pathologist_id": tpmModel.PathologistId,
				"is_active":      tpmModel.IsActive,
			})
		}
		return structures.TaskPathologistMapping{}, cErr
	}

	return mapper.MapTpm(createdTpm), nil
}

func (tpmService *TaskPathologistMappingService) updateIfTpmExists(
	tpmModel, curTpm *models.TaskPathologistMapping) (bool, models.TaskPathologistMapping, *commonStructures.CommonError) {
	if curTpm.Id == 0 {
		return false, models.TaskPathologistMapping{}, nil
	}

	if curTpm.IsActive == tpmModel.IsActive {
		if curTpm.PathologistId == tpmModel.PathologistId {
			return true, *curTpm, nil
		}
		return false, models.TaskPathologistMapping{}, &commonStructures.CommonError{
			StatusCode: http.StatusBadRequest,
			Message:    fmt.Sprintf(commonConstants.ERROR_TPM_IN_SAME_STATE, fmt.Sprint(curTpm.IsActive)),
		}
	}

	tpmModel.Id = curTpm.Id
	updatedTpm, cErr := tpmService.TaskPathDao.UpdateTaskPathMapByTaskId(tpmModel.TaskId, *tpmModel)
	if cErr != nil {
		return false, models.TaskPathologistMapping{}, cErr
	}

	return true, updatedTpm, nil
}

func (tpmService *TaskPathologistMappingService) MarkAllTaskPathMapInactiveByTaskIdWithTx(tx *gorm.DB,
	taskID uint) *commonStructures.CommonError {

	return tpmService.TaskPathDao.MarkAllTaskPathMapInactiveByTaskIdWithTx(tx, taskID)
}
