package model

import (
	"database/sql"

	"github.com/google/uuid"
)

type Refund struct {
	ID             uuid.UUID    `json:"id" db:"id" binding:"omitempty"`
	OrderID        uuid.UUID    `json:"order_id" db:"order_id" binding:"omitempty"`
	IsSellerRefund *bool        `json:"is_seller_refund" db:"is_seller_refund" binding:"omitempty"`
	IsBuyerRefund  *bool        `json:"is_buyer_refund" db:"is_buyer_refund" binding:"omitempty"`
	Reason         string       `json:"reason" db:"reason" binding:"omitempty"`
	Image          *string      `json:"image" db:"image" binding:"omitempty"`
	AcceptedAt     sql.NullTime `json:"accepted_at" db:"accepted_at" binding:"omitempty"`
	RejectedAt     sql.NullTime `json:"rejected_at" db:"rejected_at" binding:"omitempty"`
	RefundedAt     sql.NullTime `json:"refunded_at" db:"refunded_at" binding:"omitempty"`
}

type RefundOrder struct {
	ID             uuid.UUID    `json:"id" db:"id" binding:"omitempty"`
	OrderID        uuid.UUID    `json:"order_id" db:"order_id" binding:"omitempty"`
	IsSellerRefund *bool        `json:"is_seller_refund" db:"is_seller_refund" binding:"omitempty"`
	IsBuyerRefund  *bool        `json:"is_buyer_refund" db:"is_buyer_refund" binding:"omitempty"`
	Reason         string       `json:"reason" db:"reason" binding:"omitempty"`
	Image          *string      `json:"image" db:"image" binding:"omitempty"`
	AcceptedAt     sql.NullTime `json:"accepted_at" db:"accepted_at" binding:"omitempty"`
	RejectedAt     sql.NullTime `json:"rejected_at" db:"rejected_at" binding:"omitempty"`
	RefundedAt     sql.NullTime `json:"refunded_at" db:"refunded_at" binding:"omitempty"`
	Order          *OrderModel  `json:"order" db:"order" binding:"omitempty"`
}
