package body

const (
	InvalidEmailFormatMessage = "Invalid email format."
	InvalidOTPFormatMessage   = "Invalid otp format."
)

type UnprocessableEntity struct {
	Fields map[string]string `json:"fields"`
}
