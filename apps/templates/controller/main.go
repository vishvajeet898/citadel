package controller

import (
	"github.com/Orange-Health/citadel/apps/templates/service"
)

type Template struct {
	TemplateService service.TemplateServiceInterface
}

func InitTemplateController() *Template {
	return &Template{
		TemplateService: service.InitializeRemarkService(),
	}
}
