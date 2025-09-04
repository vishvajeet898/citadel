package partnerApiClient

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/Orange-Health/citadel/adapters/cache"
	"github.com/Orange-Health/citadel/clients"
	"github.com/Orange-Health/citadel/common/constants"
	"github.com/Orange-Health/citadel/common/structures"
	"github.com/Orange-Health/citadel/common/utils"
)

type PartnerApiClient struct {
	ApiClient *clients.ApiClient
	Headers   map[string]string
	Cache     cache.CacheLayer
}

func NewClient() *PartnerApiClient {
	apiClient := clients.NewClient()
	apiClient.BaseURL = PartnerApiBaseUrl
	return &PartnerApiClient{
		ApiClient: apiClient,
		Headers: map[string]string{
			"service":         constants.CitadelServiceName,
			"partner-api-key": PartnerApiKey,
		},
		Cache: cache.InitializeCache(),
	}
}

func InitializePartnerApiClient() PartnerApiClientInterface {
	return NewClient()
}

type PartnerApiClientInterface interface {
	GetPartnerById(ctx context.Context, partnerId uint) (structures.Partner, error)
}

func (partnerApiClient *PartnerApiClient) GetPartnerById(ctx context.Context, partnerId uint) (structures.Partner, error) {
	var cacheResponse interface{}
	key := fmt.Sprintf(constants.CacheKeyPartnerDetails, partnerId)
	partner := structures.Partner{}
	cacheErr := partnerApiClient.Cache.Get(ctx, key, &cacheResponse)
	if cacheErr == nil {
		responseBytes, cacheErr := json.Marshal(cacheResponse)
		if cacheErr == nil {
			cacheErr := json.Unmarshal(responseBytes, &partner)
			if cacheErr == nil {
				return partner, nil
			}
		}
	}

	if partner.Id != 0 {
		return partner, nil
	}

	filters := map[string]interface{}{
		"id": fmt.Sprintf("%d", partnerId),
	}

	partners, err := partnerApiClient.GetBulkPartners(ctx, filters)
	if err != nil {
		return structures.Partner{}, err
	}
	if len(partners) == 0 {
		return structures.Partner{}, errors.New(constants.ERROR_PARTNER_DOES_NOT_EXIST)
	}

	err = partnerApiClient.Cache.Set(ctx, key, partners[0], constants.CacheExpiry15MinutesInt)
	if err != nil {
		utils.AddLog(ctx, constants.ERROR_LEVEL, constants.ERROR_WHILE_GETTING_PARTNER_DETAILS, nil, err)
	}

	return partners[0], err
}

func (partnerApiClient *PartnerApiClient) GetBulkPartners(ctx context.Context, filters map[string]interface{}) (
	[]structures.Partner, error) {
	var response interface{}
	partners := GetPartnersApiResponse{}
	headers := partnerApiClient.Headers
	err := partnerApiClient.ApiClient.Get(ctx, &response, URL_MAP[GET_BULK_PARTNERS], filters, nil, headers,
		GetBulkPartnerRetries, GetBulkPartnerDelay)
	if err != nil {
		utils.AddLog(ctx, constants.ERROR_LEVEL, constants.ERROR_WHILE_GETTING_BULK_PARTNER_DETAILS, nil, err)
		return []structures.Partner{}, err
	}

	responseBytes, err := json.Marshal(response)
	if err != nil {
		utils.AddLog(ctx, constants.ERROR_LEVEL, constants.ERROR_WHILE_GETTING_BULK_PARTNER_DETAILS, nil, err)
		return nil, err
	}

	err = json.Unmarshal(responseBytes, &partners)
	if err != nil {
		utils.AddLog(ctx, constants.ERROR_LEVEL, constants.ERROR_WHILE_GETTING_BULK_PARTNER_DETAILS, nil, err)
		return nil, err
	}

	return partners.Data, nil
}
