package utils

import (
	"github.com/Orange-Health/citadel/common/constants"
	"github.com/Orange-Health/citadel/common/structures"
)

func IsTestInhouse(testProcessingLabId, servicingLabId uint, labIdLabMap map[uint]structures.Lab) bool {
	if testProcessingLabId == servicingLabId || labIdLabMap[testProcessingLabId].Inhouse {
		return true
	}
	return false
}

func GetCovid19MasterTestIds() []uint {
	return ConvertIntSliceToUintSlice(constants.Covid19MasterTestIds)
}
