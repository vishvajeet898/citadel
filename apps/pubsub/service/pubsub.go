package service

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"strings"
	"time"

	"github.com/Orange-Health/citadel/apps/pubsub/constants"
	commonConstants "github.com/Orange-Health/citadel/common/constants"
	commonStructures "github.com/Orange-Health/citadel/common/structures"
	commonUtils "github.com/Orange-Health/citadel/common/utils"
)

type PubsubInterface interface {
	GetReportGenerationEvent(ctx context.Context, eventPayload interface{}) (
		map[string]interface{}, map[string]interface{})
	GetContactMergeConfirmEvent(ctx context.Context, eventPayload commonStructures.MergeContactEvent) (
		map[string]interface{}, map[string]interface{})
	GetReportReadyEvent(ctx context.Context, eventPayload interface{}) (
		map[string]interface{}, map[string]interface{})
	GetAddSampleRejectedTagEvent(orderId, cityCode string) (
		map[string]interface{}, map[string]interface{})
	GetRemoveSampleRejectedTagEvent(omsRequestId string, omsOrderIds []string, removeRequestTag bool, cityCode string) (
		map[string]interface{}, map[string]interface{})
	GetSampleCollectedEvent(omsTestIds []string, collectedAt *time.Time, cityCode string) (
		map[string]interface{}, map[string]interface{})
	GetResetTestTatsEvent(testIds []string, cityCode string) (
		map[string]interface{}, map[string]interface{})
	GetUpdateTestStatusEvent(omsTestStatusMap map[string]string, omsOrderId string, checkOrderCompletion bool,
		cityCode string) (map[string]interface{}, map[string]interface{})
	GetLabEtaUpdateEvent(omsOrderId string, testIds []string, lisSyncAt *time.Time, cityCode string) (
		map[string]interface{}, map[string]interface{})
	GetCheckOrderCompletionEvent(orderId, cityCode string) (
		map[string]interface{}, map[string]interface{})
	GetLisDataEvent(ctx context.Context, attuneOrderResponse commonStructures.AttuneOrderResponse) (
		map[string]interface{}, map[string]interface{})
	GetEtsTestEvent(eventPayload commonStructures.EtsTestEvent) (
		map[string]interface{}, map[string]interface{})
}

func (s *PubsubService) GetReportGenerationEvent(ctx context.Context, eventPayload interface{}) (
	map[string]interface{}, map[string]interface{}) {
	eventType := constants.OrderReportUpdateEvent
	messageAttributes := map[string]interface{}{
		"source":     "oms",
		"contains":   constants.PubSubEventContainsMap[eventType],
		"event_type": eventType,
	}

	jsonBytes, err := json.Marshal(eventPayload)
	if err != nil {
		commonUtils.AddLog(ctx, commonConstants.ERROR_LEVEL, commonConstants.ERROR_FAILED_TO_MARSHAL_PAYLOAD, nil, err)
		return nil, nil
	}

	var messageBody map[string]interface{}
	err = json.Unmarshal(jsonBytes, &messageBody)
	if err != nil {
		commonUtils.AddLog(ctx, commonConstants.ERROR_LEVEL, commonConstants.ERROR_FAILED_TO_UNMARSHAL_JSON, nil, err)
		return nil, nil
	}

	return messageBody, messageAttributes
}

func (s *PubsubService) GetReportReadyEvent(ctx context.Context, eventPayload interface{}) (
	map[string]interface{}, map[string]interface{}) {
	eventType := constants.ReportReadyEvent
	messageAttributes := map[string]interface{}{
		"source":     strings.ToLower(commonConstants.CitadelServiceName),
		"contains":   constants.PubSubEventContainsMap[eventType],
		"event_type": eventType,
	}

	jsonBytes, err := json.Marshal(eventPayload)
	if err != nil {
		commonUtils.AddLog(ctx, commonConstants.ERROR_LEVEL, commonConstants.ERROR_FAILED_TO_MARSHAL_PAYLOAD, nil, err)
		return nil, nil
	}

	var messageBody map[string]interface{}
	err = json.Unmarshal(jsonBytes, &messageBody)
	if err != nil {
		commonUtils.AddLog(ctx, commonConstants.ERROR_LEVEL, commonConstants.ERROR_FAILED_TO_UNMARSHAL_JSON, nil, err)
		return nil, nil
	}

	return messageBody, messageAttributes
}

