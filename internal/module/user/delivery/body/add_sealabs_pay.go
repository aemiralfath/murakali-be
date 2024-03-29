package body

import (
	"murakali/pkg/httperror"
	"murakali/pkg/response"
	"net/http"
	"strings"
	"time"
)

type AddSealabsPayRequest struct {
	CardNumber     string `json:"card_number"`
	Name           string `json:"name"`
	IsDefault      bool   `json:"is_default"`
	ActiveDate     string `json:"active_date"`
	ActiveDateTime time.Time
}

func (r *AddSealabsPayRequest) Validate() (UnprocessableEntity, error) {
	unprocessableEntity := false
	entity := UnprocessableEntity{
		Fields: map[string]string{
			"card_number": "",
			"name":        "",
			"is_default":  "",
			"active_date": "",
		},
	}

	r.CardNumber = strings.TrimSpace(r.CardNumber)
	if r.CardNumber == "" {
		unprocessableEntity = true
		entity.Fields["card_number"] = FieldCannotBeEmptyMessage
	}

	r.Name = strings.TrimSpace(r.Name)
	if r.Name == "" {
		unprocessableEntity = true
		entity.Fields["name"] = FieldCannotBeEmptyMessage
	}

	if !r.IsDefault {
		unprocessableEntity = true
		entity.Fields["is_default"] = InvalidIsDefault
	}

	r.ActiveDate = strings.TrimSpace(r.ActiveDate)
	if r.ActiveDate == "" {
		unprocessableEntity = true
		entity.Fields["active_date"] = FieldCannotBeEmptyMessage
	}

	activeTime, err := time.Parse("02-01-2006 15:04:05", r.ActiveDate)
	if err != nil {
		unprocessableEntity = true
		entity.Fields["active_date"] = InvalidDateFormatMessage
	}

	r.ActiveDateTime = activeTime

	if unprocessableEntity {
		return entity, httperror.New(
			http.StatusUnprocessableEntity,
			response.UnprocessableEntityMessage,
		)
	}

	return entity, nil
}
