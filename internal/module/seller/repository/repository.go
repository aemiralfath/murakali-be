package repository

import (
	"context"
	"database/sql"
	"fmt"
	"murakali/internal/model"
	"murakali/internal/module/seller"
	"murakali/internal/module/seller/delivery/body"
	"murakali/pkg/httperror"
	"murakali/pkg/pagination"
	"murakali/pkg/response"
	"net/http"

	"github.com/go-redis/redis/v8"
)

type sellerRepo struct {
	PSQL        *sql.DB
	RedisClient *redis.Client
}

func NewSellerRepository(psql *sql.DB, client *redis.Client) seller.Repository {
	return &sellerRepo{
		PSQL:        psql,
		RedisClient: client,
	}
}

func (r *sellerRepo) GetTotalOrder(ctx context.Context, shopID, orderStatusID string) (int64, error) {
	var total int64
	if err := r.PSQL.QueryRowContext(ctx, GetTotalOrderQuery, shopID, fmt.Sprintf("%%%s%%", orderStatusID)).Scan(&total); err != nil {
		return 0, err
	}

	return total, nil
}

func (r *sellerRepo) GetShopIDByUser(ctx context.Context, userID string) (string, error) {
	var shopID string
	if err := r.PSQL.QueryRowContext(ctx, GetShopIDByUserQuery, userID).Scan(&shopID); err != nil {
		return "", err
	}

	return shopID, nil
}

func (r *sellerRepo) GetShopIDByOrder(ctx context.Context, orderID string) (string, error) {
	var shopID string
	if err := r.PSQL.QueryRowContext(ctx, GetShopIDByOrderQuery, orderID).Scan(&shopID); err != nil {
		return "", err
	}

	return shopID, nil
}

func (r *sellerRepo) GetAddressByBuyerID(ctx context.Context, userID string) (*model.Address, error) {
	var address model.Address
	if err := r.PSQL.QueryRowContext(ctx, GetAddressByBuyerIDQuery, userID).Scan(
		&address.ID,
		&address.UserID,
		&address.Name,
		&address.ProvinceID,
		&address.CityID,
		&address.Province,
		&address.City,
		&address.District,
		&address.SubDistrict,
		&address.AddressDetail,
		&address.ZipCode,
		&address.IsDefault,
		&address.IsShopDefault,
		&address.CreatedAt,
		&address.UpdatedAt); err != nil {
		return nil, err
	}

	return &address, nil
}

func (r *sellerRepo) GetAddressBySellerID(ctx context.Context, userID string) (*model.Address, error) {
	var address model.Address
	if err := r.PSQL.QueryRowContext(ctx, GetAddressBySellerIDQuery, userID).Scan(
		&address.ID,
		&address.UserID,
		&address.Name,
		&address.ProvinceID,
		&address.CityID,
		&address.Province,
		&address.City,
		&address.District,
		&address.SubDistrict,
		&address.AddressDetail,
		&address.ZipCode,
		&address.IsDefault,
		&address.IsShopDefault,
		&address.CreatedAt,
		&address.UpdatedAt); err != nil {
		return nil, err
	}

	return &address, nil
}

func (r *sellerRepo) GetBuyerIDByOrderID(ctx context.Context, orderID string) (string, error) {
	var buyerID string
	if err := r.PSQL.QueryRowContext(ctx, GetBuyerIDByOrderIDQuery, orderID).Scan(
		&buyerID); err != nil {
		return "", err
	}

	return buyerID, nil
}

func (r *sellerRepo) GetSellerIDByOrderID(ctx context.Context, orderID string) (string, error) {
	var sellerID string
	if err := r.PSQL.QueryRowContext(ctx, GetSellerIDByOrderIDQuery, orderID).Scan(
		&sellerID); err != nil {
		return "", err
	}

	return sellerID, nil
}

