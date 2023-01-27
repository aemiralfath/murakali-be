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
	GetCategories(c *gin.Context)
	UploadProductPicture(c *gin.Context)
	AddCategory(c *gin.Context)
	DeleteCategory(c *gin.Context)
	EditCategory(c *gin.Context)
	GetBanner(c *gin.Context)
	AddBanner(c *gin.Context)
	DeleteBanner(c *gin.Context)
	EditBanner(c *gin.Context)
}
