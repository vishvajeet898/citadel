package service

import (
	"context"
	"net/http"
	"strings"
	"sync"

	"github.com/Orange-Health/citadel/apps/search/constants"
	"github.com/Orange-Health/citadel/apps/search/structures"
	commonConstants "github.com/Orange-Health/citadel/common/constants"
	commonStructures "github.com/Orange-Health/citadel/common/structures"
	commonUtils "github.com/Orange-Health/citadel/common/utils"
)

func validateAndReturnSearchParameters(taskListBasicRequest structures.TaskListBasicRequest) (
	structures.TasksSearch, *commonStructures.CommonError) {

	tasksSearch := structures.TasksSearch{
		PatientName:   taskListBasicRequest.PatientName,
		ContactNumber: taskListBasicRequest.ContactNumber,
		PatientId:     taskListBasicRequest.PatientId,
		PartnerName:   taskListBasicRequest.PartnerName,
		DoctorName:    taskListBasicRequest.DoctorName,
		VisitId:       taskListBasicRequest.VisitId,
	}

	if taskListBasicRequest.OrderId != "" {
		orderId := commonUtils.ConvertStringToUint(taskListBasicRequest.OrderId)
		if orderId == 0 {
			return tasksSearch, &commonStructures.CommonError{
				StatusCode: http.StatusBadRequest,
				Message:    commonConstants.ERROR_INVALID_ORDER_ID,
			}
		}
		tasksSearch.OrderId = commonUtils.ConvertStringToUint(taskListBasicRequest.OrderId)
	}

	if taskListBasicRequest.RequestId != "" {
		requestIdString := taskListBasicRequest.RequestId
		activeCityCodes := commonConstants.ActiveCityCodes
		for _, cityCode := range activeCityCodes {
			if strings.HasPrefix(requestIdString, cityCode) {
				requestIdString = strings.TrimPrefix(requestIdString, cityCode)
				tasksSearch.CityCode = cityCode
				break
			}
		}
		requestId := commonUtils.ConvertStringToUint(requestIdString)
		if requestId == 0 {
			return tasksSearch, &commonStructures.CommonError{
				StatusCode: http.StatusBadRequest,
				Message:    commonConstants.ERROR_INVALID_REQUEST_ID,
			}
		}
		tasksSearch.RequestId = requestId
	}

	return tasksSearch, nil
}

func validateAndReturnFilters(taskListBasicRequest structures.TaskListBasicRequest) (
	structures.TasksFilters, *commonStructures.CommonError) {
	taskFilters := structures.TasksFilters{}

	if taskListBasicRequest.Status != "" {
		statuses := []string{}
		statusesString := strings.Split(taskListBasicRequest.Status, ",")
		for _, status := range statusesString {
			if !commonUtils.SliceContainsString(commonConstants.TASK_STATUSES, status) {
				return taskFilters, &commonStructures.CommonError{
					StatusCode: http.StatusBadRequest,
					Message:    commonConstants.ERROR_INVALID_TASK_STATUS,
				}
			}
			statuses = append(statuses, status)
		}
		taskFilters.Status = statuses
	}

	if taskListBasicRequest.Department != "" {
		taskFilters.Department = strings.Split(taskListBasicRequest.Department, ",")
	}

	if taskListBasicRequest.SpecialRequirement != "" {
		specialRequirements := []string{}
		specialRequirementsString := strings.Split(taskListBasicRequest.SpecialRequirement, ",")
		for _, specialRequirement := range specialRequirementsString {
			if !commonUtils.SliceContainsString(constants.SPECIAL_REQUIREMENTS, specialRequirement) {
				return taskFilters, &commonStructures.CommonError{
					StatusCode: http.StatusBadRequest,
					Message:    commonConstants.ERROR_INVALID_SPECIAL_REQUIREMENT,
				}
			}
			specialRequirements = append(specialRequirements, specialRequirement)
		}
		taskFilters.SpecialRequirement = specialRequirements
	}

	if taskListBasicRequest.OrderType != "" {
		orderTypes := []string{}
		orderTypesString := strings.Split(taskListBasicRequest.OrderType, ",")
		for _, orderType := range orderTypesString {
			if !commonUtils.SliceContainsString(commonConstants.ORDER_TYPES, orderType) {
				return taskFilters, &commonStructures.CommonError{
					StatusCode: http.StatusBadRequest,
					Message:    commonConstants.ERROR_INVALID_ORDER_TYPE,
				}
			}
			orderTypes = append(orderTypes, orderType)
		}
		taskFilters.OrderType = orderTypes
	}

	if taskListBasicRequest.LabId != "" {
		labIds := []uint{}
		labIdStrings := strings.Split(taskListBasicRequest.LabId, ",")
		for _, labIdString := range labIdStrings {
			labId := commonUtils.ConvertStringToUint(labIdString)
			if labId == 0 {
				return taskFilters, &commonStructures.CommonError{
					StatusCode: http.StatusBadRequest,
					Message:    commonConstants.ERROR_INVALID_LAB_ID,
				}
			}
			labIds = append(labIds, labId)
		}
		taskFilters.LabId = labIds
	}

	return taskFilters, nil
}

