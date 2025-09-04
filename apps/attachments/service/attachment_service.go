package service

import (
	"context"
	"net/http"
	"strings"

	"gorm.io/gorm"

	mapper "github.com/Orange-Health/citadel/apps/attachments/mapper"
	"github.com/Orange-Health/citadel/apps/attachments/structures"
	commonConstants "github.com/Orange-Health/citadel/common/constants"
	commonStructures "github.com/Orange-Health/citadel/common/structures"
	commonUtils "github.com/Orange-Health/citadel/common/utils"
	commonModels "github.com/Orange-Health/citadel/models"
)

type AttachmentServiceInterface interface {
	GetAttachmentsByTaskId(taskId uint, attachmentTypes []string) (map[string][]structures.Attachment, *commonStructures.CommonError)
	GetAttachmentDtosByTaskId(taskId uint, attachmentTypes []string) ([]commonModels.Attachment, *commonStructures.CommonError)
	AddAttachment(ctx context.Context, taskId uint, addRequestRequest structures.AddAttachmentRequest) (structures.Attachment, *commonStructures.CommonError)
	CreateAttachments(attachments []commonModels.Attachment) *commonStructures.CommonError
	CreateAttachmentsWithTx(tx *gorm.DB, attachments []commonModels.Attachment) *commonStructures.CommonError
	UpdateAttachmentsWithTx(tx *gorm.DB, attachments []commonModels.Attachment) *commonStructures.CommonError
	DeleteAttachmentsByIdsWithTx(tx *gorm.DB, attachmentIds []uint) *commonStructures.CommonError
}

func getFileUrl(fileName string) string {
	return "https://" + commonConstants.Bucket + ".s3.amazonaws.com/" + fileName
}

func validateRequestParametersInAddAttachment(attachmentType, attachmentLabel, fileName, extension string) *commonStructures.CommonError {
	if attachmentType == "" {
		return &commonStructures.CommonError{
			StatusCode: http.StatusBadRequest,
			Message:    commonConstants.ATTACHMENT_TYPE_REQUIRED,
		}
	}

	if attachmentLabel == "" {
		return &commonStructures.CommonError{
			StatusCode: http.StatusBadRequest,
			Message:    commonConstants.ATTACHMENT_LABEL_REQUIRED,
		}
	}

	if !commonUtils.SliceContainsString(commonConstants.ATTACHMENT_LABELS, attachmentLabel) {
		return &commonStructures.CommonError{
			StatusCode: http.StatusBadRequest,
			Message:    commonConstants.INVALID_ATTACHMENT_LABEL,
		}
	}

	if !commonUtils.SliceContainsString(commonConstants.ATTACHMENT_TYPES, attachmentType) {
		return &commonStructures.CommonError{
			StatusCode: http.StatusBadRequest,
			Message:    commonConstants.INVALID_ATTACHMENT_TYPE,
		}
	}

	if fileName == "" {
		return &commonStructures.CommonError{
			StatusCode: http.StatusBadRequest,
			Message:    commonConstants.FILE_NAME_REQUIRED,
		}
	}

	fileNameWords := strings.Split(fileName, "/")
	if fileNameWords[0] != "citadel" {
		return &commonStructures.CommonError{
			StatusCode: http.StatusBadRequest,
			Message:    commonConstants.INVALID_FILE_NAME,
		}

	}

	if !commonUtils.SliceContainsString(commonConstants.ALLOWED_ATTACHMENT_EXTESNIONS, extension) {
		return &commonStructures.CommonError{
			StatusCode: http.StatusBadRequest,
			Message:    commonConstants.INVALID_ATTACHMENT_EXTENSION,
		}
	}

	return nil
}

func (attachmentService *AttachmentService) GetAttachmentsByTaskId(taskId uint, attachmentTypes []string) (
	map[string][]structures.Attachment, *commonStructures.CommonError) {
	attachments, cErr := attachmentService.AttachmentDao.GetAttachmentsByTaskId(taskId, attachmentTypes)
	if cErr != nil {
		return nil, cErr
	}
	attachmentStructs := mapper.MapAttachments(attachments)
	return groupAttachments(attachmentStructs), nil
}

func groupAttachments(attachments []structures.Attachment) map[string][]structures.Attachment {

	attachmentMap := make(map[string][]structures.Attachment)
	for _, attachment := range attachments {
		attachmentMap[attachment.AttachmentLabel] = append(attachmentMap[attachment.AttachmentLabel], attachment)
	}

	return attachmentMap
}

func (attachmentService *AttachmentService) GetAttachmentDtosByTaskId(taskId uint, attachmentTypes []string) (
	[]commonModels.Attachment, *commonStructures.CommonError) {
	attachments, cErr := attachmentService.AttachmentDao.GetAttachmentsByTaskId(taskId, attachmentTypes)
	if cErr != nil {
		return nil, cErr
	}
	return attachments, nil
}

