package mapper

import (
	"github.com/Orange-Health/citadel/apps/patient_details/structures"
	commonUtils "github.com/Orange-Health/citadel/common/utils"
	commonModels "github.com/Orange-Health/citadel/models"
)

func MapPatientDetails(patientDetail commonModels.PatientDetail) structures.PatientDetail {

	patientDetailStruct := structures.PatientDetail{}

	patientDetailStruct.ID = patientDetail.Id
	patientDetailStruct.Name = patientDetail.Name
	patientDetailStruct.ExpectedDob = patientDetail.ExpectedDob
	patientDetailStruct.Dob = patientDetail.Dob
	patientDetailStruct.Gender = patientDetail.Gender
	// patientDetailStruct.Number = commonUtils.MaskPhoneNumber(patientDetail.Number)
	patientDetailStruct.Number = patientDetail.Number
	patientDetailStruct.SystemPatientID = patientDetail.SystemPatientId
	patientDetailStruct.CreatedAt = patientDetail.CreatedAt
	patientDetailStruct.UpdatedAt = patientDetail.UpdatedAt
	patientDetailStruct.DeletedAt = commonUtils.GetGoLangTimeFromGormDeletedAt(patientDetail.DeletedAt)
	patientDetailStruct.CreatedBy = patientDetail.CreatedBy
	patientDetailStruct.UpdatedBy = patientDetail.UpdatedBy
	patientDetailStruct.DeletedBy = patientDetail.DeletedBy

	return patientDetailStruct
}
