package delivery

import (
	"murakali/internal/module/auth"

	"github.com/gin-gonic/gin"
)

func MapAuthRoutes(authGroup *gin.RouterGroup, h auth.Handlers) {
	authGroup.POST("/register", h.RegisterEmail)
	authGroup.PUT("/register", h.RegisterUser)
	authGroup.POST("/verify", h.VerifyOTP)
	authGroup.GET("/verify", h.ResetPasswordVerifyOTP)
	authGroup.POST("/login", h.Login)
	authGroup.GET("/logout", h.Logout)
	authGroup.POST("/reset-password", h.ResetPasswordEmail)
	authGroup.PATCH("/reset-password", h.ResetPasswordUser)
	authGroup.GET("/refresh", h.RefreshToken)
	authGroup.GET("/unique/username/:username", h.CheckUniqueUsername)
	authGroup.GET("/unique/phone-no/:phone_no", h.CheckUniquePhoneNo)
	authGroup.GET("/google-oauth", h.GoogleAuth)
}
