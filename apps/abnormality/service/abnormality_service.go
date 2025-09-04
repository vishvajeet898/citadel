package abnormalityService

import (
	"context"

	"github.com/Orange-Health/citadel/apps/abnormality/structures"
	commonConstants "github.com/Orange-Health/citadel/common/constants"
	commonStructs "github.com/Orange-Health/citadel/common/structures"
	commonUtils "github.com/Orange-Health/citadel/common/utils"
)

type AbnormalityServiceInterface interface {
	GetInvestigationAbnormality(ctx context.Context, investigationValue string,
		investigation commonStructs.Investigation) string
	GetInvestigationAbnormalityByMasterReferenceRange(ctx context.Context, investigationValue string,
		masterReferenceRange commonStructs.MasterReferenceRange) string
	GetInvestigationAutoApprovalStatus(ctx context.Context, investigationValue string,
		investigation commonStructs.Investigation) bool
}

func (abnormalityService *AbnormalityService) GetInvestigationAbnormality(ctx context.Context, investigationValue string,
	investigation commonStructs.Investigation) string {

	referenceRangeStruct := getReferenceRangeStruct(investigationValue, investigation)

	if referenceRangeStruct.RefLabel == "" {
		return ""
	}

	return getInvestigationAbnormalityByReferenceRange(ctx, *referenceRangeStruct)
}

func (abnormalityService *AbnormalityService) GetInvestigationAbnormalityByMasterReferenceRange(ctx context.Context,
	investigationValue string, masterReferenceRange commonStructs.MasterReferenceRange) string {

	referenceRangeStruct := getReferenceRangeStructByMasterReferenceRange(investigationValue, masterReferenceRange)

	if referenceRangeStruct.RefLabel == "" {
		return ""
	}

	return getInvestigationAbnormalityByReferenceRange(ctx, *referenceRangeStruct)
}

func (abnormalityService *AbnormalityService) GetInvestigationAutoApprovalStatus(ctx context.Context,
	investigationValue string, investigation commonStructs.Investigation) bool {

	if investigationValue == "" {
		return false
	}

	referenceRangeStruct := getReferenceRangeStruct(investigationValue, investigation)

	value := commonUtils.ConvertStringToFloat32ForAbnormality(referenceRangeStruct.Value)
	amin := referenceRangeStruct.Amin
	amax := referenceRangeStruct.Amax

	switch referenceRangeStruct.RefLabel {
	case commonConstants.ReferenceRangeLabelRange:
		if amin != "" && amax != "" {
			return value >= commonUtils.ConvertStringToFloat32ForAbnormality(amin) &&
				value <= commonUtils.ConvertStringToFloat32ForAbnormality(amax)
		} else if amin != "" && amax == "" {
			return value >= commonUtils.ConvertStringToFloat32ForAbnormality(amin)
		} else if amin == "" && amax != "" {
			return value <= commonUtils.ConvertStringToFloat32ForAbnormality(amax)
		}
	case commonConstants.ReferenceRangeLabelTextual:
		return referenceRangeStruct.CompareTextualValueAutoApproval()
	}

	return false
}

func getInvestigationAbnormalityByReferenceRange(ctx context.Context, refRange structures.ReferenceRange) string {
	return getAbnormality(ctx, refRange)
}

func getAbnormality(ctx context.Context,
	referenceRangeStruct structures.ReferenceRange) string {

	switch referenceRangeStruct.RefLabel {
	case commonConstants.ReferenceRangeLabelRange:
		return getRangeAbnormality(referenceRangeStruct)
	case commonConstants.ReferenceRangeLabelTextual:
		return getTextualAbnormality(referenceRangeStruct)
	default:
		commonUtils.AddLog(ctx, commonConstants.ERROR_LEVEL, commonConstants.INVALID_REFERENCE_TYPE, map[string]interface{}{
			"ReferenceRange": referenceRangeStruct,
		}, nil)
		return ""
	}
}

func getRangeAbnormality(referenceRangeStruct structures.ReferenceRange) string {
	if referenceRangeStruct.Value != "" && referenceRangeStruct.Value[0] == '-' {
		return commonConstants.ABNORMALITY_LOWER_ABNORMAL
	}
	if findAbnormality(referenceRangeStruct.Imin, referenceRangeStruct.Imax, referenceRangeStruct.Value, true, true, true) {
		return commonConstants.ABNORMALITY_IMPROBABLE
	} else if (referenceRangeStruct.Cmin != "" && findAbnormality(referenceRangeStruct.Imin, referenceRangeStruct.Cmin, referenceRangeStruct.Value, false, true, false)) ||
		(referenceRangeStruct.Cmax != "" && findAbnormality(referenceRangeStruct.Cmax, referenceRangeStruct.Imax, referenceRangeStruct.Value, true, false, false)) {
		return commonConstants.ABNORMALITY_CRITICAL
	} else if findAbnormality(referenceRangeStruct.Cmin, referenceRangeStruct.Nmin, referenceRangeStruct.Value, false, false, false) {
		return commonConstants.ABNORMALITY_LOWER_ABNORMAL
	} else if findAbnormality(referenceRangeStruct.Nmax, referenceRangeStruct.Cmax, referenceRangeStruct.Value, false, false, false) {
		return commonConstants.ABNORMALITY_UPPER_ABNORMAL
	} else if findAbnormality(referenceRangeStruct.Nmin, referenceRangeStruct.Nmax, referenceRangeStruct.Value, true, true, false) {
		return commonConstants.ABNORMALITY_NORMAL
	} else {
		return ""
	}
}

