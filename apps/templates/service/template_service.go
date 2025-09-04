package service

import (
	"net/http"

	mapper "github.com/Orange-Health/citadel/apps/templates/mapper"
	"github.com/Orange-Health/citadel/apps/templates/structures"
	commonConstants "github.com/Orange-Health/citadel/common/constants"
	commonStructures "github.com/Orange-Health/citadel/common/structures"
	commonUtils "github.com/Orange-Health/citadel/common/utils"
)

type TemplateServiceInterface interface {
	GetTemplatesByType(templateType string) (
		[]structures.Template, *commonStructures.CommonError)
}

func (templateService *TemplateService) GetTemplatesByType(templateType string) (
	[]structures.Template, *commonStructures.CommonError) {

	if templateType != "" && !commonUtils.SliceContainsString(commonConstants.TEMPLATE_TYPES, templateType) {
		return []structures.Template{}, &commonStructures.CommonError{
			StatusCode: http.StatusBadRequest,
			Message:    commonConstants.ERROR_INVALID_TEMPLATE_TYPE,
		}
	}

	templates, err := templateService.TemplateDao.GetTemplatesByType(templateType)
	if err != nil {
		return []structures.Template{}, err
	}

	return mapper.MapToTemplateStructs(templates), nil
}
