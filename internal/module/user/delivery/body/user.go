package body

const (
	FieldCannotBeEmptyMessage = "Field cannot be empty."
	InvalidDateFormatMessage  = "Invalid date format."
)

type UnprocessableEntity struct {
	Fields map[string]string `json:"fields"`
}
