package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"

	commonUtils "github.com/Orange-Health/citadel/common/utils"
)

// @Summary		Get Audit Logs
// @Description	Get all audit logs by order id
// @Tags			audit_logs
// @Produce		json
// @Param			order_id	path		int						true	"Task ID"
// @Success		200		{object}	[]structures.AuditLog	"Audit Logs"
// @Failure		400		{object}	structures.CommonError	"Common API Response"
// @Router			/api/v1/audit_logs/tasks/{taskId} [get]
func (auditLogController *AuditLog) GetLogsByOrderId(c *gin.Context) {
	omsOrderId := c.Param("orderId")
	logs, cErr := auditLogController.AuditLogService.GetLogsByOrderId(omsOrderId)

	if cErr != nil {
		commonUtils.HandleErrorResponse(c, cErr.StatusCode, cErr.Message)
		return
	}

	c.JSON(http.StatusOK, logs)
}
