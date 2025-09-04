package structures

import (
	"time"
)

// @swagger:response InvestigationResult
type InvestigationResult struct {
	// The investigation result ID.
	// example: 1
	Id uint `json:"id"`
	// The test detail ID of the investigation result.
	// example: 1
	TestDetailId uint `json:"test_detail_id"`
	// The test detail status of the investigation result.
	// example: "approved"
	TestDetailsStatus string `json:"test_details_status"`
	// The investigation ID of the investigation result.
	// example: 1
	MasterInvestigationId uint `json:"master_investigation_id,omitempty"`
	// The investigation name of the investigation result.
	// example: "Hemoglobin"
	InvestigationName string `json:"investigation_name"`
	// The investigation value of the investigation result.
	// example: "12.5"
	InvestigationValue string `json:"investigation_value,omitempty"`
	// The investigation device value of the investigation result.
	// example: "12.5"
	DeviceValue string `json:"device_value,omitempty"`
	// The investigation result representation type of the investigation result.
	// example: "biopattern"
	ResultRepresentationType string `json:"result_representation_type,omitempty"`
	// The investigation Department of the investigation result.
	// example: "Hematology"
	Department string `json:"department,omitempty"`
	// The investigation unit of measurement of the investigation result.
	// example: "g/dL"
	Uom string `json:"uom,omitempty"`
	// The investigation reference range of the investigation result.
	// example: "12.0 - 16.0"
	ReferenceRangeText string `json:"reference_range_text,omitempty"`
	// The investigation method performed.
	// example: "Manual"
	Method string `json:"method,omitempty"`
	// The investigation method type.
	// example: "Manual"
	MethodType string `json:"method_type,omitempty"`
	// The investigation Lis-Code.
	// example: "HGB0001"
	LisCode string `json:"lis_code"`
	// The investigation abnormality.
	// example: "N"
	Abnormality string `json:"abnormality,omitempty"`
	// The investigation Approved By.
	// example: 1
	ApprovedBy uint `json:"approved_by,omitempty"`
	// The investigation Approved At.
	// example: "2021-07-01T12:00:00Z"
	ApprovedAt string `json:"approved_at,omitempty"`
	// The investigation Entered Technician.
	// example: 1
	EnteredBy uint `json:"entered_by,omitempty"`
	// The investigation result Entered At.
	// example: "2021-07-01T12:00:00Z"
	EnteredAt string `json:"entered_at,omitempty"`
	// The investigation Status.
	// example: "Approved"
	Status string `json:"status"`
	// The investigation Reportable in PDF.
	// example: true
	IsNonReportable bool `json:"is_non_reportable,omitempty"`
	// If the investigation was auto approved.
	// example: true
	IsAutoApproved bool `json:"is_auto_approved"`
	// The investigation Source.
	// example: "OH Lab BLR"
	Source string `json:"source,omitempty"`
	// The investigation Is Critical.
	// example: true
	IsCritical bool `json:"is_critical"`
	// The investigation Master Test ID.
	// example: 1
	MasterTestId uint `json:"master_test_id"`
	// Barcodes of the investigation.
	// example: "1234567890"
	Barcodes string `json:"barcodes,omitempty"`
	// The investigation Investigation Data.
	// example: {
	// 		"id": 1,
	// 		"data": "This is a investigation value in text",
	// 		"data_type": "investigation_value_text",
	// 		"created_at": "2021-07-01T12:00:00Z",
	// 		"updated_at": "2021-07-01T12:00:00Z",
	// 		"deleted_at": "2021-07-01T12:00:00Z",
	// 		"created_by": 1,
	// 		"updated_by": 1,
	// 		"deleted_by": 1
	// 	}
	InvestigationData InvestigationData `json:"investigation_data,omitempty"`

	AutoApprovalFailureReason string `json:"auto_approval_failure_reason,omitempty"`

	ProcessingLabId uint `json:"processing_lab_id,omitempty"`
	// The QC fields carries details about the Quality Check run and it's status.
	QcFlag            string `json:"qc_flag,omitempty"`
	QcLotNumber       string `json:"qc_lot_number,omitempty"`
	QcValue           string `json:"qc_value,omitempty"`
	QcWestGardWarning string `json:"qc_west_gard_warning,omitempty"`
	QcStatus          string `json:"qc_status,omitempty"`
}

// @swagger:response Remark
type Remark struct {
	// The remark ID.
	// example: 1
	Id uint `json:"id"`
	// The investigation result ID of the remark.
	// example: 1
	InvestigationResultId uint `json:"investigation_result_id"`
	// The description of the remark.
	// example: "This is a remark"
	Description string `json:"description"`
	// The remark type of the remark.
	// example: "RERUN"
	RemarkType string `json:"remark_type"`
	// The remark by of the remark.
	// example: 12
	RemarkBy uint `json:"remark_by"`
}

