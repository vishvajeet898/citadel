package task

import (
	"github.com/gin-gonic/gin"

	"github.com/Orange-Health/citadel/apps/task_metadata/controller"
)

func RouteHandler(router *gin.RouterGroup) {
	taskMetadataController := controller.InitTaskMetadataController()

	router.GET("/tasks/:taskId/details", taskMetadataController.GetTaskMetadataDetails)
	router.PATCH("/tasks/:taskId", taskMetadataController.UpdateTaskMetadata)
	router.PATCH("/task/:taskId/last-event-sent-at", taskMetadataController.UpdateLastEventSentAt)
}
