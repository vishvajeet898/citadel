package dao

import (
	"gorm.io/gorm"

	commonStructures "github.com/Orange-Health/citadel/common/structures"
	commonUtils "github.com/Orange-Health/citadel/common/utils"
	commonModels "github.com/Orange-Health/citadel/models"
)

type DataLayer interface {
	GetAttachmentsByTaskId(taskId uint, attachmentTypes []string) ([]commonModels.Attachment, *commonStructures.CommonError)
	CreateAttachments(attachments []commonModels.Attachment) *commonStructures.CommonError
	CreateAttachmentsWithTx(tx *gorm.DB, attachments []commonModels.Attachment) *commonStructures.CommonError
	UpdateAttachmentsWithTx(tx *gorm.DB, attachments []commonModels.Attachment) *commonStructures.CommonError
	DeleteAttachmentsByIdsWithTx(tx *gorm.DB, attachmentIds []uint) *commonStructures.CommonError
}

func (attachmentDao *AttachmentDao) GetAttachmentsByTaskId(taskId uint, attachmentTypes []string) ([]commonModels.Attachment,
	*commonStructures.CommonError) {
	var attachments []commonModels.Attachment
	query := attachmentDao.Db.Where("task_id = ?", taskId)
	if len(attachmentTypes) > 0 {
		query = query.Where("attachment_type IN (?)", attachmentTypes)
	}
	err := query.Find(&attachments).Error
	if err != nil {
		return []commonModels.Attachment{}, commonUtils.HandleORMError(err)
	}
	return attachments, nil
}

func (attachmentDao *AttachmentDao) CreateAttachments(attachments []commonModels.Attachment) *commonStructures.CommonError {
	err := attachmentDao.Db.Create(&attachments).Error
	if err != nil {
		return commonUtils.HandleORMError(err)
	}
	return nil
}

func (attachmentDao *AttachmentDao) CreateAttachmentsWithTx(tx *gorm.DB,
	attachments []commonModels.Attachment) *commonStructures.CommonError {
	err := tx.Create(&attachments).Error
	if err != nil {
		return commonUtils.HandleORMError(err)
	}
	return nil
}

func (attachmentDao *AttachmentDao) UpdateAttachmentsWithTx(tx *gorm.DB,
	attachments []commonModels.Attachment) *commonStructures.CommonError {
	err := tx.Save(&attachments).Error
	if err != nil {
		return commonUtils.HandleORMError(err)
	}
	return nil
}

func (attachmentDao *AttachmentDao) DeleteAttachmentsByIdsWithTx(tx *gorm.DB,
	attachmentIds []uint) *commonStructures.CommonError {
	err := tx.Where("id IN (?)", attachmentIds).Delete(&commonModels.Attachment{}).Error
	if err != nil {
		return commonUtils.HandleORMError(err)
	}
	return nil
}
