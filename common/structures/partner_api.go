package structures

type Partner struct {
	Id                uint   `json:"id"`
	PartnerName       string `json:"partnerName"`
	ReportFormat      uint   `json:"reportFormat"`
	CobrandedImageUrl string `json:"cobrandedImageUrl"`
}
