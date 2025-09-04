package attuneClient

import (
	"fmt"
	"time"

	"github.com/Orange-Health/citadel/conf"
)

const (
	LOGIN                        = "LOGIN"
	GET_PATIENT_RESULT_BY_VISIT  = "GET_PATIENT_RESULT_BY_VISIT"
	INSERT_DOCTOR_DATA_TO_ATTUNE = "INSERT_DOCTOR_DATA_TO_ATTUNE"
	SYNC_DATA_TO_ATTUNE          = "SYNC_DATA_TO_ATTUNE"
)

var ATTUNE_URLS = map[string]string{
	LOGIN:                        "/v1/Authenticate",
	GET_PATIENT_RESULT_BY_VISIT:  "/Orders/GetPatientResultDetailsbyVisitNo",
	INSERT_DOCTOR_DATA_TO_ATTUNE: "/DoctorApprove/InsertDoctorApproveDetails",
	SYNC_DATA_TO_ATTUNE:          "/Orders/pInsertPatientOrderInfo",
}

var (
	config                                 = conf.GetConfig()
	attuneConfigMap      map[string]string = config.GetStringMapString("attune")
	AttuneAccessToken    string            = ""
	AttuneTokenExpiry    int64             = 0
	AttuneTokenCreatedOn int64             = 0

	attuneLoginPayload string = fmt.Sprintf(
		"grant_type=password&username=%s&password=%s&client_id=%s&client_Secret=%s",
		attuneConfigMap["username"], attuneConfigMap["password"], attuneConfigMap["client_id"], attuneConfigMap["client_secret"])
)

const (
	LoginRetries                = 3
	LoginDelay                  = time.Millisecond * 500
	PatientResultDetailsRetries = 3
	PatientResultDetailsDelay   = time.Millisecond * 500
	SyncDataToAttuneRetries     = 1
	SyncDataToAttuneDelay       = time.Millisecond * 0
)
