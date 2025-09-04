package structures

type OmsVisitDetailsByTestIdsResponse struct {
	VisitDetails []OmsVisitDetailsMap `json:"visit_details"`
}

type OmsVisitDetailsMap struct {
	TestId  uint `json:"test_id"`
	VisitId uint `json:"visit_id"`
}

type DeltaValuesResponse struct {
	Data []DeltaValuesStruct `json:"data"`
}

type PatientPastRecordsResponse struct {
	Data []PatientPastRecords `json:"data"`
}

type VisitDetailsByTestIdsResponse struct {
	VisitDetails []VisitDetailsTestIdMapping `json:"visit_details"`
}

type VisitDetailsTestIdMapping struct {
	VisitID string `json:"visit_id"`
	TestId  uint   `json:"test_id"`
}
