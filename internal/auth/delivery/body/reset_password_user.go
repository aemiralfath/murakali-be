package body

import (
	"murakali/internal/util"
	"murakali/pkg/httperror"
	"murakali/pkg/response"
	"net/http"
)

type ResetPasswordUserRequest struct {
	Password string `json:"password"`

	IsPasswordSameOldPassword  bool
	IsPasswordContainsUsername bool
}

func (r *ResetPasswordUserRequest) Validate() (UnprocessableEntity, error) {
	unprocessableEntity := false

	entity := UnprocessableEntity{
		Fields: map[string]string{
			"password": "",
		},
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
