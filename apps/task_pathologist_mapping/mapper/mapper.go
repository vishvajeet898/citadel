package mapper

import (
	"github.com/Orange-Health/citadel/apps/task_pathologist_mapping/structures"
	commonModels "github.com/Orange-Health/citadel/models"
)

func MapTpm(tpm commonModels.TaskPathologistMapping) structures.TaskPathologistMapping {
	var tpmStruct structures.TaskPathologistMapping
	tpmStruct.TaskID = tpm.TaskId
	tpmStruct.PathologistID = tpm.PathologistId
	tpmStruct.IsActive = tpm.IsActive
	return tpmStruct
}
