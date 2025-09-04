package dao

import (
	"fmt"
	"strings"

	"gorm.io/gorm"

	"github.com/Orange-Health/citadel/apps/test_detail/mapper"
	"github.com/Orange-Health/citadel/apps/test_detail/structures"
	commonConstants "github.com/Orange-Health/citadel/common/constants"
	commonStructures "github.com/Orange-Health/citadel/common/structures"
	commonUtils "github.com/Orange-Health/citadel/common/utils"
	commonModels "github.com/Orange-Health/citadel/models"
)

type DataLayer interface {
	// Test Details
	GetTestDetailById(testDetailID uint) (
		commonModels.TestDetail, *commonStructures.CommonError)
	GetTestDetailByOmsTestId(omsTestId string) (
		commonModels.TestDetail, *commonStructures.CommonError)
	GetTestDetailByOmsTestIds(omsTestIds []string) (
		[]commonModels.TestDetail, *commonStructures.CommonError)
	GetTestDetailByIds(testDetailIDs []uint) (
		[]structures.TestDetail, *commonStructures.CommonError)
	GetTestDetailByIdsWithTx(tx *gorm.DB, testDetailIDs []uint) (
		[]structures.TestDetail, *commonStructures.CommonError)
	GetTestDetailByIdWithTx(tx *gorm.DB, testDetailID uint) (
		commonModels.TestDetail, *commonStructures.CommonError)
	GetTestDetailsByTaskId(taskID uint) (
		[]commonModels.TestDetail, *commonStructures.CommonError)
	GetTestDetailsByOmsOrderId(omsOrderId string) (
		[]commonModels.TestDetail, *commonStructures.CommonError)
	GetTestDetailsByOmsOrderIds(omsOrderIds []string) (
		[]commonModels.TestDetail, *commonStructures.CommonError)
	GetTestDetailsByOmsOrderIdWithTx(tx *gorm.DB, omsOrderId string) (
		[]commonModels.TestDetail, *commonStructures.CommonError)
	GetTestDetailsIdsByOmsOrderIdWithTx(tx *gorm.DB, omsOrderId string) (
		[]uint, *commonStructures.CommonError)
	GetTestDetailsByTestIds(testIds []uint) (
		[]commonModels.TestDetail, *commonStructures.CommonError)
	GetTestBasicDetailsForSearchScreenByTaskIds(taskIds []uint) (
		[]structures.TestBasicDetails, *commonStructures.CommonError)
	GetOmsTestIdAndStatusByOmsOrderId(omsOrderId string) (
		[]structures.TestBasicDetails, *commonStructures.CommonError)
	GetAllBasicTestDetailsByOmsOrderId(omsOrderId string) (
		[]structures.OmsTestBasicDetailsDbStruct, *commonStructures.CommonError)
	GetActiveTestDetailsByTaskId(taskID uint) (
		[]commonModels.TestDetail, *commonStructures.CommonError)
	GetTestDetailsByOmsOrderIdWithSampleStatus(omsOrderId, sampleStatus string) (
		[]commonModels.TestDetail, *commonStructures.CommonError)
	ContainsPackageTests(omsOrderId string) bool
	UpdateTestDetails(testDetails []commonModels.TestDetail) (
		[]commonModels.TestDetail, *commonStructures.CommonError)
	CreateTestDetailsWithTx(tx *gorm.DB, testDetails []commonModels.TestDetail) (
		[]commonModels.TestDetail, *commonStructures.CommonError)
	UpdateTestDetailsById(testDetails commonModels.TestDetail) *commonStructures.CommonError
	UpdateTestDetailsWithTx(tx *gorm.DB, testDetails []commonModels.TestDetail) (
		[]commonModels.TestDetail, *commonStructures.CommonError)
	UpdateTaskIdInTestDetailsWithOmsTestIdWithTx(tx *gorm.DB, omsTestIds []string,
		taskId uint) *commonStructures.CommonError
	UpdateTaskIdInTestDetailsWithOmsTestIds(centralOmsTestIds []string, taskId uint) *commonStructures.CommonError
	UpdateTestStatusesByOmsTestIdsWithTx(tx *gorm.DB, testIdStatusMap map[string]string,
		userId uint) *commonStructures.CommonError
	UpdateDuplicateTestDetailsByTaskIdWithTx(tx *gorm.DB, taskId uint, masterTestIds []uint) *commonStructures.CommonError
	DeleteTestDetailsByOmsTestIdsWithTx(tx *gorm.DB, omsTestIds []string) *commonStructures.CommonError
	DeleteTestDetailsByTaskIdAndOmsTestIdWithTx(tx *gorm.DB, taskID uint, omsTestId string) *commonStructures.CommonError
	DeleteTestDetailsByOmsOrderIdWithTx(tx *gorm.DB, omsOrderId string) *commonStructures.CommonError
	UpdateReportStatusByOmsTestIds(omsTestIds []string, oldStatus, newStatus string) *commonStructures.CommonError

	// Test Details Metadata
	GetTestDetailsMetadataByTestDetailIds(testDetailsIds []uint) (
		[]commonModels.TestDetailsMetadata, *commonStructures.CommonError)
	CreateTestDetailsMetadataWithTransaction(tx *gorm.DB, testDetailsMetadata []commonModels.TestDetailsMetadata) (
		[]commonModels.TestDetailsMetadata, *commonStructures.CommonError)
	UpdateTestDetailsMetadataWithTransaction(tx *gorm.DB, testDetailsMetadata []commonModels.TestDetailsMetadata) (
		[]commonModels.TestDetailsMetadata, *commonStructures.CommonError)
	UpdatePickedAtTimeBasedOnActiveTests(testDetailsIds []uint) *commonStructures.CommonError
	DeleteTestDetailsMetadataByTestDetailsIdsWithTx(tx *gorm.DB, testDetailsIds []uint) *commonStructures.CommonError
}

