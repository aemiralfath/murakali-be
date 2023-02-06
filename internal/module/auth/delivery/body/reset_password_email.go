package body

import (
	"murakali/pkg/httperror"
	"murakali/pkg/response"
	"net/http"
	"net/mail"
	"strings"
)

type ResetPasswordEmailRequest struct {
	Email string `json:"email"`
}

func (r *ResetPasswordEmailRequest) Validate() (UnprocessableEntity, error) {
	unprocessableEntity := false
	entity := UnprocessableEntity{
		Fields: map[string]string{
			"email": "",
		},
	}

	r.Email = strings.TrimSpace(r.Email)
	if r.Email == "" {
		unprocessableEntity = true
		entity.Fields["email"] = FieldCannotBeEmptyMessage
	}

	_, err := mail.ParseAddress(r.Email)
	if err != nil {
		unprocessableEntity = true
		entity.Fields["email"] = InvalidEmailFormatMessage
	}

	if unprocessableEntity {
		return entity, httperror.New(
			http.StatusUnprocessableEntity,
			response.UnprocessableEntityMessage,
		)
	}

	return entity, nil
}
