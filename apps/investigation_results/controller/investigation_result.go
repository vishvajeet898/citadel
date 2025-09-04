package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/Orange-Health/citadel/apps/investigation_results/structures"
	commonConstants "github.com/Orange-Health/citadel/common/constants"
	commonStructs "github.com/Orange-Health/citadel/common/structures"
	commonUtils "github.com/Orange-Health/citadel/common/utils"
)

// @Summary		Get Investigation Results
// @Description	Get all investigation results by task id
// @Tags			investigation-results
// @Produce		json
// @Param			taskId	path		int									true	"Task ID"
// @Success		200		{object}	[]structures.InvestigationResult	"Investigation Results"
// @Router			/api/v1/investigation-results/tasks/{taskId}/details [get]
func (investigationResultController *InvestigationResult) GetInvestigationResultByTaskId(c *gin.Context) {
	taskId := commonUtils.ConvertStringToUint(c.Param("taskId"))
	if taskId == 0 {
		commonUtils.HandleErrorResponse(c, http.StatusBadRequest, commonConstants.ERROR_INVALID_TASK_ID)
		return
	}

	invResults, cErr := investigationResultController.InvResService.GetInvestigationResultByTaskId(taskId)
	if cErr != nil {
		commonUtils.HandleErrorResponse(c, cErr.StatusCode, cErr.Message)
		return
	}

	c.JSON(http.StatusOK, invResults)
}

// @Summary		Get Patient Past Records
// @Description	Get all investigation results of a patient
// @Tags			investigation-results
// @Produce		json
// @Param			patient_id	path		string										true	"Patient ID"
// @Success		200			{object}	commonStructs.PatientPastRecordsApiResponse	"Patient Past Records"
// @Router			/api/v1/investigation-results/patients/past-records [get]
func (investigationResultController *InvestigationResult) GetPatientPastRecords(c *gin.Context) {
	patientId := c.Query("patient_id")
	if patientId == "" {
		commonUtils.HandleErrorResponse(c, http.StatusBadRequest, commonConstants.ERROR_INVALID_PATIENT_ID)
		return
	}

	invResults, cErr := investigationResultController.InvResService.GetPatientPastRecords(c.Request.Context(), patientId)
	if cErr != nil {
		commonUtils.HandleErrorResponse(c, cErr.StatusCode, cErr.Message)
		return
	}

	c.JSON(http.StatusOK, invResults)
}

// @Summary		Get Delta Values from Patient ID
// @Description	Get delta values from patient id
// @Tags			investigation-results
// @Produce		json
// @Param			patient_id					query		int													true	"Patient Id"
// @Param			master_investigation_ids	query		int													true	"Master Investigation Ids"
// @Param			limit						query		int													true	"Limit"
// @Success		200							{object}	map[uint][]commonStructs.DeltaValuesStructResponse	"Delta Values"
// @Router			/api/v1/investigation-results/patients/delta-values [get]
func (investigationResultController *InvestigationResult) GetDeltaValuesFromPatientId(c *gin.Context) {
	deltaValueRequest := structures.DeltaValuesRequest{
		PatientId:              c.Query("patient_id"),
		MasterInvestigationIds: commonUtils.ConvertStringToUintSlice(c.Query("master_investigation_ids")),
		Limit:                  commonUtils.ConvertStringToUint(c.Query("limit")),
	}

	if deltaValueRequest.PatientId == "" {
		commonUtils.HandleErrorResponse(c, http.StatusBadRequest, commonConstants.ERROR_INVALID_PATIENT_ID)
		return
	}

	if len(deltaValueRequest.MasterInvestigationIds) == 0 {
		commonUtils.HandleErrorResponse(c, http.StatusBadRequest, commonConstants.ERROR_INVALID_MASTER_INVESTIGATION_IDS)
		return
	}

	if deltaValueRequest.Limit == 0 {
		deltaValueRequest.Limit = 5
	}

	deltaValues, cErr := investigationResultController.InvResService.GetDeltaValuesFromPatientId(c.Request.Context(), deltaValueRequest)
	if cErr != nil {
		commonUtils.HandleErrorResponse(c, cErr.StatusCode, cErr.Message)
		return
	}

	c.JSON(http.StatusOK, deltaValues)
}

