package user

import "github.com/gin-gonic/gin"

type Handlers interface {
	EditUser(c *gin.Context)
	EditEmail(c *gin.Context)
	EditEmailUser(c *gin.Context)
}
