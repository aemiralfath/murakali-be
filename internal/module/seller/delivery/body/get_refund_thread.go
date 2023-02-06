package body

import (
	"murakali/internal/model"
	"time"
)

type GetRefundThreadResponse struct {
	UserName      string        `json:"user_name"`
	PhotoURL      *string       `json:"photo_url"`
	RefundData    *model.Refund `json:"refund_data"`
	RefundThreads []*RThread    `json:"refund_threads"`
}

type RThread struct {
	ID        string    `json:"id"`
	RefundID  string    `json:"refund_id" `
	UserID    string    `json:"user_id"`
	UserName  string    `json:"user_name"`
	ShopName  *string   `json:"shop_name"`
	PhotoURL  *string   `json:"photo_url"`
	IsSeller  bool      `json:"is_seller"`
	IsBuyer   bool      `json:"is_buyer"`
	Text      string    `json:"text"`
	CreatedAt time.Time `json:"created_at"`
}
