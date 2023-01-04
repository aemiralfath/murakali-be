package constant

const (
	RegisterTokenCookie       = "register_token"
	RefreshTokenCookie        = "refresh_token"
	ResetPasswordTokenCookie  = "reset_password_token"
	ChangePasswordTokenCookie = "change_password_token"

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

	SLPStatusPaid = "TXN_PAID"

	OrderStatusWaitingToPay      = 1
	OrderStatusWaitingForSeller  = 2
	OrderStatusWaitingForPacking = 3
	OrderStatusOnDelivery        = 4
	OrderStatusCompleted         = 5
	OrderStatusReceived          = 6
	OrderStatusCanceled          = 7
)
