package service

import (
	"context"

	commonConstants "github.com/Orange-Health/citadel/common/constants"
	commonStructures "github.com/Orange-Health/citadel/common/structures"
	commonUtils "github.com/Orange-Health/citadel/common/utils"
)

func (etsService *EtsService) filterAlreadySentTatBreachedTests(
	tatBreachedDetails []commonStructures.EtsTestEvent) []commonStructures.EtsTestEvent {

	if len(tatBreachedDetails) == 0 {
		return tatBreachedDetails
	}

	testIdTatBreachedDetailsMap, testIds := make(map[string]commonStructures.EtsTestEvent), []string{}
	for _, tatBreachedObject := range tatBreachedDetails {
		testIdTatBreachedDetailsMap[tatBreachedObject.TestID] = tatBreachedObject
		testIds = append(testIds, tatBreachedObject.TestID)
	}
	alreadySentTatBreachedTests := etsService.EtsDao.FetchActiveTatBreachedTests(testIds)

	toBeSentTatBreachedDetails := []commonStructures.EtsTestEvent{}
	for _, testId := range testIds {
		if !commonUtils.SliceContainsString(alreadySentTatBreachedTests, testId) {
			toBeSentTatBreachedDetails = append(toBeSentTatBreachedDetails, testIdTatBreachedDetailsMap[testId])
		}
	}

	return toBeSentTatBreachedDetails
}

func (etsService *EtsService) HandleTatBreachedTestsCronByEvents(ctx context.Context) {
	inhouseLabIds := etsService.CdsService.GetInhouseLabIds(ctx)
	tatBreachedDetails := etsService.EtsDao.FetchTatBreachDetails(inhouseLabIds)
	tatBreachedDetails = etsService.filterAlreadySentTatBreachedTests(tatBreachedDetails)

	if len(tatBreachedDetails) == 0 {
		return
	}

	etsService.GetPublishAndCreateEtsTestEventForCron(ctx, tatBreachedDetails)
}

func (etsService *EtsService) CreateLisAlertPayloadForEts(ctx context.Context, testIds []string,
	attuneStatus string) []commonStructures.EtsTestEvent {
	inhouseLabIds := etsService.CdsService.GetInhouseLabIds(ctx)
	etsEvents := etsService.EtsDao.FetchLisWebhookTests(testIds, inhouseLabIds)

	if attuneStatus != commonConstants.ATTUNE_TEST_STATUS_RERUN {
		etsEvents = etsService.EtsDao.KeepTestsWhichAreAlreadySent(etsEvents)
	}

	for index := range etsEvents {
		etsEvents[index].LisStatus = attuneStatus
	}

	return etsEvents
}

func (etsService *EtsService) FetchEtsTestEventBasicDetails(ctx context.Context,
	testIds []string) []commonStructures.EtsTestEvent {
	inhouseLabIds := etsService.CdsService.GetInhouseLabIds(ctx)
	etsTestEvents := etsService.EtsDao.FetchEtsTestEventsBasicDetails(testIds, inhouseLabIds)
	return etsService.EtsDao.KeepTestsWhichAreAlreadySent(etsTestEvents)
}

func (etsService *EtsService) FetchEtsTestEventDetailsWhileSampleRejection(ctx context.Context, orderId string,
	sampleNumber uint) []commonStructures.EtsTestEvent {
	inhouseLabIds := etsService.CdsService.GetInhouseLabIds(ctx)
	etsTestEvents := etsService.EtsDao.FetchEtsTestEventDetailsWhileSampleRejection(orderId, sampleNumber, inhouseLabIds)
	return etsService.EtsDao.KeepTestsWhichAreAlreadySent(etsTestEvents)
}

func (etsService *EtsService) FetchEtsTestEventDetailsWhilePartialSampleRejection(ctx context.Context, orderId string,
	sampleNumber uint, testId string) []commonStructures.EtsTestEvent {
	inhouseLabIds := etsService.CdsService.GetInhouseLabIds(ctx)
	etsTestEvents := etsService.EtsDao.FetchEtsTestEventDetailsWhilePartialSampleRejection(testId, orderId, sampleNumber,
		inhouseLabIds)
	return etsService.EtsDao.KeepTestsWhichAreAlreadySent(etsTestEvents)
}

func (etsService *EtsService) GetAndPublishEtsTestEventForSampleRejection(ctx context.Context, omsOrderId string,
	sampleNumber uint) {
	etsTestEvents := etsService.FetchEtsTestEventDetailsWhileSampleRejection(ctx, omsOrderId, sampleNumber)
	etsService.GetPublishAndUpdateEtsTestEvent(ctx, etsTestEvents, true)
}

func (etsService *EtsService) GetAndPublishEtsTestEventForPartialRejection(ctx context.Context, orderId string,
	sampleNumber uint, testId string) {
	etsTestEvents := etsService.FetchEtsTestEventDetailsWhilePartialSampleRejection(ctx, orderId, sampleNumber, testId)
	etsService.GetPublishAndUpdateEtsTestEvent(ctx, etsTestEvents, true)
}

func (etsService *EtsService) GetAndPublishEtsTestEventForLisWebhook(ctx context.Context, testIds []string,
	attuneStatus string) {
	etsTestEvents := etsService.CreateLisAlertPayloadForEts(ctx, testIds, attuneStatus)
	if attuneStatus == commonConstants.ATTUNE_TEST_STATUS_RERUN {
		etsService.GetPublishAndCreateEtsTestEventForRerunWebhook(ctx, etsTestEvents)
		return
	}
	etsService.GetPublishAndUpdateEtsTestEvent(ctx, etsTestEvents, false)
}

func (etsService *EtsService) GetAndPublishEtsTestBasicEvent(ctx context.Context, testIds []string) {
	etsTestEvents := etsService.FetchEtsTestEventBasicDetails(ctx, testIds)
	etsService.GetPublishAndUpdateEtsTestEvent(ctx, etsTestEvents, false)
}
