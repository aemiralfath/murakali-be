package body

import (
	"murakali/pkg/httperror"
	"murakali/pkg/response"
	"net/http"
	"regexp"
	"strings"
	"time"
)

type EditUserRequest struct {
	Username  string `json:"username"`
	PhoneNo   string `json:"phone_no"`
	FullName  string `json:"fullname"`
	Gender    string `json:"gender"`
	BirthDate string `json:"birth_date"`
}

func (r *EditUserRequest) Validate() (UnprocessableEntity, error) {
	unprocessableEntity := false
	entity := UnprocessableEntity{
		Fields: map[string]string{
			"username":   "",
			"phone_no":   "",
			"fullname":   "",
			"gender":     "",
			"birth_date": "",
		},
	}

	r.Username = strings.TrimSpace(r.Username)
	if r.Username == "" {
		unprocessableEntity = true
		entity.Fields["username"] = FieldCannotBeEmptyMessage
	}

	regex := regexp.MustCompile(`^8[1-9]\d{6,9}$`)
	if !regex.MatchString(r.PhoneNo) {
		unprocessableEntity = true
		entity.Fields["phone_no"] = InvalidPhoneNoFormatMessage
	}

	r.FullName = strings.TrimSpace(r.FullName)
	if r.FullName == "" {
		unprocessableEntity = true
		entity.Fields["fullname"] = FieldCannotBeEmptyMessage
	}

	r.Gender = strings.TrimSpace(r.Gender)
	if r.Gender == "" {
		unprocessableEntity = true
		entity.Fields["gender"] = FieldCannotBeEmptyMessage
	}

	r.BirthDate = strings.TrimSpace(r.BirthDate)
	if r.BirthDate == "" {
		unprocessableEntity = true
		entity.Fields["birth_date"] = FieldCannotBeEmptyMessage
	}

	today, _ := time.Parse("02-01-2006", time.Now().Format("02-01-2006"))
	BirthDate, err := time.Parse("02-01-2006", r.BirthDate)
	if err != nil {
		unprocessableEntity = true
		entity.Fields["birth_date"] = InvalidDateFormatMessage
	}

	if BirthDate.Equal(today) || BirthDate.After(today) {
		unprocessableEntity = true
		entity.Fields["birth_date"] = InvalidBirthDateAfterTodayMassage
	}

	if unprocessableEntity {
		return entity, httperror.New(
			http.StatusUnprocessableEntity,
			response.UnprocessableEntityMessage,
		)
	}

	return entity, nil
}
