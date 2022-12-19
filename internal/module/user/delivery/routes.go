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
	userGroup.GET("/address/:id", h.GetAddressByID)
	userGroup.PUT("/address/:id", h.UpdateAddressByID)
	userGroup.DELETE("/address/:id", h.DeleteAddressByID)
}
