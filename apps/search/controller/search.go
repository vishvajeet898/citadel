package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/Orange-Health/citadel/apps/search/constants"
	"github.com/Orange-Health/citadel/apps/search/structures"
	commonUtils "github.com/Orange-Health/citadel/common/utils"
)

// @Summary		Get Tasks List
// @Description	Get all tasks based on filters provided
// @Tags			searchController
// @Produce		json
// @Success		200			{object}	structures.TaskListResponse		"Tasks List"
// @Failure		400,404,500	{object}	structures.CommonAPIResponse	"Error Response"
// @Router			/api/v1/searchController/tasks [get]
func (searchController *Search) GetTasksList(c *gin.Context) {
	userId, cErr := commonUtils.GetUserIdFromContext(c)
	if cErr != nil {
		commonUtils.HandleErrorResponse(c, cErr.StatusCode, cErr.Message)
		return
	}

	taskListBasicRequest := structures.TaskListBasicRequest{
		UserId:             userId,
		Limit:              commonUtils.ConvertStringToUint(c.Query("limit")),
		Offset:             commonUtils.ConvertStringToUint(c.Query("offset")),
		TaskTypes:          c.Query("task_types"),
		PatientName:        c.Query("patient_name"),
		ContactNumber:      c.Query("contact_number"),
		PatientId:          c.Query("patient_id"),
		PartnerName:        c.Query("partner_name"),
		DoctorName:         c.Query("doctor_name"),
		OrderId:            c.Query("order_id"),
		VisitId:            c.Query("visit_id"),
		RequestId:          c.Query("request_id"),
		Status:             c.Query("status"),
		Department:         c.Query("department"),
		LabId:              c.Query("lab_id"),
		SpecialRequirement: c.Query("special_requirement"),
		OrderType:          c.Query("order_type"),
	}

	response, cErr := searchController.SearchService.GetTasksList(c.Request.Context(), taskListBasicRequest)
	if cErr != nil {
		commonUtils.HandleErrorResponse(c, cErr.StatusCode, cErr.Message)
		return
	}

	c.JSON(http.StatusOK, response)
}

func (searchController *Search) GetAmendmentTasksList(c *gin.Context) {
	userId, cErr := commonUtils.GetUserIdFromContext(c)
	if cErr != nil {
		commonUtils.HandleErrorResponse(c, cErr.StatusCode, cErr.Message)
		return
	}

	taskListBasicRequest := structures.TaskListBasicRequest{
		UserId:        userId,
		Limit:         commonUtils.ConvertStringToUint(c.Query("limit")),
		Offset:        commonUtils.ConvertStringToUint(c.Query("offset")),
		TaskTypes:     constants.TASK_TYPE_AMENDMENT,
		PatientName:   c.Query("patient_name"),
		ContactNumber: c.Query("contact_number"),
		PatientId:     c.Query("patient_id"),
		PartnerName:   c.Query("partner_name"),
		DoctorName:    c.Query("doctor_name"),
		OrderId:       c.Query("order_id"),
		VisitId:       c.Query("visit_id"),
		RequestId:     c.Query("request_id"),
	}

	response, cErr := searchController.SearchService.GetAmendmentTasksList(c.Request.Context(), taskListBasicRequest)
	if cErr != nil {
		commonUtils.HandleErrorResponse(c, cErr.StatusCode, cErr.Message)
		return
	}

	c.JSON(http.StatusOK, response)
}

func (searchController *Search) GetBarcodeDetails(c *gin.Context) {
	barcode := c.Params.ByName("barcode")
	serviceType := c.Query("serviceType")
	labId := commonUtils.ConvertStringToUint(c.Query("lab_id"))

	response, cErr := searchController.SearchService.GetBarcodeDetails(barcode, serviceType, labId)
	if cErr != nil {
		commonUtils.HandleErrorResponse(c, cErr.StatusCode, cErr.Message)
		return
	}

	c.JSON(http.StatusOK, response)
}

func (searchController *Search) GetInfoScreenDataByBarcode(c *gin.Context) {
	barcode := c.Params.ByName("barcode")
	if barcode == "" {
		commonUtils.HandleErrorResponse(c, http.StatusBadRequest, constants.ERROR_INVALID_BARCODE)
		return
	}

	labId := c.Query("lab_id")
	labIdUint := commonUtils.ConvertStringToUint(labId)

	response, cErr := searchController.SearchService.GetInfoScreenDataByBarcode(c.Request.Context(), barcode, labIdUint)
	if cErr != nil {
		commonUtils.HandleErrorResponse(c, cErr.StatusCode, cErr.Message)
		return
	}

	c.JSON(http.StatusOK, response)
}

func (searchController *Search) GetInfoScreenDataByOrderId(c *gin.Context) {
	omsOrderId := c.Params.ByName("orderId")
	if omsOrderId == "" {
		commonUtils.HandleErrorResponse(c, http.StatusBadRequest, constants.ERROR_INVALID_ORDER_ID)
		return
	}

	labId := c.Query("lab_id")
	labIdUint := commonUtils.ConvertStringToUint(labId)

	response, cErr := searchController.SearchService.GetInfoScreenDataByOrderId(c.Request.Context(), omsOrderId, labIdUint)
	if cErr != nil {
		commonUtils.HandleErrorResponse(c, cErr.StatusCode, cErr.Message)
		return
	}

	c.JSON(http.StatusOK, response)
}

func (searchController *Search) GetInfoScreenDataByVisitId(c *gin.Context) {
	visitId := c.Params.ByName("visitId")
	if visitId == "" {
		commonUtils.HandleErrorResponse(c, http.StatusBadRequest, constants.ERROR_INVALID_VISIT_ID)
		return
	}

	labId := c.Query("lab_id")
	labIdUint := commonUtils.ConvertStringToUint(labId)

	response, cErr := searchController.SearchService.GetInfoScreenDataByVisitId(c.Request.Context(), visitId, labIdUint)
	if cErr != nil {
		commonUtils.HandleErrorResponse(c, cErr.StatusCode, cErr.Message)
		return
	}

	c.JSON(http.StatusOK, response)
}
