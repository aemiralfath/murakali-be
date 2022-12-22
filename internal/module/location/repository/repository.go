package repository

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/go-redis/redis/v8"
	"murakali/internal/constant"
	"murakali/internal/module/location"
)

type locationRepo struct {
	PSQL        *sql.DB
	RedisClient *redis.Client
}

func NewLocationRepository(psql *sql.DB, client *redis.Client) location.Repository {
	return &locationRepo{PSQL: psql, RedisClient: client}
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
