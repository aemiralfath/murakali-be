package admin

import (
	"context"
	"murakali/internal/model"
	"murakali/pkg/pagination"
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
}
