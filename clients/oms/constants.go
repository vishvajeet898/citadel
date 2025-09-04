package omsClient

import (
	"github.com/Orange-Health/citadel/conf"
)

var (
	Config     = conf.GetConfig()
	OmsBaseUrl = Config.GetString("oms.base_url")
	OmsApiKey  = Config.GetString("oms.api_key")
)

const (
	REPORT_GENERATION_EVENT_DETAILS       = "REPORT_GENERATION_EVENT_DETAILS"
	UPDATE_PATIENT_DETAILS                = "UPDATE_PATIENT_DETAILS"
	DELTA_VALUES_FROM_PATIENT_IDS         = "DELTA_VALUES_FROM_PATIENT_IDS"
	PATIENT_PAST_RECORDS_FROM_PATIENT_IDS = "PATIENT_PAST_RECORDS_FROM_PATIENT_IDS"
	GET_ORDER_BY_ID                       = "GET_ORDER_BY_ID"
)

var URL_MAP = map[string]string{
	REPORT_GENERATION_EVENT_DETAILS:       "/api/citadel/order/%d/report-generation-event-details",
	UPDATE_PATIENT_DETAILS:                "/api/citadel/patient-details/%d",
	DELTA_VALUES_FROM_PATIENT_IDS:         "/api/citadel/patients/delta-values",
	PATIENT_PAST_RECORDS_FROM_PATIENT_IDS: "/api/citadel/patients/past-records",
	GET_ORDER_BY_ID:                       "/api/citadel/order/%s",
}
