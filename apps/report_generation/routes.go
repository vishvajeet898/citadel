package reportgeneration

import (
	"github.com/gin-gonic/gin"

	"github.com/Orange-Health/citadel/apps/report_generation/controller"
	commonConstants "github.com/Orange-Health/citadel/common/constants"
	tokenAuthMiddleware "github.com/Orange-Health/citadel/middlewares/token_auth"
)

func RouteHandler(router *gin.RouterGroup) {
	reportGenerationRouter := controller.InitReportGenerationController()

	router.POST("/trigger", tokenAuthMiddleware.Authenticate(commonConstants.OmsServiceName),
		reportGenerationRouter.TriggerReportGenerationEvent)
	router.POST("/trigger-order-approved-event", tokenAuthMiddleware.Authenticate(commonConstants.OmsServiceName),
		reportGenerationRouter.TriggerOrderApprovedEvent)
}
