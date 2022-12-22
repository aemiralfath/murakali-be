package repository

import (
	"context"
	"database/sql"
	"math"
	"murakali/internal/module/cart"
	"murakali/internal/module/cart/delivery/body"

	"github.com/go-redis/redis/v8"
)

type cartRepo struct {
	PSQL        *sql.DB
	RedisClient *redis.Client
}

func NewCartRepository(psql *sql.DB, client *redis.Client) cart.Repository {
	return &cartRepo{
		PSQL:        psql,
		RedisClient: client,
	}
}

func (r *cartRepo) GetTotalCart(ctx context.Context, userID string) (int64, error) {
	var total int64
	if err := r.PSQL.QueryRowContext(ctx, GetTotalCartQuery, userID).Scan(&total); err != nil {
		return 0, err
	}

	return total, nil
}

func (r *cartRepo) GetCartHoverHome(ctx context.Context, userID string, limit int) ([]*body.CartHome, error) {
	cartHomes := make([]*body.CartHome, 0)
	res, err := r.PSQL.QueryContext(
		ctx, GetCartHoverHomeQuery,
		userID, limit)

	if err != nil {
		return nil, err
	}
	defer res.Close()

	for res.Next() {
		var CartHome body.CartHome
		if errScan := res.Scan(
			&CartHome.Title,
			&CartHome.ThumbnailURL,
			&CartHome.Price,
			&CartHome.DiscountPersentase,
			&CartHome.DiscountFixPrice,
			&CartHome.MinProductPrice,
			&CartHome.MaxDiscountPrice,
			&CartHome.Quantity,
			&CartHome.VariantName,
			&CartHome.VariantType,
		); errScan != nil {
			return nil, err
		}

		if CartHome.Price >= CartHome.MinProductPrice && CartHome.DiscountPersentase > 0 {
			CartHome.ResultDiscount = math.Min(CartHome.MaxDiscountPrice,
				CartHome.Price*(CartHome.DiscountPersentase/100))
		}

		if CartHome.Price >= CartHome.MinProductPrice && CartHome.DiscountFixPrice > 0 {
			CartHome.ResultDiscount = math.Max(CartHome.ResultDiscount, CartHome.DiscountFixPrice)

			CartHome.ResultDiscount = math.Min(CartHome.MaxDiscountPrice,
				CartHome.Price*(CartHome.DiscountPersentase/100))
		}

		if CartHome.ResultDiscount > 0 {
			CartHome.SubPrice = CartHome.Price - CartHome.ResultDiscount
		}

		cartHomes = append(cartHomes, &CartHome)
	}

	if res.Err() != nil {
		return nil, err
	}

	return cartHomes, nil
}
