package structures

import "time"

type ExternalInvestigateResultUpsertItem struct {
	ContactId                           uint       `json:"contact_id" binding:"required"`
	LoincCode                           string     `json:"loinc_code" binding:"required"`
	MasterInvestigationId               uint       `json:"master_investigation_id"`
	MasterInvestigationMethodMappingId  uint       `json:"master_investigation_method_mapping_id"`
	SystemExternalInvestigationResultId uint       `json:"system_external_investigation_result_id" binding:"required"`
	SystemExternalReportId              uint       `json:"system_external_report_id" binding:"required"`
	InvestigationName                   string     `json:"investigation_name" binding:"required"`
	InvestigationValue                  string     `json:"investigation_value" binding:"required"`
	Uom                                 string     `json:"uom"`
	ReferenceRangeText                  string     `json:"reference_range_text"`
	ReportedAt                          *time.Time `json:"reported_at" binding:"required"`
	CreatedBy                           uint       `json:"created_by"`
	UpdatedBy                           uint       `json:"updated_by"`
	IsAbnormal                          bool       `json:"is_abnormal"`
	LabName                             string     `json:"lab_name"`
	Abnormality                         string     `json:"abnormality"`
}

type UpsertExternalInvestigaitonResultsReqBody struct {
	Investigations []ExternalInvestigateResultUpsertItem `json:"investigations" binding:"required,gt=0,dive"`
}

type ExternalInvestigationResultsDbFilters struct {
	Limit                              uint
	Offset                             uint
	LoincCode                          string
	ContactId                          uint
	MasterInvestigationMethodMappingId uint
}
