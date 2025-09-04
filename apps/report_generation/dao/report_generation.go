package dao

import (
	"context"

	"github.com/Orange-Health/citadel/apps/report_generation/structures"
	commonConstants "github.com/Orange-Health/citadel/common/constants"
	commonStructures "github.com/Orange-Health/citadel/common/structures"
	commonUtils "github.com/Orange-Health/citadel/common/utils"
	commonModels "github.com/Orange-Health/citadel/models"
)

type DataLayer interface {
	GetInvestigationDataForReportGeneration(ctx context.Context, omsOrderId string,
		omsTestIds []string) ([]structures.InvestigationEvent, *commonStructures.CommonError)
	GetUserDetails(ctx context.Context, userIds []string) ([]commonModels.User, *commonStructures.CommonError)
	FetchVisitToTestMapping(omsOrderId string, omsTestIds []string) ([]structures.VisitTestsStruct,
		*commonStructures.CommonError)
	GetOrderDetailsAndPatientDetailsByOmsOrderId(omsOrderId string) (commonModels.OrderDetails, commonModels.PatientDetail,
		*commonStructures.CommonError)
	GetPatientDetailsById(patientDetailsId uint) (commonModels.PatientDetail, *commonStructures.CommonError)
}

func getSelectQueryForReportGenerationData() []string {
	return []string{
		"investigation_results.id as id",
		"investigation_results.master_investigation_id as master_investigation_id",
		"investigation_results.result_representation_type as result_representation_type",
		"test_details.central_oms_test_id as test_id",
		"investigation_results.investigation_name as name",
		"investigation_results.lis_code as lis_code",
		"investigation_results.abnormality as abnormality",
		"investigation_results.department as department_name",
		"investigation_results.method as method",
		"investigation_results.investigation_value as value",
		"investigation_data.data as data",
		"investigation_results.uom as unit",
		"investigation_results.reference_range_text as reference_range_text",
		"investigation_results.investigation_status as status",
		"investigation_results.approved_at as approved_at",
		"investigation_results.approved_by as approved_by",
		"investigation_results.entered_at as result_entered_at",
		"investigation_results.entered_by as result_entered_by",
		"investigation_results.method as investigation_method",
		"investigation_results.is_abnormal as is_abnormal",
		"investigation_results.is_nabl_approved as is_nabl_approved",
		"investigation_results.source as source",
		"investigation_results.created_at as created_at",
		"investigation_results.updated_at as updated_at",
		"remarks.description as medical_remarks",
	}
}

func (dao *ReportGenerationDao) GetInvestigationDataForReportGeneration(ctx context.Context, orderId string,
	omsTestIds []string) ([]structures.InvestigationEvent, *commonStructures.CommonError) {

	investigations := []structures.InvestigationEvent{}

	err := dao.Db.WithContext(ctx).Select(getSelectQueryForReportGenerationData()).
		Model(&commonModels.InvestigationResult{}).
		Joins("INNER JOIN test_details ON investigation_results.test_details_id = test_details.id").
		Joins("INNER JOIN test_details_metadata ON test_details.id = test_details_metadata.test_details_id").
		Joins("INNER JOIN tasks ON test_details.task_id = tasks.id").
		Joins("LEFT JOIN investigation_data ON investigation_results.id = investigation_data.investigation_result_id AND investigation_data.deleted_at IS NULL").
		Joins("LEFT JOIN remarks ON investigation_results.id = remarks.investigation_result_id AND remarks.remark_type = ? AND remarks.deleted_at IS NULL", commonConstants.REMARK_TYPE_MEDICAL_REMARK).
		Where("test_details.central_oms_test_id IN (?) and test_details.status = ?", omsTestIds, commonConstants.TEST_STATUS_APPROVE).
		Where("tasks.oms_order_id = ?", orderId).
		Where("investigation_results.investigation_status IN (?)", commonConstants.INVESTIGATION_STATUSES_APPROVE).
		Scan(&investigations).Error

	if err != nil {
		return nil, commonUtils.HandleORMError(err)
	}

	return investigations, nil
}

