package controller

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/Orange-Health/citadel/apps/test_detail/structures"
	commonConstants "github.com/Orange-Health/citadel/common/constants"
	commonStructures "github.com/Orange-Health/citadel/common/structures"
	commonUtils "github.com/Orange-Health/citadel/common/utils"
)

// @Summary		GetTestDetailByTaskId
// @Description	Get Test Details for the given Task ID
// @Tags			test-details
// @Produce		json
// @Param			taskId		path		int						true	"Task ID"
// @Success		200			{object}	[]structures.TestDetail	"Common API Response"
// @Success		400,404,500	{object}	structures.CommonError	"Common API Response"
// @Router			/api/v1/test-details/tasks/{taskId} [get]
func (testDetailController *TestDetail) GetTestDetailByTaskId(c *gin.Context) {
	taskID := commonUtils.ConvertStringToUint(c.Param("taskId"))

	if taskID == 0 {
		commonUtils.HandleErrorResponse(c, http.StatusBadRequest, commonConstants.ERROR_INVALID_TASK_ID)
		return
	}

	testDetails, err := testDetailController.TestDetailService.GetTestDetailsByTaskId(taskID)
	if err != nil {
		commonUtils.HandleErrorResponse(c, err.StatusCode, err.Message)
		return
	}

	c.JSON(http.StatusOK, commonStructures.CommonAPIResponse{
		Message: fmt.Sprintf(commonConstants.FETCHED_SUCCESSFULLY, "Test Detail"),
		Data:    testDetails,
	})
}

func (testDetailController *TestDetail) GetAllBasicTestDetailsByOmsOrderId(c *gin.Context) {
	omsOrderId := c.Param("orderId")

	if omsOrderId == "" {
		commonUtils.HandleErrorResponse(c, http.StatusBadRequest, commonConstants.ERROR_ORDER_ID_CANNOT_BE_ZERO)
		return
	}

	basicTestDetails, cErr := testDetailController.TestDetailService.GetAllBasicTestDetailsByOmsOrderId(c.Request.Context(),
		omsOrderId)
	if cErr != nil {
		commonUtils.HandleErrorResponse(c, cErr.StatusCode, cErr.Message)
		return
	}

	c.JSON(http.StatusOK, structures.TestBasicDetailsByOmsOrderResponse{
		OrderTestDetails: basicTestDetails,
	})
}

func (testDetailController *TestDetail) UpdateProcessingLabForTestDetails(c *gin.Context) {
	var request structures.UpdateProcessingLabRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		commonUtils.HandleErrorResponse(c, http.StatusBadRequest, commonConstants.ERROR_INVALID_PROCESSING_LAB_REQUEST)
		return
	}

	userId, cErr := commonUtils.GetUserIdFromContext(c)
	if cErr != nil {
		commonUtils.HandleErrorResponse(c, cErr.StatusCode, cErr.Message)
		return
	}
	request.UserId = userId

	err := testDetailController.TestDetailService.UpdateProcessingLabForTestDetails(c.Request.Context(), request)
	if err != nil {
		commonUtils.HandleErrorResponse(c, err.StatusCode, err.Message)
		return
	}

	c.JSON(http.StatusOK, commonStructures.CommonAPIResponse{
		Message: fmt.Sprintf(commonConstants.UPDATED_SUCCESSFULLY, "Processing Lab for Test Details"),
	})
}
