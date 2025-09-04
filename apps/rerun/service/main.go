package service

import (
	"github.com/Orange-Health/citadel/adapters/cache"
	"github.com/Orange-Health/citadel/adapters/sentry"
	"github.com/Orange-Health/citadel/apps/rerun/dao"
	"github.com/Orange-Health/citadel/apps/users/service"
)

type RerunService struct {
	RerunDao    dao.DataLayer
	Cache       cache.CacheLayer
	Sentry      sentry.SentryLayer
	UserService service.UserServiceInterface
}

func InitializeRerunService() RerunServiceInterface {
	return &RerunService{
		RerunDao:    dao.InitializeRerunDao(),
		Cache:       cache.InitializeCache(),
		Sentry:      sentry.InitializeSentry(),
		UserService: service.InitializeUserService(),
	}
}
