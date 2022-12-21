package product

import (
	"context"
	"murakali/internal/model"

	"github.com/google/uuid"
)

type Repository interface {
	GetCategories(ctx context.Context) ([]*model.Category, error)
	GetCategoriesByName(ctx context.Context, name string) ([]*model.Category, error)
	GetCategoriesByParentID(ctx context.Context, parentID uuid.UUID) ([]*model.Category, error)
}
