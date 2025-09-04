package service

import (
	"gorm.io/gorm"

	commonStructures "github.com/Orange-Health/citadel/common/structures"
	commonModels "github.com/Orange-Health/citadel/models"
)

type CoAuthorizePathologistInterface interface {
	CreateCoAuthrizePathologistWithTx(tx *gorm.DB, coAuthorizePathologistDetail commonModels.CoAuthorizedPathologists) (
		commonModels.CoAuthorizedPathologists, *commonStructures.CommonError)
	DeleteCurrentCoAuthorizePathologistWithTx(tx *gorm.DB, taskId, userId uint) *commonStructures.CommonError
}

func (coAuthorizePathologistService *CoAuthorizePathologistService) CreateCoAuthrizePathologistWithTx(tx *gorm.DB,
	coAuthorizePathologistDetail commonModels.CoAuthorizedPathologists) (
	commonModels.CoAuthorizedPathologists, *commonStructures.CommonError) {

	coAuthorizePathologist, err := coAuthorizePathologistService.CoAuthorizePathDao.CreateCoAuthrizePathologistWithTx(tx, coAuthorizePathologistDetail)
	if err != nil {
		return commonModels.CoAuthorizedPathologists{}, err
	}

	return coAuthorizePathologist, nil
}

func (coAuthorizePathologistService *CoAuthorizePathologistService) DeleteCurrentCoAuthorizePathologistWithTx(tx *gorm.DB,
	taskId, userId uint) *commonStructures.CommonError {

	err := coAuthorizePathologistService.CoAuthorizePathDao.DeleteCurrentCoAuthorizePathologistWithTx(tx, taskId, userId)
	if err != nil {
		return err
	}

	return nil
}
