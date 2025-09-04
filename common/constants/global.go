package constants

import "time"

type TraceIdType string

const (
	TraceIdKey TraceIdType = "trace_id"
)

const (
	CitadelServiceName = "CITADEL"
)

var (
	Environment        = Config.GetString("env")
	ActiveCityCodes    = Config.GetStringSlice("active_city_codes")
	CitadelSystemId    = Config.GetUint("citadel_system_id")
	LisSystemId        = Config.GetUint("lis_system_id")
	AutoApprovalIdsMap = Config.GetStringMapString("auto_approval_ids_map")
	OmsBaseDomain      = Config.GetString("oms.base_domain")
	SuperlabBaseUrl    = Config.GetString("superlab.base_url")
)

const (
	DateLayout                               = "2006-01-02"
	DateLayoutReverse                        = "02-01-2006"
	TimeLayout                               = "15:04"
	DateTimeLayout                           = "2006-01-02 15:04"
	DateTimeInSecLayout                      = "2006-01-02 15:04:05"
	DateTimeUTCLayout                        = "2006-01-02T15:04:05Z07:00"
	DateTimeUTCLayoutWithoutTZOffset         = "2006-01-02T15:04:05Z"
	PrettyDateTimeLayout                     = "02 Jan 2006 03:04 PM"
	TimeStampLayout                          = "20060102150405"
	DateTimeUTCWithFractionSecWithoutZOffset = "2006-01-02T15:04:05.0000000Z"
	DateTimeAMPMFormat                       = "2006-01-02 | 03:04 PM"
	LocalTimeZoneLocation                    = "Asia/Kolkata"
)

const (
	InvestigationShortHand = "INV"
	GroupShortHand         = "GRP"
)

const (
	SourceAttune = "attune"

	SourceOfApprovalIM   = "IM"
	SourceOfApprovalLIMS = "LIMS"
)

const (
	SmallYes     = "yes"
	AutoVerified = "av" // AV represents AutoVerified
)

const (
	BiopatternRepresentationType = "biopattern"

	CultureInvValueTextPrefix = "<InvestigationResults>"

	InvestigationValue = "InvestigationValue"
)

// Service Names
const (
	PorteServiceName  = "PORTE"
	GrootServiceName  = "GROOT"
	OmsServiceName    = "OMS"
	EtsServiceName    = "ETS"
	HealthServiceName = "HEALTH"
)

// Content Types
const (
	ContentTypeJson                = "application/json"
	ContentTypeJsonWithCharsetUtf8 = "application/json; charset=utf-8"
)

const (
	DefaultHTTPClientTimeout = 30 * time.Second
)

var (
	TestCodesToBeSkippedInRetry = []string{"ODB002"}
	DefaultAutoApprovalCodes    = []string{"OD0246", "OD0247", "OD0248"}
)

const (
	LevenshteinDistanceThreshold = 2
)

const (
	NrlLabId = 67
)

const (
	SuperlabSystemName = "Superlab"
	CitadelSystemName  = "Citadel"
)
