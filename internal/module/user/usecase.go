package user

import (
	"context"
	"murakali/internal/module/user/delivery/body"
	"murakali/pkg/pagination"
)

type UseCase interface {
	GetAddress(ctx context.Context, userID, name string, pagination *pagination.Pagination) (*pagination.Pagination, error)
	CreateAddress(ctx context.Context, userID string, requestBody body.CreateAddressRequest) error
}
