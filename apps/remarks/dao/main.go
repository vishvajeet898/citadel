package dao

import (
	"gorm.io/gorm"

	"github.com/Orange-Health/citadel/adapters/psql"
)

type RemarkDao struct {
	Db *gorm.DB
}

func InitializeRemarkDao() DataLayer {
	return &RemarkDao{
		Db: psql.GetDbInstance(),
	}
}
