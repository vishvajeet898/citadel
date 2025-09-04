package models

import (
	"time"
)

type TestDetail struct {
	BaseModel
	OmsOrderId           string     `gorm:"column:oms_order_id;not null" json:"oms_order_id"`
	TaskId               uint       `gorm:"column:task_id" json:"task_id"`
	OmsTestId            uint       `gorm:"column:oms_test_id;not null" json:"oms_test_id"`
	CentralOmsTestId     string     `gorm:"column:central_oms_test_id;not null;type:varchar(50)" json:"central_oms_test_id"`
	CityCode             string     `gorm:"column:city_code;not null;type:varchar(10)" json:"city_code"`
	TestName             string     `gorm:"column:test_name;not null;type:varchar(255)" json:"test_name"`
	LabId                uint       `gorm:"column:lab_id;not null" json:"lab_id"`
	ProcessingLabId      uint       `gorm:"column:processing_lab_id;not null" json:"processing_lab_id"`
	IsManualReportUpload bool       `gorm:"column:is_manual_report_upload;not null" json:"is_manual_report_upload"`
	IsDuplicate          bool       `gorm:"column:is_duplicate;not null" json:"is_duplicate"`
	LisCode              string     `gorm:"column:lis_code;not null;type:varchar(20)" json:"lis_code"`
	MasterTestId         uint       `gorm:"column:master_test_id;not null" json:"master_test_id"`
	MasterPackageId      uint       `gorm:"column:master_package_id" json:"master_package_id"`
	TestType             string     `gorm:"column:test_type;not null;type:varchar(10)" json:"test_type"`
	Department           string     `gorm:"column:department;type:varchar(100)" json:"department"`
	Status               string     `gorm:"column:status;not null;type:varchar(25)" json:"status"`
	OmsStatus            string     `gorm:"column:oms_status;type:varchar(25)" json:"oms_status"`
	DoctorTat            *time.Time `gorm:"column:doctor_tat;type:timestamp" json:"doctor_tat"`
	IsAutoApproved       bool       `gorm:"column:is_auto_approved;not null" json:"is_auto_approved"`
	ReportSentAt         *time.Time `gorm:"column:report_sent_at;type:timestamp" json:"report_sent_at"`
	ApprovalSource       string     `gorm:"column:approval_source;type:varchar(20)" json:"approval_source"`
	LabEta               *time.Time `gorm:"column:lab_eta;type:timestamp" json:"lab_eta"`
	LabTat               float32    `gorm:"column:lab_tat" json:"lab_tat"`
	ReportEta            *time.Time `gorm:"column:report_eta" json:"report_eta"`
	ReportStatus         string     `gorm:"column:report_status;not null;type:varchar(25)" json:"report_status"`
	CpEnabled            bool       `gorm:"column:cp_enabled;not null" json:"cp_enabled"`
}

func (TestDetail) TableName() string {
	return "test_details"
}
