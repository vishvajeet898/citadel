package dao

import (
	"gorm.io/gorm"

	"github.com/Orange-Health/citadel/adapters/psql"
)

type InvestigationResultDao struct {
	Db *gorm.DB
}

func InitializeInvestigationResultDao() DataLayer {
	return &InvestigationResultDao{
		Db: psql.GetDbInstance(),
	}
}
