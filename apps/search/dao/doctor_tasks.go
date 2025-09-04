package dao

import (
	"gorm.io/gorm"

	"github.com/Orange-Health/citadel/apps/search/constants"
	"github.com/Orange-Health/citadel/apps/search/structures"
	commonConstants "github.com/Orange-Health/citadel/common/constants"
	commonStructures "github.com/Orange-Health/citadel/common/structures"
	commonUtils "github.com/Orange-Health/citadel/common/utils"
	commonModels "github.com/Orange-Health/citadel/models"
)

func getCommonSelectStringsForTaskDetails(taskType string) []string {
	selectStrings := []string{
		"DISTINCT(tasks.id) as task_id",
		"tasks.order_id as order_id",
		"tasks.oms_order_id as oms_order_id",
		"tasks.city_code as city_code",
		"tasks.doctor_tat as doctor_tat",
		"tasks.status as status",
		"tasks.lab_id as lab_id",
		"patient_details.name as patient_name",
		"task_metadata.contains_morphle as contains_morphle",
		"task_metadata.contains_package as contains_package",
	}

	switch taskType {
	case constants.TASK_TYPE_IN_PROGRESS:
		selectStrings = append(selectStrings, "task_pathologist_mapping.pathologist_id as picked_by")
	case constants.TASK_TYPE_CO_AUTHORIZE:
		selectStrings = append(selectStrings, "co_authorized_pathologists.co_authorized_by as co_authorized_by")
	}

	return selectStrings
}

func getSelectStringsForAmendmentTaskDetails() []string {
	selectStrings := []string{
		"tasks.id as task_id",
		"tasks.order_id as order_id",
		"tasks.city_code as city_code",
		"patient_details.name as patient_name",
		"tasks.lab_id as lab_id",
	}

	return selectStrings
}

func (searchDao *SearchDao) applyWhereClauseBasedOnTaskType(query *gorm.DB,
	taskStatuses []string, isCritical bool, taskType string, userId uint) *gorm.DB {

	testDetailStatuses := []string{
		commonConstants.TEST_STATUS_RESULT_SAVED,
		commonConstants.TEST_STATUS_RERUN_RESULT_SAVED,
		commonConstants.TEST_STATUS_WITHHELD,
		commonConstants.TEST_STATUS_CO_AUTHORIZE,
	}

	query = query.Where("tasks.status IN (?)", taskStatuses).
		Where("test_details.status IN (?)", testDetailStatuses).
		Where("test_details.cp_enabled = ?", true).
		Where("test_details.deleted_at IS NULL")

	switch taskType {
	case constants.TASK_TYPE_CRITICAL, constants.TASK_TYPE_NORMAL:
		query = query.Where("task_metadata.is_critical = ?", isCritical).
			Where("task_pathologist_mapping.is_active = ? OR task_pathologist_mapping.task_id IS NULL", false)
	case constants.TASK_TYPE_WITHHELD:
		query = query.Where("task_pathologist_mapping.is_active = ? OR task_pathologist_mapping.task_id IS NULL", false)
	case constants.TASK_TYPE_CO_AUTHORIZE:
		query = query.Joins("INNER JOIN co_authorized_pathologists ON tasks.id = co_authorized_pathologists.task_id").
			Where("co_authorized_pathologists.co_authorized_to = ?", userId).
			Where("co_authorized_pathologists.deleted_at IS NULL").
			Where("task_pathologist_mapping.is_active = ? OR task_pathologist_mapping.task_id IS NULL", false)
	case constants.TASK_TYPE_IN_PROGRESS:
		query = query.Where("task_pathologist_mapping.is_active = ?", true)
	}

	return query
}

func (searchDao *SearchDao) addClausesInQuery(query *gorm.DB,
	taskListDbRequest structures.TaskListDbRequest) *gorm.DB {
	if len(taskListDbRequest.VisitClause) > 0 {
		query = query.Joins("INNER JOIN task_visit_mapping ON tasks.id = task_visit_mapping.task_id")
		query = query.Where(taskListDbRequest.VisitClause)
	}

	if len(taskListDbRequest.PatientDetailsClause) > 0 {
		query = query.Where(taskListDbRequest.PatientDetailsClause)
	}

	if taskListDbRequest.PatientName != "" {
		query = query.Where("lower(patient_details.name) LIKE ?", taskListDbRequest.PatientName+"%")
	}

	if taskListDbRequest.PartnerName != "" {
		query = query.Where("lower(task_metadata.partner_name) LIKE ?", taskListDbRequest.PartnerName+"%")
	}

	if taskListDbRequest.DoctorName != "" {
		query = query.Where("lower(task_metadata.doctor_name) LIKE ?", taskListDbRequest.DoctorName+"%")
	}

	if len(taskListDbRequest.TasksClause) > 0 {
		query = query.Where(taskListDbRequest.TasksClause)
	}

	if len(taskListDbRequest.TaskMetadataClause) > 0 {
		query = query.Where(taskListDbRequest.TaskMetadataClause)
	}

	if len(taskListDbRequest.TestDetailsClause) > 0 {
		query = query.Where(taskListDbRequest.TestDetailsClause)
	}

	return query
}

