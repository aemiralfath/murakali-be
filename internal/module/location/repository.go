package location

import "context"

type Repository interface {
	GetProvinceRedis(ctx context.Context) (string, error)
	InsertProvinceRedis(ctx context.Context, value string) error
	GetCityRedis(ctx context.Context, provinceID int) (string, error)
	InsertCityRedis(ctx context.Context, provinceID int, value string) error
}
