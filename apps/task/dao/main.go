package dao

import (
	"gorm.io/gorm"

	"github.com/Orange-Health/citadel/adapters/psql"
)

type TaskDao struct {
	Db *gorm.DB
}

func InitializeTaskDao() DataLayer {
	return &TaskDao{
		Db: psql.GetDbInstance(),
	}
}
