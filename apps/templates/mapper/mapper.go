package mapper

import (
	"github.com/Orange-Health/citadel/apps/templates/structures"
	commonModels "github.com/Orange-Health/citadel/models"
)

func MapToTemplateStruct(template commonModels.Template) structures.Template {
	var templateStruct structures.Template
	templateStruct.ID = template.Id
	templateStruct.Title = template.Title
	templateStruct.TemplateType = template.TemplateType
	templateStruct.Description = template.Description
	templateStruct.DisplayOrder = template.DisplayOrder
	return templateStruct
}

func MapToTemplateStructs(templates []commonModels.Template) []structures.Template {
	var templateStructs []structures.Template
	for _, template := range templates {
		templateStructs = append(templateStructs, MapToTemplateStruct(template))
	}

	return templateStructs
}
