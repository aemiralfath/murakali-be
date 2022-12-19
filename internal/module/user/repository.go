package user

import (
	"context"
	"murakali/internal/model"
	"murakali/internal/module/user/delivery/body"
	"murakali/pkg/pagination"
	"murakali/pkg/postgre"
)

type Repository interface {
	GetUserByID(ctx context.Context, id string) (*model.User, error)
	GetTotalAddress(ctx context.Context, userID, name string) (int64, error)
	GetAddresses(ctx context.Context, userID, name string, pagination *pagination.Pagination) ([]*model.Address, error)
	GetAddressByID(ctx context.Context, userID, addressID string) (*model.Address, error)
	GetDefaultUserAddress(ctx context.Context, userID string) (*model.Address, error)
	GetDefaultShopAddress(ctx context.Context, userID string) (*model.Address, error)
	CreateAddress(ctx context.Context, tx postgre.Transaction, userID string, requestBody body.CreateAddressRequest) error
	UpdateAddress(ctx context.Context, tx postgre.Transaction, address *model.Address) error
	UpdateDefaultAddress(ctx context.Context, tx postgre.Transaction, status bool, address *model.Address) error
	UpdateDefaultShopAddress(ctx context.Context, tx postgre.Transaction, status bool, address *model.Address) error
	DeleteAddress(ctx context.Context, addressID string) error
}
