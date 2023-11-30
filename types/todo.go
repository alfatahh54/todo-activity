package types

import "time"

type TodoType struct {
	ID              *int       `json:"id,omitempty"`
	ActivityGroupID int        `json:"activity_group_id,omitempty" binding:"required"`
	Title           string     `json:"title,omitempty" binding:"required"`
	IsActive        bool       `json:"is_active,omitempty"`
	Priority        string     `json:"priority,omitempty" binding:"required"`
	CreatedAt       *time.Time `json:"created_at,omitempty"`
	UpdatedAt       *time.Time `json:"updated_at,omitempty"`
	DeletedAt       *time.Time `json:"deleted_at,omitempty"`
}

type TodoUpdateType struct {
	Title    string `json:"title,omitempty"`
	IsActive *bool  `json:"is_active,omitempty"`
	Priority string `json:"priority,omitempty"`
}
