package delivery

import (
	"murakali/internal/module/product"

	"github.com/gin-gonic/gin"
)

func MapProductRoutes(productGroup *gin.RouterGroup, h product.Handlers) {
	productGroup.GET("/category", h.GetCategories)
}
