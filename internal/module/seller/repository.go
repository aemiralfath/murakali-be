package seller

import (
	"context"
	"murakali/internal/model"
	"murakali/pkg/pagination"
)

type Repository interface {
	GetTotalOrder(ctx context.Context, userID string) (int64, error)
	GetOrders(ctx context.Context, userID string, pgn *pagination.Pagination) ([]*model.Order, error)
	GetShopIDByUser(ctx context.Context, userID string) (string, error)
}
