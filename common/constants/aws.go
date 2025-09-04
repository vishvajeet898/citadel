package constants

var (
	Bucket              = Config.GetString("aws.bucket")
	S3AccessKeyID       = Config.GetString("aws.S3_AWS_ACCESS_KEY_ID")
	S3SecretAccessKey   = Config.GetString("aws.S3_AWS_SECRET_ACCESS_KEY")
	SQSAccessKeyID      = Config.GetString("aws.SQS_AWS_ACCESS_KEY_ID")
	SQSSecretAccessKey  = Config.GetString("aws.SQS_AWS_SECRET_ACCESS_KEY")
	SQSQueueURL         = Config.GetString("aws.SQS_QUEUE_URL")
	StandardSQSQueueURL = Config.GetString("aws.SQS_STANDARD_QUEUE_URL")
	Region              = Config.GetString("aws.AWS_REGION")
	SNSAccessKeyID      = Config.GetString("aws.SNS_AWS_ACCESS_KEY_ID")
	SNSSecretAccessKey  = Config.GetString("aws.SNS_AWS_SECRET_ACCESS_KEY")
)

// Topic ARNs
var (
	OrderReportUpdateTopicArn    = Config.GetString("topics.order_report_update")
	OrderResultsApprovedTopicArn = Config.GetString("topics.order_results_approved")
	ContactMergeConfirmTopicArn  = Config.GetString("topics.contact_merge_confirm")
	ReportReadyTopicArn          = Config.GetString("topics.report_ready")
	OmsUpdatesTopicArn           = Config.GetString("topics.oms_updates")
	CitadelTopicArn              = Config.GetString("topics.citadel")
	EtsTestEventTopicArn         = Config.GetString("topics.ets_test_event")
)

var QueuesToPoll = []string{
	SQSQueueURL,
	StandardSQSQueueURL,
}
