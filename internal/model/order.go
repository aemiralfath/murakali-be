package model

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type Order struct {
	OrderID     string         `json:"order_id"`
	OrderStatus int            `json:"order_status"`
	TotalPrice  *float64       `json:"total_price"`
	DeliveryFee *float64       `json:"delivery_fee"`
	ResiNumber  *string        `json:"resi_no"`
	ShopID      string         `json:"shop_id"`
	ShopName    string         `json:"shop_name"`
	VoucherCode *string        `json:"voucher_code"`
	CreatedAt   time.Time      `json:"created_at" db:"created_at" binding:"omitempty"`
	Detail      []*OrderDetail `json:"detail"`
}

type OrderModel struct {
	ID            uuid.UUID    `json:"id" db:"id" binding:"omitempty"`
	TransactionID uuid.UUID    `json:"transaction_id" db:"transaction_id" binding:"omitempty"`
	ShopID        uuid.UUID    `json:"shop_id" db:"shop_id" binding:"omitempty"`
	UserID        uuid.UUID    `json:"user_id" db:"user_id" binding:"omitempty"`
	CourierID     uuid.UUID    `json:"courier_id" db:"courier_id" binding:"omitempty"`
	VoucherShopID *uuid.UUID   `json:"voucher_shop_id" db:"courier_id" binding:"omitempty"`
	OrderStatusID int          `json:"order_status_id" db:"order_status_id" binding:"omitempty"`
	TotalPrice    float64      `json:"total_price" db:"total_price" binding:"omitempty"`
	DeliveryFee   float64      `json:"delivery_fee" db:"delivery_fee" binding:"omitempty"`
	ResiNo        *string      `json:"resi_no" db:"resi_no" binding:"omitempty"`
	CreatedAt     time.Time    `json:"created_at" db:"created_at" binding:"omitempty"`
	ArrivedAt     sql.NullTime `json:"arrived_at" db:"arrived_at" binding:"omitempty"`
}
