package user

import "github.com/gin-gonic/gin"

type Handlers interface {
	GetAddress(c *gin.Context)
	CreateAddress(c *gin.Context)
}
