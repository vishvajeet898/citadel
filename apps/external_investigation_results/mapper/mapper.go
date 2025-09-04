package mapper

import (
	"github.com/Orange-Health/citadel/apps/external_investigation_results/structures"
	commonModels "github.com/Orange-Health/citadel/models"
)

func MapToModel(externalInvestigationResultsData *[]structures.ExternalInvestigateResultUpsertItem) *[]commonModels.ExternalInvestigationResult {
	var investigationResults []commonModels.ExternalInvestigationResult
	for _, investigationResultItem := range *externalInvestigationResultsData {
		investigationResult := commonModels.ExternalInvestigationResult{
			ContactId:                           investigationResultItem.ContactId,
			MasterInvestigationId:               investigationResultItem.MasterInvestigationId,
			MasterInvestigationMethodMappingId:  investigationResultItem.MasterInvestigationMethodMappingId,
			SystemExternalInvestigationResultId: investigationResultItem.SystemExternalInvestigationResultId,
			LoincCode:                           investigationResultItem.LoincCode,
			InvestigationName:                   investigationResultItem.InvestigationName,
			InvestigationValue:                  investigationResultItem.InvestigationValue,
			Uom:                                 investigationResultItem.Uom,
			ReferenceRangeText:                  investigationResultItem.ReferenceRangeText,
			ReportedAt:                          investigationResultItem.ReportedAt,
			IsAbnormal:                          investigationResultItem.IsAbnormal,
			LabName:                             investigationResultItem.LabName,
			SystemExternalReportId:              investigationResultItem.SystemExternalReportId,
			Abnormality:                         investigationResultItem.Abnormality,
		}
		investigationResult.CreatedBy = investigationResultItem.CreatedBy
		investigationResult.UpdatedBy = investigationResultItem.UpdatedBy
		investigationResults = append(investigationResults, investigationResult)
	}
	return &investigationResults
}
