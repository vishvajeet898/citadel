package controller

import (
	"github.com/Orange-Health/citadel/apps/task/service"
	workerService "github.com/Orange-Health/citadel/apps/task/worker_service"
)

type Task struct {
	TaskService       service.TaskServiceInterface
	TaskWorkerService workerService.TaskWorkerServiceInterface
}

func InitTaskController() *Task {
	return &Task{
		TaskService:       service.InitializeTaskService(),
		TaskWorkerService: workerService.InitializeWorkerService(),
	}
}
