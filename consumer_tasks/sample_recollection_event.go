package consumerTasks

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/Orange-Health/citadel/common/constants"
	"github.com/Orange-Health/citadel/common/structures"
	"github.com/Orange-Health/citadel/common/utils"
	"github.com/Orange-Health/citadel/models"
	"gorm.io/gorm"
)

func (eventProcessor *EventProcessor) SampleRecollectionEventTask(ctx context.Context, eventPayload string) error {
	sampleRecollectionPayload := structures.SampleRecollectionEvent{}
	err := json.Unmarshal(json.RawMessage(eventPayload), &sampleRecollectionPayload)
	if err != nil {
		return err
	}
	utils.AddLog(ctx, constants.INFO_LEVEL, utils.GetCurrentFunctionName(),
		map[string]interface{}{
			"event_payload": eventPayload,
			"event_type":    constants.SampleRecollectionEvent,
		}, nil)

	omsOrderIds, omsTestIds := []string{}, []string{}
	omsOrderIdToOrderDetailsMap := map[string]models.OrderDetails{}
	omsOrderIdToRecollectionTestDetailsMap, omsOrderIdToNonRecollectionTestDetailsMap, omsOrderIdToCollectLaterTestDetailsMap :=
		map[string][]models.TestDetail{}, map[string][]models.TestDetail{}, map[string][]models.TestDetail{}
	for _, test := range sampleRecollectionPayload.RecollectionOrderDetails {
		omsOrderIds = append(omsOrderIds, test.OrderId)
		omsTestIds = append(omsTestIds, test.TestIds...)
	}

	orderDetails, cErr := eventProcessor.OrderDetailsService.GetOrderDetailsByOmsOrderIds(omsOrderIds)
	if cErr != nil {
		return errors.New(cErr.Message)
	}
	for _, orderDetail := range orderDetails {
		omsOrderIdToOrderDetailsMap[orderDetail.OmsOrderId] = orderDetail
	}

	testDetails, cErr := eventProcessor.TestDetailService.GetTestDetailModelByOmsTestIds(omsTestIds)
	if cErr != nil {
		return errors.New(cErr.Message)
	}
	recollectionTestStatuses := []string{
		constants.TEST_STATUS_REJECTED,
		constants.TEST_STATUS_SAMPLE_NOT_RECEIVED,
	}
	sampleCollectLaterOmsTestIds := []string{}
	for _, testDetail := range testDetails {
		if _, ok := omsOrderIdToRecollectionTestDetailsMap[testDetail.OmsOrderId]; !ok {
			omsOrderIdToRecollectionTestDetailsMap[testDetail.OmsOrderId] = []models.TestDetail{}
		}
		if _, ok := omsOrderIdToNonRecollectionTestDetailsMap[testDetail.OmsOrderId]; !ok {
			omsOrderIdToNonRecollectionTestDetailsMap[testDetail.OmsOrderId] = []models.TestDetail{}
		}
		if _, ok := omsOrderIdToCollectLaterTestDetailsMap[testDetail.OmsOrderId]; !ok {
			omsOrderIdToCollectLaterTestDetailsMap[testDetail.OmsOrderId] = []models.TestDetail{}
		}
		if testDetail.Status == constants.TEST_STATUS_COLLECT_SAMPLE_LATER {
			sampleCollectLaterOmsTestIds = append(sampleCollectLaterOmsTestIds, testDetail.CentralOmsTestId)
			omsOrderIdToCollectLaterTestDetailsMap[testDetail.OmsOrderId] = append(
				omsOrderIdToCollectLaterTestDetailsMap[testDetail.OmsOrderId], testDetail)
		} else if utils.SliceContainsString(recollectionTestStatuses, testDetail.Status) {
			omsOrderIdToRecollectionTestDetailsMap[testDetail.OmsOrderId] = append(
				omsOrderIdToRecollectionTestDetailsMap[testDetail.OmsOrderId], testDetail)
		} else {
			omsOrderIdToNonRecollectionTestDetailsMap[testDetail.OmsOrderId] = append(
				omsOrderIdToNonRecollectionTestDetailsMap[testDetail.OmsOrderId], testDetail)
		}
	}

	collectLaterSamples, collectLaterSampleMetadatas, collectLaterTestDetails :=
		[]models.Sample{}, []models.SampleMetadata{}, []models.TestDetail{}
	if len(sampleCollectLaterOmsTestIds) > 0 {
		collectLaterSamples, collectLaterSampleMetadatas, cErr =
			eventProcessor.SampleService.GetSamplesForTests(sampleCollectLaterOmsTestIds)
		if cErr != nil {
			return errors.New(cErr.Message)
		}
		sampleIds := []uint{}
		for index := range collectLaterSampleMetadatas {
			collectLaterSampleMetadatas[index].CollectLaterReason = ""
			collectLaterSampleMetadatas[index].TaskSequence = sampleRecollectionPayload.TaskId
		}
		for index := range collectLaterSamples {
			sampleIds = append(sampleIds, collectLaterSamples[index].Id)
			collectLaterSamples[index].Status = constants.SampleDefault
			collectLaterSamples[index].UpdatedBy = constants.CitadelSystemId
		}
		collectLaterTestDetails, cErr = eventProcessor.SampleService.GetTestDetailsBySampleIds(sampleIds)
		if cErr != nil {
			return errors.New(cErr.Message)
		}
		for index := range collectLaterTestDetails {
			collectLaterTestDetails[index].Status = constants.TEST_STATUS_REQUESTED
			collectLaterTestDetails[index].UpdatedBy = constants.CitadelSystemId
		}
	}

	txErr := eventProcessor.Db.Transaction(func(tx *gorm.DB) error {
		for omsOrderId, orderDetails := range omsOrderIdToOrderDetailsMap {
			cErr := eventProcessor.SampleService.CreateSamplesForRecollectionWithTx(ctx, tx, orderDetails,
				omsOrderIdToRecollectionTestDetailsMap[omsOrderId], sampleRecollectionPayload.TaskId)
			if cErr != nil {
				return errors.New(cErr.Message)
			}

			nonRecollectedOmsTestIds := []string{}
			for _, testDetail := range omsOrderIdToNonRecollectionTestDetailsMap[omsOrderId] {
				nonRecollectedOmsTestIds = append(nonRecollectedOmsTestIds, testDetail.CentralOmsTestId)
			}

			cErr = eventProcessor.SampleService.UpdateTaskIdByOmsTestIdsWithTx(tx, nonRecollectedOmsTestIds,
				sampleRecollectionPayload.TaskId)
			if cErr != nil {
				return errors.New(cErr.Message)
			}

			if len(collectLaterSamples) > 0 {
				_, _, cErr = eventProcessor.SampleService.UpdateSamplesAndSamplesMetadataWithTx(tx,
					collectLaterSamples, collectLaterSampleMetadatas)
				if cErr != nil {
					return errors.New(cErr.Message)
				}
			}

			if len(collectLaterTestDetails) > 0 {
				_, cErr = eventProcessor.TestDetailService.UpdateTestDetailsWithTx(tx, collectLaterTestDetails)
				if cErr != nil {
					return errors.New(cErr.Message)
				}
			}

			cErr = eventProcessor.TestSampleMappingService.UpdateTsmForTestIdAndRecollectionPendingTrueWithTx(tx,
				nonRecollectedOmsTestIds)
			if cErr != nil {
				return errors.New(cErr.Message)
			}
		}
		return nil
	})
	if txErr != nil {
		return errors.New(txErr.Error())
	}

	recollectionPendingPresent, cErr := eventProcessor.TestSampleMappingService.OrdersWithRecollectionsPendingPresent(
		omsOrderIds, omsTestIds)
	if cErr != nil {
		utils.AddLog(ctx, constants.ERROR_LEVEL, utils.GetCurrentFunctionName(), nil, errors.New(cErr.Message))
	}

	omsOrderIdToTestIdToTestStatusMap := map[string]map[string]string{}
	for _, testDetail := range collectLaterTestDetails {
		if _, ok := omsOrderIdToTestIdToTestStatusMap[testDetail.OmsOrderId]; !ok {
			omsOrderIdToTestIdToTestStatusMap[testDetail.OmsOrderId] = map[string]string{}
		}
		omsOrderIdToTestIdToTestStatusMap[testDetail.OmsOrderId][testDetail.CentralOmsTestId] = constants.TEST_STATUS_REQUESTED
	}

	for omsOrderId, testIdToTestStatus := range omsOrderIdToTestIdToTestStatusMap {
		go eventProcessor.SampleService.PublishUpdateTestStatusEvent(testIdToTestStatus, omsOrderId, true,
			orderDetails[0].CityCode)
	}

	go eventProcessor.SampleService.PublishRemoveSampleRejectedTagEvent(sampleRecollectionPayload.RequestId, omsOrderIds,
		!recollectionPendingPresent, orderDetails[0].CityCode)

	return nil
}
