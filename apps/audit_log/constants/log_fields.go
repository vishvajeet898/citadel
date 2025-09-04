package constants

// Audit Log Fields contains the fields for the sample audit log
const (
	LogFieldVisitID          = "visit_id"
	LogFieldBarcode          = "barcode"
	LogFieldStatus           = "status"
	LogFieldLabID            = "lab_id"
	LogFieldDestinationLabID = "destination_lab_id"
	LogFieldRejectionReason  = "rejection_reason"
	LogFieldParentSampleID   = "parent_sample_id"
	LogFieldDeletedAt        = "deleted_at"
)

// Audit Log Fields contains the fields for the sample metadata audit log
const (
	LogFieldNotReceivedReason  = "not_received_reason"
	LogFieldCollectLaterReason = "collect_later_reason"
	LogFieldRejectingLab       = "rejecting_lab"
)
