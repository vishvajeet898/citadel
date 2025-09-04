package structures

type AutoApprovalSlackMessageAttributes struct {
	RequestID          string
	OrderId            string
	CityCode           string
	PatientName        string
	VisitID            string
	TestNames          string
	InvestigationsName string
}
