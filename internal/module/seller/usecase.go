package seller

import (
	"context"
	"murakali/internal/model"
	"murakali/internal/module/seller/delivery/body"
	"murakali/pkg/pagination"
)

type UseCase interface {
	GetOrder(ctx context.Context, userID string, pgn *pagination.Pagination) (*pagination.Pagination, error)
	ChangeOrderStatus(ctx context.Context, userID string, requestBody body.ChangeOrderStatusRequest) error
	GetOrderByOrderID(ctx context.Context, orderID string) (*model.Order, error)
}
