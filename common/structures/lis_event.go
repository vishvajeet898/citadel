package structures

type LisEvent struct {
	ContentType string `json:"content_type"`
	EntityID    string `json:"entity_id"`
	EntityName  string `json:"entity_name"`
	WebhookData string `json:"webhook_data"`
}

type LisOrderInfoEvent struct {
	OrderId           string                   `json:"OrderId"`
	OrgCode           string                   `json:"OrgCode"`
	OverAllStatus     string                   `json:"OverAllStatus"`
	VisitDocumentinfo []VisitDocumentInfoEvent `json:"VisitDocumentinfo"`
	OrderInfo         []AttuneOrderInfoEvent   `json:"OrderInfo"`
	ResultAsPdf       string                   `json:"ResultAsPdf"`
	ReportPdfFormat   string                   `json:"ReportPDFFormat"`
}

type VisitDocumentInfoEvent struct {
	VisitDocument string `json:"VisitDocument"`
}

type LisOrderUpdateDetailsEvent struct {
	LisVisitId        string                                       `json:"LisVisitId"`
	OrderInfo         map[string]map[string]LisTestUpdateInfoEvent `json:"OrderInfo"`
	PdfResult         string                                       `json:"PdfResult"`
	ReportPdfFormat   string                                       `json:"ReportPDFFormat"`
	VisitDocumentInfo []VisitDocumentInfoEvent                     `json:"VisitDocumentinfo"`
}

type DocumentWithMetadata struct {
	Url             string
	InvestigationId uint
}
