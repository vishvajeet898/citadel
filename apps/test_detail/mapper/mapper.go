package mapper

import (
	"github.com/Orange-Health/citadel/apps/test_detail/structures"
	commonUtils "github.com/Orange-Health/citadel/common/utils"
	commonModels "github.com/Orange-Health/citadel/models"
)

func MapTestDetail(testDetail commonModels.TestDetail) structures.TestDetail {
	var testDetailStruct structures.TestDetail
	testDetailStruct.ID = testDetail.Id
	testDetailStruct.TaskID = testDetail.TaskId
	testDetailStruct.OmsTestID = testDetail.OmsTestId
	testDetailStruct.CentralOmsTestId = testDetail.CentralOmsTestId
	testDetailStruct.MasterTestId = testDetail.MasterTestId
	testDetailStruct.MasterPackageId = testDetail.MasterPackageId
	testDetailStruct.OmsOrderId = testDetail.OmsOrderId
	testDetailStruct.TestName = testDetail.TestName
	testDetailStruct.Status = testDetail.Status
	testDetailStruct.DoctorTat = testDetail.DoctorTat
	testDetailStruct.IsAutoApproved = testDetail.IsAutoApproved
	testDetailStruct.ReportSentAt = testDetail.ReportSentAt
	testDetailStruct.CreatedAt = testDetail.CreatedAt
	testDetailStruct.UpdatedAt = testDetail.UpdatedAt
	testDetailStruct.DeletedAt = commonUtils.GetGoLangTimeFromGormDeletedAt(testDetail.DeletedAt)
	testDetailStruct.CreatedBy = testDetail.CreatedBy
	testDetailStruct.UpdatedBy = testDetail.UpdatedBy
	testDetailStruct.DeletedBy = testDetail.DeletedBy
	testDetailStruct.LabId = testDetail.LabId
	testDetailStruct.ProcessingLabId = testDetail.ProcessingLabId
	testDetailStruct.Department = testDetail.Department
	return testDetailStruct
}

func MapTestDetails(testDetails []commonModels.TestDetail) []structures.TestDetail {
	var testDetailsStruct []structures.TestDetail
	for _, testDetail := range testDetails {
		testDetailsStruct = append(testDetailsStruct, MapTestDetail(testDetail))
	}
	return testDetailsStruct
}
