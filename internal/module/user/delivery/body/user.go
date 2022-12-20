package body

const (
	FieldCannotBeEmptyMessage         = "Field cannot be empty."
	InvalidDateFormatMessage          = "Invalid date format."
	InvalidPhoneNoFormatMessage       = "Invalid phone no format."
	InvalidEmailFormatMessage         = "Invalid email format."
	InvalidBirthDateAfterTodayMassage = "Birth Date should be in past."
)

type UnprocessableEntity struct {
	Fields map[string]string `json:"fields"`
}
