package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/Orange-Health/citadel/apps/task/structures"
	commonConstants "github.com/Orange-Health/citadel/common/constants"
	commonStructs "github.com/Orange-Health/citadel/common/structures"
	commonUtils "github.com/Orange-Health/citadel/common/utils"
)

// @Summary		GetTaskById
// @Description	Get Task Details for the given Task ID
// @Tags			tasks
// @Produce		json
// @Param			taskId		path		int						true	"Task ID"
// @Success		200			{object}	structures.TaskDetail	"Common API Response"
// @Failure		400,404,500	{object}	structures.CommonError	"Common API Response"
// @Router			/api/v1/tasks/{taskId}/details [get]
func (taskController *Task) GetTaskDetails(c *gin.Context) {
	taskID := commonUtils.ConvertStringToUint(c.Param("taskId"))
	if taskID == 0 {
		commonUtils.HandleErrorResponse(c, http.StatusBadRequest, commonConstants.ERROR_INVALID_TASK_ID)
		return
	}

	taskDetail, err := taskController.TaskService.GetTaskById(taskID)
	if err != nil {
		commonUtils.HandleErrorResponse(c, err.StatusCode, err.Message)
		return
	}

	c.JSON(http.StatusOK, taskDetail)
}

// @Summary		GetTasksList
// @Description	Get all tasks based on filters provided
// @Tags			tasks
// @Produce		json
// @Success		200			{object}	structures.CommonAPIResponse	"Common API Response"
// @Failure		400,404,500	{object}	structures.CommonAPIResponse	"Common API Response"
// @Router			/api/v1/tasks/count [get]
func (taskController *Task) GetTasksCount(c *gin.Context) {
	taskCount, cErr := taskController.TaskService.GetTasksCount()
	if cErr != nil {
		commonUtils.HandleErrorResponse(c, cErr.StatusCode, cErr.Message)
		return
	}

	c.JSON(http.StatusOK, commonStructs.CommonAPIResponse{
		Data: taskCount,
	})
}

