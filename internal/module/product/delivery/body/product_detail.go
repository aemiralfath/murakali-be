package body

import "time"

type ProductDetailRequest struct {
}

type ProductDetailResponse struct {
	ProductInfo   *ProductInfo     `json:"products_info"`
	ProductDetail []*ProductDetail `json:"products_detail"`
}

type ProductInfo struct {
	ProductID          string     `json:"id"`
	SKU                string     `json:"sku"`
	Title              string     `json:"title"`
	Description        string     `json:"description"`
	ViewCount          int64      `json:"view_count"`
	FavoriteCount      int64      `json:"favorite_count"`
	UnitSold           float64    `json:"unit_sold"`
	ListedStatus       bool       `json:"listed_status"`
	ThumbnailURL       string     `json:"thumbnail_url"`
	RatingAVG          *float64   `json:"rating_avg"`
	MinPrice           *float64   `json:"min_price"`
	MaxPrice           *float64   `json:"max_price"`
	PromotionName      string     `json:"promotion_name"`
	DiscountPercentage *float64   `json:"discount_percentage"`
	DiscountFixPrice   *int       `json:"discount_fix_price"`
	MinProductPrice    *float64   `json:"min_product_price"`
	MaxDiscountPrice   *float64   `json:"max_discount_price"`
	Quota              *int       `json:"quota"`
	MaxQuantity        *int64     `json:"max_quantity"`
	ActiveDate         *time.Time `json:"active_date"`
	ExpiryDate         *time.Time `json:"expiry_date"`
	ParentID           string     `json:"parent_id"`
	CategoryName       string     `json:"category_name"`
	CategoryURL        string     `json:"category_url"`
}

type ProductDetail struct {
	ProductDetailID string           `json:"id"`
	Price           *float64         `json:"price"`
	Stock           *int             `json:"stock"`
	Weight          *int             `json:"weight"`
	Size            *int             `json:"size"`
	Hazardous       bool             `json:"hazardous"`
	Condition       string           `json:"condition"`
	BulkPrice       bool             `json:"bulk_price"`
	ProductURL      string           `json:"product_url"`
	Variant         []*VariantDetail `json:"variant"`
}

type VariantDetail struct {
	Name string `json:"name"`
	Type string `json:"type"`
}
