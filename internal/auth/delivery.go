package auth

import "github.com/gin-gonic/gin"

type Handlers interface {
	Register(c *gin.Context)
}
