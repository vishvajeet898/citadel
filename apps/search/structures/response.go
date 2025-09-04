package structures

import "time"

type TaskListResponse struct {
	CriticalTasks    TaskResponseStruct `json:"critical"`
	WithheldTasks    TaskResponseStruct `json:"withheld"`
	CoAuthorizeTasks TaskResponseStruct `json:"co_authorize"`
	NormalTasks      TaskResponseStruct `json:"normal"`
	InProgressTasks  TaskResponseStruct `json:"in_progress"`
}

type AmendmentTaskListResponse struct {
	AmendmentTasks []TaskDetailsStruct `json:"amendment_tasks"`
	Count          uint                `json:"count"`
}

type TaskResponseStruct struct {
	Tasks    []TaskDetailsStruct `json:"tasks"`
	ShowMore bool                `json:"show_more"`
}

type TaskDetailsStruct struct {
	TaskId         uint                `json:"task_id"`
	OrderId        uint                `json:"order_id,omitempty"`
	OmsOrderId     string              `json:"oms_order_id,omitempty"`
	CityCode       string              `json:"city_code"`
	DoctorTat      *time.Time          `json:"doctor_tat,omitempty"`
	PickedBy       uint                `json:"picked_by,omitempty"`
	CoAuthorizedBy uint                `json:"co_authorized_by,omitempty"`
	ApprovedBy     []uint              `json:"approved_by,omitempty"`
	Status         string              `json:"status,omitempty"`
	PatientName    string              `json:"patient_name"`
	TaskContents   TaskContentStruct   `json:"task_contents,omitempty"`
	ProcessingLabs []LabDetailsStruct  `json:"processing_labs,omitempty"`
	TestDetails    []TestDetailsStruct `json:"test_details,omitempty"`
}

type TaskContentStruct struct {
	ContainsMorphle bool `json:"contains_morphle"`
	ContainsPackage bool `json:"contains_package"`
}

type LabDetailsStruct struct {
	Id   uint   `json:"id"`
	Name string `json:"name"`
}

type TestDetailsStruct struct {
	Name        string `json:"name"`
	Status      string `json:"status"`
	StatusLabel string `json:"status_label"`
}

type TaskDetailsDbStruct struct {
	TaskId          uint       `json:"task_id"`
	OrderId         uint       `json:"order_id"`
	OmsOrderId      string     `json:"oms_order_id"`
	CityCode        string     `json:"city_code"`
	DoctorTat       *time.Time `json:"doctor_tat"`
	PickedBy        uint       `json:"picked_by"`
	CoAuthorizedBy  uint       `json:"co_authorized_by"`
	Status          string     `json:"status"`
	PatientName     string     `json:"patient_name"`
	ContainsMorphle bool       `json:"contains_morphle"`
	ContainsPackage bool       `json:"contains_package"`
	LabId           uint       `json:"lab_id"`
}

type InfoScreenPatientDetailsResponse struct {
	Name        string     `json:"name"`
	ExpectedDob *time.Time `json:"expected_dob,omitempty"`
	Dob         *time.Time `json:"dob,omitempty"`
	Gender      string     `json:"gender"`
}

type InfoScreenTestDetailsResponse struct {
	ID                uint                   `json:"id"`
	AlnumTestId       string                 `json:"alnum_test_id"`
	Name              string                 `json:"name"`
	ProcessingLabId   uint                   `json:"processing_lab_id,omitempty"`
	ProcessingLabName string                 `json:"processing_lab_name,omitempty"`
	Status            string                 `json:"status"`
	Department        string                 `json:"department"`
	StatusLabel       string                 `json:"status_label"`
	CreatedAt         *time.Time             `json:"created_at"`
	DeletedAt         *time.Time             `json:"deleted_at"`
	Metadata          map[string]interface{} `json:"metadata"`
}

type InfoScreenOrderDetails struct {
	OrderID            string     `json:"order_id"`
	OrderStatus        string     `json:"order_status"`
	ServicingLabId     uint       `json:"servicing_lab_id"`
	ServicingCityCode  string     `json:"servicing_city_code"`
	PatientName        string     `json:"patient_name"`
	PatientAgeYears    uint       `json:"patient_age_years"`
	PatientAgeMonths   uint       `json:"patient_age_months"`
	PatientAgeDays     uint       `json:"patient_age_days"`
	PatientExpectedDob *time.Time `json:"patient_expected_dob"`
	PatientDob         *time.Time `json:"patient_dob"`
	PatientGender      string     `json:"patient_gender"`
	RequestID          string     `json:"request_id"`
	CampId             uint       `json:"camp_id"`
	CreatedAt          *time.Time `json:"created_at"`
	DeletedAt          *time.Time `json:"deleted_at"`
}

