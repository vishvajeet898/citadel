package mapper

import (
	commonStructures "github.com/Orange-Health/citadel/common/structures"
	commonUtils "github.com/Orange-Health/citadel/common/utils"
	"github.com/Orange-Health/citadel/models"
)

func MapTsmgModelToTsmInfo(tsmm models.TestSampleMapping) commonStructures.TestSampleMappingInfo {
	return commonStructures.TestSampleMappingInfo{
		Id:                  tsmm.Id,
		OmsTestId:           tsmm.OmsTestId,
		SampleId:            tsmm.SampleId,
		SampleNumber:        tsmm.SampleNumber,
		OmsOrderId:          tsmm.OmsOrderId,
		IsRejected:          tsmm.IsRejected,
		RecollectionPending: tsmm.RecollectionPending,
		RejectionReason:     tsmm.RejectionReason,
		CreatedAt:           tsmm.CreatedAt,
		CreatedBy:           tsmm.CreatedBy,
		UpdatedAt:           tsmm.UpdatedAt,
		UpdatedBy:           tsmm.UpdatedBy,
		DeletedAt:           commonUtils.GetGoLangTimeFromGormDeletedAt(tsmm.DeletedAt),
		DeletedBy:           tsmm.DeletedBy,
	}
}

func MapBulkTsmgModelToTsmInfo(tsmms []models.TestSampleMapping) []commonStructures.TestSampleMappingInfo {
	tsmis := []commonStructures.TestSampleMappingInfo{}
	for _, tsmm := range tsmms {
		tsmis = append(tsmis, MapTsmgModelToTsmInfo(tsmm))
	}
	return tsmis
}
