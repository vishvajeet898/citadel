package consumerTasks

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/Orange-Health/citadel/common/constants"
	"github.com/Orange-Health/citadel/common/structures"
)

func (eventProcessor *EventProcessor) TestDetailsEventTask(ctx context.Context, eventPayload string) error {
	testEtaPayload := structures.TestDetailsEvent{}
	err := json.Unmarshal(json.RawMessage(eventPayload), &testEtaPayload)
	if err != nil {
		return err
	}

	testIdToTestEtaDetailsMap, omsTestIds := map[string]structures.TestDetailsEventData{}, []string{}

	for _, test := range testEtaPayload.TestDetails {
		if test.TestId == "" {
			continue
		}
		testIdToTestEtaDetailsMap[test.TestId] = test
		omsTestIds = append(omsTestIds, test.TestId)
	}
	if len(omsTestIds) == 0 {
		return nil
	}
	testDetails, cErr := eventProcessor.TestDetailService.GetTestDetailModelByOmsTestIds(omsTestIds)
	if cErr != nil {
		return errors.New(cErr.Message)
	}

	completedTestIds := []string{}

	for index := range testDetails {
		testDetails[index].LabEta = getLabEta(testIdToTestEtaDetailsMap[testDetails[index].CentralOmsTestId].LabEta)
		testDetails[index].ReportEta = getReportEta(testIdToTestEtaDetailsMap[testDetails[index].CentralOmsTestId].ReportTat)
		testDetails[index].LabTat = testIdToTestEtaDetailsMap[testDetails[index].CentralOmsTestId].LabTat
		testDetails[index].OmsStatus = getOmsTestStatus(testIdToTestEtaDetailsMap[testDetails[index].CentralOmsTestId].Status)
		if testIdToTestEtaDetailsMap[testDetails[index].CentralOmsTestId].Status == constants.TestStatusCompletedNotSentUint ||
			testIdToTestEtaDetailsMap[testDetails[index].CentralOmsTestId].Status == constants.TestStatusCompletedSentUint {
			completedTestIds = append(completedTestIds, testDetails[index].CentralOmsTestId)
		}
		if testIdToTestEtaDetailsMap[testDetails[index].CentralOmsTestId].Status == constants.TestStatusCompletedNotSentUint {
			testDetails[index].ReportStatus = constants.TEST_REPORT_STATUS_SENT
			completedTestIds = append(completedTestIds, testDetails[index].CentralOmsTestId)
		}
		if testIdToTestEtaDetailsMap[testDetails[index].CentralOmsTestId].Status == constants.TestStatusCompletedSentUint {
			testDetails[index].ReportStatus = constants.TEST_REPORT_STATUS_DELIVERED
		}
	}

	_, cErr = eventProcessor.TestDetailService.UpdateTestDetails(testDetails)
	if cErr != nil {
		return errors.New(cErr.Message)
	}

	eventProcessor.EtsService.GetAndPublishEtsTestBasicEvent(ctx, completedTestIds)

	return nil
}