func getTextualAbnormality(referenceRangeStruct structures.ReferenceRange) string {
	return referenceRangeStruct.CompareTextualValue()
}

func findAbnormality(min, max, valueStr string, minInclusive, maxInclusive, checkOuterRange bool) bool {
	value := commonUtils.ConvertStringToFloat32ForAbnormality(valueStr)
	if checkOuterRange {
		if min != "" && max != "" {
			return value <= commonUtils.ConvertStringToFloat32ForAbnormality(min) ||
				value >= commonUtils.ConvertStringToFloat32ForAbnormality(max)
		} else if min != "" && max == "" {
			return value <= commonUtils.ConvertStringToFloat32ForAbnormality(min)
		} else if min == "" && max != "" {
			return value >= commonUtils.ConvertStringToFloat32ForAbnormality(max)
		}
	} else {
		if minInclusive && maxInclusive {
			if min != "" && max != "" {
				return value >= commonUtils.ConvertStringToFloat32ForAbnormality(min) &&
					value <= commonUtils.ConvertStringToFloat32ForAbnormality(max)
			} else if min != "" && max == "" {
				return value >= commonUtils.ConvertStringToFloat32ForAbnormality(min)
			} else if min == "" && max != "" {
				return value <= commonUtils.ConvertStringToFloat32ForAbnormality(max)
			}
		} else if minInclusive && !maxInclusive {
			if min != "" && max != "" {
				return value >= commonUtils.ConvertStringToFloat32ForAbnormality(min) &&
					value < commonUtils.ConvertStringToFloat32ForAbnormality(max)
			} else if min != "" && max == "" {
				return value >= commonUtils.ConvertStringToFloat32ForAbnormality(min)
			} else if min == "" && max != "" {
				return value < commonUtils.ConvertStringToFloat32ForAbnormality(max)
			}
		} else if !minInclusive && maxInclusive {
			if min != "" && max != "" {
				return value > commonUtils.ConvertStringToFloat32ForAbnormality(min) &&
					value <= commonUtils.ConvertStringToFloat32ForAbnormality(max)
			} else if min != "" && max == "" {
				return value > commonUtils.ConvertStringToFloat32ForAbnormality(min)
			} else if min == "" && max != "" {
				return value <= commonUtils.ConvertStringToFloat32ForAbnormality(max)
			}
		} else if !minInclusive && !maxInclusive {
			if min != "" && max != "" {
				return value > commonUtils.ConvertStringToFloat32ForAbnormality(min) &&
					value < commonUtils.ConvertStringToFloat32ForAbnormality(max)
			} else if min != "" && max == "" {
				return value > commonUtils.ConvertStringToFloat32ForAbnormality(min)
			} else if min == "" && max != "" {
				return value < commonUtils.ConvertStringToFloat32ForAbnormality(max)
			}
		}
	}

	return false
}

func getReferenceRangeStruct(testValue string, investigation commonStructs.Investigation) *structures.ReferenceRange {
	return &structures.ReferenceRange{
		Value: testValue,

		Nmin:   investigation.ReferenceRange.NormalRange.MinValue,
		Nmax:   investigation.ReferenceRange.NormalRange.MaxValue,
		Nvalue: investigation.ReferenceRange.NormalRange.ResultValue,

		Cmin:   investigation.ReferenceRange.CriticalRange.MinValue,
		Cmax:   investigation.ReferenceRange.CriticalRange.MaxValue,
		Cvalue: investigation.ReferenceRange.CriticalRange.ResultValue,

		Imin:   investigation.ReferenceRange.ImprobableRange.MinValue,
		Imax:   investigation.ReferenceRange.ImprobableRange.MaxValue,
		Ivalue: investigation.ReferenceRange.ImprobableRange.ResultValue,

		Amin:   investigation.ReferenceRange.AutoApprovalRange.MinValue,
		Amax:   investigation.ReferenceRange.AutoApprovalRange.MaxValue,
		Avalue: investigation.ReferenceRange.AutoApprovalRange.ResultValue,

		RefLabel: investigation.ReferenceRange.ReferenceLabel,
	}
}

func getReferenceRangeStructByMasterReferenceRange(testValue string,
	masterReferenceRange commonStructs.MasterReferenceRange) *structures.ReferenceRange {
	return &structures.ReferenceRange{
		Value: testValue,

		Nmin:   masterReferenceRange.NormalRange.MinValue,
		Nmax:   masterReferenceRange.NormalRange.MaxValue,
		Nvalue: masterReferenceRange.NormalRange.ResultValue,

		Cmin:   masterReferenceRange.CriticalRange.MinValue,
		Cmax:   masterReferenceRange.CriticalRange.MaxValue,
		Cvalue: masterReferenceRange.CriticalRange.ResultValue,

		Imin:   masterReferenceRange.ImprobableRange.MinValue,
		Imax:   masterReferenceRange.ImprobableRange.MaxValue,
		Ivalue: masterReferenceRange.ImprobableRange.ResultValue,

		Amin:   masterReferenceRange.AutoApprovalRange.MinValue,
		Amax:   masterReferenceRange.AutoApprovalRange.MaxValue,
		Avalue: masterReferenceRange.AutoApprovalRange.ResultValue,

		RefLabel: masterReferenceRange.ReferenceLabel,
	}
}
