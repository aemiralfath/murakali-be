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
	productGroup.GET("/:product_id/picture", h.GetAllProductImage)
	productGroup.GET("/:product_id/review", h.GetProductReviews)
	productGroup.GET("/:product_id/review/rating", h.GetTotalReviewRatingByProductID)
	productGroup.GET("/", h.GetProducts)
	productGroup.Use(mw.AuthJWTMiddleware())
	productGroup.GET("/favorite", h.GetFavoriteProducts)
	productGroup.POST("/favorite/check", h.CheckProductIsFavorite)
	productGroup.POST("/favorite", h.CreateFavoriteProduct)
	productGroup.DELETE("/favorite", h.DeleteFavoriteProduct)
	productGroup.DELETE("/review/:review_id", h.DeleteProductReview)
	productGroup.POST("/:product_id/review", h.CreateProductReview)
	productGroup.Use(mw.SellerJWTMiddleware())
	productGroup.POST("/", h.CreateProduct)
	productGroup.POST("/picture", h.UploadProductPicture)
	productGroup.PUT("/status/:id", h.UpdateListedStatus)
	productGroup.PATCH("/bulk-status", h.UpdateListedStatusBulk)
	productGroup.PUT("/:id", h.UpdateProduct)
}
