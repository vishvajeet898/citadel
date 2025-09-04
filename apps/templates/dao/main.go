package dao

import (
	"gorm.io/gorm"

	"github.com/Orange-Health/citadel/adapters/psql"
)

type TemplateDao struct {
	Db *gorm.DB
}

func InitializeTemplateDao() DataLayer {
	return &TemplateDao{
		Db: psql.GetDbInstance(),
	}
}
