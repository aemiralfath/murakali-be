package cart

import (
	"context"
	"murakali/internal/module/cart/delivery/body"
	"murakali/pkg/pagination"
)

type Repository interface {
	GetTotalCart(ctx context.Context, userID string) (int64, error)
	GetCartHoverHome(ctx context.Context, userID string, limit int) ([]*body.CartHome, error)
	GetCartItems(ctx context.Context, userID string, pgn *pagination.Pagination) ([]*body.CartItemsResponse,
		[]*body.ProductDetailResponse, []*body.PromoResponse, error)
}