func (testDetailDao *TestDetailDao) GetTestDetailById(testDetailId uint) (
	commonModels.TestDetail, *commonStructures.CommonError) {

	td := commonModels.TestDetail{}
	err := testDetailDao.Db.Where("id = ?", testDetailId).First(&td).Error
	if err != nil {
		return commonModels.TestDetail{}, commonUtils.HandleORMError(err)
	}
	return td, nil
}

func (testDetailDao *TestDetailDao) GetTestDetailByOmsTestId(omsTestId string) (
	commonModels.TestDetail, *commonStructures.CommonError) {

	td := commonModels.TestDetail{}
	err := testDetailDao.Db.Where("central_oms_test_id = ?", omsTestId).First(&td).Error
	if err != nil {
		return commonModels.TestDetail{}, commonUtils.HandleORMError(err)
	}
	return td, nil
}

func (testDetailDao *TestDetailDao) GetTestDetailByOmsTestIds(omsTestIds []string) ([]commonModels.TestDetail,
	*commonStructures.CommonError) {

	testDetails := []commonModels.TestDetail{}
	err := testDetailDao.Db.Where("central_oms_test_id IN (?)", omsTestIds).Find(&testDetails).Error
	if err != nil {
		return []commonModels.TestDetail{}, commonUtils.HandleORMError(err)
	}
	return testDetails, nil
}

func (testDetailDao *TestDetailDao) GetTestDetailByIds(testDetailIDs []uint) (
	[]structures.TestDetail, *commonStructures.CommonError) {

	tds := []commonModels.TestDetail{}
	err := testDetailDao.Db.Where("oms_test_id in (?)", testDetailIDs).First(&tds).Error
	if err != nil {
		return []structures.TestDetail{}, commonUtils.HandleORMError(err)
	}
	return mapper.MapTestDetails(tds), nil
}

