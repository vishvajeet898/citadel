package models

import "time"

type SamplesAudit struct {
	Sample
	LogAction    string    `gorm:"column:log_action" json:"log_action,omitempty"`
	LogID        uint      `gorm:"column:log_id" json:"log_id,omitempty"`
	LogTimestamp time.Time `gorm:"column:log_timestamp" json:"log_timestamp,omitempty"`
}

func (SamplesAudit) TableName() string {
	return "samples_audit"
}
