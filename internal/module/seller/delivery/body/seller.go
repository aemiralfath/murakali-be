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
	InvalidPasswordSameOldPasswordMessage      = "Your new password cannot be the same as your old password."
	SellerNotFoundMessage                      = "Seller not found."
	CourierSellerAlreadyExistMessage           = "Courier Seller Already Exist"
	CourierNotFoundMessage                     = "Courier not exist"
	CourierSellerNotFoundMessage               = "Courier Seller Not Found"
	VoucherSellerNotFoundMessage               = "Voucher Seller Not Found"
	PromotionSellerNotFoundMessage             = "Promotion Seller Not Found"
	PromotionExpairedMessage                   = "Promotion is Expired"
	PromotionCannotUpdateActivedMessage        = "Promotion Active Date cannot be update. Promotion Currently Ongoing."
	PromotionCannotUpdateActivedBeforeMessage  = "Promotion Active Date cannot be update before its actived date."
	PromotionCannotUpdateExpiredAfterMessage   = "Promotion Expired Date cannot be update after its Expired Date."
	PromotionCannotUpdateExpiredBefoferMessage = "Promotion Expired Date cannot be update before current time."
	CategoryNotFoundMessage                    = "Category not found."
	ShopNotFoundMessage                        = "Shop not found."
	CodeVoucherAlreadyExist                    = "Code Voucher Already Exist"
)

type UnprocessableEntity struct {
	Fields map[string]string `json:"fields"`
}
