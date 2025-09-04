package worker

import (
	"context"

	"github.com/RichardKnop/machinery/v1"
	"github.com/RichardKnop/machinery/v1/config"
	"github.com/RichardKnop/machinery/v1/log"
	"github.com/RichardKnop/machinery/v1/tasks"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"

	"github.com/Orange-Health/citadel/adapters/cache"
	"github.com/Orange-Health/citadel/adapters/psql"
	"github.com/Orange-Health/citadel/adapters/sentry"
	"github.com/Orange-Health/citadel/adapters/sqs"
	abnormalityService "github.com/Orange-Health/citadel/apps/abnormality/service"
	attachmentsService "github.com/Orange-Health/citadel/apps/attachments/service"
	attuneService "github.com/Orange-Health/citadel/apps/attune/service"
	cdsService "github.com/Orange-Health/citadel/apps/cds/service"
	contactService "github.com/Orange-Health/citadel/apps/contact/service"
	etsService "github.com/Orange-Health/citadel/apps/ets/service"
	investigationResultsService "github.com/Orange-Health/citadel/apps/investigation_results/service"
	orderDetailsService "github.com/Orange-Health/citadel/apps/order_details/service"
	patientDetailService "github.com/Orange-Health/citadel/apps/patient_details/service"
	pubsubService "github.com/Orange-Health/citadel/apps/pubsub/service"
	receivingDeskService "github.com/Orange-Health/citadel/apps/receiving_desk/service"
	remarksService "github.com/Orange-Health/citadel/apps/remarks/service"
	reportGenerationService "github.com/Orange-Health/citadel/apps/report_generation/service"
	rerunService "github.com/Orange-Health/citadel/apps/rerun/service"
	sampleService "github.com/Orange-Health/citadel/apps/samples/service"
	taskService "github.com/Orange-Health/citadel/apps/task/service"
	taskPathMappingService "github.com/Orange-Health/citadel/apps/task_pathologist_mapping/service"
	testDetailService "github.com/Orange-Health/citadel/apps/test_detail/service"
	testSampleMappingService "github.com/Orange-Health/citadel/apps/test_sample_mapping/service"
	userService "github.com/Orange-Health/citadel/apps/users/service"
	attuneClient "github.com/Orange-Health/citadel/clients/attune"
	cdsClient "github.com/Orange-Health/citadel/clients/cds"
	healthApiClient "github.com/Orange-Health/citadel/clients/health_api"
	omsClient "github.com/Orange-Health/citadel/clients/oms"
	partnerApiClient "github.com/Orange-Health/citadel/clients/partner_api"
	reportRebrandingClient "github.com/Orange-Health/citadel/clients/report_rebranding"
	s3Client "github.com/Orange-Health/citadel/clients/s3"
	s3wrapperClient "github.com/Orange-Health/citadel/clients/s3wrapper"
	slackClient "github.com/Orange-Health/citadel/clients/slack"
	snsClient "github.com/Orange-Health/citadel/clients/sns"
	"github.com/Orange-Health/citadel/common/constants"
	commonTasks "github.com/Orange-Health/citadel/common_tasks"
	workerPeriodicTasks "github.com/Orange-Health/citadel/worker_periodic_tasks"
	workerTasks "github.com/Orange-Health/citadel/worker_tasks"
)

func StartServer(clientMode bool) (*machinery.Server, error) {
	sqsClient := sqs.InitAndGetSQSClient()
	redisMaxActive := 1
	redisMaxIdle := 0
	if !clientMode {
		redisMaxActive = constants.WorkerRedisPoolSize
		redisMaxIdle = constants.WorkerMaxIdle
	}
	var visibilityTimeout = 43200
	var cnf = &config.Config{
		Broker:        constants.WorkerBroker,
		DefaultQueue:  constants.WorkerDefaultQueue,
		ResultBackend: constants.WorkerResultBackend,
		Redis: &config.RedisConfig{
			MaxActive: redisMaxActive,
			MaxIdle:   redisMaxIdle,
		},
		SQS: &config.SQSConfig{
			Client:            sqsClient,
			WaitTimeSeconds:   20,
			VisibilityTimeout: &visibilityTimeout,
		},
	}

	if constants.Environment == "local" {
		cnf.AMQP = &config.AMQPConfig{
			Exchange:     "machinery_exchange",
			ExchangeType: "direct",
			BindingKey:   "machinery_task",
		}
	}

	server, err := machinery.NewServer(cnf)
	if err != nil {
		return nil, err
	}

	return server, err
}

