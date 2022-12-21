package model

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type Category struct {
	ID        uuid.UUID    `json:"id" db:"id" binding:"omitempty"`
	ParentID  uuid.UUID    `json:"parent_id" db:"parent_id" binding:"omitempty"`
	Name      string       `json:"name" db:"name" binding:"omitempty"`
	PhotoURL  string       `json:"photo_url" db:"photo_url" binding:"omitempty"`
	CreatedAt time.Time    `json:"created_at" db:"created_at" binding:"omitempty"`
	UpdatedAt sql.NullTime `json:"updated_at" db:"updated_at" binding:"omitempty"`
	DeletedAt sql.NullTime `json:"deleted_at" db:"deleted_at" binding:"omitempty"`
}
