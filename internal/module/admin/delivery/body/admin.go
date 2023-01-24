package body

const (
	FieldCannotBeEmptyMessage    = "Field cannot be empty."
	IDNotValidMessage            = "ID not valid."
	InvalidPhoneNoFormatMessage  = "Invalid phone no format."
	InvalidEmailFormatMessage    = "Invalid email format."
	InvalidDateFormatMessage     = "Invalid date format."
	InvalidOTPFormatMessage      = "Invalid OTP."
	InvalidPasswordFormatMessage = "Password must contain at least 8-40 characters," +
		"at least 1 number, 1 Upper case, 1 special character, and not contains username"
	InvalidPasswordSameOldPasswordMessage = "Your new password cannot be the same as your old password."
	VoucherSellerNotFoundMessage          = "Voucher Admin Not Found"
	CodeVoucherAlreadyExist               = "Code Voucher Already Exist"
	CategoryLevelInvalid                  = "Category Level Invalid"
	ProductNotFound                       = "Product not found"
	UpdateProductFailed                   = "Update product failed"
	ImageIsEmpty                          = "image cannot be empty"
)

type UnprocessableEntity struct {
	Fields map[string]string `json:"fields"`
}
