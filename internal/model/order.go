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
	ProductDetailID  string   `json:"product_detail_id"`
	ProductID        string   `json:"product_id"`
	ProductTitle     string   `json:"product_title"`
	ProductDetailURL string   `json:"product_detail_url"`
	OrderQuantity    int      `json:"order_quantity"`
	ItemPrice        *float64 `json:"order_item_price"`
	TotalPrice       *float64 `json:"order_total_price"`
}
