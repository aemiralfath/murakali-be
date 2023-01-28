package body

const (
	FieldCannotBeEmptyMessage   = "Field cannot be empty."
	InvalidQuantityValueMessage = "Quantity must greater that 0"
)

type UnprocessableEntity struct {
	Fields map[string]interface{} `json:"fields"`
}
