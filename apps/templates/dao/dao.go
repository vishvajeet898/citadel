package dao

import (
	commonStructures "github.com/Orange-Health/citadel/common/structures"
	commonUtils "github.com/Orange-Health/citadel/common/utils"
	commonModels "github.com/Orange-Health/citadel/models"
)

type DataLayer interface {
	GetTemplatesByType(templateType string) (
		[]commonModels.Template, *commonStructures.CommonError)
}

func (templateDao *TemplateDao) GetTemplatesByType(templateType string) (
	[]commonModels.Template, *commonStructures.CommonError) {

	templates := []commonModels.Template{}
	query := templateDao.Db
	if templateType != "" {
		query = query.Where("template_type = ?", templateType)
	}
	if err := query.Find(&templates).Error; err != nil {
		return templates, commonUtils.HandleORMError(err)
	}

	return templates, nil
}
