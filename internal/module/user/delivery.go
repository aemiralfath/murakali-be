package user

import (
	"github.com/gin-gonic/gin"
)

type Handlers interface {
	GetSealabsPay(c *gin.Context)
	AddSealabsPay(c *gin.Context)
	PatchSealabsPay(c *gin.Context)
	DeleteSealabsPay(c *gin.Context)
	EditUser(c *gin.Context)
	EditEmail(c *gin.Context)
	EditEmailUser(c *gin.Context)
}
