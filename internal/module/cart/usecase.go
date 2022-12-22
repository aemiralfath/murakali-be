package cart

import (
	"context"
	"murakali/internal/module/cart/delivery/body"
)

type UseCase interface {
	GetCartHoverHome(ctx context.Context, userID string, limit int) (*body.CartHomeResponse, error)
}
