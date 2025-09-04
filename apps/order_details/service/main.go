package service

import (
	"github.com/Orange-Health/citadel/adapters/cache"
	"github.com/Orange-Health/citadel/adapters/sentry"
	"github.com/Orange-Health/citadel/apps/order_details/dao"
)

type OrderDetailsService struct {
	Dao    dao.DataLayer
	Cache  cache.CacheLayer
	Sentry sentry.SentryLayer
}

func InitializeOrderDetailsService() OrderDetailsServiceInterface {
	return &OrderDetailsService{
		Dao:    dao.InitializeOrderDetailsDao(),
		Cache:  cache.InitializeCache(),
		Sentry: sentry.InitializeSentry(),
	}
}
