package delivery

import (
	"github.com/gin-gonic/gin"
	"murakali/internal/auth"
)

func MapAuthRoutes(authGroup *gin.RouterGroup, h auth.Handlers) {
	authGroup.POST("/register", h.RegisterEmail)
	authGroup.PUT("/register", h.RegisterUser)
	authGroup.POST("/verify", h.VerifyOTP)
	authGroup.POST("/login", h.Login)
}
