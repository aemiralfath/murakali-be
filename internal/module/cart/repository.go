package cart

import (
	"context"
	"murakali/internal/module/cart/delivery/body"
)

type Repository interface {
	GetTotalCart(ctx context.Context, userID string) (int64, error)
	GetCartHoverHome(ctx context.Context, userID string, limit int) ([]*body.CartHome, error)
}
