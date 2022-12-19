package user

import (
	"context"
	"murakali/internal/model"
	"murakali/internal/module/user/delivery/body"
	"murakali/pkg/pagination"
)

type UseCase interface {
	GetAddress(ctx context.Context, userID, name string, pagination *pagination.Pagination) (*pagination.Pagination, error)
	CreateAddress(ctx context.Context, userID string, requestBody body.CreateAddressRequest) error
	GetAddressByID(ctx context.Context, userID, addressID string) (*model.Address, error)
	UpdateAddressByID(ctx context.Context, userID, addressID string, requestBody body.UpdateAddressRequest) error
	DeleteAddressByID(ctx context.Context, userID, addressID string) error
}
