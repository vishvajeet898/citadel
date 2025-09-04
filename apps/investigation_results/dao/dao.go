package dao

import (
	"gorm.io/gorm"

	"github.com/Orange-Health/citadel/apps/investigation_results/structures"
	commonConstants "github.com/Orange-Health/citadel/common/constants"
	commonStructures "github.com/Orange-Health/citadel/common/structures"
	commonUtils "github.com/Orange-Health/citadel/common/utils"
	commonModels "github.com/Orange-Health/citadel/models"
)

type DataLayer interface {
	GetInvestigationsByInvestigationIds(investigationIds []uint) (
		[]commonModels.InvestigationResult, *commonStructures.CommonError)
	GetInvestigationDataByInvestigationResultsIds(invResultIds []uint) (
		[]commonModels.InvestigationData, *commonStructures.CommonError)
	GetPatientPastRecords(patientIds []string) (
		[]commonStructures.PatientPastRecords, *commonStructures.CommonError)
	GetDeltaValuesFromPatientIds(patientIds []string, masterInvestigationIds []uint) (
		[]commonStructures.DeltaValuesStruct, *commonStructures.CommonError)
	GetInvestigationResultsByTestDetailsIds(testDetailsIds []uint) (
		[]commonModels.InvestigationResult, *commonStructures.CommonError)
	GetInvestigationResultsDbStructsByTaskId(taskId uint) (
		[]structures.InvestigationResultDbResponse, *commonStructures.CommonError)
	GetInvestigationResultsDbStructsByTestDetailIds(testDetailIds []uint) (
		[]structures.InvestigationResultDbResponse, *commonStructures.CommonError)
	GetInvestigationResultModelsByTaskId(taskId uint) (
		[]commonModels.InvestigationResult, *commonStructures.CommonError)
	GetDetailsForAbnormality(taskId uint) (
		structures.BasicAbnormalityStruct, *commonStructures.CommonError)
	GetInvestigationResultsByTaskIdAndOmsTestId(taskId uint, omsTestId string) (
		[]commonModels.InvestigationResult, *commonStructures.CommonError)
	GetTestDetailsByTaskId(taskId uint) (
		[]structures.TestDetailsDbResponse, *commonStructures.CommonError)
	GetInvestigationResultsMetadataByInvestigationResultsIds(invResultIds []uint) (
		[]commonModels.InvestigationResultMetadata, *commonStructures.CommonError)

	CreateInvestigationResultsWithTx(tx *gorm.DB,
		investigationResults []commonModels.InvestigationResult) (
		[]commonModels.InvestigationResult, *commonStructures.CommonError)
	UpdateInvestigationResultsWithTx(tx *gorm.DB,
		investigationResults []commonModels.InvestigationResult) (
		[]commonModels.InvestigationResult, *commonStructures.CommonError)
	CreateInvestigationDataWithTx(tx *gorm.DB,
		invData []commonModels.InvestigationData) (
		[]commonModels.InvestigationData, *commonStructures.CommonError)
	DeleteInvestigationResultsByIdsWithTx(tx *gorm.DB,
		investigationIds []uint) *commonStructures.CommonError
	UpdateInvestigationsDataWithTx(tx *gorm.DB,
		investigationData []commonModels.InvestigationData) (
		[]commonModels.InvestigationData, *commonStructures.CommonError)
	CreateInvestigationResultsMetadataWithTx(tx *gorm.DB,
		investigationResultsMetadata []commonModels.InvestigationResultMetadata) (
		[]commonModels.InvestigationResultMetadata, *commonStructures.CommonError)
	UpdateInvestigationResultsMetadataWithTx(tx *gorm.DB,
		investigationResultsMetadata []commonModels.InvestigationResultMetadata) (
		[]commonModels.InvestigationResultMetadata, *commonStructures.CommonError)
	DeleteInvestigationResultsMetadataByIdsWithTx(tx *gorm.DB,
		investigationIds []uint) *commonStructures.CommonError
}

