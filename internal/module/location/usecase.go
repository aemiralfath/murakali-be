package location

import (
	"context"
	"murakali/internal/module/location/delivery/body"
)

type UseCase interface {
	GetProvince(ctx context.Context) (*body.ProvinceResponse, error)
	GetCity(ctx context.Context, provinceID int) (*body.CityResponse, error)
	GetSubDistrict(ctx context.Context, province, city string) (*body.SubDistrictResponse, error)
}
