package body

import (
	"murakali/pkg/httperror"
	"murakali/pkg/response"
	"net/http"
	"strings"
)

type VerifyOTPRequest struct {
	OTP string `json:"otp" form:"otp"`
}

type VerifyOTPResponse struct {
	OTP string `json:"otp" form:"otp"`
}

func (r *VerifyOTPRequest) Validate() (UnprocessableEntity, error) {
	unprocessableEntity := false
	entity := UnprocessableEntity{
		Fields: map[string]string{
			"otp": "",
		},
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
