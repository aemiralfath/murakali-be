package body

type HistoryWalletResponse struct {
	ID          string `json:"id"`
	From        string `json:"from"`
	To          string `json:"to"`
	Amount      int64  `json:"amount"`
	Description string `json:"description"`
	CreatedAt   string `json:"created_at"`
}

type GetWalletHistoryRequest struct {
	UserID string `uri:"user_id" binding:"required"`
}
