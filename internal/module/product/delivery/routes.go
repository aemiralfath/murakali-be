package delivery

import (
	"murakali/internal/middleware"
	"murakali/internal/module/product"

	"github.com/gin-gonic/gin"
)

func MapProductRoutes(productGroup *gin.RouterGroup, h product.Handlers, mw *middleware.MWManager) {
	productGroup.GET("/category", h.GetCategories)
	productGroup.GET("/banner", h.GetBanners)
	productGroup.GET("/category/:name_lvl_one", h.GetCategoriesByNameLevelOne)
	productGroup.GET("/category/:name_lvl_one/:name_lvl_two", h.GetCategoriesByNameLevelTwo)
	productGroup.GET("/category/:name_lvl_one/:name_lvl_two/:name_lvl_three", h.GetCategoriesByNameLevelThree)
	productGroup.GET("/recommended", h.GetRecommendedProducts)
	productGroup.GET("/:product_id", h.GetProductDetail)
	productGroup.GET("/:product_id/review", h.GetProductReviews)
	productGroup.GET("/:product_id/review/rating", h.GetTotalReviewRatingByProductID)
	productGroup.GET("/", h.GetProducts)
	productUserGroup := productGroup.Group("/")
	productUserGroup.Use(mw.AuthJWTMiddleware())
	productUserGroup.GET("/favorite", h.GetFavoriteProducts)
	productUserGroup.POST("/", h.CreateProduct)
	productUserGroup.PUT("/:id", h.UpdateListedStatus)
}
