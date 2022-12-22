package body

import (
	"fmt"
	"murakali/internal/util"
	"murakali/pkg/httperror"
	"murakali/pkg/response"
	"net/http"
	"strings"
)

type ChangePasswordRequest struct {
	NewPassword string `json:"password"`
}

func (r *ChangePasswordRequest) Validate() (UnprocessableEntity, error) {
	unprocessableEntity := false
	entity := UnprocessableEntity{
		Fields: map[string]string{
			"password": "",
		},
	}

	r.NewPassword = strings.TrimSpace(r.NewPassword)
	if r.NewPassword == "" {
		unprocessableEntity = true
		entity.Fields["password"] = FieldCannotBeEmptyMessage
	}
	fmt.Println(r.NewPassword)
	if !util.VerifyPassword(r.NewPassword) {
		unprocessableEntity = true
		entity.Fields["password"] = InvalidPasswordFormatMessage
	}

	if unprocessableEntity {
		return entity, httperror.New(
			http.StatusUnprocessableEntity,
			response.UnprocessableEntityMessage,
		)
	}

	return entity, nil
}
