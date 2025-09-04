package models

type Sample struct {
	BaseModel
	OmsCityCode      string `gorm:"column:oms_city_code" json:"city_code,omitempty"`
	OmsOrderId       string `gorm:"column:oms_order_id" json:"oms_order_id,omitempty"`
	OmsRequestId     string `gorm:"column:oms_request_id" json:"oms_request_id,omitempty"`
	VisitId          string `gorm:"column:visit_id" json:"visit_id,omitempty"`
	VialTypeId       uint   `gorm:"column:vial_type_id" json:"vial_type_id,omitempty"`
	Barcode          string `gorm:"column:barcode" json:"barcode,omitempty"`
	Status           string `gorm:"column:status" json:"status,omitempty"`
	LabId            uint   `gorm:"column:lab_id" json:"lab_id,omitempty"`
	DestinationLabId uint   `gorm:"column:destination_lab_id" json:"destination_lab_id,omitempty"`
	RejectionReason  string `gorm:"column:rejection_reason" json:"rejection_reason,omitempty"`
	ParentSampleId   uint   `gorm:"column:parent_sample_id" json:"parent_sample_id,omitempty"`
	SampleNumber     uint   `gorm:"column:sample_number" json:"sample_number,omitempty"`
}

func (Sample) TableName() string {
	return "samples"
}
