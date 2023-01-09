package model

import (
	"github.com/google/uuid"
	"time"
)

type WalletHistory struct {
	ID            uuid.UUID `json:"id" db:"id" binding:"omitempty"`
	TransactionID uuid.UUID `json:"transaction_id" db:"transaction_id" binding:"omitempty"`
	WalletID      uuid.UUID `json:"wallet_id" db:"wallet_id" binding:"omitempty"`
	From          string    `json:"from" db:"from" binding:"omitempty"`
	To            string    `json:"to" db:"to" binding:"omitempty"`
	Description   string    `json:"description" db:"description" binding:"omitempty"`
	Amount        float64   `json:"amount" db:"amount" binding:"omitempty"`
	CreatedAt     time.Time `json:"created_at" db:"created_at" binding:"omitempty"`
}
