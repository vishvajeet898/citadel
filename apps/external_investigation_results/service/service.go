package service

import (
	"github.com/Orange-Health/citadel/apps/external_investigation_results/mapper"
	"github.com/Orange-Health/citadel/apps/external_investigation_results/structures"
	commonStructures "github.com/Orange-Health/citadel/common/structures"
	commonModels "github.com/Orange-Health/citadel/models"
)

type ExternalInvestigationResultServiceInterface interface {
	BulkUpsertInvestigations(externalInvestigationResults *[]structures.ExternalInvestigateResultUpsertItem) (*[]commonModels.ExternalInvestigationResult, *commonStructures.CommonError)
	BulkDeleteInvestigations(systemExternalInvestigationResultIds *[]uint, deletedBy uint) *commonStructures.CommonError
	UpdateContact(sourceContactId, newContactId uint) *commonStructures.CommonError
	FetchInvestigations(
		filters structures.ExternalInvestigationResultsDbFilters,
	) (*[]commonModels.ExternalInvestigationResult, *commonStructures.CommonError)
}

func (s *ExternalInvestigationResultService) BulkUpsertInvestigations(externalInvestigationResultsData *[]structures.ExternalInvestigateResultUpsertItem) (*[]commonModels.ExternalInvestigationResult, *commonStructures.CommonError) {

	investigationResults := mapper.MapToModel(externalInvestigationResultsData)

	err := s.ExtInvResDao.UpsertInvestigations(investigationResults)
	if err != nil {
		return nil, err
	}

	return investigationResults, nil
}

func (s *ExternalInvestigationResultService) BulkDeleteInvestigations(
	systemExternalInvestigationResultIds *[]uint,
	deletedBy uint) *commonStructures.CommonError {
	err := s.ExtInvResDao.BulkDelete(*systemExternalInvestigationResultIds, deletedBy)
	if err != nil {
		return err
	}
	return nil
}

func (s *ExternalInvestigationResultService) UpdateContact(sourceContactId, newContactId uint) *commonStructures.CommonError {
	return s.ExtInvResDao.UpdateContact(sourceContactId, newContactId)
}

func (s *ExternalInvestigationResultService) FetchInvestigations(
	filters structures.ExternalInvestigationResultsDbFilters,
) (*[]commonModels.ExternalInvestigationResult, *commonStructures.CommonError) {
	return s.ExtInvResDao.GetInvestigations(filters)
}
