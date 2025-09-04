package task

import (
	"github.com/gin-gonic/gin"

	"github.com/Orange-Health/citadel/apps/task/controller"
	commonConstants "github.com/Orange-Health/citadel/common/constants"
	tokenAuthMiddleware "github.com/Orange-Health/citadel/middlewares/token_auth"
)

func RouteHandler(router *gin.RouterGroup) {
	taskController := controller.InitTaskController()

	router.GET("/:taskId/details", taskController.GetTaskDetails)
	router.GET("/count", taskController.GetTasksCount)
	router.GET("/:taskId/calling-details", taskController.GetTaskCallingDetails)

	router.PATCH("/details", taskController.UpdateAllTaskDetails)
	router.PATCH("/attune-sync/:visitId", taskController.TriggerReportGenerationAndAttuneEventByVisitId)
	router.PATCH("/:taskId/undo-release", taskController.UndoReportRelease)
	router.POST("/orders/:orderId", tokenAuthMiddleware.Authenticate(commonConstants.OmsServiceName),
		taskController.CreateUpdateTask)
}
