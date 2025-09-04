package dao

import (
	"gorm.io/gorm"

	"github.com/Orange-Health/citadel/adapters/psql"
)

type OrderDetailsDao struct {
	Db *gorm.DB
}

func InitializeOrderDetailsDao() DataLayer {
	return &OrderDetailsDao{
		Db: psql.GetDbInstance(),
	}
}
