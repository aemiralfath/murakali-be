package usecase

import (
	"context"
	"math"
	"murakali/config"
	"murakali/internal/model"
	"murakali/internal/module/seller"
	"murakali/internal/module/seller/delivery/body"
	"murakali/pkg/httperror"
	"murakali/pkg/pagination"
	"murakali/pkg/postgre"
	"murakali/pkg/response"
	"net/http"
)

type sellerUC struct {
	cfg        *config.Config
	txRepo     *postgre.TxRepo
	sellerRepo seller.Repository
}

func NewSellerUseCase(cfg *config.Config, txRepo *postgre.TxRepo, sellerRepo seller.Repository) seller.UseCase {
	return &sellerUC{cfg: cfg, txRepo: txRepo, sellerRepo: sellerRepo}
}

func (u *sellerUC) GetOrder(ctx context.Context, userID, orderStatusID string, pgn *pagination.Pagination) (*pagination.Pagination, error) {
	shopID, err := u.sellerRepo.GetShopIDByUser(ctx, userID)
	if err != nil {
		return nil, err
	}

	totalRows, err := u.sellerRepo.GetTotalOrder(ctx, shopID, orderStatusID)
	if err != nil {
		return nil, err
	}

	totalPages := int(math.Ceil(float64(totalRows) / float64(pgn.Limit)))
	pgn.TotalRows = totalRows
	pgn.TotalPages = totalPages

	orders, err := u.sellerRepo.GetOrders(ctx, shopID, orderStatusID, pgn)
	if err != nil {
		return nil, err
	}

	pgn.Rows = orders
	return pgn, nil
}

func (u *sellerUC) ChangeOrderStatus(ctx context.Context, userID string, requestBody body.ChangeOrderStatusRequest) error {
	shopIDFromUser, err := u.sellerRepo.GetShopIDByUser(ctx, userID)
	if err != nil {
		return err
	}
	shopIDFromOrder, err := u.sellerRepo.GetShopIDByOrder(ctx, requestBody.OrderID)
	if err != nil {
		return err
	}

	if shopIDFromUser != shopIDFromOrder {
		return httperror.New(http.StatusUnauthorized, response.UnauthorizedMessage)
	}

	err = u.sellerRepo.ChangeOrderStatus(ctx, requestBody)
	if err != nil {
		return err
	}
	return nil
}

func (u *sellerUC) GetOrderByOrderID(ctx context.Context, orderID string) (*model.Order, error) {
	order, err := u.sellerRepo.GetOrderByOrderID(ctx, orderID)
	if err != nil {
		return nil, err
	}

	return order, nil
}
