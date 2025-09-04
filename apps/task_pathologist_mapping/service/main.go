package service

import (
	"github.com/Orange-Health/citadel/adapters/cache"
	"github.com/Orange-Health/citadel/adapters/sentry"
	"github.com/Orange-Health/citadel/apps/task_pathologist_mapping/dao"
)

type TaskPathologistMappingService struct {
	TaskPathDao dao.DataLayer
	Cache       cache.CacheLayer
	Sentry      sentry.SentryLayer
}

func InitializeTaskPathologistMappingService() TaskPathologistMappingServiceInterface {
	return &TaskPathologistMappingService{
		TaskPathDao: dao.InitializeTaskPathMapDao(),
		Cache:       cache.InitializeCache(),
		Sentry:      sentry.InitializeSentry(),
	}
}
