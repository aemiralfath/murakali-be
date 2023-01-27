package delivery

import (
	"murakali/internal/middleware"
	"murakali/internal/module/seller"

	"github.com/gin-gonic/gin"
)

func MapSellerRoutes(sellerGroup *gin.RouterGroup, h seller.Handlers, mw *middleware.MWManager) {
	sellerGroup.GET("/:seller_id", h.GetSellerBySellerID)
	sellerGroup.GET("/:seller_id/category", h.GetCategoryBySellerID)
	sellerGroup.POST("/delivery", h.UpdateOnDeliveryOrder)
	sellerGroup.POST("/expired", h.UpdateExpiredAtOrder)

	sellerGroup.Use(mw.AuthJWTMiddleware())
	sellerGroup.Use(mw.SellerJWTMiddleware())
	sellerGroup.GET("/performance", h.GetPerformance)
	sellerGroup.GET("/information", h.GetSellerDetailInformation)
	sellerGroup.PATCH("/information", h.UpdateSellerInformation)
	sellerGroup.GET("/user/:user_id", h.GetSellerByUserID)
	sellerGroup.GET("/order", h.GetOrder)
	sellerGroup.GET("/order/:order_id", h.GetOrderByOrderID)
	sellerGroup.PATCH("/order-status", h.ChangeOrderStatus)
	sellerGroup.PATCH("/order-cancel", h.CancelOrderStatus)
	sellerGroup.GET("/courier", h.GetCourierSeller)
	sellerGroup.POST("/courier", h.CreateCourierSeller)
	sellerGroup.DELETE("/courier/:id", h.DeleteCourierSellerByID)
	sellerGroup.PATCH("/order-resi/:id", h.UpdateResiNumberInOrderSeller)
	sellerGroup.POST("/withdrawal/:id", h.WithdrawalOrderBalance)
	sellerGroup.GET("/voucher", h.GetAllVoucherSeller)
	sellerGroup.POST("/voucher", h.CreateVoucherSeller)
	sellerGroup.PUT("/voucher", h.UpdateVoucherSeller)
	sellerGroup.GET("/voucher/:id", h.DetailVoucherSeller)
	sellerGroup.DELETE("/voucher/:id", h.DeleteVoucherSeller)
	sellerGroup.GET("/product/without-promotion", h.GetProductWithoutPromotionSeller)
	sellerGroup.GET("/promotion", h.GetAllPromotionSeller)
	sellerGroup.POST("/promotion", h.CreatePromotionSeller)
	sellerGroup.PUT("/promotion", h.UpdatePromotionSeller)
	sellerGroup.GET("/promotion/:id", h.GetDetailPromotionSellerByID)
	sellerGroup.GET("/refund/:refund_id", h.GetRefundOrderSeller)
	sellerGroup.POST("/refund-thread", h.CreateRefundThreadSeller)
	sellerGroup.PATCH("/refund-accept", h.UpdateRefundAccept)
	sellerGroup.PATCH("/refund-reject", h.UpdateRefundReject)
}
