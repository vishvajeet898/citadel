package dao

import (
	"gorm.io/gorm"

	"github.com/Orange-Health/citadel/adapters/psql"
)

type ReceivingDeskDao struct {
	Db *gorm.DB
}

type DataLayer interface {
	BeginTransaction() *gorm.DB
}

func InitializeReceivingDeskDao() DataLayer {
	return &ReceivingDeskDao{
		Db: psql.GetDbInstance(),
	}
}
