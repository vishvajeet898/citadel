package models

type User struct {
	BaseModel
	UserName     string `gorm:"column:user_name;not null;type:varchar(255)" json:"name"`
	SystemUserId string `gorm:"column:system_user_id;default:null;type:varchar(50)" json:"user_id"`
	UserType     string `gorm:"column:user_type;not null;type:varchar(100)" json:"user_type"`
	Email        string `gorm:"column:email;default:null;type:varchar(255)" json:"email"`
	AttuneUserId string `gorm:"column:attune_user_id;default:null;type:varchar(50)" json:"attune_user_id"`
	AgentId      string `gorm:"column:agent_id;default:null;type:varchar(50)" json:"agent_id"`
}

func (User) TableName() string {
	return "users"
}
