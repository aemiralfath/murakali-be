package product

import "github.com/gin-gonic/gin"

type Handlers interface {
	GetProducts(c *gin.Context)
	GetCategories(c *gin.Context)
	GetBanners(c *gin.Context)
	GetCategoriesByNameLevelOne(c *gin.Context)
	GetCategoriesByNameLevelTwo(c *gin.Context)
	GetCategoriesByNameLevelThree(c *gin.Context)
	GetRecommendedProducts(c *gin.Context)
	GetProductDetail(c *gin.Context)
	GetFavoriteProducts(c *gin.Context)
}
