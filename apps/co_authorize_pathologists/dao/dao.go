package dao

import (
	"gorm.io/gorm"

	commonStructures "github.com/Orange-Health/citadel/common/structures"
	commonUtils "github.com/Orange-Health/citadel/common/utils"
	commonModels "github.com/Orange-Health/citadel/models"
)

type DataLayer interface {
	CreateCoAuthrizePathologistWithTx(tx *gorm.DB,
		coAuthorizePathologistDetail commonModels.CoAuthorizedPathologists) (
		commonModels.CoAuthorizedPathologists, *commonStructures.CommonError)
	DeleteCurrentCoAuthorizePathologistWithTx(tx *gorm.DB,
		taskId, userId uint) *commonStructures.CommonError
}

func (coAuthorizePathologistDao *CoAuthorizePathologistDao) CreateCoAuthrizePathologistWithTx(
	tx *gorm.DB, coAuthorizePathologistDetail commonModels.CoAuthorizedPathologists) (
	commonModels.CoAuthorizedPathologists, *commonStructures.CommonError) {
	err := tx.Create(&coAuthorizePathologistDetail).Error
	if err != nil {
		return commonModels.CoAuthorizedPathologists{}, commonUtils.HandleORMError(err)
	}

	return coAuthorizePathologistDetail, nil
}

func (coAuthorizePathologistDao *CoAuthorizePathologistDao) DeleteCurrentCoAuthorizePathologistWithTx(
	tx *gorm.DB, taskId, userId uint) *commonStructures.CommonError {
	currentTime := commonUtils.GetCurrentTime()
	updates := map[string]interface{}{
		"deleted_at": currentTime,
		"updated_at": currentTime,
		"deleted_by": userId,
		"updated_by": userId,
	}

	err := tx.Model(&commonModels.CoAuthorizedPathologists{}).Where("task_id = ?", taskId).Updates(updates).Error
	if err != nil {
		return commonUtils.HandleORMError(err)
	}

	return nil
}