func (testDetailDao *TestDetailDao) GetTestDetailByIdsWithTx(tx *gorm.DB, testDetailIDs []uint) (
	[]structures.TestDetail, *commonStructures.CommonError) {

	tds := []commonModels.TestDetail{}
	err := tx.Where("id in (?)", testDetailIDs).First(&tds).Error
	if err != nil {
		return []structures.TestDetail{}, commonUtils.HandleORMError(err)
	}
	return mapper.MapTestDetails(tds), nil
}

func (testDetailDao *TestDetailDao) GetTestDetailByIdWithTx(tx *gorm.DB, testDetailID uint) (
	commonModels.TestDetail, *commonStructures.CommonError) {

	td := commonModels.TestDetail{}
	err := tx.Where("id = ?", testDetailID).First(&td).Error
	if err != nil {
		return commonModels.TestDetail{}, commonUtils.HandleORMError(err)
	}
	return td, nil
}

func (testDetailDao *TestDetailDao) GetTestDetailsByTaskId(taskId uint) (
	[]commonModels.TestDetail, *commonStructures.CommonError) {

	task := commonModels.Task{}
	err := testDetailDao.Db.Where("id = ?", taskId).First(&task).Error
	if err != nil {
		return []commonModels.TestDetail{}, commonUtils.HandleORMError(err)
	}

	testDetails := []commonModels.TestDetail{}
	err = testDetailDao.Db.Find(&testDetails, "oms_order_id = ?", task.OmsOrderId).Error
	if err != nil {
		return []commonModels.TestDetail{}, commonUtils.HandleORMError(err)
	}
	return testDetails, nil
}

func (testDetailDao *TestDetailDao) GetTestDetailsByOmsOrderId(omsOrderId string) (
	[]commonModels.TestDetail, *commonStructures.CommonError) {

	testDetails := []commonModels.TestDetail{}
	if err := testDetailDao.Db.Find(&testDetails, "oms_order_id = ?", omsOrderId).Error; err != nil {
		return []commonModels.TestDetail{}, commonUtils.HandleORMError(err)
	}

	return testDetails, nil
}

func (testDetailDao *TestDetailDao) GetTestDetailsByOmsOrderIds(omsOrderIds []string) (
	[]commonModels.TestDetail, *commonStructures.CommonError) {
	testDetails := []commonModels.TestDetail{}
	if err := testDetailDao.Db.Find(&testDetails, "oms_order_id IN (?)", omsOrderIds).Error; err != nil {
		return []commonModels.TestDetail{}, commonUtils.HandleORMError(err)
	}
	return testDetails, nil
}

func (testDetailDao *TestDetailDao) GetTestDetailsByOmsOrderIdWithTx(tx *gorm.DB, omsOrderId string) (
	[]commonModels.TestDetail, *commonStructures.CommonError) {

	testDetails := []commonModels.TestDetail{}
	if err := tx.Find(&testDetails, "oms_order_id = ?", omsOrderId).Error; err != nil {
		return []commonModels.TestDetail{}, commonUtils.HandleORMError(err)
	}

	return testDetails, nil
}

func (testDetailDao *TestDetailDao) GetTestDetailsIdsByOmsOrderIdWithTx(tx *gorm.DB, omsOrderId string) (
	[]uint, *commonStructures.CommonError) {
	testDetailsIds := []uint{}
	if err := tx.Model(&commonModels.TestDetail{}).
		Where("oms_order_id = ?", omsOrderId).
		Pluck("id", &testDetailsIds).Error; err != nil {
		return []uint{}, commonUtils.HandleORMError(err)
	}
	return testDetailsIds, nil
}

func (testDetailDao *TestDetailDao) GetTestDetailsByTestIds(testIds []uint) (
	[]commonModels.TestDetail, *commonStructures.CommonError) {

	testDetails := []commonModels.TestDetail{}
	err := testDetailDao.Db.Find(&testDetails, "id IN (?)", testIds).Error
	if err != nil {
		return []commonModels.TestDetail{}, commonUtils.HandleORMError(err)
	}
	return testDetails, nil
}

