package model

import (
	"github.com/google/uuid"
)

type OrderItem struct {
	ID              uuid.UUID `json:"id" db:"id" binding:"omitempty"`
	OrderID         uuid.UUID `json:"order_id" db:"order_id" binding:"omitempty"`
	ProductDetailID uuid.UUID `json:"product_detail_id" db:"product_detail_id" binding:"omitempty"`
	Quantity        int       `json:"quantity" db:"quantity" binding:"omitempty"`
	ItemPrice       float64   `json:"item_price" db:"item_price" binding:"omitempty"`
	TotalPrice      float64   `json:"total_price" db:"total_price" binding:"omitempty"`
}
