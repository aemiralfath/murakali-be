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
	GetSellerBySellerID(ctx context.Context, sellerID string) (*body.SellerResponse, error)
	GetCourierByID(ctx context.Context, courierID string) (string, error)
	GetCourierSellerNotNullByShopAndCourierID(ctx context.Context, shopID, courierID string) (string, error)
	GetShopIDByUserID(ctx context.Context, userID string) (string, error)
	CreateCourierSeller(ctx context.Context, shopID string, courierID string) error
	GetCourierSellerByID(ctx context.Context, shopCourierID string) (string, error)
	UpdateCourierSellerByID(ctx context.Context, shopID, courierID string) error
	DeleteCourierSellerByID(ctx context.Context, shopCourierID string) error
	GetCategoryBySellerID(ctx context.Context, shopID string) ([]*body.CategoryResponse, error)
	GetAllCourier(ctx context.Context) ([]*body.CourierInfo, error)
	GetCourierSeller(ctx context.Context, userID string) ([]*body.CourierSellerRelationInfo, error)
	GetBuyerIDByOrderID(ctx context.Context, orderID string) (string, error)
	GetSellerIDByOrderID(ctx context.Context, orderID string) (string, error)
	GetAddressByBuyerID(ctx context.Context, userID string) (*model.Address, error)
	GetAddressBySellerID(ctx context.Context, userID string) (*model.Address, error)
	UpdateResiNumberInOrderSeller(ctx context.Context, noResi, orderID, shopID string) error
}