func (attachmentService *AttachmentService) CreateAttachments(
	attachments []commonModels.Attachment) *commonStructures.CommonError {
	cErr := attachmentService.AttachmentDao.CreateAttachments(attachments)
	if cErr != nil {
		return cErr
	}
	return nil
}

func (attachmentService *AttachmentService) CreateAttachmentsWithTx(tx *gorm.DB,
	attachments []commonModels.Attachment) *commonStructures.CommonError {
	if len(attachments) == 0 {
		return nil
	}
	cErr := attachmentService.AttachmentDao.CreateAttachmentsWithTx(tx, attachments)
	if cErr != nil {
		return cErr
	}
	return nil
}

func (attachmentService *AttachmentService) UpdateAttachmentsWithTx(tx *gorm.DB,
	attachments []commonModels.Attachment) *commonStructures.CommonError {
	if len(attachments) == 0 {
		return nil
	}
	cErr := attachmentService.AttachmentDao.UpdateAttachmentsWithTx(tx, attachments)
	if cErr != nil {
		return cErr
	}
	return nil
}

func (attachmentService *AttachmentService) DeleteAttachmentsByIdsWithTx(tx *gorm.DB,
	attachmentIds []uint) *commonStructures.CommonError {
	if len(attachmentIds) == 0 {
		return nil
	}
	cErr := attachmentService.AttachmentDao.DeleteAttachmentsByIdsWithTx(tx, attachmentIds)
	if cErr != nil {
		return cErr
	}
	return nil
}

func (attachmentService *AttachmentService) AddAttachment(ctx context.Context, taskId uint,
	addAttachmentRequest structures.AddAttachmentRequest) (
	structures.Attachment, *commonStructures.CommonError) {

	thumbnailPublicUrl, thumbnailReference := "", ""
	extension := commonUtils.GetFileExtension(addAttachmentRequest.FileName)
	cErr := validateRequestParametersInAddAttachment(addAttachmentRequest.AttachmentType,
		addAttachmentRequest.AttachmentLabel, addAttachmentRequest.FileName, extension)
	if cErr != nil {
		return structures.Attachment{}, cErr
	}

	fileUrl := getFileUrl(addAttachmentRequest.FileName)

	tokenizedFileUrl, err := attachmentService.S3wrapperClient.GetTokenizeOrderFilePublicUrl(ctx, fileUrl)
	if err != nil {
		return structures.Attachment{}, &commonStructures.CommonError{
			StatusCode: http.StatusInternalServerError,
			Message:    commonConstants.TOKENIZED_URL_ERROR,
		}
	}

	if commonUtils.SliceContainsString(commonConstants.IMAGE_EXTENSIONS, extension) {
		mediaResizeRequest := commonStructures.MediaResizeRequest{
			Files: []commonStructures.MediaResizeFile{
				{
					FileType:      extension,
					FileUrl:       tokenizedFileUrl,
					FileReference: addAttachmentRequest.FileName,
				},
			},
		}
		mediaResizeResponse, err := attachmentService.ReportRebrandingClient.ResizeMedia(ctx, mediaResizeRequest)
		if err != nil {
			commonUtils.AddLog(ctx, commonConstants.DEBUG_LEVEL, commonConstants.ERROR_WHILE_RESIZING_MEDIA, nil, err)
		}
		if len(mediaResizeResponse.Files) != 0 {
			thumbnailReference = mediaResizeResponse.Files[0].FileReference
			thumbnailPrivateUrl := mediaResizeResponse.Files[0].FileUrl
			thumbnailPublicUrl, err = attachmentService.S3wrapperClient.GetTokenizeOrderFilePublicUrl(ctx,
				thumbnailPrivateUrl)
			if err != nil {
				commonUtils.AddLog(ctx, commonConstants.DEBUG_LEVEL, commonConstants.TOKENIZED_URL_ERROR, nil, err)
			}
		}
	}

	attachment := commonModels.Attachment{
		TaskId: taskId,
		// InvestigationResultId: investigationResult.ID,
		Reference:          addAttachmentRequest.FileName,
		AttachmentUrl:      tokenizedFileUrl,
		AttachmentType:     addAttachmentRequest.AttachmentType,
		AttachmentLabel:    addAttachmentRequest.AttachmentLabel,
		Extension:          extension,
		ThumbnailUrl:       thumbnailPublicUrl,
		ThumbnailReference: thumbnailReference,
	}

	cErr = attachmentService.CreateAttachments([]commonModels.Attachment{attachment})
	if cErr != nil {
		return structures.Attachment{}, cErr
	}

	return mapper.MapAttachment(attachment), nil
}
