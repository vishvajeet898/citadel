package consumerTasks

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/Orange-Health/citadel/common/constants"
	"github.com/Orange-Health/citadel/common/structures"
	"github.com/Orange-Health/citadel/common/utils"
)

func (eventProcessor *EventProcessor) UpdateTaskSequenceEventTask(ctx context.Context, eventPayload string) error {
	payload := structures.UpdateTaskSequenceRequest{}
	err := json.Unmarshal(json.RawMessage(eventPayload), &payload)
	if err != nil {
		return err
	}
	utils.AddLog(ctx, constants.INFO_LEVEL, utils.GetCurrentFunctionName(), map[string]interface{}{
		"event_payload": eventPayload,
		"event_type":    constants.UpdateTaskSequence,
	}, nil)

	cErr := eventProcessor.SampleService.UpdateTaskSequenceForSample(payload)
	if cErr != nil {
		return errors.New(cErr.Message)
	}
	return nil
}
