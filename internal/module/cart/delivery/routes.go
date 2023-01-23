package delivery

import (
	"murakali/internal/middleware"
	"murakali/internal/module/cart"

	"github.com/gin-gonic/gin"
)

func MapCartRoutes(cartGroup *gin.RouterGroup, h cart.Handlers, mw *middleware.MWManager) {
	cartGroup.Use(mw.AuthJWTMiddleware())
	cartGroup.GET("/hover-home", h.GetCartHoverHome)
	cartGroup.GET("/items", h.GetCartItems)
	cartGroup.POST("/items", h.AddCartItems)
	cartGroup.PUT("/items", h.UpdateCartItems)
	cartGroup.DELETE("/items/:id", h.DeleteCartItems)

	cartGroup.GET("/voucher/:shop_id", h.GetAllVoucher)
}
