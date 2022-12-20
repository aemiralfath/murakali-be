package body

import (
	"murakali/pkg/httperror"
	"murakali/pkg/response"
	"net/http"
	"strings"
)

type RegisterMerchant struct {
	ShopName string `json:"shop_name"`
}

func (r *RegisterMerchant) Validate() (UnprocessableEntity, error) {
	unprocessableEntity := false
	entity := UnprocessableEntity{
		Fields: map[string]string{
			"shop_name": "",
		},
	}

	r.ShopName = strings.TrimSpace(r.ShopName)
	if r.ShopName == "" {
		unprocessableEntity = true
		entity.Fields["card_number"] = FieldCannotBeEmptyMessage
	}

	if unprocessableEntity {
		return entity, httperror.New(
			http.StatusUnprocessableEntity,
			response.UnprocessableEntityMessage,
		)
	}

	return entity, nil
}
