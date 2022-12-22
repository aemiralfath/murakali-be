package delivery

import (
	"murakali/internal/middleware"
	"murakali/internal/module/user"

	"github.com/gin-gonic/gin"
)

func MapUserRoutes(userGroup *gin.RouterGroup, h user.Handlers, mw *middleware.MWManager) {
	userGroup.Use(mw.AuthJWTMiddleware())
	userGroup.GET("/address", h.GetAddress)
	userGroup.POST("/address", h.CreateAddress)
	userGroup.GET("/address/:id", h.GetAddressByID)
	userGroup.PUT("/address/:id", h.UpdateAddressByID)
	userGroup.DELETE("/address/:id", h.DeleteAddressByID)
	userGroup.PUT("/profile", h.EditUser)
	userGroup.POST("/email", h.EditEmail)
	userGroup.GET("/email", h.EditEmailUser)
	userGroup.GET("/sealab-pay", h.GetSealabsPay)
	userGroup.POST("/sealab-pay", h.AddSealabsPay)
	userGroup.PATCH("/sealab-pay/:cardNumber", h.PatchSealabsPay)
	userGroup.DELETE("/sealab-pay/:cardNumber", h.DeleteSealabsPay)
	userGroup.POST("/register-merchant", h.RegisterMerchant)
	userGroup.GET("/profile", h.GetUserProfile)
	userGroup.POST("/profile/picture", h.UploadProfilePicture)
	userGroup.POST("/password", h.VerifyPasswordChange)
	userGroup.POST("/verify", h.VerifyOTP)
	userGroup.PATCH("/password", h.ChangePassword)
}
