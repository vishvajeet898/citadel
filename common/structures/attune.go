package structures

import (
	"time"
)

type AttuneAuthResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int64  `json:"expires_in"`
}

type AttuneGetSyncDataResponse struct {
	Status   string              `json:"status"`
	Response AttuneOrderResponse `json:"response"`
}

type AttuneOrderResponse struct {
	OrderId           string                 `json:"OrderId"`
	OrgCode           string                 `json:"OrgCode"`
	OverAllStatus     string                 `json:"OverAllStatus"`
	PatientInfo       AttunePatientInfo      `json:"PatientInfo"`
	PatientVisitInfo  AttunePatientVisitInfo `json:"PatientVisitInfo"`
	CorporatePatient  string                 `json:"CorporatePatient"`
	VisitDocumentinfo []VisitDocumentInfo    `json:"VisitDocumentinfo"`
	OrderInfo         []AttuneOrderInfo      `json:"OrderInfo"`
	ResultAsPdf       string                 `json:"ResultAsPdf"`
	OutsourceAsPdf    string                 `json:"OutsourceAsPdf"`
	ReportPdfFormat   string                 `json:"ReportPDFFormat"`
}

type LisOrderInfo struct {
	OrderId           string              `json:"OrderId"`
	OrgCode           string              `json:"OrgCode"`
	OverAllStatus     string              `json:"OverAllStatus"`
	VisitDocumentinfo []VisitDocumentInfo `json:"VisitDocumentinfo"`
	OrderInfo         []AttuneOrderInfo   `json:"OrderInfo"`
	ResultAsPdf       string              `json:"ResultAsPdf"`
	ReportPdfFormat   string              `json:"ReportPDFFormat"`
}

type AttunePatientInfo struct {
	Salutation        string              `json:"Salutation"`
	PatientId         string              `json:"PatientID"`
	PatientNumber     string              `json:"PatientNumber"`
	SalutationCode    string              `json:"SalutationCode"`
	FirstName         string              `json:"FirstName"`
	MiddleName        string              `json:"MiddleName"`
	LastName          string              `json:"LastName"`
	Gender            string              `json:"Gender"`
	Age               string              `json:"Age"`
	Dob               string              `json:"DOB"`
	MobileNumber      string              `json:"MobileNumber"`
	EmailId           string              `json:"EmailID"`
	UrnType           string              `json:"URNType"`
	UrnNumber         string              `json:"URNNumber"`
	MaritalStatus     string              `json:"MaritalStatus"`
	TelephoneNumber   string              `json:"TelephoneNumber"`
	ExternalPatientNo int                 `json:"ExternalPatientNo"`
	HealthHubId       int                 `json:"HealthHubId"`
	EmployeeId        int                 `json:"EmployeeId"`
	DocumentOf        string              `json:"DocumentOf"`
	AddressDetails    []AttuneAddressInfo `json:"AddressDetails"`
	Name              string              `json:"Name"`
	Email             string              `json:"Email"`
}

type AttuneAddressInfo struct {
	Address     string `json:"Address"`
	AddressType string `json:"AddressType"`
	City        string `json:"City"`
	StateId     string `json:"StateID"`
	CountryId   string `json:"CountryID"`
	Suburb      string `json:"Suburb"`
	State       string `json:"State"`
	Country     string `json:"Country"`
	PostalCode  int    `json:"PostalCode"`
	Location    int    `json:"Location"`
	PassportNo  string `json:"PassportNo"`
}

type AttunePatientVisitInfo struct {
	CompanyId                  int                 `json:"CompanyID"`
	PatientVisitId             string              `json:"PatientVisitId"`
	ExternalVisitNumber        string              `json:"ExternalVisitNumber"`
	VisitType                  string              `json:"VisitType"`
	VisitDate                  string              `json:"VisitDate"`
	CollectedDate              string              `json:"CollectedDate"`
	VatRegisterationNo         string              `json:"VATRegisterationNo"`
	ClientCode                 string              `json:"ClientCode"`
	ClientId                   string              `json:"ClientID"`
	ClientName                 string              `json:"ClientName"`
	ReferingDoctorCode         string              `json:"ReferingDoctorCode"`
	ReferingDoctorName         string              `json:"ReferingDoctorName"`
	ReferingDoctorMobileNumber string              `json:"ReferingDoctorMobileNumber"`
	IsNotification             string              `json:"IsNotification"`
	HospitalNameCode           string              `json:"HospitalNameCode"`
	HospitalName               string              `json:"HospitalName"`
	VisitHistory               string              `json:"VisitHistory"`
	VisitRemarks               string              `json:"VisitRemarks"`
	RegisterLocation           string              `json:"RegisterLocation"`
	CorporatePatient           int                 `json:"CorporatePatient"`
	WardNo                     string              `json:"WardNo"`
	SrfId                      string              `json:"srfId"`
	ReportLanguage             []map[string]string `json:"ReportLanguage"`
	Partner                    string              `json:"Partner"`
}

