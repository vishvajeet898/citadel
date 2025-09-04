package dao

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/Orange-Health/citadel/common/constants"
	commonConstants "github.com/Orange-Health/citadel/common/constants"
	commonStructures "github.com/Orange-Health/citadel/common/structures"
	commonUtils "github.com/Orange-Health/citadel/common/utils"
	"github.com/Orange-Health/citadel/models"
)

type DataLayer interface {
	GetTaskPathMapByTaskId(taskID uint) (
		models.TaskPathologistMapping, *commonStructures.CommonError)
	GetActiveTaskPathMapByTaskId(taskID uint) (
		models.TaskPathologistMapping, *commonStructures.CommonError)
	CreateTaskPathMap(tpm models.TaskPathologistMapping) (
		models.TaskPathologistMapping, *commonStructures.CommonError)
	UpdateTaskPathMapByTaskId(taskID uint, tpm models.TaskPathologistMapping) (
		models.TaskPathologistMapping, *commonStructures.CommonError)
	MarkAllTaskPathMapInactiveByTaskIdWithTx(tx *gorm.DB,
		taskID uint) *commonStructures.CommonError
}

func (taskPathMapDao *TaskPathMapDao) GetTaskPathMapByTaskId(taskID uint) (
	models.TaskPathologistMapping, *commonStructures.CommonError) {

	var tpm models.TaskPathologistMapping
	if err := taskPathMapDao.Db.Where("task_id = ?", taskID).First(&tpm).Error; err != nil {
		return tpm, commonUtils.HandleORMError(err)
	}
	return tpm, nil
}

func (taskPathMapDao *TaskPathMapDao) GetActiveTaskPathMapByTaskId(taskID uint) (
	models.TaskPathologistMapping, *commonStructures.CommonError) {

	var tpm models.TaskPathologistMapping
	if err := taskPathMapDao.Db.Where("task_id = ? AND is_active = true", taskID).First(&tpm).Error; err != nil {
		return tpm, commonUtils.HandleORMError(err)
	}
	return tpm, nil
}

func (taskPathMapDao *TaskPathMapDao) CreateTaskPathMap(tpm models.TaskPathologistMapping) (
	models.TaskPathologistMapping, *commonStructures.CommonError) {
	err := taskPathMapDao.Db.Transaction(func(tx *gorm.DB) error {
		task := models.Task{}
		err := tx.Clauses(clause.Locking{
			Strength: commonConstants.CLAUSE_UPDATE,
			Options:  commonConstants.CLAUSE_NOWAIT,
		}).First(&task, tpm.TaskId).Error

		if err != nil {
			return err
		}

		if err := tx.Model(&tpm).Create(&tpm).Error; err != nil {
			return err
		}

		tasksMap := map[string]interface{}{
			"status":          constants.TASK_STATUS_IN_PROGRESS,
			"updated_by":      tpm.UpdatedBy,
			"updated_at":      commonUtils.GetCurrentTime(),
			"previous_status": task.Status,
		}
		if err := tx.Model(&task).Where("id = ?", tpm.TaskId).Updates(tasksMap).Error; err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return tpm, commonUtils.HandleORMError(err)
	}

	return tpm, nil
}

func (taskPathMapDao *TaskPathMapDao) UpdateTaskPathMapByTaskId(taskID uint, tpm models.TaskPathologistMapping) (
	models.TaskPathologistMapping, *commonStructures.CommonError) {

	err := taskPathMapDao.Db.Transaction(func(tx *gorm.DB) error {

		task := models.Task{}
		err := tx.Clauses(clause.Locking{
			Strength: commonConstants.CLAUSE_UPDATE,
			Options:  commonConstants.CLAUSE_NOWAIT,
		}).First(&task, tpm.TaskId).Error

		if err != nil {
			return err
		}

		tpmMap := map[string]interface{}{
			"pathologist_id": tpm.PathologistId,
			"task_id":        tpm.TaskId,
			"is_active":      tpm.IsActive,
			"updated_at":     commonUtils.GetCurrentTime(),
			"updated_by":     tpm.PathologistId,
		}

		if tpm.DeletedAt != nil {
			tpmMap["deleted_at"] = tpm.DeletedAt
			tpmMap["deleted_by"] = tpm.PathologistId
		}

		if err := tx.Model(&tpm).Where("task_id = ?", taskID).Updates(tpmMap).Error; err != nil {
			tx.Rollback()
			return err
		}

		if !tpm.IsActive && task.Status == constants.TASK_STATUS_IN_PROGRESS {
			tasksMap := map[string]interface{}{
				"status":     task.PreviousStatus,
				"updated_by": tpm.PathologistId,
				"updated_at": commonUtils.GetCurrentTime(),
			}
			if err := tx.Model(&task).Where("id = ?", tpm.TaskId).Updates(tasksMap).Error; err != nil {
				return err
			}
		} else if tpm.IsActive {
			tasksMap := map[string]interface{}{
				"status":          constants.TASK_STATUS_IN_PROGRESS,
				"previous_status": task.Status,
				"updated_by":      tpm.UpdatedBy,
			}
			if err := tx.Model(&task).Where("id = ?", tpm.TaskId).Updates(tasksMap).Error; err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		return tpm, commonUtils.HandleORMError(err)
	}

	return tpm, nil
}

func (taskPathMapDao *TaskPathMapDao) MarkAllTaskPathMapInactiveByTaskIdWithTx(tx *gorm.DB,
	taskID uint) *commonStructures.CommonError {

	updates := map[string]interface{}{
		"is_active":  false,
		"updated_at": commonUtils.GetCurrentTime(),
		"updated_by": commonConstants.CitadelSystemId,
	}
	if err := tx.Model(&models.TaskPathologistMapping{}).Where("task_id = ?", taskID).
		Updates(updates).Error; err != nil {
		return commonUtils.HandleORMError(err)
	}

	return nil
}
