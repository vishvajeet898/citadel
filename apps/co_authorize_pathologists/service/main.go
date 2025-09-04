package service

import (
	"github.com/Orange-Health/citadel/adapters/cache"
	"github.com/Orange-Health/citadel/adapters/sentry"
	"github.com/Orange-Health/citadel/apps/co_authorize_pathologists/dao"
)

type CoAuthorizePathologistService struct {
	CoAuthorizePathDao dao.DataLayer
	Cache              cache.CacheLayer
	Sentry             sentry.SentryLayer
}

func InitializeCoAuthorizePathologistService() CoAuthorizePathologistInterface {
	return &CoAuthorizePathologistService{
		CoAuthorizePathDao: dao.InitializeoAuthorizePathologistDao(),
		Cache:              cache.InitializeCache(),
		Sentry:             sentry.InitializeSentry(),
	}
}
