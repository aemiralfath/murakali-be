package usecase

import (
	"context"
	"database/sql"
	"math"
	"murakali/config"
	"murakali/internal/constant"
	"murakali/internal/model"
	"murakali/internal/module/admin"
	"murakali/internal/module/admin/delivery/body"
	"murakali/pkg/httperror"
	"murakali/pkg/pagination"
	"murakali/pkg/postgre"
	"murakali/pkg/response"
	"net/http"
	"time"
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

func (u *adminUC) GetRefunds(ctx context.Context, sortFilter string, pgn *pagination.Pagination) (*pagination.Pagination, error) {
	totalRows, err := u.adminRepo.GetTotalRefunds(ctx)
	if err != nil {
		return nil, err
	}

	totalPages := int(math.Ceil(float64(totalRows) / float64(pgn.Limit)))
	pgn.TotalRows = totalRows
	pgn.TotalPages = totalPages

	refunds, err := u.adminRepo.GetRefunds(ctx, sortFilter, pgn)
	if err != nil {
		return nil, err
	}

	pgn.Rows = refunds
	return pgn, nil
}

func (u *adminUC) RefundOrder(ctx context.Context, refundID string) error {
	refund, err := u.adminRepo.GetRefundByID(ctx, refundID)
	if err != nil {
		if err == sql.ErrNoRows {
			return httperror.New(http.StatusBadRequest, response.RefundNotFound)
		}

		return err
	}

	order, err := u.adminRepo.GetOrderByID(ctx, refund.OrderID.String())
	if err != nil {
		return err
	}

	errTx := u.txRepo.WithTransaction(func(tx postgre.Transaction) error {
		refund.RefundedAt.Valid = true
		refund.RefundedAt.Time = time.Now()

		if errRefund := u.adminRepo.UpdateRefund(ctx, tx, refund); errRefund != nil {
			return errRefund
		}

		order.OrderStatusID = constant.OrderStatusRefunded
		if errStatus := u.adminRepo.UpdateOrderStatus(ctx, tx, order); errStatus != nil {
			return errStatus
		}

		orderItems, err := u.adminRepo.GetOrderItemsByOrderID(ctx, tx, order.ID.String())
		if err != nil {
			return err
		}

		for _, item := range orderItems {
			productDetailData, err := u.adminRepo.GetProductDetailByID(ctx, tx, item.ProductDetailID.String())
			if err != nil {
				return err
			}

			productDetailData.Stock += float64(item.Quantity)
			errProduct := u.adminRepo.UpdateProductDetailStock(ctx, tx, productDetailData)
			if errProduct != nil {
				return errProduct
			}
		}

		walletMarketplace, err := u.adminRepo.GetWalletByUserID(ctx, constant.AdminMarketplaceID)
		if err != nil {
			return err
		}

		walletMarketplace.Balance -= order.TotalPrice
		walletMarketplace.UpdatedAt.Valid = true
		walletMarketplace.UpdatedAt.Time = time.Now()
		if err := u.adminRepo.UpdateWalletBalance(ctx, tx, walletMarketplace); err != nil {
			return err
		}

		walletUser, err := u.adminRepo.GetWalletByUserID(ctx, order.UserID.String())
		if err != nil {
			return err
		}

		walletUser.Balance += order.TotalPrice
		walletUser.UpdatedAt.Valid = true
		walletUser.UpdatedAt.Time = time.Now()
		if err := u.adminRepo.UpdateWalletBalance(ctx, tx, walletUser); err != nil {
			return err
		}

		walletMarketplaceHistory := &model.WalletHistory{
			TransactionID: order.TransactionID,
			WalletID:      walletMarketplace.ID,
			From:          walletMarketplace.ID.String(),
			To:            walletUser.ID.String(),
			Description:   "Refund order " + order.ID.String(),
			Amount:        order.TotalPrice,
			CreatedAt:     time.Now(),
		}
		if err := u.adminRepo.InsertWalletHistory(ctx, tx, walletMarketplaceHistory); err != nil {
			return err
		}

		walletUserHistory := &model.WalletHistory{
			TransactionID: order.TransactionID,
			WalletID:      walletUser.ID,
			From:          walletMarketplace.ID.String(),
			To:            walletUser.ID.String(),
			Description:   "Refund order " + order.ID.String(),
			Amount:        order.TotalPrice,
			CreatedAt:     time.Now(),
		}
		if err := u.adminRepo.InsertWalletHistory(ctx, tx, walletUserHistory); err != nil {
			return err
		}

		return nil
	})
	if errTx != nil {
		return errTx
	}

	return nil
}

func (u *adminUC) CreateVoucher(ctx context.Context, requestBody body.CreateVoucherRequest) error {
	count, _ := u.adminRepo.CountCodeVoucher(ctx, requestBody.Code)
	if count > 0 {
		return httperror.New(http.StatusBadRequest, body.CodeVoucherAlreadyExist)
	}

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
