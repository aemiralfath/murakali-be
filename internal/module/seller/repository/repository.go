package repository

import (
	"context"
	"database/sql"
	"fmt"
	"murakali/internal/constant"
	"murakali/internal/model"
	"murakali/internal/module/seller"
	"murakali/internal/module/seller/delivery/body"
	"murakali/pkg/httperror"
	"murakali/pkg/pagination"
	"murakali/pkg/postgre"
	"murakali/pkg/response"
	"net/http"
	"time"

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
		&order.ShopPhoneNumber,
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
			&detail.ProductWeight,
			&detail.ProductDetailURL,
			&detail.OrderQuantity,
			&detail.ItemPrice,
			&detail.TotalPrice,
		); errScan != nil {
			return nil, errScan
		}
		variant := make(map[string]string, 0)
		variantResult, errVariant := r.PSQL.QueryContext(ctx, GetOrderDetailProductVariant, detail.ProductDetailID)
		if errVariant != nil {
			if errVariant != sql.ErrNoRows {
				return nil, err
			}
		}
		for variantResult.Next() {
			var varName string
			var varType string
			if errScanVariant := variantResult.Scan(
				&varName,
				&varType,
			); errScanVariant != nil {
				return nil, errScanVariant
			}
			variant[varName] = varType
		}

		detail.Variant = variant
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
				&detail.ProductWeight,
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

