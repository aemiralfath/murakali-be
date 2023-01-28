package body

import (
	"murakali/pkg/httperror"
	"murakali/pkg/response"
	"net/http"
	"strings"
)

type CancelOrderStatus struct {
	OrderID     string `json:"order_id"`
	CancelNotes string `json:"cancel_notes"`
}

func (r *CancelOrderStatus) Validate() (UnprocessableEntity, error) {
	unprocessableEntity := false
	entity := UnprocessableEntity{
		Fields: map[string]string{
			"order_id":     "",
			"cancel_notes": "",
		},
	}

	r.OrderID = strings.TrimSpace(r.OrderID)
	if r.OrderID == "" {
		unprocessableEntity = true
		entity.Fields["order_id"] = FieldCannotBeEmptyMessage
	}

	r.CancelNotes = strings.TrimSpace(r.CancelNotes)
	if r.CancelNotes == "" {
		unprocessableEntity = true
		entity.Fields["cancel_notes"] = FieldCannotBeEmptyMessage
	}

	if unprocessableEntity {
		return entity, httperror.New(
			http.StatusUnprocessableEntity,
			response.UnprocessableEntityMessage,
		)
	}

	return entity, nil
}
