package models

type TestSampleMapping struct {
	BaseModel
	OmsCityCode         string `gorm:"column:oms_city_code" json:"oms_city_code,omitempty"`
	OmsTestId           string `gorm:"column:oms_test_id" json:"oms_test_id,omitempty"`
	VialTypeId          uint   `gorm:"column:vial_type_id" json:"vial_type_id,omitempty"`
	SampleId            uint   `gorm:"column:sample_id" json:"sample_id,omitempty"`
	SampleNumber        uint   `gorm:"column:sample_number" json:"sample_number,omitempty"`
	OmsOrderId          string `gorm:"column:oms_order_id" json:"oms_order_id,omitempty"`
	RecollectionPending bool   `gorm:"column:recollection_pending" json:"recollection_pending,omitempty"`
	IsRejected          bool   `gorm:"column:is_rejected" json:"is_rejected,omitempty"`
	RejectionReason     string `gorm:"column:rejection_reason" json:"rejection_reason,omitempty"`
}

func (TestSampleMapping) TableName() string {
	return "test_sample_mapping"
}
