package structures

import (
	patientDetailStruct "github.com/Orange-Health/citadel/apps/patient_details/structures"
	commonStructures "github.com/Orange-Health/citadel/common/structures"
)

// @swagger:response taskDetailResponse
type TaskDetail struct {
	commonStructures.BaseStruct
	// The order ID of the task.
	// example: 54321
	OrderID uint `json:"order_id"`
	// The request ID of the task.
	// example: "REQ12345"
	RequestID string `json:"request_id"`
	// The OMS order ID of the task.
	// example: "OBLR54321"
	OmsOrderId string `json:"oms_order_id"`
	// The lab ID of the task.
	// example: 1
	LabID uint `json:"lab_id"`
	// The city code of the task.
	// example: "BLR"
	CityCode string `json:"city_code"`
	// The status of the task.
	// example: "Completed"
	Status string `json:"status"`
	// The order type of the task.
	// example: "D2C"
	OrderType string `json:"order_type"`
	// Whether the task is active.
	// example: true
	IsActive bool `json:"is_active"`
	// The completion time of the task.
	// example: "2022-01-01T00:00:00Z"
	CompletedAt string `json:"completed_at"`
	// The details of the patient.
	PatientDetails patientDetailStruct.PatientDetail `json:"patient_details"`
	// The Task Metadata.
	TaskMetadata TaskMetadata `json:"task_metadata"`
	// The Task Visit Details.
	TaskVisits []commonStructures.VisitDetailsForTask `json:"task_visits"`
	// Visit Count
	VisitCount int `json:"visit_count"`
}

// @swagger:response taskMetadataResponse
type TaskMetadata struct {
	commonStructures.BaseStruct
	// The ID of the task.
	// example: 1
	TaskID uint `json:"task_id"`
	// Whether the task is a package.
	// example: false
	IsPackage bool `json:"is_package"`
	// Whether the task is morphle.
	// example: false
	IsMorphle bool `json:"is_morphle"`
	// Whether the task is critical.
	// example: false
	IsCritical bool `json:"is_critical"`
	// The name of the doctor.
	// example: "Dr. John Doe"
	DoctorName string `json:"doctor_name"`
	// The number of the doctor.
	// example: "9087654321"
	DoctorNumber string `json:"doctor_number"`
	// The name of the partner.
	// example: "Partner Name"
	PartnerName string `json:"partner_name"`
	// Notes from the doctor.
	// example: "Patient is suffering from fever."
	DoctorNotes string `json:"doctor_notes"`
}

type TaskVisitMapping struct {
	commonStructures.BaseStruct
	// The ID of the task.
	// example: 1
	TaskID uint `json:"task_id"`
	// The ID of the visit.
	// example: 1
	VisitID string `json:"visit_id"`
}

type UpdateAllTaskDetailsStruct struct {
	Task   UpdateTaskStruct `json:"task"`
	UserId uint             `json:"user_id"`
}

type UpdateTaskStruct struct {
	Id          uint                      `json:"id"`
	TestDetails []UpdateTestDetailsStruct `json:"test_details"`
}

type UpdateTestDetailsStruct struct {
	Id             uint                        `json:"id"`
	Investigations []UpdateInvestigationStruct `json:"investigations"`
}

type UpdateInvestigationStruct struct {
	Id                 uint                     `json:"id"`
	Status             string                   `json:"status,omitempty"`
	Abnormality        string                   `json:"abnormality"`
	InvestigationValue string                   `json:"investigation_value,omitempty"`
	InvestigationData  string                   `json:"investigation_data,omitempty"`
	Department         string                   `json:"department,omitempty"`
	CoAuthorizeTo      uint                     `json:"co_authorize_to,omitempty"`
	RerunDetails       UpdateRerunDetailsStruct `json:"rerun_details,omitempty"`
	Remarks            UpdateRemarkStruct       `json:"remarks,omitempty"`
}

type UpdateRerunDetailsStruct struct {
	Id          uint   `json:"id,omitempty"`
	RerunReason string `json:"rerun_reason,omitempty"`
}

type UpdateRemarkStruct struct {
	MedicalRemark  UpdateRemarkStructDetails `json:"medical_remark,omitempty"`
	WithheldReason UpdateRemarkStructDetails `json:"withheld_reason,omitempty"`
}

type UpdateRemarkStructDetails struct {
	Id          uint   `json:"id,omitempty"`
	Description string `json:"description,omitempty"`
	ToDelete    bool   `json:"to_delete,omitempty"`
}

type TaskCallingDetailsResponse struct {
	MobileNumber string `json:"mobile_number"`
	AgentId      string `json:"agent_id"`
}
