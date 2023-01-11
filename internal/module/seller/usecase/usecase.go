package usecase

import (
	"context"
	"database/sql"
	"fmt"
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
	courier, err := u.sellerRepo.GetAllCourier(ctx)
	if err != nil {
		return nil, err
	}

	courierSeller, err := u.sellerRepo.GetCourierSeller(ctx, userID)
	if err != nil {
		return nil, err
	}

	resultCourierSeller := make([]*body.CourierSellerInfo, 0)
	totalData := len(courier)
	totalDataCourierSeller := len(courierSeller)
	fmt.Println("ini courier seller", courierSeller)
	for i := 0; i < totalData; i++ {
		var shopCourierIDTemp string
		var deletedAtTemp string
		for j := 0; j < totalDataCourierSeller; j++ {
			if courier[i].CourierID == courierSeller[j].CourierID {
				shopCourierIDTemp = courierSeller[j].ShopCourierID.String()
				if !courierSeller[j].DeletedAt.Time.IsZero() {
					deletedAtTemp = courierSeller[j].DeletedAt.Time.String()
				}
			}
		}
		p := &body.CourierSellerInfo{
			ShopCourierID: shopCourierIDTemp,
			CourierID:     courier[i].CourierID,
			Name:          courier[i].Name,
			Code:          courier[i].Code,
			Service:       courier[i].Service,
			Description:   courier[i].Description,
			DeletedAt:     deletedAtTemp,
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

func (u *sellerUC) CreateCourierSeller(ctx context.Context, userID, courierID string) error {
	_, err := u.sellerRepo.GetCourierByID(ctx, courierID)
	if err != nil {
		if err == sql.ErrNoRows {
			return httperror.New(http.StatusBadRequest, body.CourierNotFoundMessage)
		}
		return err
	}

	shopID, err := u.sellerRepo.GetShopIDByUserID(ctx, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return httperror.New(http.StatusBadRequest, response.ShopAddressNotFound)
		}
		return err
	}

	sellerCourierID, _ := u.sellerRepo.GetCourierSellerNotNullByShopAndCourierID(ctx, shopID, courierID)

	if sellerCourierID != "" {
		if err != nil {
			if err == sql.ErrNoRows {
				return httperror.New(http.StatusBadRequest, body.CourierSellerNotFoundMessage)
			}
			return err
		}

		err = u.sellerRepo.UpdateCourierSellerByID(ctx, shopID, courierID)
		if err == sql.ErrNoRows {
			return httperror.New(http.StatusBadRequest, body.CourierSellerAlreadyExistMessage)
		}
		return nil
	}

	err = u.sellerRepo.CreateCourierSeller(ctx, shopID, courierID)
	if err != nil {
		return err
	}
	return nil
}

func (u *sellerUC) DeleteCourierSellerByID(ctx context.Context, shopCourierID string) error {
	_, err := u.sellerRepo.GetCourierSellerByID(ctx, shopCourierID)
	if err != nil {
		if err == sql.ErrNoRows {
			return httperror.New(http.StatusBadRequest, body.CourierSellerNotFoundMessage)
		}
		return err
	}

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

func (u *sellerUC) UpdateResiNumberInOrderSeller(ctx context.Context, userID, orderID string, requestBody body.UpdateNoResiOrderSellerRequest) error {
	shopID, err := u.sellerRepo.GetShopIDByUserID(ctx, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return httperror.New(http.StatusBadRequest, response.ShopAddressNotFound)
		}
		return err
	}

	err = u.sellerRepo.UpdateResiNumberInOrderSeller(ctx, requestBody.NoResi, orderID, shopID)
	if err != nil {
		return err
	}

	return nil
}
