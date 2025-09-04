package service

import (
	"context"

	"github.com/Orange-Health/citadel/adapters/cache"
	"github.com/Orange-Health/citadel/adapters/sentry"
	"github.com/Orange-Health/citadel/apps/attune/dao"
	cdsService "github.com/Orange-Health/citadel/apps/cds/service"
	attuneClient "github.com/Orange-Health/citadel/clients/attune"
	healthApiClient "github.com/Orange-Health/citadel/clients/health_api"
	partnerApiClient "github.com/Orange-Health/citadel/clients/partner_api"
	commonStructures "github.com/Orange-Health/citadel/common/structures"
	commonModels "github.com/Orange-Health/citadel/models"
)

type AttuneService struct {
	AttuneDao        dao.DataLayer
	Cache            cache.CacheLayer
	Sentry           sentry.SentryLayer
	CdsService       cdsService.CdsServiceInterface
	AttuneClient     attuneClient.AttuneClientInterface
	PartnerApiClient partnerApiClient.PartnerApiClientInterface
	HealthApiClient  healthApiClient.HealthApiClientInterface
}

type AttuneServiceInterface interface {
	GetSampleByVisitId(visitId string) (commonModels.Sample, *commonStructures.CommonError)
	GetOrderDetailsByVisitId(visitId string) (commonModels.OrderDetails, *commonStructures.CommonError)
	GetLisSyncData(ctx context.Context, visitId, reportPdfFormat string) (commonStructures.LisSyncDetails,
		commonStructures.AttuneOrderResponse, *commonStructures.CommonError)
	SyncDataToLisByOmsOrderId(ctx context.Context, omsOrderId string, labId uint, samples []commonModels.Sample,
		samplesMetadata []commonModels.SampleMetadata, sampleIdBarcodeMap map[uint]string) (
		string, *commonStructures.CommonError)
	ModifyLisDataPostSyncByOrderId(ctx context.Context, omsOrderId string, labId uint,
		testDetail commonModels.TestDetail, sample commonStructures.SampleInfo) *commonStructures.CommonError
	CancelLisSyncData(ctx context.Context, testDetails []commonModels.TestDetail,
		visitId string, sampleInfo commonStructures.SampleInfo) *commonStructures.CommonError
	UpdateSrfIdToAttune(ctx context.Context, sample commonStructures.SampleInfo, srfId string) *commonStructures.CommonError
}

func InitializeAttuneService() AttuneServiceInterface {
	return &AttuneService{
		AttuneDao:        dao.InitializeAttuneDao(),
		Cache:            cache.InitializeCache(),
		Sentry:           sentry.InitializeSentry(),
		CdsService:       cdsService.InitializeCdsService(),
		AttuneClient:     attuneClient.InitializeAttuneClient(),
		PartnerApiClient: partnerApiClient.InitializePartnerApiClient(),
		HealthApiClient:  healthApiClient.InitializeHealthApiClient(),
	}
}
