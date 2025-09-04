package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/Orange-Health/citadel/apps/task_pathologist_mapping/structures"
	commonUtils "github.com/Orange-Health/citadel/common/utils"
)

// @Summary		Create Task Pathologist Mapping
// @Description	Create Task Pathologist Mapping
// @Tags			task-pathologist-mapping
// @Accept			json
// @Produce		json
// @Param			structures.taskPathologistMapping	body		structures.TaskPathologistMapping	true	"Task Pathologist Mapping"
// @Success		201									{object}	structures.TaskPathologistMapping
// @Router			/ [post]
func (taskPathController *TaskPathologistMapping) CreateTPM(c *gin.Context) {

	var taskPathologistMapping structures.TaskPathologistMapping
	if err := c.ShouldBindJSON(&taskPathologistMapping); err != nil {
		commonUtils.HandleErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	userId, cErr := commonUtils.GetUserIdFromContext(c)
	if cErr != nil {
		commonUtils.HandleErrorResponse(c, cErr.StatusCode, cErr.Message)
		return
	}
	taskPathologistMapping.PathologistID = userId

	createdTaskPathologistMapping, err := taskPathController.TaskPathService.CreateTaskPathMap(c.Request.Context(),
		taskPathologistMapping)
	if err != nil {
		commonUtils.HandleErrorResponse(c, err.StatusCode, err.Message)
		return
	}

	c.JSON(http.StatusCreated, createdTaskPathologistMapping)
}
