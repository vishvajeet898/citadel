package dao

import (
	"gorm.io/gorm"

	"github.com/Orange-Health/citadel/adapters/psql"
)

type AttachmentDao struct {
	Db *gorm.DB
}

func InitializeAttachmentDao() DataLayer {
	return &AttachmentDao{
		Db: psql.GetDbInstance(),
	}
}
