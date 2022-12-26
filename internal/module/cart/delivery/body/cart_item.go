package body

import "github.com/google/uuid"

type CartItemsResponse struct {
	ID      uuid.UUID          `json:"id" db:"id"`
	Shop    *ShopResponse      `json:"shop"`
	Product []*ProductResponse `json:"product"`
}

type ShopResponse struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

type ProductResponse struct {
	ID           string            `json:"id"`
	Title        string            `json:"title"`
	ThumbnailURL string            `json:"thumbnail_url"`
	ProductPrice float64           `json:"product_price"`
	ProductStock float64           `json:"product_stock"`
	Quantity     float64           `json:"quantity"`
	Variant      map[string]string `json:"variant"`
	Promo        *PromoResponse    `json:"promo"`
}

type PromoResponse struct {
	DiscountPercentage *float64 `json:"discount_percentage" db:"discount_percentage"`
	DiscountFixPrice   *float64 `json:"discount_fix_price" db:"discount_fix_price"`
	MinProductPrice    *float64 `json:"min_product_price" db:"min_product_price"`
	MaxDiscountPrice   *float64 `json:"max_discount_price" db:"max_discount_price"`
	ResultDiscount     float64  `json:"result_discount" db:"result_discount"`
	SubPrice           float64  `json:"sub_price" db:"sub_price"`
}
