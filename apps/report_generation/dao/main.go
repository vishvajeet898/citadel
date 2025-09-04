package dao

import (
	"gorm.io/gorm"

	"github.com/Orange-Health/citadel/adapters/psql"
)

type ReportGenerationDao struct {
	Db *gorm.DB
}

func InitializeReportGenerationDao() DataLayer {
	return &ReportGenerationDao{
		Db: psql.GetDbInstance(),
	}
}
