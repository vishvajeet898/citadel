package calculationsService

import (
	"github.com/Orange-Health/citadel/adapters/cache"
	"github.com/Orange-Health/citadel/adapters/sentry"
)

type CalculationsService struct {
	Cache  cache.CacheLayer
	Sentry sentry.SentryLayer
}

func InitializeCalculationsService() CalculationsServiceInterface {
	return &CalculationsService{
		Cache:  cache.InitializeCache(),
		Sentry: sentry.InitializeSentry(),
	}
}
