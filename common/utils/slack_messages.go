package utils

import (
	"fmt"

	"github.com/Orange-Health/citadel/common/constants"
	"github.com/Orange-Health/citadel/common/structures"
)

func GetSlackMessageForLisSyncFailed(patientName, cityCode string, omsOrderId, omsRequestId string,
	attuneTestSampleMap []structures.AttuneTestSampleMap, testDetailsListForAttune []structures.AttuneTestDetails) string {

	testNames, barcodes := []string{}, []string{}

	for _, test := range testDetailsListForAttune {
		testNames = append(testNames, test.TestName)
	}
	for _, testSampleMap := range attuneTestSampleMap {
		barcodes = append(barcodes, testSampleMap.Barcode)
	}

	omsOrderIdString := GetStringOrderIdWithoutStringPart(omsOrderId)
	omsOrderUrl := fmt.Sprintf("%s/request/%s/order/%s", GetOmsBaseDomain(cityCode), omsRequestId, omsOrderIdString)
	return fmt.Sprint(
		":warning: *LIS SYNC FAILED!*\n",
		"*Patient:* "+patientName+"\n",
		"*Tests:* "+ConvertStringSliceToString(testNames)+"\n",
		"*Barcodes:* "+ConvertStringSliceToString(barcodes)+"\n",
		"*Order ID:* "+omsOrderIdString+"\n",
		"*Request ID:* "+omsRequestId+"\n",
		"<", omsOrderUrl, "|Launch OMS>",
	)
}

func GetSlackMessageForLisReportPdfEmpty(patientName, visitId, cityCode string, omsOrderId, omsRequestId string) string {
	omsOrderIdString := GetStringOrderIdWithoutStringPart(omsOrderId)
	omsOrderUrl := fmt.Sprintf("%s/request/%s/order/%s", GetOmsBaseDomain(cityCode), omsRequestId, omsOrderIdString)
	return fmt.Sprint(
		":warning: *REPORT PDF EMPTY!*\n",
		"*Patient:* "+patientName+"\n",
		"*Order ID:* "+omsOrderIdString+"\n",
		"*Request ID:* "+omsRequestId+"\n",
		"*Visit ID:* "+visitId+"\n",
		"*Description:* Report PDF not found in LIS Payload\n",
		"<", omsOrderUrl, "|Launch OMS>",
	)
}

func GetSlackMessageForMissingParametersInCds(patientName, visitId, cityCode string, omsOrderId, omsRequestId string,
	missingTestsInCds []string, testCodeNameMapLis map[string]string) string {

	investigationNames := []string{}
	for _, testCode := range missingTestsInCds {
		if testName, ok := testCodeNameMapLis[testCode]; ok {
			investigationNames = append(investigationNames, testName)
		}
	}
	investigationNamesString := ConvertStringSliceToString(investigationNames)

	omsOrderIdString := GetStringOrderIdWithoutStringPart(omsOrderId)
	omsOrderUrl := fmt.Sprintf("%s/request/%s/order/%s", GetOmsBaseDomain(cityCode), omsRequestId, omsOrderIdString)
	return fmt.Sprint(":warning: *MISSING PARAMETERS IN CDS*\n",
		"*Patient:* "+patientName+"\n",
		"*Order ID:* "+omsOrderIdString+"\n",
		"*Request ID:* "+omsRequestId+"\n",
		"*Visit ID:* "+visitId+"\n",
		"*Investigations:* "+investigationNamesString+"\n",
		"<", omsOrderUrl, "|Launch OMS>")
}

func GetSlackMessageForMissingParametersInLis(patientName, visitId, cityCode string, omsOrderId, omsRequestId string,
	missingTestsInLis []string, testCodeInvestigationMap map[string]structures.Investigation) string {

	investigationNames := []string{}
	for _, testCode := range missingTestsInLis {
		if investigation, ok := testCodeInvestigationMap[testCode]; ok {
			investigationNames = append(investigationNames, investigation.InvestigationName)
		}
	}
	investigationNamesString := ConvertStringSliceToString(investigationNames)

	omsOrderIdString := GetStringOrderIdWithoutStringPart(omsOrderId)
	omsOrderUrl := fmt.Sprintf("%s/request/%s/order/%s", GetOmsBaseDomain(cityCode), omsRequestId, omsOrderIdString)
	return fmt.Sprint(":warning: *MISSING PARAMETERS IN LIS*\n",
		"*Patient:* "+patientName+"\n",
		"*Order ID:* "+omsOrderIdString+"\n",
		"*Request ID:* "+omsRequestId+"\n",
		"*Visit ID:* "+visitId+"\n",
		"*Investigations:* "+investigationNamesString+"\n",
		"<", omsOrderUrl, "|Launch OMS>")
}

func GetSlackMessageForBlankValues(patientName, visitId, cityCode string, omsOrderId, omsRequestId string,
	blankTests []string, testCodeInvestigationMap map[string]structures.Investigation) string {

	investigationNames := []string{}
	for _, testCode := range blankTests {
		if investigation, ok := testCodeInvestigationMap[testCode]; ok {
			investigationNames = append(investigationNames, investigation.InvestigationName)
		}
	}
	investigationNamesString := ConvertStringSliceToString(investigationNames)

	omsOrderIdString := GetStringOrderIdWithoutStringPart(omsOrderId)
	omsOrderUrl := fmt.Sprintf("%s/request/%s/order/%s", GetOmsBaseDomain(cityCode), omsRequestId, omsOrderIdString)
	return fmt.Sprint(":warning: *BLANK VALUES IDENTIFIED*\n",
		"*Patient:* "+patientName+"\n",
		"*Order ID:* "+omsOrderIdString+"\n",
		"*Request ID:* "+omsRequestId+"\n",
		"*Visit ID:* "+visitId+"\n",
		"*Investigations:* "+investigationNamesString+"\n",
		"<", omsOrderUrl, "|Launch OMS>")
}

