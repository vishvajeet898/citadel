package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/Orange-Health/citadel/apps/receiving_desk/structures"
	commonConstants "github.com/Orange-Health/citadel/common/constants"
	commonStructures "github.com/Orange-Health/citadel/common/structures"
	commonUtils "github.com/Orange-Health/citadel/common/utils"
)

func (rdController *ReceivingDesk) GetCollectedSamples(c *gin.Context) {
	labId, cErr := commonUtils.GetCurrentLabIdFromContext(c)
	if cErr != nil {
		commonUtils.HandleErrorResponse(c, cErr.StatusCode, cErr.Message)
		return
	}
	collectedSamplesRequest := structures.CollectedSamplesRequest{
		OmsOrderId: c.Query("order_id"),
		Barcode:    c.Query("barcode"),
		TrfId:      c.Query("trf_id"),
		LabId:      labId,
		SearchType: c.Query("search_type"),
	}

	collectedSamples, cErr := rdController.ReceivingDeskService.GetCollectedSamples(c.Request.Context(),
		collectedSamplesRequest)
	if cErr != nil {
		commonUtils.HandleErrorResponse(c, cErr.StatusCode, cErr.Message)
		return
	}

	c.JSON(http.StatusOK, collectedSamples)
}

func (rdController *ReceivingDesk) ReceiveAndSyncSamples(c *gin.Context) {
	userId, cErr := commonUtils.GetUserIdFromContext(c)
	if cErr != nil {
		commonUtils.HandleErrorResponse(c, cErr.StatusCode, cErr.Message)
		return
	}
	labId, cErr := commonUtils.GetCurrentLabIdFromContext(c)
	if cErr != nil {
		commonUtils.HandleErrorResponse(c, cErr.StatusCode, cErr.Message)
		return
	}
	receiveSamplesRequest := structures.ReceiveSamplesRequest{}
	if err := c.ShouldBindJSON(&receiveSamplesRequest); err != nil {
		commonUtils.HandleErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	receiveSamplesRequest.UserId = userId
	receiveSamplesRequest.ReceivingLabId = labId

	omsOrderIds, cErr := rdController.ReceivingDeskService.ReceiveAndSyncSamples(c.Request.Context(), receiveSamplesRequest)
	if cErr != nil {
		commonUtils.HandleErrorResponse(c, cErr.StatusCode, cErr.Message)
		return
	}

	for _, omsOrderId := range omsOrderIds {
		_ = rdController.ReceivingDeskWorkerService.UpdateTaskAndTaskMetadataPostReceiving(c.Request.Context(),
			omsOrderId)
	}

	c.JSON(http.StatusOK, commonStructures.CommonAPIResponse{
		Message: commonConstants.DONE_RESPONSE,
	})
}

func (rdController *ReceivingDesk) MarkAsNotReceived(c *gin.Context) {
	userId, cErr := commonUtils.GetUserIdFromContext(c)
	if cErr != nil {
		commonUtils.HandleErrorResponse(c, cErr.StatusCode, cErr.Message)
		return
	}
	labId, cErr := commonUtils.GetCurrentLabIdFromContext(c)
	if cErr != nil {
		commonUtils.HandleErrorResponse(c, cErr.StatusCode, cErr.Message)
		return
	}

	markAsNotReceivedRequest := commonStructures.MarkAsNotReceivedRequest{}
	if err := c.ShouldBindJSON(&markAsNotReceivedRequest); err != nil {
		commonUtils.HandleErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	markAsNotReceivedRequest.UserId = userId
	markAsNotReceivedRequest.LabId = labId

	cErr = rdController.ReceivingDeskService.MarkAsNotReceived(c.Request.Context(), markAsNotReceivedRequest, true)
	if cErr != nil {
		commonUtils.HandleErrorResponse(c, cErr.StatusCode, cErr.Message)
		return
	}

	c.JSON(http.StatusOK, commonStructures.CommonAPIResponse{
		Message: commonConstants.DONE_RESPONSE,
	})
}
