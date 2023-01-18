package body

import (
	"murakali/pkg/httperror"
	"murakali/pkg/response"
	"net/http"
	"strings"
	"time"
)

type UpdatePromotionRequest struct {
	PromotionID        string  `json:"promotion_id"`
	ProductID          string  `json:"product_id"`
	Name               string  `json:"name"`
	MaxQuantity        int     `json:"max_quantity"`
	ActivedDate        string  `json:"actived_date"`
	ExpiredDate        string  `json:"expired_date"`
	DiscountPercentage float64 `json:"discount_percentage"`
	DiscountFixPrice   float64 `json:"discount_fix_price"`
	MinProductPrice    float64 `json:"min_product_price"`
	MaxDiscountPrice   float64 `json:"max_discount_price"`

	ActiveDateTime  time.Time
	ExpiredDateTime time.Time
}

type ShopProductPromo struct {
	ShopID      string
	PromotionID string
	ProductID   string
}

func (r *UpdatePromotionRequest) Validate() (UnprocessableEntity, error) {
	unprocessableEntity := false
	entity := UnprocessableEntity{
		Fields: map[string]string{
			"promotion_id":        "",
			"product_id":          "",
			"name":                "",
			"max_quantity":        "",
			"actived_date":        "",
			"expired_date":        "",
			"discount_percentage": "",
			"discount_fix_price":  "",
			"min_product_price":   "",
			"max_discount_price":  "",
		},
	}

	r.PromotionID = strings.TrimSpace(r.PromotionID)
	if r.PromotionID == "" {
		unprocessableEntity = true
		entity.Fields["promotion_id"] = FieldCannotBeEmptyMessage
	}

	r.ProductID = strings.TrimSpace(r.ProductID)
	if r.ProductID == "" {
		unprocessableEntity = true
		entity.Fields["product_id"] = FieldCannotBeEmptyMessage
	}

	r.Name = strings.TrimSpace(r.Name)
	if r.Name == "" {
		unprocessableEntity = true
		entity.Fields["name"] = FieldCannotBeEmptyMessage
	}

	if r.MaxQuantity <= 0 {
		unprocessableEntity = true
		entity.Fields["max_quantity"] = FieldCannotBeEmptyMessage
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
