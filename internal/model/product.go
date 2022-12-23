package model

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type Product struct {
	ID            uuid.UUID    `json:"id" db:"id" binding:"omitempty"`
	CategoryID    uuid.UUID    `json:"category_id" db:"category_id" binding:"omitempty"`
	ShopID        uuid.UUID    `json:"shop_id" db:"shop_id" binding:"omitempty"`
	Sku           string       `json:"sku" db:"sku" binding:"omitempty"`
	Title         string       `json:"title" db:"title" binding:"omitempty"`
	Description   string       `json:"description" db:"description" binding:"omitempty"`
	ViewCount     int64        `json:"view_count" db:"view_count" binding:"omitempty"`
	FavoriteCount int64        `json:"favorite_count" db:"favorite_count" binding:"omitempty"`
	UnitSold      int64        `json:"unit_sold" db:"unit_sold" binding:"omitempty"`
	ListedStatus  bool         `json:"listed_status" db:"listed_status" binding:"omitempty"`
	ThumbnailUrl  string       `json:"thumbnail_url" db:"thumbnail_url" binding:"omitempty"`
	RatingAvg     float64      `json:"rating_avg" db:"rating_avg" binding:"omitempty"`
	MinPrice      float64      `json:"min_price" db:"min_price" binding:"omitempty"`
	MaxPrice      float64      `json:"max_price" db:"max_price" binding:"omitempty"`
	CreatedAt     time.Time    `json:"created_at" db:"created_at" binding:"omitempty"`
	UpdatedAt     sql.NullTime `json:"updated_at" db:"updated_at" binding:"omitempty"`
	DeletedAt     sql.NullTime `json:"deleted_at" db:"deleted_at" binding:"omitempty"`
}
