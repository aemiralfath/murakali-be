package model

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type CartItem struct {
	ID              uuid.UUID    `json:"id" db:"id" binding:"omitempty"`
	UserID          uuid.UUID    `json:"user_id" db:"user_id" binding:"omitempty"`
	ProductDetailID uuid.UUID    `json:"product_detail_id" db:"product_detail_id" binding:"omitempty"`
	Quantity        float64      `json:"quantity" db:"quantity" binding:"omitempty"`
	CreatedAt       time.Time    `json:"created_at" db:"created_at" binding:"omitempty"`
	UpdatedAt       sql.NullTime `json:"updated_at" db:"updated_at" binding:"omitempty"`
	DeletedAt       sql.NullTime `json:"deleted_at" db:"deleted_at" binding:"omitempty"`
}
