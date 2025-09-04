package service

import (
	"github.com/Orange-Health/citadel/adapters/cache"
	"github.com/Orange-Health/citadel/adapters/sentry"
	"github.com/Orange-Health/citadel/apps/task_metadata/dao"
)

type TaskMetadataService struct {
	taskMetadataDao dao.DataLayer
	Cache           cache.CacheLayer
	Sentry          sentry.SentryLayer
}

func InitializeTaskMetadataService() TaskMetadataServiceInterface {
	return &TaskMetadataService{
		taskMetadataDao: dao.InitializeTaskMetadataDao(),
		Cache:           cache.InitializeCache(),
		Sentry:          sentry.InitializeSentry(),
	}
}
