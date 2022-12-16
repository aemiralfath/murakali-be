package body

import (
	"murakali/internal/util"
	"murakali/pkg/httperror"
	"murakali/pkg/response"
	"net/http"
	"net/mail"
	"strings"
)

type ResetPasswordUserRequest struct {
	Email    string `json:"email"`
	OTP      string `json:"otp"`
	Password string `json:"password"`

	IsPasswordSameOldPassword  bool
	IsPasswordContainsUsername bool
}

func (r *ResetPasswordUserRequest) Validate() (UnprocessableEntity, error) {
	unprocessableEntity := false

	entity := UnprocessableEntity{
		Fields: map[string]string{
			"email":    "",
			"otp":      "",
			"password": "",
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

	r.OTP = strings.TrimSpace(r.OTP)
	if len(r.OTP) != 6 {
		unprocessableEntity = true
		entity.Fields["otp"] = InvalidOTPFormatMessage
	}

	if !util.VerifyPassword(r.Password) {
		unprocessableEntity = true
		entity.Fields["password"] = InvalidPasswordFormatMessage
	}

	if r.IsPasswordContainsUsername {
		unprocessableEntity = true
		entity.Fields["password"] = InvalidPasswordFormatMessage
	}

	if r.IsPasswordSameOldPassword {
		unprocessableEntity = true
		entity.Fields["password"] = InvalidPasswordSameOldPasswordMessage
	}

	if unprocessableEntity {
		return entity, httperror.New(
			http.StatusUnprocessableEntity,
			response.UnprocessableEntityMessage,
		)
	}

	return entity, nil
}
