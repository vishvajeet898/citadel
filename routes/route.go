package routes

import (
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	attachments "github.com/Orange-Health/citadel/apps/attachments"
	auditLog "github.com/Orange-Health/citadel/apps/audit_log"
	externalInvestigationResults "github.com/Orange-Health/citadel/apps/external_investigation_results"
	health "github.com/Orange-Health/citadel/apps/health"
	investigationResults "github.com/Orange-Health/citadel/apps/investigation_results"
	patientDetails "github.com/Orange-Health/citadel/apps/patient_details"
	receivingDesk "github.com/Orange-Health/citadel/apps/receiving_desk"
	reportGeneration "github.com/Orange-Health/citadel/apps/report_generation"
	samples "github.com/Orange-Health/citadel/apps/samples"
	search "github.com/Orange-Health/citadel/apps/search"
	task "github.com/Orange-Health/citadel/apps/task"
	taskMetaData "github.com/Orange-Health/citadel/apps/task_metadata"
	taskPathologistMapping "github.com/Orange-Health/citadel/apps/task_pathologist_mapping"
	template "github.com/Orange-Health/citadel/apps/templates"
	testDetail "github.com/Orange-Health/citadel/apps/test_detail"
	users "github.com/Orange-Health/citadel/apps/users"
	_ "github.com/Orange-Health/citadel/docs"
)

// @title			Citadel API
// @description	Citadel API Documentation
func initialiseRoutes(router *gin.Engine) {
	health.RouteHandler(router.Group("/ping"))
	task.RouteHandler(router.Group("/api/v1/tasks"))
	testDetail.RouteHandler(router.Group("/api/v1/test-details"))
	taskPathologistMapping.RouteHandler(router.Group("/api/v1/task-pathologist-mapping"))
	attachments.RouteHandler(router.Group("/api/v1/attachments"))
	investigationResults.RouteHandler(router.Group("/api/v1/investigation-results"))
	template.RouteHandler(router.Group("/api/v1/templates"))
	auditLog.RouteHandler(router.Group("/api/v1/audit-logs"))
	patientDetails.RouteHandler(router.Group("/api/v1/patient-details"))
	users.RouteHandler(router.Group("/api/v1/users"))
	search.RouteHandler(router.Group("/api/v1/search"))
	reportGeneration.RouteHandler(router.Group("/api/v1/report-generation"))
	taskMetaData.RouteHandler(router.Group("/api/v1/task-metadata"))
	receivingDesk.RouteHandler(router.Group("/api/v1/receiving-desk"))
	externalInvestigationResults.RouteHandler(router.Group("/api/v1/external-investigation-results"))
	samples.RouteHandler(router.Group("/api/v1/samples"))

	if gin.IsDebugging() {
		router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	}
}
