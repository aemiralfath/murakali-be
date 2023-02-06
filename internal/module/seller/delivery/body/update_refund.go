package body

import (
	"murakali/pkg/httperror"
	"murakali/pkg/response"
	"net/http"
	"strings"
)

type UpdateRefundRequest struct {
	RefundID string `json:"refund_id"`
}

func (r *UpdateRefundRequest) Validate() (UnprocessableEntity, error) {
	unprocessableEntity := false
	entity := UnprocessableEntity{
		Fields: map[string]string{
			"refund_id": "",
		},
	}

	r.RefundID = strings.TrimSpace(r.RefundID)
	if r.RefundID == "" {
		unprocessableEntity = true
		entity.Fields["refund_id"] = FieldCannotBeEmptyMessage
	}

	if unprocessableEntity {
		return entity, httperror.New(
			http.StatusUnprocessableEntity,
			response.UnprocessableEntityMessage,
		)
	}

	return entity, nil
}