// @swagger:response RerunDetails
type RerunDetails struct {
	// The Rerun Investigation ID
	// example: 1
	Id uint `json:"id"`
	// The Test Detail ID
	// example: 1
	TestDetailId uint `json:"test_detail_id"`
	// The Master Investigation ID
	// example: 1
	MasterInvestigationId uint `json:"master_investigation_id"`
	// The Investigation Value
	// example: "12.5"
	InvestigationValue string `json:"investigation_value,omitempty"`
	// The Investigation Device Value
	// example: "12.5"
	DeviceValue string `json:"device_value,omitempty"`
	// User ID of the User who triggered the Rerun
	// example: 1
	RerunTriggeredBy uint `json:"rerun_triggered_by"`
	// Time at which the Rerun was Triggered
	// example: "2021-07-01T12:00:00Z"
	RerunTriggeredAt *time.Time `json:"rerun_triggered_at"`
	// Reason for Rerun
	// example: "Invalid Value"
	RerunReason string `json:"rerun_reason"`
	// Remarks for Rerun
	// example: "Value was invalid"
	RerunRemarks string `json:"rerun_remarks"`
}

type InvestigationResultDetails struct {
	InvestigationResult InvestigationResult `json:"investigation_result"`

	RerunDetails []RerunDetails `json:"rerun_details,omitempty"`
	// The investigation remarks.
	// example: [
	// 	{
	// 		"id": 1,
	// 		"remark": "This is a remark",
	// 		"remark_type": "medical_remark",
	// 		"created_at": "2021-07-01T12:00:00Z",
	// 		"updated_at": "2021-07-01T12:00:00Z",
	// 		"deleted_at": "2021-07-01T12:00:00Z",
	// 		"created_by": 1,
	// 		"updated_by": 1,
	// 		"deleted_by": 1
	// 	}
	// ]
	Remarks []Remark `json:"remarks,omitempty"`
}

// @swagger:response InvestigationData
type InvestigationData struct {
	// The investigation data ID.
	// example: 1
	Id uint `json:"id,omitempty"`
	// Data contains the investigation data in text format
	// swagger:strfmt binary
	Data string `json:"data,omitempty"`
	// The investigation data type.
	// example: "medical_remark"
	DataType string `json:"data_type,omitempty"`
	// InvestigationResultId is the ID of the investigation result.
	// example: 1
	InvestigationResultId uint `json:"investigation_result_id,omitempty"`
}

type ModifyValueResponse struct {
	// Value of the investigation result.
	// example: "10"
	Value string `json:"value"`
	// Abnormality of the investigation result.
	// example: "upper_abnormal"
	Abnormality string `json:"abnormality"`
}

type InvestigationResultDbResponse struct {
	Id                                 uint       `json:"id"`
	TestDetailsId                      uint       `json:"test_details_id"`
	MasterInvestigationId              uint       `json:"master_investigation_id"`
	MasterInvestigationMethodMappingId uint       `json:"master_investigation_method_mapping_id"`
	MasterTestId                       uint       `json:"master_test_id"`
	TestDetailsStatus                  string     `json:"test_details_status"`
	InvestigationName                  string     `json:"investigation_name"`
	InvestigationValue                 string     `json:"investigation_value"`
	DeviceValue                        string     `json:"device_value"`
	ResultRepresentationType           string     `json:"result_representation_type"`
	Department                         string     `json:"department"`
	Uom                                string     `json:"uom"`
	Method                             string     `json:"method"`
	MethodType                         string     `json:"method_type"`
	ReferenceRangeText                 string     `json:"reference_range_text"`
	LisCode                            string     `json:"lis_code"`
	Abnormality                        string     `json:"abnormality"`
	IsAbnormal                         bool       `json:"is_abnormal"`
	ApprovedBy                         uint       `json:"approved_by"`
	ApprovedAt                         *time.Time `json:"approved_at"`
	EnteredBy                          uint       `json:"entered_by"`
	EnteredAt                          *time.Time `json:"entered_at"`
	InvestigationStatus                string     `json:"investigation_status"`
	IsAutoApproved                     bool       `json:"is_auto_approved"`
	IsCritical                         bool       `json:"is_critical"`
	Barcodes                           string     `json:"barcodes"`
	ProcessingLabId                    uint       `json:"processing_lab_id"`
	AutoApprovalFailureReason          string     `json:"auto_approval_failure_reason" `
	QcFlag                             string     `json:"qc_flag,omitempty"`
	QcLotNumber                        string     `json:"qc_lot_number,omitempty"`
	QcValue                            string     `json:"qc_value,omitempty"`
	QcWestGardWarning                  string     `json:"qc_west_gard_warning,omitempty"`
	QcStatus                           string     `json:"qc_status,omitempty"`
	CpEnabled                          bool       `json:"cp_enabled,omitempty"`
}

type TestDetailsDbResponse struct {
	Id           uint   `json:"id"`
	MasterTestId uint   `json:"master_test_id"`
	Status       string `json:"status"`
	LabId        uint   `json:"lab_id"`
	CpEnabled    bool   `json:"cp_enabled"`
	Barcodes     string `json:"barcodes"`
}
