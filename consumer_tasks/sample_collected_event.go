package consumerTasks

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/Orange-Health/citadel/common/constants"
	"github.com/Orange-Health/citadel/common/structures"
	"github.com/Orange-Health/citadel/common/utils"
)

func (eventProcessor *EventProcessor) SampleCollectedEventTask(ctx context.Context, eventPayload string) error {
	sampleCollectedPayload := structures.SampleCollectedEvent{}
	err := json.Unmarshal(json.RawMessage(eventPayload), &sampleCollectedPayload)
	if err != nil {
		return err
	}
	utils.AddLog(ctx, constants.INFO_LEVEL, utils.GetCurrentFunctionName(),
		map[string]interface{}{
			"event_payload": eventPayload,
			"event_type":    constants.SampleCollectedEvent,
		}, nil)

	sampleCollectedRequest := structures.SampleCollectedRequest{
		RequestId:      sampleCollectedPayload.RequestId,
		OmsTaskType:    sampleCollectedPayload.TaskDetails.OmsTaskType,
		OmsTaskId:      sampleCollectedPayload.TaskDetails.OmsTaskId,
		CollectedAt:    sampleCollectedPayload.TaskDetails.CollectedAt,
		CollectionType: sampleCollectedPayload.TaskDetails.CollectionType,
		IsB2c:          sampleCollectedPayload.TaskDetails.IsB2c,
		UserId:         constants.CitadelSystemId,
	}
	cErr := eventProcessor.SampleService.UpdateSampleCollected(sampleCollectedRequest)
	if cErr != nil {
		return errors.New(cErr.Message)
	}
	return nil
}
