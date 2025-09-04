package models

import (
	"time"
)

type PatientDetail struct {
	BaseModel
	Name            string     `gorm:"column:name;not null;type:varchar(255)" json:"name"`
	Dob             *time.Time `gorm:"column:dob;not null;type:date" json:"dob"`
	ExpectedDob     *time.Time `gorm:"column:expected_dob;type:date" json:"expected_dob"`
	Gender          string     `gorm:"column:gender;not null;type:varchar(10)" json:"gender"`
	Number          string     `gorm:"column:number;not null;type:varchar(20)" json:"number"`
	SystemPatientId string     `gorm:"column:system_patient_id;type:varchar(50)" json:"system_patient_id"`
}

func (PatientDetail) TableName() string {
	return "patient_details"
}
