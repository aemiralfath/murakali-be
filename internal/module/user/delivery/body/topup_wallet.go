package body

import (
	"murakali/pkg/httperror"
	"murakali/pkg/response"
	"net/http"
	"strings"
)

type TopUpWalletRequest struct {
	CardNumber string `json:"card_number"`
	Amount     int    `json:"amount"`
}

type TopUpWalletResponse struct {
	TransactionID string `json:"transaction_id"`
}

func (r *TopUpWalletRequest) Validate() (UnprocessableEntity, error) {
	unprocessableEntity := false
	entity := UnprocessableEntity{
		Fields: map[string]string{
			"card_number": "",
			"amount":      "",
		},
	}

	r.CardNumber = strings.TrimSpace(r.CardNumber)
	if r.CardNumber == "" {
		unprocessableEntity = true
		entity.Fields["card_number"] = FieldCannotBeEmptyMessage
	}

	if r.Amount < 10000 {
		unprocessableEntity = true
		entity.Fields["amount"] = TopUpAmountNotValidMessage
	}

	if unprocessableEntity {
		return entity, httperror.New(
			http.StatusUnprocessableEntity,
			response.UnprocessableEntityMessage,
		)
	}

	return entity, nil
}
