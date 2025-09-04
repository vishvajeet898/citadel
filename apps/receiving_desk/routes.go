package ReceivingDesks

import (
	"github.com/gin-gonic/gin"

	"github.com/Orange-Health/citadel/apps/receiving_desk/controller"
)

func RouteHandler(router *gin.RouterGroup) {
	receivingDeskController := controller.InitReceivingDeskController()

	router.GET("/collected-samples", receivingDeskController.GetCollectedSamples)
	router.PATCH("/receive-and-sync", receivingDeskController.ReceiveAndSyncSamples)
	router.PATCH("/mark-not-received", receivingDeskController.MarkAsNotReceived)
}
