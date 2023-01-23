package admin

import (
	"context"
	"murakali/internal/model"
	"murakali/pkg/pagination"
	"murakali/pkg/postgre"
)

type Repository interface {
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
	GetWalletByUserID(ctx context.Context, userID string) (*model.Wallet, error)
	InsertWalletHistory(ctx context.Context, tx postgre.Transaction, walletHistory *model.WalletHistory) error

	UpdateWalletBalance(ctx context.Context, tx postgre.Transaction, wallet *model.Wallet) error
}
