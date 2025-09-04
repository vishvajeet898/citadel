package structures

import "time"

type EtsTestDbEvent struct {
	TestID             string     `json:"test_id"`
	OrderID            string     `json:"order_id"`
	TestStatus         string     `json:"test_status"`
	OmsTestStatus      string     `json:"oms_test_status"`
	TestName           string     `json:"test_name"`
	TestDeletedAt      *time.Time `json:"test_deleted_at"`
	MasterTestID       uint       `json:"master_test_id"`
	LabEta             string     `json:"lab_eta"`
	LabID              uint       `json:"lab_id"`
	PatientName        string     `json:"patient_name"`
	PatientDob         *time.Time `json:"patient_dob"`
	PatientExpectedDob *time.Time `json:"patient_expected_dob"`
	PatientGender      string     `json:"patient_gender"`
	Barcode            string     `json:"barcode"`
	VialTypeID         uint       `json:"vial_type_id"`
	IsRejected         bool       `json:"is_rejected,omitempty"`
	CityCode           string     `json:"city_code"`
	LisStatus          string     `json:"lis_status,omitempty"`
}

type EtsTestEvent struct {
	TestID        string     `json:"test_id"`
	OrderID       string     `json:"order_id"`
	TestStatus    uint       `json:"test_status"`
	TestName      string     `json:"test_name"`
	TestDeletedAt *time.Time `json:"test_deleted_at"`
	MasterTestID  uint       `json:"master_test_id"`
	LabEta        string     `json:"lab_eta"`
	LabID         uint       `json:"lab_id"`
	PatientName   string     `json:"patient_name"`
	PatientAge    uint       `json:"patient_age"`
	PatientGender string     `json:"patient_gender"`
	Barcode       string     `json:"barcode"`
	VialTypeID    uint       `json:"vial_type_id"`
	IsRejected    bool       `json:"is_rejected,omitempty"`
	CityCode      string     `json:"city_code"`
	LisStatus     string     `json:"lis_status,omitempty"`
}
