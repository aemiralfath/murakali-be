package user

import (
	"context"
	"murakali/internal/model"
	"murakali/pkg/pagination"
)

type Repository interface {
	GetUserByID(ctx context.Context, id string) (*model.User, error)
	GetTotalAddress(ctx context.Context, userID, name string) (int64, error)
	GetAddresses(ctx context.Context, userID, name string, pagination *pagination.Pagination) ([]*model.Address, error)
}
