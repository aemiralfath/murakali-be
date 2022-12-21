package product

import "github.com/gin-gonic/gin"

type Handlers interface {
	GetCategories(c *gin.Context)
}
