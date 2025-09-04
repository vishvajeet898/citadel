package controller

import (
	"github.com/Orange-Health/citadel/apps/investigation_results/service"
)

type InvestigationResult struct {
	InvResService service.InvestigationResultServiceInterface
}

func InitInvestigationResultController() *InvestigationResult {
	return &InvestigationResult{
		InvResService: service.InitializeInvestigationResultService(),
	}
}
