package workerTasks

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/RichardKnop/machinery/v1/tasks"

	cacheClient "github.com/Orange-Health/citadel/adapters/cache"
	"github.com/Orange-Health/citadel/adapters/sentry"
	"github.com/Orange-Health/citadel/common/constants"
	"github.com/Orange-Health/citadel/common/structures"
	"github.com/Orange-Health/citadel/common/utils"
	consumerTasks "github.com/Orange-Health/citadel/consumer_tasks"
)

func (wt *WorkerTaskService) EventHandlerTask(eventType, traceId, contains, eventKey string) error {
	ctx := context.Background()
	startTime := time.Now()

	var eventPayloadInterface interface{}
	err := wt.Cache.Get(ctx, eventKey, &eventPayloadInterface)
	if err != nil {
		utils.AddLog(ctx, constants.ERROR_LEVEL, utils.GetCurrentFunctionName(), nil, err)
		return err
	}
	eventPayload := eventPayloadInterface.(string)
	loggingAttributes := map[string]interface{}{
		"event_type":    eventType,
		"event_payload": eventPayload,
		"trace_id":      traceId,
		"contains":      contains,
	}
	utils.AddLog(ctx, constants.DEBUG_LEVEL, utils.GetCurrentFunctionName(), loggingAttributes, nil)

	sentryInstance := sentry.InitializeSentry()

	eventProcessor := consumerTasks.EventProcessor{
		Db:                          wt.Db,
		Cache:                       wt.Cache,
		Sentry:                      wt.Sentry,
		TaskService:                 wt.TaskService,
		SampleService:               wt.SampleService,
		ReceivingDeskService:        wt.ReceivingDeskService,
		TestSampleMappingService:    wt.TestSampleMappingService,
		AttuneService:               wt.AttuneService,
		OrderDetailsService:         wt.OrderDetailsService,
		TaskPathMappingService:      wt.TaskPathMappingService,
		PatientDetailService:        wt.PatientDetailService,
		TestDetailService:           wt.TestDetailService,
		InvestigationResultsService: wt.InvestigationResultsService,
		RerunService:                wt.RerunService,
		RemarkService:               wt.RemarkService,
		ReportGenerationService:     wt.ReportGenerationService,
		UserService:                 wt.UserService,
		AbnormalityService:          wt.AbnormalityService,
		AttachmentsService:          wt.AttachmentsService,
		ContactService:              wt.ContactService,
		EtsService:                  wt.EtsService,
		CdsService:                  wt.CdsService,
		PubsubService:               wt.PubsubService,
		AttuneClient:                wt.AttuneClient,
		CdsClient:                   wt.CdsClient,
		ReportRebrandingClient:      wt.ReportRebrandingClient,
		S3wrapperClient:             wt.S3wrapperClient,
		S3Client:                    wt.S3Client,
		SnsClient:                   wt.SnsClient,
		PartnerApiClient:            wt.PartnerApiClient,
		HealthApiClient:             wt.HealthApiClient,
		SlackClient:                 wt.SlackClient,
		CommonTaskProcessor:         wt.CommonTaskProcessor,
	}

	switch eventType {
	case constants.OmsAttachmentEvent:
		err = eventProcessor.OmsAttachmentEventTask(ctx, eventPayload)
	case constants.OmsManualReportUploadEvent:
		err = eventProcessor.OmsManualReportUploadEventTask(ctx, eventPayload)
	case constants.OmsTestDeleteEvent:
		err = eventProcessor.OmsTestDeleteEventTask(ctx, eventPayload)
	case constants.OmsOrderCreatedEvent:
		err = eventProcessor.OmsOrderCreateUpdateEventTask(ctx, eventPayload)
	case constants.OmsOrderUpdatedEvent:
		err = eventProcessor.OmsOrderCreateUpdateEventTask(ctx, eventPayload)
	case constants.OmsOrderCompletedEvent:
		err = eventProcessor.OmsOrderCompletedEventTask(ctx, eventPayload)
	case constants.OmsOrderCancelledEvent:
		err = eventProcessor.OmsOrderCancelledEventTask(ctx, eventPayload)
	case constants.MergeContactEvent:
		err = eventProcessor.MergeContactEventTask(ctx, eventPayload)
	case constants.SampleCollectedEvent:
		err = eventProcessor.SampleCollectedEventTask(ctx, eventPayload)
	case constants.LisEvent:
		err = eventProcessor.LisEventTask(ctx, eventPayload)
	case constants.OrderReportPdfReadyEvent:
		err = eventProcessor.OrderReportPdfReadyEventTask(ctx, eventPayload)
	case constants.SampleRecollectionEvent:
		err = eventProcessor.SampleRecollectionEventTask(ctx, eventPayload)
	case constants.TestDetailsEvent:
		err = eventProcessor.TestDetailsEventTask(ctx, eventPayload)
	case constants.MarkSampleSnrEvent:
		err = eventProcessor.MarkSampleSnrEventTask(ctx, eventPayload)
	case constants.UpdateTaskSequence:
		err = eventProcessor.UpdateTaskSequenceEventTask(ctx, eventPayload)
	case constants.UpdateSrfIdToLisEvent:
		err = eventProcessor.UpdateSrfIdToLisEventTask(ctx, eventPayload)
	default:
		sentryInstance.LogError(ctx, "Unknown event type", nil, loggingAttributes)
	}

	if err != nil {
		utils.AddLog(ctx, constants.DEBUG_LEVEL, utils.GetCurrentFunctionName(), map[string]interface{}{
			"event_type":    eventType,
			"event_payload": eventPayload,
			"trace_id":      traceId,
			"contains":      contains,
			"redis_key":     eventKey,
		}, err)
		return err
	}

	_ = wt.Cache.Delete(ctx, eventKey)
	utils.AddLog(ctx, constants.INFO_LEVEL, utils.GetCurrentFunctionName(),
		map[string]interface{}{
			"total_time_taken": time.Since(startTime).String(),
			"event_type":       eventType,
		}, nil)

	return nil
}

