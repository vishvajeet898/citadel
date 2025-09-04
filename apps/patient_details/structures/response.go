package structures

import (
	"time"

	commonStructures "github.com/Orange-Health/citadel/common/structures"
)

// @swagger:response patientDetail
type PatientDetail struct {
	commonStructures.BaseStruct
	// The name of the patient.
	// example: "Jane Doe"
	Name string `json:"name"`
	// The gender of the patient.
	// example: "Male/Female"
	Gender string `json:"gender"`
	// The date of birth of the patient.
	// example: "1997-01-01"
	Dob *time.Time `json:"dob"`
	// The expected date of birth of the patient.
	// example: "1997-01-01"
	ExpectedDob *time.Time `json:"expected_dob"`
	// The number of the patient.
	// example: "9087654321"
	Number string `json:"number"`
	// The system patient ID of the patient.
	// example: "12345"
	SystemPatientID string `json:"system_patient_id"`
}
