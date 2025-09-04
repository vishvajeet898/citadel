package workerService

import (
	"github.com/Orange-Health/citadel/adapters/cache"
	"github.com/Orange-Health/citadel/adapters/sentry"
)

type SampleWorkerService struct {
	Cache  cache.CacheLayer
	Sentry sentry.SentryLayer
}

func InitializeWorkerService() SampleWorkerServiceInterface {
	return &SampleWorkerService{
		Cache:  cache.InitializeCache(),
		Sentry: sentry.InitializeSentry(),
	}
}
