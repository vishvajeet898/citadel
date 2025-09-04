package dao

import (
	"gorm.io/gorm"

	"github.com/Orange-Health/citadel/adapters/psql"
)

type RerunDao struct {
	Db *gorm.DB
}

func InitializeRerunDao() DataLayer {
	return &RerunDao{
		Db: psql.GetDbInstance(),
	}
}
