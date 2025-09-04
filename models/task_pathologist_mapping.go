package models

type TaskPathologistMapping struct {
	BaseModel
	PathologistId uint `gorm:"column:pathologist_id;not null" json:"pathologist_id"`
	Pathologist   User `gorm:"foreignKey:PathologistId;references:Id"`
	TaskId        uint `gorm:"column:task_id;not null" json:"task_id"`
	Task          Task `gorm:"foreignKey:TaskId;references:Id"`
	IsActive      bool `gorm:"column:is_active;not null" json:"is_active"`
}

func (TaskPathologistMapping) TableName() string {
	return "task_pathologist_mapping"
}
