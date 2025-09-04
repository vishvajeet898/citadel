package constants

var (
	ERROR_BARCODE_AND_COLLECT_LATER_REASON_EMPTY        = "both barcode and collect later reason can't be empty"
	ERROR_BARCODE_URL_AND_SCANNED_TIME_REQUIRED         = "both barcodeImageURL and barcodeScannedTime required with barcode"
	ERROR_SAMPLE_ALREADY_REJECTED                       = "sample already rejected"
	ERROR_ATLEAST_ONE_TEST_IS_REQUIRED_FOR_RECOLLECTION = "at least one test is required for recollection"
	ERROR_REQUEST_ID_REQUIRED                           = "request id is required"
	ERROR_TASK_ID_REQUIRED                              = "task id is required"
	ERROR_TASK_SEQUENCE_REQUIRED                        = "task sequence is required"
	ERROR_TEST_IDS_REQUIRED                             = "test ids are required"
	ERROR_INVALID_PARAMETERS                            = "invalid parameters provided, please check the request body and try again"
)

var ReasonsToNotTriggerFreshdeskTicket = []string{
	"Duplicate Test (NA)",
	"Sync Error / Barcode Error (NA)",
}

var (
	SampleRejectedFreshDeskSubject          = "Sample Completely Rejected, Action Needed"
	SamplePartiallyRejectedFreshDeskSubject = "Sample Partially Rejected, Action Needed"
)

const (
	CitadelLisEventGroupId = "citadel_lis_event"
)

const (
	ReverseLogisticsNormalTatForDashboard   uint = 210
	ReverseLogisticsCampTatForDashboard     uint = 320
	ReverseLogisticsInclinicTatForDashboard uint = 360
	ReverseLogisticsDelayDaysForDashboard   uint = 7
	ReverseLogisticsNormalTatForSlack       uint = 240
	ReverseLogisticsCampTatForSlack         uint = 360
	ReverseLogisticsInclinicTatForSlack     uint = 390
	ReverseLogisticsDelayDaysForSlack       uint = 2
)
