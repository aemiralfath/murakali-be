package repository

import (
	"context"
	"database/sql"
	"fmt"
	"murakali/internal/model"
	"murakali/internal/module/admin"
	"murakali/internal/module/admin/delivery/body"
	"murakali/pkg/pagination"
	"murakali/pkg/postgre"

	"github.com/go-redis/redis/v8"
)

type adminRepo struct {
	PSQL        *sql.DB
	RedisClient *redis.Client
}

func NewAdminRepository(psql *sql.DB, client *redis.Client) admin.Repository {
	return &adminRepo{
		PSQL:        psql,
		RedisClient: client,
	}
}

func (r *adminRepo) CountCodeVoucher(ctx context.Context, code string) (int64, error) {
	var total int64

	if err := r.PSQL.QueryRowContext(ctx, CountCodeVoucher, code).Scan(&total); err != nil {
		return -1, err
	}

	return total, nil
}

func (r *adminRepo) GetTotalVoucher(ctx context.Context, voucherStatusID string) (int64, error) {
	var total int64

	q := GetTotalVoucherQuery
	switch voucherStatusID {
	case "1":
		q = GetTotalVoucherQuery
	case "2":
		q += FilterVoucherWillCome
	case "3":
		q += FilterVoucherOngoing
	case "4":
		q += FilterVoucherHasEnded
	default:
		q = GetTotalVoucherQuery
	}

	if err := r.PSQL.QueryRowContext(ctx, q).Scan(&total); err != nil {
		return -1, err
	}

	return total, nil
}

func (r *adminRepo) GetTotalRefunds(ctx context.Context) (int64, error) {
	var total int64
	if err := r.PSQL.QueryRowContext(ctx, GetTotalRefundsQuery).Scan(&total); err != nil {
		return -1, err
	}

	return total, nil
}

