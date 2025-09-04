package models

type Remark struct {
	BaseModel
	InvestigationResultId uint                `gorm:"column:investigation_result_id;not null" json:"investigation_result_id"`
	InvestigationResult   InvestigationResult `gorm:"foreignKey:InvestigationResultId;references:Id"`
	Description           string              `gorm:"column:description;not null;type:text" json:"description"`
	RemarkType            string              `gorm:"column:remark_type;not null;type:varchar(50)" json:"remark_type"`
	RemarkBy              uint                `gorm:"column:remark_by;not null;type:varchar(50)" json:"remark_by"`
}

func (Remark) TableName() string {
	return "remarks"
}
