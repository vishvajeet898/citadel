package patientServiceClient

import (
	"github.com/Orange-Health/citadel/conf"
)

const (
	PATIENT_DETAILS_API  = "PATIENT_DETAILS_API"
	SIMILAR_PATIENTS_API = "SIMILAR_PATIENTS_API"
)

var URL_MAP = map[string]string{
	PATIENT_DETAILS_API:  "/api/v1/patients/",
	SIMILAR_PATIENTS_API: "/api/v1/patients/list/%s",
}

var (
	Config                       = conf.GetConfig()
	PatientServiceApiKey  string = Config.GetString("patient_service.api_key")
	PatientServiceBaseUrl string = Config.GetString("patient_service.base_url")
)
