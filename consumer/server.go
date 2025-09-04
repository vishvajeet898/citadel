package consumer

import (
	"context"
	"encoding/json"
	"sync"

	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	zl "github.com/rs/zerolog/log"

	"github.com/Orange-Health/citadel/adapters/cache"
	"github.com/Orange-Health/citadel/adapters/sentry"
	"github.com/Orange-Health/citadel/common/constants"
	"github.com/Orange-Health/citadel/common/structures"
	"github.com/Orange-Health/citadel/common/utils"
	"github.com/Orange-Health/citadel/worker"
)

type ConsumerServer struct {
	Context context.Context
}

func Start() {
	if gin.IsDebugging() {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	} else {
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}

	consumerServer := &ConsumerServer{
		Context: context.Background(),
	}

	cache.Initialize(consumerServer.Context)
	sentry.Initialize()
	awsAdapter, err := utils.GetNewPubSubAdapter(constants.SQSAccessKeyID, constants.SQSSecretAccessKey)
	if err != nil {
		zl.Error().Err(err).Msg(constants.ERROR_WHILE_CREATING_AWS_PUB_SUB_ADAPTER)
		return
	}

	queueHandlerMapping := map[string]func(*sqs.Message) error{
		constants.SQSQueueURL:         consumerServer.handler,
		constants.StandardSQSQueueURL: consumerServer.rawMessageHandler,
	}

	for {
		var wg sync.WaitGroup
		for _, queueURL := range constants.QueuesToPoll {
			wg.Add(1)
			go func(url string) {
				defer wg.Done()
				err := awsAdapter.PollMessages(url, queueHandlerMapping[url])
				if err != nil {
					utils.AddLog(consumerServer.Context, constants.ERROR_LEVEL, utils.GetCurrentFunctionName(),
						map[string]interface{}{
							"queue": url,
						}, err)
				}
			}(queueURL)
		}
		wg.Wait() // All queues are polled for 20 secs by default so this is fine for now. But if someone wants to poll a queue for a shorter time frame, this will end up waiting for the longest queue's poll time frame. Please keep this in mind if someone is changing this in future.
	}
}

func (consumerServer *ConsumerServer) handler(message *sqs.Message) error {

	body := *message.Body
	messageAttributes := extractMessageAttributes(consumerServer.Context, body)

	event := structures.EventPayload{
		EventType:    messageAttributes["event_type"].(string),
		EventPayload: body,
		TraceID:      messageAttributes["trace_id"].(string),
		Contains:     messageAttributes["contains"].(string),
		GroupId:      *message.Attributes["MessageGroupId"],
	}

	// Optional redis key
	if _, ok := messageAttributes["redis_key"]; ok {
		event.RedisKey = messageAttributes["redis_key"].(string)
	}

	return worker.SendEventHandlerToWorker(consumerServer.Context, event)
}

func (consumerServer *ConsumerServer) rawMessageHandler(message *sqs.Message) error {
	event := structures.EventPayload{
		EventType:    *message.MessageAttributes["event_type"].StringValue,
		EventPayload: *message.Body,
		TraceID:      *message.MessageAttributes["trace_id"].StringValue,
		Contains:     *message.MessageAttributes["contains"].StringValue,
	}

	// Optional redis key
	if _, ok := message.MessageAttributes["redis_key"]; ok {
		event.RedisKey = *message.MessageAttributes["redis_key"].StringValue
	}
	return worker.SendEventHandlerToWorker(consumerServer.Context, event)
}

func extractMessageAttributes(ctx context.Context, body string) map[string]interface{} {
	// deserialize message into map[string]interface{}
	messageMap := make(map[string]interface{})
	sentryInstance := sentry.InitializeSentry()
	err := json.Unmarshal([]byte(body), &messageMap)
	if err != nil {
		sentryInstance.LogError(ctx, "Error while unmarshalling message body", err, nil)
		return nil
	}

	messageAttributes, ok := messageMap["MessageAttributes"].(map[string]interface{})
	if !ok {
		sentryInstance.LogError(ctx, "MessageAttributes missing or not a map", nil, map[string]interface{}{
			"message": body,
		})
		return nil
	}

	messageAttributesMap := make(map[string]interface{})
	for key, value := range messageAttributes {
		valueMap := value.(map[string]interface{})
		valueValue, ok := valueMap["Value"].(string)
		if !ok {
			sentryInstance.LogError(ctx, "Value missing or not a string", nil, map[string]interface{}{
				"message": body,
			})
			continue
		}
		messageAttributesMap[key] = valueValue
	}

	return messageAttributesMap
}
