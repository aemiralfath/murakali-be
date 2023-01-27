package model

import (
	"time"

	"github.com/google/uuid"
)

type RefundThread struct {
	ID        uuid.UUID    `json:"id" db:"id" binding:"omitempty"`
	RefundID  uuid.UUID    `json:"refund_id" db:"refund_id" binding:"omitempty"`
	UserID    uuid.UUID    `json:"user_id" db:"user_id" binding:"omitempty"`
	IsSeller  *bool        `json:"is_seller" db:"is_seller" binding:"omitempty"`
	IsBuyer   *bool        `json:"is_buyer" db:"is_buyer" binding:"omitempty"`
	Text      string       `json:"text" db:"text" binding:"omitempty"`
	CreatedAt time.Time    `json:"created_at" db:"created_at" binding:"omitempty"`
}
