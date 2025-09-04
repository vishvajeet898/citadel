package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/Orange-Health/citadel/apps/attachments/structures"
	commonConstants "github.com/Orange-Health/citadel/common/constants"
	commonUtils "github.com/Orange-Health/citadel/common/utils"
)

// @Summary		Get Attachments
// @Description	Get all attachments by task id
// @Tags			attachments
// @Produce		json
// @Param			taskId	path		int									true	"Task ID"
// @Success		200		{object}	map[string][]structures.Attachment	"Attachments"
// @Router			/api/v1/attachments/tasks/{taskId} [get]
func (attachmentController *Attachment) GetAttachmentsByTaskId(c *gin.Context) {
	taskId := commonUtils.ConvertStringToUint(c.Param("taskId"))
	if taskId == 0 {
		commonUtils.HandleErrorResponse(c, http.StatusBadRequest, commonConstants.ERROR_INVALID_TASK_ID)
		return
	}

	attachments, cErr := attachmentController.AttachmentService.GetAttachmentsByTaskId(taskId, nil)
	if cErr != nil {
		commonUtils.HandleErrorResponse(c, cErr.StatusCode, cErr.Message)
		return
	}

	c.JSON(http.StatusOK, attachments)
}

// @Summary		Add Attachment
// @Description	Add attachment by Task ID and Investigation ID
// @Tags			attachments
// @Produce		json
// @Param			taskId	path		int						true	"Task ID"
// @Success		200		{object}	structures.Attachment	"Attachments"
// @Router			/api/v1/attachments/tasks/{taskId} [post]
func (attachmentController *Attachment) AddAttachmentByTaskId(c *gin.Context) {
	taskId := commonUtils.ConvertStringToUint(c.Param("taskId"))
	if taskId == 0 {
		commonUtils.HandleErrorResponse(c, http.StatusBadRequest, commonConstants.ERROR_INVALID_TASK_ID)
		return
	}
	addAttachmentRequest := structures.AddAttachmentRequest{}
	if err := c.ShouldBindJSON(&addAttachmentRequest); err != nil {
		commonUtils.HandleErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	attachment, cErr := attachmentController.AttachmentService.AddAttachment(c.Request.Context(),
		taskId, addAttachmentRequest)
	if cErr != nil {
		commonUtils.HandleErrorResponse(c, cErr.StatusCode, cErr.Message)
		return
	}

	c.JSON(http.StatusCreated, attachment)
}
