package partnerApiClient

import (
	"time"

	"github.com/Orange-Health/citadel/conf"
)

var (
	Config            = conf.GetConfig()
	PartnerApiKey     = Config.GetString("partner_api.api_key")
	PartnerApiBaseUrl = Config.GetString("partner_api.base_url")
)

const (
	GET_BULK_PARTNERS = "GET_BULK_PARTNERS"
)

var URL_MAP = map[string]string{
	GET_BULK_PARTNERS: "/v1/internal/partners",
}

const (
	GetBulkPartnerRetries = 3
	GetBulkPartnerDelay   = time.Millisecond * 500
)
