package service

import (
	"context"

	"github.com/Orange-Health/citadel/adapters/cache"
	"github.com/Orange-Health/citadel/adapters/sentry"
	cdsService "github.com/Orange-Health/citadel/apps/cds/service"
	"github.com/Orange-Health/citadel/apps/search/dao"
	"github.com/Orange-Health/citadel/apps/search/structures"
	taskService "github.com/Orange-Health/citadel/apps/task/service"
	testDetailService "github.com/Orange-Health/citadel/apps/test_detail/service"
	commonStructures "github.com/Orange-Health/citadel/common/structures"
)

type SearchService struct {
	SearchDao dao.DataLayer
	Cache     cache.CacheLayer
	Sentry    sentry.SentryLayer

	TaskService       taskService.TaskServiceInterface
	TestDetailService testDetailService.TestDetailServiceInterface
	CdsService        cdsService.CdsServiceInterface
}

type SearchServiceInterface interface {
	GetTasksList(ctx context.Context, taskListBasicRequest structures.TaskListBasicRequest) (
		structures.TaskListResponse, *commonStructures.CommonError)
	GetAmendmentTasksList(ctx context.Context, taskListBasicRequest structures.TaskListBasicRequest) (
		structures.AmendmentTaskListResponse, *commonStructures.CommonError)

	GetBarcodeDetails(barcode, searchType string, labId uint) (structures.BarcodeDetailsResponse,
		*commonStructures.CommonError)
	GetInfoScreenDataByBarcode(ctx context.Context, barcode string, labId uint) (structures.InfoScreenSearchResponse,
		*commonStructures.CommonError)
	GetInfoScreenDataByOrderId(ctx context.Context, omsOrderId string, labId uint) (structures.InfoScreenSearchResponse,
		*commonStructures.CommonError)
	GetInfoScreenDataByVisitId(ctx context.Context, visitId string, labId uint) (structures.InfoScreenSearchResponse,
		*commonStructures.CommonError)
}

func InitializeSearchService() SearchServiceInterface {
	return &SearchService{
		SearchDao: dao.InitializeSearchDao(),
		Cache:     cache.InitializeCache(),
		Sentry:    sentry.InitializeSentry(),

		TaskService:       taskService.InitializeTaskService(),
		TestDetailService: testDetailService.InitializeTestDetailService(),
		CdsService:        cdsService.InitializeCdsService(),
	}
}
