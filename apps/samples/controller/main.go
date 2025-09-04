package controller

import (
	"github.com/Orange-Health/citadel/apps/samples/service"
	workerService "github.com/Orange-Health/citadel/apps/samples/worker_service"
)

type Sample struct {
	SampleService       service.SampleServiceInterface
	SampleWorkerService workerService.SampleWorkerServiceInterface
}

func InitSampleController() *Sample {
	return &Sample{
		SampleService:       service.InitializeSampleService(),
		SampleWorkerService: workerService.InitializeWorkerService(),
	}
}