func (testDetailDao *TestDetailDao) GetTestBasicDetailsForSearchScreenByTaskIds(taskIds []uint) (
	[]structures.TestBasicDetails, *commonStructures.CommonError) {

	testBasicDetails := []structures.TestBasicDetails{}
	if err := testDetailDao.Db.Table(commonConstants.TableTestDetails).
		Joins("JOIN tasks ON tasks.id = test_details.task_id").
		Select("test_details.task_id, test_details.test_name, test_details.status, test_details.lab_id").
		Where("test_details.task_id IN (?)", taskIds).
		Where("test_details.cp_enabled = ?", true).
		Where("test_details.deleted_at IS NULL").
		Where("tasks.deleted_at IS NULL").
		Scan(&testBasicDetails).Error; err != nil {
		return []structures.TestBasicDetails{}, commonUtils.HandleORMError(err)
	}
	return testBasicDetails, nil
}

func (testDetailDao *TestDetailDao) GetOmsTestIdAndStatusByOmsOrderId(omsOrderId string) (
	[]structures.TestBasicDetails, *commonStructures.CommonError) {
	testBasicDetails := []structures.TestBasicDetails{}
	if err := testDetailDao.Db.Table(commonConstants.TableTestDetails).
		Select("test_details.central_oms_test_id as oms_test_id, test_details.status").
		Where("test_details.oms_order_id = ?", omsOrderId).
		Where("test_details.deleted_at IS NULL").
		Scan(&testBasicDetails).Error; err != nil {
		return []structures.TestBasicDetails{}, commonUtils.HandleORMError(err)
	}
	return testBasicDetails, nil
}

func (testDetailDao *TestDetailDao) GetAllBasicTestDetailsByOmsOrderId(omsOrderId string) (
	[]structures.OmsTestBasicDetailsDbStruct, *commonStructures.CommonError) {
	omsTestBasicDetails := []structures.OmsTestBasicDetailsDbStruct{}
	selectStrings := []string{
		"test_details.central_oms_test_id as test_id",
		"test_details.test_name as test_name",
		"test_details.status as test_status",
		"test_details.lab_id as lab_id",
		"test_details.master_test_id as master_test_id",
	}
	if err := testDetailDao.Db.Table(commonConstants.TableTestDetails).
		Select(selectStrings).
		Where("test_details.oms_order_id = ?", omsOrderId).
		Where("test_details.deleted_at IS NULL").
		Scan(&omsTestBasicDetails).Error; err != nil {
		return nil, commonUtils.HandleORMError(err)
	}

	return omsTestBasicDetails, nil
}

func (testDetailDao *TestDetailDao) GetTestDetailsByOmsOrderIdWithSampleStatus(omsOrderId, sampleStatus string) (
	[]commonModels.TestDetail, *commonStructures.CommonError) {
	testDetails := []commonModels.TestDetail{}
	err := testDetailDao.Db.Model(&commonModels.TestDetail{}).
		Joins("INNER JOIN test_sample_mapping ON test_sample_mapping.oms_test_id = test_details.central_oms_test_id").
		Joins("INNER JOIN samples ON samples.oms_order_id = test_sample_mapping.oms_order_id AND samples.sample_number = test_sample_mapping.sample_number").
		Where("test_details.oms_order_id = ? AND samples.status = ?", omsOrderId, sampleStatus).
		Where("test_details.deleted_at IS NULL").
		Where("samples.deleted_at IS NULL").
		Where("test_sample_mapping.deleted_at IS NULL").
		Find(&testDetails).Error
	if err != nil {
		return []commonModels.TestDetail{}, commonUtils.HandleORMError(err)
	}
	return testDetails, nil
}

