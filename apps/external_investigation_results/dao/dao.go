package dao

import (
	"fmt"
	"net/http"
	"time"

	"github.com/Orange-Health/citadel/apps/external_investigation_results/structures"
	commonConstants "github.com/Orange-Health/citadel/common/constants"
	commonStructures "github.com/Orange-Health/citadel/common/structures"
	commonUtils "github.com/Orange-Health/citadel/common/utils"
	commonModels "github.com/Orange-Health/citadel/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type DataLayer interface {
	UpsertInvestigations(investigations *[]commonModels.ExternalInvestigationResult) *commonStructures.CommonError
	BulkDelete(systemExternalInvestigationIds []uint, deletedBy uint) *commonStructures.CommonError
	UpdateContact(sourceContactId, newContactId uint) *commonStructures.CommonError
	GetInvestigations(
		filters structures.ExternalInvestigationResultsDbFilters,
	) (*[]commonModels.ExternalInvestigationResult, *commonStructures.CommonError)
}

func (dao *ExternalInvestigationResultDao) UpsertInvestigations(
	investigations *[]commonModels.ExternalInvestigationResult) *commonStructures.CommonError {
	err := dao.Db.Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "system_external_investigation_result_id"}},
		DoUpdates: clause.AssignmentColumns([]string{
			"contact_id",
			"master_investigation_id",
			"master_investigation_method_mapping_id",
			"loinc_code",
			"investigation_name",
			"investigation_value",
			"uom",
			"reference_range_text",
			"reported_at",
			"deleted_at",
			"deleted_by",
			"updated_by",
			"lab_name",
			"system_external_report_id",
			"is_abnormal",
			"abnormality",
		}),
	}).Create(investigations).Error
	if err != nil {
		return commonUtils.HandleORMError(err)
	}
	return nil
}

func (dao *ExternalInvestigationResultDao) BulkDelete(systemExternalInvestigationIds []uint, deletedBy uint) *commonStructures.CommonError {
	now := time.Now()
	deletedAt := gorm.DeletedAt{Time: now, Valid: true}
	updates := map[string]interface{}{
		"deleted_at": deletedAt,
		"deleted_by": deletedBy,
	}
	res := dao.Db.Model(commonModels.ExternalInvestigationResult{}).
		Where("system_external_investigation_result_id in ?", systemExternalInvestigationIds).
		Updates(updates)
	err := res.Error
	rowsAffected := res.RowsAffected
	if err != nil {
		return commonUtils.HandleORMError(err)
	}
	if rowsAffected == 0 {
		return &commonStructures.CommonError{
			StatusCode: http.StatusNotFound,
			Message:    fmt.Sprintf(commonConstants.EXT_INV_NOT_FOUND_FOR_SYSTEM_EXT_INV_IDS, systemExternalInvestigationIds),
		}
	}
	return nil
}

func (dao *ExternalInvestigationResultDao) UpdateContact(
	sourceContactId,
	newContactId uint,
) *commonStructures.CommonError {
	updates := map[string]interface{}{
		"contact_id": newContactId,
	}
	err := dao.Db.Model(commonModels.ExternalInvestigationResult{}).
		Where("contact_id = ?", sourceContactId).
		Updates(updates).Error
	if err != nil {
		return commonUtils.HandleORMError(err)
	}
	return nil
}

func (dao *ExternalInvestigationResultDao) GetInvestigations(
	filters structures.ExternalInvestigationResultsDbFilters,
) (*[]commonModels.ExternalInvestigationResult, *commonStructures.CommonError) {
	var investigations []commonModels.ExternalInvestigationResult
	db := dao.Db

	// Patient Filters
	db = db.Where("contact_id = ?", filters.ContactId)

	// Investigation Filters
	if filters.LoincCode != "" && filters.MasterInvestigationMethodMappingId != 0 {
		db = db.Where("loinc_code = ? OR master_investigation_method_mapping_id = ?", filters.LoincCode, filters.MasterInvestigationMethodMappingId)
	} else if filters.LoincCode != "" {
		db = db.Where("loinc_code = ?", filters.LoincCode)
	} else if filters.MasterInvestigationMethodMappingId != 0 {
		db = db.Where("master_investigation_method_mapping_id = ?", filters.MasterInvestigationMethodMappingId)
	}

	err := db.Order("reported_at desc").Limit(int(filters.Limit)).Offset(int(filters.Offset)).Find(&investigations).Error

	if err != nil {
		return nil, commonUtils.HandleORMError(err)
	}
	return &investigations, nil
}
