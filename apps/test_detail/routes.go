package testDetail

import (
	"github.com/gin-gonic/gin"

	"github.com/Orange-Health/citadel/apps/test_detail/controller"
)

func RouteHandler(router *gin.RouterGroup) {
	testDetailController := controller.InitTestDetailController()

	router.GET("/tasks/:taskId/details", testDetailController.GetTestDetailByTaskId)
	router.GET("/order/:orderId", testDetailController.GetAllBasicTestDetailsByOmsOrderId)

	router.PATCH("/update-processing-lab", testDetailController.UpdateProcessingLabForTestDetails)
}
