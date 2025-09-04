package structures

import (
	"time"

	orderDetailStructures "github.com/Orange-Health/citadel/apps/order_details/structures"
	testDetailStructures "github.com/Orange-Health/citadel/apps/test_detail/structures"
	commonStructures "github.com/Orange-Health/citadel/common/structures"
	commonModels "github.com/Orange-Health/citadel/models"
)

type CollectionSequence struct {
	Sequence        int
	SequenceDetails map[uint]map[uint][]uint
}

type SampleTestMappingDetails struct {
	SampleInfo             commonStructures.SampleInfo
	TestSampleMappingInfos []commonModels.TestSampleMapping
}

type HandleZeroSampleNumberRequest struct {
	TestDetail testDetailStructures.TestDetail
	Tsm        commonStructures.TestSampleMappingInfo
}

type AddBarcodesRequest struct {
	Accessions []AddBarcodeDetails `json:"accessions" binding:"required,min=1"`
}

type AddBarcodeDetails struct {
	Id                 uint       `json:"id"`
	Barcode            string     `json:"barcode"`
	CollectLaterReason string     `json:"collectLaterReason"`
	BarcodeImageURL    string     `json:"barcodeImageUrl"`
	BarcodeScannedTime *time.Time `json:"barcodeScannedTime"`
}

type UpdateAccessionBody struct {
	AccessionId     uint   `json:"accession_id"`
	Barcode         string `json:"barcode"`
	SkipReason      string `json:"skip_reason"`
	BarcodeImageURL string `json:"image_url"`
}

type SampleDetailsRequest struct {
	IsCampTask       bool     `json:"is_camp_task"`
	OrderIds         []string `json:"order_ids"`
	TaskId           uint     `json:"task_id"`
	IsAdditionalTask bool     `json:"is_additional_task"`
}

type VialStructForQuery struct {
	AccessionId      uint       `json:"accession_id"`
	OrderId          string     `json:"order_id"`
	TaskSequence     uint       `json:"task_sequence"`
	VialName         string     `json:"vial_name"`
	VialType         uint       `json:"vial_type"`
	Barcode          string     `json:"barcode"`
	ImageUrl         string     `json:"image_url"`
	ReasonForSkip    string     `json:"reason_for_skip"`
	BarcodeScannedAt *time.Time `json:"barcode_scanned_at"`
	CollectedVolume  uint       `json:"collected_volume"`
	Tests            []uint8    `json:"tests"`
}

type SampleDetailsResponse struct {
	AccessionId      uint                     `json:"accession_id"`
	TaskSequence     uint                     `json:"task_sequence"`
	VialType         uint                     `json:"vial_type"`
	Barcode          string                   `json:"barcode"`
	ImageUrl         string                   `json:"image_url"`
	ReasonForSkip    string                   `json:"reason_for_skip"`
	BarcodeScannedAt *time.Time               `json:"barcode_scanned_at"`
	CollectedVolume  uint                     `json:"collected_volume"`
	Tests            []SchedulerTestsResponse `json:"tests"`
}

type SchedulerTestsResponse struct {
	TestId string `json:"id"`
}

type OrderTestsDetail struct {
	OrderDetail            orderDetailStructures.OrderDetail        `json:"order_detail"`
	TestDetails            []testDetailStructures.TestDetail        `json:"test_details"`
	SampleInfos            []commonStructures.SampleInfo            `json:"sample_infos"`
	TestSampleMappingInfos []commonStructures.TestSampleMappingInfo `json:"test_sample_mapping_infos"`
}

type OrderDetail struct {
	Id               uint       `json:"id"`
	CreatedAt        *time.Time `json:"created_at"`
	UpdatedAt        *time.Time `json:"updated_at"`
	DeletedAt        *time.Time `json:"deleted_at"`
	CreatedBy        string     `json:"created_by"`
	UpdatedBy        string     `json:"updated_by"`
	DeletedBy        string     `json:"deleted_by"`
	OmsOrderId       string     `json:"oms_order_id"`
	OmsRequestId     string     `json:"oms_request_id"`
	OmsCityCode      string     `json:"oms_city_code"`
	PatientDetailsId uint       `json:"patient_details_id"`
	PartnerId        uint       `json:"partner_id"`
	DoctorId         uint       `json:"doctor_id"`
	TrfId            string     `json:"trf_id"`
	TrfStatus        string     `json:"trf_status"`
	ServicingLabId   uint       `json:"servicing_lab_id"`
	CollectionType   string     `json:"collection_type"`
	CollectedAt      *time.Time `json:"collected_at"`
	RequestSource    string     `json:"request_source"`
	BulkOrderId      uint       `json:"bulk_order_id"`
	RequestCreatedAt *time.Time `json:"request_created_at"`
}

