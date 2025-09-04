package models

type Template struct {
	BaseModel
	Title        string `gorm:"column:title;type:varchar(255)" json:"title"`
	TemplateType string `gorm:"column:template_type;not null;type:varchar(50)" json:"template_type"`
	Description  string `gorm:"column:description;not null;type:text" json:"description"`
	DisplayOrder int    `gorm:"column:display_order;not null;type:int" json:"display_order"`
}

func (Template) TableName() string {
	return "templates"
}
