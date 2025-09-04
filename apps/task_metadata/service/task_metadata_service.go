package service

import (
	mapper "github.com/Orange-Health/citadel/apps/task_metadata/mapper"
	"github.com/Orange-Health/citadel/apps/task_metadata/structures"
	commonStructures "github.com/Orange-Health/citadel/common/structures"
)

type TaskMetadataServiceInterface interface {
	GetTaskMetadataDetails(taskID uint) (structures.TaskMetadata, *commonStructures.CommonError)
	UpdateTaskMetadata(taskMetadata structures.TaskMetadata) (structures.TaskMetadata, *commonStructures.CommonError)
	UpdateLastEventSentAt(taskID uint) *commonStructures.CommonError
}

func (taskMetadataService *TaskMetadataService) GetTaskMetadataDetails(taskID uint) (structures.TaskMetadata, *commonStructures.CommonError) {
	taskMetadata, err := taskMetadataService.taskMetadataDao.GetTaskMetadataDetails(taskID)
	if err != nil {
		return structures.TaskMetadata{}, err
	}
	return mapper.MapTmd(taskMetadata), nil
}

func (taskMetadataService *TaskMetadataService) UpdateTaskMetadata(taskMetadataStruct structures.TaskMetadata) (structures.TaskMetadata, *commonStructures.CommonError) {

	_, err := taskMetadataService.taskMetadataDao.GetTaskMetadataDetails(taskMetadataStruct.TaskID)
	if err != nil {
		return structures.TaskMetadata{}, err
	}

	tmdModel := mapper.MapTmdModel(taskMetadataStruct)
	updatedTmd, err := taskMetadataService.taskMetadataDao.UpdateTaskMetadata(tmdModel)
	if err != nil {
		return structures.TaskMetadata{}, err
	}

	return mapper.MapTmd(updatedTmd), nil
}

func (taskMetadataService *TaskMetadataService) UpdateLastEventSentAt(taskID uint) *commonStructures.CommonError {
	err := taskMetadataService.taskMetadataDao.UpdateLastEventSentAt(taskID)
	if err != nil {
		return err
	}
	return nil
}
