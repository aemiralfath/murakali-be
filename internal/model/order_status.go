package model

import (
	"database/sql"
	"time"
)

type OrderStatus struct {
	ID        int          `json:"id" db:"id" binding:"omitempty"`
	Name      string       `json:"name" db:"name" binding:"omitempty"`
	CreatedAt time.Time    `json:"created_at" db:"created_at" binding:"omitempty"`
	UpdatedAt sql.NullTime `json:"updated_at" db:"updated_at" binding:"omitempty"`
	DeletedAt sql.NullTime `json:"deleted_at" db:"deleted_at" binding:"omitempty"`
}