func validateAndReturnTaskListRequestParameters(taskListBasicRequest structures.TaskListBasicRequest) (
	structures.TaskListRequest, *commonStructures.CommonError) {
	taskListRequest := structures.TaskListRequest{}
	if taskListBasicRequest.Limit <= 0 {
		return taskListRequest, &commonStructures.CommonError{
			StatusCode: http.StatusBadRequest,
			Message:    commonConstants.ERROR_INVALID_LIMIT,
		}
	}
	taskListRequest.Limit = taskListBasicRequest.Limit

	taskListRequest.Offset = taskListBasicRequest.Offset

	taskListRequest.UserId = taskListBasicRequest.UserId

	if taskListBasicRequest.TaskTypes == "" {
		return taskListRequest, &commonStructures.CommonError{
			StatusCode: http.StatusBadRequest,
			Message:    commonConstants.ERROR_INVALID_TASK_TYPE,
		}
	}

	taskTypes := strings.Split(taskListBasicRequest.TaskTypes, ",")
	for _, taskType := range taskTypes {
		if !commonUtils.SliceContainsString(constants.TASK_TYPES, taskType) {
			return taskListRequest, &commonStructures.CommonError{
				StatusCode: http.StatusBadRequest,
				Message:    commonConstants.ERROR_INVALID_TASK_TYPE,
			}
		}
	}
	taskListRequest.TaskTypes = taskTypes

	tasksSearch, cErr := validateAndReturnSearchParameters(taskListBasicRequest)
	if cErr != nil {
		return taskListRequest, cErr
	}
	taskListRequest.Search = tasksSearch

	taskFilters, cErr := validateAndReturnFilters(taskListBasicRequest)
	if cErr != nil {
		return taskListRequest, cErr
	}
	taskListRequest.Filters = taskFilters

	return taskListRequest, nil
}

