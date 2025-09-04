package dao

import (
	"gorm.io/gorm"

	"github.com/Orange-Health/citadel/adapters/psql"
)

type AuditLogDao struct {
	Db *gorm.DB
}

func InitializeAuditLogDao() DataLayer {
	return &AuditLogDao{
		Db: psql.GetDbInstance(),
	}
}
