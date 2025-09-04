package structures

import (
	"time"
)

type PatientPastRecords struct {
	OrderId               uint       `json:"order_id"`
	InvestigationValue    string     `json:"investigation_value,omitempty"`
	Uom                   string     `json:"uom,omitempty"`
	MasterInvestigationId uint       `json:"master_investigation_id"`
	Department            string     `json:"department"`
	InvestigationName     string     `json:"investigation_name"`
	ApprovedAt            *time.Time `json:"approved_at"`
	CityCode              string     `json:"city_code"`
}

type DeltaValuesStruct struct {
	OrderId               uint       `json:"order_id"`
	InvestigationValue    string     `json:"investigation_value,omitempty"`
	MasterInvestigationId uint       `json:"master_investigation_id"`
	ApprovedAt            *time.Time `json:"approved_at"`
	CityCode              string     `json:"city_code"`
}

type DeltaValuesStructResponse struct {
	InvestigationValue    string     `json:"investigation_value,omitempty"`
	MasterInvestigationId uint       `json:"master_investigation_id"`
	ApprovedAt            *time.Time `json:"approved_at"`
}

type PatientPastRecordsApiResponse struct {
	PatientDetails PatientPastRecordsPatientDetails `json:"patient_details"`
	Dates          []string                         `json:"dates"`
	Data           []PatientPastRecordsData         `json:"data"`
}

type PatientPastRecordsPatientDetails struct {
	Name      string `json:"name"`
	AgeYears  uint   `json:"age_years"`
	AgeMonths uint   `json:"age_months"`
	AgeDays   uint   `json:"age_days"`
	Gender    string `json:"gender"`
}

type PatientPastRecordsData struct {
	Department     string                             `json:"department"`
	Investigations []PatientPastRecordsInvestigations `json:"investigations"`
}

type PatientPastRecordsInvestigations struct {
	InvestigationId   uint                                 `json:"investigation_id"`
	InvestigationName string                               `json:"investigation_name"`
	PastRecords       map[string][]PatientPastRecordsValue `json:"past_records"`
}

type PatientPastRecordsValue struct {
	InvestigationValue string `json:"investigation_value,omitempty"`
	InvestigationData  string `json:"investigation_data,omitempty"`
	Uom                string `json:"uom,omitempty"`
}
