package service

import (
	"context"

	"gorm.io/gorm"

	"github.com/Orange-Health/citadel/adapters/cache"
	"github.com/Orange-Health/citadel/adapters/sentry"
	cdsService "github.com/Orange-Health/citadel/apps/cds/service"
	coAuthorizePathologistService "github.com/Orange-Health/citadel/apps/co_authorize_pathologists/service"
	etsService "github.com/Orange-Health/citadel/apps/ets/service"
	investigationResultsService "github.com/Orange-Health/citadel/apps/investigation_results/service"
	remarkService "github.com/Orange-Health/citadel/apps/remarks/service"
	rerunService "github.com/Orange-Health/citadel/apps/rerun/service"
	sampleService "github.com/Orange-Health/citadel/apps/samples/service"
	"github.com/Orange-Health/citadel/apps/task/dao"
	"github.com/Orange-Health/citadel/apps/task/structures"
	taskPathService "github.com/Orange-Health/citadel/apps/task_pathologist_mapping/service"
	testDetailService "github.com/Orange-Health/citadel/apps/test_detail/service"
	userService "github.com/Orange-Health/citadel/apps/users/service"
	attuneClient "github.com/Orange-Health/citadel/clients/attune"
	omsClient "github.com/Orange-Health/citadel/clients/oms"
	commonStructures "github.com/Orange-Health/citadel/common/structures"
	commonModels "github.com/Orange-Health/citadel/models"
)

type TaskService struct {
	TaskDao                       dao.DataLayer
	Cache                         cache.CacheLayer
	Sentry                        sentry.SentryLayer
	TaskPathService               taskPathService.TaskPathologistMappingServiceInterface
	SampleService                 sampleService.SampleServiceInterface
	TestDetailService             testDetailService.TestDetailServiceInterface
	InvestigationResultsService   investigationResultsService.InvestigationResultServiceInterface
	CoAuthorizePathologistService coAuthorizePathologistService.CoAuthorizePathologistInterface
	RerunService                  rerunService.RerunServiceInterface
	RemarkService                 remarkService.RemarkServiceInterface
	EtsService                    etsService.EtsServiceInterface
	CdsService                    cdsService.CdsServiceInterface
	UserService                   userService.UserServiceInterface
	OmsClient                     omsClient.OmsClientInterface
	AttuneClient                  attuneClient.AttuneClientInterface
}

type TaskServiceInterface interface {
	// Tasks
	GetTaskById(taskID uint) (structures.TaskDetail, *commonStructures.CommonError)
	GetTaskModelById(taskID uint) (commonModels.Task, *commonStructures.CommonError)
	GetTaskByOmsOrderId(omsOrderId string) (commonModels.Task, *commonStructures.CommonError)
	GetTasksCount() (uint, *commonStructures.CommonError)
	GetAmendmentTasksCount(ctx context.Context) (uint, *commonStructures.CommonError)
	CreateTaskWithTx(tx *gorm.DB, task commonModels.Task) (commonModels.Task, *commonStructures.CommonError)
	UpdateTask(task commonModels.Task) (commonModels.Task, *commonStructures.CommonError)
	UpdateTaskWithTx(tx *gorm.DB, task commonModels.Task) (commonModels.Task, *commonStructures.CommonError)
	DeleteTaskWithTx(tx *gorm.DB, taskID uint) *commonStructures.CommonError
	GetTaskIdByVisitId(visitId string) (uint, *commonStructures.CommonError)
	UpdateAllTaskDetails(ctx context.Context, taskDetails structures.UpdateAllTaskDetailsStruct) *commonStructures.CommonError
	UndoReportRelease(ctx context.Context, taskID uint) *commonStructures.CommonError
	GetTaskCallingDetails(taskId, userId uint, callingType string) (
		structures.TaskCallingDetailsResponse, *commonStructures.CommonError)
	GetQcFailedTestDataToRerun(ctx context.Context, userId uint, testDetailsIdsToRerun []uint, cityCode string,
		testDetails []commonModels.TestDetail, qcFailedInvestigations []commonModels.InvestigationResult,
		userIdToAttuneUserId map[uint]int) (
		[]commonModels.RerunInvestigationResult, map[string]commonStructures.AttuneOrderResponse, *commonStructures.CommonError)

	// Attune Task Data
	FetchAttuneDataForApprovingTests(ctx context.Context, task commonModels.Task,
		testDetailsToBeAproved []commonModels.TestDetail, investigations []commonModels.InvestigationResult,
		medicalRemarks []commonModels.Remark, approvedByUsers []commonModels.User,
	) (map[string]commonStructures.AttuneOrderResponse, *commonStructures.CommonError)

	// Task Metadata
	GetTaskMetadataByTaskId(taskID uint) (commonModels.TaskMetadata, *commonStructures.CommonError)
	UpdateTaskMetadata(taskMetadata commonModels.TaskMetadata) (commonModels.TaskMetadata, *commonStructures.CommonError)
	CreateTaskMetadataWithTx(tx *gorm.DB, taskMetadata commonModels.TaskMetadata) (commonModels.TaskMetadata, *commonStructures.CommonError)
	UpdateTaskMetadataWithTx(tx *gorm.DB, taskMetadata commonModels.TaskMetadata) (commonModels.TaskMetadata, *commonStructures.CommonError)

	// Task Visit Mapping
	GetTaskVisitMappingsByTaskId(taskID uint) ([]commonModels.TaskVisitMapping, *commonStructures.CommonError)
	CreateTaskVisitMappingWithTx(tx *gorm.DB, taskVisitMappings []commonModels.TaskVisitMapping) ([]commonModels.TaskVisitMapping, *commonStructures.CommonError)
	CreateDeleteTaskVisitMappingsWithTx(tx *gorm.DB, taskId uint, toCreateVisitIds, toDeleteVisitIds []string) *commonStructures.CommonError
}

func InitializeTaskService() TaskServiceInterface {
	return &TaskService{
		TaskDao:                       dao.InitializeTaskDao(),
		Cache:                         cache.InitializeCache(),
		Sentry:                        sentry.InitializeSentry(),
		TaskPathService:               taskPathService.InitializeTaskPathologistMappingService(),
		SampleService:                 sampleService.InitializeSampleService(),
		TestDetailService:             testDetailService.InitializeTestDetailService(),
		InvestigationResultsService:   investigationResultsService.InitializeInvestigationResultService(),
		CoAuthorizePathologistService: coAuthorizePathologistService.InitializeCoAuthorizePathologistService(),
		RerunService:                  rerunService.InitializeRerunService(),
		RemarkService:                 remarkService.InitializeRemarkService(),
		EtsService:                    etsService.InitializeEtsService(),
		CdsService:                    cdsService.InitializeCdsService(),
		UserService:                   userService.InitializeUserService(),
		OmsClient:                     omsClient.InitializeOmsClient(),
		AttuneClient:                  attuneClient.InitializeAttuneClient(),
	}
}
