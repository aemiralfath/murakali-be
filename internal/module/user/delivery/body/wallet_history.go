package body

import "time"

type HistoryWalletResponse struct {
	ID          string  `json:"id"`
	From        string  `json:"from"`
	To          string  `json:"to"`
	Amount      float64 `json:"amount"`
	Description string  `json:"description"`
	CreatedAt   string  `json:"created_at"`
}

type GetWalletHistoryRequest struct {
	UserID string `uri:"user_id" binding:"required"`
}

type DetailHistoryWalletResponse struct {
	ID          string                     `json:"id"`
	Transaction *TransactionDetailResponse `json:"transaction"`
	From        string                     `json:"from"`
	To          string                     `json:"to"`
	Amount      float64                    `json:"amount"`
	Description string                     `json:"description"`
	CreatedAt   time.Time                  `json:"created_at"`
}
