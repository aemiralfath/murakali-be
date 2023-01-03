package auth

import "github.com/gin-gonic/gin"

type Handlers interface {
	RegisterEmail(c *gin.Context)
	RegisterUser(c *gin.Context)
	VerifyOTP(c *gin.Context)
	Login(c *gin.Context)
	Logout(c *gin.Context)
	RefreshToken(c *gin.Context)
	ResetPasswordEmail(c *gin.Context)
	ResetPasswordUser(c *gin.Context)
	ResetPasswordVerifyOTP(c *gin.Context)
	CheckUniqueUsername(c *gin.Context)
	CheckUniquePhoneNo(c *gin.Context)
}
