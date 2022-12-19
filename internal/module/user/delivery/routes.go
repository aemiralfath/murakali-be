package delivery

import (
	"github.com/gin-gonic/gin"
	"murakali/internal/middleware"
	"murakali/internal/module/user"
)

func MapUserRoutes(userGroup *gin.RouterGroup, h user.Handlers, mw *middleware.MWManager) {
	userGroup.Use(mw.AuthJWTMiddleware())
	userGroup.GET("/address", h.GetAddress)
	userGroup.POST("/address", h.CreateAddress)
	userGroup.PUT("/address", h.UpdateAddress)
	userGroup.DELETE("/address", h.DeleteAddress)
}
