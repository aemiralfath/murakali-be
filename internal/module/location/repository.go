package location

import (
	"context"
	"murakali/internal/model"
)

type Repository interface {
	GetProvinceRedis(ctx context.Context) (string, error)
	InsertProvinceRedis(ctx context.Context, value string) error
	GetCityRedis(ctx context.Context, provinceID int) (string, error)
	InsertCityRedis(ctx context.Context, provinceID int, value string) error
	GetSubDistrictRedis(ctx context.Context, province, city string) (string, error)
	InsertSubDistrictRedis(ctx context.Context, province, city, value string) error
	GetUrbanRedis(ctx context.Context, province, city, subDistrict string) (string, error)
	InsertUrbanRedis(ctx context.Context, province, city, subDistrict, value string) error
	GetShopCourierID(ctx context.Context, shopID string) ([]string, error)
	GetProductCourierWhitelistID(ctx context.Context, productID string) ([]string, error)
	GetCourierByID(ctx context.Context, courierID string) (*model.Courier, error)
	GetCostRedis(ctx context.Context, key string) (*string, error)
	InsertCostRedis(ctx context.Context, key string, value string) error
	GetShopByID(ctx context.Context, shopID string) (*model.Shop, error)
	GetShopAddress(ctx context.Context, userID string) (*model.Address, error)
}
