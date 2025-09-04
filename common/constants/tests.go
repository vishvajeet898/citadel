package constants

// Test Details Statuses
const (
	TEST_STATUS_REQUESTED            = "requested"
	TEST_STATUS_RESULT_PENDING       = "result_pending"
	TEST_STATUS_RESULT_SAVED         = "result_saved"
	TEST_STATUS_REJECTED             = "rejected"
	TEST_STATUS_RERUN_REQUESTED      = "rerun_requested"
	TEST_STATUS_RERUN_RESULT_SAVED   = "rerun_result_saved"
	TEST_STATUS_WITHHELD             = "withheld"
	TEST_STATUS_CO_AUTHORIZE         = "co_authorize"
	TEST_STATUS_APPROVE              = "approve"
	TEST_STATUS_COMPLETED_NOT_SENT   = "completed_not_sent"
	TEST_STATUS_COMPLETED_SENT       = "completed_sent"
	TEST_STATUS_IN_TRANSFER          = "in_transfer"
	TEST_STATUS_TRANSFER_FAILED      = "transfer_failed"
	TEST_STATUS_SAMPLE_NOT_RECEIVED  = "sample_not_received"
	TEST_STATUS_COLLECT_SAMPLE_LATER = "collect_sample_later"

	// Derived Statuses
	TEST_STATUS_LAB_RECEIVED = "lab_received"
)

// Test Details Statuses Label Maps
var TEST_STATUSES_LABEL_MAP = map[string]string{
	TEST_STATUS_REQUESTED:            "Requested",
	TEST_STATUS_RESULT_PENDING:       "Result Pending",
	TEST_STATUS_RESULT_SAVED:         "Result Saved",
	TEST_STATUS_REJECTED:             "Rejected",
	TEST_STATUS_RERUN_REQUESTED:      "Rerun",
	TEST_STATUS_RERUN_RESULT_SAVED:   "Result Saved",
	TEST_STATUS_WITHHELD:             "Withheld",
	TEST_STATUS_CO_AUTHORIZE:         "Co-Authorize",
	TEST_STATUS_APPROVE:              "Approve",
	TEST_STATUS_COMPLETED_NOT_SENT:   "Completed Not Sent",
	TEST_STATUS_COMPLETED_SENT:       "Completed Sent",
	TEST_STATUS_IN_TRANSFER:          "In Transfer",
	TEST_STATUS_TRANSFER_FAILED:      "Transfer Failed",
	TEST_STATUS_SAMPLE_NOT_RECEIVED:  "Sample Not Received",
	TEST_STATUS_COLLECT_SAMPLE_LATER: "Collect Sample Later",
	TEST_STATUS_LAB_RECEIVED:         "Lab Received",
}

// Test Details Statuses
var TEST_STATUSES = []string{
	TEST_STATUS_REQUESTED,
	TEST_STATUS_RESULT_PENDING,
	TEST_STATUS_RESULT_SAVED,
	TEST_STATUS_REJECTED,
	TEST_STATUS_RERUN_REQUESTED,
	TEST_STATUS_RERUN_RESULT_SAVED,
	TEST_STATUS_WITHHELD,
	TEST_STATUS_CO_AUTHORIZE,
	TEST_STATUS_APPROVE,
	TEST_STATUS_COMPLETED_NOT_SENT,
	TEST_STATUS_COMPLETED_SENT,
	TEST_STATUS_IN_TRANSFER,
	TEST_STATUS_TRANSFER_FAILED,
	TEST_STATUS_SAMPLE_NOT_RECEIVED,
	TEST_STATUS_COLLECT_SAMPLE_LATER,
	TEST_STATUS_LAB_RECEIVED,
}

// Attune Test Statuses
const (
	ATTUNE_TEST_STATUS_COMPLETED = "Completed"
	ATTUNE_TEST_STATUS_APPROVED  = "Approve"
	ATTUNE_TEST_STATUS_RECHECK   = "Recheck"
	ATTUNE_TEST_STATUS_RERUN     = "Rerun"
)

// Approval Sources
const (
	TEST_APPROVAL_SOURCE_ATTUNE = "ATTUNE"
	TEST_APPROVAL_SOURCE_OH     = "OH"
)

// Test Details Statuses as uint
const (
	TestStatusRequestedUint                uint = 1
	TestStatusOrderedUint                  uint = 2
	TestStatusCompletedNotSentUint         uint = 3
	TestStatusCancelledUint                uint = 5
	TestStatusCompletedSentUint            uint = 6
	TestStatusInTransferUint               uint = 7
	TestStatusTransferFailedUint           uint = 8
	TestStatusStatusSampleNotReceivedUint  uint = 9
	TestStatusStatusCollectSampleLaterUint uint = 10
	TestStatusTechnicianSavedUint          uint = 11
	TestStatusLabReceivedUint              uint = 12
	TestStatusRerunUint                    uint = 13
)

