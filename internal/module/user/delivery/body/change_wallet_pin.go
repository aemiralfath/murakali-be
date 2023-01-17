package body

import (
	"murakali/pkg/httperror"
	"murakali/pkg/response"
	"net/http"
	"strconv"
	"strings"
)

type ChangeWalletPinStepUpRequest struct {
	Password string `json:"password"`
}

type ChangeWalletPinRequest struct {
	Pin string `json:"pin"`
}

func (r *ChangeWalletPinRequest) Validate() (UnprocessableEntity, error) {
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

func (r *ChangeWalletPinStepUpRequest) Validate() (UnprocessableEntity, error) {
	unprocessableEntity := false
	entity := UnprocessableEntity{
		Fields: map[string]string{
			"password": "",
		},
	}

	r.Password = strings.TrimSpace(r.Password)
	if r.Password == "" {
		unprocessableEntity = true
		entity.Fields["password"] = FieldCannotBeEmptyMessage
	}

	if unprocessableEntity {
		return entity, httperror.New(
			http.StatusUnprocessableEntity,
			response.UnprocessableEntityMessage,
		)
	}

	return entity, nil
}
