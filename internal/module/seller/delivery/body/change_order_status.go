package body

import (
	"murakali/pkg/httperror"
	"murakali/pkg/response"
	"net/http"
	"strings"
)

type ChangeOrderStatusRequest struct {
	OrderID       string `json:"order_id"`
	OrderStatusID string `json:"order_status_id"`
}

func (r *ChangeOrderStatusRequest) Validate() (UnprocessableEntity, error) {
	unprocessableEntity := false
	entity := UnprocessableEntity{
		Fields: map[string]string{
			"order_id":        "",
			"order_status_id": "",
		},
	}

	r.OrderID = strings.TrimSpace(r.OrderID)
	if r.OrderID == "" {
		unprocessableEntity = true
		entity.Fields["order_id"] = FieldCannotBeEmptyMessage
	}

	r.OrderStatusID = strings.TrimSpace(r.OrderStatusID)
	if r.OrderStatusID == "" {
		unprocessableEntity = true
		entity.Fields["order_status_id"] = FieldCannotBeEmptyMessage
	}

	if unprocessableEntity {
		return entity, httperror.New(
			http.StatusUnprocessableEntity,
			response.UnprocessableEntityMessage,
		)
	}

	return entity, nil
}
