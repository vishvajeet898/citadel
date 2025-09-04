package workerTasks

import (
	"gorm.io/gorm"

	"github.com/Orange-Health/citadel/adapters/cache"
	"github.com/Orange-Health/citadel/adapters/sentry"
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
	commonTasks "github.com/Orange-Health/citadel/common_tasks"
)

type WorkerTaskService struct {
	// Adapters
	Db     *gorm.DB
	Cache  cache.CacheLayer
	Sentry sentry.SentryLayer

	// Services
	OrderDetailsService         orderDetailsService.OrderDetailsServiceInterface
	SampleService               sampleService.SampleServiceInterface
	ReceivingDeskService        receivingDeskService.ReceivingDeskServiceInterface
	TestSampleMappingService    testSampleMappingService.TestSampleMappingServiceInterface
	AttuneService               attuneService.AttuneServiceInterface
	TaskService                 taskService.TaskServiceInterface
	TaskPathMappingService      taskPathMappingService.TaskPathologistMappingServiceInterface
	PatientDetailService        patientDetailService.PatientDetailServiceInterface
	TestDetailService           testDetailService.TestDetailServiceInterface
	InvestigationResultsService investigationResultsService.InvestigationResultServiceInterface
	RerunService                rerunService.RerunServiceInterface
	RemarkService               remarksService.RemarkServiceInterface
	ReportGenerationService     reportGenerationService.ReportGenerationInterface
	UserService                 userService.UserServiceInterface
	AbnormalityService          abnormalityService.AbnormalityServiceInterface
	AttachmentsService          attachmentsService.AttachmentServiceInterface
	ContactService              contactService.ContactServiceInterface
	EtsService                  etsService.EtsServiceInterface
	CdsService                  cdsService.CdsServiceInterface
	PubsubService               pubsubService.PubsubInterface

	// Clients
	AttuneClient           attuneClient.AttuneClientInterface
	CdsClient              cdsClient.CdsClientInterface
	OmsClient              omsClient.OmsClientInterface
	ReportRebrandingClient reportRebrandingClient.ReportRebrandingClientInterface
	S3wrapperClient        s3wrapperClient.S3wrapperInterface
	S3Client               s3Client.S3ClientInterface
	SnsClient              snsClient.SnsClientInterface
	PartnerApiClient       partnerApiClient.PartnerApiClientInterface
	HealthApiClient        healthApiClient.HealthApiClientInterface
	SlackClient            slackClient.SlackClientInterface

	// Common Task Processor
	CommonTaskProcessor commonTasks.CommonTaskProcessor
}
