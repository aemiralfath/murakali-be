package delivery

import (
	"murakali/internal/middleware"
	"murakali/internal/module/seller"

	"github.com/gin-gonic/gin"
)

func MapSellerRoutes(sellerGroup *gin.RouterGroup, h seller.Handlers, mw *middleware.MWManager) {
	sellerGroup.GET("/:seller_id", h.GetSellerBySellerID)
	sellerGroup.GET("/:seller_id/category", h.GetCategoryBySellerID)

	sellerGroup.Use(mw.AuthJWTMiddleware())
	sellerGroup.Use(mw.SellerJWTMiddleware())
	sellerGroup.GET("/order", h.GetOrder)
	sellerGroup.GET("/order/:order_id", h.GetOrderByOrderID)
	sellerGroup.PATCH("/order-status", h.ChangeOrderStatus)
	sellerGroup.GET("/courier", h.GetCourierSeller)
	sellerGroup.POST("/courier", h.CreateCourierSeller)
	sellerGroup.DELETE("/courier/:id", h.DeleteCourierSellerByID)
	sellerGroup.PATCH("/order-resi/:id", h.UpdateResiNumberInOrderSeller)
	sellerGroup.GET("/voucher", h.GetAllVoucherSeller)
}
