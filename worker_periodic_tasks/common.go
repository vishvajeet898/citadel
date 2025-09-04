package workerPeriodicTasks

import (
	"gorm.io/gorm"

	"github.com/Orange-Health/citadel/adapters/cache"
	"github.com/Orange-Health/citadel/adapters/sentry"
	etsService "github.com/Orange-Health/citadel/apps/ets/service"
	sampleService "github.com/Orange-Health/citadel/apps/samples/service"
)

type WorkerPeriodicTaskService struct {
	Db     *gorm.DB
	Cache  cache.CacheLayer
	Sentry sentry.SentryLayer

	EtsService    etsService.EtsServiceInterface
	SampleService sampleService.SampleServiceInterface
}