// @Summary		UpdateAllTaskDetails
// @Description	Update all task details
// @Tags			tasks
// @Accept			json
// @Produce		json
// @Param			taskDetails	body		structures.UpdateAllTaskDetailsStruct	true	"Task Details"
// @Success		200			{object}	structures.CommonAPIResponse			"Common API Response"
// @Failure		400,404,500	{object}	structures.CommonAPIResponse			"Common API Response"
// @Router			/api/v1/tasks/details [patch]
func (taskController *Task) UpdateAllTaskDetails(c *gin.Context) {
	userId, cErr := commonUtils.GetUserIdFromContext(c)
	if cErr != nil {
		commonUtils.HandleErrorResponse(c, cErr.StatusCode, cErr.Message)
		return
	}

	taskUpdateRequest := structures.UpdateAllTaskDetailsStruct{}
	if err := c.ShouldBindJSON(&taskUpdateRequest); err != nil {
		commonUtils.HandleErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	taskUpdateRequest.UserId = userId

	cErr = taskController.TaskService.UpdateAllTaskDetails(c.Request.Context(), taskUpdateRequest)
	if cErr != nil {
		commonUtils.HandleErrorResponse(c, cErr.StatusCode, cErr.Message)
		return
	}

	_ = taskController.TaskWorkerService.UpdateTaskPostSavingTask(c.Request.Context(), taskUpdateRequest.Task.Id)
	_ = taskController.TaskWorkerService.ReleaseReportTask(c.Request.Context(), taskUpdateRequest.Task.Id, false)

	c.JSON(http.StatusOK, commonStructs.CommonAPIResponse{
		Message: commonConstants.DONE_RESPONSE,
	})
}

// @Summary		UndoReleaseReport
// @Description	Undo Release Report for the given Task ID
// @Tags			tasks
// @Produce		json
// @Param			taskId		path		int								true	"Task ID"
// @Success		200			{object}	structures.CommonAPIResponse	"Common API Response"
// @Success		400,404,500	{object}	structures.CommonAPIResponse	"Common API Response"
// @Router			/api/v1/tasks/{taskId}/undo-release [patch]
func (taskController *Task) UndoReportRelease(c *gin.Context) {
	taskId := commonUtils.ConvertStringToUint(c.Param("taskId"))
	cErr := taskController.TaskService.UndoReportRelease(c.Request.Context(), taskId)
	if cErr != nil {
		commonUtils.HandleErrorResponse(c, cErr.StatusCode, cErr.Message)
		return
	}

	c.JSON(http.StatusOK, commonStructs.CommonAPIResponse{
		Message: commonConstants.DONE_RESPONSE,
	})
}

// @Summary		GetCallingTaskDetails
// @Description	Get Calling Task Details
// @Tags			tasks
// @Produce		json
// @Param			taskId			path		int										true	"Task ID"
// @Param			calling_type	query		string									false	"Calling Type"
// @Success		200				{object}	structures.TaskCallingDetailsResponse	"Calling Details Response"
// @Failure		400,404,500		{object}	structures.CommonAPIResponse			"Common API Response"
// @Router			/api/v1/tasks/{taskId}/calling-details [get]
func (taskController *Task) GetTaskCallingDetails(c *gin.Context) {
	taskId := commonUtils.ConvertStringToUint(c.Param("taskId"))
	if taskId == 0 {
		commonUtils.HandleErrorResponse(c, http.StatusBadRequest, commonConstants.ERROR_INVALID_TASK_ID)
		return
	}

	callingType := c.Query("calling_type")

	userId, cErr := commonUtils.GetUserIdFromContext(c)
	if cErr != nil {
		commonUtils.HandleErrorResponse(c, cErr.StatusCode, cErr.Message)
		return
	}

	taskDetail, err := taskController.TaskService.GetTaskCallingDetails(taskId, userId, callingType)
	if err != nil {
		commonUtils.HandleErrorResponse(c, err.StatusCode, err.Message)
		return
	}

	c.JSON(http.StatusOK, taskDetail)
}

func (taskController *Task) CreateUpdateTask(c *gin.Context) {
	orderId := c.Param("orderId")
	if orderId == "" {
		commonUtils.HandleErrorResponse(c, http.StatusBadRequest, commonConstants.ERROR_INVALID_ORDER_ID)
		return
	}

	err := taskController.TaskWorkerService.CreateUpdateTaskByOmsOrderIdTask(c.Request.Context(), orderId)
	if err != nil {
		commonUtils.AddLog(c.Request.Context(), commonConstants.ERROR_LEVEL, commonUtils.GetCurrentFunctionName(),
			map[string]interface{}{
				"err": commonConstants.ERROR_FAILED_TO_CREATE_UPDATE_TASK,
			}, err)
		return
	}

	c.JSON(http.StatusAccepted, commonStructs.CommonAPIResponse{
		Message: commonConstants.DONE_RESPONSE,
	})
}

func (taskController *Task) TriggerReportGenerationAndAttuneEventByVisitId(c *gin.Context) {
	visitId := c.Param("visitId")
	if visitId == "" {
		commonUtils.HandleErrorResponse(c, http.StatusBadRequest, commonConstants.ERROR_INVALID_VISIT_ID)
		return
	}

	taskId, cErr := taskController.TaskService.GetTaskIdByVisitId(visitId)
	if cErr != nil {
		commonUtils.HandleErrorResponse(c, cErr.StatusCode, cErr.Message)
		return
	}

	if taskId != 0 {
		_ = taskController.TaskWorkerService.ReleaseReportTask(c.Request.Context(), taskId, true)
	}

	c.JSON(http.StatusAccepted, commonStructs.CommonAPIResponse{
		Message: commonConstants.DONE_RESPONSE,
	})
}
