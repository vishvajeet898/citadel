package dao

import (
	"gorm.io/gorm"

	commonStructures "github.com/Orange-Health/citadel/common/structures"
	commonUtils "github.com/Orange-Health/citadel/common/utils"
	commonModels "github.com/Orange-Health/citadel/models"
)

type DataLayer interface {
	GetRemarksByRemarkIds(remarkIds []uint) (
		[]commonModels.Remark, *commonStructures.CommonError)
	GetRemarksByInvestigationResultIds(remarkTypes []string,
		investigationResultIds []uint) ([]commonModels.Remark, *commonStructures.CommonError)
	GetRemarksByTaskId(taskId uint) (
		[]commonModels.Remark, *commonStructures.CommonError)

	CreateRemarksWithTx(tx *gorm.DB, remarks []commonModels.Remark) (
		[]commonModels.Remark, *commonStructures.CommonError)
	UpdateRemarksWithTx(tx *gorm.DB, remarks []commonModels.Remark) (
		[]commonModels.Remark, *commonStructures.CommonError)
	DeleteRemarksWithTx(tx *gorm.DB, remarkIds []uint, userId uint) *commonStructures.CommonError
}

func (remarkDao *RemarkDao) GetRemarksByRemarkIds(remarkIds []uint) ([]commonModels.Remark, *commonStructures.CommonError) {

	var remarks []commonModels.Remark
	err := remarkDao.Db.Where("id IN (?)", remarkIds).Find(&remarks).Error
	if err != nil {
		return remarks, commonUtils.HandleORMError(err)
	}
	return remarks, nil
}

func (remarkDao *RemarkDao) GetRemarksByInvestigationResultIds(remarkTypes []string, investigationResultIds []uint) ([]commonModels.Remark, *commonStructures.CommonError) {

	var remarks []commonModels.Remark
	query := remarkDao.Db.Where("investigation_result_id IN (?)", investigationResultIds)
	if len(remarkTypes) > 0 {
		query = query.Where("remarks.remark_type IN (?)", remarkTypes)
	}

	err := query.Find(&remarks).Error

	if err != nil {
		return remarks, commonUtils.HandleORMError(err)
	}

	return remarks, nil
}

func (remarkDao *RemarkDao) GetRemarksByTaskId(taskId uint) ([]commonModels.Remark, *commonStructures.CommonError) {

	var remarks []commonModels.Remark
	err := remarkDao.Db.
		Joins("JOIN investigation_results ON investigation_results.id = remarks.investigation_result_id").
		Joins("JOIN test_details ON test_details.id = investigation_results.test_details_id").
		Where("test_details.task_id = ?", taskId).
		Where("investigation_results.deleted_at IS NULL").
		Where("test_details.deleted_at IS NULL").
		Find(&remarks).Error

	if err != nil {
		return remarks, commonUtils.HandleORMError(err)
	}
	return remarks, nil
}

func (remarkDao *RemarkDao) CreateRemarksWithTx(tx *gorm.DB, remarks []commonModels.Remark) ([]commonModels.Remark, *commonStructures.CommonError) {

	if err := tx.Create(&remarks).Error; err != nil {
		return []commonModels.Remark{}, commonUtils.HandleORMError(err)
	}

	return remarks, nil
}

func (remarkDao *RemarkDao) UpdateRemarksWithTx(tx *gorm.DB, remarks []commonModels.Remark) ([]commonModels.Remark, *commonStructures.CommonError) {

	if err := tx.Save(&remarks).Error; err != nil {
		return []commonModels.Remark{}, commonUtils.HandleORMError(err)
	}

	return remarks, nil
}

func (remarkDao *RemarkDao) DeleteRemarksWithTx(tx *gorm.DB, remarkIds []uint, userId uint) *commonStructures.CommonError {

	currentTime := commonUtils.GetCurrentTime()
	remarkUpdates := map[string]interface{}{
		"deleted_by": userId,
		"updated_by": userId,
		"deleted_at": currentTime,
		"updated_at": currentTime,
	}
	if err := tx.Model(&commonModels.Remark{}).Where("id IN (?)", remarkIds).Updates(remarkUpdates).Error; err != nil {
		return commonUtils.HandleORMError(err)
	}

	return nil
}
