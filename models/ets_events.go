package models

type EtsEvent struct {
	BaseModel
	TestID   string `json:"test_id"`
	IsActive bool   `json:"is_active"`
}

func (EtsEvent) TableName() string {
	return "ets_events"
}
