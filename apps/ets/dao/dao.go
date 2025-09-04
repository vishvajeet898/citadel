package dao

import (
	"time"

	"gorm.io/gorm"

	mappers "github.com/Orange-Health/citadel/apps/ets/mappers"
	commonConstants "github.com/Orange-Health/citadel/common/constants"
	"github.com/Orange-Health/citadel/common/structures"
	commonStructures "github.com/Orange-Health/citadel/common/structures"
	commonModels "github.com/Orange-Health/citadel/models"
)

var SelectStringsForEtsEvent = []string{
	"test_details.central_oms_test_id AS test_id",
	"test_details.oms_order_id AS order_id",
	"test_details.status AS test_status",
	"test_details.oms_status as oms_test_status",
	"test_details.test_name AS test_name",
	"test_details.deleted_at AS test_deleted_at",
	"test_details.master_test_id AS master_test_id",
	"test_details.lab_eta AS lab_eta",
	"test_details.lab_id AS lab_id",
	"patient_details.name AS patient_name",
	"patient_details.dob AS patient_dob",
	"patient_details.expected_dob AS patient_expected_dob",
	"patient_details.gender AS patient_gender",
	"samples.barcode AS barcode",
	"samples.vial_type_id AS vial_type_id",
	"test_sample_mapping.is_rejected AS is_rejected",
	"order_details.city_code AS city_code",
}

func (etsDao *EtsDao) GetEtsCommonQuery(inhouseLabIds []uint) *gorm.DB {
	query := etsDao.Db.Table("test_details").
		Joins("INNER JOIN order_details ON order_details.oms_order_id = test_details.oms_order_id").
		Joins("INNER JOIN patient_details ON patient_details.id = order_details.patient_details_id").
		Joins("INNER JOIN test_sample_mapping ON test_sample_mapping.oms_test_id = test_details.central_oms_test_id AND test_sample_mapping.oms_order_id = order_details.oms_order_id").
		Joins("INNER JOIN samples ON samples.oms_order_id = test_details.oms_order_id AND samples.sample_number = test_sample_mapping.sample_number").
		Select(SelectStringsForEtsEvent).
		Where("test_details.processing_lab_id IN (?)", inhouseLabIds)

	return query
}

func (etsDao *EtsDao) CreateEtsEvents(etsEvents []commonModels.EtsEvent) error {
	if len(etsEvents) == 0 {
		return nil
	}

	if err := etsDao.Db.Model(&commonModels.EtsEvent{}).Create(&etsEvents).Error; err != nil {
		return err
	}

	return nil
}

func (etsDao *EtsDao) MarkEventAsInactive(testIds []string) error {
	if len(testIds) == 0 {
		return nil
	}

	currentTime := time.Now()
	err := etsDao.Db.Model(&commonModels.EtsEvent{}).
		Where("test_id IN (?)", testIds).
		Updates(map[string]interface{}{
			"updated_by": commonConstants.CitadelSystemId,
			"updated_at": currentTime,
			"is_active":  false,
		}).Error
	if err != nil {
		return err
	}

	return err
}

func (etsDao *EtsDao) GetEtsEventByTestId(testId string) commonModels.EtsEvent {
	etsEvent := commonModels.EtsEvent{}
	err := etsDao.Db.Where("test_id = ? and is_active = true", testId).First(&etsEvent).Error
	if err != nil {
		return etsEvent
	}

	return etsEvent
}

func (etsDao *EtsDao) FetchTatBreachDetails(inhouseLabIds []uint) []commonStructures.EtsTestEvent {
	tatBreachedDetails := []commonStructures.EtsTestDbEvent{}
	tatBreachQuery := etsDao.GetEtsCommonQuery(inhouseLabIds)
	tatBreachedInterval := commonConstants.EtsTatBreachedInterval
	currentTime := time.Now()
	pastTime := currentTime.Add(-1 * time.Duration(5) * time.Minute)
	futureTime := currentTime.Add(time.Duration(tatBreachedInterval) * time.Minute)
	tatBreachQuery.Where("test_details.lab_eta BETWEEN ? AND ?", pastTime, futureTime).
		Where("test_details.status = ?", commonConstants.TEST_STATUS_RESULT_PENDING).
		Where("test_sample_mapping.is_rejected = ?", false).Scan(&tatBreachedDetails)

	return mappers.MapBulkEtsTestDbEventToEtsTestDbEventAndFilterOutsourcedTests(tatBreachedDetails)
}

