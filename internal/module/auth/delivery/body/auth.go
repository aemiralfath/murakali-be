package body

const (
	FieldCannotBeEmptyMessage    = "Field cannot be empty."
	InvalidEmailFormatMessage    = "Invalid email format."
	InvalidOTPFormatMessage      = "Invalid otp format."
	InvalidPhoneNoFormatMessage  = "Invalid phone no format."
	InvalidPasswordFormatMessage = "Password must contain at least 8-40 characters," +
		"at least 1 number, 1 Upper case, 1 special character, and not contains username"
	InvalidPasswordSameOldPasswordMessage = "Your new password cannot be the same as your old password."
)

type UnprocessableEntity struct {
	Fields map[string]string `json:"fields"`
}
