package s3wrapperClient

import (
	"context"
	"encoding/json"
	"time"

	"github.com/Orange-Health/citadel/clients"
	"github.com/Orange-Health/citadel/common/constants"
	"github.com/Orange-Health/citadel/common/utils"
)

type S3wrapperClient struct {
	ApiClient *clients.ApiClient
	Headers   map[string]string
}

func NewClient() *S3wrapperClient {
	apiClient := clients.NewClient()
	apiClient.BaseURL = S3wrapperBaseUrl
	return &S3wrapperClient{
		ApiClient: apiClient,
		Headers: map[string]string{
			"service": constants.CitadelServiceName,
			"api-key": S3wrapperApiKey,
		},
	}
}

func InitializeS3wrapperClient() S3wrapperInterface {
	return NewClient()
}

type S3wrapperInterface interface {
	GetTokenizeOrderFilePublicUrl(ctx context.Context, orderFileUrl string) (string, error)
}

func (s3wrapperClient *S3wrapperClient) GetTokenizeOrderFilePublicUrl(ctx context.Context, orderFileUrl string) (string, error) {
	requestBody := PublicUrlRequestBody{FileReference: orderFileUrl}
	responseData := OrderFilePublicUrlData{}
	var response interface{}
	headers := s3wrapperClient.Headers
	headers["Content-Type"] = constants.ContentTypeJsonWithCharsetUtf8

	err := s3wrapperClient.ApiClient.Post(ctx, &response, URL_MAP[ORDER_FILE_PUBLIC_URL], nil, requestBody, headers, 3, time.Millisecond*500)
	if err != nil {
		utils.AddLog(ctx, constants.ERROR_LEVEL, constants.TOKENIZED_URL_ERROR, nil, err)
		return "", err
	}

	responseBytes, err := json.Marshal(response)
	if err != nil {
		utils.AddLog(ctx, constants.ERROR_LEVEL, constants.TOKENIZED_URL_ERROR, nil, err)
		return "", err
	}

	err = json.Unmarshal(responseBytes, &responseData)
	if err != nil {
		utils.AddLog(ctx, constants.ERROR_LEVEL, constants.TOKENIZED_URL_ERROR, nil, err)
		return "", err
	}

	return responseData.MediaUrl, nil
}
