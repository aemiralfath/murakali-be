package cart

import (
	"context"
	"murakali/internal/module/cart/delivery/body"
	"murakali/pkg/pagination"
)

type UseCase interface {
	GetCartHoverHome(ctx context.Context, userID string, limit int) (*body.CartHomeResponse, error)
	GetCartItems(ctx context.Context, userID string, pgn *pagination.Pagination) (*pagination.Pagination, error)
	AddCartItems(ctx context.Context, userID string, requestBody body.AddCartItemRequest) error
	UpdateCartItems(ctx context.Context, userID string, requestBody body.CartItemRequest) error
}
