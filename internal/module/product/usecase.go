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


	GetSearchProducts(ctx context.Context, pgn *pagination.Pagination, query *body.GetSearchProductQueryRequest) (*pagination.Pagination, error) 
}
