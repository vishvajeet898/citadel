package service

import (
	"github.com/Orange-Health/citadel/adapters/cache"
	"github.com/Orange-Health/citadel/adapters/sentry"
	"github.com/Orange-Health/citadel/apps/attachments/dao"
	reportRebrandingClient "github.com/Orange-Health/citadel/clients/report_rebranding"
	s3wrapperClient "github.com/Orange-Health/citadel/clients/s3wrapper"
)

type AttachmentService struct {
	AttachmentDao          dao.DataLayer
	Cache                  cache.CacheLayer
	Sentry                 sentry.SentryLayer
	S3wrapperClient        s3wrapperClient.S3wrapperInterface
	ReportRebrandingClient reportRebrandingClient.ReportRebrandingClientInterface
}

func InitializeAttachmentService() AttachmentServiceInterface {
	return &AttachmentService{
		AttachmentDao:          dao.InitializeAttachmentDao(),
		Cache:                  cache.InitializeCache(),
		Sentry:                 sentry.InitializeSentry(),
		S3wrapperClient:        s3wrapperClient.InitializeS3wrapperClient(),
		ReportRebrandingClient: reportRebrandingClient.InitializeReportRebrandingClient(),
	}
}
