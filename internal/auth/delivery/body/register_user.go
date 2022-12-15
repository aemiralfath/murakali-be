package body

import (
	"murakali/internal/util"
	"murakali/pkg/httperror"
	"murakali/pkg/response"
	"net/http"
	"net/mail"
	"regexp"
	"strings"
)

type RegisterUserRequest struct {
	Email    string `json:"email"`
	OTP      string `json:"otp"`
	Username string `json:"username"`
	FullName string `json:"fullname"`
	Password string `json:"password"`
	PhoneNo  string `json:"phone_no"`
}

func (r *RegisterUserRequest) Validate() (UnprocessableEntity, error) {
	unprocessableEntity := false
	entity := UnprocessableEntity{
		Fields: map[string]string{
			"email":    "",
			"otp":      "",
			"username": "",
			"fullname": "",
			"password": "",
			"phone_no": "",
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

	if strings.Contains(strings.ToLower(r.Password), r.Username) {
		unprocessableEntity = true
		entity.Fields["password"] = InvalidPasswordFormatMessage
	}

	regex, err := regexp.Compile(`^8[1-9][0-9]{6,9}$`)
	if err != nil {
		return entity, err
	}

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