func GetSlackMessageForIncorrectResultTypes(patientName, visitId, cityCode string, omsOrderId, omsRequestId string,
	incorrectTests []string, testCodeInvestigationMap map[string]structures.Investigation) string {

	investigationNames := []string{}
	for _, testCode := range incorrectTests {
		if investigation, ok := testCodeInvestigationMap[testCode]; ok {
			investigationNames = append(investigationNames, investigation.InvestigationName)
		}
	}
	investigationNamesString := ConvertStringSliceToString(investigationNames)

	omsOrderIdString := GetStringOrderIdWithoutStringPart(omsOrderId)
	omsOrderUrl := fmt.Sprintf("%s/request/%s/order/%s", GetOmsBaseDomain(cityCode), omsRequestId, omsOrderIdString)
	return fmt.Sprint(":warning: *INCORRECT RESULT TYPES IDENTIFIED*\n",
		"*Patient:* "+patientName+"\n",
		"*Order ID:* "+omsOrderIdString+"\n",
		"*Request ID:* "+omsRequestId+"\n",
		"*Visit ID:* "+visitId+"\n",
		"*Investigations:* "+investigationNamesString+"\n",
		"<", omsOrderUrl, "|Launch OMS>")
}

func GetSlackMessageForDoctorDetailsMissing(patientName, visitId, cityCode string,
	omsOrderId, omsRequestId string, missingDoctorDetailsTests []string) string {

	missingDoctorDetailsTestsString := ConvertStringSliceToString(missingDoctorDetailsTests)

	omsOrderIdString := GetStringOrderIdWithoutStringPart(omsOrderId)
	omsOrderUrl := fmt.Sprintf("%s/request/%s/order/%s", GetOmsBaseDomain(cityCode), omsRequestId, omsOrderIdString)
	return fmt.Sprint(":warning: *DOCTOR DETAILS MISSING*\n",
		"*Test approved whose approval doctor details are missing:*"+missingDoctorDetailsTestsString,
		"*Patient:* "+patientName+"\n",
		"*Order ID:* "+omsOrderIdString+"\n",
		"*Request ID:* "+omsRequestId+"\n",
		"*Visit ID:* "+visitId+"\n",
		"<", omsOrderUrl, "|Launch OMS>")
}

func getAlertMessageForFailureReason(failureReason string) string {
	switch failureReason {
	case constants.AUTO_APPROVAL_FAIL_REASON_PAST_RECORD:
		return constants.AA_FAIL_ALERT_PAST_RECORD
	case constants.AUTO_APPROVAL_FAIL_REASON_REF_RANGE:
		return constants.AA_FAIL_ALERT_REF_RANGE
	}
	return ""
}

func GetSlackMessageBlocksForAutoApprovalFailure(attributes structures.AutoApprovalSlackMessageAttributes,
	failureReason string) []map[string]interface{} {
	alert := getAlertMessageForFailureReason(failureReason)
	blocks := []map[string]interface{}{
		{
			"type": "header",
			"text": map[string]interface{}{
				"type":  "plain_text",
				"text":  ":warning: *" + alert + "*",
				"emoji": true,
			},
		},
		{
			"type": "section",
			"text": map[string]interface{}{
				"type": "mrkdwn",
				"text": fmt.Sprintf("*Patient Name:* %s", attributes.PatientName),
			},
		},
		{
			"type": "section",
			"text": map[string]interface{}{
				"type": "mrkdwn",
				"text": fmt.Sprintf("*Order ID:* %s", attributes.OrderId),
			},
		},
		{
			"type": "section",
			"text": map[string]interface{}{
				"type": "mrkdwn",
				"text": fmt.Sprintf("*Request ID:* %s", attributes.RequestID),
			},
		},
		{
			"type": "section",
			"text": map[string]interface{}{
				"type": "mrkdwn",
				"text": fmt.Sprintf("*Visit ID:* %s", attributes.VisitID),
			},
		},
		{
			"type": "section",
			"text": map[string]interface{}{
				"type": "mrkdwn",
				"text": fmt.Sprintf("*City Code:* %s", attributes.CityCode),
			},
		},
		{
			"type": "section",
			"text": map[string]interface{}{
				"type": "mrkdwn",
				"text": fmt.Sprintf("*Investigations:* %s", attributes.InvestigationsName),
			},
		},
		{
			"type": "section",
			"text": map[string]interface{}{
				"type": "mrkdwn",
				"text": fmt.Sprintf("*Test Names:* %s", attributes.TestNames),
			},
		},
		{
			"type": "actions",
			"elements": []map[string]interface{}{
				{
					"type": "button",
					"text": map[string]interface{}{
						"type":  "plain_text",
						"text":  "Info Screen",
						"emoji": true,
					},
					"url":   getSuperlabInfoScreenUrl(attributes.OrderId),
					"style": "primary",
				},
			},
		},
	}
	return blocks
}
