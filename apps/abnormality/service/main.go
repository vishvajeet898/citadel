package abnormalityService

import (
	"github.com/Orange-Health/citadel/adapters/cache"
	"github.com/Orange-Health/citadel/adapters/sentry"
)

type AbnormalityService struct {
	Cache  cache.CacheLayer
	Sentry sentry.SentryLayer
}

func InitializeAbnormalityService() AbnormalityServiceInterface {
	return &AbnormalityService{
		Cache:  cache.InitializeCache(),
		Sentry: sentry.InitializeSentry(),
	}
}
