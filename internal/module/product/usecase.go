package product

import (
	"context"
	"murakali/internal/model"
	"murakali/internal/module/product/delivery/body"
	"murakali/pkg/pagination"
)

type UseCase interface {
	GetCategories(ctx context.Context) ([]*body.CategoryResponse, error)
	GetBanners(ctx context.Context) ([]*model.Banner, error)
	GetCategoriesByName(ctx context.Context, name string) ([]*body.CategoryResponse, error)
	GetRecommendedProducts(ctx context.Context, pgn *pagination.Pagination) (*pagination.Pagination, error)
	GetProductDetail(ctx context.Context, productID string) (*body.ProductDetailResponse, error)
	GetAllProductImage(ctx context.Context, productID string) ([]*body.GetImageResponse, error)
	GetProducts(ctx context.Context, pgn *pagination.Pagination, query *body.GetProductQueryRequest) (*pagination.Pagination, error)
	GetFavoriteProducts(ctx context.Context, pgn *pagination.Pagination, query *body.GetProductQueryRequest,
		userID string) (*pagination.Pagination, error)
	CheckProductIsFavorite(ctx context.Context, userID, productID string) bool
	CountSpecificFavoriteProduct(
		ctx context.Context, productID string) (int64, error)
	CreateFavoriteProduct(ctx context.Context, productID, userID string) error
	DeleteFavoriteProduct(ctx context.Context, productID, userID string) error
	GetProductReviews(ctx context.Context, pgn *pagination.Pagination, productID string,
		query *body.GetReviewQueryRequest) (*pagination.Pagination, error)
	GetTotalReviewRatingByProductID(ctx context.Context, productID string) (*body.AllRatingProduct, error)
	CreateProductReview(ctx context.Context, reqBody body.ReviewProductRequest, userID string) error
	DeleteProductReview(ctx context.Context, reviewID, userID string) error
	CreateProduct(ctx context.Context, requestBody body.CreateProductRequest, userID string) error
	UpdateListedStatus(ctx context.Context, productID string) error
	UpdateProductListedStatusBulk(ctx context.Context, product body.UpdateProductListedStatusBulkRequest) error
	UpdateProduct(ctx context.Context, requestBody body.UpdateProductRequest, userID, productID string) error
	UpdateProductMetadata(ctx context.Context) error
}
