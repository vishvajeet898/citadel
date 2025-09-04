package dao

import (
	"gorm.io/gorm"

	commonStructures "github.com/Orange-Health/citadel/common/structures"
	commonUtils "github.com/Orange-Health/citadel/common/utils"
	commonModels "github.com/Orange-Health/citadel/models"
)

type DataLayer interface {
	GetPatientDetailsById(id uint) (commonModels.PatientDetail, *commonStructures.CommonError)
	GetPatientDetailsBySystemPatientIds(systemPatientId []string) (
		*[]commonModels.PatientDetail, *commonStructures.CommonError)

	UpdatePatientDetails(patientDetails commonModels.PatientDetail) (
		commonModels.PatientDetail, *commonStructures.CommonError)
	CreatePatientDetailsWithTx(tx *gorm.DB, patientDetail commonModels.PatientDetail) (
		commonModels.PatientDetail, *commonStructures.CommonError)
	UpdatePatientDetailsWithTx(tx *gorm.DB, patientDetails commonModels.PatientDetail) (
		commonModels.PatientDetail, *commonStructures.CommonError)
}

func (patientDetailDao *PatientDetailDao) GetPatientDetailsById(id uint) (
	commonModels.PatientDetail, *commonStructures.CommonError) {

	var patientDetails commonModels.PatientDetail
	err := patientDetailDao.Db.Where("id = ?", id).First(&patientDetails).Error
	if err != nil {
		return patientDetails, commonUtils.HandleORMError(err)
	}

	return patientDetails, nil
}

func (patientDetailDao *PatientDetailDao) GetPatientDetailsBySystemPatientIds(systemPatientIds []string) (
	*[]commonModels.PatientDetail, *commonStructures.CommonError) {

	var patientDetails []commonModels.PatientDetail
	err := patientDetailDao.Db.Where("system_patient_id in ?", systemPatientIds).Find(&patientDetails).Error
	if err != nil {
		return nil, commonUtils.HandleORMError(err)
	}

	return &patientDetails, nil
}

func (patientDetailDao *PatientDetailDao) UpdatePatientDetails(patientDetails commonModels.PatientDetail) (
	commonModels.PatientDetail, *commonStructures.CommonError) {

	if err := patientDetailDao.Db.Save(&patientDetails).Error; err != nil {
		return patientDetails, commonUtils.HandleORMError(err)
	}

	return patientDetails, nil
}

func (patientDetailDao *PatientDetailDao) CreatePatientDetailsWithTx(tx *gorm.DB,
	patientDetail commonModels.PatientDetail) (
	commonModels.PatientDetail, *commonStructures.CommonError) {

	if err := tx.Create(&patientDetail).Error; err != nil {
		return patientDetail, commonUtils.HandleORMError(err)
	}

	return patientDetail, nil
}

func (patientDetailDao *PatientDetailDao) UpdatePatientDetailsWithTx(tx *gorm.DB,
	patientDetails commonModels.PatientDetail) (
	commonModels.PatientDetail, *commonStructures.CommonError) {

	if err := tx.Save(&patientDetails).Error; err != nil {
		return patientDetails, commonUtils.HandleORMError(err)
	}

	return patientDetails, nil
}
