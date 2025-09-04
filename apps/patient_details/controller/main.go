package controller

import (
	"github.com/Orange-Health/citadel/apps/patient_details/service"
)

type PatientDetail struct {
	PatientDetailService service.PatientDetailServiceInterface
}

func InitPatientDetailController() *PatientDetail {
	return &PatientDetail{
		PatientDetailService: service.InitializePatientDetailService(),
	}
}
