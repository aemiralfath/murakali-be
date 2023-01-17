package body

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type PromotionSellerResponse struct {
	ID                  uuid.UUID    `json:"id" `
	PromotionName       string       `json:"promotion_name"`
	ProductID           uuid.UUID    `json:"product_id"`
	ProductName         string       `json:"product_name"`
	ProductThumbnailURL *string      `json:"product_thumbnail_url"`
	DiscountPercentage  *float64     `json:"discount_percentage"`
	DiscountFixPrice    *float64     `json:"discount_fix_price"`
	MinProductPrice     *float64     `json:"min_product_price"`
	MaxDiscountPrice    *float64     `json:"max_discount_price"`
	Quota               int          `json:"quota"`
	MaxQuantity         int          `json:"max_quantity"`
	ActivedDate         time.Time    `json:"actived_date"`
	ExpiredDate         time.Time    `json:"expired_date"`
	CreatedAt           time.Time    `json:"created_at"`
	UpdatedAt           sql.NullTime `json:"updated_at"`
	DeletedAt           sql.NullTime `json:"deleted_at"`
}
