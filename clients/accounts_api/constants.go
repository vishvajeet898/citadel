package accountsApiClient

import (
	"github.com/Orange-Health/citadel/conf"
)

var (
	Config             = conf.GetConfig()
	AccountsApiKey     = Config.GetString("accounts_api.api_key")
	AccountsApiBaseUrl = Config.GetString("accounts_api.base_url")
)

const (
	CREATE_FRESHDESK_TICKET = "CREATE_FRESHDESK_TICKET"
)

var URL_MAP = map[string]string{
	CREATE_FRESHDESK_TICKET: "/api/v1/communicator/freshdesk/",
}
