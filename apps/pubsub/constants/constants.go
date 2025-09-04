package constants

var PubSubEventContainsMap = map[string][]string{
	OrderReportUpdateEvent:       orderReportUpdateEventContains,
	ContactMergeConfirmEvent:     contactMergeConfirmEventContains,
	ReportReadyEvent:             reportReadyEventContains,
	AddSampleRejectedTagEvent:    addSampleRejectedTagEventContains,
	RemoveSampleRejectedTagEvent: removeSampleRejectedTagEventContains,
	ResetTestTatsEvent:           resetTestTatsEventContains,
	SampleCollectedEvent:         sampleCollectedEventContains,
	UpdateTestStatusEvent:        updateTestStatusEventContains,
	LabEtaUpdateEvent:            labEtaUpdateEventContains,
	CheckOrderCompletionEvent:    checkOrderCompletionEventContains,
	CitadelLisEvent:              citadelLisEventContains,
}

const (
	OrderReportUpdateEvent       = "order.report_updated"
	ContactMergeConfirmEvent     = "contact.merge_confirm"
	ReportReadyEvent             = "order.report_ready"
	AddSampleRejectedTagEvent    = "order.add_sample_rejected_tag"
	RemoveSampleRejectedTagEvent = "order.remove_sample_rejected_tag"
	ResetTestTatsEvent           = "order.reset_test_tats"
	SampleCollectedEvent         = "order.sample_collected"
	UpdateTestStatusEvent        = "order.update_test_status"
	LabEtaUpdateEvent            = "order.lab_eta_update"
	CheckOrderCompletionEvent    = "order.check_order_completion"
	CitadelLisEvent              = "lisvisit.update"
	EtsTestEvent                 = "ets.test"
)

var (
	orderReportUpdateEventContains       = []string{"order", "visits"}
	contactMergeConfirmEventContains     = []string{"merge_contact", "master_contact", "service", "was_merge_contact_found"}
	reportReadyEventContains             = []string{"report_pdf"}
	addSampleRejectedTagEventContains    = []string{"order_id"}
	removeSampleRejectedTagEventContains = []string{"order_id"}
	resetTestTatsEventContains           = []string{"test_ids"}
	sampleCollectedEventContains         = []string{"test_ids", "collected_at", "city_code"}
	updateTestStatusEventContains        = []string{"test_ids", "sample_status"}
	labEtaUpdateEventContains            = []string{"test_ids", "lis_sync_at"}
	checkOrderCompletionEventContains    = []string{"order_id"}
	citadelLisEventContains              = []string{"lis_data", "visit_id"}
)
