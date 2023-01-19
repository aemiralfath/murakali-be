package user

import "github.com/gin-gonic/gin"

type Handlers interface {
	GetAddress(c *gin.Context)
	CreateAddress(c *gin.Context)
	GetAddressByID(c *gin.Context)
	UpdateAddressByID(c *gin.Context)
	DeleteAddressByID(c *gin.Context)
	EditUser(c *gin.Context)
	EditEmail(c *gin.Context)
	EditEmailUser(c *gin.Context)
	GetSealabsPay(c *gin.Context)
	AddSealabsPay(c *gin.Context)
	PatchSealabsPay(c *gin.Context)
	DeleteSealabsPay(c *gin.Context)
	RegisterMerchant(c *gin.Context)
	GetUserProfile(c *gin.Context)
	UploadProfilePicture(c *gin.Context)
	VerifyPasswordChange(c *gin.Context)
	VerifyOTP(c *gin.Context)
	ChangePassword(c *gin.Context)
	CreateTransaction(c *gin.Context)
	CreateSLPPayment(c *gin.Context)
	SLPPaymentCallback(c *gin.Context)
	WalletPaymentCallback(c *gin.Context)
	GetOrder(c *gin.Context)
	ActivateWallet(c *gin.Context)
	GetWallet(c *gin.Context)
	GetWalletHistory(c *gin.Context)
	TopUpWallet(c *gin.Context)
	WalletStepUp(c *gin.Context)
	CreateWalletPayment(c *gin.Context)
	ChangeWalletPinStepUp(c *gin.Context)
	ChangeWalletPin(c *gin.Context)
}
