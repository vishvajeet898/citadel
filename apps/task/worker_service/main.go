package workerService

import (
	"github.com/Orange-Health/citadel/adapters/cache"
	"github.com/Orange-Health/citadel/adapters/sentry"
)

type TaskWorkerService struct {
	Cache  cache.CacheLayer
	Sentry sentry.SentryLayer
}

func InitializeWorkerService() TaskWorkerServiceInterface {
	return &TaskWorkerService{
		Cache:  cache.InitializeCache(),
		Sentry: sentry.InitializeSentry(),
	}
}