type RejectSampleRequest struct {
	RejectionReason  string `json:"rejectionReason"`
	SampleRejectedBy string `json:"sampleRejectedBy"`
	UserId           uint   `json:"userId"`
	LabId            uint   `json:"labId"`
}

type RejectSamplePartiallyRequest struct {
	RejectionReason string `json:"rejectionReason"`
	SampleNumber    uint   `json:"sampleNumber"`
	TestId          string `json:"testId"`
	UserId          uint   `json:"userId"`
	LabId           uint   `json:"labId"`
}

type CreateSampleRequest struct {
	OrderId      string `json:"order_id"`
	VialTypeId   uint   `json:"vial_type_id"`
	SampleNumber uint   `json:"sample_number"`
	Barcode      string `json:"barcode"`
}

type UpdateSampleDetailsForRescheduleRequest struct {
	RequestId        string `json:"request_id"`
	TaskId           uint   `json:"task_id"`
	TaskType         uint   `json:"task_type"`
	OldTaskId        uint   `json:"old_task_id"`
	NewTaskId        uint   `json:"new_task_id"`
	OldTaskRequestId string `json:"old_task_request_id"`
}

type RemapSamplesRequest struct {
	TestIds     []string `json:"test_ids"`
	NewSequence uint     `json:"new_sequence"`
}

type MarkAccessionAsAccessionedRequest struct {
	RequestId string `json:"request_id"`
	SampleId  uint   `json:"sample_id"`
	IsWebhook bool   `json:"is_webhook"`
}

type AddCollectedVolumneRequest struct {
	SampleId uint `json:"sample_id"`
	Volume   uint `json:"volume"`
}

type AddVolumeResponse struct {
	RequestId string `json:"request_id"`
	OrderId   string `json:"order_id"`
}

type ForcefullyMarkSampleAsCollectedRequest struct {
	OmsOrderId    string `json:"oms_order_id"`
	SampleNumbers []uint `json:"sample_numbers"`
	CollectedAt   string `json:"collected_at"`
	UserId        uint   `json:"user_id"`
}

type DelayedReverseLogisticsSamplesDbStruct struct {
	RequestId           string     `json:"request_id"`
	OrderId             string     `json:"order_id"`
	CityCode            string     `json:"city_code"`
	TrfId               string     `json:"trf_id"`
	PatientName         string     `json:"patient_name"`
	PatientAge          *time.Time `json:"patient_age"`
	PatientGender       string     `json:"patient_gender"`
	Barcode             string     `json:"barcode"`
	VialTypeId          uint       `json:"vial_type_id"`
	SampleCollectedTime *time.Time `json:"sample_collected_time"`
}

type DelayedReverseLogisticsSamplesResponse struct {
	RequestId            string     `json:"request_id"`
	OrderId              string     `json:"order_id"`
	TrfId                string     `json:"trf_id"`
	PatientName          string     `json:"patient_name"`
	PatientAge           uint       `json:"patient_age"`
	PatientGender        string     `json:"patient_gender"`
	Barcode              string     `json:"barcode"`
	VialTypeId           uint       `json:"vial_type_id"`
	SampleCollectedTime  *time.Time `json:"sample_collected_time"`
	LogisticsMinuteSpent string     `json:"logistics_minute_spent"`
	VialTypeName         string     `json:"vial_type_name"`
	VialColor            string     `json:"vial_color"`
}

type SrfOrderIdsOrderDetails struct {
	OmsOrderId string `json:"oms_order_id"`
	CityCode   string `json:"city_code"`
}
