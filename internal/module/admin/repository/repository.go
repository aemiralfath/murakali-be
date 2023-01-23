package repository

import (
	"context"
	"database/sql"
	"fmt"
	"murakali/internal/model"
	"murakali/internal/module/admin"
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
		q = q + FilterVoucherWillCome
	case "3":
		q = q + FilterVoucherOngoing
	case "4":
		q = q + FilterVoucherHasEnded
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
		q = q + FilterVoucherWillCome
	case "3":
		q = q + FilterVoucherOngoing
	case "4":
		q = q + FilterVoucherHasEnded
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

func (r *adminRepo) GetWalletByUserID(ctx context.Context, userID string) (*model.Wallet, error) {
	var walletModel model.Wallet
	if err := r.PSQL.QueryRowContext(ctx, GetWalletByUserIDQuery, userID).Scan(&walletModel.ID, &walletModel.UserID,
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
