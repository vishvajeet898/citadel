package controller

import (
	"github.com/Orange-Health/citadel/apps/audit_log/service"
)

type AuditLog struct {
	AuditLogService service.AuditLogServiceInterface
}

func InitAuditLogController() *AuditLog {
	return &AuditLog{
		AuditLogService: service.InitializeAuditLogService(),
	}
}
