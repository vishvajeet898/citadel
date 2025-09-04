package controller

import (
	"github.com/Orange-Health/citadel/apps/task_metadata/service"
)

type TaskMetadata struct {
	TaskMetadataService service.TaskMetadataServiceInterface
}

func InitTaskMetadataController() *TaskMetadata {
	return &TaskMetadata{
		TaskMetadataService: service.InitializeTaskMetadataService(),
	}
}
