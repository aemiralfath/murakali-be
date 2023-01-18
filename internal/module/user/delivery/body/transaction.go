package body

import (
	"database/sql"
	"murakali/internal/model"

	"github.com/google/uuid"
)

type TransactionDetailResponse struct {
	ID                   uuid.UUID      `json:"id"`
	VoucherMarketplaceID *uuid.UUID     `json:"voucher_marketplace_id"`
	WalletID             *uuid.UUID     `json:"wallet_id"`
	CardNumber           *string        `json:"card_number"`
	Invoice              *string        `json:"invoice"`
	TotalPrice           float64        `json:"total_price"`
	PaidAt               sql.NullTime   `json:"paid_at"`
	CanceledAt           sql.NullTime   `json:"canceled_at"`
	ExpiredAt            sql.NullTime   `json:"expired_at"`
	Orders               []*model.Order `json:"orders"`
}
