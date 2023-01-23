package admin

import (
	"context"
	"murakali/internal/model"
	"murakali/internal/module/admin/delivery/body"
	"murakali/pkg/pagination"
)

type UseCase interface {
	GetAllVoucher(ctx context.Context, voucherStatusID, sortFilter string, pgn *pagination.Pagination) (*pagination.Pagination, error)
	CreateVoucher(ctx context.Context, requestBody body.CreateVoucherRequest) error
	UpdateVoucher(ctx context.Context, requestBody body.UpdateVoucherRequest) error
	GetDetailVoucher(ctx context.Context, voucherID string) (*model.Voucher, error)
	DeleteVoucher(ctx context.Context, voucherID string) error
	GetRefunds(ctx context.Context, sortFilter string, pgn *pagination.Pagination) (*pagination.Pagination, error)
	RefundOrder(ctx context.Context, refundID string) error
}
