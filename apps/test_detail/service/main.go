package service

import (
	"github.com/Orange-Health/citadel/adapters/cache"
	"github.com/Orange-Health/citadel/adapters/sentry"
	attuneService "github.com/Orange-Health/citadel/apps/attune/service"
	cdsService "github.com/Orange-Health/citadel/apps/cds/service"
	orderDetailService "github.com/Orange-Health/citadel/apps/order_details/service"
	"github.com/Orange-Health/citadel/apps/test_detail/dao"
)

type TestDetailService struct {
	TestDetailDao      dao.DataLayer
	Cache              cache.CacheLayer
	Sentry             sentry.SentryLayer
	OrderDetailService orderDetailService.OrderDetailsServiceInterface
	CdsService         cdsService.CdsServiceInterface
	AttuneService      attuneService.AttuneServiceInterface
}

func InitializeTestDetailService() TestDetailServiceInterface {
	return &TestDetailService{
		TestDetailDao:      dao.InitializeTestDetailDao(),
		Cache:              cache.InitializeCache(),
		Sentry:             sentry.InitializeSentry(),
		OrderDetailService: orderDetailService.InitializeOrderDetailsService(),
		CdsService:         cdsService.InitializeCdsService(),
		AttuneService:      attuneService.InitializeAttuneService(),
	}
}
