package service

import (
	"github.com/Orange-Health/citadel/adapters/cache"
	"github.com/Orange-Health/citadel/adapters/sentry"
	"github.com/Orange-Health/citadel/apps/users/dao"
)

type UserService struct {
	UserDao dao.DataLayer
	Cache   cache.CacheLayer
	Sentry  sentry.SentryLayer
}

func InitializeUserService() UserServiceInterface {
	return &UserService{
		UserDao: dao.InitializeUserDao(),
		Cache:   cache.InitializeCache(),
		Sentry:  sentry.InitializeSentry(),
	}
}
