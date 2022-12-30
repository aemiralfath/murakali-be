package model

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type Wallet struct {
	ID           uuid.UUID    `json:"id" db:"id" binding:"omitempty"`
	UserID       uuid.UUID    `json:"user_id" db:"user_id" binding:"omitempty"`
	Balance      float64      `json:"balance" db:"balance" binding:"omitempty"`
	PIN          string       `json:"pin" db:"pin" binding:"omitempty"`
	AttemptCount int          `json:"attempt_count" db:"attempt_count" binding:"omitempty"`
	AttemptAt    sql.NullTime `json:"attempt_at" db:"attempt_at" binding:"omitempty"`
	UnlockedAt   sql.NullTime `json:"unlocked_at" db:"unlocked_at" binding:"omitempty"`
	CreatedAt    time.Time    `json:"created_at" db:"created_at" binding:"omitempty"`
	UpdatedAt    sql.NullTime `json:"updated_at" db:"updated_at" binding:"omitempty"`
	DeletedAt    sql.NullTime `json:"deleted_at" db:"deleted_at" binding:"omitempty"`
	ActiveDate   sql.NullTime `json:"active_date" db:"active_date" binding:"omitempty"`
}
