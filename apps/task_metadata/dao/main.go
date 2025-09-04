package dao

import (
	"gorm.io/gorm"

	"github.com/Orange-Health/citadel/adapters/psql"
)

type TaskMetadataDao struct {
	Db *gorm.DB
}

func InitializeTaskMetadataDao() DataLayer {
	return &TaskMetadataDao{
		Db: psql.GetDbInstance(),
	}
}
