package seller

import (
	"context"
	"murakali/pkg/pagination"
)

type UseCase interface {
	GetOrder(ctx context.Context, userID string, pgn *pagination.Pagination) (*pagination.Pagination, error)
}
