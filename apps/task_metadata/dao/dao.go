package dao

import (
	commonStructures "github.com/Orange-Health/citadel/common/structures"
	commonUtils "github.com/Orange-Health/citadel/common/utils"
	"github.com/Orange-Health/citadel/models"
)

type DataLayer interface {
	GetTaskMetadataDetails(taskID uint) (models.TaskMetadata, *commonStructures.CommonError)
	UpdateTaskMetadata(taskMetadata models.TaskMetadata) (models.TaskMetadata, *commonStructures.CommonError)
	UpdateLastEventSentAt(taskID uint) *commonStructures.CommonError
}

func (d *TaskMetadataDao) GetTaskMetadataDetails(taskID uint) (models.TaskMetadata, *commonStructures.CommonError) {

	taskMetadata := models.TaskMetadata{}
	if err := d.Db.First(&taskMetadata, "task_id = ?", taskID).Error; err != nil {
		return models.TaskMetadata{}, commonUtils.HandleORMError(err)
	}

	return taskMetadata, nil
}

func (taskMetadataDao *TaskMetadataDao) UpdateTaskMetadata(taskMetadata models.TaskMetadata) (models.TaskMetadata, *commonStructures.CommonError) {

	if err := taskMetadataDao.Db.Model(&taskMetadata).Where("id = ?", taskMetadata.Id).Updates(&taskMetadata).Error; err != nil {
		return models.TaskMetadata{}, commonUtils.HandleORMError(err)
	}

	if err := taskMetadataDao.Db.First(&taskMetadata, taskMetadata.Id).Error; err != nil {
		return models.TaskMetadata{}, commonUtils.HandleORMError(err)
	}

	return taskMetadata, nil

}

func (taskMetadataDao *TaskMetadataDao) UpdateLastEventSentAt(taskID uint) *commonStructures.CommonError {

	if err := taskMetadataDao.Db.Model(&models.TaskMetadata{}).Where("task_id = ?", taskID).
		Update("last_event_sent_at", commonUtils.GetCurrentTime()).Error; err != nil {
		return commonUtils.HandleORMError(err)
	}

	return nil
}
