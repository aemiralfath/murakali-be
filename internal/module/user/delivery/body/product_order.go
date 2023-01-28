package body

import "github.com/google/uuid"

type ProductUnitSoldOrderQty struct {
	ProductID uuid.UUID `json:"product_id" db:"product_id"`
	UnitSold  int64     `json:"unit_sold" db:"unit_sold"`
	Quantity  int64     `json:"quantity" db:"quantity"`
}
