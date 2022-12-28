package seller

import "github.com/gin-gonic/gin"

type Handlers interface {
	GetOrder(c *gin.Context)
}
