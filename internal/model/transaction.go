package model

import (
	"database/sql"

	"github.com/google/uuid"
)

type Transaction struct {
	ID                   uuid.UUID    `json:"id" db:"id" binding:"omitempty"`
	VoucherMarketplaceID *uuid.UUID   `json:"voucher_marketplace_id" db:"voucher_marketplace_id" binding:"omitempty"`
	WalletID             *uuid.UUID   `json:"wallet_id" db:"wallet_id" binding:"omitempty"`
	CardNumber           *string      `json:"card_number" db:"card_number" binding:"omitempty"`
	Invoice              *string      `json:"invoice" db:"invoice" binding:"omitempty"`
	TotalPrice           float64      `json:"total_price" db:"total_price" binding:"omitempty"`
	PaidAt               sql.NullTime `json:"paid_at" db:"paid_at" binding:"omitempty"`
	CanceledAt           sql.NullTime `json:"canceled_at" db:"canceled_at" binding:"omitempty"`
	ExpiredAt            sql.NullTime `json:"expired_at" db:"expired_at" binding:"omitempty"`
}
