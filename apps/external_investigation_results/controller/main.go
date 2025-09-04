package controller

import "github.com/Orange-Health/citadel/apps/external_investigation_results/service"

type ExternalInvestigationResult struct {
	ExtInvResService service.ExternalInvestigationResultServiceInterface
}

func InitExternalInvestigationResultController() *ExternalInvestigationResult {
	return &ExternalInvestigationResult{
		ExtInvResService: service.InitializeExternalInvestigationResultService(),
	}
}
