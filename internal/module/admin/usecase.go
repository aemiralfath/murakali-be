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
	GetCategories(ctx context.Context) ([]*body.CategoryResponse, error)
	AddCategory(ctx context.Context, requestBody body.CategoryRequest) error
	DeleteCategory(ctx context.Context, categoryID string) error
	EditCategory(ctx context.Context, requestBody body.CategoryRequest) error
	GetBanner(ctx context.Context) ([]*body.BannerResponse, error)
	AddBanner(ctx context.Context, requestBody body.BannerRequest) error
	DeleteBanner(ctx context.Context, bannerID string) error
}
