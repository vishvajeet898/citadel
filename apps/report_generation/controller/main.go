package controller

import (
	"github.com/Orange-Health/citadel/apps/report_generation/service"
)

type ReportGenerationEvent struct {
	S service.ReportGenerationInterface
}

func InitReportGenerationController() *ReportGenerationEvent {
	return &ReportGenerationEvent{
		S: service.InitializeReportGenerationService(),
	}
}
