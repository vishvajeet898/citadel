package utils

import (
	"encoding/json"
	"fmt"

	"github.com/Orange-Health/citadel/common/constants"
	"github.com/Orange-Health/citadel/common/structures"
)

func ValidatePayload(lisOrderUpdateDetails structures.LisOrderUpdateDetails,
	omsTestDetails []structures.TestDetailsForLisEvent, testCodeInvestigationMap map[string]structures.Investigation,
	testCodeMimmMap map[string]uint, mimmTestCodeMap map[uint][]string, attunePathologistsUserIds []string) (
	[]string, []string, []string, []string, map[string]string, []string) {

	valueMap, testCodePresentMap, testCodeNameMap := map[string]string{}, map[string]bool{}, map[string]string{}

	valueMap, testCodePresentMap, testCodeNameMap, doctorDetailsMissingTests :=
		createValueMapAndTestCodeMapForInvestigations(lisOrderUpdateDetails, omsTestDetails,
			constants.AttuneTestStatusApprove, valueMap, testCodePresentMap, testCodeNameMap, attunePathologistsUserIds)
	valueMap, testCodePresentMap, testCodeNameMap, _ =
		createValueMapAndTestCodeMapForInvestigations(lisOrderUpdateDetails, omsTestDetails,
			constants.AttuneTestStatusCompleted, valueMap, testCodePresentMap, testCodeNameMap, attunePathologistsUserIds)

	// Missing Parameters
	missingTestsInCds, missingTestsInLis := checkForMissingParameters(testCodePresentMap, testCodeMimmMap, mimmTestCodeMap)

	// Blank Values & Incorrect Result Type
	blankTests, incorrectResultTypes := checkForBlankValuesAndIncorrectResultType(valueMap, testCodeInvestigationMap)

	return missingTestsInCds, missingTestsInLis, blankTests, incorrectResultTypes, testCodeNameMap, doctorDetailsMissingTests
}

func checkForMissingParameters(lisTestCodePresentMap map[string]bool, testCodeMimmMap map[string]uint,
	mimmTestCodeMap map[uint][]string) ([]string, []string) {

	cdsPresentMap := map[string]bool{}
	missingTestsInCds, missingTestsInLis := []string{}, []string{}
	for lisTestCode := range lisTestCodePresentMap {
		if _, keyExists := testCodeMimmMap[lisTestCode]; keyExists {
			mimm := testCodeMimmMap[lisTestCode]
			allCdsTestCodes := mimmTestCodeMap[mimm]
			for _, cdsTestCode := range allCdsTestCodes {
				cdsPresentMap[cdsTestCode] = true
			}
		} else {
			missingTestsInCds = append(missingTestsInCds, lisTestCode)
		}
	}

	for cdsTestCode := range testCodeMimmMap {
		if _, keyExists := cdsPresentMap[cdsTestCode]; !keyExists {
			missingTestsInLis = append(missingTestsInLis, cdsTestCode)
		}
	}

	return missingTestsInCds, missingTestsInLis
}

func checkForBlankValuesAndIncorrectResultType(valueMap map[string]string,
	testCodeInvestigationMap map[string]structures.Investigation) ([]string, []string) {

	blankTests, incorrectResultTypes := []string{}, []string{}

	for testCode, value := range valueMap {
		blankValue := checkIfTestValueIsEmpty(value)
		if blankValue {
			blankTests = append(blankTests, testCode)
		}
		if !blankValue && checkIfTestValueHasIncorrectResultType(value, testCodeInvestigationMap[testCode].ResultType) {
			incorrectResultTypes = append(incorrectResultTypes, testCode)
		}
	}

	return blankTests, incorrectResultTypes
}

func checkIfTestValueIsEmpty(testValue string) bool {
	return testValue == ""
}

func checkIfTestValueHasIncorrectResultType(testValue string, resultType string) bool {
	switch resultType {
	case constants.ResultTypeNumeric:
		return !(NumericResultTypeRegex.MatchString(testValue))
	case constants.ResultTypeTextual:
		return !(TextualResultTypeRegex.MatchString(testValue))
	case constants.ResultTypeSemiQuantitative:
		return !(SemiQuantitativeResultTypeRegex.MatchString(testValue))
	}

	return false
}

func createValueMapAndTestCodeMapForInvestigations(lisOrderUpdateDetails structures.LisOrderUpdateDetails,
	testDetails []structures.TestDetailsForLisEvent, attuneStatus string, valueMap map[string]string,
	testCodePresentMap map[string]bool, testCodeNameMap map[string]string, attunePathologistsUserIds []string) (
	map[string]string, map[string]bool, map[string]string, []string) {

	doctorDetailsMissingTests := []string{}
	orderInfo := lisOrderUpdateDetails.OrderInfo[attuneStatus]

	for _, test := range testDetails {
		if testResults, keyExists := orderInfo[test.TestCode]; keyExists {
			if SliceContainsString(constants.TestCodesToBeSkippedInRetry, test.TestCode) {
				continue
			}
			orderInfo := testResults.MetaData
			resultApprovedBy, resultApprovedByString := orderInfo.ResultApprovedBy, fmt.Sprint(orderInfo.ResultApprovedBy)
			if orderInfo.TestType == constants.InvestigationShortHand {
				valueMap[orderInfo.TestCode] = orderInfo.TestValue
				testCodePresentMap[orderInfo.TestCode] = true
				testCodeNameMap[orderInfo.TestCode] = orderInfo.TestName
				if attuneStatus == constants.AttuneTestStatusApprove {
					if resultApprovedBy == 0 || !SliceContainsString(attunePathologistsUserIds, resultApprovedByString) {
						doctorDetailsMissingTests = append(doctorDetailsMissingTests, orderInfo.TestName)
					}
				}
			} else if orderInfo.TestType == constants.GroupShortHand {
				index, queue := 0, []interface{}{}
				queue = append(queue, orderInfo.OrderContentListInfo)

				if attuneStatus == constants.AttuneTestStatusApprove {
					if resultApprovedBy == 0 || !SliceContainsString(attunePathologistsUserIds, resultApprovedByString) {
						doctorDetailsMissingTests = append(doctorDetailsMissingTests, orderInfo.TestName)
					}
				}

				for index < len(queue) {
					currentNode := queue[index]
					index += 1
					marshalledNode, err := json.Marshal(currentNode)
					if err != nil {
						continue
					}

					attuneOrderContentListInfo := []structures.AttuneOrderContentListInfo{}
					err = json.Unmarshal(marshalledNode, &attuneOrderContentListInfo)
					if err != nil {
						continue
					}

					for _, orderInfo := range attuneOrderContentListInfo {
						switch orderInfo.TestType {
						case constants.InvestigationShortHand:
							valueMap[orderInfo.TestCode] = orderInfo.TestValue
							testCodePresentMap[orderInfo.TestCode] = true
							testCodeNameMap[orderInfo.TestCode] = orderInfo.TestName
						case constants.GroupShortHand:
							queue = append(queue, orderInfo.ParameterListInfo)
						}
					}
				}
			}
		}
	}

	return valueMap, testCodePresentMap, testCodeNameMap, doctorDetailsMissingTests
}
