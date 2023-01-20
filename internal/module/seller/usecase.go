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
	GetSellerByUserID(ctx context.Context, userID string) (*body.SellerResponse, error)
	CreateCourierSeller(ctx context.Context, userID string, courierID string) error
	DeleteCourierSellerByID(ctx context.Context, shopCourierID string) error
	GetCategoryBySellerID(ctx context.Context, shopID string) ([]*body.CategoryResponse, error)
	UpdateResiNumberInOrderSeller(ctx context.Context, userID, orderID string, requestBody body.UpdateNoResiOrderSellerRequest) error
	UpdateOnDeliveryOrder(ctx context.Context) error
	UpdateExpiredAtOrder(ctx context.Context) error
	GetAllVoucherSeller(ctx context.Context, userID string, pgn *pagination.Pagination) (*pagination.Pagination, error)
	CreateVoucherSeller(ctx context.Context, userID string, requestBody body.CreateVoucherRequest) error
	UpdateVoucherSeller(ctx context.Context, userID string, requestBody body.UpdateVoucherRequest) error
	GetDetailVoucherSeller(ctx context.Context, voucherIDShopID *body.VoucherIDShopID) (*model.Voucher, error)
	DeleteVoucherSeller(ctx context.Context, voucherIDShopID *body.VoucherIDShopID) error
	GetAllPromotionSeller(ctx context.Context, userID string, promoStatusID string, pgn *pagination.Pagination) (*pagination.Pagination, error)
	CreatePromotionSeller(ctx context.Context, userID string, requestBody body.CreatePromotionRequest) (int, error)
	UpdatePromotionSeller(ctx context.Context, userID string, requestBody body.UpdatePromotionRequest) error
	GetDetailPromotionSellerByID(ctx context.Context, shopProductPromo *body.ShopProductPromo) (*body.PromotionDetailSeller, error)
}
