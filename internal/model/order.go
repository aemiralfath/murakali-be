package model

import (
	"time"
)

type Order struct {
	OrderID     string         `json:"order_id"`
	OrderStatus int            `json:"order_status"`
	TotalPrice  *float64       `json:"total_price"`
	DeliveryFee *float64       `json:"delivery_fee"`
	ResiNumber  string         `json:"resi_no"`
	ShopID      string         `json:"shop_id"`
	ShopName    string         `json:"shop_name"`
	VoucherCode string         `json:"voucher_code"`
	CreatedAt   time.Time      `json:"created_at" db:"created_at" binding:"omitempty"`
	Detail      []*OrderDetail `json:"detail"`
}
