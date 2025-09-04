package dao

import (
	"gorm.io/gorm"

	"github.com/Orange-Health/citadel/adapters/psql"
)

type CdsDao struct {
	Db *gorm.DB
}

type DataLayer interface {
}

func InitializeCdsDao() DataLayer {
	return &CdsDao{
		Db: psql.GetDbInstance(),
	}
}
