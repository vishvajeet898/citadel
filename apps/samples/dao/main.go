package dao

import (
	"context"

	"gorm.io/gorm"

	"github.com/Orange-Health/citadel/adapters/psql"
	"github.com/Orange-Health/citadel/apps/samples/structures"
	commonStructures "github.com/Orange-Health/citadel/common/structures"
	commonModels "github.com/Orange-Health/citadel/models"
)

type SampleDao struct {
	Db *gorm.DB
}

type DataLayer interface {
	GetSampleForCollectedB2C(omsRequestId string, taskId uint, statusArray []string,
		appendTaskSequence, useSampleNumber bool) ([]commonStructures.SampleInfo, *commonStructures.CommonError)
	GetSamplesByOmsOrderId(omsOrderId string) ([]commonStructures.SampleInfo, *commonStructures.CommonError)
	GetSamplesByOmsOrderIds(omsOrderIds []string) ([]commonStructures.SampleInfo, *commonStructures.CommonError)
	GetSampleByOrderIdAndSampleNumber(omsOrderId string, sampleNumber uint) (commonStructures.SampleInfo,
		*commonStructures.CommonError)
	GetSamplesByOmsOrderIdAndSampleNumbers(omsOrderId string, sampleNumbers []uint) ([]commonModels.Sample,
		[]commonModels.SampleMetadata, *commonStructures.CommonError)
	GetSamplesForTests(testIds []string) ([]commonModels.Sample, []commonModels.SampleMetadata, *commonStructures.CommonError)
	GetSampleDetailsForScheduler(sampleDetailsRequest structures.SampleDetailsRequest) (
		[]structures.VialStructForQuery, *commonStructures.CommonError)
	GetSampleDataBySampleNumberAndTestId(sampleNumber uint, testId string) ([]commonModels.Sample, []commonModels.SampleMetadata,
		*commonStructures.CommonError)
	GetAllSampleTestsBySampleNumber(sampleNumber uint, omsOrderId string) ([]commonModels.TestDetail, *commonStructures.CommonError)
	GetAllSampleTestsBySampleNumbers(sampleNumbers []uint, omsOrderId string) ([]commonModels.TestDetail,
		*commonStructures.CommonError)
	GetCollectedSamples(omsOrderId string, labId uint) ([]commonModels.Sample, []commonModels.SampleMetadata,
		*commonStructures.CommonError)
	GetSampleByBarcodeForReceiving(barcode string) (commonModels.Sample, *commonStructures.CommonError)
	GetSampleDataBySampleId(sampleId uint) (commonModels.Sample, commonModels.SampleMetadata,
		*commonStructures.CommonError)
	GetSampleDataByBarcodeForRejection(barcode string) ([]commonModels.Sample, []commonModels.SampleMetadata,
		*commonStructures.CommonError)
	GetSamplesDataBySampleIds(sampleIds []uint) ([]commonModels.Sample, []commonModels.SampleMetadata,
		*commonStructures.CommonError)
	GetVisitDetailsForTaskByOmsOrderId(omsOrderId string) ([]commonStructures.VisitDetailsForTask,
		*commonStructures.CommonError)
	BarcodesExistsInSystem(barcodes []string) (bool, *commonStructures.CommonError)
	GetTestDetailsBySampleIds(sampleIds []uint) ([]commonModels.TestDetail, *commonStructures.CommonError)
	GetAllTestsAndSampleMappingsBySampleNumbers(sampleNumbers []uint, omsOrderId string) (
		[]commonModels.TestDetail, []commonModels.TestSampleMapping, *commonStructures.CommonError)
	GetTestDetailsForLisEventByVisitId(visitId string) ([]commonStructures.TestDetailsForLisEvent,
		*commonStructures.CommonError)
	GetOmsTestDetailsByVisitId(visitId string) ([]commonStructures.OmsTestDetailsForLis, *commonStructures.CommonError)
	GetVisitIdsByOmsOrderId(omsOrderId string) ([]string, *commonStructures.CommonError)
	GetVisitLabMapByOmsTestIds(omsTestIds []string) (map[string]uint, *commonStructures.CommonError)
	GetSampleByVisitId(visitId string) (commonModels.Sample, *commonStructures.CommonError)
	GetCovidTestSamples(omsOrderId string) []commonStructures.SampleInfo
	GetMaxSampleNumberByOmsOrderId(omsOrderId string) uint
	GetMaxSampleNumberByOmsOrderIdWithTx(tx *gorm.DB, omsOrderId string) uint
	IsSampleCollected(omsOrderId string) (bool, *commonStructures.CommonError)

	BeginTransaction() *gorm.DB
	CreateSampleWithTx(tx *gorm.DB, sampleInfo commonStructures.SampleInfo) (commonStructures.SampleInfo,
		*commonStructures.CommonError)
	AssignTaskSequenceToSamples(tx *gorm.DB, id uint, alnumTestIds []string, isAdditionalTask uint,
		omsRequestId string) *commonStructures.CommonError

	UpdateBulkSamples(samples []commonStructures.SampleInfo) ([]commonStructures.SampleInfo, *commonStructures.CommonError)
	UpdateSamplesAndSamplesMetadata(samples []commonModels.Sample, samplesMetadata []commonModels.SampleMetadata) (
		[]commonModels.Sample, []commonModels.SampleMetadata, *commonStructures.CommonError)
	UpdateSampleAndSampleMetadataWithTx(tx *gorm.DB, sample commonModels.Sample, sampleMetadata commonModels.SampleMetadata) (
		commonModels.Sample, commonModels.SampleMetadata, *commonStructures.CommonError)
	UpdateSamplesAndSamplesMetadataWithTx(tx *gorm.DB, samples []commonModels.Sample,
		samplesMetadata []commonModels.SampleMetadata) (
		[]commonModels.Sample, []commonModels.SampleMetadata, *commonStructures.CommonError)
	UpdateTaskIdByOmsTestIdsWithTx(tx *gorm.DB, omsTestIds []string, taskId uint) *commonStructures.CommonError
	ReMarkSampleDefaultEmedicNotCollected(omsRequestId string, taskSequence uint) *commonStructures.CommonError
	MarkSampleAsEmedicNotCollected(omsRequestId string) *commonStructures.CommonError
	CollectionPortalMarkAccessionAsAccessioned(isWebhook bool, sampleId uint, omsRequestId string) *commonStructures.CommonError
	AddCollectedVolumeToSample(sampleId, volume uint) *commonStructures.CommonError
	UpdateTaskSequenceForSample(omsRequestId string, taskId uint, testIds []string) *commonStructures.CommonError
	UpdateSampleDetailsForReschedule(requestBody structures.UpdateSampleDetailsForRescheduleRequest) *commonStructures.CommonError
	RemapSamplesToNewTaskSequence(orderIdToSampleUpdateMap map[string][]uint, newSequence uint) *commonStructures.CommonError

	DeleteSampleByOrderIdAndSampleNumberWithTx(tx *gorm.DB, omsOrderId string,
		sampleNumber uint) *commonStructures.CommonError
	DeleteAllSamplesDataByOmsOrderIdWithTx(tx *gorm.DB, omsOrderId string) *commonStructures.CommonError
	RemoveSamplesNotLinkedToAnyTests(omsOrderId string) *commonStructures.CommonError

	// Outbound Samples
	CreateInterlabSamplesWithTx(ctx context.Context, tx *gorm.DB, samples []commonModels.Sample,
		samplesMetadata []commonModels.SampleMetadata) *commonStructures.CommonError
	GetInterlabSamplesByParentSampleId(ctx context.Context, parentSampleId uint) ([]commonStructures.SampleInfo,
		*commonStructures.CommonError)

	GetSamplesForDelayedReverseLogistics(normalTat, campTat, inclinicTat, days uint) (
		[]structures.DelayedReverseLogisticsSamplesDbStruct, *commonStructures.CommonError)
	GetSamplesForDelayedInterlabLogistics() ([]commonStructures.SampleInfo, *commonStructures.CommonError)
	GetSrfOrderIds(cityCode string) []string
}

func InitializeSampleDao() DataLayer {
	return &SampleDao{
		Db: psql.GetDbInstance(),
	}
}
