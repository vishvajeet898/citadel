package commonTasks

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/Orange-Health/citadel/common/constants"
	"github.com/Orange-Health/citadel/common/structures"
	"github.com/Orange-Health/citadel/common/utils"
	"github.com/Orange-Health/citadel/models"
)

func (ctp *CommonTaskProcessor) addDocumentsByLisEvent(ctx context.Context, taskId uint,
	documents []structures.DocumentWithMetadata, folderPath string, attachmentType string) error {
	if taskId == 0 {
		return nil
	}

	attachments, _ := ctp.AttachmentsService.GetAttachmentDtosByTaskId(taskId, nil)
	attachmentsReferenceMap := map[string]bool{}
	createAttachments := []models.Attachment{}
	for _, attachment := range attachments {
		attachmentsReferenceMap[attachment.Reference] = true
	}

	if len(documents) > 0 {
		for _, document := range documents {
			documentUrl := strings.ReplaceAll(document.Url, " ", "")
			referenceNumber := fetchReferenceNumberFromDocument(documentUrl)
			if referenceNumber == "" {
				continue
			}

			currentTimeInMilliseconds := utils.GetCurrentTimeInMilliseconds()
			filepath := "/tmp_testing_" + referenceNumber + "_" + fmt.Sprint(currentTimeInMilliseconds)
			file, extension, err := downloadFile(documentUrl, filepath)
			if err != nil {
				continue
			}
			defer file.Close()

			if file == nil || extension == "" {
				continue
			}

			fileName := "citadel/" + folderPath + "/" + referenceNumber + "." + extension

			fileBlob, err := os.Open(filepath)
			if err != nil {
				continue
			}
			defer fileBlob.Close()
			awsFilePath, _ := ctp.S3Client.UploadLocalFile(fileBlob, fileName, constants.Bucket)

			err = os.Remove(filepath)
			if err != nil {
				utils.AddLog(ctx, constants.ERROR_LEVEL, constants.ERROR_FAILED_TO_DELETE_FILE, nil, err)
			}

			publicUrl, err := ctp.S3wrapperClient.GetTokenizeOrderFilePublicUrl(ctx, awsFilePath)
			if err != nil {
				continue
			}

			if _, ok := attachmentsReferenceMap[fileName]; !ok {
				createAttachments = append(createAttachments, models.Attachment{
					TaskId:                taskId,
					Reference:             fileName,
					InvestigationResultId: document.InvestigationId,
					AttachmentUrl:         publicUrl,
					AttachmentType:        attachmentType,
					AttachmentLabel:       constants.AttachmentLabelLabDocs,
					Extension:             extension,
				})
				attachmentsReferenceMap[fileName] = true
			}
		}
	}
	if len(createAttachments) > 0 {
		createAttachments = ctp.AddThumbnailUrl(ctx, createAttachments)
		cErr := ctp.AttachmentsService.CreateAttachments(createAttachments)
		if cErr != nil {
			return errors.New(cErr.Message)
		}
	}

	return nil
}

func (ctp *CommonTaskProcessor) AddVisitDocumentTaskByLisEvent(ctx context.Context, visitId string, taskId uint,
	visitDocuments []string) error {
	documents := make([]structures.DocumentWithMetadata, len(visitDocuments))
	for index, documentUrl := range visitDocuments {
		documents[index] = structures.DocumentWithMetadata{
			Url:             documentUrl,
			InvestigationId: 0,
		}
	}
	return ctp.addDocumentsByLisEvent(ctx, taskId, documents, constants.FolderPathVisitDocument, constants.AttachmentTypeVisitDocument)
}

func (ctp *CommonTaskProcessor) AddTestDocumentTaskByLisEvent(ctx context.Context, taskId uint,
	testDocumentMap map[string][]structures.TestDocumentInfoResponse) error {
	if taskId == 0 {
		return nil
	}

	var documents []structures.DocumentWithMetadata
	for _, testDocs := range testDocumentMap {
		for _, testDoc := range testDocs {
			documents = append(documents, structures.DocumentWithMetadata{
				Url:             testDoc.TestDocument,
				InvestigationId: testDoc.InvestigationId,
			})
		}
	}
	return ctp.addDocumentsByLisEvent(ctx, taskId, documents, constants.FolderPathTestDocument, constants.AttachmentTypeTestDocument)
}

func downloadFile(url, filepath string) (*os.File, string, error) {
	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return nil, "", err
	}

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return nil, "", err
	}
	defer resp.Body.Close()

	headers := resp.Header
	contentDisposition := headers.Get("Content-Disposition")
	contents := strings.Split(contentDisposition, ".")
	extension := contents[len(contents)-1]
	extension = strings.ReplaceAll(extension, "\"", "")

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return nil, "", err
	}

	return out, extension, nil
}

func fetchReferenceNumberFromDocument(visitDocument string) string {
	references := strings.Split(visitDocument, "=")
	if len(references) > 1 {
		return references[1]
	}

	return ""
}

func (ctp *CommonTaskProcessor) AddThumbnailUrl(ctx context.Context, attachments []models.Attachment) []models.Attachment {
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

	resizeMediaResponse, _ := ctp.ReportRebrandingClient.ResizeMedia(ctx, resizeMediaRequest)

	for _, file := range resizeMediaResponse.Files {
		attachmentIdToResizedFileMap[int(file.Id)] = file
	}

	for key, attachment := range attachmentIdToAttachmentsMap {
		updatedAttachment := attachment
		if resizedFile, ok := attachmentIdToResizedFileMap[key]; ok {
			if resizedFile.FileUrl != "" {
				thumbnailPublicUrl, err := ctp.S3wrapperClient.GetTokenizeOrderFilePublicUrl(ctx,
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
