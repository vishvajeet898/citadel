package mapper

import (
	"github.com/Orange-Health/citadel/apps/remarks/structures"
	commonUtils "github.com/Orange-Health/citadel/common/utils"
	"github.com/Orange-Health/citadel/models"
)

func MapRemark(remarks models.Remark) structures.Remark {
	var mappedRemark structures.Remark
	mappedRemark.ID = remarks.Id
	mappedRemark.InvestigationResultId = remarks.InvestigationResultId
	mappedRemark.Description = remarks.Description
	mappedRemark.RemarkType = remarks.RemarkType
	mappedRemark.RemarkBy = remarks.RemarkBy
	mappedRemark.CreatedAt = remarks.CreatedAt
	mappedRemark.UpdatedAt = remarks.UpdatedAt
	mappedRemark.DeletedAt = commonUtils.GetGoLangTimeFromGormDeletedAt(remarks.DeletedAt)
	mappedRemark.CreatedBy = remarks.CreatedBy
	mappedRemark.UpdatedBy = remarks.UpdatedBy
	mappedRemark.DeletedBy = remarks.DeletedBy
	return mappedRemark
}

func MapRemarks(remarks []models.Remark) []structures.Remark {
	var mappedRemarks []structures.Remark
	for _, remark := range remarks {
		mappedRemarks = append(mappedRemarks, MapRemark(remark))
	}
	return mappedRemarks
}
