package body

import (
	"murakali/pkg/httperror"
	"murakali/pkg/response"
	"net/http"
	"strconv"
	"strings"
)

type ActivateWalletRequest struct {
	Pin string `json:"pin"`
}

func (r *ActivateWalletRequest) Validate() (UnprocessableEntity, error) {
	unprocessableEntity := false
	entity := UnprocessableEntity{
		Fields: map[string]string{
			"pin": "",
		},
	}

	r.Pin = strings.TrimSpace(r.Pin)
	if r.Pin == "" {
		unprocessableEntity = true
		entity.Fields["pin"] = FieldCannotBeEmptyMessage
	}

	if len(r.Pin) != 6 {
		unprocessableEntity = true
		entity.Fields["pin"] = InvalidPinFormatMessage
	}

	if _, err := strconv.Atoi(r.Pin); err != nil {
		unprocessableEntity = true
		entity.Fields["pin"] = InvalidPinFormatMessage
	}

	if unprocessableEntity {
		return entity, httperror.New(
			http.StatusUnprocessableEntity,
			response.UnprocessableEntityMessage,
		)
	}

	return entity, nil
}
