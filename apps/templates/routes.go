package template

import (
	"github.com/gin-gonic/gin"

	"github.com/Orange-Health/citadel/apps/templates/controller"
)

func RouteHandler(router *gin.RouterGroup) {
	templateController := controller.InitTemplateController()

	router.GET("", templateController.GetTemplatesByType)
}