func addWhereClauseParameters(taskListRequest structures.TaskListRequest,
	taskListDbRequest *structures.TaskListDbRequest) {
	patientDetailsClause := map[string]interface{}{}
	tasksClause := map[string]interface{}{}
	visitClause := map[string]interface{}{}
	taskMetadataClause := map[string]interface{}{}
	testDetailsClause := map[string]interface{}{}

	searchCriterion := taskListRequest.Search
	if searchCriterion.PatientName != "" {
		taskListDbRequest.PatientName = strings.ToLower(searchCriterion.PatientName)
	}

	if searchCriterion.PartnerName != "" {
		taskListDbRequest.PartnerName = strings.ToLower(searchCriterion.PartnerName)
	}

	if searchCriterion.DoctorName != "" {
		doctorName := strings.ToLower(searchCriterion.DoctorName)
		doctorName = strings.TrimSpace(strings.TrimPrefix(doctorName, "dr."))
		taskListDbRequest.DoctorName = "dr. " + doctorName
	}

	if searchCriterion.ContactNumber != "" {
		patientDetailsClause["patient_details.number"] = searchCriterion.ContactNumber
	}

	if searchCriterion.PatientId != "" {
		patientDetailsClause["patient_details.system_patient_id"] = searchCriterion.PatientId
	}

	if searchCriterion.OrderId != 0 {
		tasksClause["tasks.order_id"] = searchCriterion.OrderId
	}

	if searchCriterion.VisitId != "" {
		visitClause["task_visit_mapping.visit_id"] = searchCriterion.VisitId
	}

	if searchCriterion.RequestId != 0 {
		tasksClause["tasks.request_id"] = searchCriterion.RequestId
	}

	if searchCriterion.CityCode != "" {
		tasksClause["tasks.city_code"] = searchCriterion.CityCode
	}

	filters := taskListRequest.Filters
	if len(filters.Status) > 0 {
		tasksClause["tasks.status"] = filters.Status
	}

	if len(filters.Department) > 0 {
		testDetailsClause["test_details.department"] = filters.Department
	}

	if len(filters.LabId) > 0 {
		tasksClause["tasks.lab_id"] = filters.LabId
	}

	if len(filters.SpecialRequirement) > 0 {
		for _, specialRequirement := range filters.SpecialRequirement {
			switch specialRequirement {
			case constants.SPECIAL_REQUIREMENT_MORPHLE:
				taskMetadataClause["task_metadata.contains_morphle"] = true
			case constants.SPECIAL_REQUIREMENT_PACKAGE:
				taskMetadataClause["task_metadata.contains_package"] = true
			}
		}
	}

	if len(filters.OrderType) > 0 {
		tasksClause["tasks.order_type"] = filters.OrderType
	}

	taskListDbRequest.PatientDetailsClause = patientDetailsClause
	taskListDbRequest.TasksClause = tasksClause
	taskListDbRequest.VisitClause = visitClause
	taskListDbRequest.TaskMetadataClause = taskMetadataClause
	taskListDbRequest.TestDetailsClause = testDetailsClause
}

func (searchService *SearchService) getTestDetailsAndCreateTaskDetails(ctx context.Context, limit uint,
	taskDbDetails []structures.TaskDetailsDbStruct) (
	[]structures.TaskDetailsStruct, bool, *commonStructures.CommonError,
) {
	taskDetails := []structures.TaskDetailsStruct{}
	showMore := false

	if len(taskDbDetails) == 0 {
		return taskDetails, showMore, nil
	}

	if len(taskDbDetails) > int(limit) {
		showMore = true
		taskDbDetails = taskDbDetails[:limit]
	}

	taskIds := []uint{}
	for _, task := range taskDbDetails {
		taskIds = append(taskIds, task.TaskId)
	}

	basicTestDetails, cErr := searchService.TestDetailService.GetTestBasicDetailsForSearchScreenByTaskIds(taskIds)
	if cErr != nil {
		return taskDetails, showMore, cErr
	}

	labIdLabMap := searchService.CdsService.GetLabIdLabMap(ctx)

	taskTestDetailsMap := map[uint][]structures.TestDetailsStruct{}
	taskProcessingLabIdsMap, taskProcessingLabsMap := map[uint][]uint{}, map[uint][]structures.LabDetailsStruct{}
	for _, testDetail := range basicTestDetails {
		if _, ok := taskTestDetailsMap[testDetail.TaskId]; !ok {
			taskTestDetailsMap[testDetail.TaskId] = []structures.TestDetailsStruct{}
		}
		if _, ok := taskProcessingLabsMap[testDetail.TaskId]; !ok {
			taskProcessingLabsMap[testDetail.TaskId] = []structures.LabDetailsStruct{}
		}
		if _, ok := taskProcessingLabIdsMap[testDetail.TaskId]; !ok {
			taskProcessingLabIdsMap[testDetail.TaskId] = []uint{}
		}
		taskTestDetailsMap[testDetail.TaskId] = append(taskTestDetailsMap[testDetail.TaskId], structures.TestDetailsStruct{
			Name:        testDetail.TestName,
			Status:      testDetail.Status,
			StatusLabel: commonConstants.TEST_STATUSES_LABEL_MAP[testDetail.Status],
		})
		if commonUtils.SliceContainsString(commonConstants.TEST_STATUSES_WITH_PROCESSING_LAB, testDetail.Status) &&
			!commonUtils.SliceContainsUint(taskProcessingLabIdsMap[testDetail.TaskId], testDetail.LabId) {
			taskProcessingLabsMap[testDetail.TaskId] = append(taskProcessingLabsMap[testDetail.TaskId],
				structures.LabDetailsStruct{
					Id:   testDetail.LabId,
					Name: labIdLabMap[testDetail.LabId].LabName,
				})
			taskProcessingLabIdsMap[testDetail.TaskId] = append(taskProcessingLabIdsMap[testDetail.TaskId], testDetail.LabId)
		}
	}

	for _, task := range taskDbDetails {
		taskDetails = append(taskDetails, structures.TaskDetailsStruct{
			TaskId:         task.TaskId,
			OrderId:        task.OrderId,
			OmsOrderId:     task.OmsOrderId,
			CityCode:       task.CityCode,
			DoctorTat:      task.DoctorTat,
			PickedBy:       task.PickedBy,
			CoAuthorizedBy: task.CoAuthorizedBy,
			Status:         task.Status,
			PatientName:    task.PatientName,
			TaskContents: structures.TaskContentStruct{
				ContainsMorphle: task.ContainsMorphle,
				ContainsPackage: task.ContainsPackage,
			},
			TestDetails:    taskTestDetailsMap[task.TaskId],
			ProcessingLabs: taskProcessingLabsMap[task.TaskId],
		})
	}

	return taskDetails, showMore, nil
}

