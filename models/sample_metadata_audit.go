package models

import "time"

type SampleMetadataAudit struct {
	SampleMetadata
	LogAction    string    `gorm:"column:log_action" json:"log_action,omitempty"`
	LogID        uint      `gorm:"column:log_id" json:"log_id,omitempty"`
	LogTimestamp time.Time `gorm:"column:log_timestamp" json:"log_timestamp,omitempty"`
}

func (SampleMetadataAudit) TableName() string {
	return "sample_metadata_audit"
}
