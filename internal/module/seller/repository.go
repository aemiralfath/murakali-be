package seller

import (
	"context"
	"murakali/internal/model"
	"murakali/internal/module/seller/delivery/body"
	"murakali/pkg/pagination"
)

type Repository interface {
	GetTotalOrder(ctx context.Context, userID string) (int64, error)
	GetOrders(ctx context.Context, userID string, pgn *pagination.Pagination) ([]*model.Order, error)
	GetShopIDByUser(ctx context.Context, userID string) (string, error)
	GetShopIDByOrder(ctx context.Context, OrderID string) (string, error)
	ChangeOrderStatus(ctx context.Context, requestBody body.ChangeOrderStatusRequest) error
}
