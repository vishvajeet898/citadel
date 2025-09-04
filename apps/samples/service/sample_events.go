package service

import (
	"context"
	"errors"
	"time"

	commonConstants "github.com/Orange-Health/citadel/common/constants"
	commonStructures "github.com/Orange-Health/citadel/common/structures"
	commonUtils "github.com/Orange-Health/citadel/common/utils"
)

func (sampleService *SampleService) PublishSampleCollectedEvent(omsTestIds []string, collectedAt *time.Time, cityCode string) {
	ctx := context.Background()
	messageBody, messageAttributes := sampleService.PubsubService.GetSampleCollectedEvent(omsTestIds, collectedAt, cityCode)
	cErr := sampleService.SnsClient.PublishTo(ctx, messageBody, messageAttributes, "", commonConstants.OmsUpdatesTopicArn,
		"")
	if cErr != nil {
		commonUtils.AddLog(ctx, commonConstants.ERROR_LEVEL, commonUtils.GetCurrentFunctionName(), nil,
			errors.New(cErr.Message))
	}
}

func (sampleService *SampleService) PublishResetTatsEvent(omsTestIds []string, cityCode string) {
	ctx := context.Background()
	messageBody, messageAttributes := sampleService.PubsubService.GetResetTestTatsEvent(omsTestIds, cityCode)
	cErr := sampleService.SnsClient.PublishTo(ctx, messageBody, messageAttributes, "", commonConstants.OmsUpdatesTopicArn,
		"")
	if cErr != nil {
		commonUtils.AddLog(ctx, commonConstants.ERROR_LEVEL, commonUtils.GetCurrentFunctionName(), nil,
			errors.New(cErr.Message))
	}
}

func (sampleService *SampleService) PublishAddSampleRejectedTagEvent(omsOrderId, servicingCityCode string) {
	ctx := context.Background()
	messageBody, messageAttributes := sampleService.PubsubService.GetAddSampleRejectedTagEvent(omsOrderId, servicingCityCode)
	if cErr := sampleService.SnsClient.PublishTo(ctx, messageBody, messageAttributes, "", commonConstants.OmsUpdatesTopicArn,
		""); cErr != nil {
		commonUtils.AddLog(ctx, commonConstants.ERROR_LEVEL, commonUtils.GetCurrentFunctionName(), nil,
			errors.New(cErr.Message))
	}
}

func (sampleService *SampleService) PublishRemoveSampleRejectedTagEvent(omsRequestId string, omsOrderIds []string,
	removeRequestTag bool, cityCode string) {
	ctx := context.Background()
	messageBody, messageAttributes := sampleService.PubsubService.GetRemoveSampleRejectedTagEvent(omsRequestId,
		omsOrderIds, removeRequestTag, cityCode)
	if cErr := sampleService.SnsClient.PublishTo(ctx, messageBody, messageAttributes, "", commonConstants.OmsUpdatesTopicArn,
		""); cErr != nil {
		commonUtils.AddLog(ctx, commonConstants.ERROR_LEVEL, commonUtils.GetCurrentFunctionName(), nil,
			errors.New(cErr.Message))
	}
}

func (sampleService *SampleService) PublishUpdateTestStatusEvent(omsTestStatusMap map[string]string, omsOrderId string,
	checkOrderCompletion bool, cityCode string) {
	ctx := context.Background()
	messageBody, messageAttributes := sampleService.PubsubService.GetUpdateTestStatusEvent(omsTestStatusMap, omsOrderId,
		checkOrderCompletion, cityCode)
	cErr := sampleService.SnsClient.PublishTo(ctx, messageBody, messageAttributes, "", commonConstants.OmsUpdatesTopicArn,
		"")
	if cErr != nil {
		commonUtils.AddLog(ctx, commonConstants.ERROR_LEVEL, commonUtils.GetCurrentFunctionName(), nil,
			errors.New(cErr.Message))
	}
}

func (sampleService *SampleService) PublishLisDataEvent(attuneOrderResponse commonStructures.AttuneOrderResponse) {
	ctx := context.Background()
	messageBody, messageAttributes := sampleService.PubsubService.GetLisDataEvent(ctx, attuneOrderResponse)
	cErr := sampleService.SnsClient.PublishTo(ctx, messageBody, messageAttributes, "", commonConstants.CitadelTopicArn, "")
	if cErr != nil {
		commonUtils.AddLog(ctx, commonConstants.ERROR_LEVEL, commonUtils.GetCurrentFunctionName(), nil,
			errors.New(cErr.Message))
	}
}
