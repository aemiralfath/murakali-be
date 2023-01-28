package body

import (
	"database/sql"
	"murakali/pkg/httperror"
	"murakali/pkg/response"
	"net/http"
	"strings"

	"github.com/google/uuid"
)

type CourierSellerResponse struct {
	Rows []*CourierSellerInfo `json:"rows"`
}

type CourierSellerInfo struct {
	ShopCourierID string    `json:"shop_courier_id" db:"shop_courier_id" `
	CourierID     uuid.UUID `json:"courier_id" db:"courier_id" `
	Name          string    `json:"name" db:"name" `
	Code          string    `json:"code" db:"code" `
	Service       string    `json:"service" db:"service" `
	Description   string    `json:"description" db:"description" `
	DeletedAt     string    `json:"deleted_at" db:"deleted_at" `
}

type CourierSellerRelationInfo struct {
	ShopCourierID uuid.UUID    `json:"shop_courier_id" db:"shop_courier_id" `
	CourierID     uuid.UUID    `json:"courier_id" db:"courier_id" `
	DeletedAt     sql.NullTime `json:"deleted_at" db:"deleted_at" `
}

type CourierInfo struct {
	CourierID   uuid.UUID `json:"courier_id" db:"courier_id" `
	Name        string    `json:"name" db:"name" `
	Code        string    `json:"code" db:"code" `
	Service     string    `json:"service" db:"service" `
	Description string    `json:"description" db:"description" `
}

type CourierSellerRequest struct {
	CourierID string `json:"courier_id"`
}

func (r *CourierSellerRequest) Validate() (UnprocessableEntity, error) {
	unprocessableEntity := false
	entity := UnprocessableEntity{
		Fields: map[string]string{
			"courier_id": "",
		},
	}

	r.CourierID = strings.TrimSpace(r.CourierID)
	if r.CourierID == "" {
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
