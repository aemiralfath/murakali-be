package body

import (
	"murakali/pkg/httperror"
	"murakali/pkg/response"
	"net/http"
	"strings"

	"github.com/google/uuid"
)

type CreateRefundThreadRequest struct {
	RefundID  string       `json:"refund_id"`
	IsSeller  *bool        `json:"is_seller"`
	IsBuyer   *bool        `json:"is_buyer"`
	Text      string       `json:"text"`
}

func (r *CreateRefundThreadRequest) Validate() (UnprocessableEntity, error) {
	unprocessableEntity := false
	entity := UnprocessableEntity{
		Fields: map[string]string{
			"refund_id": "",
			"is_seller": "",
			"is_buyer": "",
			"text": "",
		},
	}

	r.RefundID = strings.TrimSpace(r.RefundID)
	if r.RefundID == "" {
		unprocessableEntity = true
		entity.Fields["refund_id"] = FieldCannotBeEmptyMessage
	}

	_, err := uuid.Parse(r.RefundID)
	if err != nil {
		unprocessableEntity = true
		entity.Fields["refund_id"] = IDNotValidMessage
	}

	if r.IsSeller == nil {
		unprocessableEntity = true
		entity.Fields["is_seller"] = FieldCannotBeEmptyMessage
	}

	if r.IsBuyer == nil {
		unprocessableEntity = true
		entity.Fields["is_buyer"] = FieldCannotBeEmptyMessage
	}

	r.Text = strings.TrimSpace(r.Text)
	if r.Text == "" {
		unprocessableEntity = true
		entity.Fields["text"] = FieldCannotBeEmptyMessage
	}

	if unprocessableEntity {
		return entity, httperror.New(
			http.StatusUnprocessableEntity,
			response.UnprocessableEntityMessage,
		)
	}

	return entity, nil
}
