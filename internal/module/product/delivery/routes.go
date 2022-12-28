package delivery

import (
	"murakali/internal/module/product"

	"github.com/gin-gonic/gin"
)

func MapProductRoutes(productGroup *gin.RouterGroup, h product.Handlers) {
	productGroup.GET("/category", h.GetCategories)
	productGroup.GET("/banner", h.GetBanners)
	productGroup.GET("/category/:name_lvl_one", h.GetCategoriesByNameLevelOne)
	productGroup.GET("/category/:name_lvl_one/:name_lvl_two", h.GetCategoriesByNameLevelTwo)
	productGroup.GET("/category/:name_lvl_one/:name_lvl_two/:name_lvl_three", h.GetCategoriesByNameLevelThree)
	productGroup.GET("/recommended", h.GetRecommendedProducts)
	productGroup.GET("/:product_id", h.GetProductDetail)
}
