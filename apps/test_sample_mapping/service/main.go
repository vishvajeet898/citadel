package service

import (
	"github.com/Orange-Health/citadel/adapters/cache"
	"github.com/Orange-Health/citadel/adapters/sentry"
	cdsService "github.com/Orange-Health/citadel/apps/cds/service"
	"github.com/Orange-Health/citadel/apps/test_sample_mapping/dao"
)

type TestSampleMappingService struct {
	TestSampleMappingDao dao.DataLayer
	Cache                cache.CacheLayer
	Sentry               sentry.SentryLayer
	CdsService           cdsService.CdsServiceInterface
}

func InitializeTestSampleMappingService() TestSampleMappingServiceInterface {
	return &TestSampleMappingService{
		TestSampleMappingDao: dao.InitializeTestSampleMappingDao(),
		Cache:                cache.InitializeCache(),
		Sentry:               sentry.InitializeSentry(),
		CdsService:           cdsService.InitializeCdsService(),
	}
}
