package body

import (
	"murakali/pkg/httperror"
	"murakali/pkg/response"
	"net/http"
	"strings"
)

type AddCartItemRequest struct {
	ProductDetailID string  `json:"product_detail_id"`
	Quantity        float64 `json:"quantity"`
}

func (r *AddCartItemRequest) Validate() (UnprocessableEntity, error) {
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
