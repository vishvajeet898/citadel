package models

import "time"

type TestDetailsMetadata struct {
	BaseModel
	TestDetailsId    uint       `gorm:"column:test_details_id;not null" json:"test_details_id"`
	Barcodes         string     `gorm:"column:barcodes;not null;type:varchar(255)" json:"barcodes"`
	IsCritical       bool       `gorm:"column:is_critical;not null" json:"is_critical"`
	IsCompletedInOms bool       `gorm:"column:is_completed_in_oms" json:"is_completed_in_oms"`
	PickedAt         *time.Time `gorm:"column:picked_at;type:timestamp" json:"picked_at"`
}

func (TestDetailsMetadata) TableName() string {
	return "test_details_metadata"
}
