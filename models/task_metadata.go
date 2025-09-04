package models

import (
	"time"
)

type TaskMetadata struct {
	BaseModel
	TaskId          uint       `gorm:"column:task_id;not null" json:"task_id"`
	Task            Task       `gorm:"foreignKey:TaskId;references:Id"`
	ContainsPackage bool       `gorm:"column:contains_package;not null" json:"contains_package"`
	ContainsMorphle bool       `gorm:"column:contains_morphle;not null" json:"contains_morphle"`
	IsCritical      bool       `gorm:"column:is_critical;not null" json:"is_critical"`
	DoctorName      string     `gorm:"column:doctor_name;type:varchar(255)" json:"doctor_name"`
	DoctorNumber    string     `gorm:"column:doctor_number;type:varchar(15)" json:"doctor_number"`
	DoctorNotes     string     `gorm:"column:doctor_notes;type:text" json:"doctor_notes"`
	PartnerName     string     `gorm:"column:partner_name;type:varchar(255)" json:"partner_name"`
	LastEventSentAt *time.Time `gorm:"column:last_event_sent_at;type:timestamp" json:"last_event_sent_at"`
}

func (TaskMetadata) TableName() string {
	return "task_metadata"
}
