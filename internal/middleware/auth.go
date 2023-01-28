package middleware

import (
	"murakali/pkg/jwt"
	"murakali/pkg/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (mw *MWManager) AuthJWTMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		claim, err := jwt.ExtractJWTFromRequest(c.Request, mw.cfg.JWT.JwtSecretKey)
		if err != nil {
			response.ErrorResponse(c.Writer, response.ForbiddenMessage, http.StatusForbidden)
			c.Abort()
			return
		}

		if claim["role_id"] == nil {
			response.ErrorResponse(c.Writer, response.ForbiddenMessage, http.StatusForbidden)
			c.Abort()
			return
		}

		mw.log.Infof("body middleware bearerHeader %s", claim["id"].(string))
		c.Set("userID", claim["id"].(string))
		c.Set("roleID", claim["role_id"].(float64))
		c.Next()
	}
}
