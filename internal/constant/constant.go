package constant

const (
	AdminMarketplaceID = "4df967a8-5b05-4d2a-bb72-da3921dce8fb"

	RegisterTokenCookie        = "register_token"
	RefreshTokenCookie         = "refresh_token"
	WalletTokenCookie          = "wallet_token"
	ChangeWalletPinTokenCookie = "change_wallet_pin_token"
	ResetPasswordTokenCookie   = "reset_password_token"
	ChangePasswordTokenCookie  = "change_password_token"

	ProvinceKey    = "location:province"
	CityKey        = "location:city"
	SubDistrictKey = "location:subdistrict"
	UrbanKey       = "location:urban"
	OtpKey         = "user:otp"
	OtpDuration    = "30m"
	AddressDefault = "true"

	RoleUser   = 1
	RoleSeller = 2

	ImgMaxSize = 500000

	SLPStatusPaid      = "TXN_PAID"
	SlPMessagePaid     = "Payment successful"
	SLPStatusCanceled  = "TXN_FAILED"
	SLPMessageCanceled = "Transaction Canceled by user"

	FALSE = "false"
	ASC   = "asc"
	DESC  = "desc"

	LoginOauth    = "/login"
	RegisterOauth = "/register"

	OrderStatusWaitingToPay     = 1
	OrderStatusWaitingForSeller = 2
	OrderStatusProcessed        = 3
	OrderStatusOnDelivery       = 4
	OrderStatusDelivered        = 6
	OrderStatusReceived         = 5
	OrderStatusCompleted        = 7
	OrderStatusCanceled         = 8
	OrderStatusRefunded         = 9
)
