package accountsApiClient

import (
	// "context"

	"context"

	"github.com/Orange-Health/citadel/adapters/cache"
	"github.com/Orange-Health/citadel/clients"
	"github.com/Orange-Health/citadel/common/constants"
	"github.com/Orange-Health/citadel/common/utils"
)

type AccountsApiClient struct {
	ApiClient *clients.ApiClient
	Headers   map[string]string
	Cache     cache.CacheLayer
}

func NewClient() *AccountsApiClient {
	apiClient := clients.NewClient()
	apiClient.BaseURL = AccountsApiBaseUrl
	return &AccountsApiClient{
		ApiClient: apiClient,
		Headers: map[string]string{
			"service": constants.CitadelServiceName,
			"api-key": AccountsApiKey,
		},
		Cache: cache.InitializeCache(),
	}
}

func InitializeAccountsApiClient() AccountsApiClientInterface {
	return NewClient()
}

type AccountsApiClientInterface interface {
	CreateFreshDeskTicket(ctx context.Context, groupId uint, creatorGroup, subject, description string)
}

func (accountsApiClient *AccountsApiClient) CreateFreshDeskTicket(ctx context.Context, groupId uint,
	creatorGroup, subject, description string) {
	var response interface{}

	headers := accountsApiClient.Headers
	headers["Content-Type"] = constants.ContentTypeJsonWithCharsetUtf8

	body := map[string]interface{}{
		"subject":       subject,
		"description":   description,
		"group_id":      groupId,
		"creator_group": creatorGroup,
	}

	err := accountsApiClient.ApiClient.Post(ctx, &response, URL_MAP[CREATE_FRESHDESK_TICKET], nil, body, headers, 1, 0)
	if err != nil {
		utils.AddLog(ctx, constants.ERROR_LEVEL, constants.ERROR_WHILE_CREATING_FRESHDESK_TICKET, body, err)
	}
}
