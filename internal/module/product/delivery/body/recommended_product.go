package body

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type RecommendedProductRequest struct {
}

type RecommendedProductResponse struct {
	Limit    int         `json:"limit"`
	Products []*Products `json:"products"`
}

type Products struct {
	ID                        uuid.UUID    `json:"id" db:"id"`
	Title                     string       `json:"title" db:"title"`
	UnitSold                  int64        `json:"unit_sold" db:"unit_sold"`
	RatingAVG                 float64      `json:"rating_avg" db:"rating_avg"`
	ThumbnailURL              string       `json:"thumbnail_url" db:"thumbnail_url"`
	MinPrice                  float64      `json:"min_price" db:"min_price"`
	MaxPrice                  float64      `json:"max_price" db:"max_price"`
	ViewCount                 int64        `json:"view_count" db:"view_count"`
	SubPrice                  float64      `json:"sub_price" db:"sub_price"`
	PromoDiscountPercentage   *float64     `json:"promo_discount_percentage" db:"promo_discount_percentage"`
	PromoDiscountFixPrice     *float64     `json:"promo_discount_fix_price" db:"promo_discount_fix_price"`
	PromoMinProductPrice      *float64     `json:"promo_min_product_price" db:"promo_min_product_price"`
	PromoMaxDiscountPrice     *float64     `json:"promo_max_discount_price" db:"promo_max_discount_price"`
	ResultDiscount            *float64     `json:"result_discount" db:"result_discount"`
	VoucherDiscountPercentage *float64     `json:"voucher_discount_percentage" db:"voucher_discount_percentage"`
	VoucherDiscountFixPrice   *float64     `json:"voucher_discount_fix_price" db:"voucher_discount_fix_price"`
	ShopName                  string       `json:"shop_name" db:"shop_name"`
	CategoryName              string       `json:"category_name" db:"category_name"`
	ShopProvince              string       `json:"province" db:"province"`
	SKU                       *string      `json:"sku" db:"sku"`
	ListedStatus              bool         `json:"listed_status" db:"listed_status"`
	CreatedAt                 time.Time    `json:"created_at" db:"created_at"`
	UpdatedAt                 sql.NullTime `json:"updated_at" db:"updated_at"`
}
