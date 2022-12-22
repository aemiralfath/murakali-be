package delivery

import (
	"murakali/internal/middleware"
	"murakali/internal/module/cart"

	"github.com/gin-gonic/gin"
)

func MapCartRoutes(cartGroup *gin.RouterGroup, h cart.Handlers, mw *middleware.MWManager) {
	cartGroup.Use(mw.AuthJWTMiddleware())
	cartGroup.GET("/hover-home", h.GetCartHoverHome)
}
