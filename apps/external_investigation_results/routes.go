package externalInvestigationResults

import (
	"github.com/gin-gonic/gin"

	"github.com/Orange-Health/citadel/apps/external_investigation_results/controller"
	commonConstants "github.com/Orange-Health/citadel/common/constants"
	tokenAuthMiddleware "github.com/Orange-Health/citadel/middlewares/token_auth"
)

func RouteHandler(router *gin.RouterGroup) {
	externalInvestigationResultController := controller.InitExternalInvestigationResultController()

	router.POST("/bulk-upsert", tokenAuthMiddleware.Authenticate(commonConstants.GrootServiceName), externalInvestigationResultController.BulkUpsertInvestigations)
	router.DELETE("/bulk-delete", tokenAuthMiddleware.Authenticate(commonConstants.GrootServiceName), externalInvestigationResultController.BulkDeleteInvestigations)
	router.GET("/", tokenAuthMiddleware.Authenticate(commonConstants.HealthServiceName), externalInvestigationResultController.FetchInvestigations)
}
