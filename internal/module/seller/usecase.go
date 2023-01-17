package seller

import (
	"context"
	"murakali/internal/model"
	"murakali/internal/module/seller/delivery/body"
	"murakali/pkg/pagination"
)

type UseCase interface {
	GetOrder(ctx context.Context, userID, orderStatusID string, pgn *pagination.Pagination) (*pagination.Pagination, error)
	ChangeOrderStatus(ctx context.Context, userID string, requestBody body.ChangeOrderStatusRequest) error
	GetOrderByOrderID(ctx context.Context, orderID string) (*model.Order, error)
	GetCourierSeller(ctx context.Context, userID string) (*body.CourierSellerResponse, error)
	GetSellerBySellerID(ctx context.Context, sellerID string) (*body.SellerResponse, error)
	CreateCourierSeller(ctx context.Context, userID string, courierID string) error
	DeleteCourierSellerByID(ctx context.Context, shopCourierID string) error
	GetCategoryBySellerID(ctx context.Context, shopID string) ([]*body.CategoryResponse, error)
	UpdateResiNumberInOrderSeller(ctx context.Context, userID, orderID string, requestBody body.UpdateNoResiOrderSellerRequest) error
	GetAllVoucherSeller(ctx context.Context, userID string, pgn *pagination.Pagination) (*pagination.Pagination, error)
}
