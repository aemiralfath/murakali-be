package body

import (
	"murakali/pkg/httperror"
	"murakali/pkg/response"
	"net/http"
	"strings"
	"time"
)

type CreatePromotionRequest struct {
	Name               string   `json:"name"`
	ActivedDate        string   `json:"actived_date"`
	ExpiredDate        string   `json:"expired_date"`
	ProductPromotion   []ProductPromotionData `json:"product_promotion"`

	ActiveDateTime  time.Time
	ExpiredDateTime time.Time
}

type ProductPromotionData struct {
	ProductID string `json:"product_id"`
	Quota              int      `json:"quota"`
	MaxQuantity        int      `json:"max_quantity"`
	DiscountPercentage float64  `json:"discount_percentage"`
	DiscountFixPrice   float64  `json:"discount_fix_price"`
	MinProductPrice    float64  `json:"min_product_price"`
	MaxDiscountPrice   float64  `json:"max_discount_price"`

}

type ProductPromotion struct {
	ProductID   *string
	PromotionID *string
}

type ShopProduct struct {
	ShopID    string
	ProductID string
}

func (r *CreatePromotionRequest) Validate() (UnprocessableEntity, error) {
	unprocessableEntity := false
	entity := UnprocessableEntity{
		Fields: map[string]string{
			"name":                "",
			"actived_date":        "",
			"expired_date":        "",
			"product_promotion": "",
		},
	}

	r.Name = strings.TrimSpace(r.Name)
	if r.Name == "" {
		unprocessableEntity = true
		entity.Fields["name"] = FieldCannotBeEmptyMessage
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

	if r.ProductPromotion == nil {
		unprocessableEntity = true
		entity.Fields["product_promotion"] = FieldCannotBeEmptyMessage
	}
	if unprocessableEntity {
		return entity, httperror.New(
			http.StatusUnprocessableEntity,
			response.UnprocessableEntityMessage,
		)
	}

	return entity, nil
}
