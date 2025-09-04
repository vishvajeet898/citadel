package service

import (
	"gorm.io/gorm"

	commonStructures "github.com/Orange-Health/citadel/common/structures"
	commonModels "github.com/Orange-Health/citadel/models"
)

type RerunServiceInterface interface {
	GetRerunDetailsByTaskId(taskId uint) (
		[]commonModels.RerunInvestigationResult, *commonStructures.CommonError)
	GetRerunDetailsByTestDetailsIds(testDetailsIds []uint) (
		[]commonModels.RerunInvestigationResult, *commonStructures.CommonError)
	CreateRerunInvestigationResultsWithTx(tx *gorm.DB,
		rerunDetails []commonModels.RerunInvestigationResult) (
		[]commonModels.RerunInvestigationResult, *commonStructures.CommonError)
	UpdateRerunInvestigationResultsWithTx(tx *gorm.DB,
		rerunDetails []commonModels.RerunInvestigationResult) (
		[]commonModels.RerunInvestigationResult, *commonStructures.CommonError)
}

func (rerunService *RerunService) GetRerunDetailsByTaskId(taskId uint) (
	[]commonModels.RerunInvestigationResult, *commonStructures.CommonError) {

	return rerunService.RerunDao.GetRerunForTaskId(taskId)
}

func (rerunService *RerunService) GetRerunDetailsByTestDetailsIds(testDetailsIds []uint) (
	[]commonModels.RerunInvestigationResult, *commonStructures.CommonError) {

	return rerunService.RerunDao.GetRerunByTestDetailsIds(testDetailsIds)
}

func (rerunService *RerunService) CreateRerunInvestigationResultsWithTx(tx *gorm.DB,
	rerunDetails []commonModels.RerunInvestigationResult) (
	[]commonModels.RerunInvestigationResult, *commonStructures.CommonError) {

	return rerunService.RerunDao.CreateRerunInvestigationResultsWithTx(tx, rerunDetails)
}

func (rerunService *RerunService) UpdateRerunInvestigationResultsWithTx(tx *gorm.DB,
	rerunDetails []commonModels.RerunInvestigationResult) (
	[]commonModels.RerunInvestigationResult, *commonStructures.CommonError) {
	if len(rerunDetails) > 0 {
		rr, err := rerunService.RerunDao.UpdateRerunInvestigationResultsWithTx(tx, rerunDetails)
		if err != nil {
			return []commonModels.RerunInvestigationResult{}, err
		}
		return rr, nil
	}
	return []commonModels.RerunInvestigationResult{}, nil
}
