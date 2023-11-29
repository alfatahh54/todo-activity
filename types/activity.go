package types

import "time"

type ActivityType struct {
	ID        *int       `json:"id,omitempty"`
	Email     string     `json:"email,omitempty" binding:"required,email"`
	Title     string     `json:"title,omitempty" binding:"required"`
	CreatedAt *time.Time `json:"created_at,omitempty"`
	UpdatedAt *time.Time `json:"updated_at,omitempty"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}

type ActivityUpdateType struct {
	Email string `json:"email,omitempty"`
	Title string `json:"title,omitempty"`
}
