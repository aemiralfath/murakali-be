package response

import (
	"encoding/json"
	"net/http"
)

const (
	BadRequestMessage          = "Invalid request."
	UnprocessableEntityMessage = "Request field not valid."
	InternalServerErrorMessage = "Something is wrong, pls try again later."
	NotFoundMessage            = "Route does not exist, please check again your route path."
	UnauthorizedMessage        = "Invalid Credentials."
	ForbiddenMessage           = "Forbidden"

	AddressIsDefaultMessage        = "Address is default."
	EmailAlreadyExistMessage       = "User already registered."
	EmailAlreadyChangedMessage     = "Email already changed."
	EmailSamePreviousEmailMessage  = "This email same as your current email."
	EmailNotExistMessage           = "User not registered."
	UserNotVerifyMessage           = "User not verify."
	UserAlreadyExistMessage        = "User already exist."
	UserNameAlreadyExistMessage    = "Username already exist."
	PhoneNoAlreadyExistMessage     = "Phone no already exist."
	UserNotExistMessage            = "User not exist."
	AddressNotExistMessage         = "Address not exist."
	PasswordSameOldPasswordMessage = "Your new password cannot be the same as your old password."
	PasswordContainUsernameMessage = "Password contains username."
	OTPAlreadyExpiredMessage       = "OTP already expired."
	OTPIsNotValidMessage           = "OTP is not valid."
	UserAlreadyHaveShop            = "User already have shop."
	ShopAlreadyExists              = "Shop already exists."
	QuantityReachedMaximum         = "Quantity has reached the maximum limit!"
	ProductDetailNotExistMessage   = "Product Detail not exist."
	ProductNotExistMessage         = "Product not exist."
	ProductAlreadyHasPromoMessage  = "Product Already has Promotion"
	ProductAlreadyInFavMessage     = "Product already in favorite."
	PictureSizeTooBig              = "Picture size too big"
	TransactionIDNotExist          = "Transaction not exist."
	TransactionAlreadyExpired      = "Transaction already expired."
	TransactionAlreadyFinished     = "Transaction already finished."
	SelectShippingCourier          = "Select shipping Courier"
	UnknownShop                    = "Unknown shop."
	CartItemNotExist               = "Cart Item not exist."
	ProductQuantityNotAvailable    = "Product quantity not available."
	ShopAddressNotFound            = "Shop address not found."
	UserNotHaveShop                = "User not register shop"
	DefaultAddressNotFound         = "Default address not found."
	ShopCourierNotExist            = "Shop courier not exist."
	WalletAlreadyActivated         = "Wallet already activated."
	WalletIsNotActivated           = `Wallet is not activated.`
	SealabsCardNotFound            = "Sealabs pay card not valid."
	SealabsCardAlreadyExist        = "Sealabs card already exist."
	WalletIsBlocked                = "Wallet is temporarily blocked, please wait."
	WalletPinIsInvalid             = "Wallet pin is invalid."
	WalletBalanceNotEnough         = "Insufficient wallet balance, please top up!"
	InvalidPaymentMethod           = "Invalid payment method."
	OrderNotExistMessage           = "Order not exist."
	OrderNotCompletedMessage       = "Order not completed."
	OrderAlreadyWithdrawMessage    = "Order already withdraw."
	InvalidPasswordMessage         = "Invalid password."
	TransactionNotFound            = "Transaction not found."
	RefundNotFound                 = "Refund not found."
	RefundAlreadyFinished          = "Refund already finished."
	RefundRejected                 = "Refund is rejected."
	WalletHistoryNotFound          = "Wallet history not found."
	OrderNotWaitingForSeller       = "Order not waiting for seller"
	VoucherMarketplaceNotFound     = "Voucher Marketplace Not Found"
	VoucherShopNotFound            = "Voucher Shop Not Found"
	OrderUnderProgressRefund       = "Order is Under Progress Refunding"
	InvalidRefund                  = "Invalid Refund"
	OrderCannotToRefund            = "Order Cannot to Refund"
	OrderHasAcceptedToRefund       = "Order Has Accepted to Refund"
	OrderRefundHasBeenFinished     = "Order Refund Has Been Finished"
)

type JSONResponse struct {
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

func returnJSONResponse(w http.ResponseWriter, message string, data interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_ = json.NewEncoder(w).Encode(JSONResponse{
		Message: message,
		Data:    data,
	})
}

func SuccessResponse(w http.ResponseWriter, data interface{}, statusCode ...int) {
	code := http.StatusOK
	if len(statusCode) > 0 {
		code = statusCode[0]
	}

	returnJSONResponse(w, "success", data, code)
}

func ErrorResponse(w http.ResponseWriter, message string, statusCode ...int) {
	code := http.StatusBadRequest
	if len(statusCode) > 0 {
		code = statusCode[0]
	}

	returnJSONResponse(w, message, nil, code)
}

func ErrorResponseData(w http.ResponseWriter, data interface{}, message string, statusCode ...int) {
	code := http.StatusBadRequest
	if len(statusCode) > 0 {
		code = statusCode[0]
	}

	returnJSONResponse(w, message, data, code)
}
