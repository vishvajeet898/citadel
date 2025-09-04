package dao

import (
	"gorm.io/gorm"

	"github.com/Orange-Health/citadel/adapters/psql"
)

type TestSampleMappingDao struct {
	Db *gorm.DB
}

func InitializeTestSampleMappingDao() DataLayer {
	return &TestSampleMappingDao{
		Db: psql.GetDbInstance(),
	}
}
