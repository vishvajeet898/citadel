package structures

import (
	"time"

	commonStructures "github.com/Orange-Health/citadel/common/structures"
)

// @swagger:response TestDetail
type TestDetail struct {
	commonStructures.BaseStruct
	// The task ID of the test.
	// example: 2345
	TaskID uint `json:"task_id"`
	// The OMS test ID of the test.
	// example: 54321
	OmsTestID uint `json:"oms_test_id"`
	// The Central OMS test ID of the test.
	// example: TB54321
	CentralOmsTestId string `gorm:"column:central_oms_test_id;not null" json:"central_oms_test_id"`
	// The test name of the test.
	// example: "CBC"
	TestName string `json:"test_name"`
	// The test type of the test.
	// example: "Open"
	Status string `json:"status"`
	// Doctor TAT to complete the test result review.
	// example: "2022-10-01T00:00:00"
	DoctorTat *time.Time `json:"doctor_tat"`
	// The test result of the test.
	// example: "N"
	IsAutoApproved bool `json:"is_auto_approved"`
	// The test Report Sent At.
	// example: "2022-10-01T00:00:00"
	ReportSentAt *time.Time `json:"report_sent_at"`
	// The test is manual report upload.
	// example: true
	IsManualReportUpload bool `json:"is_manual_report_upload"`
	// The master test id of the test.
	// example: 1234
	MasterTestId uint `json:"master_test_id"`
	// Master Package Id of the test.
	// example: 1234
	MasterPackageId uint `json:"master_package_id"`
	// The OMS order id of the test.
	// example: 12345
	OmsOrderId      string `gorm:"column:oms_order_id;not null" json:"oms_order_id"`
	Department      string `gorm:"column:department;not null" json:"department"`
	LabId           uint   `gorm:"column:lab_id;not null" json:"lab_id"`
	ProcessingLabId uint   `gorm:"column:processing_lab_id;not null" json:"processing_lab_id"`
	TestType        string `gorm:"column:test_type;not null" json:"test_type"`
	LisCode         string `gorm:"column:lis_code;not null" json:"lis_code"`
	ReportStatus    string `gorm:"column:report_status;not null" json:"report_status"`
	// Need to change this to OMS City Code
	CityCode string `gorm:"column:city_code;not null" json:"city_code"`
}

type TestBasicDetails struct {
	Id             uint       `json:"id,omitempty"`
	TaskId         uint       `json:"task_id,omitempty"`
	LabId          uint       `json:"lab_id,omitempty"`
	OmsOrderId     string     `json:"oms_order_id,omitempty"`
	OmsTestId      string     `json:"oms_test_id,omitempty"`
	TestName       string     `json:"test_name,omitempty"`
	Status         string     `json:"status,omitempty"`
	DoctorTat      *time.Time `json:"doctor_tat,omitempty"`
	IsAutoApproved bool       `json:"is_auto_approved,omitempty"`
	ReportSentAt   *time.Time `json:"report_sent_at,omitempty"`
}

type LabIdRequest struct {
	Id                  uint `json:"id"`
	LabId               uint `json:"lab_id"`
	RecomputeAccessions bool `json:"recompute_accessions"`
}

type TestBasicDetailsByOmsOrderResponse struct {
	OrderTestDetails TestBasicDetailsByOmsOrderId `json:"order_test_details"`
}

type TestBasicDetailsByOmsOrderId struct {
	InhouseTests    []OmsTestBasicDetails `json:"inhouse_tests"`
	OutsourcedTests []OmsTestBasicDetails `json:"outsourced_tests"`
	CityCode        string                `json:"city_code"`
}

type OmsTestBasicDetails struct {
	OrderID      string `json:"order_id"`
	RequestID    string `json:"request_id"`
	TestId       string `json:"test_id"`
	TestName     string `json:"test_name"`
	TestStatus   string `json:"test_status"`
	MasterTestId uint   `json:"master_test_id"`
	LabId        uint   `json:"lab_id"`
	Inhouse      bool   `json:"inhouse"`
}

type OmsTestBasicDetailsDbStruct struct {
	TestId       string `json:"test_id"`
	TestName     string `json:"test_name"`
	TestStatus   string `json:"test_status"`
	MasterTestId uint   `json:"master_test_id"`
	LabId        uint   `json:"lab_id"`
}

type UpdateProcessingLabRequest struct {
	TestDetails []UpdateProcessingLabTestDetails `json:"test_details"`
	UserId      uint                             `json:"user_id"`
}

type UpdateProcessingLabTestDetails struct {
	OmsTestId       string `json:"oms_test_id"`
	ProcessingLabId uint   `json:"processing_lab_id"`
}
