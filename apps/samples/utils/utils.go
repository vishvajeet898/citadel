package utils

import (
	commonConstants "github.com/Orange-Health/citadel/common/constants"
	commonModels "github.com/Orange-Health/citadel/models"
)

func GetMasterTestToProcessingLabMap(testDetails []commonModels.TestDetail) map[uint]uint {
	masterTestIdToProcessingLabIdMap := make(map[uint]uint)
	for _, testDetail := range testDetails {
		masterTestIdToProcessingLabIdMap[testDetail.MasterTestId] = testDetail.ProcessingLabId
	}
	return masterTestIdToProcessingLabIdMap
}

func CreateSampleNumberMaps(testDetails []commonModels.TestDetail,
	testSampleMappings []commonModels.TestSampleMapping) map[uint][]commonModels.TestDetail {

	omsTestIdToTestDetailsMap := map[string]commonModels.TestDetail{}

	for _, testDetail := range testDetails {
		omsTestIdToTestDetailsMap[testDetail.CentralOmsTestId] = testDetail
	}

	sampleNumberToTestsResponseMap := map[uint][]commonModels.TestDetail{}
	for _, testSampleMapping := range testSampleMappings {
		testDetail := omsTestIdToTestDetailsMap[testSampleMapping.OmsTestId]
		if sampleNumberToTestsResponseMap[testSampleMapping.SampleNumber] == nil {
			sampleNumberToTestsResponseMap[testSampleMapping.SampleNumber] = []commonModels.TestDetail{}
		}
		sampleNumberToTestsResponseMap[testSampleMapping.SampleNumber] = append(
			sampleNumberToTestsResponseMap[testSampleMapping.SampleNumber], testDetail)
	}

	return sampleNumberToTestsResponseMap
}

func GetGroupIdAndNameForFreshdeskTicketBasedOnReason(reason string) (uint, string) {
	if reason == "Results on Hold – No Repeat (NA)" {
		return commonConstants.FreshDeskCRTGroupId, commonConstants.FreshDeskCRTCreatorName
	}

	return commonConstants.FreshDeskCentralLogisticsGroupId, commonConstants.FreshDeskCentralLogisticsCreatorName
}