func (r *sellerRepo) GetOrderByOrderID(ctx context.Context, orderID string) (*model.Order, error) {
	var order model.Order
	if err := r.PSQL.QueryRowContext(ctx, GetOrderByOrderID, orderID).Scan(
		&order.OrderID,
		&order.OrderStatus,
		&order.TotalPrice,
		&order.DeliveryFee,
		&order.ResiNumber,
		&order.ShopID,
		&order.ShopName,
		&order.SellerName,
		&order.VoucherCode,
		&order.CreatedAt,
		&order.Invoice,
		&order.CourierName,
		&order.CourierCode,
		&order.CourierService,
		&order.CourierDescription,
		&order.BuyerUsername,
		&order.BuyerPhoneNumber,
	); err != nil {
		return nil, err
	}

	orderDetail := make([]*model.OrderDetail, 0)

	res, err := r.PSQL.QueryContext(
		ctx, GetOrderDetailQuery, order.OrderID)

	if err != nil {
		return nil, err
	}
	for res.Next() {
		var detail model.OrderDetail
		if errScan := res.Scan(
			&detail.ProductDetailID,
			&detail.ProductID,
			&detail.ProductTitle,
			&detail.ProductDetailURL,
			&detail.OrderQuantity,
			&detail.ItemPrice,
			&detail.TotalPrice,
		); errScan != nil {
			return nil, errScan
		}
		orderDetail = append(orderDetail, &detail)
	}

	order.Detail = orderDetail
	return &order, nil
}

func (r *sellerRepo) GetOrders(ctx context.Context, shopID, orderStatusID string, pgn *pagination.Pagination) ([]*model.Order, error) {
	orders := make([]*model.Order, 0)

	res, err := r.PSQL.QueryContext(
		ctx, GetOrdersQuery,
		shopID,
		fmt.Sprintf("%%%s%%", orderStatusID),
		pgn.GetLimit(),
		pgn.GetOffset())

	if err != nil {
		return nil, err
	}
	defer res.Close()

	for res.Next() {
		var order model.Order
		if errScan := res.Scan(
			&order.OrderID,
			&order.OrderStatus,
			&order.TotalPrice,
			&order.DeliveryFee,
			&order.ResiNumber,
			&order.ShopID,
			&order.ShopName,
			&order.VoucherCode,
			&order.CreatedAt,
		); errScan != nil {
			return nil, err
		}

		orderDetail := make([]*model.OrderDetail, 0)

		res2, err2 := r.PSQL.QueryContext(
			ctx, GetOrderDetailQuery, order.OrderID)

		if err2 != nil {
			return nil, err2
		}

		for res2.Next() {
			var detail model.OrderDetail
			if errScan := res2.Scan(
				&detail.ProductDetailID,
				&detail.ProductID,
				&detail.ProductTitle,
				&detail.ProductDetailURL,
				&detail.OrderQuantity,
				&detail.ItemPrice,
				&detail.TotalPrice,
			); errScan != nil {
				return nil, err
			}
			orderDetail = append(orderDetail, &detail)
		}

		order.Detail = orderDetail
		orders = append(orders, &order)
	}
	if res.Err() != nil {
		return nil, err
	}
	return orders, nil
}

func (r *sellerRepo) ChangeOrderStatus(ctx context.Context, requestBody body.ChangeOrderStatusRequest) error {
	_, err := r.PSQL.ExecContext(
		ctx, ChangeOrderStatusQuery, requestBody.OrderStatusID, requestBody.OrderID)
	if err != nil {
		return err
	}

	return nil
}

func (r *sellerRepo) GetCourierSeller(ctx context.Context, userID string) ([]*body.CourierSellerRelationInfo, error) {
	courierSeller := make([]*body.CourierSellerRelationInfo, 0)

	res, err := r.PSQL.QueryContext(
		ctx, GetCourierSellerQuery,
		userID,
	)

	if err != nil {
		return nil, err
	}
	defer res.Close()

	for res.Next() {
		var courierSellerData body.CourierSellerRelationInfo
		if errScan := res.Scan(
			&courierSellerData.ShopCourierID,
			&courierSellerData.CourierID,
			&courierSellerData.DeletedAt,
		); errScan != nil {
			return nil, err
		}
		courierSeller = append(courierSeller, &courierSellerData)
	}
	if res.Err() != nil {
		return nil, err
	}
	return courierSeller, err
}

func (r *sellerRepo) GetAllCourier(ctx context.Context) ([]*body.CourierInfo, error) {
	courier := make([]*body.CourierInfo, 0)

	res, err := r.PSQL.QueryContext(
		ctx, GetAllCourierQuery,
	)

	if err != nil {
		return nil, err
	}
	defer res.Close()

	for res.Next() {
		var courierData body.CourierInfo
		if errScan := res.Scan(
			&courierData.CourierID,
			&courierData.Name,
			&courierData.Code,
			&courierData.Service,
			&courierData.Description,
		); errScan != nil {
			return nil, err
		}

		courier = append(courier, &courierData)
	}

	if res.Err() != nil {
		return nil, err
	}

	return courier, err
}

