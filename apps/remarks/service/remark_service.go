package service

import (
	"gorm.io/gorm"

	mapper "github.com/Orange-Health/citadel/apps/remarks/mapper"
	"github.com/Orange-Health/citadel/apps/remarks/structures"
	commonStructures "github.com/Orange-Health/citadel/common/structures"
	commonModels "github.com/Orange-Health/citadel/models"
)

type RemarkServiceInterface interface {
	GetRemarksStructByTaskId(taskId uint) (
		[]structures.Remark, *commonStructures.CommonError)
	GetRemarkByTaskId(taskId uint) (
		[]commonModels.Remark, *commonStructures.CommonError)
	GetRemarksByRemarkIds(remarkIds []uint) (
		[]commonModels.Remark, *commonStructures.CommonError)
	GetRemarksByInvestigationResultIds(remarkTypes []string,
		investigationResultIds []uint) ([]commonModels.Remark, *commonStructures.CommonError)

	CreateRemarksWithTx(tx *gorm.DB, remarks []commonModels.Remark) (
		[]commonModels.Remark, *commonStructures.CommonError)
	UpdateRemarksWithTx(tx *gorm.DB, remarks []commonModels.Remark) (
		[]commonModels.Remark, *commonStructures.CommonError)
	DeleteRemarksWithTx(tx *gorm.DB, remarkIds []uint, userId uint) *commonStructures.CommonError
}

func (remarkService *RemarkService) GetRemarksByRemarkIds(remarkIds []uint) (
	[]commonModels.Remark, *commonStructures.CommonError) {
	remarks, err := remarkService.RemarkDao.GetRemarksByRemarkIds(remarkIds)
	if err != nil {
		return []commonModels.Remark{}, err
	}

	return remarks, nil
}

func (remarkService *RemarkService) GetRemarksByInvestigationResultIds(remarkTypes []string,
	investigationResultIds []uint) ([]commonModels.Remark, *commonStructures.CommonError) {
	if len(investigationResultIds) == 0 {
		return []commonModels.Remark{}, nil
	}

	remarks, err := remarkService.RemarkDao.GetRemarksByInvestigationResultIds(remarkTypes, investigationResultIds)
	if err != nil {
		return []commonModels.Remark{}, err
	}

	return remarks, nil
}

func (remarkService *RemarkService) GetRemarksStructByTaskId(taskId uint) (
	[]structures.Remark, *commonStructures.CommonError) {

	remarks, err := remarkService.RemarkDao.GetRemarksByTaskId(taskId)
	if err != nil {
		return []structures.Remark{}, err
	}

	return mapper.MapRemarks(remarks), nil
}

func (remarkService *RemarkService) GetRemarkByTaskId(taskId uint) (
	[]commonModels.Remark, *commonStructures.CommonError) {

	remarks, err := remarkService.RemarkDao.GetRemarksByTaskId(taskId)
	if err != nil {
		return []commonModels.Remark{}, err
	}

	return remarks, nil
}

func (remarkService *RemarkService) CreateRemarksWithTx(tx *gorm.DB, remarks []commonModels.Remark) (
	[]commonModels.Remark, *commonStructures.CommonError) {
	if len(remarks) == 0 {
		return []commonModels.Remark{}, nil
	}

	createdRemarks, err := remarkService.RemarkDao.CreateRemarksWithTx(tx, remarks)
	if err != nil {
		return []commonModels.Remark{}, err
	}

	return createdRemarks, nil
}

func (remarkService *RemarkService) UpdateRemarksWithTx(tx *gorm.DB, remarks []commonModels.Remark) (
	[]commonModels.Remark, *commonStructures.CommonError) {
	if len(remarks) == 0 {
		return []commonModels.Remark{}, nil
	}

	updatedRemarks, err := remarkService.RemarkDao.UpdateRemarksWithTx(tx, remarks)
	if err != nil {
		return []commonModels.Remark{}, err
	}

	return updatedRemarks, nil
}

func (remarkService *RemarkService) DeleteRemarksWithTx(tx *gorm.DB, remarkIds []uint, userId uint) *commonStructures.CommonError {
	err := remarkService.RemarkDao.DeleteRemarksWithTx(tx, remarkIds, userId)
	return err
}
