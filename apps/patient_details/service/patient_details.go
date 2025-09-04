package service

import (
	"context"
	"encoding/json"
	"net/http"

	"gorm.io/gorm"

	mapper "github.com/Orange-Health/citadel/apps/patient_details/mapper"
	"github.com/Orange-Health/citadel/apps/patient_details/structures"
	commonConstants "github.com/Orange-Health/citadel/common/constants"
	commonStructures "github.com/Orange-Health/citadel/common/structures"
	commonUtils "github.com/Orange-Health/citadel/common/utils"
	commonModels "github.com/Orange-Health/citadel/models"
)

type PatientDetailServiceInterface interface {
	GetPatientDetailsResponseById(id uint) (
		structures.PatientDetail, *commonStructures.CommonError)
	GetPatientDetailsById(patientDetailsId uint) (
		commonModels.PatientDetail, *commonStructures.CommonError)
	GetPatientDetailsBySystemPatientIds(systemPatientId []string) (
		*[]commonModels.PatientDetail, *commonStructures.CommonError)

	CreatePatientDetailsWithTx(tx *gorm.DB, patientDetail commonModels.PatientDetail) (
		commonModels.PatientDetail, *commonStructures.CommonError)
	UpdatePatientDetailsWithTx(tx *gorm.DB, patientDetail commonModels.PatientDetail) (
		commonModels.PatientDetail, *commonStructures.CommonError)
	UpdatePatientDetails(ctx context.Context, patientDetailsRequest structures.PatientDetailsUpdateRequest) (
		structures.PatientDetail, *commonStructures.CommonError)
}

func (patientDetailService *PatientDetailService) GetPatientDetailsResponseById(id uint) (
	structures.PatientDetail, *commonStructures.CommonError) {

	patientDetails, err := patientDetailService.PatientDetailDao.GetPatientDetailsById(id)
	if err != nil {
		return structures.PatientDetail{}, err
	}

	return mapper.MapPatientDetails(patientDetails), nil
}

func (patientDetailService *PatientDetailService) GetPatientDetailsBySystemPatientIds(

	systemPatientId []string,
) (*[]commonModels.PatientDetail, *commonStructures.CommonError) {
	return patientDetailService.PatientDetailDao.GetPatientDetailsBySystemPatientIds(systemPatientId)
}

func (patientDetailService *PatientDetailService) GetFieldsToBeUpdated(currentPatientDetails *commonModels.PatientDetail, patientDetailsRequest structures.PatientDetailsUpdateRequest) {
	if patientDetailsRequest.Name != "" {
		currentPatientDetails.Name = patientDetailsRequest.Name
	}
	if patientDetailsRequest.Gender != "" {
		currentPatientDetails.Gender = patientDetailsRequest.Gender
	}
	if patientDetailsRequest.ExpectedDob != "" {
		expectedDob := commonUtils.StringToTime(patientDetailsRequest.ExpectedDob, commonConstants.DateLayout)
		currentPatientDetails.ExpectedDob = &expectedDob
	}
	if patientDetailsRequest.Dob != "" {
		dob := commonUtils.StringToTime(patientDetailsRequest.Dob, commonConstants.DateLayout)
		currentPatientDetails.Dob = &dob
		currentPatientDetails.ExpectedDob = &dob
	}
	if patientDetailsRequest.Number != "" {
		currentPatientDetails.Number = patientDetailsRequest.Number
	}
}

func (PatientDetailService *PatientDetailService) CreateRequestForUpdatingPatientInOms(
	patientDetails commonModels.PatientDetail) structures.OmsPatientDetailsUpdateRequest {
	omsPatientDetailsRequest := structures.OmsPatientDetailsUpdateRequest{
		Name:   patientDetails.Name,
		Gender: patientDetails.Gender,
		Number: patientDetails.Number,
	}

	ageYears, ageMonths, ageDays := commonUtils.GetAgeYearsMonthsAndDaysFromDob(*patientDetails.ExpectedDob)
	if patientDetails.Dob != nil {
		ageYears, ageMonths, ageDays = commonUtils.GetAgeYearsMonthsAndDaysFromDob(*patientDetails.Dob)
	}
	omsPatientDetailsRequest.AgeYears = ageYears
	omsPatientDetailsRequest.AgeMonths = ageMonths
	omsPatientDetailsRequest.AgeDays = ageDays

	return omsPatientDetailsRequest
}

func (patientDetailService *PatientDetailService) UpdatePatientDetails(ctx context.Context,
	patientDetailsRequest structures.PatientDetailsUpdateRequest) (
	structures.PatientDetail, *commonStructures.CommonError) {

	currentPatientDetails, cErr := patientDetailService.PatientDetailDao.GetPatientDetailsById(patientDetailsRequest.Id)
	if cErr != nil {
		return structures.PatientDetail{}, cErr
	}

	patientDetailService.GetFieldsToBeUpdated(&currentPatientDetails, patientDetailsRequest)

	omsPatientDetails := patientDetailService.CreateRequestForUpdatingPatientInOms(currentPatientDetails)

	omsResponse, err := patientDetailService.OmsClient.UpdatePatientDetails(ctx, patientDetailsRequest.OrderId,
		patientDetailsRequest.CityCode, omsPatientDetails)
	if err != nil {
		return structures.PatientDetail{}, &commonStructures.CommonError{
			Message:    err.Error(),
			StatusCode: http.StatusInternalServerError,
		}
	}

	omsResponseBytes, err := json.Marshal(omsResponse)
	if err != nil {
		return structures.PatientDetail{}, &commonStructures.CommonError{
			Message:    err.Error(),
			StatusCode: http.StatusInternalServerError,
		}
	}

	omsPatientDetailsResponse := structures.OmsPatientDetailsUpdateRequest{}
	err = json.Unmarshal(omsResponseBytes, &omsPatientDetailsResponse)
	if err != nil {
		return structures.PatientDetail{}, &commonStructures.CommonError{
			Message:    err.Error(),
			StatusCode: http.StatusInternalServerError,
		}
	}

	currentPatientDetails.SystemPatientId = omsPatientDetailsResponse.SystemPatientId

	updatedPatientDetails, cErr := patientDetailService.PatientDetailDao.UpdatePatientDetails(currentPatientDetails)
	if cErr != nil {
		return structures.PatientDetail{}, cErr
	}

	return mapper.MapPatientDetails(updatedPatientDetails), nil
}

func (patientDetailService *PatientDetailService) GetPatientDetailsById(patientDetailsId uint) (
	commonModels.PatientDetail, *commonStructures.CommonError) {
	return patientDetailService.PatientDetailDao.GetPatientDetailsById(patientDetailsId)
}

func (patientDetailService *PatientDetailService) CreatePatientDetailsWithTx(tx *gorm.DB, patientDetail commonModels.PatientDetail) (commonModels.PatientDetail, *commonStructures.CommonError) {
	return patientDetailService.PatientDetailDao.CreatePatientDetailsWithTx(tx, patientDetail)
}

func (patientDetailService *PatientDetailService) UpdatePatientDetailsWithTx(tx *gorm.DB, patientDetail commonModels.PatientDetail) (commonModels.PatientDetail, *commonStructures.CommonError) {
	return patientDetailService.PatientDetailDao.UpdatePatientDetailsWithTx(tx, patientDetail)
}
