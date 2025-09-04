package models

type InvestigationResultMetadata struct {
	BaseModel
	InvestigationResultId uint   `gorm:"column:investigation_result_id" json:"investigation_result_id"`
	QcFlag                string `gorm:"column:qc_flag" json:"qc_flag"`
	QcLotNumber           string `gorm:"column:qc_lot_number" json:"qc_lot_number"`
	QcValue               string `gorm:"column:qc_value" json:"qc_value"`
	QcWestGardWarning     string `gorm:"column:qc_west_gard_warning" json:"qc_west_gard_warning"`
	QcStatus              string `gorm:"column:qc_status" json:"qc_status"`
}

func (InvestigationResultMetadata) TableName() string {
	return "investigation_results_metadata"
}
