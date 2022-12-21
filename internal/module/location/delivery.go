package location

import "github.com/gin-gonic/gin"

type Handlers interface {
	GetProvince(c *gin.Context)
	GetCity(c *gin.Context)
}
