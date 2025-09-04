package structures

import (
	commonStructures "github.com/Orange-Health/citadel/common/structures"
)

// @swagger:response Templete
type Template struct {
	commonStructures.BaseStruct
	// The title of the Templete.
	// example: "This is a medical remark"
	Title string `json:"title"`
	// The type of the Templete.
	// example: "medical_remark"
	TemplateType string `json:"template_type"`
	// The description of the Templete.
	// example: "This is a medical remark"
	Description string `json:"description"`
	// The display order of the Templete.
	// example: 1
	DisplayOrder int `json:"display_order"`
}