func getGroupIdByEventType(eventType, groupId string) string {
	if utils.SliceContainsString(constants.VariableGroupIdEventNames, eventType) {
		finalGroupId := constants.EventNameGroupIdMap[eventType]
		if finalGroupId == "" {
			return fmt.Sprintf("%s:%v", groupId, time.Now().Unix())
		}
		return fmt.Sprintf(finalGroupId, groupId)
	} else if groupId == "" {
		return fmt.Sprintf("%s:%v", constants.EventNameGroupIdMap[eventType], time.Now().Unix())
	}
	return fmt.Sprintf("%s:%v", groupId, time.Now().Unix())
}

func CreateEventHandlerSignature(ctx context.Context, event structures.EventPayload) (*tasks.Signature, error) {
	sentryInstance := sentry.InitializeSentry()
	var eventPayload string
	loggingAttributes := map[string]interface{}{
		"event_type":    event.EventType,
		"trace_id":      event.TraceID,
		"contains":      event.Contains,
		"redis_key":     event.RedisKey,
		"event_payload": event.EventPayload,
		"group_id":      event.GroupId,
	}
	// If redis key exists, remove the body passed in the event payload and fetch the payload from redis in the task
	if event.RedisKey != "" {
		event.EventPayload = ""

		awsAdapter, err := utils.GetNewPubSubAdapter(constants.SNSAccessKeyID, constants.SNSSecretAccessKey)
		if err != nil {
			utils.AddLog(ctx, constants.DEBUG_LEVEL, utils.GetCurrentFunctionName(), nil, err)
			return nil, err
		}

		eventPayload, err = awsAdapter.FetchValueFromRedis(event.RedisKey)
		if err != nil {
			sentryInstance.LogError(ctx, constants.ERROR_WHILE_FETCHING_PAYLOAD_FROM_REDIS, err, loggingAttributes)
			return nil, err
		}
	} else {
		eventPayload = event.EventPayload

		var eventPayloadMap map[string]interface{}
		err := json.Unmarshal([]byte(eventPayload), &eventPayloadMap)
		if err != nil {
			sentryInstance.LogError(ctx, constants.ERROR_WHILE_PARSING_EVENT_PAYLOAD, err, loggingAttributes)
			return nil, err
		}

		eventPayload = eventPayloadMap["Message"].(string)
		if eventPayload == "" {
			sentryInstance.LogError(ctx, constants.ERROR_EMPTY_MESSAGE_IN_EVENT_PAYLOAD, nil, loggingAttributes)
			return nil, errors.New(constants.ERROR_EMPTY_MESSAGE_IN_EVENT_PAYLOAD)
		}
	}

	cache := cacheClient.InitializeCache()
	epoch := time.Now().UnixNano()
	eventKey := fmt.Sprintf("%s_%d", event.EventType, epoch)
	err := cache.Set(ctx, eventKey, eventPayload, constants.CacheExpiry1HourInt)
	if err != nil {
		utils.AddLog(ctx, constants.ERROR_LEVEL, utils.GetCurrentFunctionName(), nil, err)
		return nil, err
	}

	retryTimeout := constants.DefaultRetryTimout
	retryCount := constants.DefaultRetryCount
	switch event.EventType {
	case constants.OmsOrderCreatedEvent, constants.OmsOrderUpdatedEvent, constants.OmsTestDeleteEvent:
		retryTimeout = constants.OmsOrderCreateUpdateRetryTimeout
	case constants.SampleCollectedEvent:
		retryTimeout = constants.OmsCollectionRetryTimeout
		retryCount = constants.OmsCollectionRetryCount
	}

	scheduleTask := &tasks.Signature{
		Name: constants.ConsumerEvents,
		Args: []tasks.Arg{
			{
				Type:  "string",
				Value: event.EventType,
				Name:  "event_type",
			},
			{
				Type:  "string",
				Value: event.TraceID,
				Name:  "trace_id",
			},
			{
				Type:  "string",
				Value: event.Contains,
				Name:  "contains",
			},
			{
				Type:  "string",
				Value: eventKey,
				Name:  "event_key",
			},
		},
		RoutingKey:           constants.WorkerDefaultQueue,
		RetryCount:           retryCount,
		RetryTimeout:         retryTimeout,
		BrokerMessageGroupId: getGroupIdByEventType(event.EventType, event.GroupId),
	}
	return scheduleTask, nil
}