func (s *PubsubService) GetContactMergeConfirmEvent(ctx context.Context, eventPayload commonStructures.MergeContactEvent) (
	map[string]interface{}, map[string]interface{}) {
	eventType := constants.ContactMergeConfirmEvent
	service := strings.ToLower(commonConstants.CitadelServiceName)
	messageAttributes := map[string]interface{}{
		"source":     service,
		"contains":   constants.PubSubEventContainsMap[eventType],
		"event_type": eventType,
	}
	eventPayload.Service = service

	jsonBytes, err := json.Marshal(eventPayload)
	if err != nil {
		commonUtils.AddLog(ctx, commonConstants.ERROR_LEVEL, commonConstants.ERROR_FAILED_TO_MARSHAL_PAYLOAD, nil, err)
		return nil, nil
	}

	var messageBody map[string]interface{}
	err = json.Unmarshal(jsonBytes, &messageBody)
	if err != nil {
		commonUtils.AddLog(ctx, commonConstants.ERROR_LEVEL, commonConstants.ERROR_FAILED_TO_UNMARSHAL_JSON, nil, err)
		return nil, nil
	}

	return messageBody, messageAttributes
}

func (s *PubsubService) GetAddSampleRejectedTagEvent(orderId, cityCode string) (map[string]interface{},
	map[string]interface{}) {
	eventType := constants.AddSampleRejectedTagEvent
	messageAttributes := map[string]interface{}{
		"source":     strings.ToLower(commonConstants.CitadelServiceName),
		"contains":   constants.PubSubEventContainsMap[eventType],
		"event_type": eventType,
	}

	eventPayload := map[string]interface{}{
		"order_id":            orderId,
		"servicing_city_code": cityCode,
	}

	return eventPayload, messageAttributes
}

func (s *PubsubService) GetRemoveSampleRejectedTagEvent(omsRequestId string, omsOrderIds []string,
	removeRequestTag bool, cityCode string) (map[string]interface{}, map[string]interface{}) {
	eventType := constants.RemoveSampleRejectedTagEvent
	messageAttributes := map[string]interface{}{
		"source":     strings.ToLower(commonConstants.CitadelServiceName),
		"contains":   constants.PubSubEventContainsMap[eventType],
		"event_type": eventType,
	}

	eventPayload := map[string]interface{}{
		"request_id":          omsRequestId,
		"order_ids":           omsOrderIds,
		"remove_request_tag":  removeRequestTag,
		"servicing_city_code": cityCode,
	}

	return eventPayload, messageAttributes
}

func (s *PubsubService) GetSampleCollectedEvent(omsTestIds []string, collectedAt *time.Time, cityCode string) (
	map[string]interface{}, map[string]interface{}) {
	eventType := constants.SampleCollectedEvent
	messageAttributes := map[string]interface{}{
		"source":     strings.ToLower(commonConstants.CitadelServiceName),
		"contains":   constants.PubSubEventContainsMap[eventType],
		"event_type": eventType,
	}

	eventPayload := map[string]interface{}{
		"test_ids":            omsTestIds,
		"collected_at":        collectedAt,
		"servicing_city_code": cityCode,
	}

	return eventPayload, messageAttributes
}

func (s *PubsubService) GetResetTestTatsEvent(testIds []string, cityCode string) (map[string]interface{},
	map[string]interface{}) {
	eventType := constants.ResetTestTatsEvent
	messageAttributes := map[string]interface{}{
		"source":     strings.ToLower(commonConstants.CitadelServiceName),
		"contains":   constants.PubSubEventContainsMap[eventType],
		"event_type": eventType,
	}

	eventPayload := map[string]interface{}{
		"test_ids":            testIds,
		"servicing_city_code": cityCode,
	}

	return eventPayload, messageAttributes
}

func (s *PubsubService) GetUpdateTestStatusEvent(omsTestStatusMap map[string]string, omsOrderId string,
	checkOrderCompletion bool, cityCode string) (map[string]interface{}, map[string]interface{}) {

	eventType := constants.UpdateTestStatusEvent
	messageAttributes := map[string]interface{}{
		"source":     strings.ToLower(commonConstants.CitadelServiceName),
		"contains":   constants.PubSubEventContainsMap[eventType],
		"event_type": eventType,
	}

	eventPayload := map[string]interface{}{
		"test_status_map":        omsTestStatusMap,
		"order_id":               omsOrderId,
		"check_order_completion": checkOrderCompletion,
		"servicing_city_code":    cityCode,
	}

	return eventPayload, messageAttributes
}

