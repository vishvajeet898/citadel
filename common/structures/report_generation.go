package structures

import "time"

type ReportGenerationEvent struct {
	Order        OrderEvent         `json:"order"`
	Visits       []VisitEvent       `json:"visits"`
	IndexDetails ReportIndexDetails `json:"index_details"`
}

type OrderEvent struct {
	ID             uint        `json:"id"`
	RequestId      uint        `json:"request_id"`
	CreatedOn      string      `json:"created_on"`
	PartnerID      uint        `json:"partner_id"`
	PatientDetails interface{} `json:"patient_details"`
	ReferredBy     string      `json:"referred_by"`
	CityCode       string      `json:"city_code"`
}

type VisitEvent struct {
	Id    string      `json:"id"`
	Date  string      `json:"date"`
	Tests []TestEvent `json:"tests"`
}

type TestEvent struct {
	Id                uint                 `json:"id"`
	MasterTestId      uint                 `json:"master_test_id"`
	Name              string               `json:"name"`
	Status            string               `json:"status"`
	BlockTemplateId   uint                 `json:"block_template_id"`
	IsPanel           bool                 `json:"is_panel"`
	DepartmerntName   string               `json:"department_name"`
	SampleType        string               `json:"sample_type"`
	SampleCollectedAt string               `json:"sample_collected_at"`
	SampleReceivedAt  string               `json:"sample_received_at"`
	ApprovedAt        string               `json:"approved_at"`
	ApprovedBy        string               `json:"approved_by"`
	Investigations    []InvestigationEvent `json:"investigations"`
	ReportReleasedAt  string               `json:"report_released_at"`
	ReportReleasedBy  int                  `json:"report_released_by"`
}

type InvestigationEvent struct {
	Id                       uint       `json:"id"`
	MasterInvestigationId    uint       `json:"master_investigation_id"`
	ResultRepresentationType string     `json:"result_representation_type"`
	TestId                   uint       `json:"test_id"`
	OrderId                  uint       `json:"order_id"`
	Name                     string     `json:"name"`
	DepartmentName           string     `json:"department_name"`
	Method                   string     `json:"method"` // Check if this is required
	Value                    string     `json:"value"`
	Data                     string     `json:"data,omitempty"`
	Unit                     string     `json:"unit"`
	ReferenceRangeText       string     `json:"reference_range_text"`
	Status                   string     `json:"status"`
	ApprovedAt               *time.Time `json:"approved_at"`
	ApprovedBy               string     `json:"approved_by"`
	ResultEnteredAt          *time.Time `json:"result_entered_at"`
	ResultEnteredBy          string     `json:"result_entered_by"`
	InvestigationMethod      string     `json:"investigation_method"`
	InterpretationNotes      string     `json:"interpretation_notes"`
	MedicalRemarks           string     `json:"medical_remarks"`
	IsAbnormal               bool       `json:"is_abnormal"`
	IsNablApproved           bool       `json:"is_nabl_approved"`
	Source                   string     `json:"source"`
	CreatedAt                string     `json:"created_at"`
	UpdatedAt                string     `json:"updated_at"`
}

type ReportIndexDetails struct {
	IsPackageOrder bool                `json:"is_package_order"`
	PackageIds     []uint              `json:"package_ids"`
	TestDetails    []ReportTestDetails `json:"test_details"`
}

type ReportTestDetails struct {
	MasterTestId   uint   `json:"test_id"`
	TestStatus     uint   `json:"test_status"`
	TestEta        string `json:"test_eta"`
	TestHoldReason string `json:"test_hold_reason"`
	IsNew          bool   `json:"is_new"`
	MappingId      uint   `json:"mapping_id"`
}