func (r *sellerRepo) GetSellerBySellerID(ctx context.Context, sellerID string) (*body.SellerResponse, error) {
	var sellerData body.SellerResponse
	if err := r.PSQL.QueryRowContext(ctx, GetShopIDByShopIDQuery, sellerID).Scan(
		&sellerData.ID,
		&sellerData.UserID,
		&sellerData.Name,
		&sellerData.TotalProduct,
		&sellerData.TotalRating,
		&sellerData.RatingAVG,
		&sellerData.CreatedAt,
		&sellerData.PhotoURL,
	); err != nil {
		return nil, err
	}

	return &sellerData, nil
}

func (r *sellerRepo) GetCategoryBySellerID(ctx context.Context, shopID string) ([]*body.CategoryResponse, error) {
	categories := make([]*body.CategoryResponse, 0)

	res, err := r.PSQL.QueryContext(
		ctx, GetCategoryBySellerIDQuery, shopID)

	if err != nil {
		return nil, err
	}
	defer res.Close()

	for res.Next() {
		var category body.CategoryResponse
		if errScan := res.Scan(
			&category.ID,
			&category.Name,
		); errScan != nil {
			return nil, err
		}

		categories = append(categories, &category)
	}
	if res.Err() != nil {
		return nil, err
	}
	return categories, nil
}

func (r *sellerRepo) GetShopIDByUserID(ctx context.Context, userID string) (string, error) {
	var ID string
	if err := r.PSQL.QueryRowContext(ctx, GetShopIDByUserIDQuery, userID).Scan(&ID); err != nil {
		return "", err
	}
	return ID, nil
}

func (r *sellerRepo) GetCourierSellerNotNullByShopAndCourierID(ctx context.Context, shopID, courierID string) (string, error) {
	var ID string
	if err := r.PSQL.QueryRowContext(ctx, GetCourierSellerNotNullByShopAndCourierIDQuery, shopID, courierID).Scan(&ID); err != nil {
		return "", err
	}
	return ID, nil
}

func (r *sellerRepo) GetCourierByID(ctx context.Context, courierID string) (string, error) {
	var ID string
	if err := r.PSQL.QueryRowContext(ctx, GetCourierByIDQuery, courierID).Scan(&ID); err != nil {
		return "", err
	}
	return ID, nil
}

func (r *sellerRepo) CreateCourierSeller(ctx context.Context, shopID, courierID string) error {
	if _, err := r.PSQL.ExecContext(ctx, CreateCourierSellerQuery,
		shopID,
		courierID); err != nil {
		return err
	}
	return nil
}

func (r *sellerRepo) GetCourierSellerByID(ctx context.Context, shopCourierID string) (string, error) {
	var ID string
	if err := r.PSQL.QueryRowContext(ctx, GetCourierSellerByIDQuery, shopCourierID).Scan(&ID); err != nil {
		return "", err
	}
	return ID, nil
}

func (r *sellerRepo) UpdateCourierSellerByID(ctx context.Context, shopID, courierID string) error {
	_, err := r.PSQL.ExecContext(ctx, UpdateCourierSellerQuery, courierID, shopID)
	if err != nil {
		return err
	}
	return nil
}

func (r *sellerRepo) DeleteCourierSellerByID(ctx context.Context, shopCourierID string) error {
	_, err := r.PSQL.ExecContext(ctx, DeleteCourierSellerQuery, shopCourierID)
	if err != nil {
		return err
	}
	return nil
}

func (r *sellerRepo) UpdateResiNumberInOrderSeller(ctx context.Context, noResi, orderID, shopID string) error {
	temp, err := r.PSQL.ExecContext(ctx,
		UpdateResiNumberInOrderSellerQuery,
		noResi, orderID, shopID)
	if err != nil {
		return err
	}

	rowsAffected, _ := temp.RowsAffected()
	if rowsAffected == 0 {
		return httperror.New(http.StatusNotFound, response.OrderNotExistMessage)
	}
	return nil
}
