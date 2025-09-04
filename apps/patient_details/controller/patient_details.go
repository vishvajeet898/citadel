package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/Orange-Health/citadel/apps/patient_details/structures"
	commonConstants "github.com/Orange-Health/citadel/common/constants"
	commonUtils "github.com/Orange-Health/citadel/common/utils"
)

func (patientDetailController *PatientDetail) GetPatientDetailsById(c *gin.Context) {
	patientId := commonUtils.ConvertStringToUint(c.Param("patientId"))
	if patientId == 0 {
		commonUtils.HandleErrorResponse(c, http.StatusBadRequest, commonConstants.ERROR_INVALID_PATIENT_ID)
		return
	}

	patientDetails, cErr := patientDetailController.PatientDetailService.GetPatientDetailsResponseById(patientId)
	if cErr != nil {
		commonUtils.HandleErrorResponse(c, cErr.StatusCode, cErr.Message)
		return
	}

	c.JSON(http.StatusOK, patientDetails)
}

func (patientDetailsController *PatientDetail) UpdatePatientDetails(c *gin.Context) {
	patientDetailsRequest := structures.PatientDetailsUpdateRequest{}

	if err := c.ShouldBindJSON(&patientDetailsRequest); err != nil {
		commonUtils.HandleErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	patientDetails, cErr := patientDetailsController.PatientDetailService.UpdatePatientDetails(c.Request.Context(), patientDetailsRequest)
	if cErr != nil {
		commonUtils.HandleErrorResponse(c, cErr.StatusCode, cErr.Message)
		return
	}

	c.JSON(http.StatusOK, patientDetails)
}