func Worker(taskServer *machinery.Server) error {
	consumerTag, processes := constants.CitadelWorker, constants.WorkerProcesses

	// The second argument is a consumer tag
	// Ideally, each worker should have a unique tag (worker1, worker2 etc)
	worker := taskServer.NewWorker(consumerTag, processes)

	// Here we inject some custom code for error handling,
	// start and end of task hooks, useful for metrics for example.
	errorHandler := func(err error) {
		log.ERROR.Println("Error:", err)
	}

	preTaskHandler := func(signature *tasks.Signature) {
		log.INFO.Println("Starting task:", signature.Name)
	}

	postTaskHandler := func(signature *tasks.Signature) {
		log.INFO.Println("Ending task:", signature.Name)
	}

	worker.SetPostTaskHandler(postTaskHandler)
	worker.SetErrorHandler(errorHandler)
	worker.SetPreTaskHandler(preTaskHandler)

	if err := worker.Launch(); err != nil {
		log.ERROR.Fatalf("error while launching worker %v\n", err.Error())
		return err
	}

	return nil
}

func Start() error {
	if gin.IsDebugging() {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	} else {
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}

	// initiating a db connection
	psql.Initialize()
	db := psql.GetDbInstance()

	ctx := context.Background()

	// initializing sentry
	sentry.Initialize()

	// initilizing cache
	cache.Initialize(ctx)

	taskServer, err := StartServer(false)
	if err != nil {
		log.ERROR.Fatalf("Error starting task server: %v\n", err.Error())
		return err
	}

	cacheLayer := cache.InitializeCache()
	sentryLayer := sentry.InitializeSentry()
	orderDetailsLayer := orderDetailsService.InitializeOrderDetailsService()
	sampleServiceLayer := sampleService.InitializeSampleService()
	testSampleMappingServiceLayer := testSampleMappingService.InitializeTestSampleMappingService()
	attuneServiceLayer := attuneService.InitializeAttuneService()
	taskServiceLayer := taskService.InitializeTaskService()
	taskPathMappingServiceLayer := taskPathMappingService.InitializeTaskPathologistMappingService()
	patientServiceLayer := patientDetailService.InitializePatientDetailService()
	testServiceLayer := testDetailService.InitializeTestDetailService()
	investigationResultsServiceLayer := investigationResultsService.InitializeInvestigationResultService()
	rerunService := rerunService.InitializeRerunService()
	remarkServiceLayer := remarksService.InitializeRemarkService()
	reportGenerationServiceLayer := reportGenerationService.InitializeReportGenerationService()
	userServiceLayer := userService.InitializeUserService()
	receivingDeskServiceLayer := receivingDeskService.InitializeReceivingDeskService()
	abnormalityServiceLayer := abnormalityService.InitializeAbnormalityService()
	attachmentsServiceLayer := attachmentsService.InitializeAttachmentService()
	contactServiceLayer := contactService.InitializeContactService()
	etsServiceLayer := etsService.InitializeEtsService()
	cdsServiceLayer := cdsService.InitializeCdsService()
	pubsubServiceLayer := pubsubService.InitializePubsubService()
	attuneClientLayer := attuneClient.InitializeAttuneClient()
	cdsClientLayer := cdsClient.InitializeCdsClient()
	omsClientLayer := omsClient.InitializeOmsClient()
	reportRebrandingClientLayer := reportRebrandingClient.InitializeReportRebrandingClient()
	s3wrapperClientLayer := s3wrapperClient.InitializeS3wrapperClient()
	s3ClientLayer := s3Client.InitializeS3Client()
	partnerApiClientLayer := partnerApiClient.InitializePartnerApiClient()
	healthApiClientLayer := healthApiClient.InitializeHealthApiClient()
	snsClientLayer := snsClient.InitializeSnsClient()
	slackClientLayer := slackClient.InitializeSlackClient()

	commonTasksProcessor := commonTasks.CommonTaskProcessor{
		Db:                          db,
		Cache:                       cacheLayer,
		Sentry:                      sentryLayer,
		OrderDetailsService:         orderDetailsLayer,
		SampleService:               sampleServiceLayer,
		AttuneService:               attuneServiceLayer,
		TaskService:                 taskServiceLayer,
		TestDetailService:           testServiceLayer,
		InvestigationResultsService: investigationResultsServiceLayer,
		AttachmentsService:          attachmentsServiceLayer,
		RemarkService:               remarkServiceLayer,
		UserService:                 userServiceLayer,
		CdsService:                  cdsServiceLayer,
		ReportGenerationService:     reportGenerationServiceLayer,
		AttuneClient:                attuneClientLayer,
		S3Client:                    s3ClientLayer,
		S3wrapperClient:             s3wrapperClientLayer,
		ReportRebrandingClient:      reportRebrandingClientLayer,
	}

	wtService := workerTasks.WorkerTaskService{
		Db:                          db,
		Cache:                       cacheLayer,
		Sentry:                      sentryLayer,
		OrderDetailsService:         orderDetailsLayer,
		SampleService:               sampleServiceLayer,
		ReceivingDeskService:        receivingDeskServiceLayer,
		TestSampleMappingService:    testSampleMappingServiceLayer,
		AttuneService:               attuneServiceLayer,
		TaskService:                 taskServiceLayer,
		TaskPathMappingService:      taskPathMappingServiceLayer,
		PatientDetailService:        patientServiceLayer,
		TestDetailService:           testServiceLayer,
		InvestigationResultsService: investigationResultsServiceLayer,
		RerunService:                rerunService,
		RemarkService:               remarkServiceLayer,
		ReportGenerationService:     reportGenerationServiceLayer,
		UserService:                 userServiceLayer,
		AbnormalityService:          abnormalityServiceLayer,
		AttachmentsService:          attachmentsServiceLayer,
		ContactService:              contactServiceLayer,
		EtsService:                  etsServiceLayer,
		CdsService:                  cdsServiceLayer,
		PubsubService:               pubsubServiceLayer,
		AttuneClient:                attuneClientLayer,
		CdsClient:                   cdsClientLayer,
		OmsClient:                   omsClientLayer,
		ReportRebrandingClient:      reportRebrandingClientLayer,
		S3wrapperClient:             s3wrapperClientLayer,
		S3Client:                    s3ClientLayer,
		SnsClient:                   snsClientLayer,
		PartnerApiClient:            partnerApiClientLayer,
		HealthApiClient:             healthApiClientLayer,
		SlackClient:                 slackClientLayer,
		CommonTaskProcessor:         commonTasksProcessor,
	}
	wtpService := workerPeriodicTasks.WorkerPeriodicTaskService{
		Db:     db,
		Cache:  cacheLayer,
		Sentry: sentryLayer,

		EtsService:    etsServiceLayer,
		SampleService: sampleServiceLayer,
	}

	// Register tasks
	tasksToBeRegistered := map[string]interface{}{
		constants.ConsumerEvents:                   wtService.EventHandlerTask,
		constants.ReleaseReportTask:                wtService.ReleaseReportTask,
		constants.UpdateTaskPostSavingTask:         wtService.UpdateTaskPostSavingTask,
		constants.SampleRejectionTask:              wtService.SampleRejectionTask,
		constants.UpdateTaskPostReceivingTask:      wtService.UpdateTaskPostReceivingTask,
		constants.CreateUpdateTaskByOmsOrderIdTask: wtService.CreateUpdateTaskByOmsOrderIdTask,

		// Periodic tasks
		constants.StaleTasksPeriodicTask:                        wtpService.StaleTasksPeriodicTask,
		constants.EtsTatBreachedPeriodicTaskName:                wtpService.EtsTatBreachedPeriodicTask,
		constants.SlackAlertForReverseLogisticsPeriodicTaskName: wtpService.SlackAlertForDelayedReverseLogisticsPeriodicTask,
	}

	err = taskServer.RegisterTasks(tasksToBeRegistered)
	if err != nil {
		return err
	}

	return Worker(taskServer)
}
