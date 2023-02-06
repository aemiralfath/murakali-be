package body

import (
	"murakali/pkg/httperror"
	"murakali/pkg/response"
	"net/http"
	"net/mail"
	"strings"
	"time"
)

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	AccessToken string    `json:"access_token"`
	ExpiredAt   time.Time `json:"expired_at"`
}

func (r *LoginRequest) Validate() (UnprocessableEntity, error) {
	unprocessableEntity := false
	entity := UnprocessableEntity{
		Fields: map[string]string{
			"email":    "",
			"password": "",
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
