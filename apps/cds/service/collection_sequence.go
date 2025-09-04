package service

import (
	"context"
	"fmt"
	"net/http"

	commonConstants "github.com/Orange-Health/citadel/common/constants"
	commonStructures "github.com/Orange-Health/citadel/common/structures"
	commonUtils "github.com/Orange-Health/citadel/common/utils"
)

func getCacheKeyForCollectionSequence(request commonStructures.CollectionSequenceRequest) string {
	masterTestIdsString := commonUtils.ConvertUintSliceToString(request.MasterTestIds)
	return fmt.Sprintf(commonConstants.CacheKeyCollectionSequence, request.CityCode, masterTestIdsString)
}

func (cdsService *CdsService) GetCollectionSequence(ctx context.Context, request commonStructures.CollectionSequenceRequest) (
	commonStructures.CollectionSequenceResponse, *commonStructures.CommonError) {
	// collectionSequenceResponse := commonStructures.CollectionSequenceResponse{}
	// cacheKey := getCacheKeyForCollectionSequence(request)
	// err := cdsService.Cache.Get(ctx, cacheKey, &collectionSequenceResponse)
	// if err == nil && len(collectionSequenceResponse.Collections) > 0 {
	// 	return collectionSequenceResponse, nil
	// }
	collectionSequenceResponse, err := cdsService.CdsClient.GetCollectionSequences(ctx, request)
	if err != nil {
		return commonStructures.CollectionSequenceResponse{}, &commonStructures.CommonError{
			Message:    commonConstants.ERROR_IN_GETTING_COLLECTION_SEQUENCE,
			StatusCode: http.StatusInternalServerError,
		}
	}
	// _ = cdsService.Cache.Set(ctx, cacheKey, collectionSequenceResponse, commonConstants.CacheExpiry15MinutesInt)
	return collectionSequenceResponse, nil
}
