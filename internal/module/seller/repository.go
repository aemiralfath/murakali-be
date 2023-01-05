package seller

import (
	"context"
	"murakali/internal/model"
	"murakali/internal/module/seller/delivery/body"
	"murakali/pkg/pagination"
)

type Repository interface {
	GetTotalOrder(ctx context.Context, userID, orderStatusID string) (int64, error)
	GetOrders(ctx context.Context, userID, orderStatusID string, pgn *pagination.Pagination) ([]*model.Order, error)
	GetShopIDByUser(ctx context.Context, userID string) (string, error)
	GetShopIDByOrder(ctx context.Context, OrderID string) (string, error)
	ChangeOrderStatus(ctx context.Context, requestBody body.ChangeOrderStatusRequest) error
	GetOrderByOrderID(ctx context.Context, OrderID string) (*model.Order, error)
	GetCourierSeller(ctx context.Context, userID string) ([]*body.CourierSellerInfo, error)
	GetSellerBySellerID(ctx context.Context, sellerID string) (*body.SellerResponse, error)
}
