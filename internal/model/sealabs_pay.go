package model

import (
	"database/sql"
	"github.com/google/uuid"
	"time"
)

type SealabsPay struct {
	CardNumber 	string `json:"card_number" db:"card_number" binding:"omitempty"`
	UserID    	uuid.UUID    `json:"user_id" db:"user_id" binding:"omitempty"`
	Name 		string `json:"name" db:"name" binding:"omitempty"`
	IsDefault 	bool `json:"is_default" db:"is_default" binding:"omitempty"`
	ActiveDate 	time.Time `json:"active_date" db:"active_date" binding:"omitempty"`
	CreatedAt 	time.Time    `json:"created_at" db:"created_at" binding:"omitempty"`
	UpdatedAt 	sql.NullTime `json:"updated_at" db:"updated_at" binding:"omitempty"`
	DeletedAt 	sql.NullTime `json:"deleted_at" db:"deleted_at" binding:"omitempty"`

}