func (searchDao *SearchDao) getTaskListCommonQuery(taskListDbRequest structures.TaskListDbRequest,
	taskStatuses []string, isCritical bool, taskType string) *gorm.DB {

	query := searchDao.Db.Model(&commonModels.Task{}).
		Joins("INNER JOIN task_metadata ON tasks.id = task_metadata.task_id").
		Joins("INNER JOIN patient_details ON tasks.patient_details_id = patient_details.id").
		Joins("INNER JOIN test_details ON tasks.id = test_details.task_id").
		Joins("LEFT JOIN task_pathologist_mapping ON tasks.id = task_pathologist_mapping.task_id").
		Select(getCommonSelectStringsForTaskDetails(taskType))

	query = searchDao.addClausesInQuery(query, taskListDbRequest)

	query = searchDao.applyWhereClauseBasedOnTaskType(query, taskStatuses, isCritical, taskType, taskListDbRequest.UserId)

	query = query.Limit(int(taskListDbRequest.Limit + 1)).
		Offset(int(taskListDbRequest.Offset))

	query = query.Order("tasks.doctor_tat")

	return query
}

func (searchDao *SearchDao) getAmendmentTaskListQuery(taskListDbRequest structures.TaskListDbRequest) *gorm.DB {
	query := searchDao.Db.Model(&commonModels.Task{}).
		Joins("INNER JOIN patient_details ON tasks.patient_details_id = patient_details.id")

	if len(taskListDbRequest.TaskMetadataClause) > 0 ||
		taskListDbRequest.PartnerName != "" ||
		taskListDbRequest.DoctorName != "" {
		query = query.Joins("INNER JOIN task_metadata ON tasks.id = task_metadata.task_id")
	}
	query = query.Select(getSelectStringsForAmendmentTaskDetails())

	query = searchDao.addClausesInQuery(query, taskListDbRequest)

	query = query.Where("tasks.status = ?", commonConstants.TASK_STATUS_COMPLETED).
		Limit(int(taskListDbRequest.Limit)).
		Offset(int(taskListDbRequest.Offset))

	query = query.Order("tasks.completed_at DESC")

	return query
}

func (searchDao *SearchDao) GetCriticalTaskDetails(taskListDbRequest structures.TaskListDbRequest) (
	[]structures.TaskDetailsDbStruct, *commonStructures.CommonError) {

	taskDetails := []structures.TaskDetailsDbStruct{}
	taskStatuses := []string{commonConstants.TASK_STATUS_PENDING}
	err := searchDao.getTaskListCommonQuery(taskListDbRequest, taskStatuses, true, constants.TASK_TYPE_CRITICAL).
		Scan(&taskDetails).Error

	if err != nil {
		return taskDetails, commonUtils.HandleORMError(err)
	}

	return taskDetails, nil
}

func (searchDao *SearchDao) GetWithheldTaskDetails(taskListDbRequest structures.TaskListDbRequest) (
	[]structures.TaskDetailsDbStruct, *commonStructures.CommonError) {

	taskDetails := []structures.TaskDetailsDbStruct{}
	taskStatuses := []string{commonConstants.TASK_STATUS_WITHHELD_APPROVAL}
	err := searchDao.getTaskListCommonQuery(taskListDbRequest, taskStatuses, false, constants.TASK_TYPE_WITHHELD).
		Scan(&taskDetails).Error

	if err != nil {
		return taskDetails, commonUtils.HandleORMError(err)
	}

	return taskDetails, nil
}

func (searchDao *SearchDao) GetCoAuthorizedTaskDetails(taskListDbRequest structures.TaskListDbRequest) (
	[]structures.TaskDetailsDbStruct, *commonStructures.CommonError) {

	taskDetails := []structures.TaskDetailsDbStruct{}
	taskStatuses := []string{commonConstants.TASK_STATUS_CO_AUTHORIZE}
	err := searchDao.getTaskListCommonQuery(taskListDbRequest, taskStatuses, false, constants.TASK_TYPE_CO_AUTHORIZE).
		Scan(&taskDetails).Error

	if err != nil {
		return taskDetails, commonUtils.HandleORMError(err)
	}

	return taskDetails, nil
}

func (searchDao *SearchDao) GetNormalTaskDetails(taskListDbRequest structures.TaskListDbRequest) (
	[]structures.TaskDetailsDbStruct, *commonStructures.CommonError) {

	taskDetails := []structures.TaskDetailsDbStruct{}
	taskStatuses := []string{commonConstants.TASK_STATUS_PENDING}
	err := searchDao.getTaskListCommonQuery(taskListDbRequest, taskStatuses, false, constants.TASK_TYPE_NORMAL).
		Scan(&taskDetails).Error

	if err != nil {
		return taskDetails, commonUtils.HandleORMError(err)
	}

	return taskDetails, nil
}

func (searchDao *SearchDao) GetInProgressTaskDetails(taskListDbRequest structures.TaskListDbRequest) (
	[]structures.TaskDetailsDbStruct, *commonStructures.CommonError) {

	taskDetails := []structures.TaskDetailsDbStruct{}
	taskStatuses := []string{commonConstants.TASK_STATUS_IN_PROGRESS}
	err := searchDao.getTaskListCommonQuery(taskListDbRequest, taskStatuses, false, constants.TASK_TYPE_IN_PROGRESS).
		Scan(&taskDetails).Error

	if err != nil {
		return taskDetails, commonUtils.HandleORMError(err)
	}

	return taskDetails, nil
}

func (searchDao *SearchDao) GetAmendmentTaskDetails(taskListDbRequest structures.TaskListDbRequest) (
	[]structures.TaskDetailsDbStruct, *commonStructures.CommonError) {

	taskDetails := []structures.TaskDetailsDbStruct{}
	err := searchDao.getAmendmentTaskListQuery(taskListDbRequest).
		Scan(&taskDetails).Error

	if err != nil {
		return taskDetails, commonUtils.HandleORMError(err)
	}

	return taskDetails, nil
}
