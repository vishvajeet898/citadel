package models

import (
	"time"
)

type Task struct {
	BaseModel
	OrderId          uint          `gorm:"column:order_id;not null" json:"order_id"`
	RequestId        uint          `gorm:"column:request_id;not null" json:"request_id"`
	OmsOrderId       string        `gorm:"column:oms_order_id" json:"oms_order_id"`     // Central order id
	OmsRequestId     string        `gorm:"column:oms_request_id" json:"oms_request_id"` // Central Request Id
	LabId            uint          `gorm:"column:lab_id;not null" json:"lab_id"`
	CityCode         string        `gorm:"column:city_code;not null;type:varchar(10)" json:"city_code"`
	Status           string        `gorm:"column:status;not null;type:varchar(50)" json:"status"`
	PreviousStatus   string        `gorm:"column:previous_status;type:varchar(50)" json:"previous_status"`
	OrderType        string        `gorm:"column:order_type;not null;type:varchar(50)" json:"order_type"`
	PatientDetailsId uint          `gorm:"column:patient_details_id;not null" json:"patient_details_id"`
	PatientDetail    PatientDetail `gorm:"foreignKey:PatientDetailsId;references:Id"`
	DoctorTat        *time.Time    `gorm:"column:doctor_tat;type:timestamp" json:"doctor_tat"`
	IsActive         bool          `gorm:"column:is_active;not null" json:"is_active"` // TODO: Drop this column
	CompletedAt      *time.Time    `gorm:"column:completed_at;type:timestamp" json:"completed_at"`
}

func (Task) TableName() string {
	return "tasks"
}
