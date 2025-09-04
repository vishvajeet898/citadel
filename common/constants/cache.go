package constants

import "time"

var (
	PubSubRedisAddr        string = Config.GetString("redis.pubsub.address")
	PubSubRedisPass        string = Config.GetString("redis.pubsub.password")
	PubSubRedisDb          int    = Config.GetInt("redis.pubsub.db")
	PubSubRedisPoolSize    int    = Config.GetInt("redis.pubsub.poolsize")
	PubSubRedisMinIdleConn int    = Config.GetInt("redis.pubsub.minidle")
)

// Cache Keys
const (
	// Event Keys
	ReportReleaseKey             = "report_release::task_id:%d"
	UpdateTaskPostSavingTaskKey  = "update_task_post_saving_task::task:%d"
	AmendmendTaskCountKey        = "amendment_task_count"
	LisEventKey                  = "lis_event::visit_id:%s"
	LisEventReportKey            = "lis_event_report::order_id:%s::test_ids:%s::city_code:%s"
	ReportPdfReadyEventKey       = "report_pdf_ready_event::order:%s::city_code:%s"
	OmsLisEventKey               = "oms_lis_event::order:%s::city_code:%s"
	OmsAttachmentEventKey        = "oms_attachment_event::order:%s::city_code:%s"
	OmsTestCompletionEventKey    = "oms_test_completion_event::order:%d::city_code:%s"
	OmsManualUploadKey           = "oms_manual_upload::order:%s::city_code:%s"
	OmsRerunEventKey             = "oms_rerun_event::order:%s::city_code:%s"
	OmsCreateUpdateOrderEventKey = "oms_create_update_order_event::order:%s"
	OmsCompletedOrderEventKey    = "oms_completed_order_event::order:%d::city_code:%s"
	OmsCancelledOrderEventKey    = "oms_cancelled_order_event::order:%d::city_code:%s"
	OmsSampleEventKey            = "oms_sample_event::sample:%d::city_code:%s"
	OmsSampleDeleteKey           = "oms_sample_delete::sample:%d::city_code:%s"
	CreateUpdateTaskTaskKey      = "create_update_task::order_id:%s"

	// Normal Keys
	CacheKeyDepartmentMapping       = "department_mapping"
	CacheKeyInvestigationDetails    = "investigation_details_%s_%s_%d_%s_%s"
	CacheKeyDependentInvestigations = "dependent_investigations_%s_%d_%d_%s_%s"
	CacheKeyTaskPathlogistMapping   = "task_pathologist_mapping_%d"
	CacheKeyVialsAll                = "master_vial_types:all"
	CacheKeyVialsMap                = "master_vial_types:map"
	CacheKeyVialsId                 = "master_vial_types:%v"
	CacheKeyLabsAll                 = "labs:all"
	CacheKeyLabsMap                 = "labs:map"
	CacheKeyLabsId                  = "labs:%v"
	CacheKeyOutSourceLabIds         = "outsource_labs:ids"
	CacheKeyInHouseLabIds           = "inhouse_labs:ids"
	CacheKeyPartnerDetails          = "partner_details_%d"
	CacheKeyCurrentInProgressSample = "current_in_progress_sample_%s"
	CacheKeySyncToLisVisitCount     = "sync_to_lis_visit_count:%s"
	CacheKeyMasterTestAll           = "master_test:all"
	CacheKeyMasterTestMap           = "master_test:map"
	CacheKeyMasterTestId            = "master_test:%v"
	CacheKeyNrlEnabledMasterTestIds = "nrl_enabled_master_test_ids"
	CacheKeyCollectionSequence      = "collection_sequence:%s_%s"
	CacheKeyDoctorDetails           = "doctor_details_%d"
	CacheKeyReceiveAndSync          = "receive_and_sync:%s"
	CacheKeySrfOrderIds             = "srf_order_ids:%s"
	CacheKeySrfOrderIdsAll          = "srf_order_ids:all"
)

// Cache Expiry Time Duration
const (
	CacheExpiry24Hours = 24 * time.Hour
)

// Cache Expiry Time Duration in int
const (
	CacheExpiry1MinutesInt  = 1 * 60
	CacheExpiry2MinutesInt  = 2 * 60
	CacheExpiry5MinutesInt  = 5 * 60
	CacheExpiry10MinutesInt = 10 * 60
	CacheExpiry15MinutesInt = 15 * 60

	CacheExpiry1HourInt = 1 * 60 * 60

	CacheExpiry1DayInt   = 1 * 24 * 60 * 60
	CacheExpiry10DaysInt = 10 * 24 * 60 * 60
)

var (
	TaskPathologistMappingKeyExpiry = 60 * StaleTaskDuration
)
