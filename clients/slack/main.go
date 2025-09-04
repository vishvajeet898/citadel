package slackClient

import (
	"context"
	"time"

	"github.com/Orange-Health/citadel/adapters/cache"
	"github.com/Orange-Health/citadel/clients"
	"github.com/Orange-Health/citadel/common/constants"
	"github.com/Orange-Health/citadel/common/utils"
)

type SlackClient struct {
	ApiClient *clients.ApiClient
	Headers   map[string]string
	Cache     cache.CacheLayer
}

func NewClient() *SlackClient {
	apiClient := clients.NewClient()
	apiClient.BaseURL = SlackBaseUrl
	return &SlackClient{
		ApiClient: apiClient,
		Headers: map[string]string{
			"service":       constants.CitadelServiceName,
			"Content-Type":  constants.ContentTypeJsonWithCharsetUtf8,
			"Authorization": "Bearer " + SlackToken,
		},
		Cache: cache.InitializeCache(),
	}
}

func InitializeSlackClient() SlackClientInterface {
	return NewClient()
}

type SlackClientInterface interface {
	SendToSlackDirectly(ctx context.Context, channel string, blocks []map[string]interface{}) error
}

func (slackClient *SlackClient) SendToSlackDirectly(ctx context.Context, channel string,
	blocks []map[string]interface{}) error {

	// Override channel for non-prod environments
	if constants.Environment != "prod" {
		// Use a default staging channel if the environment is not prod
		channel = constants.SlackStagingCommunicationChannel
	}

	payload := map[string]interface{}{
		"channel": channel,
		"blocks":  blocks,
	}

	var response interface{}
	err := slackClient.ApiClient.Post(ctx, &response, CHAT_POST_MESSAGE, nil, payload, slackClient.Headers, 3,
		time.Millisecond*500)
	if err != nil {
		utils.AddLog(ctx, constants.ERROR_LEVEL, constants.ERROR_FAILED_TO_SEND_SLACK_MESSAGE, nil, err)
		return err
	}

	return nil
}
