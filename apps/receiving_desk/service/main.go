package service

import (
	"context"

	"github.com/Orange-Health/citadel/adapters/cache"
	"github.com/Orange-Health/citadel/adapters/sentry"
	attuneService "github.com/Orange-Health/citadel/apps/attune/service"
	cdsService "github.com/Orange-Health/citadel/apps/cds/service"
	orderDetailsService "github.com/Orange-Health/citadel/apps/order_details/service"
	patientDetailService "github.com/Orange-Health/citadel/apps/patient_details/service"
	pubsubService "github.com/Orange-Health/citadel/apps/pubsub/service"
	"github.com/Orange-Health/citadel/apps/receiving_desk/dao"
	"github.com/Orange-Health/citadel/apps/receiving_desk/structures"
	sampleService "github.com/Orange-Health/citadel/apps/samples/service"
	taskService "github.com/Orange-Health/citadel/apps/task/service"
	testDetailsService "github.com/Orange-Health/citadel/apps/test_detail/service"
	healthApiClient "github.com/Orange-Health/citadel/clients/health_api"
	partnerApiClient "github.com/Orange-Health/citadel/clients/partner_api"
	snsClient "github.com/Orange-Health/citadel/clients/sns"
	commonStructures "github.com/Orange-Health/citadel/common/structures"
)

type ReceivingDeskService struct {
	ReceivingDeskDao     dao.DataLayer
	Cache                cache.CacheLayer
	Sentry               sentry.SentryLayer
	OrderDetailsService  orderDetailsService.OrderDetailsServiceInterface
	SampleService        sampleService.SampleServiceInterface
	PatientDetailService patientDetailService.PatientDetailServiceInterface
	TaskService          taskService.TaskServiceInterface
	TestDetailsService   testDetailsService.TestDetailServiceInterface
	AttuneService        attuneService.AttuneServiceInterface
	CdsService           cdsService.CdsServiceInterface
	PubsubService        pubsubService.PubsubInterface
	HealthApiClient      healthApiClient.HealthApiClientInterface
	PartnerApiClient     partnerApiClient.PartnerApiClientInterface
	SnsClient            snsClient.SnsClientInterface
}

type ReceivingDeskServiceInterface interface {
	GetCollectedSamples(ctx context.Context, collectedSamplesRequest structures.CollectedSamplesRequest) (
		structures.CollectedSamplesResponse, *commonStructures.CommonError)
	ReceiveAndSyncSamples(ctx context.Context, receiveSamplesRequest structures.ReceiveSamplesRequest) (
		[]string, *commonStructures.CommonError)
	MarkAsNotReceived(ctx context.Context, markAsNotReceivedRequest commonStructures.MarkAsNotReceivedRequest,
		sendEvents bool) *commonStructures.CommonError
}

func InitializeReceivingDeskService() ReceivingDeskServiceInterface {
	return &ReceivingDeskService{
		ReceivingDeskDao:     dao.InitializeReceivingDeskDao(),
		Cache:                cache.InitializeCache(),
		Sentry:               sentry.InitializeSentry(),
		OrderDetailsService:  orderDetailsService.InitializeOrderDetailsService(),
		SampleService:        sampleService.InitializeSampleService(),
		PatientDetailService: patientDetailService.InitializePatientDetailService(),
		TaskService:          taskService.InitializeTaskService(),
		TestDetailsService:   testDetailsService.InitializeTestDetailService(),
		AttuneService:        attuneService.InitializeAttuneService(),
		CdsService:           cdsService.InitializeCdsService(),
		PubsubService:        pubsubService.InitializePubsubService(),
		HealthApiClient:      healthApiClient.InitializeHealthApiClient(),
		PartnerApiClient:     partnerApiClient.InitializePartnerApiClient(),
		SnsClient:            snsClient.InitializeSnsClient(),
	}
}
