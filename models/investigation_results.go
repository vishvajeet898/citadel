package models

import (
	"time"
)

type InvestigationResult struct {
	BaseModel
	TestDetailsId                      uint       `gorm:"column:test_details_id" json:"test_details_id"`
	MasterInvestigationId              uint       `gorm:"column:master_investigation_id" json:"master_investigation_id"`
	MasterInvestigationMethodMappingId uint       `gorm:"column:master_investigation_method_mapping_id" json:"master_investigation_method_mapping_id"`
	InvestigationName                  string     `gorm:"column:investigation_name;type:varchar(100)" json:"investigation_name"`
	InvestigationValue                 string     `gorm:"column:investigation_value;type:varchar(100)" json:"investigation_value"`
	DeviceValue                        string     `gorm:"column:device_value;type:varchar(255)" json:"device_value"`
	ResultRepresentationType           string     `gorm:"column:result_representation_type;type:varchar(50)" json:"result_representation_type"`
	Department                         string     `gorm:"column:department;type:varchar(100)" json:"department"`
	Uom                                string     `gorm:"column:uom;type:varchar(100)" json:"uom"`
	Method                             string     `gorm:"column:method;type:varchar(100)" json:"method"`
	MethodType                         string     `gorm:"column:method_type;type:varchar(50)" json:"method_type"`
	ReferenceRangeText                 string     `gorm:"column:reference_range_text;type:varchar(500)" json:"reference_range_text"`
	LisCode                            string     `gorm:"column:lis_code;type:varchar(100)" json:"lis_code"`
	Abnormality                        string     `gorm:"column:abnormality;type:varchar(10)" json:"abnormality"`
	IsAbnormal                         bool       `gorm:"column:is_abnormal" json:"is_abnormal"`
	ApprovedBy                         uint       `gorm:"column:approved_by" json:"approved_by"`
	ApprovedAt                         *time.Time `gorm:"column:approved_at" json:"approved_at"`
	EnteredBy                          uint       `gorm:"column:entered_by" json:"entered_by"`
	EnteredAt                          *time.Time `gorm:"column:entered_at" json:"entered_at"`
	InvestigationStatus                string     `gorm:"column:investigation_status;type:varchar(50)" json:"investigation_status"`
	IsAutoApproved                     bool       `gorm:"column:is_auto_approved" json:"is_auto_approved"`
	IsNonReportable                    bool       `gorm:"column:is_non_reportable" json:"is_non_reportable"`
	AutoVerified                       bool       `gorm:"column:auto_verified" json:"auto_verified"`
	IsNablApproved                     bool       `gorm:"column:is_nabl_approved" json:"is_nabl_approved"`
	Source                             string     `gorm:"column:source;type:varchar(50)" json:"source"`
	IsCritical                         bool       `gorm:"column:is_critical" json:"is_critical"`
	ApprovalSource                     string     `gorm:"column:approval_source;type:varchar(20)" json:"approval_source"`
	AutoApprovalFailureReason          string     `gorm:"column:auto_approval_failure_reason;type:varchar(255)" json:"auto_approval_failure_reason"`
}

func (InvestigationResult) TableName() string {
	return "investigation_results"
}
