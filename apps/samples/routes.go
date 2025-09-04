package samples

import (
	"github.com/gin-gonic/gin"

	"github.com/Orange-Health/citadel/apps/samples/controller"
	commonConstants "github.com/Orange-Health/citadel/common/constants"
	tokenAuthMiddleware "github.com/Orange-Health/citadel/middlewares/token_auth"
)

func RouteHandler(router *gin.RouterGroup) {
	sampleController := controller.InitSampleController()

	router.GET("/order-details", tokenAuthMiddleware.Authenticate(commonConstants.OmsServiceName),
		sampleController.GetSampleDetailsByOmsOrderIds)
	router.GET("/order-details/:orderId", tokenAuthMiddleware.Authenticate(commonConstants.OmsServiceName),
		sampleController.GetSampleDetailsByOrderId)
	router.GET("/request-details/:requestId", tokenAuthMiddleware.Authenticate(commonConstants.OmsServiceName),
		sampleController.GetSampleDetailsByRequestId)
	router.GET("/oms-test-details/:visitId", tokenAuthMiddleware.Authenticate(commonConstants.OmsServiceName),
		sampleController.GetOmsTestDetailsByVisitId)
	router.GET("/visit-details/:visitId", sampleController.GetLisSyncDataByVisitId)
	router.GET("/delayed-samples-reverse-logistics", tokenAuthMiddleware.Authenticate(commonConstants.OmsServiceName),
		sampleController.GetSamplesForDelayedReverseLogisticsDashboard)
	router.GET("/srf-order-ids", tokenAuthMiddleware.Authenticate(commonConstants.OmsServiceName),
		sampleController.GetSrfOrderIds)

	router.POST("/remap-samples", tokenAuthMiddleware.Authenticate(commonConstants.OmsServiceName),
		sampleController.RemapSamplesToNewTaskSequence)
	router.POST("/update-task-sequence", tokenAuthMiddleware.Authenticate(commonConstants.OmsServiceName),
		sampleController.UpdateTaskSequenceForSample)
	router.POST("/update-sample-details-for-reschedule", tokenAuthMiddleware.Authenticate(commonConstants.OmsServiceName),
		sampleController.UpdateSampleDetailsForReschedule)
	router.POST("/scheduler/sample-details", tokenAuthMiddleware.Authenticate(commonConstants.OmsServiceName),
		sampleController.GetSampleDetailsForScheduler)
	router.POST("/scheduler/update-sample-details", tokenAuthMiddleware.Authenticate(commonConstants.OmsServiceName),
		sampleController.UpdateSampleDetailsPostTaskCompletion)

	router.PATCH("/collected", sampleController.UpdateSampleCollected)
	router.PATCH("/forcefully-mark-collected", sampleController.ForcefullyMarkSampleAsCollected)
	router.PATCH("/mark-sample-accessioned", tokenAuthMiddleware.Authenticate(commonConstants.OmsServiceName),
		sampleController.CollectionPortalMarkAccessionAsAccessioned)
	router.PATCH("/collect-volume", tokenAuthMiddleware.Authenticate(commonConstants.OmsServiceName),
		sampleController.AddCollectedVolumeToSample)
	router.PATCH("/add-barcode", tokenAuthMiddleware.Authenticate(commonConstants.OmsServiceName),
		sampleController.AddBarcodeDetails)
	router.PATCH("/add-barcode-for-orangers", tokenAuthMiddleware.Authenticate(commonConstants.OmsServiceName),
		sampleController.AddBarcodeDetailsForOrangers)
	router.PATCH("/barcode/:barcode/reject", sampleController.RejectSampleByBarcode)
	router.PATCH("/test/reject-sample-partially", sampleController.RejectSamplePartially)
}
