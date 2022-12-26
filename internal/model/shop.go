package model

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type Shop struct {
	ID           uuid.UUID    `json:"id" db:"id" binding:"omitempty"`
	UserID       uuid.UUID    `json:"user_id" db:"user_id" binding:"omitempty"`
	Name         string       `json:"name" db:"name" binding:"omitempty"`
	TotalProduct int          `json:"total_product" db:"total_product" binding:"omitempty"`
	TotalRating  float64      `json:"total_rating" db:"total_rating" binding:"omitempty"`
	RatingAvg    float64      `json:"rating_avg" db:"rating_avg" binding:"omitempty"`
	CreatedAt    time.Time    `json:"created_at" db:"created_at" binding:"omitempty"`
	UpdatedAt    sql.NullTime `json:"updated_at" db:"updated_at" binding:"omitempty"`
	DeletedAt    sql.NullTime `json:"deleted_at" db:"deleted_at" binding:"omitempty"`
}
