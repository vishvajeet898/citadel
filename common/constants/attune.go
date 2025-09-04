package constants

import (
	"time"
)

var (
	AttuneAPIBaseUrl            = Config.GetString("attune.api_base_url")
	AttuneIncomingAPIKey        = Config.GetString("attune.incoming_api_key") // Keeping it for now
	AttuneOrgCodes              = Config.GetStringMapString("attune.org_codes")
	QcEnabledLabIds             = Config.GetIntSlice("attune.qc_enabled_lab_ids")
	AttunePatientDefaultContact = "9999999999"
)

const (
	AttuneReportWithoutStationery     = "NonStationeryPDF"       // for reports without branding headers
	AttuneReportWithStationery        = "StationeryPDF"          // for reports with branding headers
	AttuneReportWithCobrandStationery = "StationeryPDFCoBranded" // for reports with co-branding headers
)

var AttuneReportdfTypes = []string{
	AttuneReportWithoutStationery,
	AttuneReportWithStationery,
}

const (
	FileTypeReport           = "report"
	FileTypeReportHeaderless = "report_headerless"
	FileTypeReportCoBranded  = "report_co_branded"
)

const (
	ReportFormatNonBranded          uint = 1
	ReportFormatOrangeHealthBranded uint = 2
	ReportFormatCoBranded           uint = 3
)

// Test Statuses
const (
	AttuneTestStatusCompleted = "Completed"
	AttuneTestStatusApprove   = "Approve"
	AttuneTestStatusCancel    = "Cancel"
	AttuneTestStatusRetest    = "Retest"
	AttuneTestStatusRerun     = "Rerun"
	AttuneTestStatusOrdered   = "Ordered"
)

const (
	AttuneAddressLengthLimit    = 250
	LisMissingValuesMaxRetries  = 3
	LisMissingValuesBackoffTime = 500 * time.Millisecond
	LisSyncMaxRetries           = 3
	LisSyncBackoffTime          = 500 * time.Millisecond
	LisMissingPdfMaxRetries     = 3
	LisMissingPdfBackoffTime    = 1000 * time.Millisecond
)

// Attune Order Details Sync Message Type
const (
	AttuneMessageTypeNew      = "NW"
	AttuneMessageTypeModified = "MO"
	AttuneMessageTypeCancel   = "CA"
)

const (
	AttuneClientName = "GENERAL"
)
