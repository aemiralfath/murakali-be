package user

import (
	"github.com/gin-gonic/gin"
)

type Handlers interface {
	GetSealabsPay(c *gin.Context)
	AddSealabsPay(c *gin.Context)
}
