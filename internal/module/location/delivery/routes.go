package delivery

import (
	"github.com/gin-gonic/gin"
	"murakali/internal/module/location"
)

func MapAuthRoutes(locationGroup *gin.RouterGroup, h location.Handlers) {
	locationGroup.GET("/province", h.GetProvince)
}
