package repository

import (
	"context"
	"database/sql"
	"fmt"
	"murakali/internal/model"
	"murakali/internal/module/admin"
	"murakali/pkg/pagination"

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

func (r *adminRepo) DeleteVoucher(ctx context.Context, voucherID string) error {
	_, err := r.PSQL.ExecContext(ctx, DeleteVoucherQuery, voucherID)
	if err != nil {
		return err
	}
	return nil
}