func CreateEventHandlerSignatureForRawMessage(ctx context.Context, event structures.EventPayload) (*tasks.Signature, error) {
	cache := cacheClient.InitializeCache()
	sentryInstance := sentry.InitializeSentry()
	eventPayload := event.EventPayload
	eventBodyMap := map[string]interface{}{}
	json.Unmarshal([]byte(eventPayload), &eventBodyMap)
	groupId := constants.EventNameGroupIdMap[event.EventType]
	if entityId, ok := eventBodyMap["entity_id"].(string); ok {
		groupId = fmt.Sprintf("%s:%s", constants.EventNameGroupIdMap[event.EventType], entityId)
	} else {
		groupId = fmt.Sprintf("%s:%v", groupId, time.Now().Unix())
	}
	utils.AddLog(ctx, constants.DEBUG_LEVEL, utils.GetCurrentFunctionName(),
		map[string]interface{}{
			"event_type": event.EventType,
			"group_id":   groupId,
		}, nil)
	loggingAttributes := map[string]interface{}{
		"event_type":    event.EventType,
		"trace_id":      event.TraceID,
		"contains":      event.Contains,
		"redis_key":     event.RedisKey,
		"event_payload": event.EventPayload,
	}

	if event.RedisKey != "" {
		event.EventPayload = ""

		awsAdapter, err := utils.GetNewPubSubAdapter(constants.SNSAccessKeyID, constants.SNSSecretAccessKey)
		if err != nil {
			utils.AddLog(ctx, constants.DEBUG_LEVEL, utils.GetCurrentFunctionName(), nil, err)
			return nil, err
		}

		eventPayload, err = awsAdapter.FetchValueFromRedis(event.RedisKey)
		if err != nil {
			sentryInstance.LogError(ctx, constants.ERROR_WHILE_FETCHING_PAYLOAD_FROM_REDIS, err, loggingAttributes)
			return nil, err
		}
	}

	epoch := time.Now().UnixNano()
	eventKey := fmt.Sprintf("%s_%d", event.EventType, epoch)
	err := cache.Set(ctx, eventKey, eventPayload, constants.CacheExpiry1HourInt)
	if err != nil {
		utils.AddLog(ctx, constants.ERROR_LEVEL, utils.GetCurrentFunctionName(), nil, err)
		return nil, err
	}

	scheduleTask := &tasks.Signature{
		Name: constants.ConsumerEvents,
		Args: []tasks.Arg{
			{
				Type:  "string",
				Value: event.EventType,
				Name:  "event_type",
			},
			{
				Type:  "string",
				Value: event.TraceID,
				Name:  "trace_id",
			},
			{
				Type:  "string",
				Value: event.Contains,
				Name:  "contains",
			},
			{
				Type:  "string",
				Value: eventKey,
				Name:  "event_key",
			},
		},
		RoutingKey:           constants.WorkerDefaultQueue,
		RetryCount:           3,
		RetryTimeout:         10,
		BrokerMessageGroupId: groupId,
	}
	return scheduleTask, nil
}
