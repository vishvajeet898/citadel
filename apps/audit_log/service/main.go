package service

import (
	"github.com/Orange-Health/citadel/adapters/cache"
	"github.com/Orange-Health/citadel/adapters/sentry"
	"github.com/Orange-Health/citadel/apps/audit_log/dao"
	userService "github.com/Orange-Health/citadel/apps/users/service"
)

type AuditLogService struct {
	AuditLogDao dao.DataLayer
	Cache       cache.CacheLayer
	Sentry      sentry.SentryLayer
	UserService userService.UserServiceInterface
}

func InitializeAuditLogService() AuditLogServiceInterface {
	return &AuditLogService{
		AuditLogDao: dao.InitializeAuditLogDao(),
		Cache:       cache.InitializeCache(),
		Sentry:      sentry.InitializeSentry(),
		UserService: userService.InitializeUserService(),
	}
}
