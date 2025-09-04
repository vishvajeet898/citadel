package models

import (
	"time"

	"github.com/google/uuid"
)

type OrderDetails struct {
	BaseModel
	OmsOrderId       string     `gorm:"column:oms_order_id;not null" json:"oms_order_id"`
	OmsRequestId     string     `gorm:"column:oms_request_id;not null" json:"oms_request_id"`
	Uuid             uuid.UUID  `gorm:"column:uuid" json:"uuid"`
	CityCode         string     `gorm:"column:city_code;not null" json:"city_code"`
	PatientDetailsId uint       `gorm:"column:patient_details_id;not null" json:"patient_details_id"`
	OrderStatus      string     `gorm:"column:order_status;not null" json:"order_status"`
	PartnerId        uint       `gorm:"column:partner_id" json:"partner_id"`
	DoctorId         uint       `gorm:"column:doctor_id" json:"doctor_id"`
	TrfId            string     `gorm:"column:trf_id;type:varchar(255)" json:"trf_id"`
	ServicingLabId   uint       `gorm:"column:servicing_lab_id;not null" json:"servicing_lab_id"`
	CollectionType   uint       `gorm:"column:collection_type;not null" json:"collection_type"`
	RequestSource    string     `gorm:"column:request_source" json:"request_source"`
	BulkOrderId      uint       `gorm:"column:bulk_order_id" json:"bulk_order_id"`
	CampId           uint       `gorm:"column:camp_id" json:"camp_id"`
	ReferredBy       string     `gorm:"column:referred_by" json:"referred_by"`
	CollectedOn      *time.Time `gorm:"column:collected_on" json:"collected_on"`
	SrfId            string     `gorm:"column:srf_id;type:varchar(255)" json:"srf_id"`
}

func (OrderDetails) TableName() string {
	return "order_details"
}
