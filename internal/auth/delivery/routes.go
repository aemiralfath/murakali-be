package delivery

import (
	"github.com/gin-gonic/gin"
	"murakali/internal/auth"
)

func MapAuthRoutes(authGroup *gin.RouterGroup, h auth.Handlers) {
	authGroup.POST("/register", h.Register)
	authGroup.PATCH("/verify", h.VerifyOTP)
}
