package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/Orange-Health/citadel/apps/task_metadata/structures"
	commonConstants "github.com/Orange-Health/citadel/common/constants"
	commonStructures "github.com/Orange-Health/citadel/common/structures"
	commonUtils "github.com/Orange-Health/citadel/common/utils"
)

// @Summary		GetTaskMetadataDetails
// @Description	Get Task Metadata Details for the given Task ID
// @Tags			task-metadata
// @Produce		json
// @Param			taskId		path		int						true	"Task ID"
// @Success		200			{object}	structures.TaskMetadata	"Common API Response"
// @Success		400,404,500	{object}	structures.CommonError	"Common API Response"
// @Router			/api/v1/task-metadata/task/{taskId}/details [get]
func (taskMetadataController *TaskMetadata) GetTaskMetadataDetails(ctx *gin.Context) {
	taskID := commonUtils.ConvertStringToUint(ctx.Param("taskId"))
	if taskID == 0 {
		commonUtils.HandleErrorResponse(ctx, http.StatusBadRequest, commonConstants.ERROR_INVALID_TASK_ID)
		return
	}
	taskMetadata, err := taskMetadataController.TaskMetadataService.GetTaskMetadataDetails(taskID)
	if err != nil {
		commonUtils.HandleErrorResponse(ctx, err.StatusCode, err.Message)
		return
	}

	ctx.JSON(http.StatusOK, taskMetadata)
}

// @Summary		UpdateTaskMetadata
// @Description	Update Task Metadata Details for the given Task ID
// @Tags			task-metadata
// @Accept			json
// @Produce		json
// @Param			taskId			path		int						true	"Task ID"
// @Param			taskMetadata	body		models.TaskMetadata		true	"Task Metadata Object"
// @Success		200				{object}	models.TaskMetadata		"Common API Response"
// @Success		400,404,500		{object}	structures.CommonError	"Common API Response"
// @Router			/api/v1/task-metadata/tasks/:taskId [patch]
func (taskMetadataController *TaskMetadata) UpdateTaskMetadata(ctx *gin.Context) {
	var tmd structures.TaskMetadata
	taskID := commonUtils.ConvertStringToUint(ctx.Param("taskId"))
	if taskID == 0 {
		commonUtils.HandleErrorResponse(ctx, http.StatusBadRequest, commonConstants.ERROR_INVALID_TASK_ID)
		return
	}

	if err := ctx.BindJSON(&tmd); err != nil {
		commonUtils.HandleErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	tmd.TaskID = taskID
	taskMetadata, err := taskMetadataController.TaskMetadataService.UpdateTaskMetadata(tmd)
	if err != nil {
		commonUtils.HandleErrorResponse(ctx, err.StatusCode, err.Message)
		return
	}

	ctx.JSON(http.StatusOK, taskMetadata)
}

// @Summary		UpdateLastEventSentAt
// @Description	Update Last Event Sent At for the given Task ID
// @Tags			task-metadata
// @Produce		json
// @Param			taskId		path		int						true	"Task ID"
// @Success		200			{object}	string					"Common API Response"
// @Success		400,404,500	{object}	structures.CommonError	"Common API Response"
// @Router			/api/v1/task-metadata/task/{taskId}/last-event-sent-at [patch]
func (taskMetadataController *TaskMetadata) UpdateLastEventSentAt(c *gin.Context) {
	taskID := commonUtils.ConvertStringToUint(c.Param("taskId"))
	if taskID == 0 {
		commonUtils.HandleErrorResponse(c, http.StatusBadRequest, commonConstants.ERROR_INVALID_TASK_ID)
		return
	}

	err := taskMetadataController.TaskMetadataService.UpdateLastEventSentAt(taskID)
	if err != nil {
		commonUtils.HandleErrorResponse(c, err.StatusCode, err.Message)
		return
	}

	c.JSON(http.StatusOK, commonStructures.CommonAPIResponse{
		Message: commonConstants.DONE_RESPONSE,
	})
}
