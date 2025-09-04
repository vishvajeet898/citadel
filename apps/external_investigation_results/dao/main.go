package dao

import (
	"gorm.io/gorm"

	"github.com/Orange-Health/citadel/adapters/psql"
)

type ExternalInvestigationResultDao struct {
	Db *gorm.DB
}

func InitializeExternalInvestigationResultDao() DataLayer {
	return &ExternalInvestigationResultDao{
		Db: psql.GetDbInstance(),
	}
}
