package constants

// Production Slack channels
var (
	SlackReportSyncFailureChannel   = Config.GetString("slack.report_sync_failures_channel")
	SlackMissingParametersChannel   = Config.GetString("slack.missing_parameters_channel")
	SlackSampleLossRiskAlertChannel = Config.GetString("slack.sample_loss_risk_alert_channel")
	SlackAutoApprovalFailureChannel = Config.GetString("slack.auto_approval_failures_channel")
	SlackPdfGenerationChannel       = Config.GetString("slack.pdf_generation_channel")
)

// Staging Slack channels
var (
	SlackStagingCommunicationChannel = "staging-citadel-slack-communication"
)

const (
	MaxBlocksForSlackMessage = 45 // Maximum is 50 but geerally some of the blocks are taken by headers
)
