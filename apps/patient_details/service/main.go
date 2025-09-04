package service

import (
	"github.com/Orange-Health/citadel/adapters/cache"
	"github.com/Orange-Health/citadel/adapters/sentry"
	"github.com/Orange-Health/citadel/apps/patient_details/dao"
	omsClient "github.com/Orange-Health/citadel/clients/oms"
)

type PatientDetailService struct {
	PatientDetailDao dao.DataLayer
	Cache            cache.CacheLayer
	Sentry           sentry.SentryLayer
	OmsClient        omsClient.OmsClientInterface
}

func InitializePatientDetailService() PatientDetailServiceInterface {
	return &PatientDetailService{
		PatientDetailDao: dao.InitializePatientDetailDao(),
		Cache:            cache.InitializeCache(),
		Sentry:           sentry.InitializeSentry(),
		OmsClient:        omsClient.InitializeOmsClient(),
	}
}
