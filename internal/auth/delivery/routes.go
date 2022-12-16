package delivery

import (
	"murakali/internal/auth"

	"github.com/gin-gonic/gin"
)

func MapAuthRoutes(authGroup *gin.RouterGroup, h auth.Handlers) {
	authGroup.POST("/register", h.RegisterEmail)
	authGroup.PUT("/register", h.RegisterUser)
	authGroup.POST("/verify", h.VerifyOTP)
	authGroup.GET("/verify", h.VerifyOTP)
	authGroup.POST("/login", h.Login)
	authGroup.POST("/reset-password", h.ResetPasswordEmail)
	authGroup.PATCH("/reset-password", h.ResetPasswordUser)
	authGroup.GET("/refresh", h.RefreshToken)
}