func getSelectStringsForInvestigationResults() []string {
	return []string{
		"investigation_results.id as id",
		"investigation_results.test_details_id as test_details_id",
		"investigation_results.master_investigation_id as master_investigation_id",
		"investigation_results.master_investigation_method_mapping_id as master_investigation_method_mapping_id",
		"investigation_results.investigation_name as investigation_name",
		"investigation_results.investigation_value as investigation_value",
		"investigation_results.device_value as device_value",
		"investigation_results.result_representation_type as result_representation_type",
		"investigation_results.department as department",
		"investigation_results.uom as uom",
		"investigation_results.method as method",
		"investigation_results.method_type as method_type",
		"investigation_results.investigation_status as investigation_status",
		"investigation_results.reference_range_text as reference_range_text",
		"investigation_results.lis_code as lis_code",
		"investigation_results.abnormality as abnormality",
		"investigation_results.is_abnormal as is_abnormal",
		"investigation_results.approved_by as approved_by",
		"investigation_results.approved_at as approved_at",
		"investigation_results.entered_by as entered_by",
		"investigation_results.entered_at as entered_at",
		"investigation_results.is_auto_approved as is_auto_approved",
		"investigation_results.is_critical as is_critical",
		"investigation_results.auto_approval_failure_reason as auto_approval_failure_reason",
		"investigation_results_metadata.qc_flag as qc_flag",
		"investigation_results_metadata.qc_lot_number as qc_lot_number",
		"investigation_results_metadata.qc_value as qc_value",
		"investigation_results_metadata.qc_west_gard_warning as qc_west_gard_warning",
		"investigation_results_metadata.qc_status as qc_status",
		"test_details.master_test_id as master_test_id",
		"test_details.status as test_details_status",
		"test_details.lab_id as processing_lab_id",
		"test_details_metadata.barcodes as barcodes",
	}
}

func getSelectStringsForInvestigationResultsV1() []string {
	return []string{
		"id", "test_details_id", "master_investigation_id", "master_investigation_method_mapping_id", "investigation_name",
		"investigation_value", "device_value", "result_representation_type", "department", "uom", "method",
		"method_type", "investigation_status", "reference_range_text", "lis_code", "abnormality", "is_abnormal",
		"approved_by", "approved_at", "entered_by", "entered_at", "is_auto_approved", "is_critical",
		"auto_approval_failure_reason",
	}
}

func (ird *InvestigationResultDao) GetInvestigationsByInvestigationIds(
	investigationIds []uint) ([]commonModels.InvestigationResult, *commonStructures.CommonError) {

	var invResults []commonModels.InvestigationResult
	err := ird.Db.Where("id IN (?)", investigationIds).Find(&invResults).Error
	if err != nil {
		return nil, commonUtils.HandleORMError(err)
	}
	return invResults, nil
}

func (ird *InvestigationResultDao) GetPatientPastRecords(
	patientIds []string) ([]commonStructures.PatientPastRecords, *commonStructures.CommonError) {

	var patientPastRecords []commonStructures.PatientPastRecords
	selectStrings := []string{
		"investigation_results.investigation_value as investigation_value",
		"investigation_results.uom as uom",
		"investigation_results.master_investigation_id as master_investigation_id",
		"investigation_results.department as department",
		"investigation_results.investigation_name as investigation_name",
		"investigation_results.approved_at as approved_at",
		"tasks.order_id as order_id",
		"tasks.city_code as city_code",
	}
	if err := ird.Db.Table(commonConstants.TableInvestigationResults).
		Select(selectStrings).
		Joins("INNER JOIN test_details ON test_details.id = investigation_results.test_details_id").
		Joins("INNER JOIN tasks ON tasks.id = test_details.task_id").
		Joins("INNER JOIN patient_details ON tasks.patient_details_id = patient_details.id").
		Where("patient_details.system_patient_id IN (?)", patientIds).
		Where("tasks.status = ?", commonConstants.TASK_STATUS_COMPLETED).
		Where("investigation_results.investigation_status IN (?)", commonConstants.INVESTIGATION_STATUSES_APPROVE).
		Where("investigation_results.approved_at IS NOT NULL").
		Find(&patientPastRecords).Error; err != nil {
		return nil, commonUtils.HandleORMError(err)
	}

	return patientPastRecords, nil
}

