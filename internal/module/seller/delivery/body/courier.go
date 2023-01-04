package body

import (


	"github.com/google/uuid"
)

type CourierSellerResponse struct {
	Rows []*CourierSellerInfo `json:"rows"`
}


type CourierSellerInfo struct {
	ShopCourierID          uuid.UUID    `json:"id" db:"shop_courier_id" `
	Name        string       `json:"name" db:"name" `
	Code        string       `json:"code" db:"code" `
	Service     string       `json:"service" db:"service" `
	Description string       `json:"description" db:"description" `
}
