package consumerTasks

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/Orange-Health/citadel/common/constants"
	"github.com/Orange-Health/citadel/common/structures"
	"github.com/Orange-Health/citadel/common/utils"
)

func (eventProcessor *EventProcessor) OrderReportPdfReadyEventTask(ctx context.Context, eventPayload string) error {
	reportReadyEvent := structures.ReportReadyEvent{}
	err := json.Unmarshal(json.RawMessage(eventPayload), &reportReadyEvent)
	if err != nil {
		eventProcessor.Sentry.LogError(ctx, constants.ERROR_FAILED_TO_UNMARSHAL_JSON, err, nil)
		return err
	}

	redisKey := fmt.Sprintf(constants.ReportPdfReadyEventKey, reportReadyEvent.ReportPdfEvent.OrderID,
		reportReadyEvent.ReportPdfEvent.CityCode)
	keyExists, err := eventProcessor.Cache.Exists(ctx, redisKey)
	if err != nil || keyExists {
		utils.AddLog(ctx, constants.DEBUG_LEVEL, utils.GetCurrentFunctionName(), nil, err)
		return errors.New(constants.ERROR_LIS_EVENT_TASK_IN_PROGRESS)
	}

	err = eventProcessor.Cache.Set(ctx, redisKey, true, constants.CacheExpiry10MinutesInt)
	if err != nil {
		utils.AddLog(ctx, constants.DEBUG_LEVEL, utils.GetCurrentFunctionName(), nil, err)
		return err
	}

	defer func() {
		err := eventProcessor.Cache.Delete(ctx, redisKey)
		if err != nil {
			utils.AddLog(ctx, constants.DEBUG_LEVEL, utils.GetCurrentFunctionName(), nil, err)
		}
	}()
	if reportReadyEvent.ReportPdfEvent.IsDummyReport {
		publicURL, err := eventProcessor.S3wrapperClient.GetTokenizeOrderFilePublicUrl(
			ctx,
			reportReadyEvent.ReportPdfEvent.ReportPdfBrandedURL,
		)
		if err != nil {
			publicURL = reportReadyEvent.ReportPdfEvent.ReportPdfBrandedURL // fallack if not getting public URL
		}
		slackBlocks := []map[string]interface{}{
			{
				"type": "section",
				"text": map[string]interface{}{
					"type": "mrkdwn",
					"text": fmt.Sprintf("*:alphabet-white-d: :alphabet-white-u: :alphabet-white-m: :alphabet-white-m: :alphabet-white-y: DUMMY Report Ready* 🧪\n*Order ID:* %s\n*Report URL:* <%s|View Report>",
						reportReadyEvent.ReportPdfEvent.OrderID,
						publicURL),
				},
			},
		}
		err = eventProcessor.SlackClient.SendToSlackDirectly(ctx, constants.SlackPdfGenerationChannel, slackBlocks)
		if err != nil {
			utils.AddLog(ctx, constants.ERROR_LEVEL, "FailedToSendSlackMessage", map[string]interface{}{
				"order_id": reportReadyEvent.ReportPdfEvent.OrderID,
				"error":    err.Error(),
			}, nil)
		}
		return nil
	}

	cErr := eventProcessor.TestDetailService.UpdateReportStatusByOmsTestIds(reportReadyEvent.ReportPdfEvent.TestIds,
		constants.TEST_REPORT_STATUS_QUEUED, constants.TEST_REPORT_STATUS_CREATED)
	if cErr != nil {
		return errors.New(cErr.Message)
	}

	reportReadyEvent.ServicingCityCode = reportReadyEvent.ReportPdfEvent.CityCode
	messageBody, messageAttributes := eventProcessor.PubsubService.GetReportReadyEvent(ctx, reportReadyEvent)
	cErr = eventProcessor.SnsClient.PublishTo(ctx, messageBody, messageAttributes,
		fmt.Sprint(reportReadyEvent.ReportPdfEvent.OrderID), constants.ReportReadyTopicArn,
		reportReadyEvent.ReportPdfEvent.JobUUID)
	if cErr != nil {
		utils.AddLog(ctx, constants.ERROR_LEVEL, utils.GetCurrentFunctionName(), nil, errors.New(cErr.Message))
		return errors.New(cErr.Message)
	}

	cErr = eventProcessor.TestDetailService.UpdateReportStatusByOmsTestIds(reportReadyEvent.ReportPdfEvent.TestIds,
		constants.TEST_REPORT_STATUS_CREATED, constants.TEST_REPORT_STATUS_SYNCED)
	if cErr != nil {
		return errors.New(cErr.Message)
	}

	return nil
}
