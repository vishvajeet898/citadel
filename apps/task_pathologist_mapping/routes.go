package tpm

import (
	"github.com/gin-gonic/gin"

	"github.com/Orange-Health/citadel/apps/task_pathologist_mapping/controller"
)

func RouteHandler(router *gin.RouterGroup) {
	taskPathologistMappingController := controller.InitTaskPathologistMappingController()

	router.POST("/", taskPathologistMappingController.CreateTPM)
}
