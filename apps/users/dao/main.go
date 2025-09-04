package dao

import (
	"gorm.io/gorm"

	"github.com/Orange-Health/citadel/adapters/psql"
)

type UserDao struct {
	Db *gorm.DB
}

func InitializeUserDao() DataLayer {
	return &UserDao{
		Db: psql.GetDbInstance(),
	}
}
