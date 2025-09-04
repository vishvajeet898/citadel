package constants

var (
	LisEvent                   = "lisvisit.update"
	OrderReportPdfReadyEvent   = "order.report_pdf_ready"
	OmsCompletedEvent          = "citadel.lis_completed"
	OmsApprovedEvent           = "citadel.lis_approved"
	OmsAttachmentEvent         = "citadel.attachment"
	OmsTestDeleteEvent         = "citadel.test_deleted"
	OmsManualReportUploadEvent = "citadel.manual_report_upload"
	OmsOrderCreatedEvent       = "order.created"
	OmsOrderUpdatedEvent       = "order.updated"
	OmsOrderCompletedEvent     = "order.completed"
	OmsOrderCancelledEvent     = "order.cancelled"
	MergeContactEvent          = "contact.merge"
	SampleCollectedEvent       = "citadel.sample_collected"
	SampleRecollectionEvent    = "citadel.sample_recollection"
	TestDetailsEvent           = "citadel.test_details"
	MarkSampleSnrEvent         = "citadel.mark_sample_snr"
	UpdateTaskSequence         = "citadel.update_task_sequence"
	UpdateSrfIdToLisEvent      = "citadel.update_srf_id_to_lis"
)

var (
	DoctorTatDuration = Config.GetInt("doctor_tat_duration")
)

var (
	DefaultRetryTimout               = 10
	OmsOrderCreateUpdateRetryTimeout = 3
	OmsCollectionRetryTimeout        = 10
)

var (
	DefaultRetryCount       = 3
	OmsCollectionRetryCount = 5
)

var EventNameGroupIdMap = map[string]string{
	LisEvent:                   "lis_event",
	OrderReportPdfReadyEvent:   "report_updates",
	OmsAttachmentEvent:         "order_updates_%s",
	OmsTestDeleteEvent:         "order_updates_%s",
	OmsManualReportUploadEvent: "order_updates_%s",
	OmsOrderCreatedEvent:       "order_updates_%s",
	OmsOrderUpdatedEvent:       "order_updates_%s",
	OmsOrderCompletedEvent:     "order_updates_%s",
	OmsOrderCancelledEvent:     "order_updates_%s",
	MergeContactEvent:          "merge_contact_event",
	SampleCollectedEvent:       "order_updates_%s",
	SampleRecollectionEvent:    "order_updates_%s",
	TestDetailsEvent:           "order_updates_%s",
	MarkSampleSnrEvent:         "order_updates_%s",
	UpdateTaskSequence:         "order_updates_%s",
	UpdateSrfIdToLisEvent:      "order_updates_%s",
}

var VariableGroupIdEventNames = []string{
	OmsTestDeleteEvent, OmsManualReportUploadEvent, OmsOrderCreatedEvent, OmsOrderUpdatedEvent,
	OmsOrderCompletedEvent, OmsOrderCancelledEvent, SampleCollectedEvent, SampleRecollectionEvent, TestDetailsEvent,
	MarkSampleSnrEvent, UpdateTaskSequence, UpdateSrfIdToLisEvent,
}
