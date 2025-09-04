package consumerTasks

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/Orange-Health/citadel/common/constants"
	"github.com/Orange-Health/citadel/common/structures"
	"github.com/Orange-Health/citadel/common/utils"
)

func (eventProcessor *EventProcessor) UpdateSrfIdToLisEventTask(ctx context.Context, eventPayload string) error {
	payload := structures.UpdateSrfIdToLisRequest{}
	err := json.Unmarshal(json.RawMessage(eventPayload), &payload)
	if err != nil {
		return err
	}
	utils.AddLog(ctx, constants.INFO_LEVEL, utils.GetCurrentFunctionName(), map[string]interface{}{
		"event_payload": eventPayload,
		"event_type":    constants.UpdateSrfIdToLisEvent,
	}, nil)

	orderDetails, cErr := eventProcessor.OrderDetailsService.GetOrderDetailsByOmsOrderId(payload.AlnumOrderId)
	if cErr != nil {
		return errors.New(cErr.Message)
	}
	if orderDetails.Id == 0 {
		return nil // No order details found, nothing to update
	}

	if orderDetails.SrfId != payload.SrfId {
		orderDetails.SrfId = payload.SrfId
		orderDetails.UpdatedBy = constants.CitadelSystemId
		_, cErr = eventProcessor.OrderDetailsService.UpdateOrderDetails(orderDetails)
		if cErr != nil {
			return errors.New(cErr.Message)
		}
	}

	cErr = eventProcessor.SampleService.UpdateSrfIdToLis(ctx, orderDetails)
	if cErr != nil {
		return errors.New(cErr.Message)
	}
	return nil
}