func (dao *ReportGenerationDao) GetUserDetails(ctx context.Context, userIds []string) (
	[]commonModels.User, *commonStructures.CommonError) {

	users := []commonModels.User{}
	if err := dao.Db.WithContext(ctx).Where("id IN (?)", userIds).Find(&users).Error; err != nil {
		return nil, commonUtils.HandleORMError(err)
	}

	return users, nil
}

func (dao *ReportGenerationDao) FetchVisitToTestMapping(orderId string, omsTestIds []string) ([]structures.VisitTestsStruct,
	*commonStructures.CommonError) {
	visitIdTests := []structures.VisitTestsStruct{}
	testStatusStrings := []string{commonConstants.TEST_STATUS_COMPLETED_NOT_SENT, commonConstants.TEST_STATUS_COMPLETED_SENT,
		commonConstants.TEST_STATUS_APPROVE}
	selectStrings := []string{
		"samples.visit_id",
		"samples.barcode",
		"sample_metadata.received_at",
		"sample_metadata.collected_at",
		"samples.sample_number",
		"samples.vial_type_id",
		"test_details.department as department_name",
		"test_details.*",
	}

	query := dao.Db.Table(commonConstants.TableSamples).
		Select(selectStrings).
		Joins("INNER JOIN sample_metadata ON samples.id = sample_metadata.sample_id").
		Joins("INNER JOIN test_sample_mapping ON samples.oms_order_id = test_sample_mapping.oms_order_id AND samples.sample_number = test_sample_mapping.sample_number").
		Joins("INNER JOIN test_details ON test_sample_mapping.oms_test_id = test_details.central_oms_test_id").
		Where("samples.oms_order_id = ?", orderId).
		Where("samples.visit_id != ''").
		Where("samples.deleted_at IS NULL").
		Where("sample_metadata.deleted_at IS NULL").
		Where("test_sample_mapping.deleted_at IS NULL").
		Where("test_details.deleted_at IS NULL")

	if len(omsTestIds) > 0 {
		query = query.Where("test_details.status IN (?) OR test_details.central_oms_test_id in (?)", testStatusStrings,
			omsTestIds)
	} else {
		query = query.Where("test_details.status IN (?)", testStatusStrings)
	}

	if err := query.Scan(&visitIdTests).Error; err != nil {
		return visitIdTests, commonUtils.HandleORMError(err)
	}

	return visitIdTests, nil
}

func (dao *ReportGenerationDao) GetOrderDetailsAndPatientDetailsByOmsOrderId(omsOrderId string) (commonModels.OrderDetails,
	commonModels.PatientDetail, *commonStructures.CommonError) {
	orderDetails, patientDetails := commonModels.OrderDetails{}, commonModels.PatientDetail{}
	if err := dao.Db.Table(commonConstants.TableOrderDetails).
		Find(&orderDetails, "oms_order_id = ?", omsOrderId).Error; err != nil {
		return orderDetails, patientDetails, commonUtils.HandleORMError(err)
	}

	if err := dao.Db.Table(commonConstants.TablePatientDetails).
		Find(&patientDetails, "id = ?", orderDetails.PatientDetailsId).Error; err != nil {
		return orderDetails, patientDetails, commonUtils.HandleORMError(err)
	}

	return orderDetails, patientDetails, nil
}

func (dao *ReportGenerationDao) GetPatientDetailsById(patientDetailsId uint) (commonModels.PatientDetail,
	*commonStructures.CommonError) {
	patientDetails := commonModels.PatientDetail{}
	if err := dao.Db.Table(commonConstants.TablePatientDetails).
		Find(&patientDetails, "id = ?", patientDetailsId).Error; err != nil {
		return patientDetails, commonUtils.HandleORMError(err)
	}

	return patientDetails, nil
}
