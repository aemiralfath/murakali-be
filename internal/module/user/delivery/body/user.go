package body

const (
	FieldCannotBeEmptyMessage         = "Field cannot be empty."
	InvalidPhoneNoFormatMessage       = "Invalid phone no format."
	InvalidEmailFormatMessage         = "Invalid email format."
	InvalidDateFormatMessage          = "Invalid date format."
	InvalidBirthDateAfterTodayMassage = "Birth Date should be in past."
)

type UnprocessableEntity struct {
	Fields map[string]string `json:"fields"`
}