type AttuneOrderInfo struct {
	TestID               string                       `json:"TestID"`
	TestCode             string                       `json:"TestCode"`
	OrderedDate          time.Time                    `json:"OrderedDate"`
	TestType             string                       `json:"TestType"`
	TestName             string                       `json:"TestName"`
	TestValueType        string                       `json:"TestValueType"`
	TestValue            string                       `json:"TestValue"`
	UOMCode              string                       `json:"UOMCode"`
	DeviceActualValue    string                       `json:"DeviceActualValue"`
	DeviceID             string                       `json:"DeviceID"`
	MethodName           string                       `json:"MethodName"`
	DepartmentName       string                       `json:"DepartmentName"`
	ReferenceRange       string                       `json:"ReferenceRange"`
	ResultCapturedAt     string                       `json:"ResultCapturedAt"`
	ResultCapturedBy     int                          `json:"ResultCapturedBy"`
	ResultApprovedAt     *string                      `json:"ResultApprovedAt"`
	ResultApprovedBy     int                          `json:"ResultApprovedBy"`
	ApproverName         string                       `json:"ApproverName"`
	TestStatus           string                       `json:"TestStatus"`
	MedicalRemarks       string                       `json:"MedicalRemarks"`
	TechnicalRemarks     string                       `json:"TechnicalRemarks"`
	RerunReason          string                       `json:"RerunReason"`
	RerunRemarks         string                       `json:"RerunRemarks"`
	RerunTime            string                       `json:"RerunTime"`
	CreatedAt            string                       `json:"CreatedAt"`
	UpdatedAt            string                       `json:"UpdatedAt"`
	UserID               int                          `json:"UserID"`
	UserName             string                       `json:"UserName"`
	SampleName           string                       `json:"SampleName"`
	BarcodeNumber        string                       `json:"BarcodeNumber"`
	IMDevice             string                       `json:"IMDevice"`
	IMDeviceFlag         string                       `json:"IMDeviceFlag"`
	SampleCollectedTime  string                       `json:"SampleCollectedTime"`
	SampleReceivedTime   string                       `json:"SampleReceivedTime"`
	SampleRejectedTime   string                       `json:"SampleRejectedTime"`
	SampleCollectedBy    string                       `json:"SampleCollectedBy"`
	SampleReceivedBy     string                       `json:"SampleReceivedBy"`
	SampleRejectedBy     string                       `json:"SampleRejectedBy"`
	OrderContentListInfo []AttuneOrderContentListInfo `json:"OrderContentListInfo"`
	TestDocumentInfo     []AttuneTestDocumentInfo     `json:"TestDocumentInfo"`
	QcFlag               string                       `json:"QcFlag"`
	QcLotNumber          string                       `json:"QcLotNumber"`
	QcValue              string                       `json:"QcValue"`
	QcWestGardWarning    string                       `json:"QcWestGardWarning"`
	QcStatus             string                       `json:"QcStatus"`
}

type AttuneOrderContentListInfo struct {
	TestCode            string                       `json:"TestCode"`
	OrderedDate         time.Time                    `json:"OrderedDate"`
	TestType            string                       `json:"TestType"`
	TestID              string                       `json:"TestID"`
	TestName            string                       `json:"TestName"`
	TestValueType       string                       `json:"TestValueType"`
	TestValue           string                       `json:"TestValue"`
	UOMCode             string                       `json:"UOMCode"`
	DeviceActualValue   string                       `json:"DeviceActualValue"`
	MethodName          string                       `json:"MethodName"`
	DepartmentName      string                       `json:"DepartmentName"`
	ReferenceRange      string                       `json:"ReferenceRange"`
	ResultCapturedAt    time.Time                    `json:"ResultCapturedAt"`
	ResultCapturedBy    int                          `json:"ResultCapturedBy"`
	ResultApprovedAt    *string                      `json:"ResultApprovedAt"`
	ResultApprovedBy    int                          `json:"ResultApprovedBy"`
	ApproverName        string                       `json:"ApproverName"`
	TestStatus          string                       `json:"TestStatus"`
	MedicalRemarks      string                       `json:"MedicalRemarks"`
	TechnicalRemarks    string                       `json:"TechnicalRemarks"`
	DeviceID            string                       `json:"DeviceID"`
	IMDevice            string                       `json:"IMDevice"`
	IMDeviceFlag        string                       `json:"IMDeviceFlag"`
	RerunReason         string                       `json:"RerunReason"`
	RerunRemarks        string                       `json:"RerunRemarks"`
	RerunTime           string                       `json:"RerunTime"`
	CreatedAt           string                       `json:"CreatedAt"`
	UpdatedAt           string                       `json:"UpdatedAt"`
	UserID              int                          `json:"UserID"`
	UserName            string                       `json:"UserName"`
	SampleName          string                       `json:"SampleName"`
	BarcodeNumber       string                       `json:"BarcodeNumber"`
	SampleCollectedTime *time.Time                   `json:"SampleCollectedTime"`
	SampleReceivedTime  *time.Time                   `json:"SampleReceivedTime"`
	SampleRejectedTime  *time.Time                   `json:"SampleRejectedTime"`
	SampleCollectedBy   string                       `json:"SampleCollectedBy"`
	SampleReceivedBy    string                       `json:"SampleReceivedBy"`
	SampleRejectedBy    string                       `json:"SampleRejectedBy"`
	ParameterListInfo   []AttuneOrderContentListInfo `json:"ParameterListInfo"`
	TestDocumentInfo    []AttuneTestDocumentInfo     `json:"TestDocumentInfo"`
	QcFlag              string                       `json:"QcFlag"`
	QcLotNumber         string                       `json:"QcLotNumber"`
	QcValue             string                       `json:"QcValue"`
	QcWestGardWarning   string                       `json:"QcWestGardWarning"`
	QcStatus            string                       `json:"QcStatus"`
}

