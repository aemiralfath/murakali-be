package body

import (
	"murakali/pkg/httperror"
	"murakali/pkg/response"
	"net/http"
	"net/mail"
	"strings"
)

type VerifyOTPRequest struct {
	Email string `json:"email" form:"email"`
	OTP   string `json:"otp" form:"otp"`
}

type VerifyOTPResponse struct {
	Email string `json:"email" form:"email"`
	OTP   string `json:"otp" form:"otp"`
}

func (r *VerifyOTPRequest) Validate() (UnprocessableEntity, error) {
	unprocessableEntity := false
	entity := UnprocessableEntity{
		Fields: map[string]string{
			"email": "",
			"otp":   "",
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

	if unprocessableEntity {
		return entity, httperror.New(
			http.StatusUnprocessableEntity,
			response.UnprocessableEntityMessage,
		)
	}

	return entity, nil
}
