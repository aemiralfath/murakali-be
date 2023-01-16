package body

import (
	"database/sql"
	"murakali/internal/model"

	"github.com/google/uuid"
)

type GetTransactionByUserIDResponse struct {
	ID                 uuid.UUID           `json:"id"`
	VoucherMarketplace *model.Voucher      `json:"voucher_marketplace"`
	WalletID           *uuid.UUID          `json:"wallet_id"`
	CardNumber         *string             `json:"card_number"`
	Invoice            *string             `json:"invoice"`
	TotalPrice         float64             `json:"total_price"`
	ExpiredAt          sql.NullTime        `json:"expired_at"`
	Orders             []*model.OrderModel `json:"orders"`
}
