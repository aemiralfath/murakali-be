package body

import (
	"database/sql"
	"time"
)

type PromotionDetailSeller struct {
	PromotionID             string       `json:"promotion_id"`
	PromotionName           string       `json:"promotion_name"`
	ProductID               string       `json:"product_id"`
	ProductName             string       `json:"product_name"`
	MinPrice                float64      `json:"min_price"`
	MaxPrice                float64      `json:"max_price"`
	ProductThumbnailURL     string       `json:"product_thumbnail_url"`
	DiscountPercentage      *float64     `json:"discount_percentage"`
	DiscountFixPrice        *float64     `json:"discount_fix_price"`
	MinProductPrice         *float64     `json:"min_product_price"`
	MaxDiscountPrice        *float64     `json:"max_discount_price"`
	Quota                   int          `json:"quota"`
	MaxQuantity             int          `json:"max_quantity"`
	ProductSubMinPrice      float64      `json:"product_sub_min_price"`
	ProductSubMaxPrice      float64      `json:"product_sub_max_price"`
	ProductMinDiscountPrice float64      `json:"product_min_discount_price"`
	ProductMaxDiscountPrice float64      `json:"product_max_discount_price"`
	ActivedDate             time.Time    `json:"actived_date"`
	ExpiredDate             time.Time    `json:"expired_date"`
	CreatedAt               time.Time    `json:"created_at"`
	UpdatedAt               sql.NullTime `json:"updated_at"`
	DeletedAt               sql.NullTime `json:"deleted_at"`
}
