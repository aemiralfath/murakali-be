package body

import (
	"github.com/google/uuid"
	"murakali/pkg/httperror"
	"murakali/pkg/response"
	"net/http"
	"strings"
)

type CreateSLPPaymentRequest struct {
	TransactionID string `json:"transaction_id"`
}

type CreateSLPPaymentResponse struct {
	RedirectURL string `json:"redirect_url"`
}

type SLPPaymentResponse struct {
	Code    string   `json:"code"`
	Message string   `json:"message"`
	Data    struct{} `json:"data"`
}

func (r *CreateSLPPaymentRequest) Validate() (UnprocessableEntity, error) {
	unprocessableEntity := false
	entity := UnprocessableEntity{
		Fields: map[string]string{
			"transaction_id": "",
		},
	}

	r.TransactionID = strings.TrimSpace(r.TransactionID)
	if r.TransactionID == "" {
		unprocessableEntity = true
		entity.Fields["transaction_id"] = FieldCannotBeEmptyMessage
	}

	_, err := uuid.Parse(r.TransactionID)
	if err != nil {
		unprocessableEntity = true
		entity.Fields["transaction_id"] = IDNotValidMessage
	}

	if unprocessableEntity {
		return entity, httperror.New(
			http.StatusUnprocessableEntity,
			response.UnprocessableEntityMessage,
		)
	}

	return entity, nil
}
