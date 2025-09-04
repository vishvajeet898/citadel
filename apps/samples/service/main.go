package service

import (
	"context"
	"time"

	"gorm.io/gorm"

	"github.com/Orange-Health/citadel/adapters/cache"
	"github.com/Orange-Health/citadel/adapters/sentry"
	attuneService "github.com/Orange-Health/citadel/apps/attune/service"
	cdsService "github.com/Orange-Health/citadel/apps/cds/service"
	etsService "github.com/Orange-Health/citadel/apps/ets/service"
	orderDetailsService "github.com/Orange-Health/citadel/apps/order_details/service"
	patientDetailsService "github.com/Orange-Health/citadel/apps/patient_details/service"
	pubsubService "github.com/Orange-Health/citadel/apps/pubsub/service"
	"github.com/Orange-Health/citadel/apps/samples/dao"
	"github.com/Orange-Health/citadel/apps/samples/structures"
	testDetailsService "github.com/Orange-Health/citadel/apps/test_detail/service"
	tsmService "github.com/Orange-Health/citadel/apps/test_sample_mapping/service"
	accountsApiClient "github.com/Orange-Health/citadel/clients/accounts_api"
	slackClient "github.com/Orange-Health/citadel/clients/slack"
	snsClient "github.com/Orange-Health/citadel/clients/sns"
	commonStructures "github.com/Orange-Health/citadel/common/structures"
	commonModels "github.com/Orange-Health/citadel/models"
)

type SampleService struct {
	SampleDao                dao.DataLayer
	Cache                    cache.CacheLayer
	Sentry                   sentry.SentryLayer
	CdsService               cdsService.CdsServiceInterface
	OrderDetailsService      orderDetailsService.OrderDetailsServiceInterface
	PatientDetailsService    patientDetailsService.PatientDetailServiceInterface
	TestDetailsService       testDetailsService.TestDetailServiceInterface
	TestSampleMappingService tsmService.TestSampleMappingServiceInterface
	AttuneService            attuneService.AttuneServiceInterface
	EtsService               etsService.EtsServiceInterface
	PubsubService            pubsubService.PubsubInterface
	SnsClient                snsClient.SnsClientInterface
	AccountsApiClient        accountsApiClient.AccountsApiClientInterface
	SlackClient              slackClient.SlackClientInterface
}

