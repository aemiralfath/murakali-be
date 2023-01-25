package model

type Discount struct {
	DiscountPercentage *float64 `json:"discount_percentage" db:"discount_percentage" binding:"omitempty"`
	DiscountFixPrice   *float64 `json:"discount_fix_price" db:"discount_fix_price" binding:"omitempty"`
	MinProductPrice    *float64 `json:"min_product_price" db:"min_product_price" binding:"omitempty"`
	MaxDiscountPrice   *float64 `json:"max_discount_price" db:"max_discount_price" binding:"omitempty"`
}