func (testDetailDao *TestDetailDao) GetActiveTestDetailsByTaskId(taskID uint) (
	[]commonModels.TestDetail, *commonStructures.CommonError) {

	testDetailStatuses := []string{
		commonConstants.TEST_STATUS_RESULT_SAVED,
		commonConstants.TEST_STATUS_RERUN_RESULT_SAVED,
	}
	testDetails := []commonModels.TestDetail{}
	err := testDetailDao.Db.
		Find(&testDetails, "task_id = ? AND status IN (?) AND deleted_at IS NULL", taskID, testDetailStatuses).Error
	if err != nil {
		return []commonModels.TestDetail{}, commonUtils.HandleORMError(err)
	}
	return testDetails, nil
}

func (testDetailDao *TestDetailDao) ContainsPackageTests(omsOrderId string) bool {
	testDetail := commonModels.TestDetail{}
	testDetailDao.Db.Model(&commonModels.TestDetail{}).
		Where("oms_order_id = ? AND master_package_id != 0", omsOrderId).
		Find(&testDetail)
	return testDetail.Id != 0
}

func (testDetailDao *TestDetailDao) UpdateTestDetails(testDetails []commonModels.TestDetail) (
	[]commonModels.TestDetail, *commonStructures.CommonError) {

	err := testDetailDao.Db.Save(&testDetails).Error
	if err != nil {
		return []commonModels.TestDetail{}, commonUtils.HandleORMError(err)
	}
	return testDetails, nil
}

func (testDetailDao *TestDetailDao) CreateTestDetailsWithTx(tx *gorm.DB,
	testDetails []commonModels.TestDetail) (
	[]commonModels.TestDetail, *commonStructures.CommonError) {

	err := tx.Create(&testDetails).Error
	if err != nil {
		return []commonModels.TestDetail{}, commonUtils.HandleORMError(err)
	}
	return testDetails, nil
}

func (testDetailDao *TestDetailDao) UpdateTestDetailsWithTx(tx *gorm.DB,
	testDetails []commonModels.TestDetail) (
	[]commonModels.TestDetail, *commonStructures.CommonError) {

	err := tx.Save(&testDetails).Error
	if err != nil {
		return []commonModels.TestDetail{}, commonUtils.HandleORMError(err)
	}
	return testDetails, nil
}

func (testDetailDao *TestDetailDao) UpdateTaskIdInTestDetailsWithOmsTestIdWithTx(tx *gorm.DB, omsTestIds []string,
	taskId uint) *commonStructures.CommonError {
	updatesMap := map[string]interface{}{
		"task_id":    taskId,
		"status":     commonConstants.TEST_STATUS_RESULT_PENDING,
		"updated_at": commonUtils.GetCurrentTime(),
		"updated_by": commonConstants.CitadelSystemId,
	}

	if err := tx.Model(&commonModels.TestDetail{}).
		Where("central_oms_test_id IN (?)", omsTestIds).
		Updates(updatesMap).Error; err != nil {
		return commonUtils.HandleORMError(err)
	}

	return nil
}

func (testDetailDao *TestDetailDao) UpdateTaskIdInTestDetailsWithOmsTestIds(centralOmsTestIds []string,
	taskId uint) *commonStructures.CommonError {
	if len(centralOmsTestIds) == 0 {
		return nil
	}

	updatesMap := map[string]interface{}{
		"task_id":    taskId,
		"status":     commonConstants.TEST_STATUS_RESULT_PENDING,
		"updated_at": commonUtils.GetCurrentTime(),
		"updated_by": commonConstants.CitadelSystemId,
	}

	if err := testDetailDao.Db.Model(&commonModels.TestDetail{}).
		Where("central_oms_test_id IN (?)", centralOmsTestIds).
		Updates(updatesMap).Error; err != nil {
		return commonUtils.HandleORMError(err)
	}

	return nil
}

