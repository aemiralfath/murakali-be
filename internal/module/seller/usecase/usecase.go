package usecase

import (
	"context"
	"database/sql"
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

func (u *sellerUC) GetCourierSeller(ctx context.Context, userID string) (*body.CourierSellerResponse, error) {
	courierSeller, err := u.sellerRepo.GetCourierSeller(ctx, userID)
	if err != nil {
		return nil, err
	}

	resultCourierSeller := make([]*body.CourierSellerInfo, 0)
	totalData := len(courierSeller)
	for i := 0; i < totalData; i++ {
		p := &body.CourierSellerInfo{
			ShopCourierID: courierSeller[i].ShopCourierID,
			Name:          courierSeller[i].Name,
			Code:          courierSeller[i].Code,
			Service:       courierSeller[i].Service,
			Description:   courierSeller[i].Description,
		}

		resultCourierSeller = append(resultCourierSeller, p)
	}

	csr := &body.CourierSellerResponse{}

	csr.Rows = resultCourierSeller

	return csr, nil
}
func (u *sellerUC) GetSellerBySellerID(ctx context.Context, sellerID string) (*body.SellerResponse, error) {
	sellerData, err := u.sellerRepo.GetSellerBySellerID(ctx, sellerID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, httperror.New(http.StatusNotFound, body.SellerNotFoundMessage)
		}
		return nil, err
	}

	return sellerData, nil
}

func (u *sellerUC) CreateCourierSeller(ctx context.Context, userID string, courierId string) error {
	shopID, err := u.sellerRepo.GetShopIDByUserID(ctx, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return httperror.New(http.StatusBadRequest, response.UserNotExistMessage)
		}
		return err
	}

	sellerCourierID, _ := u.sellerRepo.GetCourierSellerByID(ctx, shopID, courierId)
	if sellerCourierID != "" {
		return httperror.New(http.StatusBadRequest, body.CourierSellerAlreadyExistMessage)
	}

	err = u.sellerRepo.CreateCourierSeller(ctx, shopID, courierId)
	if err != nil {
		return err
	}
	return nil
}

func (u *sellerUC) DeleteCourierSellerByID(ctx context.Context, shopCourierID string) error {
	if err := u.sellerRepo.DeleteCourierSellerByID(ctx, shopCourierID); err != nil {
		if err == sql.ErrNoRows {
			return httperror.New(http.StatusNotFound, body.CourierSellerNotFoundMessage)
		}
		return err
	}
	return nil
}
func (u *sellerUC) GetCategoryBySellerID(ctx context.Context, shopID string) ([]*body.CategoryResponse, error) {
	categories, err := u.sellerRepo.GetCategoryBySellerID(ctx, shopID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, httperror.New(http.StatusNotFound, body.CategoryNotFoundMessage)
		}
		return nil, err
	}

	return categories, nil
}
