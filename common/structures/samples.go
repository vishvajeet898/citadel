package structures

import "time"

type SampleInfo struct {
	Id                          uint       `json:"id"`
	OmsCityCode                 string     `json:"oms_city_code"`
	OmsOrderId                  string     `json:"oms_order_id"`
	OmsRequestId                string     `json:"oms_request_id"`
	LabId                       uint       `json:"lab_id"`
	VisitId                     string     `json:"visit_id"`
	VialTypeId                  uint       `json:"vial_type_id"`
	DestinationLabId            uint       `json:"destination_lab_id,omitempty"`
	CollectionSequenceNumber    uint       `json:"collection_sequence_number,omitempty"`
	Barcode                     string     `json:"barcode,omitempty"`
	Status                      string     `json:"status"`
	SampleId                    uint       `json:"sample_id,omitempty"`
	SampleNumber                uint       `json:"sample_number,omitempty"`
	ParentSampleId              uint       `json:"parent_sample_id,omitempty"`
	RejectionReason             string     `json:"rejection_reason,omitempty"`
	NotReceivedReason           string     `json:"not_received_reason,omitempty"`
	CollectLaterReason          string     `json:"collectLater_reason,omitempty"`
	TaskSequence                uint       `json:"task_sequence,omitempty"`
	BarcodeImageUrl             string     `json:"barcode_image_url,omitempty"`
	LastUpdatedAt               *time.Time `json:"last_updated_at,omitempty"`
	LisSyncAt                   *time.Time `json:"lis_sync_at,omitempty"`
	BarcodeScannedAt            *time.Time `json:"barcode_scanned_at,omitempty"`
	CollectedAt                 *time.Time `json:"collected_at,omitempty"`
	ReceivedAt                  *time.Time `json:"received_at,omitempty"`
	RejectedAt                  *time.Time `json:"rejected_at,omitempty"`
	NotReceivedAt               *time.Time `json:"not_received_at,omitempty"`
	AccessionedAt               *time.Time `json:"accessioned_at,omitempty"`
	TransferTime                uint       `json:"transfer_time,omitempty"`
	TransferredAt               *time.Time `json:"transferred_at,omitempty"`
	OutsourcedAt                *time.Time `json:"outsourced_at,omitempty"`
	RejectingLab                uint       `json:"rejecting_lab,omitempty"`
	CollectedVolume             uint       `json:"collected_volume,omitempty"`
	TransferredSampleReceivedAt *time.Time `json:"transferred_sample_received_at,omitempty"`
	CreatedAt                   *time.Time `json:"created_at"`
	CreatedBy                   uint       `json:"created_by"`
	UpdatedAt                   *time.Time `json:"updated_at"`
	UpdatedBy                   uint       `json:"updated_by"`
	DeletedBy                   uint       `json:"deleted_by,omitempty"`
	DeletedAt                   *time.Time `json:"deleted_at,omitempty"`
}

type SampleDetailsWithTestInfo struct {
	SampleInfo
	IsRejected      bool `json:"is_rejected"`
	IsQnsApplicable bool `json:"isQnsApplicable"`
	OmsTestId       uint `json:"oms_test_id"`
}

type SampleCollectedRequest struct {
	RequestId      string `json:"request_id"`
	OmsTaskType    uint   `json:"oms_task_type"`
	OmsTaskId      uint   `json:"oms_task_id"`
	CollectedAt    string `json:"collected_at"`
	CollectionType string `json:"collection_type"`
	IsB2c          bool   `json:"is_b2c"`
	UserId         uint   `json:"user_id"`
}

type UpdateSampleDetailsPostTaskCompletionRequest struct {
	RequestId      string `json:"request_id"`
	OmsTaskType    uint   `json:"oms_task_type"`
	OmsTaskId      uint   `json:"oms_task_id"`
	CollectedAt    string `json:"collected_at"`
	CollectionType string `json:"collection_type"`
	IsB2c          bool   `json:"is_b2c"`
	TaskSequence   uint   `json:"task_sequence"`
	UserId         uint   `json:"user_id"`
}

type MarkAsNotReceivedRequest struct {
	UserId            uint   `json:"user_id"`
	LabId             uint   `json:"lab_id"`
	SampleId          uint   `json:"sample_id"`
	NotReceivedReason string `json:"not_received_reason"`
}

type UpdateTaskSequenceRequest struct {
	RequestId string   `json:"request_id"`
	TaskId    uint     `json:"task_id"`
	TestIds   []string `json:"test_ids"`
}

type UpdateSrfIdToLisRequest struct {
	AlnumOrderId string `json:"alnum_order_id"`
	SrfId        string `json:"srf_id"`
}
