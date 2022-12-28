package body

import "time"

type ProductDetailRequest struct {
}

type ProductDetailResponse struct {
	ProductInfo   *ProductInfo     `json:"products_info"`
	PromotionInfo *PromotionInfo   `json:"promotions_info"`
	ProductDetail []*ProductDetail `json:"products_detail"`
}

type ProductInfo struct {
	ProductID     string   `json:"id"`
	SKU           string   `json:"sku"`
	Title         string   `json:"title"`
	Description   string   `json:"description"`
	ViewCount     int64    `json:"view_count"`
	FavoriteCount int64    `json:"favorite_count"`
	UnitSold      float64  `json:"unit_sold"`
	ListedStatus  bool     `json:"listed_status"`
	ThumbnailURL  string   `json:"thumbnail_url"`
	RatingAVG     *float64 `json:"rating_avg"`
	MinPrice      *float64 `json:"min_price"`
	MaxPrice      *float64 `json:"max_price"`
	CategoryName  string   `json:"category_name"`
	CategoryURL   string   `json:"category_url"`
}

type PromotionInfo struct {
	PromotionName               string     `json:"promotion_name"`
	PromotionDiscountPercentage *float64   `json:"promotion_discount_percentage"`
	PromotionDiscountFixPrice   *int       `json:"promotion_discount_fix_price"`
	PromotionMinProductPrice    *float64   `json:"promotion_min_product_price"`
	PromotionMaxDiscountPrice   *float64   `json:"promotion_max_discount_price"`
	PromotionQuota              *int       `json:"promotion_quota"`
	PromotionMaxQuantity        *int64     `json:"promotion_max_quantity"`
	PromotionActiveDate         *time.Time `json:"promotion_active_date"`
	PromotionExpiryDate         *time.Time `json:"promotion_expiry_date"`
}

type ProductDetail struct {
	ProductDetailID string            `json:"id"`
	NormalPrice     *float64          `json:"normal_price"`
	Stock           *int              `json:"stock"`
	Weight          *int              `json:"weight"`
	Size            *int              `json:"size"`
	Hazardous       bool              `json:"hazardous"`
	Condition       string            `json:"condition"`
	BulkPrice       bool              `json:"bulk_price"`
	ProductURL      string            `json:"product_url"`
	Variant         map[string]string `json:"variant"`
}

type VariantDetail struct {
	Name string `json:"name"`
	Type string `json:"type"`
}
