package investigation_results

import (
	"github.com/gin-gonic/gin"

	"github.com/Orange-Health/citadel/apps/investigation_results/controller"
)

// RouteHandler handles the routes for investigation results.
func RouteHandler(router *gin.RouterGroup) {
	investigationResultController := controller.InitInvestigationResultController()

	router.GET("/tasks/:taskId/details", investigationResultController.GetInvestigationResultByTaskId)
	router.GET("/tasks/:taskId/investigation-abnormality", investigationResultController.GetInvestigationAbnormality)
	router.POST("/tasks/:taskId/modify-value", investigationResultController.GetDerivedInvestigationsAndAbnormality)
	router.GET("/patients/past-records", investigationResultController.GetPatientPastRecords)
	router.GET("/patients/delta-values", investigationResultController.GetDeltaValuesFromPatientId)
}
