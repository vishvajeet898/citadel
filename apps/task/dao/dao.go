package dao

import (
	"gorm.io/gorm"

	commonConstants "github.com/Orange-Health/citadel/common/constants"
	commonStructures "github.com/Orange-Health/citadel/common/structures"
	commonUtils "github.com/Orange-Health/citadel/common/utils"
	commonModels "github.com/Orange-Health/citadel/models"
)

type DataLayer interface {
	// Transactions
	GetDbTransactionObject() *gorm.DB

	// Tasks
	GetTaskById(taskId uint) (commonModels.Task, *commonStructures.CommonError)
	GetTaskByOmsOrderId(omsOrderId string) (commonModels.Task, *commonStructures.CommonError)
	GetTasksCount() (uint, *commonStructures.CommonError)
	GetAmendmentTasksCount() (uint, *commonStructures.CommonError)
	GetTaskWithPatientDetailsById(taskId uint) (commonModels.Task, *commonStructures.CommonError)
	UpdateTask(task commonModels.Task) (commonModels.Task, *commonStructures.CommonError)
	CreateTaskWithTx(tx *gorm.DB, task commonModels.Task) (commonModels.Task, *commonStructures.CommonError)
	UpdateTaskWithTx(tx *gorm.DB, task commonModels.Task) (commonModels.Task, *commonStructures.CommonError)
	DeleteTaskWithTx(tx *gorm.DB, taskId uint) *commonStructures.CommonError

	// Task Metadata
	GetTaskMetadataByTaskId(taskID uint) (commonModels.TaskMetadata, *commonStructures.CommonError)
	UpdateTaskMetadata(taskMetadata commonModels.TaskMetadata) (commonModels.TaskMetadata, *commonStructures.CommonError)
	CreateTaskMetadataWithTx(tx *gorm.DB, taskMetadata commonModels.TaskMetadata) (commonModels.TaskMetadata, *commonStructures.CommonError)
	UpdateTaskMetadataWithTx(tx *gorm.DB, taskMetadata commonModels.TaskMetadata) (commonModels.TaskMetadata, *commonStructures.CommonError)

	// Task Visit Mapping
	GetTaskVisitMappingsByTaskId(taskID uint) ([]commonModels.TaskVisitMapping, *commonStructures.CommonError)
	GetTaskVisitMappingsByVisitId(visitId string) (commonModels.TaskVisitMapping, *commonStructures.CommonError)
	CreateTaskVisitMappingWithTransaction(tx *gorm.DB, taskVisitMappings []commonModels.TaskVisitMapping) ([]commonModels.TaskVisitMapping, *commonStructures.CommonError)
	DeleteTaskVisitMappingsByTaskIdAndVisitIds(tx *gorm.DB, taskId uint, toDeleteVisitIds []string) *commonStructures.CommonError
}

func (taskDao *TaskDao) GetDbTransactionObject() *gorm.DB {
	return taskDao.Db.Begin()
}

func (taskDao *TaskDao) GetTaskById(taskId uint) (commonModels.Task, *commonStructures.CommonError) {

	task := commonModels.Task{}
	if err := taskDao.Db.Find(&task, taskId).Error; err != nil {
		return task, commonUtils.HandleORMError(err)
	}
	return task, nil
}

func (taskDao *TaskDao) GetTaskByOmsOrderId(omsOrderId string) (commonModels.Task, *commonStructures.CommonError) {
	task := commonModels.Task{}
	if err := taskDao.Db.Find(&task, "oms_order_id = ?", omsOrderId).Error; err != nil {
		return task, commonUtils.HandleORMError(err)
	}
	return task, nil
}

func (taskDao *TaskDao) GetTasksCount() (uint, *commonStructures.CommonError) {
	var count int64
	taskStatuses := []string{
		commonConstants.TASK_STATUS_PENDING,
		commonConstants.TASK_STATUS_WITHHELD_APPROVAL,
		commonConstants.TASK_STATUS_IN_PROGRESS,
		commonConstants.TASK_STATUS_CO_AUTHORIZE,
	}

	testDetailStatuses := []string{
		commonConstants.TEST_STATUS_RESULT_SAVED,
		commonConstants.TEST_STATUS_RERUN_RESULT_SAVED,
		commonConstants.TEST_STATUS_WITHHELD,
		commonConstants.TEST_STATUS_CO_AUTHORIZE,
	}

	err := taskDao.Db.Model(&commonModels.Task{}).
		Joins("INNER JOIN test_details ON tasks.id = test_details.task_id").
		Where("tasks.status IN (?)", taskStatuses).
		Where("test_details.status IN (?)", testDetailStatuses).
		Where("test_details.cp_enabled = ?", true).
		Distinct("tasks.id").Count(&count).Error
	if err != nil {
		return 0, commonUtils.HandleORMError(err)
	}
	return uint(count), nil
}

func (taskDao *TaskDao) GetAmendmentTasksCount() (uint, *commonStructures.CommonError) {
	var count int64

	err := taskDao.Db.Model(&commonModels.Task{}).
		Where("tasks.status = ?", commonConstants.TASK_STATUS_COMPLETED).
		Count(&count).Error
	if err != nil {
		return 0, commonUtils.HandleORMError(err)
	}
	return uint(count), nil
}

