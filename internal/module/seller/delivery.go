package seller

import "github.com/gin-gonic/gin"

type Handlers interface {
	GetOrder(c *gin.Context)
	ChangeOrderStatus(c *gin.Context)
	GetOrderByOrderID(c *gin.Context)
	GetCourierSeller(c *gin.Context)
	GetSellerBySellerID(c *gin.Context)
	GetCategoryBySellerID(c *gin.Context)
}
