package body

import (
	"murakali/internal/util"
	"murakali/pkg/httperror"
	"murakali/pkg/response"
	"net/http"
	"regexp"
	"strings"
)

type RegisterUserRequest struct {
	Username string `json:"username"`
	FullName string `json:"fullname"`
	Password string `json:"password"`
	PhoneNo  string `json:"phone_no"`
}

func (r *RegisterUserRequest) Validate() (UnprocessableEntity, error) {
	unprocessableEntity := false
	entity := UnprocessableEntity{
		Fields: map[string]string{
			"username": "",
			"fullname": "",
			"password": "",
			"phone_no": "",
		},
	}

	r.Username = strings.TrimSpace(r.Username)
	if r.Username == "" {
		unprocessableEntity = true
		entity.Fields["username"] = FieldCannotBeEmptyMessage
	}

	r.FullName = strings.TrimSpace(r.FullName)
	if r.FullName == "" {
		unprocessableEntity = true
		entity.Fields["fullname"] = FieldCannotBeEmptyMessage
	}

	if !util.VerifyPassword(r.Password) {
		unprocessableEntity = true
		entity.Fields["password"] = InvalidPasswordFormatMessage
	}

	if strings.Contains(strings.ToLower(r.Password), strings.ToLower(r.Username)) {
		unprocessableEntity = true
		entity.Fields["password"] = InvalidPasswordFormatMessage
	}

	regex := regexp.MustCompile(`^8[1-9]\d{6,9}$`)
	if !regex.MatchString(r.PhoneNo) {
		unprocessableEntity = true
		entity.Fields["phone_no"] = InvalidPhoneNoFormatMessage
	}

	if unprocessableEntity {
		return entity, httperror.New(
			http.StatusUnprocessableEntity,
			response.UnprocessableEntityMessage,
		)
	}

	return entity, nil
}
