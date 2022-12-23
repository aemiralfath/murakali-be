package model

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type Voucher struct {
	ID                 uuid.UUID    `json:"id" db:"id" binding:"omitempty"`
	ShopID             uuid.UUID    `json:"shop_id" db:"shop_id" binding:"omitempty"`
	Code               string       `json:"code" db:"code" binding:"omitempty"`
	Quota              int          `json:"quota" db:"quota" binding:"omitempty"`
	ActivedDate        time.Time    `json:"actived_date" db:"actived_date" binding:"omitempty"`
	ExpiredDate        time.Time    `json:"expired_date" db:"expired_date" binding:"omitempty"`
	DiscountPercentage *float64     `json:"discount_percentage" db:"discount_percentage" binding:"omitempty"`
	DiscountFixPrice   *float64     `json:"discount_fix_price" db:"discount_fix_price" binding:"omitempty"`
	MinProductPrice    *float64     `json:"min_product_price" db:"min_product_price" binding:"omitempty"`
	MaxDiscountPrice   *float64     `json:"max_discount_price" db:"max_discount_price" binding:"omitempty"`
	MaxQuantity        int          `json:"max_quantity" db:"max_quantity" binding:"omitempty"`
	CreatedAt          time.Time    `json:"created_at" db:"created_at" binding:"omitempty"`
	UpdatedAt          sql.NullTime `json:"updated_at" db:"updated_at" binding:"omitempty"`
	DeletedAt          sql.NullTime `json:"deleted_at" db:"deleted_at" binding:"omitempty"`
}
