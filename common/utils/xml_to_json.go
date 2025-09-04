package utils

import (
	"context"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/Orange-Health/citadel/common/constants"
	"github.com/Orange-Health/citadel/common/structures"
)

func ParseXMLToJSON(ctx context.Context, data string) (string, error) {
	if !strings.HasPrefix(data, constants.CultureInvValueTextPrefix) {
		AddLog(ctx, constants.DEBUG_LEVEL, GetCurrentFunctionName(), nil, errors.New("data does not start with expected struct"))
		return data, nil
	}

	xmlBytes := []byte(data)

	var investigationResults structures.AttuneInvestigationResults
	err := xml.Unmarshal(xmlBytes, &investigationResults)
	if err != nil {
		AddLog(ctx, constants.DEBUG_LEVEL, GetCurrentFunctionName(), nil, err)
		return data, err
	}

	if len(investigationResults.InvestigationDetails) == 0 {
		AddLog(ctx, constants.DEBUG_LEVEL, GetCurrentFunctionName(), nil, errors.New("no investigation details found"))
		return data, nil
	}

	jsonInvestigationResults := structures.TransformedInvestigationResults{
		ReportStatus:       investigationResults.InvestigationDetails[0].ReportStatus,
		ClinicalHistory:    investigationResults.InvestigationDetails[0].ClinicalHistory,
		Gross:              investigationResults.InvestigationDetails[0].Gross,
		SampleType:         investigationResults.InvestigationDetails[0].SampleName,
		IncubationPeriod:   investigationResults.InvestigationDetails[0].InCubPeriod,
		CultureStainResult: investigationResults.InvestigationDetails[0].CultureStainType,
		CultureReport:      investigationResults.InvestigationDetails[0].CultureReport,
		ColonyCount:        investigationResults.InvestigationDetails[0].Colonycount,
		ResistanceDetected: investigationResults.InvestigationDetails[0].ResistanceDetected,
		StainResultDetails: []structures.StainResultDetail{},
		GrowthData:         []structures.MicrorganismGrowthData{},
	}

	for _, stainDetail := range investigationResults.InvestigationDetails[0].StainDetails[0].Stain {
		jsonInvestigationResults.StainResultDetails = append(jsonInvestigationResults.StainResultDetails,
			structures.StainResultDetail{
				Type:   stainDetail.Type,
				Result: stainDetail.Result,
			})
	}

	microOrganismDataMap := map[string]structures.MicrorganismGrowthData{}

	for _, organDetail := range investigationResults.InvestigationDetails[0].OrganDetails[0].Organ {
		if _, ok := microOrganismDataMap[organDetail.Name]; !ok {
			microOrganismDataMap[organDetail.Name] = structures.MicrorganismGrowthData{
				Microorganism:      organDetail.Name,
				ColonyCount:        organDetail.ColonyCount,
				Family:             organDetail.Family,
				FamilyDisplayOrder: nil,
				Sensitivity:        []structures.Sensitivity{},
			}
		}
		micValue, _ := strconv.Atoi(organDetail.Zone)
		displayOrder, _ := strconv.Atoi(organDetail.NameSeq)
		sensitivity := structures.Sensitivity{
			Antimicrobial: organDetail.DrugName,
			Sensitivity:   organDetail.Sensitivity,
			Level:         organDetail.Level,
			MicValue:      micValue,
			DisplayOrder:  displayOrder,
		}
		temp := microOrganismDataMap[organDetail.Name]
		temp.Sensitivity = append(microOrganismDataMap[organDetail.Name].Sensitivity, sensitivity)
		microOrganismDataMap[organDetail.Name] = temp
	}
	growthData := []structures.MicrorganismGrowthData{}
	for _, value := range microOrganismDataMap {
		growthData = append(growthData, value)
	}

	jsonInvestigationResults.GrowthData = growthData

	jsonBytes, err := json.Marshal(jsonInvestigationResults)
	if err != nil {
		AddLog(ctx, constants.ERROR_LEVEL, GetCurrentFunctionName(), map[string]interface{}{
			"investigation_details": fmt.Sprintf("%+v", investigationResults),
		}, err)
		return data, err
	}

	return string(jsonBytes), nil
}