func (testDetailDao *TestDetailDao) UpdateTestStatusesByOmsTestIdsWithTx(tx *gorm.DB, testIdStatusMap map[string]string,
	userId uint) *commonStructures.CommonError {
	cases, testIds := []string{}, []string{}

	for id, status := range testIdStatusMap {
		cases = append(cases, fmt.Sprintf("WHEN '%s' THEN '%s'", id, status))
		testIds = append(testIds, fmt.Sprintf("'%s'", id))
	}

	cases = append(cases, "ELSE status")

	query := fmt.Sprintf(
		"UPDATE test_details SET status = CASE central_oms_test_id %s END, updated_by = %d WHERE central_oms_test_id IN (%s);",
		strings.Join(cases, " "), userId, strings.Join(testIds, ", "),
	)

	if err := tx.Exec(query).Error; err != nil {
		return commonUtils.HandleORMError(err)
	}

	return nil
}

func (testDetailDao *TestDetailDao) UpdateDuplicateTestDetailsByTaskIdWithTx(tx *gorm.DB, taskId uint,
	masterTestIds []uint) *commonStructures.CommonError {
	updates := map[string]interface{}{
		"updated_at":   commonUtils.GetCurrentTime(),
		"updated_by":   commonConstants.CitadelSystemId,
		"is_duplicate": true,
	}

	err := tx.Model(&commonModels.TestDetail{}).
		Where("task_id = ?", taskId).
		Where("master_test_id IN (?)", masterTestIds).
		Where("is_duplicate = false").
		Updates(updates).Error
	if err != nil {
		return commonUtils.HandleORMError(err)
	}
	return nil
}

func (testDetail *TestDetailDao) UpdateTestDetailsById(testDetails commonModels.TestDetail) *commonStructures.CommonError {
	err := testDetail.Db.Save(&testDetails).Error
	if err != nil {
		return commonUtils.HandleORMError(err)
	}
	return nil
}

func (TestDetailDao *TestDetailDao) DeleteTestDetailsByOmsTestIdsWithTx(tx *gorm.DB,
	omsTestIds []string) *commonStructures.CommonError {

	currentTime := commonUtils.GetCurrentTime()
	updates := map[string]interface{}{
		"deleted_at": currentTime,
		"updated_at": currentTime,
		"deleted_by": commonConstants.CitadelSystemId,
		"updated_by": commonConstants.CitadelSystemId,
	}

	err := tx.Model(&commonModels.TestDetail{}).
		Where("central_oms_test_id IN (?)", omsTestIds).
		Updates(updates).Error
	if err != nil {
		return commonUtils.HandleORMError(err)
	}
	return nil
}

func (testDetailDao *TestDetailDao) DeleteTestDetailsByTaskIdAndOmsTestIdWithTx(tx *gorm.DB,
	taskID uint, omsTestId string) *commonStructures.CommonError {

	currentTime := commonUtils.GetCurrentTime()
	updates := map[string]interface{}{
		"deleted_at": currentTime,
		"updated_at": currentTime,
		"deleted_by": commonConstants.CitadelSystemId,
		"updated_by": commonConstants.CitadelSystemId,
	}

	err := tx.Model(&commonModels.TestDetail{}).
		Where("task_id = ? AND central_oms_test_id = ?", taskID, omsTestId).
		Updates(updates).Error
	if err != nil {
		return commonUtils.HandleORMError(err)
	}
	return nil
}

func (testDetailDao *TestDetailDao) DeleteTestDetailsByOmsOrderIdWithTx(tx *gorm.DB,
	omsOrderId string) *commonStructures.CommonError {

	currentTime := commonUtils.GetCurrentTime()
	updates := map[string]interface{}{
		"deleted_at": currentTime,
		"updated_at": currentTime,
		"deleted_by": commonConstants.CitadelSystemId,
		"updated_by": commonConstants.CitadelSystemId,
	}

	if err := tx.Model(&commonModels.TestDetail{}).
		Where("oms_order_id = ?", omsOrderId).
		Updates(updates).Error; err != nil {
		return commonUtils.HandleORMError(err)
	}

	return nil
}

