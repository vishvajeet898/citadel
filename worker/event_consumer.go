package worker

import (
	"context"

	"github.com/Orange-Health/citadel/common/constants"
	"github.com/Orange-Health/citadel/common/structures"
	"github.com/Orange-Health/citadel/common/utils"
	workerTasks "github.com/Orange-Health/citadel/worker_tasks"
	"github.com/RichardKnop/machinery/v1"
)

func SendEventHandlerToWorker(ctx context.Context, event structures.EventPayload) error {
	server, err := StartServer(true)
	if err != nil {
		return err
	}

	eventTypes := []string{
		constants.OmsAttachmentEvent, constants.OmsManualReportUploadEvent, constants.OmsTestDeleteEvent,
		constants.OmsOrderCreatedEvent, constants.OmsOrderUpdatedEvent, constants.OmsOrderCompletedEvent,
		constants.OmsOrderCancelledEvent, constants.MergeContactEvent, constants.SampleCollectedEvent,
		constants.LisEvent, constants.OrderReportPdfReadyEvent, constants.SampleRecollectionEvent,
		constants.TestDetailsEvent, constants.MarkSampleSnrEvent, constants.UpdateTaskSequence,
		constants.UpdateSrfIdToLisEvent,
	}
	if !utils.SliceContainsString(eventTypes, event.EventType) {
		return nil
	}

	loggerAttributes := map[string]interface{}{
		"event_type":    event.EventType,
		"event_payload": event.EventPayload,
		"trace_id":      event.TraceID,
		"contains":      event.Contains,
		"redis_key":     event.RedisKey,
		"group_id":      event.GroupId,
	}

	if event.EventType == constants.LisEvent {
		err = SendEventHandlerToWorkerForRawMessageMessageEvents(ctx, server, event, loggerAttributes)
		return err
	}
	scheduleTask, err := workerTasks.CreateEventHandlerSignature(ctx, event)
	if scheduleTask == nil || err != nil {
		utils.AddLog(ctx, constants.ERROR_LEVEL, utils.GetCurrentFunctionName(), loggerAttributes, err)
		return err
	}
	_, err = server.SendTask(scheduleTask)
	if err != nil {
		utils.AddLog(ctx, constants.ERROR_LEVEL, utils.GetCurrentFunctionName(), loggerAttributes, err)
		return err
	}
	utils.AddLog(ctx, constants.DEBUG_LEVEL, utils.GetCurrentFunctionName(), loggerAttributes, nil)
	return nil
}

func SendEventHandlerToWorkerForRawMessageMessageEvents(ctx context.Context, server *machinery.Server,
	event structures.EventPayload, loggerAttributes map[string]interface{}) error {

	scheduleTask, err := workerTasks.CreateEventHandlerSignatureForRawMessage(ctx, event)
	if scheduleTask == nil || err != nil {
		utils.AddLog(ctx, constants.ERROR_LEVEL, utils.GetCurrentFunctionName(), loggerAttributes, err)
		return err
	}
	_, err = server.SendTask(scheduleTask)
	if err != nil {
		utils.AddLog(ctx, constants.ERROR_LEVEL, utils.GetCurrentFunctionName(), loggerAttributes, err)
		return err
	}
	utils.AddLog(ctx, constants.DEBUG_LEVEL, utils.GetCurrentFunctionName(), loggerAttributes, nil)
	return nil
}
