package controller

import (
	"github.com/Orange-Health/citadel/apps/receiving_desk/service"
	workerService "github.com/Orange-Health/citadel/apps/receiving_desk/worker_service"
)

type ReceivingDesk struct {
	ReceivingDeskService       service.ReceivingDeskServiceInterface
	ReceivingDeskWorkerService workerService.ReceivingDeskWorkerServiceInterface
}

func InitReceivingDeskController() *ReceivingDesk {
	return &ReceivingDesk{
		ReceivingDeskService:       service.InitializeReceivingDeskService(),
		ReceivingDeskWorkerService: workerService.InitializeWorkerService(),
	}
}
