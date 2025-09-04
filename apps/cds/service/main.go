package service

import (
	"context"

	"github.com/Orange-Health/citadel/adapters/cache"
	"github.com/Orange-Health/citadel/adapters/sentry"
	"github.com/Orange-Health/citadel/apps/cds/dao"
	cdsClient "github.com/Orange-Health/citadel/clients/cds"
	commonStructures "github.com/Orange-Health/citadel/common/structures"
)

type CdsService struct {
	CdsDao dao.DataLayer
	Cache  cache.CacheLayer
	Sentry sentry.SentryLayer

	CdsClient cdsClient.CdsClientInterface
}

type CdsServiceInterface interface {
	// Vials
	GetMasterVialType(ctx context.Context, vialTypeId uint) (commonStructures.MasterVialType, *commonStructures.CommonError)
	GetMasterVialTypeAsMap(ctx context.Context) map[uint]commonStructures.MasterVialType
	GetMasterVialTypes(ctx context.Context) []commonStructures.MasterVialType

	// Labs
	GetAllLabs(ctx context.Context) []commonStructures.Lab
	GetLabIdLabMap(ctx context.Context) map[uint]commonStructures.Lab
	GetLabById(ctx context.Context, labId uint) (commonStructures.Lab, *commonStructures.CommonError)
	GetInhouseLabIds(ctx context.Context) []uint
	GetOutsourceLabIds(ctx context.Context) []uint

	// Master Tests
	GetMasterTestById(ctx context.Context, id uint) (commonStructures.CdsTestMaster, error)
	GetMasterTestsByIds(ctx context.Context, ids []uint) []commonStructures.CdsTestMaster
	GetMasterTestAsMap(ctx context.Context) map[uint]commonStructures.CdsTestMaster
	GetMasterTests(ctx context.Context) []commonStructures.CdsTestMaster
	GetDeduplicatedTestsAndPackages(ctx context.Context, testIds, packageIds []uint) (
		commonStructures.TestDeduplicationResponse, error)
	GetNrlCpEnabledMasterTestIds(ctx context.Context) []uint

	// Collection Sequence
	GetCollectionSequence(ctx context.Context, request commonStructures.CollectionSequenceRequest) (
		commonStructures.CollectionSequenceResponse, *commonStructures.CommonError)
}

func InitializeCdsService() CdsServiceInterface {
	return &CdsService{
		CdsDao: dao.InitializeCdsDao(),
		Cache:  cache.InitializeCache(),
		Sentry: sentry.InitializeSentry(),

		CdsClient: cdsClient.InitializeCdsClient(),
	}
}
