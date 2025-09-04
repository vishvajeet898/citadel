package structures

import (
	"time"

	"github.com/google/uuid"
)

type OmsLisEvent struct {
	Request OmsRequestDetails `json:"request"`
	Order   OmsOrderDetails   `json:"order"`
	Visits  []OmsVisitDetails `json:"visits"`
	Patient OmsPatientDetails `json:"patient"`
	Tests   OmsTestDetails    `json:"tests"`
}

type OmsRequestDetails struct {
	Id             uint   `json:"id"`
	AlnumRequestId string `json:"alnum_request_id"`
}

type OmsOrderDetails struct {
	Id              uint   `json:"id"`
	AlnumOrderId    string `json:"alnum_order_id"`
	LabId           uint   `json:"lab_id"`
	CityCode        string `json:"city_code"`
	OrderType       string `json:"order_type"`
	ContainsPackage bool   `json:"contains_package"`
	DoctorName      string `json:"doctor_name"`
	DoctorNumber    string `json:"doctor_number"`
	DoctorNotes     string `json:"doctor_notes"`
	PartnerName     string `json:"partner_name"`
}

type OmsVisitDetails struct {
	Id string `json:"id"`
}

type OmsPatientDetails struct {
	Name        string     `json:"name"`
	AgeYears    float32    `json:"age_years"`
	AgeMonths   uint       `json:"age_months"`
	AgeDays     uint       `json:"age_days"`
	Gender      string     `json:"gender"`
	Number      string     `json:"number"`
	PatientId   string     `json:"patient_id"`
	ExpectedDob *time.Time `json:"expected_dob"`
	Dob         *time.Time `json:"dob"`
}

type OmsTestStruct struct {
	TestCode     string   `json:"test_code"`
	TestId       string   `json:"test_id"`
	MasterTestId uint     `json:"master_test_id"`
	TestName     string   `json:"test_name"`
	TestType     string   `json:"test_type"`
	Barcodes     []string `json:"barcodes"`
}

type OmsTestDetails struct {
	OrderInfo   map[string]LisTestUpdateInfo `json:"order_info"`
	TestDetails []OmsTestStruct              `json:"test_details"`
}

type LisTestUpdateInfoEvent struct {
	TestCode string               `json:"TestCode"`
	TestName string               `json:"TestName"`
	MetaData AttuneOrderInfoEvent `json:"MetaData"`
}

type AttuneOrderInfoEvent struct {
	TestID               string                            `json:"TestID"`
	TestCode             string                            `json:"TestCode"`
	TestType             string                            `json:"TestType"`
	TestName             string                            `json:"TestName"`
	TestValueType        string                            `json:"TestValueType"`
	TestValue            string                            `json:"TestValue"`
	DeviceActualValue    string                            `json:"DeviceActualValue"`
	QCLotNumber          string                            `json:"QCLotNumber"`
	QCStatus             string                            `json:"QCStatus"`
	QCWestGardWarning    string                            `json:"QCWestGardWarning"`
	QCValue              string                            `json:"QCValue"`
	QCFlag               string                            `json:"QCFlag"`
	UOMCode              string                            `json:"UOMCode"`
	MethodName           string                            `json:"MethodName"`
	DepartmentName       string                            `json:"DepartmentName"`
	ReferenceRange       string                            `json:"ReferenceRange"`
	IsAbnormal           string                            `json:"IsAbnormal"`
	ResultCapturedAt     string                            `json:"ResultCapturedAt"`
	ResultCapturedBy     int                               `json:"ResultCapturedBy"`
	AutoaApprovalFlag    string                            `json:"AutoaApprovalFlag"`
	ResultApprovedAt     string                            `json:"ResultApprovedAt"`
	ResultApprovedBy     int                               `json:"ResultApprovedBy"`
	TestStatus           string                            `json:"TestStatus"`
	MedicalRemarks       string                            `json:"MedicalRemarks"`
	TechnicalRemarks     string                            `json:"TechnicalRemarks"`
	RerunReason          string                            `json:"RerunReason"`
	RerunRemarks         string                            `json:"RerunRemarks"`
	RerunTime            string                            `json:"RerunTime"`
	DeviceID             string                            `json:"DeviceID"`
	IMDevice             string                            `json:"IMDevice"`
	IMDeviceFlag         string                            `json:"IMDeviceFlag"`
	UserID               int                               `json:"UserID"`
	SourceofApproval     string                            `json:"SourceofApproval"`
	TestDocumentInfo     []AttuneTestDocumentInfo          `json:"TestDocumentInfo"`
	OrderContentListInfo []AttuneOrderContentListInfoEvent `json:"OrderContentListInfo"`
	QcFlag               string                            `json:"QcFlag"`
	QcLotNumber          string                            `json:"QcLotNumber"`
	QcValue              string                            `json:"QcValue"`
	QcWestGardWarning    string                            `json:"QcWestGardWarning"`
	QcStatus             string                            `json:"QcStatus"`
}

