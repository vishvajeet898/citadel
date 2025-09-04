package mapper

import (
	"github.com/Orange-Health/citadel/apps/attachments/structures"
	commonUtils "github.com/Orange-Health/citadel/common/utils"
	commonModels "github.com/Orange-Health/citadel/models"
)

func MapAttachment(attachment commonModels.Attachment) structures.Attachment {
	var attachmentStruct structures.Attachment

	attachmentStruct.ID = attachment.Id
	attachmentStruct.TaskId = attachment.TaskId
	attachmentStruct.InvestigationResultId = attachment.InvestigationResultId
	attachmentStruct.AttachmentUrl = attachment.AttachmentUrl
	attachmentStruct.ThumbnailUrl = attachment.ThumbnailUrl
	attachmentStruct.ThumbnailReference = attachment.ThumbnailReference
	attachmentStruct.AttachmentType = attachment.AttachmentType
	attachmentStruct.AttachmentLabel = attachment.AttachmentLabel
	attachmentStruct.Reference = attachment.Reference
	attachmentStruct.Extension = attachment.Extension
	attachmentStruct.CreatedAt = attachment.CreatedAt
	attachmentStruct.UpdatedAt = attachment.UpdatedAt
	attachmentStruct.DeletedAt = commonUtils.GetGoLangTimeFromGormDeletedAt(attachment.DeletedAt)
	attachmentStruct.CreatedBy = attachment.CreatedBy
	attachmentStruct.UpdatedBy = attachment.UpdatedBy
	attachmentStruct.DeletedBy = attachment.DeletedBy

	return attachmentStruct
}

func MapAttachments(attachments []commonModels.Attachment) []structures.Attachment {
	var attachmentStructs []structures.Attachment
	for _, attachment := range attachments {
		attachmentStructs = append(attachmentStructs, MapAttachment(attachment))
	}
	return attachmentStructs
}
