package dao

import (
	"gorm.io/gorm"

	"github.com/Orange-Health/citadel/adapters/psql"
)

type CoAuthorizePathologistDao struct {
	Db *gorm.DB
}

func InitializeoAuthorizePathologistDao() DataLayer {
	return &CoAuthorizePathologistDao{
		Db: psql.GetDbInstance(),
	}
}
