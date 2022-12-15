package model

import (
	"github.com/google/uuid"
	"time"
)

type EmailHistory struct {
	ID        uuid.UUID `json:"id" db:"id" binding:"omitempty"`
	Email     string    `json:"email" db:"email" binding:"omitempty"`
	CreatedAt time.Time `json:"created_at" db:"created_at" binding:"omitempty"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at" binding:"omitempty"`
}
