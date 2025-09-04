package structures

type EventPayload struct {
	EventType    string
	EventPayload string
	TraceID      string
	Contains     string
	RedisKey     string
	GroupId      string
}

type ReportReadyEvent struct {
	ReportPdfEvent    ReportPdfEvent `json:"report_pdf"`
	ServicingCityCode string         `json:"servicing_city_code"`
}

type ReportPdfEvent struct {
	JobUUID                string   `json:"job_uuid"`
	OrderID                string   `json:"order_id"`
	CityCode               string   `json:"city_code"`
	TestIds                []string `json:"test_ids"`
	AttuneFiles            bool     `json:"attune_files"`
	VisitId                string   `json:"visit_id"`
	ReportPdfBrandedURL    string   `json:"report_pdf_branded_url"`
	ReportPdfCobrandedURL  string   `json:"report_pdf_cobranded_url"`
	ReportPdfHeaderlessURL string   `json:"report_pdf_headerless_url"`
	IsDummyReport          bool     `json:"is_dummy_report"`
}

type AbnormalityStruct struct {
	MasterInvestigationId uint   `json:"master_investigation_id"`
	Abnormality           string `json:"abnormality"`
}

type AbnormalityUpdateEvent struct {
	OrderId       uint                `json:"order_id"`
	CityCode      string              `json:"city_code"`
	Service       string              `json:"service"`
	Abnormalities []AbnormalityStruct `json:"abnormalities"`
}

type MasterContact struct {
	Id               uint   `json:"id"`
	MatchAlgoVersion uint   `json:"match_algo_version"`
	Uuid             string `json:"uuid"`
	Salutation       string `json:"salutation"`
	FirstName        string `json:"first_name"`
	LastName         string `json:"last_name"`
	Gender           string `json:"gender"`
	ContactMergeLog  string `json:"contact_merge_log"`
	DateOfBirth      string `json:"date_of_birth"`
}

type MergeContactEvent struct {
	MergeContact         string        `json:"merge_contact"`
	MasterContact        MasterContact `json:"master_contact"`
	Service              string        `json:"service"`
	WasMergeContactFound bool          `json:"was_merge_contact_found"`
}
