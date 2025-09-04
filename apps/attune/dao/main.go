package dao

import (
	"time"

	"gorm.io/gorm"

	"github.com/Orange-Health/citadel/adapters/psql"
	commonStructures "github.com/Orange-Health/citadel/common/structures"
	commonModels "github.com/Orange-Health/citadel/models"
)

type AttuneDao struct {
	Db *gorm.DB
}

type DataLayer interface {
	GetOrderDetailsByVisitId(visitID string) (commonModels.OrderDetails, *commonStructures.CommonError)
	GetSampleByVisitId(visitId string) (commonModels.Sample, *commonStructures.CommonError)
	GetOrderDetailsAndPatientDetailsByVisitId(visitId string) (commonModels.OrderDetails, commonModels.PatientDetail,
		*commonStructures.CommonError)
	GetSampleCollectedAtByVisitId(visitId string) (*time.Time, *commonStructures.CommonError)
	GetSampleCollectedAtBySampleIds(sampleIds []uint) (*time.Time, *commonStructures.CommonError)
	GetOrderDetailsAndPatientDetailsByOmsOrderId(omsOrderId string) (commonModels.OrderDetails, commonModels.PatientDetail,
		*commonStructures.CommonError)
	GetAttuneTestSampleMapByOmsOrderId(omsOrderId string, sampleIds []uint) ([]commonStructures.AttuneTestSampleMapSnakeCase,
		*commonStructures.CommonError)
	GetTestDetailsForSyncingToAttune(omsOrderId string, sampleIds []uint, labId uint) ([]commonModels.TestDetail,
		*commonStructures.CommonError)
}

func InitializeAttuneDao() DataLayer {
	return &AttuneDao{
		Db: psql.GetDbInstance(),
	}
}