func (taskDao *TaskDao) GetTaskWithPatientDetailsById(taskId uint) (commonModels.Task, *commonStructures.CommonError) {

	task := commonModels.Task{}
	if err := taskDao.Db.Preload(commonConstants.PatientDetails).
		Find(&task, taskId).Error; err != nil {
		return task, commonUtils.HandleORMError(err)
	}
	return task, nil
}

func (taskDao *TaskDao) GetTaskMetadataByTaskId(taskId uint) (commonModels.TaskMetadata, *commonStructures.CommonError) {

	taskMetadata := commonModels.TaskMetadata{}
	err := taskDao.Db.Find(&taskMetadata, "task_id = ?", taskId).Error
	if err != nil {
		return taskMetadata, commonUtils.HandleORMError(err)
	}
	return taskMetadata, nil
}

func (taskDao *TaskDao) UpdateTask(task commonModels.Task) (commonModels.Task, *commonStructures.CommonError) {

	err := taskDao.Db.Save(&task).Error
	if err != nil {
		return task, commonUtils.HandleORMError(err)
	}
	return task, nil
}

func (taskDao *TaskDao) CreateTaskWithTx(tx *gorm.DB, task commonModels.Task) (commonModels.Task, *commonStructures.CommonError) {

	err := tx.Create(&task).Error
	if err != nil {
		return task, commonUtils.HandleORMError(err)
	}
	return task, nil
}

func (taskDao *TaskDao) UpdateTaskWithTx(tx *gorm.DB, task commonModels.Task) (commonModels.Task, *commonStructures.CommonError) {

	err := tx.Save(&task).Error
	if err != nil {
		return task, commonUtils.HandleORMError(err)
	}
	return task, nil
}

func (taskDao *TaskDao) DeleteTaskWithTx(tx *gorm.DB, taskId uint) *commonStructures.CommonError {

	updates := map[string]interface{}{
		"deleted_at": commonUtils.GetCurrentTime(),
		"updated_at": commonUtils.GetCurrentTime(),
		"deleted_by": commonConstants.CitadelSystemId,
		"updated_by": commonConstants.CitadelSystemId,
	}
	err := tx.Model(&commonModels.Task{}).
		Where("id = ?", taskId).
		Updates(updates).Error
	if err != nil {
		return commonUtils.HandleORMError(err)
	}
	return nil
}

func (taskDao *TaskDao) UpdateTaskMetadata(taskMetadata commonModels.TaskMetadata) (commonModels.TaskMetadata, *commonStructures.CommonError) {
	err := taskDao.Db.Save(&taskMetadata).Error
	if err != nil {
		return taskMetadata, commonUtils.HandleORMError(err)
	}
	return taskMetadata, nil
}

func (taskDao *TaskDao) CreateTaskMetadataWithTx(tx *gorm.DB, taskMetadata commonModels.TaskMetadata) (commonModels.TaskMetadata, *commonStructures.CommonError) {

	err := tx.Create(&taskMetadata).Error
	if err != nil {
		return taskMetadata, commonUtils.HandleORMError(err)
	}
	return taskMetadata, nil
}

func (taskDao *TaskDao) UpdateTaskMetadataWithTx(tx *gorm.DB, taskMetadata commonModels.TaskMetadata) (commonModels.TaskMetadata, *commonStructures.CommonError) {

	err := tx.Save(&taskMetadata).Error
	if err != nil {
		return taskMetadata, commonUtils.HandleORMError(err)
	}
	return taskMetadata, nil
}

func (taskDao *TaskDao) GetTaskVisitMappingsByTaskId(taskID uint) ([]commonModels.TaskVisitMapping, *commonStructures.CommonError) {

	taskVisitMappings := []commonModels.TaskVisitMapping{}
	err := taskDao.Db.Find(&taskVisitMappings, "task_id = ?", taskID).Error
	if err != nil {
		return taskVisitMappings, commonUtils.HandleORMError(err)
	}
	return taskVisitMappings, nil
}

func (taskDao *TaskDao) GetTaskVisitMappingsByVisitId(visitId string) (commonModels.TaskVisitMapping, *commonStructures.CommonError) {
	taskVisitMapping := commonModels.TaskVisitMapping{}
	err := taskDao.Db.Find(&taskVisitMapping, "visit_id = ?", visitId).Error
	if err != nil {
		return taskVisitMapping, commonUtils.HandleORMError(err)
	}
	return taskVisitMapping, nil
}

func (taskDao *TaskDao) CreateTaskVisitMappingWithTransaction(tx *gorm.DB, taskVisitMappings []commonModels.TaskVisitMapping) ([]commonModels.TaskVisitMapping, *commonStructures.CommonError) {

	err := tx.Create(&taskVisitMappings).Error
	if err != nil {
		return taskVisitMappings, commonUtils.HandleORMError(err)
	}
	return taskVisitMappings, nil
}

func (taskDao *TaskDao) DeleteTaskVisitMappingsByTaskIdAndVisitIds(tx *gorm.DB, taskId uint, toDeleteVisitIds []string) *commonStructures.CommonError {

	err := tx.Where("task_id = ? AND visit_id IN (?)", taskId, toDeleteVisitIds).Delete(&commonModels.TaskVisitMapping{}).Error
	if err != nil {
		return commonUtils.HandleORMError(err)
	}
	return nil
}
