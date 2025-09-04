package structures

import (
	commonStructures "github.com/Orange-Health/citadel/common/structures"
)

// @swagger:response Attachment
type Attachment struct {
	commonStructures.BaseStruct
	// The task ID of the attachment.
	// example: 1
	TaskId uint `json:"task_id"`
	// The investigation result ID of the attachment.
	// example: 1
	InvestigationResultId uint `json:"investigation_result_id"`
	// The reference URL of the attachment.
	// example: "https://www.google.com"
	Reference string `json:"reference"`
	// The attachment URL of the attachment.
	// example: "citadel/demo.pdf"
	AttachmentUrl string `json:"attachment_url"`
	// The thumbnail URL of the attachment.
	// example: "https://www.google.com"
	ThumbnailUrl string `json:"thumbnail_url"`
	// The thumbnail reference of the attachment.
	// example: "citadel/resized-media/demo.pdf"
	ThumbnailReference string `json:"thumbnail_reference"`
	// The attachment type of the attachment.
	// example: "PDF"
	AttachmentType string `json:"attachment_type"`
	// The attachment label of the attachment.
	// example: "Report"
	AttachmentLabel string `json:"attachment_label"`
	// Is the attachment reportable in PDF.
	// example: true
	IsReportable bool `json:"is_reportable"`
	// The extension of the attachment.
	// example: "pdf"
	Extension string `json:"extension"`
}

type AddAttachmentRequest struct {
	AttachmentType  string `json:"attachment_type"`
	AttachmentLabel string `json:"attachment_label"`
	FileName        string `json:"file_name"`
}
