package product

import (
	"context"
	"murakali/internal/model"
	"murakali/internal/module/product/delivery/body"
)

type UseCase interface {
	GetCategories(ctx context.Context) ([]*body.CategoryResponse, error)
	GetBanners(ctx context.Context) ([]*model.Banner, error)
	GetCategoriesByName(ctx context.Context, name string) ([]*body.CategoryResponse, error)
}
