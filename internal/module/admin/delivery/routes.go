package delivery

import (
	"murakali/internal/middleware"
	"murakali/internal/module/admin"

	"github.com/gin-gonic/gin"
)

func MapAdminRoutes(adminGroup *gin.RouterGroup, h admin.Handlers, mw *middleware.MWManager) {

	adminGroup.Use(mw.AuthJWTMiddleware())
	adminGroup.Use(mw.AdminJWTMiddleware())

	adminGroup.GET("/voucher", h.GetAllVoucher)
	adminGroup.POST("/voucher", h.CreateVoucher)
	adminGroup.PUT("/voucher", h.UpdateVoucher)
	adminGroup.GET("/voucher/:id", h.GetDetailVoucher)
	adminGroup.DELETE("/voucher/:id", h.DeleteVoucher)

}
