package workerService

import (
	"github.com/Orange-Health/citadel/adapters/cache"
	"github.com/Orange-Health/citadel/adapters/sentry"
)

type ReceivingDeskWorkerService struct {
	Cache  cache.CacheLayer
	Sentry sentry.SentryLayer
}

func InitializeWorkerService() ReceivingDeskWorkerServiceInterface {
	return &ReceivingDeskWorkerService{
		Cache:  cache.InitializeCache(),
		Sentry: sentry.InitializeSentry(),
	}
}
