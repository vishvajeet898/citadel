package service

import (
	"context"
	"fmt"
	"net/http"

	commonConstants "github.com/Orange-Health/citadel/common/constants"
	commonStructures "github.com/Orange-Health/citadel/common/structures"
)

func (cdsService *CdsService) GetAllLabs(ctx context.Context) []commonStructures.Lab {
	labs := []commonStructures.Lab{}
	cacheKey := commonConstants.CacheKeyLabsAll
	err := cdsService.Cache.Get(ctx, cacheKey, &labs)
	if err == nil && labs != nil {
		return labs
	}
	labs, cErr := cdsService.CdsClient.GetLabMasters(ctx)
	if cErr != nil {
		return []commonStructures.Lab{}
	}
	_ = cdsService.Cache.Set(ctx, cacheKey, labs, commonConstants.CacheExpiry15MinutesInt)
	return labs
}

func (cdsService *CdsService) GetLabIdLabMap(ctx context.Context) map[uint]commonStructures.Lab {
	labIdLabMap := map[uint]commonStructures.Lab{}
	cacheKey := commonConstants.CacheKeyLabsMap
	err := cdsService.Cache.Get(ctx, cacheKey, &labIdLabMap)
	if err == nil && len(labIdLabMap) > 0 {
		return labIdLabMap
	}
	labs := cdsService.GetAllLabs(ctx)
	for _, lab := range labs {
		labIdLabMap[lab.Id] = lab
	}
	_ = cdsService.Cache.Set(ctx, cacheKey, labIdLabMap, commonConstants.CacheExpiry15MinutesInt)
	return labIdLabMap
}

func (cdsService *CdsService) GetLabById(ctx context.Context, labId uint) (commonStructures.Lab,
	*commonStructures.CommonError) {
	lab := commonStructures.Lab{}
	cacheKey := fmt.Sprintf(commonConstants.CacheKeyLabsId, labId)
	err := cdsService.Cache.Get(ctx, cacheKey, &lab)
	if err == nil {
		return lab, nil
	}

	labIdLabMap := cdsService.GetLabIdLabMap(ctx)
	lab = labIdLabMap[labId]
	_ = cdsService.Cache.Set(ctx, cacheKey, lab, commonConstants.CacheExpiry15MinutesInt)
	if lab.Id == 0 {
		return lab, &commonStructures.CommonError{
			Message:    commonConstants.ERROR_LAB_NOT_FOUND,
			StatusCode: http.StatusInternalServerError,
		}
	}
	return lab, nil
}

func (cdsService *CdsService) GetInhouseLabIds(ctx context.Context) []uint {
	inhouseLabIds := []uint{}
	cacheKey := commonConstants.CacheKeyInHouseLabIds
	err := cdsService.Cache.Get(ctx, cacheKey, &inhouseLabIds)
	if err == nil && inhouseLabIds != nil {
		return inhouseLabIds
	}
	labs := cdsService.GetAllLabs(ctx)
	for _, lab := range labs {
		if lab.Inhouse {
			inhouseLabIds = append(inhouseLabIds, lab.Id)
		}
	}
	_ = cdsService.Cache.Set(ctx, cacheKey, inhouseLabIds, commonConstants.CacheExpiry15MinutesInt)
	return inhouseLabIds
}

func (cdsService *CdsService) GetOutsourceLabIds(ctx context.Context) []uint {
	outsourceLabIds := []uint{}
	cacheKey := commonConstants.CacheKeyOutSourceLabIds
	err := cdsService.Cache.Get(ctx, cacheKey, &outsourceLabIds)
	if err == nil && outsourceLabIds != nil {
		return outsourceLabIds
	}
	labs := cdsService.GetAllLabs(ctx)
	for _, lab := range labs {
		if !lab.Inhouse {
			outsourceLabIds = append(outsourceLabIds, lab.Id)
		}
	}
	_ = cdsService.Cache.Set(ctx, cacheKey, outsourceLabIds, commonConstants.CacheExpiry15MinutesInt)
	return outsourceLabIds
}