func (ird *InvestigationResultDao) GetDeltaValuesFromPatientIds(patientIds []string,
	masterInvestigationIds []uint) ([]commonStructures.DeltaValuesStruct, *commonStructures.CommonError) {

	deltaValuesStructs := []commonStructures.DeltaValuesStruct{}
	selectStrings := []string{
		"investigation_results.investigation_value as investigation_value",
		"investigation_results.approved_at as approved_at",
		"investigation_results.master_investigation_id as master_investigation_id",
		"tasks.order_id as order_id",
		"tasks.city_code as city_code",
	}
	err := ird.Db.Table(commonConstants.TableInvestigationResults).
		Select(selectStrings).
		Joins("INNER JOIN test_details ON investigation_results.test_details_id = test_details.id").
		Joins("INNER JOIN tasks ON tasks.id = test_details.task_id").
		Joins("INNER JOIN patient_details ON tasks.patient_details_id = patient_details.id").
		Where("patient_details.system_patient_id IN (?)", patientIds).
		Where("tasks.status = ?", commonConstants.TASK_STATUS_COMPLETED).
		Where("investigation_results.master_investigation_id IN (?)", masterInvestigationIds).
		Where("investigation_results.investigation_status IN (?)", commonConstants.INVESTIGATION_STATUSES_APPROVE).
		Scan(&deltaValuesStructs).Error
	if err != nil {
		return deltaValuesStructs, commonUtils.HandleORMError(err)
	}
	return deltaValuesStructs, nil
}

func (ird *InvestigationResultDao) GetInvestigationResultsByTestDetailsIds(testDetailsIds []uint) (
	[]commonModels.InvestigationResult, *commonStructures.CommonError) {

	investigationResults := []commonModels.InvestigationResult{}
	err := ird.Db.Find(&investigationResults, "test_details_id IN (?)", testDetailsIds).Error
	if err != nil {
		return []commonModels.InvestigationResult{}, commonUtils.HandleORMError(err)
	}
	return investigationResults, nil
}

func (ird *InvestigationResultDao) GetInvestigationResultsMetadataByInvestigationResultsIds(invResultIds []uint) (
	[]commonModels.InvestigationResultMetadata, *commonStructures.CommonError) {
	investigationResultsMetadata := []commonModels.InvestigationResultMetadata{}
	err := ird.Db.Where("investigation_result_id IN (?)", invResultIds).Find(&investigationResultsMetadata).Error
	if err != nil {
		return investigationResultsMetadata, commonUtils.HandleORMError(err)
	}
	return investigationResultsMetadata, nil
}

func (ird *InvestigationResultDao) CreateInvestigationResultsWithTx(tx *gorm.DB,
	investigationResults []commonModels.InvestigationResult) (
	[]commonModels.InvestigationResult, *commonStructures.CommonError) {

	err := tx.Create(&investigationResults).Error
	if err != nil {
		return []commonModels.InvestigationResult{}, commonUtils.HandleORMError(err)
	}
	return investigationResults, nil
}

func (ird *InvestigationResultDao) UpdateInvestigationResultsWithTx(tx *gorm.DB,
	investigationResults []commonModels.InvestigationResult) (
	[]commonModels.InvestigationResult, *commonStructures.CommonError) {

	err := tx.Save(&investigationResults).Error
	if err != nil {
		return []commonModels.InvestigationResult{}, commonUtils.HandleORMError(err)
	}
	return investigationResults, nil
}

func (ird *InvestigationResultDao) DeleteInvestigationResultsByIdsWithTx(tx *gorm.DB,
	investigationIds []uint) *commonStructures.CommonError {

	currentTime := commonUtils.GetCurrentTime()
	updates := map[string]interface{}{
		"deleted_at": currentTime,
		"updated_at": currentTime,
		"deleted_by": commonConstants.CitadelSystemId,
		"updated_by": commonConstants.CitadelSystemId,
	}

	err := tx.Model(&commonModels.InvestigationResult{}).
		Where("id IN (?)", investigationIds).
		Updates(updates).Error
	if err != nil {
		return commonUtils.HandleORMError(err)
	}

	return nil
}