func (etsDao *EtsDao) FetchActiveTatBreachedTests(testIds []string) []string {
	if len(testIds) == 0 {
		return []string{}
	}

	activeTatBreachedTests := []string{}
	etsDao.Db.Model(&commonModels.EtsEvent{}).
		Where("test_id IN (?) and is_active = true", testIds).
		Pluck("test_id", &activeTatBreachedTests)

	return activeTatBreachedTests
}

func (etsDao *EtsDao) FetchLisWebhookTests(testIds []string, inhouseLabIds []uint) []commonStructures.EtsTestEvent {
	if len(testIds) == 0 {
		return []commonStructures.EtsTestEvent{}
	}

	lisWebhookTests := []commonStructures.EtsTestDbEvent{}
	lisWebhookQuery := etsDao.GetEtsCommonQuery(inhouseLabIds)
	lisWebhookQuery.Where("test_details.central_oms_test_id IN (?)", testIds).
		Where("test_sample_mapping.is_rejected = ?", false).Scan(&lisWebhookTests)

	return mappers.MapBulkEtsTestDbEventToEtsTestDbEventAndFilterOutsourcedTests(lisWebhookTests)
}

func (etsDao *EtsDao) FetchEtsTestEventsBasicDetails(testIds []string, inhouseLabIds []uint) []commonStructures.EtsTestEvent {
	if len(testIds) == 0 {
		return []commonStructures.EtsTestEvent{}
	}

	etsTestEvents := []commonStructures.EtsTestDbEvent{}
	etsTestEventsQuery := etsDao.GetEtsCommonQuery(inhouseLabIds)
	etsTestEventsQuery.Where("test_details.central_oms_test_id IN (?)", testIds).Scan(&etsTestEvents)

	return mappers.MapBulkEtsTestDbEventToEtsTestDbEventAndFilterOutsourcedTests(etsTestEvents)
}

func (etsDao *EtsDao) FetchEtsTestEventDetailsWhileSampleRejection(orderId string, sampleNumber uint,
	inhouseLabIds []uint) []commonStructures.EtsTestEvent {
	etsTestEvents := []commonStructures.EtsTestDbEvent{}
	etsDao.GetEtsCommonQuery(inhouseLabIds).
		Where("test_details.oms_order_id = ? and samples.sample_number = ?", orderId, sampleNumber).Scan(&etsTestEvents)
	return mappers.MapBulkEtsTestDbEventToEtsTestDbEventAndFilterOutsourcedTests(etsTestEvents)
}

func (etsDao *EtsDao) FetchEtsTestEventDetailsWhilePartialSampleRejection(testId, orderId string, sampleNumber uint,
	inhouseLabIds []uint) []commonStructures.EtsTestEvent {
	etsTestEvents := []commonStructures.EtsTestDbEvent{}
	etsDao.GetEtsCommonQuery(inhouseLabIds).
		Where("test_details.central_oms_test_id = ? AND test_details.oms_order_id = ?", testId, orderId).
		Where("samples.sample_number = ?", sampleNumber).
		Scan(&etsTestEvents)
	return mappers.MapBulkEtsTestDbEventToEtsTestDbEventAndFilterOutsourcedTests(etsTestEvents)
}

func (etsDao *EtsDao) KeepTestsWhichAreAlreadySent(
	etsTestEvents []commonStructures.EtsTestEvent) []commonStructures.EtsTestEvent {
	if len(etsTestEvents) == 0 {
		return etsTestEvents
	}

	testIdTatBreachedDetailsMap, testIds := map[string]commonStructures.EtsTestEvent{}, []string{}
	for _, tatBreachedObject := range etsTestEvents {
		testIdTatBreachedDetailsMap[tatBreachedObject.TestID] = tatBreachedObject
		testIds = append(testIds, tatBreachedObject.TestID)
	}

	alreadySentTests := []string{}
	etsDao.Db.Model(&commonModels.EtsEvent{}).
		Where("test_id IN (?) and is_active = ?", testIds, true).
		Pluck("test_id", &alreadySentTests)

	toBeSentTestEvents := []structures.EtsTestEvent{}
	for _, testId := range alreadySentTests {
		toBeSentTestEvents = append(toBeSentTestEvents, testIdTatBreachedDetailsMap[testId])
	}

	return toBeSentTestEvents
}
