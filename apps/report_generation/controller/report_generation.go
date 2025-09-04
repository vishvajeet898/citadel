package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/Orange-Health/citadel/apps/report_generation/constants"
	commonConstants "github.com/Orange-Health/citadel/common/constants"
	commonStructs "github.com/Orange-Health/citadel/common/structures"
	commonUtils "github.com/Orange-Health/citadel/common/utils"
)

// @Summary		Trigger Report Generation Event
// @Description	Trigger Report Generation Event for a given order ID
// @Tags			Report Generation
// @Produce		json
// @Param			order_id	query		int								true	"Order ID"
// @Param			test_ids	query		array							false	"Test IDs"
// @Param			city_code	query		string							true	"City Code"
// @Success		200			{object}	[]structures.CommonAPIResponse	"Common API Response"
// @Failure		400,500		{object}	[]structures.CommonAPIResponse	"Common API Response"
// @Router			/api/v1/report-generation/trigger [post]
func (a *ReportGenerationEvent) TriggerReportGenerationEvent(c *gin.Context) {
	queryParams := c.Request.URL.Query()
	omsOrderId := queryParams.Get("order_id")
	omsTestIds := commonUtils.ConvertStringToStringSlice(queryParams.Get("test_ids"))
	if omsOrderId == "" {
		commonUtils.HandleErrorResponse(c, http.StatusBadRequest, commonConstants.ERROR_ORDER_ID_NOT_FOUND)
		return
	}

	cErr := a.S.TriggerReportGenerationEvent(c.Request.Context(), omsOrderId, omsTestIds)
	if cErr != nil {
		commonUtils.HandleErrorResponse(c, cErr.StatusCode, cErr.Message)
		return
	}

	c.JSON(http.StatusOK, commonStructs.CommonAPIResponse{
		Message: constants.REPORT_GENERATION_EVENT_TRIGGERED,
	})
}

func (a *ReportGenerationEvent) TriggerOrderApprovedEvent(c *gin.Context) {
	queryParams := c.Request.URL.Query()
	omsOrderId := queryParams.Get("order_id")
	isDummyReport := queryParams.Get("is_dummy_report") == "true"

	if omsOrderId == "" {
		commonUtils.HandleErrorResponse(c, http.StatusBadRequest, commonConstants.ERROR_ORDER_ID_NOT_FOUND)
		return
	}

	payload, cErr := a.S.TriggerOrderApprovedEvent(c.Request.Context(), omsOrderId, isDummyReport)
	if cErr != nil {
		commonUtils.HandleErrorResponse(c, cErr.StatusCode, cErr.Message)
		return
	}

	c.JSON(http.StatusOK, commonStructs.CommonAPIResponse{
		Message: constants.ORDER_APPROVED_EVENT_TRIGGERED,
		Data:    payload,
	})
}
