package body

const (
	FieldCannotBeEmptyMessage = "Field cannot be empty."
	IDNotValidMessage         = "ID not valid."
)

type UnprocessableEntity struct {
	Fields map[string]string `json:"fields"`
}
