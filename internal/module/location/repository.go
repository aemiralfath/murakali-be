package location

import "context"

type Repository interface {
	GetProvinceRedis(ctx context.Context) (string, error)
	InsertProvinceRedis(ctx context.Context, value string) error
}