type AttuneSyncResponse struct {
	Code    string `json:"Code"`
	Message string `json:"Message"`
}

type AttuneBookingResponse struct {
	Status   string             `json:"status"`
	Response AttuneSyncResponse `json:"response"`
}

type VisitDocumentInfo struct {
	VisitDocument string `json:"VisitDocument"`
}

type LisOrderUpdateDetails struct {
	LisVisitId        string                                  `json:"LisVisitId"`
	OrderInfo         map[string]map[string]LisTestUpdateInfo `json:"OrderInfo"`
	PdfResult         string                                  `json:"PdfResult"`
	ReportPdfFormat   string                                  `json:"ReportPDFFormat"`
	VisitDocumentInfo []VisitDocumentInfo                     `json:"VisitDocumentinfo"`
}

type LisTestUpdateInfo struct {
	TestCode string          `json:"TestCode"`
	TestName string          `json:"TestName"`
	MetaData AttuneOrderInfo `json:"MetaData"`
}

type AttuneTestDocumentInfo struct {
	TestDocument string
}

type LisSyncDetails struct {
	VisitID       string          `json:"visitId"`
	PatientName   string          `json:"patientName"`
	SyncedTests   []LisTestInfo   `json:"syncedTestList"`
	SyncedSamples []LisSampleInfo `json:"syncedSampleList"`
	LisSyncTime   string          `json:"lisSyncTime"`
}

type LisTestInfo struct {
	TestCode string `json:"TestCode"`
	TestName string `json:"TestName"`
	TestType string `json:"TestType"`
}

type LisSampleInfo struct {
	SampleName string `json:"sampleName"`
	Barcode    string `json:"barcode"`
}

type AttuneTestDetails struct {
	Status              string
	TestCode            string
	TestName            string
	TestType            string
	TestID              int
	BarcodeNo           []AttuneSampleDetails
	IsStat              string
	Price               string
	RateCardcode        string
	TestClinicalHistory string
	Instructions        string
	Remarks             string
	CreatedAt           string
	UpdatedAt           string
}

type AttuneSampleDetails struct {
	SID         string
	SampleID    string
	ContainerID string
}

type AttuneSyncDataToLisRequest struct {
	MessageType        string
	OrderID            string `json:"OrderId"`
	OrgCode            string
	PatientInfo        AttunePatientInfo
	PatientVisitInfo   AttunePatientVisitInfo
	TestDetailsList    []AttuneTestDetails
	BillingInfo        AttuneBillingInfo
	PaymentDetailsList []AttunePaymentDetails
}

type AttuneSyncDataToLisAfterSyncRequest struct {
	MessageType     string
	OrderID         string `json:"OrderId"`
	OrgCode         string
	TestDetailsList []AttuneTestDetails
}

type AttuneBillingInfo struct {
	PaymentStatus string
	GrossAmount   int
	Discount      int
	NetAmount     int
	DueAmount     int
}

type AttunePaymentDetails struct {
	PaymentType    string
	ReferenceNo    string
	AmountReceived int
}

type AttuneTestSampleMapSnakeCase struct {
	TestId     string `json:"test_id"`
	SampleId   uint   `json:"sample_id"`
	Barcode    string `json:"barcode"`
	VialTypeId uint   `json:"vial_type_id"`
}

type AttuneTestSampleMap struct {
	TestId     string `json:"testId"`
	Barcode    string `json:"barcode"`
	VialTypeId uint   `json:"vialTypeId"`
}

type TestDocumentInfoResponse struct {
	InvestigationId uint   `json:"investigation_id"`
	TestCode        string `json:"test_code"`
	TestDocument    string `json:"test_document"`
}
