package admin

import (
	"context"
	"murakali/internal/model"
	"murakali/internal/module/admin/delivery/body"
	"murakali/pkg/pagination"
	"murakali/pkg/postgre"
)

type Repository interface {
	CountCodeVoucher(ctx context.Context, code string) (int64, error)
	GetTotalVoucher(ctx context.Context, voucherStatusID string) (int64, error)
	GetTotalRefunds(ctx context.Context) (int64, error)
	GetVoucherByID(ctx context.Context, voucherID string) (*model.Voucher, error)
	GetAllVoucher(ctx context.Context, voucherStatusID, sortFilter string, pgn *pagination.Pagination) ([]*model.Voucher, error)
	GetRefunds(ctx context.Context, sortFilter string, pgn *pagination.Pagination) ([]*model.RefundOrder, error)
	CreateVoucher(ctx context.Context, voucherShop *model.Voucher) error
	UpdateVoucher(ctx context.Context, voucherShop *model.Voucher) error
	DeleteVoucher(ctx context.Context, voucherID string) error
	GetRefundByID(ctx context.Context, refundID string) (*model.Refund, error)
	GetOrderByID(ctx context.Context, orderID string) (*model.OrderModel, error)
	UpdateRefund(ctx context.Context, tx postgre.Transaction, refund *model.Refund) error
	UpdateOrderStatus(ctx context.Context, tx postgre.Transaction, order *model.OrderModel) error
	GetOrderItemsByOrderID(ctx context.Context, tx postgre.Transaction, orderID string) ([]*model.OrderItem, error)
	GetProductDetailByID(ctx context.Context, tx postgre.Transaction, productDetailID string) (*model.ProductDetail, error)
	UpdateProductDetailStock(ctx context.Context, tx postgre.Transaction, productDetailData *model.ProductDetail) error
	GetWalletByUserID(ctx context.Context, tx postgre.Transaction, userID string) (*model.Wallet, error)
	InsertWalletHistory(ctx context.Context, tx postgre.Transaction, walletHistory *model.WalletHistory) error
	UpdateWalletBalance(ctx context.Context, tx postgre.Transaction, wallet *model.Wallet) error
	GetCategories(ctx context.Context) ([]*body.CategoryResponse, error)
	AddCategory(ctx context.Context, requestBody body.CategoryRequest) error
	DeleteCategory(ctx context.Context, categoryID string) error
	EditCategory(ctx context.Context, requestBody body.CategoryRequest) error
	CountProductCategory(ctx context.Context, userid string) (int, error)
	CountCategoryParent(ctx context.Context, userid string) (int, error)
	GetBanner(ctx context.Context) ([]*body.BannerResponse, error)
	AddBanner(ctx context.Context, requestBody body.BannerRequest) error
	DeleteBanner(ctx context.Context, bannerID string) error
	EditBanner(ctx context.Context, requestBody body.BannerIDRequest) error
}
