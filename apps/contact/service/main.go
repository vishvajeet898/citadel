package service

import (
	"github.com/Orange-Health/citadel/adapters/sentry"
	externalInvestigationResultsService "github.com/Orange-Health/citadel/apps/external_investigation_results/service"
	pubsubService "github.com/Orange-Health/citadel/apps/pubsub/service"
	snsClient "github.com/Orange-Health/citadel/clients/sns"
)

type ContactService struct {
	Sentry                             sentry.SentryLayer
	SnsClient                          snsClient.SnsClientInterface
	PubsubService                      pubsubService.PubsubInterface
	ExternalInvestigationResultService externalInvestigationResultsService.ExternalInvestigationResultServiceInterface
}

func InitializeContactService() ContactServiceInterface {
	return &ContactService{
		Sentry:                             sentry.InitializeSentry(),
		SnsClient:                          snsClient.InitializeSnsClient(),
		PubsubService:                      pubsubService.InitializePubsubService(),
		ExternalInvestigationResultService: externalInvestigationResultsService.InitializeExternalInvestigationResultService(),
	}
}
