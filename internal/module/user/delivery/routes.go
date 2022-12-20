package delivery

import (
	"murakali/internal/middleware"
	"murakali/internal/module/user"

	"github.com/gin-gonic/gin"
)

func MapUserRoutes(userGroup *gin.RouterGroup, h user.Handlers, mw *middleware.MWManager) {
	userGroup.Use(mw.AuthJWTMiddleware())
	userGroup.GET("/sealab-pay", h.GetSealabsPay)
	userGroup.POST("/sealab-pay", h.AddSealabsPay)
	userGroup.PATCH("/sealab-pay", h.PatchSealabsPay)
	userGroup.DELETE("/sealab-pay", h.DeleteSealabsPay)

	userGroup.PUT("/profile", h.EditUser)
	userGroup.POST("/email", h.EditEmail)
	userGroup.GET("/email", h.EditEmailUser)
}
