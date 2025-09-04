package service

import (
	"github.com/Orange-Health/citadel/adapters/cache"
	"github.com/Orange-Health/citadel/adapters/sentry"
	"github.com/Orange-Health/citadel/apps/external_investigation_results/dao"
	patientDetailsService "github.com/Orange-Health/citadel/apps/patient_details/service"
)

type ExternalInvestigationResultService struct {
	ExtInvResDao          dao.DataLayer
	Cache                 cache.CacheLayer
	Sentry                sentry.SentryLayer
	PatientDetailsService patientDetailsService.PatientDetailServiceInterface
}

func InitializeExternalInvestigationResultService() ExternalInvestigationResultServiceInterface {
	return &ExternalInvestigationResultService{
		ExtInvResDao:          dao.InitializeExternalInvestigationResultDao(),
		Cache:                 cache.InitializeCache(),
		Sentry:                sentry.InitializeSentry(),
		PatientDetailsService: patientDetailsService.InitializePatientDetailService(),
	}
}
