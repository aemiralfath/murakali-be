package usecase

import (
	"context"
	"database/sql"
	"math"
	"murakali/config"
	"murakali/internal/model"
	"murakali/internal/module/admin"
	"murakali/internal/module/admin/delivery/body"
	"murakali/pkg/httperror"
	"murakali/pkg/pagination"
	"murakali/pkg/postgre"
	"net/http"
)

type adminUC struct {
	cfg       *config.Config
	txRepo    *postgre.TxRepo
	adminRepo admin.Repository
}

func NewAdminUseCase(cfg *config.Config, txRepo *postgre.TxRepo, adminRepo admin.Repository) admin.UseCase {
	return &adminUC{cfg: cfg, txRepo: txRepo, adminRepo: adminRepo}
}

func (u *adminUC) GetAllVoucher(ctx context.Context, voucherStatusID, sortFilter string, pgn *pagination.Pagination) (*pagination.Pagination, error) {
	totalRows, err := u.adminRepo.GetTotalVoucher(ctx, voucherStatusID)
	if err != nil {
		return nil, err
	}

	totalPages := int(math.Ceil(float64(totalRows) / float64(pgn.Limit)))
	pgn.TotalRows = totalRows
	pgn.TotalPages = totalPages

	ShopVouchers, err := u.adminRepo.GetAllVoucher(ctx, voucherStatusID, sortFilter, pgn)
	if err != nil {
		return nil, err
	}

	pgn.Rows = ShopVouchers

	return pgn, nil
}

func (u *adminUC) CreateVoucher(ctx context.Context, requestBody body.CreateVoucherRequest) error {
	voucherShop := &model.Voucher{
		Code:               requestBody.Code,
		Quota:              requestBody.Quota,
		ActivedDate:        requestBody.ActiveDateTime,
		ExpiredDate:        requestBody.ExpiredDateTime,
		DiscountPercentage: &requestBody.DiscountPercentage,
		DiscountFixPrice:   &requestBody.DiscountFixPrice,
		MinProductPrice:    &requestBody.MinProductPrice,
		MaxDiscountPrice:   &requestBody.MaxDiscountPrice,
	}

	err := u.adminRepo.CreateVoucher(ctx, voucherShop)
	if err != nil {
		return err
	}

	return nil
}

func (u *adminUC) UpdateVoucher(ctx context.Context, requestBody body.UpdateVoucherRequest) error {
	voucherShop, errVoucher := u.adminRepo.GetVoucherByID(ctx, requestBody.VoucherID)
	if errVoucher != nil {
		if errVoucher == sql.ErrNoRows {
			return httperror.New(http.StatusBadRequest, body.VoucherSellerNotFoundMessage)
		}

		return errVoucher
	}

	voucherShop.Quota = requestBody.Quota
	voucherShop.ActivedDate = requestBody.ActiveDateTime
	voucherShop.ExpiredDate = requestBody.ExpiredDateTime
	voucherShop.DiscountPercentage = &requestBody.DiscountPercentage
	voucherShop.DiscountFixPrice = &requestBody.DiscountFixPrice
	voucherShop.MinProductPrice = &requestBody.MinProductPrice
	voucherShop.MaxDiscountPrice = &requestBody.MaxDiscountPrice

	err := u.adminRepo.UpdateVoucher(ctx, voucherShop)
	if err != nil {
		return err
	}

	return nil
}

func (u *adminUC) GetDetailVoucher(ctx context.Context, voucherID string) (*model.Voucher, error) {
	voucherShop, errVoucher := u.adminRepo.GetVoucherByID(ctx, voucherID)
	if errVoucher != nil {
		if errVoucher == sql.ErrNoRows {
			return nil, httperror.New(http.StatusBadRequest, body.VoucherSellerNotFoundMessage)
		}

		return nil, errVoucher
	}

	return voucherShop, nil
}

func (u *adminUC) DeleteVoucher(ctx context.Context, voucherID string) error {

	_, errVoucher := u.adminRepo.GetVoucherByID(ctx, voucherID)
	if errVoucher != nil {
		if errVoucher == sql.ErrNoRows {
			return httperror.New(http.StatusBadRequest, body.VoucherSellerNotFoundMessage)
		}

		return errVoucher
	}

	if err := u.adminRepo.DeleteVoucher(ctx, voucherID); err != nil {
		return err
	}

	return nil
}