type AttuneOrderContentListInfoEvent struct {
	TestCode          string                   `json:"TestCode"`
	TestType          string                   `json:"TestType"`
	TestName          string                   `json:"TestName"`
	TestValue         string                   `json:"TestValue"`
	UOMCode           string                   `json:"UOMCode"`
	MethodName        string                   `json:"MethodName"`
	ReferenceRange    string                   `json:"ReferenceRange"`
	DepartmentName    string                   `json:"DepartmentName"`
	TestStatus        string                   `json:"TestStatus"`
	MedicalRemarks    string                   `json:"MedicalRemarks"`
	TechnicalRemarks  string                   `json:"TechnicalRemarks"`
	DeviceID          string                   `json:"DeviceID"`
	IMDevice          string                   `json:"IMDevice"`
	IMDeviceFlag      string                   `json:"IMDeviceFlag"`
	CreatedAt         string                   `json:"CreatedAt"`
	UpdatedAt         string                   `json:"UpdatedAt"`
	UserID            int                      `json:"UserID"`
	UserName          string                   `json:"UserName"`
	Nonreportable     string                   `json:"Nonreportable"`
	ApproverName      string                   `json:"ApproverName"`
	ResultApprovedAt  string                   `json:"ResultApprovedAt"`
	ResultApprovedBy  int                      `json:"ResultApprovedBy"`
	ResultCapturedAt  string                   `json:"ResultCapturedAt"`
	ResultCapturedBy  int                      `json:"ResultCapturedBy"`
	DeviceActualValue string                   `json:"DeviceActualValue"`
	RerunReason       string                   `json:"RerunReason"`
	RerunRemarks      string                   `json:"RerunRemarks"`
	RerunTime         string                   `json:"RerunTime"`
	ParameterListInfo []interface{}            `json:"ParameterListInfo"`
	QcFlag            string                   `json:"QcFlag"`
	QcLotNumber       string                   `json:"QcLotNumber"`
	QcValue           string                   `json:"QcValue"`
	QcWestGardWarning string                   `json:"QcWestGardWarning"`
	QcStatus          string                   `json:"QcStatus"`
	TestDocumentInfo  []AttuneTestDocumentInfo `json:"TestDocumentInfo"`
}

type InitialTestDetails struct {
	DoctorTat      *time.Time `json:"doctor_tat"`
	IsAutoApproved bool       `json:"is_auto_approved"`
	TestStatus     string     `json:"test_status"`
	ApprovalSource string     `json:"approval_source"`
	IsCritical     bool       `json:"is_critical"`
}

type OmsAttachmentEvent struct {
	AlnumOrderId string                `json:"alnum_order_id"`
	CityCode     string                `json:"city_code"`
	Attachments  []OmsAttachmentStruct `json:"attachments"`
}

type OmsAttachmentStruct struct {
	LisCode       string `json:"lis_code"`
	FileUrl       string `json:"file_url"`
	FileReference string `json:"file_reference"`
	FileType      string `json:"file_type"`
	FileLabel     string `json:"file_label"`
	FileExtension string `json:"file_extension"`
}

type OmsTestDeleteEvent struct {
	AlnumOrderId string `json:"alnum_order_id"`
	CityCode     string `json:"city_code"`
	AlnumTestId  string `json:"alnum_test_id"`
}

type OmsManualReportUploadEvent struct {
	AlnumOrderId string   `json:"alnum_order_id"`
	CityCode     string   `json:"city_code"`
	AlnumTestIds []string `json:"alnum_test_ids"`
}

type OmsRerunEvent struct {
	AlnumOrderId string         `json:"alnum_order_id"`
	CityCode     string         `json:"city_code"`
	Tests        OmsTestDetails `json:"tests"`
}

type OmsOrderCreateUpdateEvent struct {
	CityCode                  string                  `json:"city_code"`
	Order                     OmsOrderModelDetails    `json:"order"`
	Request                   OmsRequestModelsDetails `json:"request"`
	Tests                     []OmsTestModelDetails   `json:"tests"`
	Tasks                     []OmsTaskModelDetails   `json:"tasks"`
	SynchronizeTasks          bool                    `json:"synchronize_tasks"`
	TaskTestsMapping          [][]TestsJsonStruct     `json:"task_tests_mapping"`
	IsPartialCancellationFlow bool                    `json:"is_partial_cancellation_flow"`
}

