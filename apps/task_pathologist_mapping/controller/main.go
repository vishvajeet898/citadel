package controller

import (
	"github.com/Orange-Health/citadel/apps/task_pathologist_mapping/service"
)

type TaskPathologistMapping struct {
	TaskPathService service.TaskPathologistMappingServiceInterface
}

func InitTaskPathologistMappingController() *TaskPathologistMapping {
	return &TaskPathologistMapping{
		TaskPathService: service.InitializeTaskPathologistMappingService(),
	}
}
