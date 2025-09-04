package commonTasks

import (
	"gorm.io/gorm"

	"github.com/Orange-Health/citadel/adapters/cache"
	"github.com/Orange-Health/citadel/adapters/sentry"
	attachmentsService "github.com/Orange-Health/citadel/apps/attachments/service"
	attuneService "github.com/Orange-Health/citadel/apps/attune/service"
	cdsService "github.com/Orange-Health/citadel/apps/cds/service"
	investigationResultsService "github.com/Orange-Health/citadel/apps/investigation_results/service"
	orderDetailsService "github.com/Orange-Health/citadel/apps/order_details/service"
	pubsubService "github.com/Orange-Health/citadel/apps/pubsub/service"
	remarkService "github.com/Orange-Health/citadel/apps/remarks/service"
	reportGenerationService "github.com/Orange-Health/citadel/apps/report_generation/service"
	sampleService "github.com/Orange-Health/citadel/apps/samples/service"
	taskService "github.com/Orange-Health/citadel/apps/task/service"
	testDetailService "github.com/Orange-Health/citadel/apps/test_detail/service"
	userService "github.com/Orange-Health/citadel/apps/users/service"
	attuneClient "github.com/Orange-Health/citadel/clients/attune"
	reportRebrandingClient "github.com/Orange-Health/citadel/clients/report_rebranding"
	s3Client "github.com/Orange-Health/citadel/clients/s3"
	s3wrapperClient "github.com/Orange-Health/citadel/clients/s3wrapper"
)

type CommonTaskProcessor struct {
	Db     *gorm.DB
	Cache  cache.CacheLayer
	Sentry sentry.SentryLayer

	OrderDetailsService         orderDetailsService.OrderDetailsServiceInterface
	TaskService                 taskService.TaskServiceInterface
	TestDetailService           testDetailService.TestDetailServiceInterface
	SampleService               sampleService.SampleServiceInterface
	AttuneService               attuneService.AttuneServiceInterface
	InvestigationResultsService investigationResultsService.InvestigationResultServiceInterface
	AttachmentsService          attachmentsService.AttachmentServiceInterface
	RemarkService               remarkService.RemarkServiceInterface
	UserService                 userService.UserServiceInterface
	ReportGenerationService     reportGenerationService.ReportGenerationInterface
	CdsService                  cdsService.CdsServiceInterface
	PubsubService               pubsubService.PubsubInterface

	AttuneClient           attuneClient.AttuneClientInterface
	S3Client               s3Client.S3ClientInterface
	S3wrapperClient        s3wrapperClient.S3wrapperInterface
	ReportRebrandingClient reportRebrandingClient.ReportRebrandingClientInterface
}
