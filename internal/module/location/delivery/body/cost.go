package body

import (
	"github.com/google/uuid"
	"murakali/internal/model"
	"murakali/pkg/httperror"
	"murakali/pkg/response"
	"net/http"
	"strings"
)

type GetShippingCostRequest struct {
	Destination int      `json:"destination"`
	Weight      int      `json:"weight"`
	ShopID      string   `json:"shop_id"`
	ProductIDS  []string `json:"product_ids"`
}

type GetShippingCostResponse struct {
	ShippingOption []*model.Cost `json:"shipping_option"`
}

func (r *GetShippingCostRequest) Validate() (UnprocessableEntity, error) {
	unprocessableEntity := false
	entity := UnprocessableEntity{
		Fields: map[string]interface{}{
			"destination": "",
			"weight":      "",
			"shop_id":     "",
			"product_ids": "",
		},
	}

	if r.Destination == 0 {
		unprocessableEntity = true
		entity.Fields["destination"] = FieldCannotBeEmptyMessage
	}

	if r.Weight <= 0 {
		unprocessableEntity = true
		entity.Fields["weight"] = FieldCannotBeEmptyMessage
	}

	r.ShopID = strings.TrimSpace(r.ShopID)
	if _, err := uuid.Parse(r.ShopID); err != nil {
		unprocessableEntity = true
		entity.Fields["shop_id"] = InvalidIDMessage
	}

	if len(r.ProductIDS) == 0 {
		unprocessableEntity = true
		entity.Fields["product_ids"] = FieldCannotBeEmptyMessage
	}

	for i, productID := range r.ProductIDS {
		r.ProductIDS[i] = strings.TrimSpace(productID)
		if _, err := uuid.Parse(r.ProductIDS[i]); err != nil {
			unprocessableEntity = true
			entity.Fields["shop_id"] = r.ProductIDS[i] + InvalidIDMessage
		}
	}

	if unprocessableEntity {
		return entity, httperror.New(
			http.StatusUnprocessableEntity,
			response.UnprocessableEntityMessage,
		)
	}

	return entity, nil
}
