package model

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type ProductDetail struct {
	ID        uuid.UUID    `json:"id" db:"id" binding:"omitempty"`
	ProductID uuid.UUID    `json:"product_id" db:"product_id" binding:"omitempty"`
	Price     float64      `json:"price" db:"price" binding:"omitempty"`
	Stock     float64      `json:"stock" db:"stock" binding:"omitempty"`
	Weight    float64      `json:"weight" db:"weight" binding:"omitempty"`
	Size      float64      `json:"size" db:"size" binding:"omitempty"`
	Hazardous bool         `json:"hazardous" db:"hazardous" binding:"omitempty"`
	Condition string       `json:"condition" db:"condition" binding:"omitempty"`
	BulkPrice bool         `json:"bulk_price" db:"bulk_price" binding:"omitempty"`
	CreatedAt time.Time    `json:"created_at" db:"created_at" binding:"omitempty"`
	UpdatedAt sql.NullTime `json:"updated_at" db:"updated_at" binding:"omitempty"`
	DeletedAt sql.NullTime `json:"deleted_at" db:"deleted_at" binding:"omitempty"`
}
