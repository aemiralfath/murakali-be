package user

import "github.com/gin-gonic/gin"

type Handlers interface {
	UserEdit(c *gin.Context)
}
