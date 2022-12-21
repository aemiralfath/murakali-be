package product

import (
	"context"
	"murakali/internal/model"
)

type Repository interface {
	GetCategories(ctx context.Context) ([]*model.Category, error)
}
