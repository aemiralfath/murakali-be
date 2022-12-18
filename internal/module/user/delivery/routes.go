package delivery

import (
	"github.com/gin-gonic/gin"
	"murakali/internal/middleware"
	"murakali/internal/module/user"
	"murakali/pkg/response"
	"net/http"
)

func MapUserRoutes(userGroup *gin.RouterGroup, h user.Handlers, mw *middleware.MWManager) {
	userGroup.Use(mw.AuthJWTMiddleware())

	// TODO delete this endpoint later
	userGroup.GET("/", func(context *gin.Context) {
		response.SuccessResponse(context.Writer, nil, http.StatusOK)
	})
}
