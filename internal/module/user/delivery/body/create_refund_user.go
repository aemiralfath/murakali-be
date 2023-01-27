package body

import (
	"murakali/pkg/httperror"
	"murakali/pkg/response"
	"net/http"
	"strings"
)

type CreateRefundUserRequest struct {
	OrderID        string  `json:"order_id"`
	Reason         string  `json:"reason"`
	Image          *string `json:"image"`
	IsSellerRefund bool
	IsBuyerRefund  bool
}

func (r *CreateRefundUserRequest) Validate() (UnprocessableEntity, error) {
	unprocessableEntity := false
	entity := UnprocessableEntity{
		Fields: map[string]string{
			"order_id": "",
			"reason":   "",
			"image":    "",
		},
	}

	r.OrderID = strings.TrimSpace(r.OrderID)
	if r.OrderID == "" {
		unprocessableEntity = true
		entity.Fields["order_id"] = FieldCannotBeEmptyMessage
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
