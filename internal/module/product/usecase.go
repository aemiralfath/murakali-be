package product

import (
	"context"
	"murakali/internal/model"
)

type UseCase interface {
	GetCategories(ctx context.Context) ([]*model.Category, error)
}
