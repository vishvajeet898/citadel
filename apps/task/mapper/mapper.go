package mapper

import (
	patientDetailStruct "github.com/Orange-Health/citadel/apps/patient_details/structures"
	"github.com/Orange-Health/citadel/apps/task/structures"
	commonConstants "github.com/Orange-Health/citadel/common/constants"
	commonStructures "github.com/Orange-Health/citadel/common/structures"
	commonUtils "github.com/Orange-Health/citadel/common/utils"
	commonModels "github.com/Orange-Health/citadel/models"
)

func GetTaskDetails(task commonModels.Task, taskMetadata commonModels.TaskMetadata,
	visits []commonStructures.VisitDetailsForTask) (structures.TaskDetail, *commonStructures.CommonError) {

	var completedAtStr string
	var cErr *commonStructures.CommonError

	if task.CompletedAt != nil && !task.CompletedAt.IsZero() {
		completedAtStr = task.CompletedAt.Format(commonConstants.DateTimeLayout)
	}

	tds := MapTaskDetail(task, completedAtStr)

	tds.PatientDetails = MapPatientDetails(task.PatientDetail)

	tds.TaskMetadata = MapTaskMetadata(taskMetadata)

	tds.TaskVisits = visits
	tds.VisitCount = len(visits)

	return tds, cErr
}

func MapTaskDetail(task commonModels.Task, completedAtStr string) structures.TaskDetail {

	var taskStruct structures.TaskDetail
	taskStruct.ID = task.Id
	taskStruct.OrderID = task.OrderId
	taskStruct.OmsOrderId = task.OmsOrderId
	taskStruct.RequestID = task.OmsRequestId
	taskStruct.LabID = task.LabId
	taskStruct.CityCode = task.CityCode
	taskStruct.Status = task.Status
	taskStruct.OrderType = task.OrderType
	taskStruct.IsActive = task.IsActive
	taskStruct.CompletedAt = completedAtStr
	taskStruct.CreatedAt = task.CreatedAt
	taskStruct.UpdatedAt = task.UpdatedAt
	taskStruct.DeletedAt = commonUtils.GetGoLangTimeFromGormDeletedAt(task.DeletedAt)
	taskStruct.CreatedBy = task.CreatedBy
	taskStruct.UpdatedBy = task.UpdatedBy
	taskStruct.DeletedBy = task.DeletedBy

	return taskStruct
}

func MapTaskMetadata(taskMetadata commonModels.TaskMetadata) (doctorDetail structures.TaskMetadata) {
	var taskMetadataStruct structures.TaskMetadata
	taskMetadataStruct.ID = taskMetadata.Id
	taskMetadataStruct.TaskID = taskMetadata.TaskId
	taskMetadataStruct.IsPackage = taskMetadata.ContainsPackage
	taskMetadataStruct.IsMorphle = taskMetadata.ContainsMorphle
	taskMetadataStruct.IsCritical = taskMetadata.IsCritical
	taskMetadataStruct.DoctorName = taskMetadata.DoctorName
	// taskMetadataStruct.DoctorNumber = commonUtils.MaskPhoneNumber(taskMetadata.DoctorNumber)
	taskMetadataStruct.DoctorNumber = taskMetadata.DoctorNumber
	taskMetadataStruct.PartnerName = taskMetadata.PartnerName
	taskMetadataStruct.DoctorNotes = taskMetadata.DoctorNotes
	return taskMetadataStruct
}

func MapPatientDetails(patientDetail commonModels.PatientDetail) patientDetailStruct.PatientDetail {
	patientDetailsResp := patientDetailStruct.PatientDetail{}
	patientDetailsResp.ID = patientDetail.Id
	patientDetailsResp.Name = patientDetail.Name
	patientDetailsResp.ExpectedDob = patientDetail.ExpectedDob
	patientDetailsResp.Dob = patientDetail.Dob
	patientDetailsResp.Gender = patientDetail.Gender
	// patientDetailsResp.Number = commonUtils.MaskPhoneNumber(patientDetail.Number)
	patientDetailsResp.Number = patientDetail.Number
	patientDetailsResp.SystemPatientID = patientDetail.SystemPatientId

	return patientDetailsResp
}
