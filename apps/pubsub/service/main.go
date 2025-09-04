package service

import (
	"github.com/Orange-Health/citadel/adapters/cache"
	"github.com/Orange-Health/citadel/adapters/sentry"
)

type PubsubService struct {
	Cache  cache.CacheLayer
	Sentry sentry.SentryLayer
}

func InitializePubsubService() PubsubInterface {
	return &PubsubService{
		Cache:  cache.InitializeCache(),
		Sentry: sentry.InitializeSentry(),
	}
}
