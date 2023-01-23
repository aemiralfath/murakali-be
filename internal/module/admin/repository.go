package admin

import (
	"context"
	"murakali/internal/model"
	"murakali/pkg/pagination"
)

type Repository interface {
	CountCodeVoucher(ctx context.Context, code string) (int64, error)
	GetTotalVoucher(ctx context.Context, voucherStatusID string) (int64, error)
	GetVoucherByID(ctx context.Context, voucherID string) (*model.Voucher, error)
	GetAllVoucher(ctx context.Context, voucherStatusID, sortFilter string, pgn *pagination.Pagination) ([]*model.Voucher, error)
	CreateVoucher(ctx context.Context, voucherShop *model.Voucher) error
	UpdateVoucher(ctx context.Context, voucherShop *model.Voucher) error
	DeleteVoucher(ctx context.Context, voucherID string) error
}
