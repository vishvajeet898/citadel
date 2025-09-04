package dao

import (
	"gorm.io/gorm"

	"github.com/Orange-Health/citadel/adapters/psql"
)

type TestDetailDao struct {
	Db *gorm.DB
}

func InitializeTestDetailDao() DataLayer {
	return &TestDetailDao{
		Db: psql.GetDbInstance(),
	}
}
