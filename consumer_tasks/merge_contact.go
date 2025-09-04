package consumerTasks

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/Orange-Health/citadel/common/structures"
)

func (eventProcessor *EventProcessor) MergeContactEventTask(ctx context.Context, eventPayload string) error {
	mergeContactPayload := structures.MergeContactEvent{}
	err := json.Unmarshal(json.RawMessage(eventPayload), &mergeContactPayload)
	if err != nil {
		return err
	}
	commonError := eventProcessor.ContactService.Merge(ctx, mergeContactPayload)
	if commonError != nil {
		return errors.New(commonError.Message)
	}
	return nil
}
