package consumerTasks

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"gorm.io/gorm"

	"github.com/Orange-Health/citadel/common/constants"
	"github.com/Orange-Health/citadel/common/structures"
	"github.com/Orange-Health/citadel/common/utils"
	"github.com/Orange-Health/citadel/models"
)

func (eventProcessor *EventProcessor) OmsAttachmentEventTask(ctx context.Context, eventPayload string) error {
	omsAttachmentEvent := structures.OmsAttachmentEvent{}
	err := json.Unmarshal(json.RawMessage(eventPayload), &omsAttachmentEvent)
	if err != nil {
		eventProcessor.Sentry.LogError(ctx, constants.ERROR_FAILED_TO_UNMARSHAL_JSON, err, nil)
		return nil
	}

	redisKey := fmt.Sprintf(constants.OmsAttachmentEventKey, omsAttachmentEvent.AlnumOrderId, omsAttachmentEvent.CityCode)
	keyExists, err := eventProcessor.Cache.Exists(ctx, redisKey)
	if err != nil || keyExists {
		utils.AddLog(ctx, constants.DEBUG_LEVEL, utils.GetCurrentFunctionName(), nil, err)
		return errors.New(constants.ERROR_ATTACHMENT_EVENT_TASK_IN_PROGRESS)
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

	task, cErr := eventProcessor.TaskService.GetTaskByOmsOrderId(omsAttachmentEvent.AlnumOrderId)
	if cErr != nil || task.Id == 0 {
		return nil
	}

	currentAttachments, cErr := eventProcessor.AttachmentsService.GetAttachmentDtosByTaskId(task.Id, nil)
	if cErr != nil {
		err := errors.New(cErr.Message)
		utils.AddLog(ctx, constants.ERROR_LEVEL, utils.GetCurrentFunctionName(), map[string]interface{}{
			"error": err,
		}, nil)
		return nil
	}

	createAttachments, updateAttachments, deleteAttachmentsIds := eventProcessor.getAttachmentDtos(task.Id,
		currentAttachments, omsAttachmentEvent.Attachments)
	if len(createAttachments) > 0 {
		createAttachments = eventProcessor.addThumbnailUrl(ctx, createAttachments)
	}
	if len(updateAttachments) > 0 {
		updateAttachments = eventProcessor.addThumbnailUrl(ctx, updateAttachments)
	}

	if len(createAttachments) > 0 || len(updateAttachments) > 0 {

		err := eventProcessor.Db.Transaction(func(tx *gorm.DB) error {
			cErr := eventProcessor.AttachmentsService.CreateAttachmentsWithTx(tx, createAttachments)
			if cErr != nil {
				err := errors.New(cErr.Message)
				utils.AddLog(ctx, constants.ERROR_LEVEL, utils.GetCurrentFunctionName(), nil, err)
				return err
			}

			cErr = eventProcessor.AttachmentsService.UpdateAttachmentsWithTx(tx, updateAttachments)
			if cErr != nil {
				err := errors.New(cErr.Message)
				utils.AddLog(ctx, constants.ERROR_LEVEL, utils.GetCurrentFunctionName(), nil, err)
				return err
			}

			cErr = eventProcessor.AttachmentsService.DeleteAttachmentsByIdsWithTx(tx, deleteAttachmentsIds)
			if cErr != nil {
				err := errors.New(cErr.Message)
				utils.AddLog(ctx, constants.ERROR_LEVEL, utils.GetCurrentFunctionName(), nil, err)
				return err
			}

			return nil
		})

		if err != nil {
			utils.AddLog(ctx, constants.ERROR_LEVEL, utils.GetCurrentFunctionName(), nil, err)
			return nil
		}
	}

	return nil
}

func (eventProcessor *EventProcessor) getAttachmentDtos(taskId uint, currentAttachments []models.Attachment,
	omsAttachments []structures.OmsAttachmentStruct) (createAttachmentsFinal []models.Attachment,
	updateAttachmentsFinal []models.Attachment, deleteAttachmentsIdsFinal []uint) {
	createAttachmentsFinal = make([]models.Attachment, 0)
	updateAttachmentsFinal = make([]models.Attachment, 0)
	deleteAttachmentsIdsFinal = make([]uint, 0)

	currentNonLabDocsAttachments, currentOutsourceDocsAttachments := []models.Attachment{}, []models.Attachment{}
	for _, currentAttachment := range currentAttachments {
		switch currentAttachment.AttachmentLabel {
		case constants.AttachmentLabelNonLabDocs:
			currentNonLabDocsAttachments = append(currentNonLabDocsAttachments, currentAttachment)
		case constants.AttachmentLabelOutsourceDocs:
			currentOutsourceDocsAttachments = append(currentOutsourceDocsAttachments, currentAttachment)
		}
	}

	omsNonLabDocs, omsOutsourceDocs := []structures.OmsAttachmentStruct{}, []structures.OmsAttachmentStruct{}
	for _, omsAttachment := range omsAttachments {
		switch omsAttachment.FileLabel {
		case constants.AttachmentLabelNonLabDocs:
			omsNonLabDocs = append(omsNonLabDocs, omsAttachment)
		case constants.AttachmentLabelOutsourceDocs:
			omsOutsourceDocs = append(omsOutsourceDocs, omsAttachment)
		}
	}

	createAttachments, updateAttachments, deleteAttachmentsIds :=
		eventProcessor.handleNonLabDocs(taskId, currentNonLabDocsAttachments, omsNonLabDocs)
	createAttachmentsFinal = append(createAttachmentsFinal, createAttachments...)
	updateAttachmentsFinal = append(updateAttachmentsFinal, updateAttachments...)
	deleteAttachmentsIdsFinal = append(deleteAttachmentsIdsFinal, deleteAttachmentsIds...)

	createAttachments, updateAttachments, deleteAttachmentsIds =
		eventProcessor.handleOutsourceDocs(taskId, currentOutsourceDocsAttachments, omsOutsourceDocs)
	createAttachmentsFinal = append(createAttachmentsFinal, createAttachments...)
	updateAttachmentsFinal = append(updateAttachmentsFinal, updateAttachments...)
	deleteAttachmentsIdsFinal = append(deleteAttachmentsIdsFinal, deleteAttachmentsIds...)

	return
}

func (eventProcessor *EventProcessor) handleNonLabDocs(taskId uint, currentNonLabDocsAttachments []models.Attachment,
	omsNonLabDocs []structures.OmsAttachmentStruct) (createAttachments []models.Attachment, updateAttachments []models.Attachment,
	deleteAttachmentsIds []uint) {
	createAttachments = make([]models.Attachment, 0)
	updateAttachments = make([]models.Attachment, 0)
	deleteAttachmentsIds = make([]uint, 0)

	omsAttachmentReferenceMap := make(map[string]bool)
	for _, omsAttachment := range omsNonLabDocs {
		omsAttachmentReferenceMap[omsAttachment.FileReference] = true
	}

	for _, currentAttachment := range currentNonLabDocsAttachments {
		if _, ok := omsAttachmentReferenceMap[currentAttachment.Reference]; !ok {
			deleteAttachmentsIds = append(deleteAttachmentsIds, currentAttachment.Id)
		}
	}

	for _, omsAttachment := range omsNonLabDocs {
		attachmentExists := false
		for _, currentAttachment := range currentNonLabDocsAttachments {
			if currentAttachment.Reference == omsAttachment.FileReference {
				attachmentExists = true
				currentAttachment.AttachmentUrl = omsAttachment.FileUrl
				updateAttachments = append(updateAttachments, currentAttachment)
				break
			}
		}

		if !attachmentExists {
			createAttachments = append(createAttachments, models.Attachment{
				TaskId:          taskId,
				Reference:       omsAttachment.FileReference,
				AttachmentUrl:   omsAttachment.FileUrl,
				AttachmentType:  omsAttachment.FileType,
				AttachmentLabel: omsAttachment.FileLabel,
				Extension:       omsAttachment.FileExtension,
			})
		}
	}
	return
}

func (eventProcessor *EventProcessor) handleOutsourceDocs(taskId uint, currentOutsourceDocsAttachments []models.Attachment,
	omsOutsourceDocs []structures.OmsAttachmentStruct) (createAttachments []models.Attachment, updateAttachments []models.Attachment,
	deleteAttachmentsIds []uint) {
	createAttachments = make([]models.Attachment, 0)
	updateAttachments = make([]models.Attachment, 0)
	deleteAttachmentsIds = make([]uint, 0)

	omsAttachmentReferenceMap := make(map[string]bool)
	for _, omsAttachment := range omsOutsourceDocs {
		omsAttachmentReferenceMap[omsAttachment.FileReference] = true
	}

	for _, currentAttachment := range currentOutsourceDocsAttachments {
		if _, ok := omsAttachmentReferenceMap[currentAttachment.Reference]; !ok {
			deleteAttachmentsIds = append(deleteAttachmentsIds, currentAttachment.Id)
		}
	}

	for _, omsAttachment := range omsOutsourceDocs {
		attachmentExists := false
		for _, currentAttachment := range currentOutsourceDocsAttachments {
			if currentAttachment.Reference == omsAttachment.FileReference {
				attachmentExists = true
				currentAttachment.AttachmentUrl = omsAttachment.FileUrl
				updateAttachments = append(updateAttachments, currentAttachment)
				break
			}
		}

		if !attachmentExists {
			createAttachments = append(createAttachments, models.Attachment{
				TaskId:          taskId,
				Reference:       omsAttachment.FileReference,
				AttachmentUrl:   omsAttachment.FileUrl,
				AttachmentType:  omsAttachment.FileType,
				AttachmentLabel: omsAttachment.FileLabel,
				Extension:       omsAttachment.FileExtension,
			})
		}
	}
	return
}

func (eventProcessor *EventProcessor) addThumbnailUrl(ctx context.Context, attachments []models.Attachment) []models.Attachment {
	if len(attachments) == 0 {
		return attachments
	}

	attachmentIdToAttachmentsMap, attachmentIdToResizedFileMap :=
		make(map[int]models.Attachment), make(map[int]structures.MediaResizeFile)
	updatedAttachments := make([]models.Attachment, 0)

	for index := range attachments {
		if utils.SliceContainsString(constants.IMAGE_EXTENSIONS, attachments[index].Extension) {
			attachmentIdToAttachmentsMap[index] = attachments[index]
		} else {
			updatedAttachments = append(updatedAttachments, attachments[index])
		}
	}

	resizeMediaRequest := structures.MediaResizeRequest{}
	resizeMediaFiles := make([]structures.MediaResizeFile, 0)

	for key, attachment := range attachmentIdToAttachmentsMap {
		resizeMediaFiles = append(resizeMediaFiles, structures.MediaResizeFile{
			Id:            uint(key),
			FileType:      attachment.Extension,
			FileUrl:       attachment.AttachmentUrl,
			FileReference: attachment.Reference,
		})
	}
	resizeMediaRequest.Files = resizeMediaFiles

	resizeMediaResponse, _ := eventProcessor.ReportRebrandingClient.ResizeMedia(ctx, resizeMediaRequest)

	for _, file := range resizeMediaResponse.Files {
		attachmentIdToResizedFileMap[int(file.Id)] = file
	}

	for key, attachment := range attachmentIdToAttachmentsMap {
		updatedAttachment := attachment
		if resizedFile, ok := attachmentIdToResizedFileMap[key]; ok {
			if resizedFile.FileUrl != "" {
				thumbnailPublicUrl, err := eventProcessor.S3wrapperClient.GetTokenizeOrderFilePublicUrl(ctx,
					resizedFile.FileUrl)
				if err == nil {
					updatedAttachment.ThumbnailUrl = thumbnailPublicUrl
					updatedAttachment.ThumbnailReference = resizedFile.FileReference
				}
			}
		}
		updatedAttachments = append(updatedAttachments, updatedAttachment)
	}

	return updatedAttachments
}
