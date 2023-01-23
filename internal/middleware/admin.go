package middleware

import (
	"murakali/internal/constant"
	"murakali/pkg/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (mw *MWManager) AdminJWTMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		roleID, exist := c.Get("roleID")
		if !exist || roleID.(float64) != constant.RoleAdmin {
			response.ErrorResponse(c.Writer, response.ForbiddenMessage, http.StatusForbidden)
			c.Abort()
			return
		}
		c.Next()
	}
}
