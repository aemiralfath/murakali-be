package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"murakali/internal/constant"
	"murakali/internal/model"
	"murakali/internal/module/user"
	"murakali/internal/module/user/delivery/body"
	"murakali/pkg/pagination"
	"time"

	"murakali/pkg/postgre"

	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
)

type userRepo struct {
	PSQL        *sql.DB
	RedisClient *redis.Client
}

func NewUserRepository(psql *sql.DB, client *redis.Client) user.Repository {
	return &userRepo{
		PSQL:        psql,
		RedisClient: client,
	}
}

func (r *userRepo) CreateAddress(ctx context.Context, tx postgre.Transaction, userID string, requestBody body.CreateAddressRequest) error {
	_, err := tx.ExecContext(
		ctx,
		CreateAddressQuery,
		userID,
		requestBody.Name,
		requestBody.ProvinceID,
		requestBody.CityID,
		requestBody.Province,
		requestBody.City,
		requestBody.District,
		requestBody.SubDistrict,
		requestBody.AddressDetail,
		requestBody.ZipCode,
		requestBody.IsDefault,
		requestBody.IsShopDefault)
	if err != nil {
		return err
	}

	return nil
}

func (r *userRepo) GetAddressByID(ctx context.Context, userID, addressID string) (*model.Address, error) {
	var address model.Address
	if err := r.PSQL.QueryRowContext(ctx, GetAddressByIDQuery, addressID, userID).Scan(
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

func (r *userRepo) GetDefaultUserAddress(ctx context.Context, userID string) (*model.Address, error) {
	var address model.Address
	if err := r.PSQL.QueryRowContext(ctx, GetDefaultAddressQuery, userID, true).Scan(
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

func (r *userRepo) GetDefaultShopAddress(ctx context.Context, userID string) (*model.Address, error) {
	var address model.Address
	if err := r.PSQL.QueryRowContext(ctx, GetDefaultShopAddressQuery, userID, true).
		Scan(&address.ID,
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

func (r *userRepo) GetTotalAddress(ctx context.Context, userID, name string) (int64, error) {
	var total int64
	if err := r.PSQL.QueryRowContext(ctx, GetTotalAddressQuery, userID, fmt.Sprintf("%%%s%%", name)).Scan(&total); err != nil {
		return 0, err
	}

	return total, nil
}

func (r *userRepo) GetTotalAddressDefault(ctx context.Context, userID, name string, isDefault, isShopDefault bool) (int64, error) {
	var total int64
	if err := r.PSQL.QueryRowContext(
		ctx,
		GetTotalAddressDefaultQuery,
		userID,
		fmt.Sprintf("%%%s%%", name),
		isDefault,
		isShopDefault).Scan(&total); err != nil {
		return 0, err
	}

	return total, nil
}

func (r *userRepo) GetTotalOrder(ctx context.Context, userID, orderStatusID string) (int64, error) {
	var total int64
	if err := r.PSQL.QueryRowContext(ctx, GetTotalOrderQuery, userID, fmt.Sprintf("%%%s%%", orderStatusID)).Scan(&total); err != nil {
		return 0, err
	}

	return total, nil
}

func (r *userRepo) GetTotalTransactionByUserID(ctx context.Context, userID string) (int64, error) {
	var total int64
	if err := r.PSQL.QueryRowContext(ctx, GetTotalTransactionByUserIDQuery, userID).Scan(&total); err != nil {
		if err == sql.ErrNoRows {
			return 0, nil
		}
		return 0, err
	}

	return total, nil
}

func (r *userRepo) UpdateAddress(ctx context.Context, tx postgre.Transaction, address *model.Address) error {
	_, err := tx.ExecContext(ctx, UpdateAddressByIDQuery,
		address.Name,
		address.ProvinceID,
		address.CityID,
		address.Province,
		address.City,
		address.District,
		address.SubDistrict,
		address.AddressDetail,
		address.ZipCode,
		address.IsDefault,
		address.IsShopDefault,
		address.UpdatedAt,
		address.ID,
	)
	if err != nil {
		return err
	}

	return nil
}

func (r *userRepo) UpdateDefaultAddress(ctx context.Context, tx postgre.Transaction, status bool, address *model.Address) error {
	_, err := tx.ExecContext(ctx, UpdateDefaultAddressQuery, status, address.ID)
	if err != nil {
		return err
	}

	return nil
}

func (r *userRepo) UpdateDefaultShopAddress(ctx context.Context, tx postgre.Transaction, status bool, address *model.Address) error {
	_, err := tx.ExecContext(ctx, UpdateDefaultShopAddressQuery, status, address.ID)
	if err != nil {
		return err
	}

	return nil
}

func (r *userRepo) DeleteAddress(ctx context.Context, addressID string) error {
	_, err := r.PSQL.ExecContext(ctx, DeleteAddressByIDQuery, addressID)
	if err != nil {
		return err
	}

	return nil
}

func (r *userRepo) GetAddresses(ctx context.Context, userID, name string, isDefault, isShopDefault bool,
	pgn *pagination.Pagination) ([]*model.Address, error) {
	addresses := make([]*model.Address, 0)
	res, err := r.PSQL.QueryContext(
		ctx, GetAddressesDefaultQuery,
		userID,
		fmt.Sprintf("%%%s%%", name),
		isDefault,
		isShopDefault,
		pgn.GetSort(),
		pgn.GetLimit(),
		pgn.GetOffset())

	if err != nil {
		return nil, err
	}
	defer res.Close()

	for res.Next() {
		var address model.Address
		if errScan := res.Scan(
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
			&address.UpdatedAt); errScan != nil {
			return nil, errScan
		}

		addresses = append(addresses, &address)
	}

	if res.Err() != nil {
		return nil, res.Err()
	}

	return addresses, nil
}

func (r *userRepo) GetAllAddresses(ctx context.Context, userID, name string, pgn *pagination.Pagination) ([]*model.Address, error) {
	addresses := make([]*model.Address, 0)
	res, err := r.PSQL.QueryContext(
		ctx, GetAllAddressesQuery,
		userID,
		fmt.Sprintf("%%%s%%", name),
		pgn.GetSort(),
		pgn.GetLimit(),
		pgn.GetOffset())

	if err != nil {
		return nil, err
	}
	defer res.Close()

	for res.Next() {
		var address model.Address
		if errScan := res.Scan(
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
			&address.UpdatedAt); errScan != nil {
			return nil, err
		}

		addresses = append(addresses, &address)
	}

	if res.Err() != nil {
		return nil, err
	}

	return addresses, nil
}

func (r *userRepo) GetOrderByOrderID(ctx context.Context, orderID string) (*model.Order, error) {
	var order model.Order
	var strShopAddress, strBuyerAddress string
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
		&strShopAddress,
		&strBuyerAddress,
	); err != nil {
		return nil, err
	}

	json.Unmarshal([]byte(strBuyerAddress), &order.BuyerAddress)
	json.Unmarshal([]byte(strShopAddress), &order.SellerAddress)
	orderDetail := make([]*model.OrderDetail, 0)

	res, err := r.PSQL.QueryContext(
		ctx, GetOrderDetailQuery2, order.OrderID)

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

func (r *userRepo) GetBuyerIDByOrderID(ctx context.Context, orderID string) (string, error) {
	var buyerID string
	if err := r.PSQL.QueryRowContext(ctx, GetBuyerIDByOrderIDQuery, orderID).Scan(
		&buyerID); err != nil {
		return "", err
	}

	return buyerID, nil
}

func (r *userRepo) GetSellerIDByOrderID(ctx context.Context, orderID string) (string, error) {
	var sellerID string
	if err := r.PSQL.QueryRowContext(ctx, GetSellerIDByOrderIDQuery, orderID).Scan(
		&sellerID); err != nil {
		return "", err
	}

	return sellerID, nil
}

func (r *userRepo) GetAddressByBuyerID(ctx context.Context, userID string) (*model.Address, error) {
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

func (r *userRepo) GetAddressBySellerID(ctx context.Context, userID string) (*model.Address, error) {
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

func (r *userRepo) GetCostRedis(ctx context.Context, key string) (*string, error) {
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

func (r *userRepo) InsertCostRedis(ctx context.Context, key, value string) error {
	if err := r.RedisClient.Set(ctx, key, value, 0); err.Err() != nil {
		return err.Err()
	}

	return nil
}

func (r *userRepo) GetOrders(ctx context.Context, userID, orderStatusID string, pgn *pagination.Pagination) ([]*model.Order, error) {
	orders := make([]*model.Order, 0)

	res, err := r.PSQL.QueryContext(
		ctx, GetOrdersQuery,
		userID,
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

func (r *userRepo) GetOrderDetailByTransactionID(ctx context.Context, transactionID string) ([]*model.Order, error) {
	orders := make([]*model.Order, 0)

	res, err := r.PSQL.QueryContext(ctx, GetOrdersByTransactionIDQuery, transactionID)

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
			return nil, errScan
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

		if res2.Err() != nil {
			return nil, res2.Err()
		}

		order.Detail = orderDetail
		orders = append(orders, &order)
	}
	if res.Err() != nil {
		return nil, err
	}
	return orders, nil
}

func (r *userRepo) GetOrdersByTransactionID(ctx context.Context, transactionID, userID string) ([]*model.Order, error) {
	orders := make([]*model.Order, 0)
	res, err := r.PSQL.QueryContext(
		ctx, GetOrdersByTransactionIDQuery,
		transactionID)

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
			return nil, errScan
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

		if res2.Err() != nil {
			return nil, res2.Err()
		}

		order.Detail = orderDetail

		orders = append(orders, &order)
	}

	if res.Err() != nil {
		return nil, res.Err()
	}
	return orders, nil
}

func (r *userRepo) GetTransactionByUserID(ctx context.Context, userID string, status int, pgn *pagination.Pagination) ([]*model.Transaction, error) {
	transactions := make([]*model.Transaction, 0)

	var query string
	switch status {
	case constant.OrderStatusWaitingToPay:
		query = GetTransactionByUserIDNotPaidQuery
	default:
		query = GetTransactionByUserIDQuery
	}

	res, err := r.PSQL.QueryContext(
		ctx, query,
		userID,
		pgn.GetLimit(),
		pgn.GetOffset())

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
			return nil, err
		}

		transactions = append(transactions, &transaction)
	}

	if res.Err() != nil {
		return nil, res.Err()
	}

	return transactions, nil
}

func (r *userRepo) GetTransactionByID(ctx context.Context, transactionID string) (*model.Transaction, error) {
	var transaction model.Transaction
	if err := r.PSQL.QueryRowContext(ctx, GetTransactionByIDQuery, transactionID).
		Scan(
			&transaction.ID,
			&transaction.VoucherMarketplaceID,
			&transaction.WalletID,
			&transaction.CardNumber,
			&transaction.Invoice,
			&transaction.TotalPrice,
			&transaction.PaidAt,
			&transaction.CanceledAt,
			&transaction.ExpiredAt); err != nil {
		return nil, err
	}

	return &transaction, nil
}

func (r *userRepo) GetUserByID(ctx context.Context, id string) (*model.User, error) {
	var userModel model.User
	if err := r.PSQL.QueryRowContext(ctx, GetUserByIDQuery, id).
		Scan(&userModel.ID, &userModel.RoleID, &userModel.Email, &userModel.Username, &userModel.PhoneNo,
			&userModel.FullName, &userModel.Gender, &userModel.BirthDate, &userModel.IsVerify, &userModel.PhotoURL); err != nil {
		return nil, err
	}

	return &userModel, nil
}

func (r *userRepo) GetUserPasswordByID(ctx context.Context, id string) (*model.User, error) {
	var userModel model.User
	if err := r.PSQL.QueryRowContext(ctx, GetUserPasswordByIDQuery, id).
		Scan(&userModel.ID, &userModel.RoleID, &userModel.Email, &userModel.Password, &userModel.Username, &userModel.PhoneNo,
			&userModel.FullName, &userModel.Gender, &userModel.BirthDate, &userModel.IsVerify, &userModel.PhotoURL); err != nil {
		return nil, err
	}

	return &userModel, nil
}

func (r *userRepo) GetPasswordByID(ctx context.Context, id string) (string, error) {
	var password string
	if err := r.PSQL.QueryRowContext(ctx, GetPasswordByIDQuery, id).
		Scan(&password); err != nil {
		return "", err
	}

	return password, nil
}

func (r *userRepo) CheckEmailHistory(ctx context.Context, email string) (*model.EmailHistory, error) {
	var emailHistory model.EmailHistory
	if err := r.PSQL.QueryRowContext(ctx, CheckEmailHistoryQuery, email).
		Scan(&emailHistory.ID, &emailHistory.Email); err != nil {
		return nil, err
	}

	return &emailHistory, nil
}
func (r *userRepo) GetUserByUsername(ctx context.Context, username string) (*model.User, error) {
	var userModel model.User
	if err := r.PSQL.QueryRowContext(ctx, GetUserByUsernameQuery, username).
		Scan(&userModel.ID, &userModel.Email, &userModel.Username, &userModel.IsVerify); err != nil {
		return nil, err
	}

	return &userModel, nil
}

func (r *userRepo) GetWalletByUserID(ctx context.Context, userID string) (*model.Wallet, error) {
	var walletModel model.Wallet
	if err := r.PSQL.QueryRowContext(ctx, GetWalletByUserIDQuery, userID).Scan(&walletModel.ID, &walletModel.UserID,
		&walletModel.Balance, &walletModel.PIN, &walletModel.AttemptCount,
		&walletModel.AttemptAt, &walletModel.UnlockedAt, &walletModel.ActiveDate); err != nil {
		return nil, err
	}

	return &walletModel, nil
}

func (r *userRepo) GetWalletHistoryByWalletID(ctx context.Context, pgn *pagination.Pagination,
	walletID string) ([]*body.HistoryWalletResponse, error) {
	queryOrderBySomething := fmt.Sprintf(OrderBySomething, pgn.GetSort(), pgn.GetLimit(),
		pgn.GetOffset())

	walletHistory := make([]*body.HistoryWalletResponse, 0)

	res, err := r.PSQL.QueryContext(
		ctx, GetWalletHistoryUserQuery+queryOrderBySomething,
		walletID,
	)

	if err != nil {
		return nil, err
	}

	for res.Next() {
		var history body.HistoryWalletResponse
		if errScan := res.Scan(
			&history.ID,
			&history.From,
			&history.To,
			&history.Amount,
			&history.Description,
			&history.CreatedAt,
		); errScan != nil {
			return nil, err
		}
		walletHistory = append(walletHistory, &history)
	}

	if res.Err() != nil {
		return nil, err
	}
	return walletHistory, nil
}

func (r *userRepo) GetWalletHistoryByID(ctx context.Context, id string) (*model.WalletHistory, error) {
	var walletHistory model.WalletHistory
	if err := r.PSQL.QueryRowContext(ctx, GetWalletHistoryByIDQuery, id).
		Scan(&walletHistory.ID, &walletHistory.TransactionID, &walletHistory.WalletID, &walletHistory.From, &walletHistory.To,
			&walletHistory.Amount, &walletHistory.Description, &walletHistory.CreatedAt); err != nil {
		return nil, err
	}

	return &walletHistory, nil
}

func (r *userRepo) GetTotalWalletHistoryByWalletID(ctx context.Context, walletID string) (int64, error) {
	var total int64
	if err := r.PSQL.QueryRowContext(ctx, GetTotalWalletHistoryUserQuery, walletID).
		Scan(&total); err != nil {
		return 0, err
	}

	return total, nil
}

func (r *userRepo) GetUserByPhoneNo(ctx context.Context, phoneNo string) (*model.User, error) {
	var userModel model.User
	if err := r.PSQL.QueryRowContext(ctx, GetUserByPhoneNoQuery, phoneNo).
		Scan(&userModel.ID, &userModel.Email, &userModel.PhoneNo, &userModel.IsVerify); err != nil {
		return nil, err
	}

	return &userModel, nil
}
func (r *userRepo) UpdateUserField(ctx context.Context, userModel *model.User) error {
	_, err := r.PSQL.ExecContext(
		ctx, UpdateUserFieldQuery, userModel.Username, userModel.FullName,
		userModel.PhoneNo, userModel.BirthDate, userModel.Gender, userModel.UpdatedAt, userModel.Email)
	if err != nil {
		return err
	}

	return nil
}

func (r *userRepo) UpdateUserEmail(ctx context.Context, tx postgre.Transaction, userModel *model.User) error {
	_, err := tx.ExecContext(
		ctx, UpdateUserEmailQuery, userModel.Email, userModel.UpdatedAt, userModel.ID)
	if err != nil {
		return err
	}

	return nil
}

func (r *userRepo) UpdateWalletPin(ctx context.Context, wallet *model.Wallet) error {
	_, err := r.PSQL.ExecContext(ctx, UpdateWalletPinQuery, wallet.PIN, wallet.ID)
	if err != nil {
		return err
	}

	return nil
}

func (r *userRepo) CreateWallet(ctx context.Context, wallet *model.Wallet) error {
	_, err := r.PSQL.ExecContext(ctx, CreateWalletQuery, wallet.UserID, wallet.Balance, wallet.PIN, wallet.AttemptCount, wallet.ActiveDate)
	if err != nil {
		return err
	}

	return nil
}

func (r *userRepo) CreateEmailHistory(ctx context.Context, tx postgre.Transaction, email string) error {
	_, err := tx.ExecContext(ctx, CreateEmailHistoryQuery, email)
	if err != nil {
		return err
	}

	return nil
}

func (r *userRepo) InsertNewOTPKey(ctx context.Context, email, otp string) error {
	key := fmt.Sprintf("%s:%s", constant.OtpKey, email)

	duration, err := time.ParseDuration(constant.OtpDuration)
	if err != nil {
		return err
	}

	if err := r.RedisClient.Set(ctx, key, otp, duration); err.Err() != nil {
		return err.Err()
	}

	return nil
}

func (r *userRepo) GetOTPValue(ctx context.Context, email string) (string, error) {
	key := fmt.Sprintf("%s:%s", constant.OtpKey, email)

	res := r.RedisClient.Get(ctx, key)
	if res.Err() != nil {
		return "", res.Err()
	}

	value, err := res.Result()
	if err != nil {
		return "", err
	}

	return value, nil
}

func (r *userRepo) DeleteOTPValue(ctx context.Context, email string) (int64, error) {
	key := fmt.Sprintf("%s:%s", constant.OtpKey, email)

	res := r.RedisClient.Del(ctx, key)
	if res.Err() != nil {
		return -1, res.Err()
	}

	value, err := res.Result()
	if err != nil {
		return -1, err
	}

	return value, nil
}

func (r *userRepo) GetSealabsPay(ctx context.Context, userid string) ([]*model.SealabsPay, error) {
	responses := make([]*model.SealabsPay, 0)
	res, err := r.PSQL.QueryContext(ctx, GetSealabsPayByIdQuery, userid)

	if err != nil {
		return nil, err
	}
	defer res.Close()

	for res.Next() {
		var response model.SealabsPay
		if err := res.Scan(
			&response.CardNumber, &response.UserID, &response.Name, &response.IsDefault,
			&response.ActiveDate, &response.CreatedAt, &response.UpdatedAt, &response.DeletedAt); err != nil {
			return nil, err
		}
		responses = append(responses, &response)
	}
	return responses, nil
}

func (r *userRepo) CheckDefaultSealabsPay(ctx context.Context, userid string) (*string, error) {
	var temp *string
	if err := r.PSQL.QueryRowContext(ctx, CheckDefaultSealabsPayQuery, userid).
		Scan(&temp); err != nil {
		return nil, err
	}

	return temp, nil
}

func (r *userRepo) PatchSealabsPay(ctx context.Context, cardNumber string) error {
	if _, err := r.PSQL.ExecContext(ctx, PatchSealabsPayQuery, cardNumber); err != nil {
		return err
	}
	return nil
}

func (r *userRepo) CheckUserSealabsPay(ctx context.Context, userid string) (int, error) {
	var temp int
	if err := r.PSQL.QueryRowContext(ctx, CheckUserSealabsPayQuery, userid).
		Scan(&temp); err != nil {
		return 0, err
	}

	return temp, nil
}

func (r *userRepo) CheckDeletedSealabsPay(ctx context.Context, cardNumber string) (int, error) {
	var temp int
	if err := r.PSQL.QueryRowContext(ctx, CheckDeletedSealabsPayQuery, cardNumber).
		Scan(&temp); err != nil {
		return 0, err
	}

	return temp, nil
}

func (r *userRepo) DeleteSealabsPay(ctx context.Context, cardNumber string) error {
	if _, err := r.PSQL.ExecContext(ctx, DeleteSealabsPayQuery, cardNumber); err != nil {
		return err
	}
	return nil
}

func (r *userRepo) SetDefaultSealabsPayTrans(ctx context.Context, tx postgre.Transaction, cardNumber *string) error {
	if _, err := tx.ExecContext(ctx, SetDefaultSealabsPayTransQuery, cardNumber); err != nil {
		return err
	}

	return nil
}

func (r *userRepo) SetDefaultSealabsPay(ctx context.Context, cardNumber, userid string) error {
	if _, err := r.PSQL.ExecContext(ctx, SetDefaultSealabsPayQuery, cardNumber, userid); err != nil {
		return err
	}

	return nil
}

func (r *userRepo) AddSealabsPayTrans(ctx context.Context, tx postgre.Transaction, request body.AddSealabsPayRequest, userid string) error {
	if _, err := tx.ExecContext(ctx, CreateSealabsPayQuery, request.CardNumber, userid,
		request.Name, request.IsDefault, request.ActiveDateTime); err != nil {
		return err
	}

	return nil
}

func (r *userRepo) UpdateUserSealabsPayTrans(ctx context.Context, tx postgre.Transaction, request body.AddSealabsPayRequest, userid string) error {
	if _, err := tx.ExecContext(ctx, UpdateUserSealabsPayQuery, userid, request.Name, request.CardNumber); err != nil {
		return err
	}

	return nil
}

func (r *userRepo) AddSealabsPay(ctx context.Context, request body.AddSealabsPayRequest, userid string) error {
	if _, err := r.PSQL.ExecContext(ctx, CreateSealabsPayQuery, request.CardNumber, userid,
		request.Name, request.IsDefault, request.ActiveDateTime); err != nil {
		return err
	}

	return nil
}

func (r *userRepo) UpdateUserSealabsPay(ctx context.Context, request body.AddSealabsPayRequest, userid string) error {
	if _, err := r.PSQL.ExecContext(ctx, UpdateUserSealabsPayQuery, userid, request.Name, request.CardNumber); err != nil {
		return err
	}

	return nil
}

func (r *userRepo) CheckShopByID(ctx context.Context, userID string) (int64, error) {
	var result int64
	err := r.PSQL.QueryRowContext(ctx, CheckShopByIdQuery, userID).Scan(&result)
	if err != nil {
		return -1, err
	}

	return result, nil
}

func (r *userRepo) CheckShopUnique(ctx context.Context, shopName string) (int64, error) {
	var result int64
	err := r.PSQL.QueryRowContext(ctx, CheckShopUniqueQuery, shopName).Scan(&result)
	if err != nil {
		return -1, err
	}

	return result, nil
}

func (r *userRepo) AddShop(ctx context.Context, userID, shopName string) error {
	if _, err := r.PSQL.ExecContext(ctx, AddShopQuery, userID, shopName); err != nil {
		return err
	}

	return nil
}

func (r *userRepo) UpdateRole(ctx context.Context, userID string) error {
	if _, err := r.PSQL.ExecContext(ctx, UpdateRoleQuery, userID); err != nil {
		return err
	}

	return nil
}

func (r *userRepo) UpdateProfileImage(ctx context.Context, imgURL, userID string) error {
	if _, err := r.PSQL.ExecContext(ctx, UpdateProfileImageQuery, imgURL, userID); err != nil {
		return err
	}

	return nil
}

func (r *userRepo) UpdatePasswordByID(ctx context.Context, userID, newPassword string) error {
	_, err := r.PSQL.ExecContext(ctx, UpdatePasswordQuery, newPassword, userID)
	if err != nil {
		return err
	}
	return nil
}

func (r *userRepo) GetWalletUser(ctx context.Context, walletID string) (*model.Wallet, error) {
	var walletUser model.Wallet
	if err := r.PSQL.QueryRowContext(ctx, GetWalletUserQuery, walletID).Scan(
		&walletUser.ID,
		&walletUser.UserID,
		&walletUser.Balance,
		&walletUser.AttemptCount,
		&walletUser.AttemptAt,
		&walletUser.UnlockedAt,
		&walletUser.ActiveDate); err != nil {
		return nil, err
	}

	return &walletUser, nil
}

func (r *userRepo) GetSealabsPayUser(ctx context.Context, userID, cardNumber string) (*model.SealabsPay, error) {
	var sealabsPayUser model.SealabsPay
	if err := r.PSQL.QueryRowContext(ctx, GetSealabsPayUserQuery, userID, cardNumber).Scan(
		&sealabsPayUser.CardNumber,
		&sealabsPayUser.UserID,
		&sealabsPayUser.Name,
		&sealabsPayUser.IsDefault,
		&sealabsPayUser.ActiveDate); err != nil {
		return nil, err
	}

	return &sealabsPayUser, nil
}

func (r *userRepo) GetVoucherMarketplaceByID(ctx context.Context, voucherMarketplaceID string) (*model.Voucher, error) {
	var VoucherMarketplace model.Voucher
	if err := r.PSQL.QueryRowContext(ctx, GetVoucherMarketplaceByIDQuery, voucherMarketplaceID).Scan(
		&VoucherMarketplace.ID,
		&VoucherMarketplace.ShopID,
		&VoucherMarketplace.Code,
		&VoucherMarketplace.Quota,
		&VoucherMarketplace.ActivedDate,
		&VoucherMarketplace.ExpiredDate,
		&VoucherMarketplace.DiscountPercentage,
		&VoucherMarketplace.DiscountFixPrice,
		&VoucherMarketplace.MinProductPrice,
		&VoucherMarketplace.MaxDiscountPrice); err != nil {
		return nil, err
	}

	return &VoucherMarketplace, nil
}

func (r *userRepo) GetShopByID(ctx context.Context, shopID string) (*model.Shop, error) {
	var shopCart model.Shop
	if err := r.PSQL.QueryRowContext(ctx, GetShopByIDQuery, shopID).Scan(
		&shopCart.ID,
		&shopCart.Name,
		&shopCart.UserID); err != nil {
		return nil, err
	}

	return &shopCart, nil
}

func (r *userRepo) GetVoucherShopByID(ctx context.Context, voucherShopID, shopID string) (*model.Voucher, error) {
	var VoucherShop model.Voucher
	if err := r.PSQL.QueryRowContext(ctx, GetVoucherShopByIDQuery, voucherShopID, shopID).Scan(
		&VoucherShop.ID,
		&VoucherShop.ShopID,
		&VoucherShop.Code,
		&VoucherShop.Quota,
		&VoucherShop.ActivedDate,
		&VoucherShop.ExpiredDate,
		&VoucherShop.DiscountPercentage,
		&VoucherShop.DiscountFixPrice,
		&VoucherShop.MinProductPrice,
		&VoucherShop.MaxDiscountPrice); err != nil {
		return nil, err
	}

	return &VoucherShop, nil
}

func (r *userRepo) GetCourierShopByID(ctx context.Context, courierID, shopID string) (*model.Courier, error) {
	var CourierShop model.Courier
	if err := r.PSQL.QueryRowContext(ctx, GetCourierShopByIDQuery, courierID, shopID).Scan(
		&CourierShop.ID,
		&CourierShop.Name,
		&CourierShop.Code,
		&CourierShop.Service,
		&CourierShop.Description); err != nil {
		return nil, err
	}

	return &CourierShop, nil
}

func (r *userRepo) GetProductDetailByID(ctx context.Context, tx postgre.Transaction, productDetailID string) (*model.ProductDetail, error) {
	var pd model.ProductDetail
	if err := tx.QueryRowContext(ctx, GetProductDetailByIDQuery, productDetailID).Scan(
		&pd.ID,
		&pd.ProductID,
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

func (r *userRepo) CreateTransaction(ctx context.Context, tx postgre.Transaction,
	transactionData *model.Transaction) (*uuid.UUID, error) {
	var transactionID *uuid.UUID
	if err := tx.QueryRowContext(ctx, CreateTransactionQuery,
		transactionData.VoucherMarketplaceID,
		transactionData.WalletID,
		transactionData.CardNumber,
		transactionData.Invoice,
		transactionData.TotalPrice,
		transactionData.ExpiredAt).Scan(&transactionID); err != nil {
		return nil, err
	}

	return transactionID, nil
}

func (r *userRepo) CreateOrder(ctx context.Context, tx postgre.Transaction, orderData *model.OrderModel) (*uuid.UUID, error) {
	var orderID *uuid.UUID
	if err := tx.QueryRowContext(ctx, CreateOrderQuery,
		orderData.TransactionID,
		orderData.ShopID,
		orderData.UserID,
		orderData.CourierID,
		orderData.VoucherShopID,
		orderData.OrderStatusID,
		orderData.TotalPrice,
		orderData.DeliveryFee,
		orderData.BuyerAddress,
		orderData.ShopAddress).Scan(&orderID); err != nil {
		return nil, err
	}

	return orderID, nil
}

func (r *userRepo) CreateOrderItem(ctx context.Context, tx postgre.Transaction, item *model.OrderItem) (*uuid.UUID, error) {
	var orderItemID *uuid.UUID
	if err := tx.QueryRowContext(ctx, CreateOrderItemQuery,
		item.OrderID,
		item.ProductDetailID,
		item.Quantity,
		item.ItemPrice,
		item.TotalPrice,
		item.Note).Scan(&orderItemID); err != nil {
		return nil, err
	}

	return orderItemID, nil
}

func (r *userRepo) GetCartItemUser(ctx context.Context, userID, productDetailID string) (*model.CartItem, error) {
	var CartItemResult model.CartItem
	if err := r.PSQL.QueryRowContext(ctx, GetCartItemUserQuery,
		userID, productDetailID).Scan(
		&CartItemResult.ID,
		&CartItemResult.UserID,
		&CartItemResult.ProductDetailID,
		&CartItemResult.Quantity); err != nil {
		return nil, err
	}

	return &CartItemResult, nil
}

func (r *userRepo) UpdateProductDetailStock(ctx context.Context, tx postgre.Transaction,
	productDetailData *model.ProductDetail) error {
	_, err := tx.ExecContext(ctx, UpdateProductDetailStockQuery, productDetailData.Stock, productDetailData.ID)
	if err != nil {
		return err
	}

	return nil
}

func (r *userRepo) UpdateWalletBalance(ctx context.Context, tx postgre.Transaction, wallet *model.Wallet) error {
	_, err := tx.ExecContext(ctx, UpdateWalletBalanceQuery, wallet.Balance, wallet.UpdatedAt, wallet.ID)
	if err != nil {
		return err
	}

	return nil
}

func (r *userRepo) UpdateWallet(ctx context.Context, wallet *model.Wallet) error {
	_, err := r.PSQL.ExecContext(ctx, UpdateWalletQuery, wallet.AttemptCount, wallet.AttemptAt, wallet.UnlockedAt, wallet.ID)
	if err != nil {
		return err
	}

	return nil
}

func (r *userRepo) InsertWalletHistory(ctx context.Context, tx postgre.Transaction, walletHistory *model.WalletHistory) error {
	_, err := tx.ExecContext(ctx, CreateWalletHistoryQuery, walletHistory.TransactionID, walletHistory.WalletID,
		walletHistory.From, walletHistory.To, walletHistory.Description, walletHistory.Amount, walletHistory.CreatedAt)
	if err != nil {
		return err
	}
	return nil
}

func (r *userRepo) ChangeOrderStatus(ctx context.Context, requestBody body.ChangeOrderStatusRequest) error {
	_, err := r.PSQL.ExecContext(
		ctx, UpdateOrderByID, requestBody.OrderStatusID, requestBody.OrderID)
	if err != nil {
		return err
	}

	return nil
}

func (r *userRepo) GetOrderByTransactionID(ctx context.Context, transactionID string) ([]*model.OrderModel, error) {
	orders := make([]*model.OrderModel, 0)
	res, err := r.PSQL.QueryContext(ctx, GetOrderByTransactionID, transactionID)
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

func (r *userRepo) UpdateTransaction(ctx context.Context, tx postgre.Transaction, transactionData *model.Transaction) error {
	_, err := tx.ExecContext(ctx, UpdateTransactionByID, transactionData.PaidAt, transactionData.CanceledAt,
		transactionData.CardNumber, transactionData.ID)
	if err != nil {
		return err
	}

	return nil
}

func (r *userRepo) UpdateOrder(ctx context.Context, tx postgre.Transaction, orderData *model.OrderModel) error {
	_, err := tx.ExecContext(ctx, UpdateOrderByID, orderData.OrderStatusID, orderData.ID)
	if err != nil {
		return err
	}

	return nil
}

func (r *userRepo) DeleteCartItemByID(ctx context.Context, tx postgre.Transaction, cartItemData *model.CartItem) error {
	_, err := tx.ExecContext(ctx, DeleteCartItemByIDQuery, cartItemData.ID.String())
	if err != nil {
		return err
	}
	return nil
}

func (r *userRepo) GetProductPromotionByProductID(ctx context.Context, productID string) (*model.Promotion, error) {
	var promo model.Promotion
	if err := r.PSQL.QueryRowContext(ctx, GetProductPromotionByProductIDQuery, productID).
		Scan(
			&promo.ID,
			&promo.Name,
			&promo.ProductID,
			&promo.DiscountPercentage,
			&promo.DiscountFixPrice,
			&promo.MinProductPrice,
			&promo.MaxDiscountPrice,
			&promo.Quota,
			&promo.MaxQuantity,
			&promo.ActivedDate,
			&promo.ExpiredDate,
			&promo.CreatedAt,
			&promo.UpdatedAt,
			&promo.DeletedAt); err != nil {
		return nil, err
	}

	return &promo, nil
}

func (r *userRepo) UpdateVoucherQuota(ctx context.Context, tx postgre.Transaction, upVoucher *model.Voucher) error {
	_, err := tx.ExecContext(ctx, UpdateVoucherQuotaQuery, upVoucher.Quota, upVoucher.ID)
	if err != nil {
		return err
	}
	return nil
}

func (r *userRepo) UpdatePromotionQuota(ctx context.Context, tx postgre.Transaction, promo *model.Promotion) error {
	_, err := tx.ExecContext(ctx, UpdatePromotionQuotaQuery, promo.Quota, promo.ID)
	if err != nil {
		return err
	}
	return nil
}
