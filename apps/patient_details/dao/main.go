package dao

import (
	"gorm.io/gorm"

	"github.com/Orange-Health/citadel/adapters/psql"
)

type PatientDetailDao struct {
	Db *gorm.DB
}

func InitializePatientDetailDao() DataLayer {
	return &PatientDetailDao{
		Db: psql.GetDbInstance(),
	}
}