func (ird *InvestigationResultDao) UpdateInvestigationsDataWithTx(tx *gorm.DB,
	investigationData []commonModels.InvestigationData) ([]commonModels.InvestigationData, *commonStructures.CommonError) {

	err := tx.Save(&investigationData).Error
	if err != nil {
		return nil, commonUtils.HandleORMError(err)
	}
	return investigationData, nil
}

func (ird *InvestigationResultDao) GetInvestigationResultsDbStructsByTaskId(taskId uint) (
	[]structures.InvestigationResultDbResponse, *commonStructures.CommonError) {

	investigationResults := []structures.InvestigationResultDbResponse{}
	err := ird.Db.Table(commonConstants.TableTestDetails).
		Select(getSelectStringsForInvestigationResults()).
		Joins("INNER JOIN test_details_metadata ON test_details_metadata.test_details_id = test_details.id").
		Joins("JOIN investigation_results ON test_details.id = investigation_results.test_details_id and investigation_results.deleted_at IS NULL").
		Joins("LEFT JOIN investigation_results_metadata ON investigation_results_metadata.investigation_result_id = investigation_results.id").
		Where("test_details.task_id = ?", taskId).
		Where("investigation_results.deleted_at IS NULL").
		Where("investigation_results_metadata.deleted_at IS NULL").
		Where("test_details.deleted_at IS NULL").
		Find(&investigationResults).Error
	if err != nil {
		return nil, commonUtils.HandleORMError(err)
	}

	return investigationResults, nil
}

func (ird *InvestigationResultDao) GetInvestigationResultsDbStructsByTestDetailIds(testDetailIds []uint) (
	[]structures.InvestigationResultDbResponse, *commonStructures.CommonError) {
	investigationResults := []structures.InvestigationResultDbResponse{}
	if err := ird.Db.Table(commonConstants.TableInvestigationResults).
		Select(getSelectStringsForInvestigationResultsV1()).
		Where("test_details_id IN (?)", testDetailIds).
		Find(&investigationResults).Error; err != nil {
		return nil, commonUtils.HandleORMError(err)
	}
	return investigationResults, nil
}

func (ird *InvestigationResultDao) GetTestDetailsByTaskId(taskId uint) (
	[]structures.TestDetailsDbResponse, *commonStructures.CommonError) {

	selectStrings := []string{
		"test_details.id as id",
		"test_details.master_test_id as master_test_id",
		"test_details.status as status",
		"test_details.lab_id as lab_id",
		"test_details.cp_enabled as cp_enabled",
		"test_details_metadata.barcodes as barcodes",
	}

	testDetails := []structures.TestDetailsDbResponse{}
	if err := ird.Db.Table(commonConstants.TableTestDetails).Select(selectStrings).
		Joins("INNER JOIN test_details_metadata on test_details.id = test_details_metadata.test_details_id").
		Where("task_id = ?", taskId).Find(&testDetails).Error; err != nil {
		return nil, commonUtils.HandleORMError(err)
	}

	return testDetails, nil
}

func (ird *InvestigationResultDao) GetInvestigationResultModelsByTaskId(taskId uint) (
	[]commonModels.InvestigationResult, *commonStructures.CommonError) {

	investigationResults := []commonModels.InvestigationResult{}
	err := ird.Db.Model(&commonModels.InvestigationResult{}).
		Joins("JOIN test_details ON test_details.id = investigation_results.test_details_id").
		Where("test_details.task_id = ?", taskId).
		Where("investigation_results.deleted_at IS NULL").
		Where("test_details.deleted_at IS NULL").
		Find(&investigationResults).Error
	if err != nil {
		return nil, commonUtils.HandleORMError(err)
	}
	return investigationResults, nil
}

