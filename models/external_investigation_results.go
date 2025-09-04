package models

import "time"

type ExternalInvestigationResult struct {
	BaseModel
	ContactId                           uint       `gorm:"contact_id;index;not null" json:"contact_id"`
	MasterInvestigationId               uint       `gorm:"column:master_investigation_id;index;null" json:"master_investigation_id"`
	MasterInvestigationMethodMappingId  uint       `gorm:"column:master_investigation_method_mapping_id;index;null" json:"master_investigation_method_mapping_id"`
	LoincCode                           string     `gorm:"column:loinc_code;index;not null" json:"loinc_code"`
	SystemExternalInvestigationResultId uint       `gorm:"column:system_external_investigation_result_id;index;not null" json:"system_external_investigation_result_id"`
	InvestigationName                   string     `gorm:"column:investigation_name;type:varchar(100)" json:"investigation_name"`
	InvestigationValue                  string     `gorm:"column:investigation_value;type:varchar(100)" json:"investigation_value"`
	Uom                                 string     `gorm:"column:uom;type:varchar(100)" json:"uom"`
	ReferenceRangeText                  string     `gorm:"column:reference_range_text;type:varchar(500)" json:"reference_range_text"`
	ReportedAt                          *time.Time `gorm:"column:reported_at;index;not null" json:"reported_at"`
	IsAbnormal                          bool       `gorm:"column:is_abnormal;index;not null" json:"is_abnormal"`
	LabName                             string     `gorm:"column:lab_name;null" json:"lab_name"`
	SystemExternalReportId              uint       `gorm:"column:system_external_report_id;index;not null" json:"system_external_report_id"`
	Abnormality                         string     `gorm:"column:abnormality;index;null" json:"abnormality"`
}

func (ExternalInvestigationResult) TableName() string {
	return "external_investigation_results"
}
