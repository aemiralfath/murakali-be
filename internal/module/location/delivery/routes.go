package delivery

import (
	"github.com/gin-gonic/gin"
	"murakali/internal/module/location"
)

func MapAuthRoutes(locationGroup *gin.RouterGroup, h location.Handlers) {
	locationGroup.GET("/province", h.GetProvince)
	locationGroup.GET("/province/city", h.GetCity)
	locationGroup.GET("/province/city/subdistrict", h.GetSubDistrict)
}