func (r *sellerRepo) GetOrdersOnDelivery(ctx context.Context) ([]*model.OrderModel, error) {
	orders := make([]*model.OrderModel, 0)
	res, err := r.PSQL.QueryContext(ctx, GetOrderOnDeliveryQuery, constant.OrderStatusOnDelivery)
	if err != nil {
		return orders, err
	}
	defer res.Close()

	for res.Next() {
		var order model.OrderModel
		if errScan := res.Scan(
			&order.ID,
			&order.OrderStatusID,
			&order.ArrivedAt,
		); errScan != nil {
			return orders, err
		}
		orders = append(orders, &order)
	}
	if res.Err() != nil {
		return orders, err
	}

	return orders, err
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

func (r *sellerRepo) GetSellerByUserID(ctx context.Context, sellerID string) (*body.SellerResponse, error) {
	var sellerData body.SellerResponse
	if err := r.PSQL.QueryRowContext(ctx, GetShopDetailIDByUserIDQuery, sellerID).Scan(
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

func (r *sellerRepo) UpdateResiNumberInOrderSeller(ctx context.Context, noResi, orderID, shopID string, arriveAt time.Time) error {
	temp, err := r.PSQL.ExecContext(ctx,
		UpdateResiNumberInOrderSellerQuery,
		noResi, arriveAt, constant.OrderStatusOnDelivery, orderID, shopID)
	if err != nil {
		return err
	}

	rowsAffected, _ := temp.RowsAffected()
	if rowsAffected == 0 {
		return httperror.New(http.StatusNotFound, response.OrderNotExistMessage)
	}
	return nil
}

func (r *sellerRepo) GetCostRedis(ctx context.Context, key string) (*string, error) {
	res := r.RedisClient.Get(ctx, key)
	if res.Err() != nil {
		return nil, res.Err()
	}

	value, err := res.Result()
	if err != nil {
		return nil, err
	}

	return &value, nil
}

func (r *sellerRepo) InsertCostRedis(ctx context.Context, key, value string) error {
	if err := r.RedisClient.Set(ctx, key, value, 0); err.Err() != nil {
		return err.Err()
	}

	return nil
}

func (r *sellerRepo) GetTotalVoucherSeller(ctx context.Context, shopID string) (int64, error) {
	var total int64
	if err := r.PSQL.QueryRowContext(ctx, GetTotalVoucherSellerQuery, shopID).Scan(&total); err != nil {
		return -1, err
	}

	return total, nil
}

func (r *sellerRepo) GetAllVoucherSeller(ctx context.Context, shopID string) ([]*model.Voucher, error) {
	var shopVouchers []*model.Voucher
	res, err := r.PSQL.QueryContext(ctx, GetAllVoucherSellerQuery, shopID)
	if err != nil {
		return nil, err
	}
	defer res.Close()

	for res.Next() {
		var voucher model.Voucher
		if errScan := res.Scan(
			&voucher.ID,
			&voucher.ShopID,
			&voucher.Code,
			&voucher.Quota,
			&voucher.ActivedDate,
			&voucher.ExpiredDate,
			&voucher.DiscountPercentage,
			&voucher.DiscountFixPrice,
			&voucher.MinProductPrice,
			&voucher.MaxDiscountPrice,
			&voucher.CreatedAt,
			&voucher.UpdatedAt,
			&voucher.DeletedAt,
		); errScan != nil {
			return nil, err
		}

		shopVouchers = append(shopVouchers, &voucher)
	}

	if res.Err() != nil {
		return nil, err
	}

	return shopVouchers, nil
}

func (r *sellerRepo) CreateVoucherSeller(ctx context.Context, voucherShop *model.Voucher) error {
	if _, err := r.PSQL.ExecContext(ctx, CreateVoucherSellerQuery,
		voucherShop.ShopID,
		voucherShop.Code,
		voucherShop.Quota,
		voucherShop.ActivedDate,
		voucherShop.ExpiredDate,
		voucherShop.DiscountPercentage,
		voucherShop.DiscountFixPrice,
		voucherShop.MinProductPrice,
		voucherShop.MaxDiscountPrice); err != nil {
		return err
	}
	return nil
}

func (r *sellerRepo) UpdateVoucherSeller(ctx context.Context, voucherShop *model.Voucher) error {
	if _, err := r.PSQL.ExecContext(ctx, UpdateVoucherSellerQuery,
		voucherShop.Quota,
		voucherShop.ActivedDate,
		voucherShop.ExpiredDate,
		voucherShop.DiscountPercentage,
		voucherShop.DiscountFixPrice,
		voucherShop.MinProductPrice,
		voucherShop.MaxDiscountPrice,
		voucherShop.ID); err != nil {
		return err
	}
	return nil
}

func (r *sellerRepo) DeleteVoucherSeller(ctx context.Context, voucherIDShopID *body.VoucherIDShopID) error {
	_, err := r.PSQL.ExecContext(ctx, DeleteVoucherSellerQuery, voucherIDShopID.VoucherID, voucherIDShopID.ShopID)
	if err != nil {
		return err
	}
	return nil
}

func (r *sellerRepo) GetAllVoucherSellerByIDAndShopID(ctx context.Context, voucherIDShopID *body.VoucherIDShopID) (*model.Voucher, error) {
	var voucher model.Voucher
	if err := r.PSQL.QueryRowContext(ctx, GetAllVoucherSellerByIDandShopIDQuery, voucherIDShopID.VoucherID, voucherIDShopID.ShopID).Scan(
		&voucher.ID,
		&voucher.ShopID,
		&voucher.Code,
		&voucher.Quota,
		&voucher.ActivedDate,
		&voucher.ExpiredDate,
		&voucher.DiscountPercentage,
		&voucher.DiscountFixPrice,
		&voucher.MinProductPrice,
		&voucher.MaxDiscountPrice,
		&voucher.CreatedAt,
		&voucher.UpdatedAt,
		&voucher.DeletedAt,
	); err != nil {
		return nil, err
	}

	return &voucher, nil
}

func (r *sellerRepo) UpdateOrder(ctx context.Context, tx postgre.Transaction, orderData *model.OrderModel) error {
	_, err := tx.ExecContext(ctx, UpdateOrderByID, orderData.OrderStatusID, orderData.ID)
	if err != nil {
		return err
	}

	return nil
}

func (r *sellerRepo) UpdateTransaction(ctx context.Context, tx postgre.Transaction, transactionData *model.Transaction) error {
	_, err := tx.ExecContext(ctx, UpdateTransactionByID, transactionData.PaidAt, transactionData.CanceledAt, transactionData.ID)
	if err != nil {
		return err
	}

	return nil
}

func (r *sellerRepo) GetOrderByTransactionID(ctx context.Context, tx postgre.Transaction, transactionID string) ([]*model.OrderModel, error) {
	orders := make([]*model.OrderModel, 0)
	res, err := tx.QueryContext(ctx, GetOrderByTransactionID, transactionID)
	if err != nil {
		return nil, err
	}
	defer res.Close()

	for res.Next() {
		var order model.OrderModel
		if errScan := res.Scan(
			&order.ID,
			&order.TransactionID,
			&order.UserID,
			&order.ShopID,
			&order.CourierID,
			&order.VoucherShopID,
			&order.OrderStatusID,
			&order.TotalPrice,
			&order.DeliveryFee,
			&order.ResiNo,
			&order.CreatedAt,
			&order.ArrivedAt); errScan != nil {
			return nil, errScan
		}

		orders = append(orders, &order)
	}

	if res.Err() != nil {
		return nil, res.Err()
	}

	return orders, nil
}

func (r *sellerRepo) GetTransactionsExpired(ctx context.Context) ([]*model.Transaction, error) {
	transactions := make([]*model.Transaction, 0)
	res, err := r.PSQL.QueryContext(ctx, GetTransactionsExpiredQuery)
	if err != nil {
		return nil, err
	}
	defer res.Close()

	for res.Next() {
		var transaction model.Transaction
		if errScan := res.Scan(
			&transaction.ID,
			&transaction.VoucherMarketplaceID,
			&transaction.WalletID,
			&transaction.CardNumber,
			&transaction.Invoice,
			&transaction.TotalPrice,
			&transaction.PaidAt,
			&transaction.CanceledAt,
			&transaction.ExpiredAt); errScan != nil {
			return nil, errScan
		}

		transactions = append(transactions, &transaction)
	}

	if res.Err() != nil {
		return nil, res.Err()
	}

	return transactions, nil
}

func (r *sellerRepo) GetOrderItemsByOrderID(ctx context.Context, tx postgre.Transaction, orderID string) ([]*model.OrderItem, error) {
	orderItems := make([]*model.OrderItem, 0)
	res, err := tx.QueryContext(ctx, GetOrderItemsByOrderIDQuery, orderID)
	if err != nil {
		return nil, err
	}
	defer res.Close()

	for res.Next() {
		var orderItem model.OrderItem
		if errScan := res.Scan(
			&orderItem.ID,
			&orderItem.OrderID,
			&orderItem.ProductDetailID,
			&orderItem.Quantity,
			&orderItem.ItemPrice,
			&orderItem.TotalPrice); errScan != nil {
			return nil, err
		}

		orderItems = append(orderItems, &orderItem)
	}

	if res.Err() != nil {
		return nil, err
	}

	return orderItems, nil
}

func (r *sellerRepo) GetProductDetailByID(ctx context.Context, tx postgre.Transaction, productDetailID string) (*model.ProductDetail, error) {
	var pd model.ProductDetail
	if err := tx.QueryRowContext(ctx, GetProductDetailByIDQuery, productDetailID).Scan(
		&pd.ID,
		&pd.Price,
		&pd.Stock,
		&pd.Size,
		&pd.Weight,
		&pd.Hazardous,
		&pd.Condition,
		&pd.BulkPrice); err != nil {
		return nil, err
	}

	return &pd, nil
}

func (r *sellerRepo) UpdateProductDetailStock(ctx context.Context, tx postgre.Transaction,
	productDetailData *model.ProductDetail) error {
	_, err := tx.ExecContext(ctx, UpdateProductDetailStockQuery, productDetailData.Stock, productDetailData.ID)
	if err != nil {
		return err
	}

	return nil
}
