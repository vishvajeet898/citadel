package models

import (
	"time"

	"gorm.io/gorm"
)

type BaseModel struct {
	Id        uint            `gorm:"primaryKey;autoIncrement:true" json:"id"`
	CreatedAt *time.Time      `gorm:"column:created_at;autoCreateTime;not null" json:"created_at"`
	UpdatedAt *time.Time      `gorm:"column:updated_at;autoUpdateTime;not null" json:"updated_at"`
	DeletedAt *gorm.DeletedAt `gorm:"column:deleted_at;index" json:"deleted_at"`
	CreatedBy uint            `gorm:"column:created_by;not null" json:"created_by"`
	UpdatedBy uint            `gorm:"column:updated_by;not null" json:"updated_by"`
	DeletedBy uint            `gorm:"column:deleted_by" json:"deleted_by"`
}
