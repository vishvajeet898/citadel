package structures

type MasterInvestigationResponse struct {
	InvestigationDetails InvestigationDetails `json:"investigation_details"`
}

type InvestigationDetails struct {
	Investigations []Investigation `json:"investigations"`
	Panels         []Panel         `json:"panels"`
}

type Investigation struct {
	InvestigationId              uint                 `json:"investigation_id,omitempty"`
	InvestigationName            string               `json:"investigation_name,omitempty"`
	InvestigationMethodMappingId uint                 `json:"investigation_method_mapping_id,omitempty"`
	LisCode                      string               `json:"lis_code,omitempty"`
	Status                       string               `json:"status,omitempty"`
	DepartmentName               string               `json:"department_name,omitempty"`
	Method                       string               `json:"method,omitempty"`
	Unit                         string               `json:"unit,omitempty"`
	ResultRepresentationType     string               `json:"result_representation_type,omitempty"`
	IsNablAccredited             bool                 `json:"is_nabl_accredited,omitempty"`
	ReferenceRange               MasterReferenceRange `json:"reference_range,omitempty"`
	ResultType                   string               `json:"result_type,omitempty"`
	RcvPositive                  float64              `json:"rcv_positive,omitempty"`
	RcvNegative                  float64              `json:"rcv_negative,omitempty"`
	PastValueThresholdDays       uint                 `json:"past_value_threshold_days,omitempty"`
}

type Panel struct {
	PanelName        string          `json:"panel_name,omitempty"`
	TestId           uint            `json:"test_id,omitempty"`
	PanelTestCode    string          `json:"panel_test_code,omitempty"`
	Status           string          `json:"status,omitempty"`
	PanelDisplayName string          `json:"panel_display_name,omitempty"`
	Investigations   []Investigation `json:"investigations"`
	Panels           []Panel         `json:"panels"`
}

type InvestigationMasterRequest struct {
	TestIds []uint `json:"test_ids"`
	LabId   uint   `json:"lab_id"`
}

type MasterReferenceRange struct {
	NormalRange       Range  `json:"normal_range,omitempty"`
	CriticalRange     Range  `json:"critical_range,omitempty"`
	ImprobableRange   Range  `json:"improbable_range,omitempty"`
	AutoApprovalRange Range  `json:"auto_approval_range,omitempty"`
	ReferenceLabel    string `json:"reference_label,omitempty"`
}

type Range struct {
	MinAge             uint   `json:"min_age,omitempty"`
	MaxAge             uint   `json:"max_age,omitempty"`
	AgeUom             string `json:"age_uom,omitempty"`
	Gender             string `json:"gender,omitempty"`
	MinValue           string `json:"min_value,omitempty"`
	MaxValue           string `json:"max_value,omitempty"`
	ResultValue        string `json:"result_value,omitempty"`
	ReferenceRangeText string `json:"reference_range_text,omitempty"`
}

type DependentInvestigationsApiResponse struct {
	Investigation         DependentInvestigationResponse   `json:"investigation,omitempty"`
	DerivedInvestigations []DependentInvestigationResponse `json:"derived_investigations,omitempty"`
}

type DependentInvestigationResponse struct {
	InvestigationId           uint                 `json:"investigation_id,omitempty"`
	DependentInvestigationIds []uint               `json:"dependent_investigation_ids,omitempty"`
	Formula                   string               `json:"formula,omitempty"`
	Decimals                  uint                 `json:"decimals,omitempty"`
	ReferenceRange            MasterReferenceRange `json:"reference_range,omitempty" gorm:"-"`
}
