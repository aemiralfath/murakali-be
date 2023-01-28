package body

import (
	"murakali/pkg/httperror"
	"murakali/pkg/response"
	"net/http"
	"strings"
	"time"
)

type CreateVoucherRequest struct {
	Code               string  `json:"code"`
	Quota              int     `json:"quota"`
	ActivedDate        string  `json:"actived_date"`
	ExpiredDate        string  `json:"expired_date"`
	DiscountPercentage float64 `json:"discount_percentage"`
	DiscountFixPrice   float64 `json:"discount_fix_price"`
	MinProductPrice    float64 `json:"min_product_price"`
	MaxDiscountPrice   float64 `json:"max_discount_price"`

	ActiveDateTime  time.Time
	ExpiredDateTime time.Time
}

func (r *CreateVoucherRequest) Validate() (UnprocessableEntity, error) {
	unprocessableEntity := false
	entity := UnprocessableEntity{
		Fields: map[string]string{
			"code":                "",
			"actived_date":        "",
			"expired_date":        "",
			"discount_percentage": "",
			"discount_fix_price":  "",
			"min_product_price":   "",
			"max_discount_price":  "",
		},
	}

	r.Code = strings.TrimSpace(r.Code)
	if r.Code == "" {
		unprocessableEntity = true
		entity.Fields["code"] = FieldCannotBeEmptyMessage
	}

	if r.Quota < 0 {
		unprocessableEntity = true
		entity.Fields["quota"] = FieldCannotBeEmptyMessage
	}

	activeTime, err := time.Parse("02-01-2006 15:04:05", r.ActivedDate)
	if err != nil {
		unprocessableEntity = true
		entity.Fields["actived_date"] = InvalidDateFormatMessage
	}
	r.ActiveDateTime = activeTime

	expireTime, err := time.Parse("02-01-2006 15:04:05", r.ExpiredDate)
	if err != nil {
		unprocessableEntity = true
		entity.Fields["expired_date"] = InvalidDateFormatMessage
	}
	r.ExpiredDateTime = expireTime

	if r.DiscountPercentage <= 0 && r.DiscountFixPrice <= 0 {
		unprocessableEntity = true
		entity.Fields["discount_percentage"] = FieldCannotBeEmptyMessage
		entity.Fields["discount_fix_price"] = FieldCannotBeEmptyMessage
	}

	if r.MinProductPrice <= 0 {
		unprocessableEntity = true
		entity.Fields["min_product_price"] = FieldCannotBeEmptyMessage
	}

	if r.MaxDiscountPrice <= 0 {
		unprocessableEntity = true
		entity.Fields["max_discount_price"] = FieldCannotBeEmptyMessage
	}

	if unprocessableEntity {
		return entity, httperror.New(
			http.StatusUnprocessableEntity,
			response.UnprocessableEntityMessage,
		)
	}

	return entity, nil
}
