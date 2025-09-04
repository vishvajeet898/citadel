package controller

import (
	"github.com/Orange-Health/citadel/apps/attachments/service"
)

type Attachment struct {
	AttachmentService service.AttachmentServiceInterface
}

func InitAttachmentController() *Attachment {
	return &Attachment{
		AttachmentService: service.InitializeAttachmentService(),
	}
}
