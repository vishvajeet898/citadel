package service

import (
	"context"
	"time"

	commonConstants "github.com/Orange-Health/citadel/common/constants"
	commonStructures "github.com/Orange-Health/citadel/common/structures"
	commonUtils "github.com/Orange-Health/citadel/common/utils"
	commonModels "github.com/Orange-Health/citadel/models"
)

func (etsService *EtsService) GetPublishAndCreateEtsTestEventForCron(ctx context.Context,
	etsEventDetails []commonStructures.EtsTestEvent) {
	testIdsToBeCreated := []string{}
	for _, etsEventDetail := range etsEventDetails {
		if commonUtils.SliceContainsString(testIdsToBeCreated, etsEventDetail.TestID) || etsEventDetail.Barcode == "" {
			continue
		}
		labEta, err := time.Parse(commonConstants.DateTimeUTCLayoutWithoutTZOffset, etsEventDetail.LabEta)
		if err != nil {
			continue
		}
		etsEventDetail.LabEta = labEta.Format(commonConstants.DateTimeUTCLayoutWithoutTZOffset)
		body, messageAttributes := etsService.PubsubService.GetEtsTestEvent(etsEventDetail)
		cErr := etsService.SnsClient.PublishTo(ctx, body, messageAttributes, commonConstants.EtsTestEventGroupId,
			commonConstants.EtsTestEventTopicArn, commonUtils.GetDeduplicationIdForTestEvent(etsEventDetail.TestID))
		if cErr == nil {
			testIdsToBeCreated = append(testIdsToBeCreated, etsEventDetail.TestID)
		}
	}

	if len(testIdsToBeCreated) > 0 {
		etsEvents := []commonModels.EtsEvent{}
		for _, testId := range testIdsToBeCreated {
			etsEvent := commonModels.EtsEvent{
				TestID:   testId,
				IsActive: true,
			}
			etsEvent.CreatedBy = commonConstants.CitadelSystemId
			etsEvent.UpdatedBy = commonConstants.CitadelSystemId
			etsEvents = append(etsEvents, etsEvent)
		}
		if len(etsEvents) > 0 {
			err := etsService.EtsDao.CreateEtsEvents(etsEvents)
			if err != nil {
				commonUtils.AddLog(ctx, commonConstants.ERROR_LEVEL, commonUtils.GetCurrentFunctionName(), nil, err)
				return
			}
		}
	}
}

func (etsService *EtsService) GetPublishAndUpdateEtsTestEvent(ctx context.Context,
	etsEventDetails []commonStructures.EtsTestEvent, updateDb bool) {
	testIdsToBeUpdated := []string{}
	for _, etsEventDetail := range etsEventDetails {
		if commonUtils.SliceContainsString(testIdsToBeUpdated, etsEventDetail.TestID) || etsEventDetail.Barcode == "" {
			continue
		}
		labEta, err := time.Parse(commonConstants.DateTimeUTCLayoutWithoutTZOffset, etsEventDetail.LabEta)
		if err != nil {
			continue
		}
		etsEventDetail.LabEta = labEta.Format(commonConstants.DateTimeUTCLayoutWithoutTZOffset)
		body, messageAttributes := etsService.PubsubService.GetEtsTestEvent(etsEventDetail)
		cErr := etsService.SnsClient.PublishTo(ctx, body, messageAttributes, commonConstants.EtsTestEventGroupId,
			commonConstants.EtsTestEventTopicArn, commonUtils.GetDeduplicationIdForTestEvent(etsEventDetail.TestID))
		if cErr == nil {
			testIdsToBeUpdated = append(testIdsToBeUpdated, etsEventDetail.TestID)
		}
	}

	if updateDb && len(testIdsToBeUpdated) > 0 {
		err := etsService.EtsDao.MarkEventAsInactive(testIdsToBeUpdated)
		if err != nil {
			commonUtils.AddLog(ctx, commonConstants.ERROR_LEVEL, commonUtils.GetCurrentFunctionName(), nil, err)
			return
		}
	}
}

func (etsService *EtsService) GetPublishAndCreateEtsTestEventForRerunWebhook(ctx context.Context,
	etsEventDetails []commonStructures.EtsTestEvent) {
	for _, etsEventDetail := range etsEventDetails {
		labEta, err := time.Parse(commonConstants.DateTimeUTCLayoutWithoutTZOffset, etsEventDetail.LabEta)
		if err != nil {
			continue
		}
		etsEventDetail.LabEta = labEta.Format(commonConstants.DateTimeUTCLayoutWithoutTZOffset)
		body, messageAttributes := etsService.PubsubService.GetEtsTestEvent(etsEventDetail)
		cErr := etsService.SnsClient.PublishTo(ctx, body, messageAttributes, commonConstants.EtsTestEventGroupId,
			commonConstants.EtsTestEventTopicArn, commonUtils.GetDeduplicationIdForTestEvent(etsEventDetail.TestID))
		if cErr == nil {
			etsEvent := etsService.EtsDao.GetEtsEventByTestId(etsEventDetail.TestID)
			if etsEvent.TestID == "" {
				etsEvent = commonModels.EtsEvent{
					TestID:   etsEventDetail.TestID,
					IsActive: true,
				}
				etsEvent.CreatedBy = commonConstants.CitadelSystemId
				etsEvent.UpdatedBy = commonConstants.CitadelSystemId
				if err := etsService.EtsDao.CreateEtsEvents([]commonModels.EtsEvent{etsEvent}); err != nil {
					commonUtils.AddLog(ctx, commonConstants.ERROR_LEVEL, commonUtils.GetCurrentFunctionName(), nil, err)
					return
				}
			}
		}
	}
}
