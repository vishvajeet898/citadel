package utils

import (
	"fmt"
	"time"

	"github.com/Orange-Health/citadel/common/constants"
)

func GetEnteredAtTime(resultCapturedAt string) *time.Time {
	var enteredAt time.Time
	if resultCapturedAt != "" {
		if enteredAtTime, err := time.Parse(constants.DateTimeUTCLayoutWithoutTZOffset, resultCapturedAt); err == nil {
			enteredAt = enteredAtTime
		}
	}
	return &enteredAt
}

func GetApprovedAtTime(resultApprovedAt string) *time.Time {
	var approvedAt time.Time
	if resultApprovedAt != "" {
		if approvedAtTime, err := time.Parse(constants.DateTimeUTCLayoutWithoutTZOffset, resultApprovedAt); err == nil {
			approvedAt = approvedAtTime
		}
	}
	return &approvedAt
}

func GetAttuneOrgCodeByLabId(labId uint) string {
	orgCodes := constants.AttuneOrgCodes
	return orgCodes[fmt.Sprint(labId)]
}

func IsQcEnabledForLab(labId uint) bool {
	return SliceContainsInt(constants.QcEnabledLabIds, int(labId))
}
