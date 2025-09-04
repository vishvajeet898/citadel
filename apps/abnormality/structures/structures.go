package structures

import (
	"strings"

	commonConstants "github.com/Orange-Health/citadel/common/constants"
)

type ReferenceRange struct {
	Value    string
	Nmin     string
	Nmax     string
	Cmin     string
	Cmax     string
	Imin     string
	Imax     string
	Amin     string
	Amax     string
	Nvalue   string
	Cvalue   string
	Ivalue   string
	Avalue   string
	RefLabel string
}

func (rr *ReferenceRange) CompareTextualValue() string {
	if strings.Contains(rr.Nvalue, "|") {
		values := strings.Split(rr.Nvalue, "|")
		for _, value := range values {
			if strings.EqualFold(strings.TrimSpace(value), rr.Value) {
				return commonConstants.ABNORMALITY_NORMAL
			}
		}
	} else {
		if strings.EqualFold(rr.Nvalue, rr.Value) {
			return commonConstants.ABNORMALITY_NORMAL
		}
	}
	return commonConstants.ABNORMALITY_UPPER_ABNORMAL
}

func (rr *ReferenceRange) CompareTextualValueAutoApproval() bool {
	if strings.Contains(rr.Avalue, "|") {
		values := strings.Split(rr.Avalue, "|")
		for _, value := range values {
			if strings.EqualFold(strings.TrimSpace(value), rr.Value) {
				return true
			}
		}
	} else {
		if strings.EqualFold(rr.Avalue, rr.Value) {
			return true
		}
	}
	return false
}
