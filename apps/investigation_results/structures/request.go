package structures

import "time"

type DeltaValuesRequest struct {
	PatientId              string `json:"patient_id"`
	MasterInvestigationIds []uint `json:"master_investigation_ids"`
	Limit                  uint   `json:"limit"`
}

type BasicAbnormalityStruct struct {
	LabId              uint       `json:"lab_id"`
	CityCode           string     `json:"city_code"`
	PatientDob         *time.Time `json:"patient_dob"`
	PatientExpectedDob *time.Time `json:"patient_expected_dob"`
	PatientGender      string     `json:"patient_gender"`
}

type ModifyValueRequest struct {
	MasterInvestigationId uint   `json:"master_investigation_id"`
	InvestigationValue    string `json:"investigation_value"`
}

type ModifyValueApiRequest struct {
	CurrentInvestigation ModifyValueRequest   `json:"current_investigation"`
	PastInvestigations   []ModifyValueRequest `json:"past_investigations"`
}