type OmsOrderModelDetails struct {
	Id               uint    `json:"id"`
	OriginalOrderId  uint    `json:"original_order_id"`
	AlnumOrderId     string  `json:"alnum_order_id"`
	TrfId            string  `json:"trf_id"`
	Status           uint    `json:"status"`
	TrfStatus        uint    `json:"trf_status"`
	PatientName      string  `json:"patient_name"`
	PatientNumber    string  `json:"patient_number"`
	PatientEmail     string  `json:"patient_email"`
	PatientAge       float32 `json:"patient_age"`
	PatientAgeMonths uint    `json:"patient_age_months"`
	PatientAgeDays   uint    `json:"patient_age_days"`
	PatientGender    string  `json:"patient_gender"`
	PatientId        string  `json:"patient_id"`
	PatientPk        uint    `json:"patient_pk"`
	OhTeamNotes      string  `json:"oh_team_notes"`
	PartnerId        uint    `json:"partner_id"`
	SystemDoctorId   uint    `json:"system_doctor_id"`
	OrderType        string  `json:"order_type"`
	CityCode         string  `json:"city_code"`
	CollectedOn      string  `json:"collected_on"`
	Notes            string  `json:"notes"`
	ReferredBy       string  `json:"referred_by"`
	SrfId            string  `json:"srf_id"`
}

type OmsRequestModelsDetails struct {
	Id                uint      `json:"id"`
	AlnumRequestId    string    `json:"alnum_request_id"`
	Uuid              uuid.UUID `json:"uuid"`
	Status            uint      `json:"status"`
	CustomerId        uint      `json:"customer_id"`
	DoctorId          uint      `json:"doctor_id"`
	SystemDoctorId    uint      `json:"system_doctor_id"`
	CampId            uint      `json:"camp_id"`
	ServicingLabId    uint      `json:"servicing_lab_id"`
	ServicingCityCode string    `json:"servicing_city_code"`
	CollectionType    uint      `json:"collection_type"`
	Source            string    `json:"source"`
	BulkOrderDetailId uint      `json:"bulk_order_detail_id"`
	CreatedOn         string    `json:"created_on"`
}

type OmsTestModelDetails struct {
	Id                  uint    `json:"id"`
	AlnumTestId         string  `json:"alnum_test_id"`
	TestName            string  `json:"test_name"`
	LabId               uint    `json:"lab_id"`
	Status              uint    `json:"status"`
	MasterTestId        uint    `json:"master_test_id" `
	MasterPackageId     uint    `json:"master_package_id"`
	TestCode            string  `json:"test_code"`
	TestType            string  `json:"test_type"`
	LabTat              float32 `json:"lab_tat"`
	TestingFrequency    string  `json:"testing_frequency"`
	ReportTat           string  `json:"report_tat"`
	LabEta              string  `json:"lab_eta"`
	ResultDataUrl       string  `json:"result_data_url"`
	CancellationReason  string  `json:"cancellation_reason"`
	TestAddUpdateType   string  `json:"test_add_update_type"`
	MappingSampleNumber uint    `json:"mapping_sample_number"`
}

type OmsTaskModelDetails struct {
	Id                        uint       `json:"id"`
	IsAdditionalTask          uint       `json:"is_additional_task,omitempty"`
	TaskMarkedRnrTime         *time.Time `json:"task_marked_rnr_time,omitempty"`
	TaskMarkedDescheduledTime *time.Time `json:"task_marked_descheduled_time,omitempty"`
	DeletedOn                 *time.Time `json:"deleted_on,omitempty"`
}

type TestsJsonStruct struct {
	Tests []TestExpandedJsonInDb `json:"tests"`
}

type TestExpandedJsonInDb struct {
	AlnumTestId string `json:"alnumTestId"`
}

type SampleCollectedEvent struct {
	RequestId   string         `json:"request_id"`
	TaskDetails OmsTaskDetails `json:"task_details"`
	TestIds     []string       `json:"test_ids"`
}

type OmsTaskDetails struct {
	OmsTaskType    uint   `json:"oms_task_type"`
	OmsTaskId      uint   `json:"oms_task_id"`
	CollectedAt    string `json:"collected_at"`
	CollectionType string `json:"collection_type"`
	IsB2c          bool   `json:"is_b2c"`
}

type SampleRecollectionEvent struct {
	RequestId                string                     `json:"request_id"`
	TaskId                   uint                       `json:"task_id"`
	RecollectionOrderDetails []RecollectionOrderDetails `json:"recollection_order_details"`
}

type RecollectionOrderDetails struct {
	OrderId string   `json:"order_id"`
	TestIds []string `json:"test_ids"`
}

type TestDetailsEvent struct {
	TestDetails []TestDetailsEventData `json:"test_details"`
}

type TestDetailsEventData struct {
	TestId    string  `json:"test_id"`
	LabTat    float32 `json:"lab_tat"`
	ReportTat string  `json:"report_tat"`
	LabEta    string  `json:"lab_eta"`
	Status    uint    `json:"status"`
}

type MarkSampleSnrEvent struct {
	AlnumTestIds []string `json:"alnum_test_ids"`
}
