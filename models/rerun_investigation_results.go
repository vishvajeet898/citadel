package models

import (
	"time"
)

type RerunInvestigationResult struct {
	BaseModel
	TestDetailsId            uint       `gorm:"column:test_details_id;not null" json:"test_details_id"`
	TestDetail               TestDetail `gorm:"foreignKey:TestDetailsId;references:Id"`
	MasterInvestigationId    uint       `gorm:"column:master_investigation_id;not null" json:"master_investigation_id"`
	InvestigationName        string     `gorm:"column:investigation_name;not null;type:varchar(255)" json:"investigation_name"`
	InvestigationValue       string     `gorm:"column:investigation_value;not null;type:varchar(255)" json:"investigation_value"`
	DeviceValue              string     `gorm:"column:device_value;not null;type:varchar(255)" json:"device_value"`
	ResultRepresentationType string     `gorm:"column:result_representation_type;not null;type:varchar(255)" json:"result_representation_type"`
	LisCode                  string     `gorm:"column:lis_code;not null;type:varchar(255)" json:"lis_code"`
	RerunTriggeredBy         uint       `gorm:"column:rerun_triggered_by;not null" json:"rerun_triggered_by"`
	RerunTriggeredByUser     User       `gorm:"foreignKey:RerunTriggeredBy;references:Id"`
	RerunTriggeredAt         *time.Time `gorm:"column:rerun_triggered_at;type:timestamp" json:"rerun_triggered_at"`
	RerunReason              string     `gorm:"column:rerun_reason;not null;type:varchar(255)" json:"rerun_reason"`
	RerunRemarks             string     `gorm:"column:rerun_remarks;not null;type:text" json:"rerun_remarks"`
	EnteredBy                uint       `gorm:"column:entered_by;not null" json:"entered_by"`
	EnteredAt                *time.Time `gorm:"column:entered_at;type:timestamp" json:"entered_at"`
}

func (RerunInvestigationResult) TableName() string {
	return "rerun_investigation_results"
}
