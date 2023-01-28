package body

import (
	"murakali/pkg/httperror"
	"murakali/pkg/response"
	"net/http"
	"strings"
)

type ResetPasswordVerifyOTPRequest struct {
	Code string `form:"code"`
}

func (r *ResetPasswordVerifyOTPRequest) Validate() (UnprocessableEntity, error) {
	unprocessableEntity := false
	entity := UnprocessableEntity{
		Fields: map[string]string{
			"email": "",
			"code":  "",
		},
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
