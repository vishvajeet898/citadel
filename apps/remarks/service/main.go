package service

import (
	"github.com/Orange-Health/citadel/adapters/cache"
	"github.com/Orange-Health/citadel/adapters/sentry"
	"github.com/Orange-Health/citadel/apps/remarks/dao"
)

type RemarkService struct {
	RemarkDao dao.DataLayer
	Cache     cache.CacheLayer
	Sentry    sentry.SentryLayer
}

func InitializeRemarkService() RemarkServiceInterface {
	return &RemarkService{
		RemarkDao: dao.InitializeRemarkDao(),
		Cache:     cache.InitializeCache(),
		Sentry:    sentry.InitializeSentry(),
	}
}
