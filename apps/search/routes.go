package search

import (
	"github.com/gin-gonic/gin"

	"github.com/Orange-Health/citadel/apps/search/controller"
	commonConstants "github.com/Orange-Health/citadel/common/constants"
	tokenAuthMiddleware "github.com/Orange-Health/citadel/middlewares/token_auth"
)

func RouteHandler(router *gin.RouterGroup) {
	searchController := controller.InitSearchController()

	// Superlab Doctor Tasks
	router.GET("/tasks", searchController.GetTasksList)
	router.GET("/amendment-tasks", searchController.GetAmendmentTasksList)

	// Superlab Info Screen
	router.GET("info-screen/barcode/:barcode", tokenAuthMiddleware.Authenticate(commonConstants.EtsServiceName),
		searchController.GetInfoScreenDataByBarcode)
	router.GET("info-screen/order/:orderId", tokenAuthMiddleware.Authenticate(commonConstants.EtsServiceName),
		searchController.GetInfoScreenDataByOrderId)
	router.GET("info-screen/visit/:visitId", tokenAuthMiddleware.Authenticate(commonConstants.EtsServiceName),
		searchController.GetInfoScreenDataByVisitId)

	router.GET("barcode-details/:barcode",
		tokenAuthMiddleware.MutipleAuthenticate([]string{commonConstants.PorteServiceName, commonConstants.EtsServiceName}),
		searchController.GetBarcodeDetails)
}
