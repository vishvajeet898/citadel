package structures

import (
	commonStructures "github.com/Orange-Health/citadel/common/structures"
)

// @swagger:response Remark
type Remark struct {
	commonStructures.BaseStruct
	// The investigation result ID of the remark.
	// example: 1
	InvestigationResultId uint `json:"investigation_result_id"`
	// The description of the remark.
	// example: "This is a remark"
	Description string `json:"description"`
	// The remark type of the remark.
	// example: "RERUN"
	RemarkType string `json:"remark_type"`
	// The remark by of the remark.
	// example: 12
	RemarkBy uint `json:"remark_by"`
}
