package body

import (
	"murakali/pkg/httperror"
	"murakali/pkg/response"
	"net/http"
	"net/mail"
	"strings"
)

type ResetPasswordVerifyOTPRequest struct {
	Email string `form:"email"`
	Code  string `form:"code"`
}

type ResetPasswordVerifyOTPResponse struct {
	Email string `form:"email"`
	Code  string `form:"code"`
}

func (r *ResetPasswordVerifyOTPRequest) Validate() (UnprocessableEntity, error) {
	unprocessableEntity := false
	entity := UnprocessableEntity{
		Fields: map[string]string{
			"email": "",
			"code":  "",
		},
	}

	r.Email = strings.TrimSpace(r.Email)
	if r.Email == "" {
		unprocessableEntity = true
		entity.Fields["email"] = InvalidEmailFormatMessage
	}

	_, err := mail.ParseAddress(r.Email)
	if err != nil {
		unprocessableEntity = true
		entity.Fields["email"] = InvalidEmailFormatMessage
	}

	r.Code = strings.TrimSpace(r.Code)
	if r.Code == "" {
		unprocessableEntity = true
		entity.Fields["code"] = InvalidOTPFormatMessage
	}

	if unprocessableEntity {
		return entity, httperror.New(
			http.StatusUnprocessableEntity,
			response.UnprocessableEntityMessage,
		)
	}

	return entity, nil
}
