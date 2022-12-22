package body

const (
	FieldCannotBeEmptyMessage = "Field cannot be empty."
)

type UnprocessableEntity struct {
	Fields map[string]interface{} `json:"fields"`
}
