package product

import (
	"context"
	"murakali/internal/model"
	"murakali/internal/module/product/delivery/body"
	"murakali/pkg/pagination"

	"github.com/google/uuid"
)

type Repository interface {
	GetCategories(ctx context.Context) ([]*model.Category, error)
	GetBanners(ctx context.Context) ([]*model.Banner, error)
	GetCategoriesByName(ctx context.Context, name string) ([]*model.Category, error)
	GetCategoriesByParentID(ctx context.Context, parentID uuid.UUID) ([]*model.Category, error)
	GetRecommendedProducts(ctx context.Context, pgn *pagination.Pagination) ([]*body.Products, []*model.Promotion, []*model.Voucher, error)
	GetTotalProduct(ctx context.Context) (int64, error)
}
