package models

type InvestigationData struct {
	BaseModel
	InvestigationResultId uint   `gorm:"column:investigation_result_id;not null" json:"investigation_result_id"`
	Data                  string `gorm:"column:data" json:"data"`
	DataType              string `gorm:"column:data_type;not null;type:varchar(100)" json:"data_type"`
}

func (InvestigationData) TableName() string {
	return "investigation_data"
}
