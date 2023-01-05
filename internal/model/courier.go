package model

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type Courier struct {
	ID          uuid.UUID    `json:"id" db:"id" binding:"omitempty"`
	Name        string       `json:"name" db:"name" binding:"omitempty"`
	Code        string       `json:"code" db:"code" binding:"omitempty"`
	Service     string       `json:"service" db:"service" binding:"omitempty"`
	Description string       `json:"description" db:"description" binding:"omitempty"`
	CreatedAt   time.Time    `json:"created_at" db:"created_at" binding:"omitempty"`
	UpdatedAt   sql.NullTime `json:"updated_at" db:"updated_at" binding:"omitempty"`
}
