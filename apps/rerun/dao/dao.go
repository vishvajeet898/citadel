package dao

import (
	commonStructures "github.com/Orange-Health/citadel/common/structures"
	commonUtils "github.com/Orange-Health/citadel/common/utils"
	commonModels "github.com/Orange-Health/citadel/models"
	"gorm.io/gorm"
)

type DataLayer interface {
	GetRerunForTaskId(taskId uint) (
		[]commonModels.RerunInvestigationResult, *commonStructures.CommonError)
	GetRerunByTestDetailsIds(testDetailsIds []uint) (
		[]commonModels.RerunInvestigationResult, *commonStructures.CommonError)
	CreateRerunInvestigationResultsWithTx(tx *gorm.DB,
		rerunDetails []commonModels.RerunInvestigationResult) (
		[]commonModels.RerunInvestigationResult, *commonStructures.CommonError)
	UpdateRerunInvestigationResultsWithTx(tx *gorm.DB,
		rerunDetails []commonModels.RerunInvestigationResult) (
		[]commonModels.RerunInvestigationResult, *commonStructures.CommonError)
}

func (rerunDao *RerunDao) GetRerunForTaskId(taskId uint) (
	[]commonModels.RerunInvestigationResult, *commonStructures.CommonError) {

	reruns := []commonModels.RerunInvestigationResult{}
	if err := rerunDao.Db.
		Joins("JOIN test_details ON test_details.id = rerun_investigation_results.test_details_id").
		Where("test_details.task_id = ?", taskId).Find(&reruns).Error; err != nil {
		return reruns, commonUtils.HandleORMError(err)
	}

	return reruns, nil
}

func (rerunDao *RerunDao) GetRerunByTestDetailsIds(testDetailsIds []uint) (
	[]commonModels.RerunInvestigationResult, *commonStructures.CommonError) {

	rerunInvestigationResults := []commonModels.RerunInvestigationResult{}
	if err := rerunDao.Db.
		Where("test_details_id IN ?", testDetailsIds).Find(&rerunInvestigationResults).Error; err != nil {
		return rerunInvestigationResults, commonUtils.HandleORMError(err)
	}

	return rerunInvestigationResults, nil
}

func (rerunDao *RerunDao) CreateRerunInvestigationResultsWithTx(tx *gorm.DB,
	rerunDetails []commonModels.RerunInvestigationResult) (
	[]commonModels.RerunInvestigationResult, *commonStructures.CommonError) {

	for _, rerunDetail := range rerunDetails {
		if err := tx.Create(&rerunDetail).Error; err != nil {
			return nil, commonUtils.HandleORMError(err)
		}
	}

	return rerunDetails, nil
}

func (rerunDao *RerunDao) UpdateRerunInvestigationResultsWithTx(tx *gorm.DB, rerunDetails []commonModels.RerunInvestigationResult) (
	[]commonModels.RerunInvestigationResult, *commonStructures.CommonError) {
	if err := tx.Save(&rerunDetails).Error; err != nil {
		return nil, commonUtils.HandleORMError(err)
	}
	return rerunDetails, nil
}
