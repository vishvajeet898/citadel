package service

import (
	"context"
	"fmt"
	"net/http"

	commonConstants "github.com/Orange-Health/citadel/common/constants"
	commonStructures "github.com/Orange-Health/citadel/common/structures"
)

func (cdsService *CdsService) GetMasterVialType(ctx context.Context, vialTypeId uint) (commonStructures.MasterVialType,
	*commonStructures.CommonError) {
	vial := commonStructures.MasterVialType{}
	cacheKey := fmt.Sprintf(commonConstants.CacheKeyVialsId, vialTypeId)
	err := cdsService.Cache.Get(ctx, cacheKey, &vial)
	if err == nil && vial.Id > 0 {
		return vial, nil
	}
	vialMap := cdsService.GetMasterVialTypeAsMap(ctx)
	vial, ok := vialMap[vialTypeId]
	if ok {
		_ = cdsService.Cache.Set(ctx, cacheKey, vial, commonConstants.CacheExpiry15MinutesInt)
		return vial, nil
	}
	return commonStructures.MasterVialType{}, &commonStructures.CommonError{
		Message:    commonConstants.ERROR_VIAL_NOT_FOUND,
		StatusCode: http.StatusNotFound,
	}
}

func (cdsService *CdsService) GetMasterVialTypeAsMap(ctx context.Context) map[uint]commonStructures.MasterVialType {
	vialMap := map[uint]commonStructures.MasterVialType{}
	cacheKey := commonConstants.CacheKeyVialsMap
	err := cdsService.Cache.Get(ctx, cacheKey, &vialMap)
	if err == nil && len(vialMap) > 0 {
		return vialMap
	}
	masterVialTypes := cdsService.GetMasterVialTypes(ctx)
	for _, vial := range masterVialTypes {
		vialMap[vial.Id] = vial
	}
	_ = cdsService.Cache.Set(ctx, cacheKey, vialMap, commonConstants.CacheExpiry15MinutesInt)
	return vialMap
}

func (cdsService *CdsService) GetMasterVialTypes(ctx context.Context) []commonStructures.MasterVialType {
	vials := []commonStructures.MasterVialType{}
	cacheKey := commonConstants.CacheKeyVialsAll
	err := cdsService.Cache.Get(ctx, cacheKey, &vials)
	if err == nil && len(vials) > 0 {
		return vials
	}
	vials, err = cdsService.CdsClient.GetMasterVialTypes(ctx)
	if err != nil {
		return []commonStructures.MasterVialType{}
	}
	_ = cdsService.Cache.Set(ctx, cacheKey, vials, commonConstants.CacheExpiry15MinutesInt)
	return vials
}
