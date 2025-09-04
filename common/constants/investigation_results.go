package constants

// InvestigationResult Statuses
const (
	INVESTIGATION_STATUS_PENDING      = "pending"
	INVESTIGATION_STATUS_APPROVE      = "approve"
	INVESTIGATION_STATUS_RERUN        = "rerun"
	INVESTIGATION_STATUS_WITHHELD     = "withheld"
	INVESTIGATION_STATUS_CO_AUTHORIZE = "co_authorize"
	INVESTIGATION_STATUS_OMS_APPROVE  = "Approve"
	INVESTIGATION_STATUS_OMS_PENDING  = "Entered"
	INVESTIGATION_STATUS_OMS_RERUN    = "Rerun"
)

var INVESTIGATION_STATUSES_APPROVE = []string{
	INVESTIGATION_STATUS_APPROVE,
	INVESTIGATION_STATUS_OMS_APPROVE,
}

var INVESTIGATION_STATUSES_PENDING = []string{
	INVESTIGATION_STATUS_PENDING,
	INVESTIGATION_STATUS_OMS_PENDING,
}

var INVESTIGATION_STATUSES_RERUN = []string{
	INVESTIGATION_STATUS_RERUN,
	INVESTIGATION_STATUS_OMS_RERUN,
}

const (
	COLUMN_TEST_DETAILS_ID = "test_details_id"
)

var INVESTIGATION_RESULT_STATUSES = []string{
	INVESTIGATION_STATUS_PENDING,
	INVESTIGATION_STATUS_APPROVE,
	INVESTIGATION_STATUS_RERUN,
	INVESTIGATION_STATUS_WITHHELD,
	INVESTIGATION_STATUS_CO_AUTHORIZE,
	INVESTIGATION_STATUS_OMS_APPROVE,
	INVESTIGATION_STATUS_OMS_PENDING,
	INVESTIGATION_STATUS_OMS_RERUN,
}

// Investigation Abnormality Statuses
const (
	ABNORMALITY_NORMAL         = "normal"
	ABNORMALITY_UPPER_ABNORMAL = "upper_abnormal"
	ABNORMALITY_LOWER_ABNORMAL = "lower_abnormal"
	ABNORMALITY_CRITICAL       = "critical"
	ABNORMALITY_IMPROBABLE     = "improbable"

	EMPTY_INVESTIGATION_VALUE = "investigation Value is Empty for %+v"
	INVALID_REFERENCE_TYPE    = "invalid Reference Label : %s"
)

// Correct Department Name Map
var DEPARTMENT_NAME_MAP = map[string]string{
	"haematology":           "Haematology",
	"hematology":            "Haematology",
	"immunohistochemistry":  "Immunohistochemistry",
	"immunology":            "Immunology",
	"microbiology":          "Microbiology",
	"radiology":             "Radiology",
	"molecular biology":     "Molecular Biology",
	"others":                "Others",
	"outsource":             "Outsource",
	"coagulation":           "Coagulation",
	"biochemistry":          "Biochemistry",
	"flowcytometry":         "Flowcytometry",
	"immunofluorometry":     "Immunofluorometry",
	"molecular pathology":   "Molecular Pathology",
	"clinical pathology":    "Clinical Pathology",
	"cytogenetics":          "Cytogenetics",
	"cytology":              "Cytology",
	"clinical chemistry":    "Clinical Chemistry",
	"serology":              "Serology",
	"vitals":                "Vitals",
	"histopathology":        "Histopathology",
	"immunology & serology": "Immunology & Serology",
	"advanced serology":     "Advanced Serology",
	"na":                    "NA",
}

// MethodTypes
const (
	METHOD_TYPE_MANUAL          = "manual"
	METHOD_TYPE_DEVICE_MEASURED = "device_measured"
	METHOD_TYPE_CALCULATED      = "calculated"
)

const (
	APPROVAL_SOURCE_NA   = "NA"
	APPROVAL_SOURCE_IM   = "IM"
	APPROVAL_SOURCE_LIMS = "LIMS"
	APPROVAL_SOURCE_OH   = "OH"
)

const (
	AUTO_APPROVAL_FAIL_REASON_NA                          = "na"
	AUTO_APPROVAL_FAIL_REASON_PAST_RECORD                 = "past_record"
	AUTO_APPROVAL_FAIL_REASON_REF_RANGE                   = "ref_range"
	AUTO_APPROVAL_FAIL_REASON_INVALID_INVESTIGATION_VALUE = "invalid_investigation_value"
	AUTO_APPROVAL_FAIL_REASON_IM_DEVICE                   = "im_device"
	AUTO_APPROVAL_FAIL_REASON_MANUAL_INPUT                = "manual_input"
	AUTO_APPROVAL_FAIL_REASON_QC_FAILED                   = "qc_failed"
)

var (
	DeltaCheckWhitelistedMasterInvIds = Config.GetIntSlice("investigation.delta_check_whitelisted_master_investigation_ids")
	EnableAutoApprovalSlackAlerts     = Config.GetBool("investigation.enable_auto_approval_slack_alerts")
)

// Master Investigation Result Types
const (
	ResultTypeNumeric          = "numeric"
	ResultTypeTextual          = "textual"
	ResultTypeSemiQuantitative = "semi_quantitative"
)

const (
	ReferenceRangeLabelRange   = "range"
	ReferenceRangeLabelTextual = "textual"
)

// Abnormality Statuses
var (
	OhAbnormalityStringSlice = []string{
		ABNORMALITY_LOWER_ABNORMAL,
		ABNORMALITY_UPPER_ABNORMAL,
		ABNORMALITY_CRITICAL,
		ABNORMALITY_IMPROBABLE,
	}

	OhCriticalityStringSlice = []string{
		ABNORMALITY_CRITICAL,
		ABNORMALITY_IMPROBABLE,
	}
)

// Slack Alert Names
const (
	AA_FAIL_ALERT_PAST_RECORD = "Result outside delta range!"
	AA_FAIL_ALERT_REF_RANGE   = "Result outside AA range, within delta variation!"
)

// Debug/Error Messages
const (
	SKIP_DELTA_CHECK_WHITELISTED    = "Skipping delta check for whitelisted investigation id: %d"
	EMPTY_ENTERED_TIME              = "enteredTime is empty for investigation: %s"
	LAST_INVESTIGATION_RESULT       = "Last investigation result for %s: %v"
	DAYS_DIFF_LAST_INV_CURRENT_INV  = "Difference in days between last investigation result and current entered time: %d"
	SKIP_DELTA_CHECK_OLD_PAST_VALUE = "Skipping delta check for investigation: %d as the last result is older than the past value threshold days"
)
