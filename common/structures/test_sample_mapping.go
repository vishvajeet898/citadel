package structures

import "time"

type TestSampleMappingInfo struct {
	Id                  uint       `json:"id"`
	OmsCityCode         string     `json:"omsCityCode"`
	SampleId            uint       `json:"sampleId"`
	TestDetailsId       uint       `json:"testDetailsId"`
	OmsTestId           string     `json:"omsTestId"`
	SampleNumber        uint       `json:"sampleNumber"`
	VialTypeId          uint       `json:"vialTypeId"`
	OmsOrderId          string     `json:"omsOrderId"`
	IsRejected          bool       `json:"isRejected"`
	RecollectionPending bool       `json:"recollectionPending"`
	RejectionReason     string     `json:"rejectionReason"`
	CreatedAt           *time.Time `json:"createdAt"`
	CreatedBy           uint       `json:"createdBy"`
	UpdatedAt           *time.Time `json:"updatedAt"`
	UpdatedBy           uint       `json:"updatedBy"`
	DeletedAt           *time.Time `json:"deletedAt"`
	DeletedBy           uint       `json:"deletedBy"`
}
