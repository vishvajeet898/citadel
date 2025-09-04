package dao

import (
	"gorm.io/gorm"

	"github.com/Orange-Health/citadel/adapters/psql"

	commonStructures "github.com/Orange-Health/citadel/common/structures"
	commonModels "github.com/Orange-Health/citadel/models"
)

type EtsDao struct {
	Db *gorm.DB
}

type DataLayer interface {
	CreateEtsEvents(etsEvents []commonModels.EtsEvent) error
	MarkEventAsInactive(testIds []string) error
	GetEtsEventByTestId(testId string) commonModels.EtsEvent

	FetchTatBreachDetails(inhouseLabIds []uint) []commonStructures.EtsTestEvent
	FetchLisWebhookTests(testIds []string, inhouseLabIds []uint) []commonStructures.EtsTestEvent
	FetchEtsTestEventsBasicDetails(testIds []string, inhouseLabIds []uint) []commonStructures.EtsTestEvent
	FetchEtsTestEventDetailsWhileSampleRejection(omsOrderId string, sampleNumber uint,
		inhouseLabIds []uint) []commonStructures.EtsTestEvent
	FetchEtsTestEventDetailsWhilePartialSampleRejection(testId, orderId string, sampleNumber uint,
		inhouseLabIds []uint) []commonStructures.EtsTestEvent
	KeepTestsWhichAreAlreadySent(etsEvents []commonStructures.EtsTestEvent) []commonStructures.EtsTestEvent
	FetchActiveTatBreachedTests(testIds []string) []string
}

func InitializeEtsDao() DataLayer {
	return &EtsDao{
		Db: psql.GetDbInstance(),
	}
}
