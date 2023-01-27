package delivery

import (
	"murakali/internal/middleware"
	"murakali/internal/module/admin"

	"github.com/gin-gonic/gin"
)

func MapAdminRoutes(adminGroup *gin.RouterGroup, h admin.Handlers, mw *middleware.MWManager) {
	adminGroup.GET("/banner", h.GetBanner)
	adminGroup.Use(mw.AuthJWTMiddleware())
	adminGroup.Use(mw.AdminJWTMiddleware())
	adminGroup.GET("/voucher", h.GetAllVoucher)
	adminGroup.POST("/voucher", h.CreateVoucher)
	adminGroup.PUT("/voucher", h.UpdateVoucher)
	adminGroup.GET("/voucher/:id", h.GetDetailVoucher)
	adminGroup.DELETE("/voucher/:id", h.DeleteVoucher)

	adminGroup.GET("/refund", h.GetRefunds)
	adminGroup.POST("/refund/:id", h.RefundOrder)

	adminGroup.GET("/category", h.GetCategories)
	adminGroup.POST("/category", h.AddCategory)
	adminGroup.PUT("/category", h.EditCategory)
	adminGroup.DELETE("/category/:id", h.DeleteCategory)

	adminGroup.POST("/banner", h.AddBanner)
	adminGroup.PUT("/banner", h.EditBanner)
	adminGroup.DELETE("/banner/:id", h.DeleteBanner)

	adminGroup.POST("/picture", h.UploadProductPicture)
}
