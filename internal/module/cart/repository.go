package cart

import (
	"context"
	"murakali/internal/model"
	"murakali/internal/module/cart/delivery/body"
	"murakali/pkg/pagination"
)

type Repository interface {
	GetUserByID(ctx context.Context, id string) (*model.User, error)
	GetProductDetailByID(ctx context.Context, productDetailID string) (*model.ProductDetail, error)
	GetCartProductDetail(ctx context.Context, userID, productDetailID string) (*model.CartItem, error)
	CreateCart(ctx context.Context, userID, productDetailID string, quantity float64) (*model.CartItem, error)
	UpdateCartByID(ctx context.Context, cartItem *model.CartItem) error
	GetTotalCart(ctx context.Context, userID string) (int64, error)
	GetCartHoverHome(ctx context.Context, userID string, limit int) ([]*body.CartHome, error)
	GetCartItems(ctx context.Context, userID string, pgn *pagination.Pagination) ([]*body.CartItemsResponse,
		[]*body.ProductDetailResponse, []*body.PromoResponse, error)
}