type InfoScreenTestDetails struct {
	ID                uint       `json:"id"`
	AlnumTestId       string     `json:"alnum_test_id"`
	TestName          string     `json:"test_name"`
	ProcessingLabId   uint       `json:"processing_lab_id"`
	Status            string     `json:"status"`
	Department        string     `json:"department"`
	OmsStatus         string     `json:"oms_status"`
	LabTat            string     `json:"lab_tat"`
	Barcode           string     `json:"barcode"`
	SampleNumber      uint       `json:"sample_number"`
	ReportReceivedAt  *time.Time `json:"report_received_at"`
	TechnicianSavedAt *time.Time `json:"technician_saved_at"`
	IsSampleRejected  bool       `json:"is_sample_rejected"`
	SampleStatus      string     `json:"sample_status"`
	AttuneStatus      string     `json:"attune_status"`
	CreatedAt         *time.Time `json:"created_at"`
	DeletedAt         *time.Time `json:"deleted_at"`
}

type InfoScreenBasicVisitDetails struct {
	VisitID          string     `json:"visit_id"`
	Barcode          string     `json:"barcode"`
	VialTypeID       uint       `json:"vial_type_id"`
	CurrentStatus    string     `json:"current_status"`
	SampleNumber     uint       `json:"sample_number"`
	SampleId         uint       `json:"sample_id"`
	ParentSampleId   uint       `json:"parent_sample_id"`
	LabId            uint       `json:"lab_id"`
	DestinationLabId uint       `json:"destination_lab_id"`
	CreatedAt        *time.Time `json:"created_at"`
	DeletedAt        *time.Time `json:"deleted_at"`
	CollectedAt      *time.Time `json:"collected_at"`
	ReceivedAt       *time.Time `json:"received_at"`
	TransferredAt    *time.Time `json:"transferred_at"`
	OutsourcedAt     *time.Time `json:"outsourced_at"`
	RejectedAt       *time.Time `json:"rejected_at"`
	NotReceivedAt    *time.Time `json:"not_received_at"`
	LisSyncAt        *time.Time `json:"lis_sync_at"`
	BarcodeScannedAt *time.Time `json:"barcode_scanned_at"`
}

type InfoScreenSearchResponse struct {
	Patient InfoScreenPatientDetailsResponse `json:"patient"`
	Request InfoScreenRequestDetailsResponse `json:"request"`
	Order   InfoScreenOrderDetailsResponse   `json:"order"`
	Visits  []InfoScreenVisitDetailsResponse `json:"visits"`
}

type InfoScreenRequestDetailsResponse struct {
	ID                string                 `json:"id"`
	ServicingLabId    uint                   `json:"servicing_lab_id"`
	ServicingCityCode string                 `json:"servicing_city_code"`
	ServicingLabName  string                 `json:"servicing_lab_name"`
	MetaData          map[string]interface{} `json:"metadata,omitempty"`
}

type InfoScreenOrderDetailsResponse struct {
	ID        string     `json:"id"`
	Status    string     `json:"status"`
	CreatedAt *time.Time `json:"created_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}

type InfoScreenVisitDetailsResponse struct {
	ID       string                            `json:"id"`
	Samples  []InfoScreenSampleDetailsResponse `json:"samples"`
	Metadata map[string]interface{}            `json:"metadata"`
}

type InfoScreenSampleDetailsResponse struct {
	SampleId         uint                            `json:"sample_id"`
	ParentSampleId   uint                            `json:"parent_sample_id,omitempty"`
	LabId            uint                            `json:"lab_id"`
	CurrentLabId     uint                            `json:"current_lab_id"`
	LabName          string                          `json:"lab_name"`
	CurrentLabName   string                          `json:"current_lab_name"`
	CollectedAt      *time.Time                      `json:"collected_at,omitempty"`
	ReceivedAt       *time.Time                      `json:"received_at,omitempty"`
	TransferredAt    *time.Time                      `json:"transferred_at,omitempty"`
	OutsourcedAt     *time.Time                      `json:"outsourced_at,omitempty"`
	RejectedAt       *time.Time                      `json:"rejected_at,omitempty"`
	NotReceivedAt    *time.Time                      `json:"not_received_at,omitempty"`
	LisSyncAt        *time.Time                      `json:"lis_sync_at,omitempty"`
	BarcodeScannedAt *time.Time                      `json:"barcode_scanned_at,omitempty"`
	SampleNumber     uint                            `json:"sample_number,omitempty"`
	VialTypeID       uint                            `json:"vial_type_id"`
	Barcode          string                          `json:"barcode,omitempty"`
	CurrentStatus    string                          `json:"current_status"`
	CreatedAt        *time.Time                      `json:"created_at"`
	DeletedAt        *time.Time                      `json:"deleted_at,omitempty"`
	Metadata         map[string]interface{}          `json:"metadata"`
	Tests            []InfoScreenTestDetailsResponse `json:"tests"`
}

type BarcodeDetailsResponse struct {
	Barcode          string                           `json:"barcode"`
	VialTypeID       uint                             `json:"vial_type_id"`
	Patient          InfoScreenPatientDetailsResponse `json:"patient"`
	Tests            []InfoScreenTestDetailsResponse  `json:"tests"`
	CreatedAt        *time.Time                       `json:"created_at"`
	DeletedAt        *time.Time                       `json:"deleted_at"`
	VisitID          string                           `json:"visit_id"`
	OrderID          string                           `json:"order_id"`
	SampleReceivedAt *time.Time                       `json:"sample_received_at"`
	Metadata         map[string]interface{}           `json:"metadata"`
}
