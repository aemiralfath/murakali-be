package repository

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/go-redis/redis/v8"
	"murakali/internal/constant"
	"murakali/internal/model"
	"murakali/internal/module/location"
)

type locationRepo struct {
	PSQL        *sql.DB
	RedisClient *redis.Client
}

func NewLocationRepository(psql *sql.DB, client *redis.Client) location.Repository {
	return &locationRepo{PSQL: psql, RedisClient: client}
}

func (r *locationRepo) GetShopCourierID(ctx context.Context, shopID string) ([]string, error) {
	courierIDS := make([]string, 0)
	res, err := r.PSQL.QueryContext(ctx, GetShopCourierID, shopID)
	if err != nil {
		return nil, err
	}
	defer res.Close()

	for res.Next() {
		var shopCourierID string
		if errScan := res.Scan(&shopCourierID); errScan != nil {
			return nil, errScan
		}

		courierIDS = append(courierIDS, shopCourierID)
	}

	if res.Err() != nil {
		return nil, res.Err()
	}

	return courierIDS, nil
}

func (r *locationRepo) GetProductCourierWhitelistID(ctx context.Context, productID string) ([]string, error) {
	courierIDS := make([]string, 0)
	res, err := r.PSQL.QueryContext(ctx, GetProductCourierWhitelistID, productID)
	if err != nil {
		return nil, err
	}
	defer res.Close()

	for res.Next() {
		var productCourierID string
		if errScan := res.Scan(&productCourierID); errScan != nil {
			return nil, errScan
		}

		courierIDS = append(courierIDS, productCourierID)
	}

	if res.Err() != nil {
		return nil, res.Err()
	}

	return courierIDS, nil
}

func (r *locationRepo) GetCourierByID(ctx context.Context, courierID string) (*model.Courier, error) {
	var courier model.Courier
	if err := r.PSQL.QueryRowContext(ctx, GetCourierByID, courierID).Scan(
		&courier.ID, &courier.Name, &courier.Code, &courier.Service,
		&courier.Description, &courier.CreatedAt, &courier.UpdatedAt); err != nil {
		return nil, err
	}

	return &courier, nil
}

func (r *locationRepo) GetShopByID(ctx context.Context, shopID string) (*model.Shop, error) {
	var shop model.Shop
	if err := r.PSQL.QueryRowContext(ctx, GetShopByID, shopID).Scan(&shop.ID, &shop.UserID); err != nil {
		return nil, err
	}
	return &shop, nil
}

func (r *locationRepo) GetShopAddress(ctx context.Context, userID string) (*model.Address, error) {
	var shopAddress model.Address
	if err := r.PSQL.QueryRowContext(ctx, GetShopAddress, userID).Scan(&shopAddress.ID, &shopAddress.UserID, &shopAddress.CityID); err != nil {
		return nil, err
	}
	return &shopAddress, nil
}

func (r *locationRepo) InsertProvinceRedis(ctx context.Context, value string) error {
	if err := r.RedisClient.Set(ctx, constant.ProvinceKey, value, 0); err.Err() != nil {
		return err.Err()
	}

	return nil
}

func (r *locationRepo) InsertCityRedis(ctx context.Context, provinceID int, value string) error {
	if err := r.RedisClient.Set(ctx, fmt.Sprintf("%s:%d", constant.CityKey, provinceID), value, 0); err.Err() != nil {
		return err.Err()
	}

	return nil
}

func (r *locationRepo) InsertSubDistrictRedis(ctx context.Context, province, city, value string) error {
	if err := r.RedisClient.Set(ctx, fmt.Sprintf("%s:%s:%s", constant.SubDistrictKey, province, city), value, 0); err.Err() != nil {
		return err.Err()
	}

	return nil
}

func (r *locationRepo) InsertUrbanRedis(ctx context.Context, province, city, subDistrict, value string) error {
	if err := r.RedisClient.Set(ctx, fmt.Sprintf("%s:%s:%s:%s", constant.UrbanKey, province, city, subDistrict), value, 0); err.Err() != nil {
		return err.Err()
	}

	return nil
}

func (r *locationRepo) InsertCostRedis(ctx context.Context, key, value string) error {
	if err := r.RedisClient.Set(ctx, key, value, 0); err.Err() != nil {
		return err.Err()
	}

	return nil
}

func (r *locationRepo) GetProvinceRedis(ctx context.Context) (string, error) {
	res := r.RedisClient.Get(ctx, constant.ProvinceKey)
	if res.Err() != nil {
		return "", res.Err()
	}

	value, err := res.Result()
	if err != nil {
		return "", err
	}

	return value, nil
}

func (r *locationRepo) GetCityRedis(ctx context.Context, provinceID int) (string, error) {
	res := r.RedisClient.Get(ctx, fmt.Sprintf("%s:%d", constant.CityKey, provinceID))
	if res.Err() != nil {
		return "", res.Err()
	}

	value, err := res.Result()
	if err != nil {
		return "", err
	}

	return value, nil
}

func (r *locationRepo) GetSubDistrictRedis(ctx context.Context, province, city string) (string, error) {
	res := r.RedisClient.Get(ctx, fmt.Sprintf("%s:%s:%s", constant.SubDistrictKey, province, city))
	if res.Err() != nil {
		return "", res.Err()
	}

	value, err := res.Result()
	if err != nil {
		return "", err
	}

	return value, nil
}

func (r *locationRepo) GetUrbanRedis(ctx context.Context, province, city, subDistrict string) (string, error) {
	res := r.RedisClient.Get(ctx, fmt.Sprintf("%s:%s:%s:%s", constant.UrbanKey, province, city, subDistrict))
	if res.Err() != nil {
		return "", res.Err()
	}

	value, err := res.Result()
	if err != nil {
		return "", err
	}

	return value, nil
}

func (r *locationRepo) GetCostRedis(ctx context.Context, key string) (*string, error) {
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
