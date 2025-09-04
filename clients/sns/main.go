package snsClient

import (
	"context"
	"net/http"

	"github.com/Orange-Health/citadel/adapters/sentry"
	"github.com/Orange-Health/citadel/common/constants"
	"github.com/Orange-Health/citadel/common/structures"
	"github.com/Orange-Health/citadel/common/utils"
)

type SnsClient struct {
	Sentry sentry.SentryLayer
}

func NewSnsClient() *SnsClient {
	return &SnsClient{
		Sentry: sentry.InitializeSentry(),
	}
}

func InitializeSnsClient() SnsClientInterface {
	return NewSnsClient()
}

type SnsClientInterface interface {
	PublishTo(ctx context.Context, message interface{}, messageAttributes map[string]interface{},
		messageGroupId string, topic string, messageDeduplicationId string) *structures.CommonError
}

func (s *SnsClient) PublishTo(ctx context.Context, message interface{}, messageAttributes map[string]interface{},
	messageGroupId string, topic string, messageDeduplicationId string) *structures.CommonError {
	awsAdapter, err := utils.GetNewPubSubAdapter(constants.SNSAccessKeyID, constants.SNSSecretAccessKey)

	if err != nil {
		s.Sentry.LogError(ctx, constants.ERROR_WHILE_CREATING_AWS_PUB_SUB_ADAPTER, err, nil)
		utils.AddLog(ctx, constants.ERROR_LEVEL, constants.ERROR_WHILE_CREATING_AWS_PUB_SUB_ADAPTER, nil, err)
		return &structures.CommonError{
			Message:    constants.ERROR_WHILE_CREATING_AWS_PUB_SUB_ADAPTER,
			StatusCode: http.StatusInternalServerError,
		}
	}

	loggerAttributes := map[string]interface{}{
		"message":                  message,
		"message_attributes":       messageAttributes,
		"message_group_id":         messageGroupId,
		"topic":                    topic,
		"message_deduplication_id": messageDeduplicationId,
	}
	utils.AddLog(ctx, constants.DEBUG_LEVEL, utils.GetCurrentFunctionName(), loggerAttributes, nil)

	err = awsAdapter.Publish(topic, messageGroupId, messageDeduplicationId, message, messageAttributes)
	if err != nil {
		s.Sentry.LogError(ctx, constants.ERROR_WHILE_PUBLISHING_MESSAGE, err, nil)
		return &structures.CommonError{
			Message:    err.Error(),
			StatusCode: http.StatusInternalServerError,
		}
	}
	utils.AddLog(ctx, constants.DEBUG_LEVEL, constants.MESSAGE_PUBLISHED_SUCCESSFULLY, nil, nil)
	return nil
}
