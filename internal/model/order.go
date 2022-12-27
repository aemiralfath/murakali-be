package model

import (
	"time"
)

type Orders struct {
	Order *Order `json:"order"`
}

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

type OrderDetail struct {
	ProductDetailID  string   `json:"id"`
	ProductID        string   `json:"product_id"`
	ProductTitle     string   `json:"title"`
	ProductDetailURL string   `json:"url"`
	OrderQuantity    int      `json:"quantity"`
	ItemPrice        *float64 `json:"item_price"`
	TotalPrice       *float64 `json:"total_price"`
}
