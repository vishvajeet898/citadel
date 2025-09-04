package models

import "time"

type CoAuthorizedPathologists struct {
	BaseModel
	TaskId                uint                `gorm:"column:task_id;not null" json:"task_id"`
	Task                  Task                `gorm:"foreignKey:TaskId;references:Id"`
	InvestigationResultId uint                `gorm:"column:investigation_result_id;null" json:"investigation_result_id"`
	InvestigationResult   InvestigationResult `gorm:"foreignKey:InvestigationResultId;references:Id"`
	CoAuthorizedBy        uint                `gorm:"column:co_authorized_by;null" json:"co_authorized_by"`
	CoAuthorizedByUser    User                `gorm:"foreignKey:CoAuthorizedBy;references:Id"`
	CoAuthorizedTo        uint                `gorm:"column:co_authorized_to;null" json:"co_authorized_to"`
	CoAuthorizedToUser    User                `gorm:"foreignKey:CoAuthorizedTo;references:Id"`
	CoAuthorizedAt        *time.Time          `gorm:"column:co_authorized_at" json:"co_authorized_at"`
}

func (CoAuthorizedPathologists) TableName() string {
	return "co_authorized_pathologists"
}
