package product

import (
	"context"
	"murakali/internal/model"
	"murakali/pkg/pagination"

	"github.com/google/uuid"
)

type Repository interface {
	GetCategories(ctx context.Context) ([]*model.Category, error)
	GetBanners(ctx context.Context) ([]*model.Banner, error)
	GetCategoriesByName(ctx context.Context, name string) ([]*model.Category, error)
	GetCategoriesByParentID(ctx context.Context, parentID uuid.UUID) ([]*model.Category, error)
	GetRecommendedProducts(ctx context.Context, pgn *pagination.Pagination) ([]*model.Product, []*model.Promotion, []*model.Voucher, []*model.Shop, []*model.Category, error)
	GetTotalProduct(ctx context.Context) (int64, error)
}
