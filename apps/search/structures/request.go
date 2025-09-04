package structures

type TaskListBasicRequest struct {
	TaskTypes          string `json:"task_types"`
	UserId             uint   `json:"user_id"`
	Limit              uint   `json:"limit"`
	Offset             uint   `json:"offset"`
	PatientName        string `json:"patient_name,omitempty"`
	ContactNumber      string `json:"contact_number,omitempty"`
	PatientId          string `json:"patient_id,omitempty"`
	PartnerName        string `json:"partner_name,omitempty"`
	DoctorName         string `json:"doctor_name,omitempty"`
	OrderId            string `json:"order_id,omitempty"`
	VisitId            string `json:"visit_id,omitempty"`
	RequestId          string `json:"request_id,omitempty"`
	CityCode           string `json:"city_code,omitempty"`
	Status             string `json:"status,omitempty"`
	Department         string `json:"department,omitempty"`
	LabId              string `json:"lab_id,omitempty"`
	SpecialRequirement string `json:"special_requirement,omitempty"`
	OrderType          string `json:"order_type,omitempty"`
}

type TaskListRequest struct {
	TaskTypes []string
	Search    TasksSearch
	Filters   TasksFilters
	Limit     uint
	Offset    uint
	UserId    uint
}

type TasksSearch struct {
	PatientName   string
	ContactNumber string
	PatientId     string
	PartnerName   string
	DoctorName    string
	OrderId       uint
	VisitId       string
	RequestId     uint
	CityCode      string
}

type TasksFilters struct {
	Status             []string
	Department         []string
	LabId              []uint
	SpecialRequirement []string
	OrderType          []string
}

type TaskListDbRequest struct {
	UserId               uint
	Limit                uint
	Offset               uint
	PatientName          string
	PartnerName          string
	DoctorName           string
	PatientDetailsClause map[string]interface{}
	TasksClause          map[string]interface{}
	VisitClause          map[string]interface{}
	TaskMetadataClause   map[string]interface{}
	TestDetailsClause    map[string]interface{}
}
