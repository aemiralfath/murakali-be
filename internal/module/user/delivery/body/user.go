package body

const (
	FieldCannotBeEmptyMessage         = "Field cannot be empty."
	IDNotValidMessage                 = "ID not valid."
	InvalidPhoneNoFormatMessage       = "Invalid phone no format."
	InvalidEmailFormatMessage         = "Invalid email format."
	InvalidDateFormatMessage          = "Invalid date format."
	InvalidBirthDateAfterTodayMassage = "Birth Date should be in past."
	InvalidOTPFormatMessage           = "Invalid OTP."
	InvalidPasswordFormatMessage      = "Password must contain at least 8-40 characters," +
		"at least 1 number, 1 Upper case, 1 special character, and not contains username"
	InvalidPaymentMethod       = "Invalid payment Method."
	InvalidSignatureMessage    = "Invalid Signature."
	InvalidPinFormatMessage    = "Invalid pin format."
	TopUpAmountNotValidMessage = "Top up at least 10000"
)

type UnprocessableEntity struct {
	Fields map[string]string `json:"fields"`
}
