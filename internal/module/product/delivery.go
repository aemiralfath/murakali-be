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
	GetAllProductImage(c *gin.Context)
	GetFavoriteProducts(c *gin.Context)
	CheckProductIsFavorite(c *gin.Context)
	CountSpecificFavoriteProduct(c *gin.Context)
	CreateFavoriteProduct(c *gin.Context)
	DeleteFavoriteProduct(c *gin.Context)
	GetProductReviews(c *gin.Context)
	GetTotalReviewRatingByProductID(c *gin.Context)
	CreateProductReview(c *gin.Context)
	DeleteProductReview(c *gin.Context)
	CreateProduct(c *gin.Context)
	UpdateListedStatus(c *gin.Context)
	UpdateListedStatusBulk(c *gin.Context)
	UpdateProduct(c *gin.Context)
	UploadProductPicture(c *gin.Context)
}
