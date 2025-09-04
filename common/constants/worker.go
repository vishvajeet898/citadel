package constants

var (
	WorkerBroker           = Config.GetString("worker.broker")
	WorkerDefaultQueue     = Config.GetString("worker.default_queue")
	WorkerDefaultQueueFifo = Config.GetString("worker.default_queue_fifo")
	WorkerResultBackend    = Config.GetString("worker.result_backend")
	WorkerRedisPoolSize    = Config.GetInt("worker.redis_pool_size")
	WorkerMaxIdle          = Config.GetInt("worker.redis_max_idle")
	WorkerProcesses        = Config.GetInt("worker.processes")
	CitadelWorker          = "citadel_worker"
)

const (
	ConsumerEvents                   = "consumer_events"
	ReleaseReportTask                = "release_report_task"
	UpdateTaskPostSavingTask         = "update_task_post_saving_task"
	SampleRejectionTask              = "sample_reject_task"
	UpdateTaskPostReceivingTask      = "update_task_post_receiving_task"
	CreateUpdateTaskByOmsOrderIdTask = "create_update_task_by_oms_order_id_task"

	// Periodic Tasks
	StaleTasksPeriodicTask                        = "stale_tasks_periodic_task"
	EtsTatBreachedPeriodicTaskName                = "ets_tat_breached_periodic_task"
	SlackAlertForReverseLogisticsPeriodicTaskName = "slack_alert_for_reverse_logistics_periodic_task"
)

// periodic tasks frequency
var (
	StaleTasksPeriodicTaskFrequency                    = Config.GetString("cron_frequencies.stale_tasks_periodic_task")
	EtsTatBreachedPeriodicTaskFrequency                = Config.GetString("cron_frequencies.ets_tat_breached_frequency")
	SlackAlertForReverseLogisticsPeriodicTaskFrequency = Config.GetString("cron_frequencies.slack_alert_for_reverse_logistics_frequency")
)

var (
	StaleTaskDuration = Config.GetInt("stale_task_duration")
)
