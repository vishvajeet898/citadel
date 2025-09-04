package structures

import (
	"time"

	"github.com/google/uuid"
)

type ReportGenerationEvent struct {
	Order             OrderEvent         `json:"order"`
	Visits            []VisitEvent       `json:"visits"`
	IndexDetails      ReportIndexDetails `json:"index_details"`
	ServicingCityCode string             `json:"servicing_city_code"`
	IsDummyReport     bool               `json:"is_dummy_report"`
}

type OrderEvent struct {
	Id             string      `json:"id"`
	AlnumOrderId   string      `json:"alnum_order_id"`
	Token          uuid.UUID   `json:"token"`
	RequestId      string      `json:"request_id"`
	CreatedOn      string      `json:"created_on"`
	PartnerId      uint        `json:"partner_id"`
	PatientDetails interface{} `json:"patient_details"`
	ReferredBy     string      `json:"referred_by"`
	ServicingLabId uint        `json:"servicing_lab_id"`
	CityCode       string      `json:"city_code"`
}

type VisitEvent struct {
	Id    string      `json:"id"`
	Date  string      `json:"date"`
	LabId uint        `json:"lab_id"`
	Tests []TestEvent `json:"tests"`
}

type TestEvent struct {
	Id                string               `json:"id"`
	MasterTestId      uint                 `json:"master_test_id"`
	Name              string               `json:"name"`
	Status            string               `json:"status"`
	BlockTemplateId   uint                 `json:"block_template_id"`
	IsPanel           bool                 `json:"is_panel"`
	LabId             uint                 `json:"lab_id"`
	DepartmerntName   string               `json:"department_name"`
	SampleType        string               `json:"sample_type"`
	SampleId          string               `json:"sample_id"`
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
	TestId                   string     `json:"test_id"`
	OrderId                  string     `json:"order_id"`
	Name                     string     `json:"name"`
	LisCode                  string     `json:"lis_code"`
	Abnormality              string     `json:"abnormality"`
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
	TestDocument             []string   `json:"test_document,omitempty" gorm:"-"`
}

type VisitTestsStruct struct {
	VisitId          string `json:"visit_id"`
	Barcode          string `json:"barcode"`
	ReceivedAt       string `json:"received_at"`
	CollectedAt      string `json:"collected_at"`
	SampleNumber     string `json:"sample_number"`
	SampleName       string `json:"sample_name"`
	DepartmentName   string `json:"department_name"`
	Id               uint   `json:"id"`
	CentralOmsTestId string `json:"central_oms_test_id"`
	OmsOrderID       string `json:"oms_order_id"`
	TestName         string `json:"test_name"`
	MasterTestId     uint   `json:"master_test_id"`
	LisCode          string `json:"lis_code"`
	LabId            uint   `json:"lab_id"`
	ProcessingLabId  uint   `json:"processing_lab_id"`
	TestType         string `json:"test_type"`
	VialTypeId       uint   `json:"vial_type_id"`
}

type PatientDetailsEvent struct {
	Id        string  `json:"id"`
	Name      string  `json:"name"`
	Number    string  `json:"number"`
	Gender    string  `json:"gender"`
	AgeYears  float32 `json:"age_years"`
	AgeMonths uint    `json:"age_months"`
	AgeDays   uint    `json:"age_days"`
}

type ReportIndexDetails struct {
	IsCollectionTypeInClinic bool                `json:"is_collection_type_in_clinic"`
	IsPackageOrder           bool                `json:"is_package_order"`
	PackageIds               []uint              `json:"package_ids"`
	TestDetails              []ReportTestDetails `json:"test_details"`
}

type ReportTestDetails struct {
	MasterTestId   uint   `json:"test_id"`
	TestStatus     uint   `json:"test_status"`
	TestEta        string `json:"test_eta"`
	TestHoldReason string `json:"test_hold_reason"`
	IsNew          bool   `json:"is_new"`
	MappingId      uint   `json:"mapping_id"`
}