func (ird *InvestigationResultDao) GetDetailsForAbnormality(taskId uint) (
	structures.BasicAbnormalityStruct, *commonStructures.CommonError) {

	basicAbnormalityStruct := structures.BasicAbnormalityStruct{}
	selectStings := []string{
		"tasks.lab_id as lab_id",
		"tasks.city_code as city_code",
		"patient_details.dob as patient_dob",
		"patient_details.expected_dob as patient_expected_dob",
		"patient_details.gender as patient_gender",
	}

	if err := ird.Db.Table(commonConstants.TableTasks).
		Joins("INNER JOIN patient_details ON patient_details.id = tasks.patient_details_id").
		Select(selectStings).
		Where("tasks.id = ?", taskId).
		Scan(&basicAbnormalityStruct).Error; err != nil {
		return basicAbnormalityStruct, commonUtils.HandleORMError(err)
	}

	return basicAbnormalityStruct, nil
}

func (ird *InvestigationResultDao) GetInvestigationResultsByTaskIdAndOmsTestId(taskId uint, omsTestId string) (
	[]commonModels.InvestigationResult, *commonStructures.CommonError) {

	investigationResults := []commonModels.InvestigationResult{}
	err := ird.Db.Model(&commonModels.InvestigationResult{}).
		Joins("INNER JOIN test_details ON test_details.id = investigation_results.test_details_id").
		Where("test_details.task_id = ?", taskId).
		Where("test_details.central_oms_test_id = ?", omsTestId).
		Find(&investigationResults).Error
	if err != nil {
		return nil, commonUtils.HandleORMError(err)
	}
	return investigationResults, nil
}

func (ird *InvestigationResultDao) GetInvestigationDataByInvestigationResultsIds(invResultIds []uint) (
	[]commonModels.InvestigationData, *commonStructures.CommonError) {

	invData := []commonModels.InvestigationData{}
	err := ird.Db.Where("investigation_result_id IN (?)", invResultIds).Find(&invData).Error
	if err != nil {
		return invData, commonUtils.HandleORMError(err)
	}
	return invData, nil
}

func (ird *InvestigationResultDao) CreateInvestigationDataWithTx(tx *gorm.DB,
	invData []commonModels.InvestigationData) ([]commonModels.InvestigationData, *commonStructures.CommonError) {

	err := tx.Create(&invData).Error
	if err != nil {
		return []commonModels.InvestigationData{}, commonUtils.HandleORMError(err)
	}
	return invData, nil
}

func (ird *InvestigationResultDao) CreateInvestigationResultsMetadataWithTx(tx *gorm.DB,
	investigationResultsMetadata []commonModels.InvestigationResultMetadata) (
	[]commonModels.InvestigationResultMetadata, *commonStructures.CommonError) {
	err := tx.Create(&investigationResultsMetadata).Error
	if err != nil {
		return []commonModels.InvestigationResultMetadata{}, commonUtils.HandleORMError(err)
	}
	return investigationResultsMetadata, nil
}

func (ird *InvestigationResultDao) UpdateInvestigationResultsMetadataWithTx(tx *gorm.DB,
	investigationResultsMetadata []commonModels.InvestigationResultMetadata) (
	[]commonModels.InvestigationResultMetadata, *commonStructures.CommonError) {
	err := tx.Save(&investigationResultsMetadata).Error
	if err != nil {
		return []commonModels.InvestigationResultMetadata{}, commonUtils.HandleORMError(err)
	}
	return investigationResultsMetadata, nil
}

func (ird *InvestigationResultDao) DeleteInvestigationResultsMetadataByIdsWithTx(tx *gorm.DB,
	investigationIds []uint) *commonStructures.CommonError {
	currentTime := commonUtils.GetCurrentTime()
	updates := map[string]interface{}{
		"deleted_at": currentTime,
		"updated_at": currentTime,
		"deleted_by": commonConstants.CitadelSystemId,
		"updated_by": commonConstants.CitadelSystemId,
	}

	err := tx.Model(&commonModels.InvestigationResultMetadata{}).
		Where("id IN (?)", investigationIds).
		Updates(updates).Error
	if err != nil {
		return commonUtils.HandleORMError(err)
	}

	return nil
}