func (r *adminRepo) GetAllVoucher(ctx context.Context, voucherStatusID, sortFilter string, pgn *pagination.Pagination) ([]*model.Voucher, error) {
	var vouchers []*model.Voucher

	q := GetAllVoucherQuery
	switch voucherStatusID {
	case "1":
		q = GetAllVoucherQuery
	case "2":
		q += FilterVoucherWillCome
	case "3":
		q += FilterVoucherOngoing
	case "4":
		q += FilterVoucherHasEnded
	default:
		q = GetAllVoucherQuery
	}
	queryOrderBySomething := fmt.Sprintf(OrderBySomething, sortFilter, pgn.GetLimit(),
		pgn.GetOffset())
	res, err := r.PSQL.QueryContext(ctx, q+queryOrderBySomething)
	if err != nil {
		return nil, err
	}
	defer res.Close()

	for res.Next() {
		var voucher model.Voucher
		if errScan := res.Scan(
			&voucher.ID,
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

		vouchers = append(vouchers, &voucher)
	}

	if res.Err() != nil {
		return nil, err
	}

	return vouchers, nil
}

func (r *adminRepo) GetRefunds(ctx context.Context, sortFilter string, pgn *pagination.Pagination) ([]*model.RefundOrder, error) {
	refunds := make([]*model.RefundOrder, 0)

	queryOrderBySomething := fmt.Sprintf(OrderBySomething, sortFilter, pgn.GetLimit(), pgn.GetOffset())
	res, err := r.PSQL.QueryContext(ctx, GetRefundsQuery+queryOrderBySomething)
	if err != nil {
		return nil, err
	}
	defer res.Close()

	for res.Next() {
		refund := model.RefundOrder{}
		refund.Order = &model.OrderModel{}
		if errScan := res.Scan(
			&refund.ID,
			&refund.OrderID,
			&refund.IsSellerRefund,
			&refund.IsBuyerRefund,
			&refund.Reason,
			&refund.Image,
			&refund.AcceptedAt,
			&refund.RejectedAt,
			&refund.RefundedAt,
			&refund.Order.ID,
			&refund.Order.TransactionID,
			&refund.Order.ShopID,
			&refund.Order.UserID,
			&refund.Order.CourierID,
			&refund.Order.VoucherShopID,
			&refund.Order.OrderStatusID,
			&refund.Order.TotalPrice,
			&refund.Order.DeliveryFee,
			&refund.Order.ResiNo,
			&refund.Order.CreatedAt,
			&refund.Order.ArrivedAt); errScan != nil {
			return nil, errScan
		}

		refunds = append(refunds, &refund)
	}
	if res.Err() != nil {
		return nil, err
	}

	return refunds, nil
}

func (r *adminRepo) GetVoucherByID(ctx context.Context, voucherID string) (*model.Voucher, error) {
	var voucher model.Voucher
	if err := r.PSQL.QueryRowContext(ctx, GetVoucherByID, voucherID).Scan(
		&voucher.ID,
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

func (r *adminRepo) GetRefundByID(ctx context.Context, refundID string) (*model.Refund, error) {
	var refund model.Refund
	if err := r.PSQL.QueryRowContext(ctx, GetRefundByIDQuery, refundID).Scan(
		&refund.ID,
		&refund.OrderID,
		&refund.IsSellerRefund,
		&refund.IsBuyerRefund,
		&refund.Reason,
		&refund.Image,
		&refund.AcceptedAt,
		&refund.RejectedAt,
		&refund.RefundedAt); err != nil {
		return nil, err
	}

	return &refund, nil
}

func (r *adminRepo) GetOrderByID(ctx context.Context, orderID string) (*model.OrderModel, error) {
	var order model.OrderModel
	if err := r.PSQL.QueryRowContext(ctx, GetOrderByOrderIDQuery, orderID).Scan(
		&order.ID,
		&order.OrderStatusID,
		&order.UserID,
		&order.TransactionID,
		&order.TotalPrice,
		&order.DeliveryFee,
		&order.ResiNo,
		&order.CreatedAt); err != nil {
		return nil, err
	}

	return &order, nil
}

func (r *adminRepo) CreateVoucher(ctx context.Context, voucherShop *model.Voucher) error {
	if _, err := r.PSQL.ExecContext(ctx, CreateVoucherQuery,
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

func (r *adminRepo) UpdateVoucher(ctx context.Context, voucherShop *model.Voucher) error {
	if _, err := r.PSQL.ExecContext(ctx, UpdateVoucherQuery,
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

func (r *adminRepo) UpdateRefund(ctx context.Context, tx postgre.Transaction, refund *model.Refund) error {
	if _, err := tx.ExecContext(ctx, UpdateRefundQuery, refund.RefundedAt, refund.ID); err != nil {
		return err
	}

	return nil
}

func (r *adminRepo) UpdateOrderStatus(ctx context.Context, tx postgre.Transaction, orderData *model.OrderModel) error {
	_, err := tx.ExecContext(ctx, UpdateOrderByID, orderData.OrderStatusID, orderData.ID)
	if err != nil {
		return err
	}

	return nil
}

func (r *adminRepo) GetOrderItemsByOrderID(ctx context.Context, tx postgre.Transaction, orderID string) ([]*model.OrderItem, error) {
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

func (r *adminRepo) GetProductDetailByID(ctx context.Context, tx postgre.Transaction, productDetailID string) (*model.ProductDetail, error) {
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

func (r *adminRepo) GetWalletByUserID(ctx context.Context, tx postgre.Transaction, userID string) (*model.Wallet, error) {
	var walletModel model.Wallet
	if err := tx.QueryRowContext(ctx, GetWalletByUserIDQuery, userID).Scan(&walletModel.ID, &walletModel.UserID,
		&walletModel.Balance, &walletModel.PIN, &walletModel.AttemptCount,
		&walletModel.AttemptAt, &walletModel.UnlockedAt, &walletModel.ActiveDate); err != nil {
		return nil, err
	}

	return &walletModel, nil
}

func (r *adminRepo) InsertWalletHistory(ctx context.Context, tx postgre.Transaction, walletHistory *model.WalletHistory) error {
	_, err := tx.ExecContext(ctx, CreateWalletHistoryQuery, walletHistory.TransactionID, walletHistory.WalletID,
		walletHistory.From, walletHistory.To, walletHistory.Description, walletHistory.Amount, walletHistory.CreatedAt)
	if err != nil {
		return err
	}
	return nil
}

func (r *adminRepo) UpdateWalletBalance(ctx context.Context, tx postgre.Transaction, wallet *model.Wallet) error {
	_, err := tx.ExecContext(ctx, UpdateWalletBalanceQuery, wallet.Balance, wallet.UpdatedAt, wallet.ID)
	if err != nil {
		return err
	}

	return nil
}

func (r *adminRepo) UpdateProductDetailStock(ctx context.Context, tx postgre.Transaction,
	productDetailData *model.ProductDetail) error {
	_, err := tx.ExecContext(ctx, UpdateProductDetailStockQuery, productDetailData.Stock, productDetailData.ID)
	if err != nil {
		return err
	}

	return nil
}

func (r *adminRepo) DeleteVoucher(ctx context.Context, voucherID string) error {
	_, err := r.PSQL.ExecContext(ctx, DeleteVoucherQuery, voucherID)
	if err != nil {
		return err
	}
	return nil
}

func (r *adminRepo) GetCategories(ctx context.Context) ([]*body.CategoryResponse, error) {
	var categories = make([]*body.CategoryResponse, 0)
	res, err := r.PSQL.QueryContext(ctx, GetCategoriesQuery)
	if err != nil {
		return nil, err
	}

	for res.Next() {
		category := body.CategoryResponse{}
		if errScan := res.Scan(&category.CategoryID, &category.ParentID, &category.Name, &category.PhotoURL, &category.Level); errScan != nil {
			return nil, errScan
		}
		categories = append(categories, &category)
	}

	return categories, nil
}

func (r *adminRepo) AddCategory(ctx context.Context, requestBody body.CategoryRequest) error {
	if _, err := r.PSQL.ExecContext(ctx, AddCategoryQuery,
		&requestBody.ParentIDValue, requestBody.Name, requestBody.PhotoURL); err != nil {
		return err
	}
	return nil
}

func (r *adminRepo) DeleteCategory(ctx context.Context, categoryID string) error {
	_, err := r.PSQL.ExecContext(ctx, DeleteCategoryQuery, categoryID)
	if err != nil {
		return err
	}
	return nil
}

func (r *adminRepo) EditCategory(ctx context.Context, requestBody body.CategoryRequest) error {
	_, err := r.PSQL.ExecContext(ctx, EditCategoryQuery, requestBody.ParentIDValue, requestBody.Name, requestBody.PhotoURL, requestBody.ID)
	if err != nil {
		return err
	}
	return nil
}

func (r *adminRepo) CountProductCategory(ctx context.Context, userid string) (int, error) {
	var temp int
	if err := r.PSQL.QueryRowContext(ctx, CountProductCategoryQuery, userid).
		Scan(&temp); err != nil {
		return 0, err
	}

	return temp, nil
}

func (r *adminRepo) CountCategoryParent(ctx context.Context, userid string) (int, error) {
	var temp int
	if err := r.PSQL.QueryRowContext(ctx, CountCategoryParentQuery, userid).
		Scan(&temp); err != nil {
		return 0, err
	}

	return temp, nil
}

func (r *adminRepo) GetBanner(ctx context.Context) ([]*body.BannerResponse, error) {
	var banners = make([]*body.BannerResponse, 0)
	res, err := r.PSQL.QueryContext(ctx, GetBannerQuery)
	if err != nil {
		return nil, err
	}

	for res.Next() {
		banner := body.BannerResponse{}
		if errScan := res.Scan(&banner.ID, &banner.Title, &banner.Content, &banner.ImageURL, &banner.PageURL, &banner.IsActive); errScan != nil {
			return nil, errScan
		}
		banners = append(banners, &banner)
	}

	return banners, nil
}

func (r *adminRepo) AddBanner(ctx context.Context, requestBody body.BannerRequest) error {
	if _, err := r.PSQL.ExecContext(ctx, AddBannerQuery,
		&requestBody.Title, requestBody.Content, requestBody.ImageURL, requestBody.PageURL, requestBody.IsActive); err != nil {
		return err
	}
	return nil
}

func (r *adminRepo) DeleteBanner(ctx context.Context, bannerID string) error {
	_, err := r.PSQL.ExecContext(ctx, DeleteBannerQuery, bannerID)
	if err != nil {
		return err
	}
	return nil
}
