package body

import "murakali/internal/model"

type GetRefundThreadResponse struct {
	RefundData    *model.Refund         `json:"refund_data"`
	RefundThreads []*model.RefundThread `json:"refund_threads"`
}