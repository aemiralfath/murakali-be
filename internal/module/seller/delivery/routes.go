package delivery

import (
	"murakali/internal/middleware"
	"murakali/internal/module/seller"

	"github.com/gin-gonic/gin"
)

func MapSellerRoutes(sellerGroup *gin.RouterGroup, h seller.Handlers, mw *middleware.MWManager) {
	sellerGroup.Use(mw.AuthJWTMiddleware())
	sellerGroup.GET("/order", h.GetOrder)
	sellerGroup.PATCH("/orderstatus", h.ChangeOrderStatus)

}
