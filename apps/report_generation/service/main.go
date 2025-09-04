package service

import (
	"github.com/Orange-Health/citadel/adapters/cache"
	"github.com/Orange-Health/citadel/adapters/sentry"
	attachmentsService "github.com/Orange-Health/citadel/apps/attachments/service"
	cdsService "github.com/Orange-Health/citadel/apps/cds/service"
	orderDetailsService "github.com/Orange-Health/citadel/apps/order_details/service"
	pubsubService "github.com/Orange-Health/citadel/apps/pubsub/service"
	"github.com/Orange-Health/citadel/apps/report_generation/dao"
	taskService "github.com/Orange-Health/citadel/apps/task/service"
	testDetailService "github.com/Orange-Health/citadel/apps/test_detail/service"
	healthApiClient "github.com/Orange-Health/citadel/clients/health_api"
	omsClient "github.com/Orange-Health/citadel/clients/oms"
	snsClient "github.com/Orange-Health/citadel/clients/sns"
)

type ReportGenerationService struct {
	Dao                 dao.DataLayer
	Cache               cache.CacheLayer
	Sentry              sentry.SentryLayer
	OrderDetailsService orderDetailsService.OrderDetailsServiceInterface
	TestDetailService   testDetailService.TestDetailServiceInterface
	TaskService         taskService.TaskServiceInterface
	AttachmentService   attachmentsService.AttachmentServiceInterface
	CdsService          cdsService.CdsServiceInterface
	PubsubService       pubsubService.PubsubInterface
	SnsClient           snsClient.SnsClientInterface
	OmsClient           omsClient.OmsClientInterface
	HealthApiClient     healthApiClient.HealthApiClientInterface
}

func InitializeReportGenerationService() ReportGenerationInterface {
	return &ReportGenerationService{
		Dao:                 dao.InitializeReportGenerationDao(),
		Cache:               cache.InitializeCache(),
		Sentry:              sentry.InitializeSentry(),
		OrderDetailsService: orderDetailsService.InitializeOrderDetailsService(),
		TestDetailService:   testDetailService.InitializeTestDetailService(),
		TaskService:         taskService.InitializeTaskService(),
		AttachmentService:   attachmentsService.InitializeAttachmentService(),
		CdsService:          cdsService.InitializeCdsService(),
		PubsubService:       pubsubService.InitializePubsubService(),
		SnsClient:           snsClient.InitializeSnsClient(),
		OmsClient:           omsClient.InitializeOmsClient(),
		HealthApiClient:     healthApiClient.InitializeHealthApiClient(),
	}
}
