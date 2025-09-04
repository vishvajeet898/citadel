package structures

import (
	commonStructures "github.com/Orange-Health/citadel/common/structures"
)

type TaskMetadata struct {
	commonStructures.BaseStruct
	TaskID       uint   `json:"task_id,omitempty"`
	IsPackage    bool   `json:"is_package,omitempty"`
	IsMorphle    bool   `json:"is_morphle,omitempty"`
	IsCritical   bool   `json:"is_critical,omitempty"`
	DoctorName   string `json:"doctor_name,omitempty"`
	DoctorNumber string `json:"doctor_number,omitempty"`
	DoctorNotes  string `json:"doctor_notes,omitempty"`
}