func (searchService *SearchService) getCriticalTasksList(ctx context.Context,
	taskListDbRequest structures.TaskListDbRequest) ([]structures.TaskDetailsStruct, bool, *commonStructures.CommonError) {

	taskDetails, showMore := []structures.TaskDetailsStruct{}, false
	taskDbDetails, cErr := searchService.SearchDao.GetCriticalTaskDetails(taskListDbRequest)
	if cErr != nil {
		return taskDetails, false, cErr
	}

	taskDetails, showMore, cErr = searchService.getTestDetailsAndCreateTaskDetails(ctx, taskListDbRequest.Limit, taskDbDetails)
	if cErr != nil {
		return taskDetails, showMore, cErr
	}

	return taskDetails, showMore, nil
}

func (searchService *SearchService) getWithheldTasksList(ctx context.Context,
	taskListDbRequest structures.TaskListDbRequest) ([]structures.TaskDetailsStruct, bool, *commonStructures.CommonError) {

	taskDetails, showMore := []structures.TaskDetailsStruct{}, false
	taskDbDetails, cErr := searchService.SearchDao.GetWithheldTaskDetails(taskListDbRequest)
	if cErr != nil {
		return taskDetails, showMore, cErr
	}

	taskDetails, showMore, cErr = searchService.getTestDetailsAndCreateTaskDetails(ctx, taskListDbRequest.Limit, taskDbDetails)
	if cErr != nil {
		return taskDetails, showMore, cErr
	}

	return taskDetails, showMore, nil
}

func (searchService *SearchService) getCoAuthorizedTasksList(ctx context.Context,
	taskListDbRequest structures.TaskListDbRequest) ([]structures.TaskDetailsStruct, bool, *commonStructures.CommonError) {

	taskDetails, showMore := []structures.TaskDetailsStruct{}, false
	taskDbDetails, cErr := searchService.SearchDao.GetCoAuthorizedTaskDetails(taskListDbRequest)
	if cErr != nil {
		return taskDetails, showMore, cErr
	}

	taskDetails, showMore, cErr = searchService.getTestDetailsAndCreateTaskDetails(ctx, taskListDbRequest.Limit, taskDbDetails)
	if cErr != nil {
		return taskDetails, showMore, cErr
	}

	return taskDetails, showMore, nil
}

func (searchService *SearchService) getNormalTasksList(ctx context.Context,
	taskListDbRequest structures.TaskListDbRequest) ([]structures.TaskDetailsStruct, bool, *commonStructures.CommonError) {

	taskDetails, showMore := []structures.TaskDetailsStruct{}, false
	taskDbDetails, cErr := searchService.SearchDao.GetNormalTaskDetails(taskListDbRequest)
	if cErr != nil {
		return taskDetails, showMore, cErr
	}

	taskDetails, showMore, cErr = searchService.getTestDetailsAndCreateTaskDetails(ctx, taskListDbRequest.Limit, taskDbDetails)
	if cErr != nil {
		return taskDetails, showMore, cErr
	}

	return taskDetails, showMore, nil
}

