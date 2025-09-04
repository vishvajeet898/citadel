package models

type Attachment struct {
	BaseModel
	TaskId                uint   `gorm:"column:task_id;not null" json:"task_id"`
	Task                  Task   `gorm:"foreignKey:TaskId;references:Id"`
	InvestigationResultId uint   `gorm:"column:investigation_result_id" json:"investigation_result_id"`
	Reference             string `gorm:"column:reference;not null;type:varchar(100)" json:"reference"`
	AttachmentUrl         string `gorm:"column:attachment_url;not null;type:varchar(255)" json:"attachment_url"`
	ThumbnailUrl          string `gorm:"column:thumbnail_url;not null;type:varchar(255)" json:"thumbnail_url"`
	ThumbnailReference    string `gorm:"column:thumbnail_reference;not null;type:varchar(100)" json:"thumbnail_reference"`
	AttachmentType        string `gorm:"column:attachment_type;not null;type:varchar(50)" json:"attachment_type"`   // "report", "test_document", "visit_document", "prescription"
	AttachmentLabel       string `gorm:"column:attachment_label;not null;type:varchar(50)" json:"attachment_label"` // "lab_docs", "non_lab_docs", "outsource_docs"
	IsReportable          bool   `gorm:"column:is_reportable;" json:"is_reportable"`
	Extension             string `gorm:"column:extension;not null;type:varchar(10)" json:"extension"`
}

func (Attachment) TableName() string {
	return "attachments"
}
