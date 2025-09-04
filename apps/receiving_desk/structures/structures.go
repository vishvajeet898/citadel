package structures

import "time"

type CollectedSamplesRequest struct {
	OmsOrderId string `json:"order_id"`
	Barcode    string `json:"barcode"`
	TrfId      string `json:"trf_id"`
	LabId      uint   `json:"lab_id"`
	SearchType string `json:"search_type"`
}

type CollectedSamplesResponse struct {
	CollectedSamples []CollectedSample `json:"samples"`
}

type CollectedSample struct {
	Id                 uint                     `json:"id"`
	ParentSampleId     uint                     `json:"parent_sample_id"`
	OrderId            string                   `json:"order_id"`
	RequestId          string                   `json:"request_id"`
	ServicingCityCode  string                   `json:"servicing_city_code"`
	Barcode            string                   `json:"barcode"`
	TrfId              string                   `json:"trf_id"`
	VialType           uint                     `json:"vial_type"`
	SampleNumber       uint                     `json:"sample_number"`
	Status             string                   `json:"status"`
	LabDetails         LabDetailsRdResponse     `json:"lab_details"`
	CollectLaterReason string                   `json:"collect_later_reason,omitempty"`
	NotReceivedReason  string                   `json:"not_received_reason,omitempty"`
	PatientDetails     PatientDetailsRdResponse `json:"patient_details"`
	Tests              []TestDetailsRdResponse  `json:"tests"`
}

type LabDetailsRdResponse struct {
	Id      uint   `json:"id"`
	Name    string `json:"name"`
	LabType string `json:"lab_type"`
}

type PatientDetailsRdResponse struct {
	Dob         *time.Time `json:"dob"`
	ExpectedDob *time.Time `json:"expected_dob"`
	Name        string     `json:"name"`
	Age         uint       `json:"age"`
	Gender      string     `json:"gender"`
}

type TestDetailsRdResponse struct {
	Id              string `json:"id"`
	Name            string `json:"name"`
	ProcessingLabId uint   `json:"processing_lab_id"`
}

type ReceiveSamplesRequest struct {
	UserId         uint            `json:"user_id"`
	ReceivingLabId uint            `json:"receiving_lab_id"`
	Samples        []SampleDetails `json:"samples"`
}

type SampleDetails struct {
	Id      uint   `json:"id"`
	Barcode string `json:"barcode"`
}
