package product

import (
	"context"
	"murakali/internal/module/product/delivery/body"
)

type UseCase interface {
	GetCategories(ctx context.Context) ([]*body.CategoryResponse, error)
	GetCategoriesByName(ctx context.Context, name string) ([]*body.CategoryResponse, error)
}
