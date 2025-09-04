package audit_log

import (
	"github.com/gin-gonic/gin"

	"github.com/Orange-Health/citadel/apps/audit_log/controller"
)

func RouteHandler(router *gin.RouterGroup) {
	auditLogController := controller.InitAuditLogController()

	router.GET("/orders/:orderId", auditLogController.GetLogsByOrderId)
}
