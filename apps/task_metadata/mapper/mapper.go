package mapper

import (
	"time"

	"gorm.io/gorm"

	"github.com/Orange-Health/citadel/apps/task_metadata/structures"
	commonUtils "github.com/Orange-Health/citadel/common/utils"
	commonModels "github.com/Orange-Health/citadel/models"
)

func MapTmdModel(tmdStruct structures.TaskMetadata) commonModels.TaskMetadata {
	t := time.Now()
	createdAt := &t
	updatedAt := &t
	var deletedAt gorm.DeletedAt

	if tmdStruct.CreatedAt != nil {
		createdAt = tmdStruct.CreatedAt
	}

	if tmdStruct.UpdatedAt != nil {
		updatedAt = tmdStruct.UpdatedAt
	}

	if tmdStruct.DeletedAt != nil {
		deletedAt.Time = *tmdStruct.DeletedAt
	}

	var tmdModel commonModels.TaskMetadata
	tmdModel.Id = tmdStruct.ID
	tmdModel.TaskId = tmdStruct.TaskID
	tmdModel.ContainsMorphle = tmdStruct.IsMorphle
	tmdModel.ContainsPackage = tmdStruct.IsPackage
	tmdModel.IsCritical = tmdStruct.IsCritical
	tmdModel.DoctorName = tmdStruct.DoctorName
	tmdModel.DoctorNumber = tmdStruct.DoctorNumber
	tmdModel.DoctorNotes = tmdStruct.DoctorNotes
	tmdModel.CreatedAt = createdAt
	tmdModel.UpdatedAt = updatedAt
	tmdModel.DeletedAt = &deletedAt
	tmdModel.CreatedBy = tmdStruct.CreatedBy
	tmdModel.UpdatedBy = tmdStruct.UpdatedBy
	tmdModel.DeletedBy = tmdStruct.DeletedBy

	return tmdModel
}

func MapTmd(tmd commonModels.TaskMetadata) structures.TaskMetadata {
	tmdStruct := structures.TaskMetadata{
		TaskID:       tmd.TaskId,
		IsMorphle:    tmd.ContainsMorphle,
		IsPackage:    tmd.ContainsPackage,
		IsCritical:   tmd.IsCritical,
		DoctorName:   tmd.DoctorName,
		DoctorNumber: tmd.DoctorNumber,
		DoctorNotes:  tmd.DoctorNotes,
	}

	tmdStruct.ID = tmd.Id
	tmdStruct.CreatedAt = tmd.CreatedAt
	tmdStruct.UpdatedAt = tmd.UpdatedAt
	tmdStruct.DeletedAt = commonUtils.GetGoLangTimeFromGormDeletedAt(tmd.DeletedAt)
	tmdStruct.CreatedBy = tmd.CreatedBy
	tmdStruct.UpdatedBy = tmd.UpdatedBy
	tmdStruct.DeletedBy = tmd.DeletedBy

	return tmdStruct
}
