package consumerTasks

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/Orange-Health/citadel/common/constants"
	"github.com/Orange-Health/citadel/common/structures"
	"github.com/Orange-Health/citadel/common/utils"
)

func (eventProcessor *EventProcessor) OmsOrderCompletedEventTask(ctx context.Context,
	eventPayload string) error {

	omsOrderCreateUpdateEvent := structures.OmsOrderCreateUpdateEvent{}
	err := json.Unmarshal(json.RawMessage(eventPayload), &omsOrderCreateUpdateEvent)
	if err != nil {
		eventProcessor.Sentry.LogError(ctx, constants.ERROR_FAILED_TO_UNMARSHAL_JSON, err, nil)
		return err
	}

	redisKey := fmt.Sprintf(constants.OmsCompletedOrderEventKey, omsOrderCreateUpdateEvent.Order.Id,
		omsOrderCreateUpdateEvent.CityCode)
	keyExists, err := eventProcessor.Cache.Exists(ctx, redisKey)
	if err != nil || keyExists {
		utils.AddLog(ctx, constants.DEBUG_LEVEL, utils.GetCurrentFunctionName(), nil, err)
		return errors.New(constants.ERROR_COMPLETED_ORDER_TASK_IN_PROGRESS)
	}

	err = eventProcessor.Cache.Set(ctx, redisKey, true, constants.CacheExpiry10MinutesInt)
	if err != nil {
		utils.AddLog(ctx, constants.DEBUG_LEVEL, utils.GetCurrentFunctionName(), nil, err)
		return err
	}

	defer func() {
		err := eventProcessor.Cache.Delete(ctx, redisKey)
		if err != nil {
			utils.AddLog(ctx, constants.DEBUG_LEVEL, utils.GetCurrentFunctionName(), nil, err)
		}
	}()

	utils.AddLog(ctx, constants.INFO_LEVEL, utils.GetCurrentFunctionName(),
		map[string]interface{}{
			"event_payload": eventPayload,
			"event_type":    constants.OmsOrderCompletedEvent,
		}, nil)

	cErr := eventProcessor.ProcessOmsOrderCompletedEvent(ctx, omsOrderCreateUpdateEvent)
	if cErr != nil {
		return errors.New(cErr.Message)
	}

	return nil
}

func (eventProcessor *EventProcessor) ProcessOmsOrderCompletedEvent(ctx context.Context,
	omsOrderEvent structures.OmsOrderCreateUpdateEvent) *structures.CommonError {

	orderDetails, cErr := eventProcessor.OrderDetailsService.GetOrderDetailsByOmsOrderId(omsOrderEvent.Order.AlnumOrderId)
	if cErr != nil && cErr.Message != constants.ERROR_ORDER_ID_NOT_FOUND {
		return cErr
	}

	if orderDetails.Id == 0 {
		return nil
	}

	orderDetails.OrderStatus = constants.OrderStatusMapUint[omsOrderEvent.Order.Status]
	_, cErr = eventProcessor.OrderDetailsService.UpdateOrderDetails(orderDetails)
	if cErr != nil {
		return cErr
	}

	return nil
}