func (testDetailDao *TestDetailDao) UpdateReportStatusByOmsTestIds(omsTestIds []string,
	oldStatus, newStatus string) *commonStructures.CommonError {

	updates := map[string]interface{}{
		"report_status": newStatus,
		"updated_at":    commonUtils.GetCurrentTime(),
		"updated_by":    commonConstants.CitadelSystemId,
	}

	if err := testDetailDao.Db.Model(&commonModels.TestDetail{}).
		Where("central_oms_test_id IN (?) AND report_status = ?", omsTestIds, oldStatus).
		Updates(updates).Error; err != nil {
		return commonUtils.HandleORMError(err)
	}

	return nil
}

func (testDetailDao *TestDetailDao) GetTestDetailsMetadataByTestDetailIds(
	testDetailsIds []uint) (
	[]commonModels.TestDetailsMetadata, *commonStructures.CommonError) {

	testDetailsMetadata := []commonModels.TestDetailsMetadata{}
	err := testDetailDao.Db.Find(&testDetailsMetadata, "test_details_id IN (?)", testDetailsIds).Error
	if err != nil {
		return []commonModels.TestDetailsMetadata{}, commonUtils.HandleORMError(err)
	}
	return testDetailsMetadata, nil
}

func (testDetailDao *TestDetailDao) CreateTestDetailsMetadataWithTransaction(tx *gorm.DB,
	testDetailsMetadata []commonModels.TestDetailsMetadata) (
	[]commonModels.TestDetailsMetadata, *commonStructures.CommonError) {

	err := tx.Create(&testDetailsMetadata).Error
	if err != nil {
		return []commonModels.TestDetailsMetadata{}, commonUtils.HandleORMError(err)
	}
	return testDetailsMetadata, nil
}

func (testDetailDao *TestDetailDao) UpdateTestDetailsMetadataWithTransaction(tx *gorm.DB,
	testDetailsMetadata []commonModels.TestDetailsMetadata) (
	[]commonModels.TestDetailsMetadata, *commonStructures.CommonError) {

	err := tx.Save(&testDetailsMetadata).Error
	if err != nil {
		return []commonModels.TestDetailsMetadata{}, commonUtils.HandleORMError(err)
	}
	return testDetailsMetadata, nil
}

func (testDetailDao *TestDetailDao) UpdatePickedAtTimeBasedOnActiveTests(testDetailsIds []uint) *commonStructures.CommonError {

	currentTime := commonUtils.GetCurrentTime()
	updates := map[string]interface{}{
		"picked_at":  currentTime,
		"updated_at": currentTime,
		"updated_by": commonConstants.CitadelSystemId,
	}

	err := testDetailDao.Db.Model(&commonModels.TestDetailsMetadata{}).
		Where("test_details_id IN (?)", testDetailsIds).
		Updates(updates).Error
	if err != nil {
		return commonUtils.HandleORMError(err)
	}
	return nil
}

func (testDetailDao *TestDetailDao) DeleteTestDetailsMetadataByTestDetailsIdsWithTx(tx *gorm.DB,
	testDetailsIds []uint) *commonStructures.CommonError {
	currentTime := commonUtils.GetCurrentTime()
	updates := map[string]interface{}{
		"deleted_at": currentTime,
		"updated_at": currentTime,
		"deleted_by": commonConstants.CitadelSystemId,
		"updated_by": commonConstants.CitadelSystemId,
	}

	err := tx.Model(&commonModels.TestDetailsMetadata{}).
		Where("test_details_id IN (?)", testDetailsIds).
		Updates(updates).Error
	if err != nil {
		return commonUtils.HandleORMError(err)
	}

	return nil
}
