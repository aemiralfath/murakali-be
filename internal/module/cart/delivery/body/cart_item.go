package body

import (
	"murakali/pkg/httperror"
	"murakali/pkg/response"
	"net/http"
	"strings"

	"github.com/google/uuid"
)

type CartItemsResponse struct {
	ID             uuid.UUID                `json:"id" db:"id"`
	Shop           *ShopResponse            `json:"shop"`
	Weight         float64                  `json:"weight"`
	ProductDetails []*ProductDetailResponse `json:"product_details"`
}

type ShopResponse struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

type ProductDetailResponse struct {
	ID           string            `json:"id"`
	Title        string            `json:"title"`
	ThumbnailURL string            `json:"thumbnail_url"`
	ProductPrice float64           `json:"product_price"`
	ProductStock float64           `json:"product_stock"`
	Quantity     float64           `json:"quantity"`
	Weight       float64           `json:"weight"`
	Variant      map[string]string `json:"variant"`
	Promo        *PromoResponse    `json:"promo"`
}

type PromoResponse struct {
	DiscountPercentage *float64 `json:"discount_percentage" db:"discount_percentage"`
	DiscountFixPrice   *float64 `json:"discount_fix_price" db:"discount_fix_price"`
	MinProductPrice    *float64 `json:"min_product_price" db:"min_product_price"`
	MaxDiscountPrice   *float64 `json:"max_discount_price" db:"max_discount_price"`
	ResultDiscount     float64  `json:"result_discount" db:"result_discount"`
	SubPrice           float64  `json:"sub_price" db:"sub_price"`
}

type CartItemRequest struct {
	ProductDetailID string  `json:"product_detail_id"`
	Quantity        float64 `json:"quantity"`
}

func (r *CartItemRequest) Validate() (UnprocessableEntity, error) {
	unprocessableEntity := false
	entity := UnprocessableEntity{
		Fields: map[string]interface{}{
			"product_detail_id": "",
			"quantity":          "",
		},
	}

	r.ProductDetailID = strings.TrimSpace(r.ProductDetailID)
	if r.ProductDetailID == "" {
		unprocessableEntity = true
		entity.Fields["product_detail_id"] = FieldCannotBeEmptyMessage
	}

	if r.Quantity < 1 {
		unprocessableEntity = true
		entity.Fields["quantity"] = InvalidQuantityValueMessage
	}

	if unprocessableEntity {
		return entity, httperror.New(
			http.StatusUnprocessableEntity,
			response.UnprocessableEntityMessage,
		)
	}

	return entity, nil
}
