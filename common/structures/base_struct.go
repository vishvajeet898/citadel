package structures

import "time"

type BaseStruct struct {
	ID        uint       `json:"id"`
	CreatedAt *time.Time `json:"created_at,omitempty"`
	UpdatedAt *time.Time `json:"updated_at,omitempty"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
	CreatedBy uint       `json:"created_by,omitempty"`
	UpdatedBy uint       `json:"updated_by,omitempty"`
	DeletedBy uint       `json:"deleted_by,omitempty"`
}
