package attachments

import (
	"github.com/gin-gonic/gin"

	"github.com/Orange-Health/citadel/apps/attachments/controller"
)

func RouteHandler(router *gin.RouterGroup) {
	attachmentController := controller.InitAttachmentController()

	router.GET("/tasks/:taskId", attachmentController.GetAttachmentsByTaskId)

	router.POST("/tasks/:taskId", attachmentController.AddAttachmentByTaskId)
}