type SampleServiceInterface interface {
	GetSampleDetailsByOmsOrderIdAndLabId(omsOrderId string) (structures.OrderTestsDetail,
		*commonStructures.CommonError)
	GetSampleDetailsByOmsRequestId(omsRequestId string) ([]structures.OrderTestsDetail, *commonStructures.CommonError)
	GetSampleDetailsByOmsOrderIds(orderIds []string) ([]structures.OrderTestsDetail,
		*commonStructures.CommonError)
	GetSamplesForTests(testIds []string) ([]commonModels.Sample, []commonModels.SampleMetadata, *commonStructures.CommonError)
	GetCollectedSamples(omsOrderId string, labId uint) ([]commonModels.Sample, []commonModels.SampleMetadata,
		*commonStructures.CommonError)
	GetSampleByBarcodeForReceiving(barcode string) (commonModels.Sample, *commonStructures.CommonError)
	GetSampleDataBySampleId(sampleId uint) (commonModels.Sample, commonModels.SampleMetadata, *commonStructures.CommonError)
	GetSamplesDataBySampleIds(sampleIds []uint) ([]commonModels.Sample, []commonModels.SampleMetadata,
		*commonStructures.CommonError)
	GetSampleDetailsForScheduler(sampleDetailsRequest structures.SampleDetailsRequest) (
		map[string][]structures.SampleDetailsResponse, *commonStructures.CommonError)
	GetVisitIdsByOmsOrderId(omsOrderId string) ([]string, *commonStructures.CommonError)
	GetVisitLabMapByOmsTestIds(omsTestIds []string) (map[string]uint, *commonStructures.CommonError)
	BarcodesExistsInSystem(barcodes []string) (bool, *commonStructures.CommonError)
	GetVisitDetailsForTaskByOmsOrderId(omsOrderId string) ([]commonStructures.VisitDetailsForTask,
		*commonStructures.CommonError)
	IsSampleCollected(omsOrderId string) (bool, *commonStructures.CommonError)
	GetTestDetailsForLisEventByVisitId(visitd string) ([]commonStructures.TestDetailsForLisEvent,
		*commonStructures.CommonError)
	GetTestDetailsBySampleIds(sampleIds []uint) ([]commonModels.TestDetail, *commonStructures.CommonError)
	GetAllTestsAndSampleMappingsBySampleNumbers(sampleNumbers []uint, omsOrderId string) (
		[]commonModels.TestDetail, []commonModels.TestSampleMapping, *commonStructures.CommonError)
	GetAllSampleTestsBySampleNumber(sampleNumber uint, omsOrderId string) ([]commonModels.TestDetail,
		*commonStructures.CommonError)
	GetOmsTestDetailsByVisitId(visitId string) ([]commonStructures.OmsTestDetailsForLis, *commonStructures.CommonError)
	GetLisSyncDataByVisitId(ctx context.Context, visitId string) (commonStructures.LisSyncDetails,
		*commonStructures.CommonError)

	CreateSamplesWithTx(ctx context.Context, tx *gorm.DB, omsOrderId string) *commonStructures.CommonError
	CreateSamplesForRecollectionWithTx(ctx context.Context, tx *gorm.DB, orderDetails commonModels.OrderDetails,
		testDetails []commonModels.TestDetail, omsTaskId uint) *commonStructures.CommonError
	CreateSamplesAndTestDetailsWithTx(ctx context.Context, tx *gorm.DB, omsOrderId string,
		createTestDetails, labIdChangeTestDetails []commonModels.TestDetail) *commonStructures.CommonError
	CreateUpdateSamplesPostCollectionWithTx(ctx context.Context, tx *gorm.DB, omsOrderId string,
		createTestDetails, updateTestDetails []commonModels.TestDetail, omsTests []commonStructures.OmsTestModelDetails) (
		map[string]*time.Time, *commonStructures.CommonError)
	SynchronizeTasksWithSamplesWithTx(tx *gorm.DB, omsRequestId string, tasks []commonStructures.OmsTaskModelDetails,
		taskTestsMapping [][]commonStructures.TestsJsonStruct) *commonStructures.CommonError
	UpdateSampleCollected(sampleCollectionRequest commonStructures.SampleCollectedRequest) *commonStructures.CommonError
	ForcefullyMarkSampleAsCollected(request structures.ForcefullyMarkSampleAsCollectedRequest) *commonStructures.CommonError
	AddBarcodeDetails(barcodeDetails structures.AddBarcodesRequest) (map[string]string, *commonStructures.CommonError)
	AddBarcodeDetailsForOrangers(accessionBody structures.UpdateAccessionBody) (
		map[string]string, *commonStructures.CommonError)
	MarkSampleAsEmedicNotCollected(omsRequestId string) *commonStructures.CommonError
	UpdateSampleDetailsPostTaskCompletion(ctx context.Context,
		requestBody commonStructures.UpdateSampleDetailsPostTaskCompletionRequest) *commonStructures.CommonError
	UpdateTaskSequenceForSample(reqBody commonStructures.UpdateTaskSequenceRequest) *commonStructures.CommonError
	RemapSamplesToNewTaskSequence(requestBody structures.RemapSamplesRequest) *commonStructures.CommonError
	UpdateSampleDetailsForReschedule(reqBody structures.UpdateSampleDetailsForRescheduleRequest) *commonStructures.CommonError
	UpdateTaskIdByOmsTestIdsWithTx(tx *gorm.DB, omsTestIds []string, taskId uint) *commonStructures.CommonError
	UpdateSamplesAndSamplesMetadataWithTx(tx *gorm.DB, samples []commonModels.Sample,
		samplesMetadata []commonModels.SampleMetadata) (
		[]commonModels.Sample, []commonModels.SampleMetadata, *commonStructures.CommonError)
	CollectionPortalMarkAccessionAsAccessioned(
		requestBody structures.MarkAccessionAsAccessionedRequest) *commonStructures.CommonError
	AddCollectedVolumeToSample(requestBody structures.AddCollectedVolumneRequest) (
		structures.AddVolumeResponse, *commonStructures.CommonError)
	RemoveSamplesNotLinkedToAnyTests(omsOrderId string) *commonStructures.CommonError
	RejectSampleByBarcode(ctx context.Context, barcode string, requestBody structures.RejectSampleRequest) (string, []string,
		*commonStructures.CommonError)
	RejectSamplePartiallyBySampleNumberAndTestId(ctx context.Context, requestBody structures.RejectSamplePartiallyRequest) (
		string, []string, *commonStructures.CommonError)
	DeleteSamplesAndTestDetailsWithTx(ctx context.Context, tx *gorm.DB, orderId string,
		toDeleteTestDetails []commonModels.TestDetail) *commonStructures.CommonError
	DeleteTestSampleMappingForDeletedTestIds(ctx context.Context, tx *gorm.DB, omsOrderId string,
		toDeleteTestDetails []commonModels.TestDetail) *commonStructures.CommonError
	DeleteAllSamplesDataByOmsOrderIdWithTx(tx *gorm.DB, omsOrderId string) *commonStructures.CommonError
	UpdateSrfIdToLis(ctx context.Context, orderDetails commonModels.OrderDetails) *commonStructures.CommonError

	// Publish events functions
	PublishResetTatsEvent(omsTestIds []string, cityCode string)
	PublishUpdateTestStatusEvent(omsTestStatusMap map[string]string, omsOrderId string, checkOrderCompletion bool,
		cityCode string)
	PublishRemoveSampleRejectedTagEvent(omsRequestId string, omsOrderIds []string, removeRequestTag bool, cityCode string)
	PublishLisDataEvent(attuneOrderResponse commonStructures.AttuneOrderResponse)

	// Outbound Samples
	CreateInterlabSamplesWithTx(ctx context.Context, tx *gorm.DB, samples []commonModels.Sample,
		samplesMetadata []commonModels.SampleMetadata, sampleNumberToLabIdMap map[uint]uint) *commonStructures.CommonError

	GetSamplesForDelayedReverseLogisticsDashboard(ctx context.Context,
		cityCode string) []structures.DelayedReverseLogisticsSamplesResponse
	GetSrfOrderIds(ctx context.Context, cityCode string) []string
	SendSlackAlertForDelayedReverseLogisticsTat(ctx context.Context)
	SendSlackAlertForDelayedInterlabLogisticsTat(ctx context.Context)
}

func InitializeSampleService() SampleServiceInterface {
	return &SampleService{
		SampleDao:                dao.InitializeSampleDao(),
		Cache:                    cache.InitializeCache(),
		Sentry:                   sentry.InitializeSentry(),
		OrderDetailsService:      orderDetailsService.InitializeOrderDetailsService(),
		PatientDetailsService:    patientDetailsService.InitializePatientDetailService(),
		TestDetailsService:       testDetailsService.InitializeTestDetailService(),
		CdsService:               cdsService.InitializeCdsService(),
		TestSampleMappingService: tsmService.InitializeTestSampleMappingService(),
		AttuneService:            attuneService.InitializeAttuneService(),
		EtsService:               etsService.InitializeEtsService(),
		PubsubService:            pubsubService.InitializePubsubService(),
		SnsClient:                snsClient.InitializeSnsClient(),
		AccountsApiClient:        accountsApiClient.InitializeAccountsApiClient(),
		SlackClient:              slackClient.InitializeSlackClient(),
	}
}