func (s *PubsubService) GetLabEtaUpdateEvent(omsOrderId string, testIds []string, lisSyncAt *time.Time, cityCode string) (
	map[string]interface{}, map[string]interface{}) {
	eventType := constants.LabEtaUpdateEvent
	messageAttributes := map[string]interface{}{
		"source":     strings.ToLower(commonConstants.CitadelServiceName),
		"contains":   constants.PubSubEventContainsMap[eventType],
		"event_type": eventType,
	}

	eventPayload := map[string]interface{}{
		"order_id":            omsOrderId,
		"test_ids":            testIds,
		"lis_sync_at":         lisSyncAt,
		"servicing_city_code": cityCode,
	}

	return eventPayload, messageAttributes
}

func (s *PubsubService) GetCheckOrderCompletionEvent(orderId, cityCode string) (
	map[string]interface{}, map[string]interface{}) {
	eventType := constants.CheckOrderCompletionEvent
	messageAttributes := map[string]interface{}{
		"source":     strings.ToLower(commonConstants.CitadelServiceName),
		"contains":   constants.PubSubEventContainsMap[eventType],
		"event_type": eventType,
	}

	eventPayload := map[string]interface{}{
		"order_id":            orderId,
		"servicing_city_code": cityCode,
	}

	return eventPayload, messageAttributes
}

func (s *PubsubService) GetLisDataEvent(ctx context.Context, attuneOrderResponse commonStructures.AttuneOrderResponse) (
	map[string]interface{},
	map[string]interface{}) {
	eventType := constants.CitadelLisEvent
	messageAttributes := map[string]interface{}{
		"source":     strings.ToLower(commonConstants.CitadelServiceName),
		"contains":   constants.PubSubEventContainsMap[eventType],
		"event_type": eventType,
	}

	jsonBytes, err := json.Marshal(attuneOrderResponse)
	if err != nil {
		commonUtils.AddLog(ctx, commonConstants.ERROR_LEVEL, commonConstants.ERROR_FAILED_TO_MARSHAL_PAYLOAD, nil, err)
		return nil, nil
	}

	base64BodyString := base64.StdEncoding.EncodeToString(jsonBytes)
	messageBody := map[string]interface{}{
		"entity_id":    attuneOrderResponse.OrderId,
		"entity_name":  eventType,
		"webhook_data": base64BodyString,
	}

	return messageBody, messageAttributes
}

func (s *PubsubService) GetEtsTestEvent(eventPayload commonStructures.EtsTestEvent) (map[string]interface{},
	map[string]interface{}) {
	eventType := constants.EtsTestEvent
	messageMetadataDetails := map[string]interface{}{
		"datatype": map[string]string{
			"test_id":         "int",
			"order_id":        "int",
			"test_status":     "int",
			"test_name":       "text",
			"test_deleted_at": "timestamp",
			"master_test_id":  "int",
			"lab_eta":         "timestamp",
			"lab_id":          "int",
			"patient_name":    "text",
			"patient_age":     "float",
			"patient_gender":  "text",
			"barcode":         "text",
			"vial_type_id":    "int",
			"is_rejected":     "boolean",
			"city_code":       "text",
			"lis_status":      "text",
		},
	}
	messageMetadataDetailsBytes, err := json.Marshal(messageMetadataDetails)
	if err != nil {
		commonUtils.AddLog(context.Background(), commonConstants.ERROR_LEVEL,
			commonConstants.ERROR_FAILED_TO_MARSHAL_PAYLOAD, nil, err)
		return nil, nil
	}
	messageAttributes := map[string]interface{}{
		"source":     commonConstants.CitadelServiceName,
		"contains":   string(messageMetadataDetailsBytes),
		"event_type": eventType,
	}

	jsonBytes, err := json.Marshal(eventPayload)
	if err != nil {
		commonUtils.AddLog(context.Background(), commonConstants.ERROR_LEVEL,
			commonConstants.ERROR_FAILED_TO_MARSHAL_PAYLOAD, nil, err)
		return nil, nil
	}

	var messageBody map[string]interface{}
	err = json.Unmarshal(jsonBytes, &messageBody)
	if err != nil {
		commonUtils.AddLog(context.Background(), commonConstants.ERROR_LEVEL,
			commonConstants.ERROR_FAILED_TO_UNMARSHAL_JSON, nil, err)
		return nil, nil
	}

	return messageBody, messageAttributes
}
