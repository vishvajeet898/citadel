package service

import (
	"context"
	"errors"

	commonConstants "github.com/Orange-Health/citadel/common/constants"
	commonStructures "github.com/Orange-Health/citadel/common/structures"
	commonUtils "github.com/Orange-Health/citadel/common/utils"
)

type ContactServiceInterface interface {
	Merge(
		ctx context.Context,
		eventPayload commonStructures.MergeContactEvent,
	) *commonStructures.CommonError
}

func (s *ContactService) publishMergeConfirmEvent(ctx context.Context, eventPayload commonStructures.MergeContactEvent) *commonStructures.CommonError {
	messageBody, messageAttributes := s.PubsubService.GetContactMergeConfirmEvent(ctx, eventPayload)
	return s.SnsClient.PublishTo(ctx, messageBody, messageAttributes, "", commonConstants.ContactMergeConfirmTopicArn, "")
}

func (s *ContactService) mergeExternalInvestigations(masterContactId, mergeContactId uint) (bool, *commonStructures.CommonError) {
	err := s.ExternalInvestigationResultService.UpdateContact(
		mergeContactId,
		masterContactId,
	)
	if err != nil {
		return false, err
	}
	return false, nil
}

func (s *ContactService) Merge(
	ctx context.Context,
	eventPayload commonStructures.MergeContactEvent,
) *commonStructures.CommonError {
	wasMergeContactFound := false
	mergeContactId := commonUtils.ConvertStringToUint(eventPayload.MergeContact)
	_, err := s.mergeExternalInvestigations(eventPayload.MasterContact.Id, mergeContactId)
	if err != nil {
		s.Sentry.LogError(ctx, commonConstants.ERROR_FAILED_TO_MERGE_INVESTIGATIONS, errors.New(err.Message), nil)
		return err
	}
	eventPayload.WasMergeContactFound = wasMergeContactFound
	return s.publishMergeConfirmEvent(ctx, eventPayload)
}
