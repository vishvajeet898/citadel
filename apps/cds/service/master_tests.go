package service

import (
	"context"
	"errors"
	"fmt"

	commonConstants "github.com/Orange-Health/citadel/common/constants"
	commonStructures "github.com/Orange-Health/citadel/common/structures"
)

func (cdsService *CdsService) GetMasterTestById(ctx context.Context, id uint) (commonStructures.CdsTestMaster, error) {
	masterTest := commonStructures.CdsTestMaster{}
	cacheKey := fmt.Sprintf(commonConstants.CacheKeyMasterTestId, id)
	err := cdsService.Cache.Get(ctx, cacheKey, &masterTest)
	if err == nil && masterTest.Id != 0 {
		return masterTest, nil
	}
	masterTestMap := cdsService.GetMasterTestAsMap(ctx)
	masterTest, ok := masterTestMap[id]
	if !ok {
		return masterTest, errors.New(commonConstants.ERROR_MASTER_TEST_NOT_FOUND)
	}
	_ = cdsService.Cache.Set(ctx, cacheKey, masterTest, commonConstants.CacheExpiry15MinutesInt)
	return masterTest, err
}

func (cdsService *CdsService) GetMasterTestsByIds(ctx context.Context, ids []uint) []commonStructures.CdsTestMaster {
	masterTests := []commonStructures.CdsTestMaster{}
	masterTestMap := cdsService.GetMasterTestAsMap(ctx)
	for _, id := range ids {
		if masterTest, ok := masterTestMap[id]; ok {
			masterTests = append(masterTests, masterTest)
		}
	}
	return masterTests
}

func (cdsService *CdsService) GetMasterTestAsMap(ctx context.Context) map[uint]commonStructures.CdsTestMaster {
	masterTestMap := map[uint]commonStructures.CdsTestMaster{}
	cacheKey := commonConstants.CacheKeyMasterTestMap
	err := cdsService.Cache.Get(ctx, cacheKey, &masterTestMap)
	if err == nil && len(masterTestMap) > 0 {
		return masterTestMap
	}
	masterTests := cdsService.GetMasterTests(ctx)
	for _, masterTest := range masterTests {
		masterTestMap[masterTest.Id] = masterTest
	}
	_ = cdsService.Cache.Set(ctx, cacheKey, masterTestMap, commonConstants.CacheExpiry15MinutesInt)
	return masterTestMap
}

func (cdsService *CdsService) GetMasterTests(ctx context.Context) []commonStructures.CdsTestMaster {
	masterTests := []commonStructures.CdsTestMaster{}
	cacheKey := commonConstants.CacheKeyMasterTestAll
	err := cdsService.Cache.Get(ctx, cacheKey, &masterTests)
	if err == nil && len(masterTests) > 0 {
		return masterTests
	}
	masterTests, err = cdsService.CdsClient.GetMasterTests(ctx, false, false, false, false, true, true, []uint{}, false)
	if err != nil {
		return masterTests
	}
	_ = cdsService.Cache.Set(ctx, cacheKey, masterTests, commonConstants.CacheExpiry15MinutesInt)
	return masterTests
}

func (cdsService *CdsService) GetDeduplicatedTestsAndPackages(ctx context.Context, testIds, packageIds []uint) (
	commonStructures.TestDeduplicationResponse, error) {
	return cdsService.CdsClient.GetDeduplicatedTestsAndPackages(ctx, testIds, packageIds)
}

func (cdsService *CdsService) GetNrlCpEnabledMasterTestIds(ctx context.Context) []uint {
	nrlCpEnabledMasterTestIds := []uint{}

	cacheKey := commonConstants.CacheKeyNrlEnabledMasterTestIds
	err := cdsService.Cache.Get(ctx, cacheKey, &nrlCpEnabledMasterTestIds)
	if err == nil && len(nrlCpEnabledMasterTestIds) > 0 {
		return nrlCpEnabledMasterTestIds
	}

	nrlCpEnabledMasterTestIds, err = cdsService.CdsClient.GetNrlCpEnabledMasterTestIds(ctx)
	if err != nil {
		return nrlCpEnabledMasterTestIds
	}

	_ = cdsService.Cache.Set(ctx, cacheKey, nrlCpEnabledMasterTestIds, commonConstants.CacheExpiry15MinutesInt)
	return nrlCpEnabledMasterTestIds
}
