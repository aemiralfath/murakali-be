package admin

import "github.com/gin-gonic/gin"

type Handlers interface {
	GetAllVoucher(c *gin.Context)
	CreateVoucher(c *gin.Context)
	UpdateVoucher(c *gin.Context)
	DeleteVoucher(c *gin.Context)
	GetDetailVoucher(c *gin.Context)
	GetRefunds(c *gin.Context)
	RefundOrder(c *gin.Context)
}