func (searchService *SearchService) getInProgressTasksList(ctx context.Context,
	taskListDbRequest structures.TaskListDbRequest) ([]structures.TaskDetailsStruct, bool, *commonStructures.CommonError) {

	taskDetails, showMore := []structures.TaskDetailsStruct{}, false
	taskDbDetails, cErr := searchService.SearchDao.GetInProgressTaskDetails(taskListDbRequest)
	if cErr != nil {
		return taskDetails, showMore, cErr
	}

	taskDetails, showMore, cErr = searchService.getTestDetailsAndCreateTaskDetails(ctx, taskListDbRequest.Limit, taskDbDetails)
	if cErr != nil {
		return taskDetails, showMore, cErr
	}

	return taskDetails, showMore, nil
}

func (searchService *SearchService) getTaskListAndCreateResponse(ctx context.Context,
	taskListRequest structures.TaskListRequest,
	taskListDbRequest structures.TaskListDbRequest) (
	structures.TaskListResponse, []*commonStructures.CommonError) {

	var cErr *commonStructures.CommonError
	var errList []*commonStructures.CommonError
	criticalTasks, withheldTasks := []structures.TaskDetailsStruct{}, []structures.TaskDetailsStruct{}
	coAuthorizedTasks, normalTasks := []structures.TaskDetailsStruct{}, []structures.TaskDetailsStruct{}
	inProgressTasks := []structures.TaskDetailsStruct{}
	criticalShowMore, withheldShowMore := false, false
	coAuthorizedShowMore, normalShowMore, inProgressShowMore := false, false, false

	wg, mu := sync.WaitGroup{}, sync.Mutex{}
	wg.Add(len(taskListRequest.TaskTypes))

	for _, taskType := range taskListRequest.TaskTypes {
		go func(taskType string, taskListDbRequest structures.TaskListDbRequest) {
			defer wg.Done()

			switch taskType {
			case constants.TASK_TYPE_CRITICAL:
				criticalTasks, criticalShowMore, cErr = searchService.getCriticalTasksList(ctx, taskListDbRequest)
				if cErr != nil {
					mu.Lock()
					errList = append(errList, cErr)
					mu.Unlock()
				}
			case constants.TASK_TYPE_WITHHELD:
				withheldTasks, withheldShowMore, cErr = searchService.getWithheldTasksList(ctx, taskListDbRequest)
				if cErr != nil {
					mu.Lock()
					errList = append(errList, cErr)
					mu.Unlock()
				}
			case constants.TASK_TYPE_CO_AUTHORIZE:
				coAuthorizedTasks, coAuthorizedShowMore, cErr = searchService.getCoAuthorizedTasksList(ctx, taskListDbRequest)
				if cErr != nil {
					mu.Lock()
					errList = append(errList, cErr)
					mu.Unlock()
				}
			case constants.TASK_TYPE_NORMAL:
				normalTasks, normalShowMore, cErr = searchService.getNormalTasksList(ctx, taskListDbRequest)
				if cErr != nil {
					mu.Lock()
					errList = append(errList, cErr)
					mu.Unlock()
				}
			case constants.TASK_TYPE_IN_PROGRESS:
				inProgressTasks, inProgressShowMore, cErr = searchService.getInProgressTasksList(ctx, taskListDbRequest)
				if cErr != nil {
					mu.Lock()
					errList = append(errList, cErr)
					mu.Unlock()
				}
			}
		}(taskType, taskListDbRequest)
	}

	wg.Wait()

	if len(errList) > 0 {
		commonUtils.AddLog(ctx, commonConstants.ERROR_LEVEL, commonConstants.ERROR_WHILE_FETCHING_TASKS, map[string]interface{}{
			"errors": errList,
		}, nil)
		return structures.TaskListResponse{}, errList
	}

	response := structures.TaskListResponse{
		CriticalTasks: structures.TaskResponseStruct{
			Tasks:    criticalTasks,
			ShowMore: criticalShowMore,
		},
		WithheldTasks: structures.TaskResponseStruct{
			Tasks:    withheldTasks,
			ShowMore: withheldShowMore,
		},
		CoAuthorizeTasks: structures.TaskResponseStruct{
			Tasks:    coAuthorizedTasks,
			ShowMore: coAuthorizedShowMore,
		},
		NormalTasks: structures.TaskResponseStruct{
			Tasks:    normalTasks,
			ShowMore: normalShowMore,
		},
		InProgressTasks: structures.TaskResponseStruct{
			Tasks:    inProgressTasks,
			ShowMore: inProgressShowMore,
		},
	}

	return response, nil
}

