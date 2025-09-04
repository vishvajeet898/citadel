package service

import (
	"github.com/Orange-Health/citadel/adapters/cache"
	"github.com/Orange-Health/citadel/adapters/sentry"
	abnormalityService "github.com/Orange-Health/citadel/apps/abnormality/service"
	calculationsService "github.com/Orange-Health/citadel/apps/calculations/service"
	cdsService "github.com/Orange-Health/citadel/apps/cds/service"
	"github.com/Orange-Health/citadel/apps/investigation_results/dao"
	remarksService "github.com/Orange-Health/citadel/apps/remarks/service"
	rerunService "github.com/Orange-Health/citadel/apps/rerun/service"
	cdsClient "github.com/Orange-Health/citadel/clients/cds"
	omsClient "github.com/Orange-Health/citadel/clients/oms"
	patientServiceClient "github.com/Orange-Health/citadel/clients/patient_service"
)

type InvestigationResultService struct {
	InvResDao            dao.DataLayer
	Cache                cache.CacheLayer
	Sentry               sentry.SentryLayer
	AbnormalityService   abnormalityService.AbnormalityServiceInterface
	CalculationsService  calculationsService.CalculationsServiceInterface
	RerunService         rerunService.RerunServiceInterface
	RemarkService        remarksService.RemarkServiceInterface
	CdsService           cdsService.CdsServiceInterface
	CdsClient            cdsClient.CdsClientInterface
	OmsClient            omsClient.OmsClientInterface
	PatientServiceClient patientServiceClient.PatientServiceClientInterface
}

func InitializeInvestigationResultService() InvestigationResultServiceInterface {
	return &InvestigationResultService{
		InvResDao:            dao.InitializeInvestigationResultDao(),
		Cache:                cache.InitializeCache(),
		Sentry:               sentry.InitializeSentry(),
		AbnormalityService:   abnormalityService.InitializeAbnormalityService(),
		CalculationsService:  calculationsService.InitializeCalculationsService(),
		RerunService:         rerunService.InitializeRerunService(),
		RemarkService:        remarksService.InitializeRemarkService(),
		CdsService:           cdsService.InitializeCdsService(),
		CdsClient:            cdsClient.InitializeCdsClient(),
		OmsClient:            omsClient.InitializeOmsClient(),
		PatientServiceClient: patientServiceClient.InitializePatientServiceClient(),
	}
}
