package cart

import "github.com/gin-gonic/gin"

type Handlers interface {
	GetCartHoverHome(c *gin.Context)
	GetCartItems(c *gin.Context)
}
