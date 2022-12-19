package user

import "github.com/gin-gonic/gin"

type Handlers interface {
	GetAddress(c *gin.Context)
	CreateAddress(c *gin.Context)
	GetAddressByID(c *gin.Context)
	UpdateAddressByID(c *gin.Context)
	DeleteAddressByID(c *gin.Context)
}
