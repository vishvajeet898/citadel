package mapper

import (
	"time"

	commonConstants "github.com/Orange-Health/citadel/common/constants"
	commonStructures "github.com/Orange-Health/citadel/common/structures"
	commonUtils "github.com/Orange-Health/citadel/common/utils"
)

func getTestStatusUint(testStatus, omsTestStatus string) uint {
	switch omsTestStatus {
	case commonConstants.TEST_STATUS_COMPLETED_NOT_SENT:
		testStatus = commonConstants.TEST_STATUS_COMPLETED_NOT_SENT
	case commonConstants.TEST_STATUS_COMPLETED_SENT:
		testStatus = commonConstants.TEST_STATUS_COMPLETED_SENT
	}
	switch testStatus {
	case commonConstants.TEST_STATUS_REQUESTED:
		return commonConstants.TestStatusRequestedUint
	case commonConstants.TEST_STATUS_RESULT_PENDING,
		commonConstants.TEST_STATUS_RERUN_REQUESTED:
		return commonConstants.TestStatusRequestedUint
	case commonConstants.TEST_STATUS_COMPLETED_NOT_SENT,
		commonConstants.TEST_STATUS_RESULT_SAVED,
		commonConstants.TEST_STATUS_RERUN_RESULT_SAVED,
		commonConstants.TEST_STATUS_WITHHELD,
		commonConstants.TEST_STATUS_CO_AUTHORIZE,
		commonConstants.TEST_STATUS_APPROVE:
		return commonConstants.TestStatusCompletedNotSentUint
	case commonConstants.TEST_STATUS_COMPLETED_SENT:
		return commonConstants.TestStatusCompletedSentUint
	case commonConstants.TEST_STATUS_IN_TRANSFER:
		return commonConstants.TestStatusInTransferUint
	case commonConstants.TEST_STATUS_TRANSFER_FAILED:
		return commonConstants.TestStatusInTransferUint
	case commonConstants.TEST_STATUS_SAMPLE_NOT_RECEIVED:
		return commonConstants.TestStatusStatusSampleNotReceivedUint
	case commonConstants.TEST_STATUS_COLLECT_SAMPLE_LATER:
		return commonConstants.TestStatusStatusCollectSampleLaterUint
	default:
		return 0
	}
}

func getPatientAge(dob *time.Time, expectedDob *time.Time) uint {

	if dob != nil {
		ageYears, _, _ := commonUtils.GetAgeYearsMonthsAndDaysFromDob(*dob)
		return uint(ageYears)
	}
	if expectedDob != nil {
		ageYears, _, _ := commonUtils.GetAgeYearsMonthsAndDaysFromDob(*expectedDob)
		return uint(ageYears)
	}

	return 0
}

func MapBulkEtsTestDbEventToEtsTestDbEventAndFilterOutsourcedTests(tatBreachedDetails []commonStructures.EtsTestDbEvent) []commonStructures.EtsTestEvent {
	etsTestEvents := []commonStructures.EtsTestEvent{}
	for _, testEvent := range tatBreachedDetails {
		status := getTestStatusUint(testEvent.TestStatus, testEvent.OmsTestStatus)
		if status == 0 {
			continue
		}
		etsTestEvent := commonStructures.EtsTestEvent{
			TestID:        testEvent.TestID,
			OrderID:       testEvent.OrderID,
			TestStatus:    status,
			TestName:      testEvent.TestName,
			TestDeletedAt: testEvent.TestDeletedAt,
			MasterTestID:  testEvent.MasterTestID,
			LabEta:        testEvent.LabEta,
			LabID:         testEvent.LabID,
			PatientName:   testEvent.PatientName,
			PatientAge:    getPatientAge(testEvent.PatientDob, testEvent.PatientExpectedDob),
			PatientGender: testEvent.PatientGender,
			Barcode:       testEvent.Barcode,
			VialTypeID:    testEvent.VialTypeID,
			IsRejected:    testEvent.IsRejected,
			CityCode:      testEvent.CityCode,
			LisStatus:     testEvent.LisStatus,
		}
		etsTestEvents = append(etsTestEvents, etsTestEvent)
	}
	return etsTestEvents
}
