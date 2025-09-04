package dao

import (
	"gorm.io/gorm"

	"github.com/Orange-Health/citadel/adapters/psql"
)

type TaskPathMapDao struct {
	Db *gorm.DB
}

func InitializeTaskPathMapDao() DataLayer {
	return &TaskPathMapDao{
		Db: psql.GetDbInstance(),
	}
}
