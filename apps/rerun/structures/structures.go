package structures

import (
	"time"

	commonStructures "github.com/Orange-Health/citadel/common/structures"
)

// @swagger:response RerunDetails
type RerunDetails struct {
	commonStructures.BaseStruct
	// The Test Detail ID
	// example: 1
	TestDetailId uint `json:"test_detail_id"`
	// The Master Investigation ID
	// example: 1
	MasterInvestigationId uint `json:"master_investigation_id"`
	// The Investigation Name
	// example: "Hemoglobin"
	InvestigationName string `json:"investigation_name"`
	// The Investigation Value
	// example: "12.5"
	InvestigationValue string `json:"investigation_value"`
	// The Result Representation Type of Investigation
	// example: "numeric"
	ResultRepresentationType string `json:"result_representation_type"`
	// Lis Code of Investigation
	// example: "OD0112"
	LisCode string `json:"lis_code"`
	// User ID of the User who triggered the Rerun
	// example: 1
	RerunTriggeredBy uint `json:"rerun_triggered_by"`
	// Time at which the Rerun was Triggered
	// example: "2021-07-01T12:00:00Z"
	RerunTriggeredAt *time.Time `json:"rerun_triggered_at"`
	// Reason for Rerun
	// example: "Invalid Value"
	RerunReason string `json:"rerun_reason"`
	// Remarks for Rerun
	// example: "Value was invalid"
	RerunRemarks string `json:"rerun_remarks"`
}
