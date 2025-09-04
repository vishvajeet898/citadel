package patient_details

import (
	"github.com/gin-gonic/gin"

	"github.com/Orange-Health/citadel/apps/patient_details/controller"
)

func RouteHandler(router *gin.RouterGroup) {
	patientDetailController := controller.InitPatientDetailController()

	router.GET("/:patientId", patientDetailController.GetPatientDetailsById)

	router.PATCH("/", patientDetailController.UpdatePatientDetails)
}
