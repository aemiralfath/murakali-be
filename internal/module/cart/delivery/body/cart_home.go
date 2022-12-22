package body

import (
	"murakali/pkg/httperror"
	"murakali/pkg/response"
	"net/http"
)

type CartHomeRequest struct {
	Limit int `form:"limit"`
}

type CartHomeResponse struct {
	Limit     int         `json:"limit"`
	TotalItem int64       `json:"total_item"`
	CartHomes []*CartHome `json:"cart_items"`
}

type CartHome struct {
	Title              string  `json:"title" db:"title"`
	ThumbnailURL       string  `json:"thumbnail_url" db:"thumbnail_url"`
	Price              float64 `json:"price" db:"price"`
	DiscountPersentase float64 `json:"discount_percentage" db:"discount_percentage"`
	DiscountFixPrice   float64 `json:"discount_fix_price" db:"discount_fix_price"`
	MinProductPrice    float64 `json:"min_product_price" db:"min_product_price"`
	MaxDiscountPrice   float64 `json:"max_discount_price" db:"max_discount_price"`
	ResultDiscount     float64 `json:"result_discount" db:"result_discount"`
	SubPrice           float64 `json:"sub_price" db:"sub_price"`
	Quantity           int     `json:"quantity" db:"quantity"`
	VariantName        string  `json:"variant_name" db:"variant_name"`
	VariantType        string  `json:"variant_type" db:"variant_type"`
}

func (r *CartHomeRequest) Validate() (UnprocessableEntity, error) {
	unprocessableEntity := false
	entity := UnprocessableEntity{
		Fields: map[string]interface{}{
			"limit": -1,
		},
	}

	if r.Limit < 0 {
		unprocessableEntity = true
		entity.Fields["limit"] = FieldCannotBeEmptyMessage
	}

	if unprocessableEntity {
		return entity, httperror.New(
			http.StatusUnprocessableEntity,
			response.UnprocessableEntityMessage,
		)
	}

	return entity, nil
}
