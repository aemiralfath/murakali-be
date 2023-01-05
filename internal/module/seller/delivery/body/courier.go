package body

import (
	"github.com/google/uuid"
	"murakali/pkg/httperror"
	"murakali/pkg/response"
	"net/http"
	"strings"
)



type CourierSellerResponse struct {
	Rows []*CourierSellerInfo `json:"rows"`
}

type CourierSellerInfo struct {
	ShopCourierID uuid.UUID `json:"id" db:"shop_courier_id" `
	Name          string    `json:"name" db:"name" `
	Code          string    `json:"code" db:"code" `
	Service       string    `json:"service" db:"service" `
	Description   string    `json:"description" db:"description" `
}


type CourierSellerRequest struct {
	CourierID          string `json:"courier_id"`
}


func (r *CourierSellerRequest) Validate() (UnprocessableEntity, error) {
	unprocessableEntity := false
	entity := UnprocessableEntity{
		Fields: map[string]string{
			"courier_id": "",
		},
	}

	r.CourierID  = strings.TrimSpace(r.CourierID )
	if r.CourierID  == "" {
		unprocessableEntity = true
		entity.Fields["courier_id"] = FieldCannotBeEmptyMessage
	}

	if unprocessableEntity {
		return entity, httperror.New(
			http.StatusUnprocessableEntity,
			response.UnprocessableEntityMessage,
		)
	}

	return entity, nil
}
