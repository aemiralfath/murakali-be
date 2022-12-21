package delivery

import (
	"murakali/internal/middleware"
	"murakali/internal/module/cart"
	"murakali/pkg/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

func MapCartRoutes(cartGroup *gin.RouterGroup, h cart.Handlers, mw *middleware.MWManager) {
	// cartGroup.Use(mw.AuthJWTMiddleware())
	cartGroup.GET("/", func(ctx *gin.Context) {
		response.SuccessResponse(ctx.Writer, nil, http.StatusOK)
	})
}
