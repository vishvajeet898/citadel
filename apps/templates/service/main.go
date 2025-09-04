package service

import (
	"github.com/Orange-Health/citadel/adapters/cache"
	"github.com/Orange-Health/citadel/adapters/sentry"
	"github.com/Orange-Health/citadel/apps/templates/dao"
)

type TemplateService struct {
	TemplateDao dao.DataLayer
	Cache       cache.CacheLayer
	Sentry      sentry.SentryLayer
}

func InitializeRemarkService() TemplateServiceInterface {
	return &TemplateService{
		TemplateDao: dao.InitializeTemplateDao(),
		Cache:       cache.InitializeCache(),
		Sentry:      sentry.InitializeSentry(),
	}
}
