package controller

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/Orange-Health/citadel/apps/external_investigation_results/constants"
	"github.com/Orange-Health/citadel/apps/external_investigation_results/structures"
	commonConstants "github.com/Orange-Health/citadel/common/constants"
	commonStructures "github.com/Orange-Health/citadel/common/structures"
	commonUtils "github.com/Orange-Health/citadel/common/utils"
)

func (extInvResController *ExternalInvestigationResult) BulkUpsertInvestigations(c *gin.Context) {
	var data structures.UpsertExternalInvestigaitonResultsReqBody
	if err := c.ShouldBindJSON(&data); err != nil {
		commonUtils.HandleErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	investigations, err := extInvResController.ExtInvResService.BulkUpsertInvestigations(&data.Investigations)
	if err != nil {
		commonUtils.HandleErrorResponse(c, err.StatusCode, err.Message)
		return
	}
	c.JSON(http.StatusOK, commonStructures.CommonAPIResponse{
		Message: commonConstants.DONE_RESPONSE,
		Data:    investigations,
	})
}

func (extInvResController *ExternalInvestigationResult) BulkDeleteInvestigations(c *gin.Context) {
	systemExternalInvestigationIdsParam := c.Query("system_external_investigation_ids")
	deletedByParam := c.Query("deleted_by")
	if systemExternalInvestigationIdsParam == "" {
		commonUtils.HandleErrorResponse(c, http.StatusBadRequest, commonConstants.SYSTEM_EXTERNAL_INVESTIGATION_IDS_REQUIRED)
		return
	}
	if deletedByParam == "" {
		commonUtils.HandleErrorResponse(c, http.StatusBadRequest, commonConstants.DELETED_BY_QUERY_REQUIRED)
		return
	}
	systemExternalInvestigationIds := commonUtils.ConvertStringToUintSlice(systemExternalInvestigationIdsParam)
	deletedBy, convErr := strconv.Atoi(deletedByParam)
	if convErr != nil {
		commonUtils.HandleErrorResponse(c, http.StatusBadRequest, convErr.Error())
		return
	}
	err := extInvResController.ExtInvResService.BulkDeleteInvestigations(&systemExternalInvestigationIds, uint(deletedBy))
	if err != nil {
		commonUtils.HandleErrorResponse(c, err.StatusCode, err.Message)
		return
	}
	c.Status(http.StatusNoContent)
}

func (extInvController *ExternalInvestigationResult) FetchInvestigations(c *gin.Context) {
	loincCode := c.Query("loinc_code")
	masterInvestigationMethodMappingId := c.Query("master_investigation_method_mapping_id")
	contactIdStr := c.Query("contact_id")
	limit := c.Query("limit")
	offset := c.Query("offset")
	if contactIdStr == "" {
		commonUtils.HandleErrorResponse(c, http.StatusBadRequest, commonConstants.CONTACT_ID_QUERY_REQUIRED)
		return
	}
	if limit == "" {
		limit = constants.DEFAULT_LIMIT
	}
	if offset == "" {
		offset = constants.DEFAULT_OFFSET
	}
	contactId := commonUtils.ConvertStringToUint(contactIdStr)
	filters := structures.ExternalInvestigationResultsDbFilters{
		Limit:                              commonUtils.ConvertStringToUint(limit),
		Offset:                             commonUtils.ConvertStringToUint(offset),
		ContactId:                          contactId,
		LoincCode:                          loincCode,
		MasterInvestigationMethodMappingId: commonUtils.ConvertStringToUint(masterInvestigationMethodMappingId),
	}
	investigations, err := extInvController.ExtInvResService.FetchInvestigations(filters)
	if err != nil {
		commonUtils.HandleErrorResponse(c, err.StatusCode, err.Message)
		return
	}
	c.JSON(http.StatusOK, commonStructures.CommonAPIResponse{
		Message: commonConstants.OK_RESPONSE,
		Data:    investigations,
	})
}
