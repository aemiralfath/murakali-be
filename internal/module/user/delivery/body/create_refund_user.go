package body

import (
	"murakali/pkg/httperror"
	"murakali/pkg/response"
	"net/http"
	"strings"
)

type CreateRefundUserRequest struct {
	OrderID        string  `json:"order_id"`
	IsSellerRefund *bool   `json:"is_seller_refund"`
	IsBuyerRefund  *bool   `json:"is_buyer_refund"`
	Reason         string  `json:"reason"`
	Image          *string `json:"image"`
}

func (r *CreateRefundUserRequest) Validate() (UnprocessableEntity, error) {
	unprocessableEntity := false
	entity := UnprocessableEntity{
		Fields: map[string]string{
			"order_id":         "",
			"is_seller_refund": "",
			"is_buyer_refund":  "",
			"reason":           "",
			"image":            "",
		},
	}

	r.OrderID = strings.TrimSpace(r.OrderID)
	if r.OrderID == "" {
		unprocessableEntity = true
		entity.Fields["order_id"] = FieldCannotBeEmptyMessage
	}

	if r.IsSellerRefund == nil {
		unprocessableEntity = true
		entity.Fields["is_seller_refund"] = FieldCannotBeEmptyMessage
	}

	if r.IsBuyerRefund == nil {
		unprocessableEntity = true
		entity.Fields["is_buyer_refund"] = FieldCannotBeEmptyMessage
	}

	r.Reason = strings.TrimSpace(r.Reason)
	if r.Reason == "" {
		unprocessableEntity = true
		entity.Fields["reason"] = FieldCannotBeEmptyMessage
	}

	if r.Image == nil {
		unprocessableEntity = true
		entity.Fields["image"] = FieldCannotBeEmptyMessage
	}

	if r.Image != nil {
		*r.Image = strings.TrimSpace(*r.Image)
		if *r.Image == "" {
			unprocessableEntity = true
			entity.Fields["image"] = FieldCannotBeEmptyMessage
		}
	}

	if unprocessableEntity {
		return entity, httperror.New(
			http.StatusUnprocessableEntity,
			response.UnprocessableEntityMessage,
		)
	}

	return entity, nil
}
