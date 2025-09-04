package service

import (
	"context"

	"github.com/Orange-Health/citadel/adapters/cache"
	"github.com/Orange-Health/citadel/adapters/sentry"
	cdsService "github.com/Orange-Health/citadel/apps/cds/service"
	"github.com/Orange-Health/citadel/apps/ets/dao"
	pubsubService "github.com/Orange-Health/citadel/apps/pubsub/service"
	snsClient "github.com/Orange-Health/citadel/clients/sns"
)

type EtsService struct {
	EtsDao dao.DataLayer
	Cache  cache.CacheLayer
	Sentry sentry.SentryLayer

	CdsService cdsService.CdsServiceInterface

	PubsubService pubsubService.PubsubInterface
	SnsClient     snsClient.SnsClientInterface
}

type EtsServiceInterface interface {
	HandleTatBreachedTestsCronByEvents(ctx context.Context)
	GetAndPublishEtsTestEventForSampleRejection(ctx context.Context, omsOrderId string, sampleNumber uint)
	GetAndPublishEtsTestEventForPartialRejection(ctx context.Context, omsOrderId string, sampleNumber uint, testId string)
	GetAndPublishEtsTestEventForLisWebhook(ctx context.Context, testIds []string, attuneStatus string)
	GetAndPublishEtsTestBasicEvent(ctx context.Context, testIds []string)
}

func InitializeEtsService() EtsServiceInterface {
	return &EtsService{
		EtsDao:        dao.InitializeEtsDao(),
		Cache:         cache.InitializeCache(),
		Sentry:        sentry.InitializeSentry(),
		CdsService:    cdsService.InitializeCdsService(),
		PubsubService: pubsubService.InitializePubsubService(),
		SnsClient:     snsClient.InitializeSnsClient(),
	}
}