var TestStatusUintToStringMap = map[uint]string{
	TestStatusRequestedUint:                TEST_STATUS_REQUESTED,
	TestStatusOrderedUint:                  TEST_STATUS_REQUESTED,
	TestStatusCompletedNotSentUint:         TEST_STATUS_COMPLETED_NOT_SENT,
	TestStatusCancelledUint:                TEST_STATUS_REJECTED,
	TestStatusCompletedSentUint:            TEST_STATUS_COMPLETED_SENT,
	TestStatusInTransferUint:               TEST_STATUS_IN_TRANSFER,
	TestStatusTransferFailedUint:           TEST_STATUS_TRANSFER_FAILED,
	TestStatusStatusSampleNotReceivedUint:  TEST_STATUS_SAMPLE_NOT_RECEIVED,
	TestStatusStatusCollectSampleLaterUint: TEST_STATUS_COLLECT_SAMPLE_LATER,
}

var OmsTestStatusToUintMap = map[string]uint{
	TEST_STATUS_REQUESTED:            TestStatusRequestedUint,
	TEST_STATUS_RESULT_PENDING:       TestStatusRequestedUint,
	TEST_STATUS_RESULT_SAVED:         TestStatusRequestedUint,
	TEST_STATUS_RERUN_REQUESTED:      TestStatusRequestedUint,
	TEST_STATUS_RERUN_RESULT_SAVED:   TestStatusRequestedUint,
	TEST_STATUS_WITHHELD:             TestStatusRequestedUint,
	TEST_STATUS_CO_AUTHORIZE:         TestStatusRequestedUint,
	TEST_STATUS_APPROVE:              TestStatusCompletedNotSentUint,
	TEST_STATUS_COMPLETED_NOT_SENT:   TestStatusCompletedNotSentUint,
	TEST_STATUS_REJECTED:             TestStatusCancelledUint,
	TEST_STATUS_COMPLETED_SENT:       TestStatusCompletedSentUint,
	TEST_STATUS_IN_TRANSFER:          TestStatusInTransferUint,
	TEST_STATUS_TRANSFER_FAILED:      TestStatusTransferFailedUint,
	TEST_STATUS_SAMPLE_NOT_RECEIVED:  TestStatusStatusSampleNotReceivedUint,
	TEST_STATUS_COLLECT_SAMPLE_LATER: TestStatusStatusCollectSampleLaterUint,
}

var TEST_STATUSES_WITH_PROCESSING_LAB = []string{
	TEST_STATUS_RESULT_SAVED,
	TEST_STATUS_RERUN_RESULT_SAVED,
	TEST_STATUS_WITHHELD,
	TEST_STATUS_CO_AUTHORIZE,
}

var OmsCompletedTestStatusesUint = []uint{
	TestStatusCompletedNotSentUint,
	TestStatusCompletedSentUint,
}

var OmsCompletedTestStatuses = []string{
	TEST_STATUS_COMPLETED_NOT_SENT,
	TEST_STATUS_COMPLETED_SENT,
}

const (
	TEST_REPORT_STATUS_NOT_READY             = "not_ready"
	TEST_REPORT_STATUS_QUEUED                = "queued"
	TEST_REPORT_STATUS_CREATED               = "created"
	TEST_REPORT_STATUS_CREATION_FAILED       = "creation_failed"
	TEST_REPORT_STATUS_SYNCED                = "synced"
	TEST_REPORT_STATUS_SYNC_FAILED           = "sync_failed"
	TEST_REPORT_STATUS_SENT                  = "sent"
	TEST_REPORT_STATUS_DELIVERED             = "delivered"
	TEST_REPORT_STATUS_PRELIMINARY_QUEUED    = "preliminary_queued"
	TEST_REPORT_STATUS_PRELIMINARY_DELIVERED = "preliminary_delivered"
)

var (
	TEST_REPORT_STATUSES_NOT_NEW = []string{
		TEST_REPORT_STATUS_SENT,
		TEST_REPORT_STATUS_DELIVERED,
	}
)

var (
	DutyDoctorMap        = Config.GetStringMap("duty_doctor")
	Covid19MasterTestIds = Config.GetIntSlice("master_tests.covid_rtpcr_ids")
)
