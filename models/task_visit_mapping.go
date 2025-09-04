package models

type TaskVisitMapping struct {
	BaseModel
	TaskId  uint   `gorm:"column:task_id;not null" json:"task_id"`
	Task    Task   `gorm:"foreignKey:TaskId;references:Id"`
	VisitId string `gorm:"column:visit_id;not null;type:varchar(255)" json:"visit_id"`
}

func (TaskVisitMapping) TableName() string {
	return "task_visit_mapping"
}
