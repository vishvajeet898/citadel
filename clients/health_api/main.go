package healthApiClient

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

type HealthApiClient struct {
	ApiClient *clients.ApiClient
	Headers   map[string]string
	Cache     cache.CacheLayer
}

func NewClient() *HealthApiClient {
	apiClient := clients.NewClient()
	apiClient.BaseURL = HealthApiBaseUrl
	return &HealthApiClient{
		ApiClient: apiClient,
		Headers: map[string]string{
			"service": constants.CitadelServiceName,
			"api-key": HealthApiKey,
		},
		Cache: cache.InitializeCache(),
	}
}

func InitializeHealthApiClient() HealthApiClientInterface {
	return NewClient()
}

type HealthApiClientInterface interface {
	GetDoctorById(ctx context.Context, id uint) (structures.Doctor, error)
	SendGenericSlackMessage(ctx context.Context, channel, message string)
}

func (healthApiClient *HealthApiClient) SendGenericSlackMessage(ctx context.Context, channel, message string) {
	var response interface{}

	headers := healthApiClient.Headers
	headers["Content-Type"] = constants.ContentTypeJsonWithCharsetUtf8

	body := map[string]interface{}{
		"channel": channel,
		"message": message,
	}

	if err := healthApiClient.ApiClient.Post(ctx, &response, URL_MAP[SEND_GENERIC_SLACK_MESSAGE], nil, body, headers,
		SendSlackMessageRetries, SendSlackMessageDelay); err != nil {
		utils.AddLog(ctx, constants.ERROR_LEVEL, constants.ERROR_WHILE_SENDING_GENERIC_SLACK_MESSAGE, body, err)
	}
}

func (healthApiClient *HealthApiClient) GetDoctorById(ctx context.Context, id uint) (structures.Doctor, error) {
	doctor := structures.Doctor{}

	key := fmt.Sprintf(constants.CacheKeyDoctorDetails, id)
	cacheErr := healthApiClient.Cache.Get(ctx, key, &doctor)
	if cacheErr == nil {
		responseBytes, cacheErr := json.Marshal(doctor)
		if cacheErr == nil {
			cacheErr := json.Unmarshal(responseBytes, &doctor)
			if cacheErr == nil {
				return doctor, nil
			}
		}
	}

	if doctor.Id != 0 {
		return doctor, nil
	}

	doctors, err := healthApiClient.GetDoctors(ctx, fmt.Sprint(id))
	if err != nil {
		utils.AddLog(ctx, constants.ERROR_LEVEL, constants.ERROR_WHILE_GETTING_DOCTORS, nil, err)
		return doctor, err
	}
	if len(doctors) == 0 {
		return doctor, errors.New(constants.ERROR_DOCTOR_DOES_NOT_EXIST)
	}
	err = healthApiClient.Cache.Set(ctx, key, doctors[0], constants.CacheExpiry15MinutesInt)
	if err != nil {
		utils.AddLog(ctx, constants.ERROR_LEVEL, constants.ERROR_WHILE_GETTING_DOCTORS, nil, err)
	}
	return doctors[0], nil
}

func (healthApiClient *HealthApiClient) GetDoctors(ctx context.Context, ids string) ([]structures.Doctor, error) {
	var response interface{}
	headers := healthApiClient.Headers
	headers["Content-Type"] = constants.ContentTypeJsonWithCharsetUtf8

	queryParams := map[string]interface{}{
		"ids": ids,
	}

	err := healthApiClient.ApiClient.Get(ctx, &response, URL_MAP[GET_DOCTORS], queryParams, nil, headers,
		GetDoctorRetries, GetDoctorDelay)
	if err != nil {
		utils.AddLog(ctx, constants.ERROR_LEVEL, constants.ERROR_WHILE_GETTING_DOCTORS, nil, err)
		return []structures.Doctor{}, err
	}

	doctorsResponse := DoctorApiResponse{}

	responseBytes, err := json.Marshal(response)
	if err != nil {
		utils.AddLog(ctx, constants.ERROR_LEVEL, constants.ERROR_WHILE_GETTING_DOCTORS, nil, err)
		return nil, err
	}

	err = json.Unmarshal(responseBytes, &doctorsResponse)
	if err != nil {
		utils.AddLog(ctx, constants.ERROR_LEVEL, constants.ERROR_WHILE_GETTING_BULK_PARTNER_DETAILS, nil, err)
		return nil, err
	}

	return doctorsResponse.Result, nil

}
