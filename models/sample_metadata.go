package models

import "time"

type SampleMetadata struct {
	BaseModel
	OmsCityCode              string     `gorm:"column:oms_city_code" json:"oms_city_code,omitempty"`
	OmsOrderId               string     `gorm:"column:oms_order_id" json:"oms_order_id,omitempty"`
	SampleId                 uint       `gorm:"column:sample_id" json:"sample_id,omitempty"`
	CollectionSequenceNumber uint       `gorm:"column:collection_sequence_number" json:"collection_sequence_number,omitempty"`
	LastUpdatedAt            *time.Time `gorm:"column:last_updated_at" json:"last_updated_at,omitempty"`
	TransferredAt            *time.Time `gorm:"column:transferred_at" json:"transferred_at,omitempty"`
	OutsourcedAt             *time.Time `gorm:"column:outsourced_at" json:"outsourced_at,omitempty"`
	CollectedAt              *time.Time `gorm:"column:collected_at" json:"collected_at,omitempty"`
	ReceivedAt               *time.Time `gorm:"column:received_at" json:"received_at,omitempty"`
	AccessionedAt            *time.Time `gorm:"column:accessioned_at" json:"accessioned_at,omitempty"`
	RejectedAt               *time.Time `gorm:"column:rejected_at" json:"rejected_at,omitempty"`
	NotReceivedAt            *time.Time `gorm:"column:not_received_at" json:"not_received_at,omitempty"`
	LisSyncAt                *time.Time `gorm:"column:lis_sync_at" json:"lis_sync_at,omitempty"`
	BarcodeScannedAt         *time.Time `gorm:"column:barcode_scanned_at" json:"barcode_scanned_at,omitempty"`
	NotReceivedReason        string     `gorm:"column:not_received_reason" json:"not_received_reason,omitempty"`
	BarcodeImageUrl          string     `gorm:"column:barcode_image_url" json:"barcode_image_url,omitempty"`
	TaskSequence             uint       `gorm:"column:task_sequence" json:"task_sequence,omitempty"`
	CollectLaterReason       string     `gorm:"column:collect_later_reason" json:"collect_later_reason,omitempty"`
	RejectingLab             uint       `gorm:"column:rejecting_lab" json:"rejecting_lab,omitempty"`
	CollectedVolume          uint       `gorm:"column:collected_volume" json:"collected_volume,omitempty"`
}

func (SampleMetadata) TableName() string {
	return "sample_metadata"
}