func (searchService *SearchService) GetTasksList(ctx context.Context, taskListBasicRequest structures.TaskListBasicRequest) (
	structures.TaskListResponse, *commonStructures.CommonError) {

	taskListRequest, cErr := validateAndReturnTaskListRequestParameters(taskListBasicRequest)
	if cErr != nil {
		return structures.TaskListResponse{}, cErr
	}

	taskListDbRequest := structures.TaskListDbRequest{
		UserId: taskListRequest.UserId,
		Limit:  taskListRequest.Limit,
		Offset: taskListRequest.Offset,
	}

	addWhereClauseParameters(taskListRequest, &taskListDbRequest)

	response, errList := searchService.getTaskListAndCreateResponse(ctx, taskListRequest, taskListDbRequest)
	if len(errList) > 0 {
		return structures.TaskListResponse{}, errList[0]
	}

	return response, nil
}

func (searchService *SearchService) getAmendmentTasksListAndCreateResponse(ctx context.Context,
	taskListDbRequest structures.TaskListDbRequest) (
	structures.AmendmentTaskListResponse, *commonStructures.CommonError) {

	var cErr *commonStructures.CommonError
	var cErrList []*commonStructures.CommonError
	amendmentTasks, count := []structures.TaskDetailsStruct{}, uint(0)

	wg := sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()
		taskDbDetails, cErr := searchService.SearchDao.GetAmendmentTaskDetails(taskListDbRequest)
		if cErr != nil {
			cErrList = append(cErrList, cErr)
		}

		for _, task := range taskDbDetails {
			amendmentTasks = append(amendmentTasks, structures.TaskDetailsStruct{
				TaskId:      task.TaskId,
				OrderId:     task.OrderId,
				CityCode:    task.CityCode,
				PatientName: task.PatientName,
				TaskContents: structures.TaskContentStruct{
					ContainsMorphle: task.ContainsMorphle,
					ContainsPackage: task.ContainsPackage,
				},
			})
		}
	}()

	go func() {
		defer wg.Done()

		count, cErr = searchService.TaskService.GetAmendmentTasksCount(ctx)
		if cErr != nil {
			cErrList = append(cErrList, cErr)
		}
	}()

	wg.Wait()

	if len(cErrList) > 0 {
		commonUtils.AddLog(ctx, commonConstants.ERROR_LEVEL, commonConstants.ERROR_WHILE_FETCHING_TASKS, map[string]interface{}{
			"errors": cErrList,
		}, nil)
		return structures.AmendmentTaskListResponse{}, cErrList[0]
	}

	response := structures.AmendmentTaskListResponse{
		AmendmentTasks: amendmentTasks,
		Count:          count,
	}

	return response, nil
}

func (searchService *SearchService) GetAmendmentTasksList(ctx context.Context, taskListBasicRequest structures.TaskListBasicRequest) (
	structures.AmendmentTaskListResponse, *commonStructures.CommonError) {
	taskListRequest, cErr := validateAndReturnTaskListRequestParameters(taskListBasicRequest)
	if cErr != nil {
		return structures.AmendmentTaskListResponse{}, cErr
	}

	taskListDbRequest := structures.TaskListDbRequest{
		UserId: taskListRequest.UserId,
		Limit:  taskListRequest.Limit,
		Offset: taskListRequest.Offset,
	}

	addWhereClauseParameters(taskListRequest, &taskListDbRequest)

	response, cErr := searchService.getAmendmentTasksListAndCreateResponse(ctx, taskListDbRequest)
	if cErr != nil {
		return structures.AmendmentTaskListResponse{}, cErr
	}

	return response, nil
}
