package cart

import "github.com/gin-gonic/gin"

type Handlers interface {
	GetCartHoverHome(c *gin.Context)
	GetCartItems(c *gin.Context)
	AddCartItems(c *gin.Context)
	UpdateCartItems(c *gin.Context)
	DeleteCartItems(c *gin.Context)
	GetAllVoucher(c *gin.Context)
}