// @Summary		Get Investigation Abnormality
// @Description	Get investigation abnormality by task id, investigation code and value
// @Tags			investigation-results
// @Produce		json
// @Param			taskId				path		int								true	"Task ID"
// @Param			investigation_code	query		string							true	"Investigation Code"
// @Param			investigation_value	query		string							true	"Investigation Value"
// @Success		200					{object}	structures.CommonAPIResponse	"Abnormality"
// @Router			/api/v1/investigation-results/tasks/{taskId}/investigation-abnormality [get]
func (investigationResultController *InvestigationResult) GetInvestigationAbnormality(c *gin.Context) {
	taskId := commonUtils.ConvertStringToUint(c.Param("taskId"))
	if taskId == 0 {
		commonUtils.HandleErrorResponse(c, http.StatusBadRequest, commonConstants.ERROR_INVALID_TASK_ID)
		return
	}

	queryParams := c.Request.URL.Query()
	investigationCode := queryParams.Get("investigation_code")
	if investigationCode == "" {
		commonUtils.HandleErrorResponse(c, http.StatusBadRequest, commonConstants.ERROR_INVALID_INVESTIGATION_CODE)
		return
	}

	investigationValue := queryParams.Get("investigation_value")
	if investigationValue == "" {
		commonUtils.HandleErrorResponse(c, http.StatusBadRequest, commonConstants.ERROR_INVALID_INVESTIGATION_VALUE)
		return
	}

	abnormality, cErr := investigationResultController.InvResService.GetInvestigationAbnormality(c.Request.Context(), taskId, investigationCode, investigationValue)
	if cErr != nil {
		commonUtils.HandleErrorResponse(c, cErr.StatusCode, cErr.Message)
		return
	}

	c.JSON(http.StatusOK, commonStructs.CommonAPIResponse{
		Data: abnormality,
	})
}

// @Summary		Get Derived Investigations and Abnormality
// @Description	Get derived investigations and abnormality by task id, investigation code and value
// @Tags			investigation-results
// @Produce		json
// @Param			taskId					path		int									true	"Task ID"
// @Param			ModifyValueApiRequest	body		structures.ModifyValueApiRequest	true	"Modify Value API Request"
// @Success		200						{object}	structures.CommonAPIResponse		"Abnormality"
// @Router			/api/v1/investigation-results/tasks/{taskId}/modify-value [post]
func (investigationResultController *InvestigationResult) GetDerivedInvestigationsAndAbnormality(c *gin.Context) {
	taskId := commonUtils.ConvertStringToUint(c.Param("taskId"))
	if taskId == 0 {
		commonUtils.HandleErrorResponse(c, http.StatusBadRequest, commonConstants.ERROR_INVALID_TASK_ID)
		return
	}

	modifyValueApiRequest := structures.ModifyValueApiRequest{}
	if err := c.ShouldBindJSON(&modifyValueApiRequest); err != nil {
		commonUtils.HandleErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	masterInvestigationId := modifyValueApiRequest.CurrentInvestigation.MasterInvestigationId
	if masterInvestigationId == 0 {
		commonUtils.HandleErrorResponse(c, http.StatusBadRequest, commonConstants.ERROR_INVALID_MASTER_INVESTIGATION_ID)
		return
	}

	investigationValue := modifyValueApiRequest.CurrentInvestigation.InvestigationValue
	if investigationValue == "" {
		commonUtils.HandleErrorResponse(c, http.StatusBadRequest, commonConstants.ERROR_INVALID_INVESTIGATION_VALUE)
		return
	}

	modifyValueResponse, cErr := investigationResultController.InvResService.GetDerivedInvestigationsAndAbnormality(
		c.Request.Context(), taskId, masterInvestigationId, investigationValue, modifyValueApiRequest.PastInvestigations)
	if cErr != nil {
		commonUtils.HandleErrorResponse(c, cErr.StatusCode, cErr.Message)
		return
	}

	c.JSON(http.StatusOK, commonStructs.CommonAPIResponse{
		Data: modifyValueResponse,
	})
}
