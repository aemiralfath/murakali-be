package body

const (
	InvalidEmailFormatMessage = "Invalid email format."
)

type UnprocessableEntity struct {
	Fields map[string]string `json:"fields"`
}
