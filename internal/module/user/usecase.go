package user

import (
	"context"
	"murakali/pkg/pagination"
)

type UseCase interface {
	GetAddress(ctx context.Context, userID, name string, pagination *pagination.Pagination) (*pagination.Pagination, error)
}
