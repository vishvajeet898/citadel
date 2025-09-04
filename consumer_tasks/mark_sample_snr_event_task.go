package consumerTasks

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/Orange-Health/citadel/common/constants"
	"github.com/Orange-Health/citadel/common/structures"
	"github.com/Orange-Health/citadel/common/utils"
)

func (eventProcessor *EventProcessor) MarkSampleSnrEventTask(ctx context.Context, eventPayload string) error {
	payload := structures.MarkSampleSnrEvent{}
	err := json.Unmarshal(json.RawMessage(eventPayload), &payload)
	if err != nil {
		return err
	}

	utils.AddLog(ctx, constants.INFO_LEVEL, utils.GetCurrentFunctionName(),
		map[string]interface{}{
			"event_payload": eventPayload,
			"event_type":    constants.MarkSampleSnrEvent,
		}, nil)

	samples, _, cErr := eventProcessor.SampleService.GetSamplesForTests(payload.AlnumTestIds)
	if cErr != nil {
		return errors.New(cErr.Message)
	}
	if len(samples) == 0 {
		return nil // No samples found for the provided test IDs
	}

	for _, sample := range samples {
		markSampleNotReceivedRequest := structures.MarkAsNotReceivedRequest{
			UserId:            constants.CitadelSystemId,
			LabId:             sample.LabId,
			SampleId:          sample.Id,
			NotReceivedReason: constants.AutoCancellationDueToRNRReason,
		}
		cErr := eventProcessor.ReceivingDeskService.MarkAsNotReceived(ctx, markSampleNotReceivedRequest, false)
		if cErr != nil {
			return errors.New(cErr.Message)
		}
	}
	return nil
}
