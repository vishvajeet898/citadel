package healthApiClient

import (
	"time"

	"github.com/Orange-Health/citadel/conf"
)

var (
	Config           = conf.GetConfig()
	HealthApiKey     = Config.GetString("health_api.api_key")
	HealthApiBaseUrl = Config.GetString("health_api.base_url")
)

const (
	GET_DOCTORS                = "GET_DOCTORS"
	SEND_GENERIC_SLACK_MESSAGE = "SEND_GENERIC_SLACK_MESSAGE"
)

var URL_MAP = map[string]string{
	GET_DOCTORS:                "/api/v1/doctor/list",
	SEND_GENERIC_SLACK_MESSAGE: "/api/v1/ops/communication/slack/",
}

const (
	SendSlackMessageRetries = 1
	SendSlackMessageDelay   = 0
	GetDoctorRetries        = 3
	GetDoctorDelay          = time.Millisecond * 500
)
