package body

import (
	"database/sql"
	"murakali/internal/model"
	"murakali/pkg/httperror"
	"murakali/pkg/response"
	"net/http"
	"strings"

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

type GetTransactionByIDResponse struct {
	ID                 uuid.UUID      `json:"id"`
	VoucherMarketplace *model.Voucher `json:"voucher_marketplace"`
	WalletID           *uuid.UUID     `json:"wallet_id"`
	CardNumber         *string        `json:"card_number"`
	Invoice            *string        `json:"invoice"`
	TotalPrice         float64        `json:"total_price"`
	ExpiredAt          sql.NullTime   `json:"expired_at"`
	Orders             []*model.Order `json:"orders"`
}

type GetTransactionByIDRequest struct {
	TransactionID string `uri:"transaction_id" binding:"required"`
}

type ChangeTransactionPaymentMethodReq struct {
	TransactionID string `json:"transaction_id"`
	CardNumber    string `json:"card_number"`
}

func (r *ChangeTransactionPaymentMethodReq) Validate() (UnprocessableEntity, error) {
	unprocessableEntity := false
	entity := UnprocessableEntity{
		Fields: map[string]string{
			"transaction_id": "",
			"card_number":    "",
		},
	}

	r.TransactionID = strings.TrimSpace(r.TransactionID)
	r.CardNumber = strings.TrimSpace(r.CardNumber)

	if r.TransactionID == "" {
		unprocessableEntity = true
		entity.Fields["transaction_id"] = FieldCannotBeEmptyMessage
	}

	if r.CardNumber == "" {
		unprocessableEntity = true
		entity.Fields["card_number"] = FieldCannotBeEmptyMessage
	}

	if unprocessableEntity {
		return entity, httperror.New(
			http.StatusUnprocessableEntity,
			response.UnprocessableEntityMessage,
		)
	}

	return entity, nil
}
