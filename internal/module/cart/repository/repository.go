package repository

import (
	"context"
	"database/sql"
	"math"
	"murakali/internal/model"
	"murakali/internal/module/cart"
	"murakali/internal/module/cart/delivery/body"
	"murakali/pkg/pagination"

	"github.com/go-redis/redis/v8"
	"github.com/lib/pq"
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

func (r *cartRepo) GetUserByID(ctx context.Context, id string) (*model.User, error) {
	var userModel model.User
	if err := r.PSQL.QueryRowContext(ctx, GetUserByIDQuery, id).
		Scan(&userModel.ID, &userModel.RoleID, &userModel.Email, &userModel.Username, &userModel.PhoneNo,
			&userModel.FullName, &userModel.Gender, &userModel.BirthDate, &userModel.IsVerify, &userModel.PhotoURL); err != nil {
		return nil, err
	}

	return &userModel, nil
}

func (r *cartRepo) GetProductDetailByID(ctx context.Context, productDetailID string) (*model.ProductDetail, error) {
	var ProductDetailData model.ProductDetail
	if err := r.PSQL.QueryRowContext(ctx, GetProductDetailByIDQuery, productDetailID).
		Scan(&ProductDetailData.ID, &ProductDetailData.ProductID, &ProductDetailData.Price,
			&ProductDetailData.Stock, &ProductDetailData.Weight, &ProductDetailData.Size); err != nil {
		return nil, err
	}

	return &ProductDetailData, nil
}

func (r *cartRepo) GetCartProductDetail(ctx context.Context, userID, productDetailID string) (*model.CartItem, error) {
	var cartModel model.CartItem
	if err := r.PSQL.QueryRowContext(ctx, GetCartProductDetailQuery, userID, productDetailID).
		Scan(&cartModel.ID, &cartModel.UserID, &cartModel.ProductDetailID, &cartModel.Quantity); err != nil {
		return nil, err
	}

	return &cartModel, nil
}

func (r *cartRepo) CreateCart(ctx context.Context, userID, productDetailID string, quantity float64) (*model.CartItem, error) {
	var cartModel model.CartItem
	if err := r.PSQL.QueryRowContext(ctx, CreateCartQuery, userID, productDetailID, quantity).
		Scan(&cartModel.ID); err != nil {
		return nil, err
	}

	return &cartModel, nil
}

func (r *cartRepo) UpdateCartByID(ctx context.Context, cartItem *model.CartItem) error {
	_, err := r.PSQL.ExecContext(ctx, UpdateCartByIDQuery, cartItem.Quantity, cartItem.UpdatedAt, cartItem.ID.String())
	if err != nil {
		return err
	}

	return nil
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
		var cartItem body.CartHome
		if errScan := res.Scan(
			&cartItem.Title,
			&cartItem.ThumbnailURL,
			&cartItem.Price,
			&cartItem.DiscountPercentage,
			&cartItem.DiscountFixPrice,
			&cartItem.MinProductPrice,
			&cartItem.MaxDiscountPrice,
			&cartItem.Quantity,
			&cartItem.VariantName,
			&cartItem.VariantType,
		); errScan != nil {
			return nil, err
		}

		if cartItem.Price >= cartItem.MinProductPrice && cartItem.DiscountPercentage > 0 {
			cartItem.ResultDiscount = math.Min(cartItem.MaxDiscountPrice,
				cartItem.Price*(cartItem.DiscountPercentage/100))
		}

		if cartItem.Price >= cartItem.MinProductPrice && cartItem.DiscountFixPrice > 0 {
			cartItem.ResultDiscount = math.Max(cartItem.ResultDiscount, cartItem.DiscountFixPrice)

			cartItem.ResultDiscount = math.Min(cartItem.MaxDiscountPrice,
				cartItem.Price*(cartItem.DiscountPercentage/100))
		}

		if cartItem.ResultDiscount > 0 {
			cartItem.SubPrice = cartItem.Price - cartItem.ResultDiscount
		}

		cartHomes = append(cartHomes, &cartItem)
	}

	if res.Err() != nil {
		return nil, err
	}

	return cartHomes, nil
}

func (r *cartRepo) GetCartItems(ctx context.Context, userID string, pgn *pagination.Pagination) ([]*body.CartItemsResponse,
	[]*body.ProductDetailResponse, []*body.PromoResponse, error) {
	cartItems := make([]*body.CartItemsResponse, 0)
	products := make([]*body.ProductDetailResponse, 0)
	promos := make([]*body.PromoResponse, 0)
	res, err := r.PSQL.QueryContext(
		ctx, GetCartItemsQuery,
		userID,
		pgn.GetLimit(),
		pgn.GetOffset())

	if err != nil {
		return nil, nil, nil, err
	}
	defer res.Close()

	for res.Next() {
		var cartItem body.CartItemsResponse
		var productData body.ProductDetailResponse
		var promo body.PromoResponse
		var shop body.ShopResponse
		var VariantName []string
		var VariantType []string
		if errScan := res.Scan(
			&cartItem.ID,
			&productData.Quantity,
			&productData.ID,
			&productData.Title,
			&shop.ID,
			&shop.Name,
			&productData.ThumbnailURL,
			&productData.ProductPrice,
			&productData.ProductStock,
			&promo.DiscountPercentage,
			&promo.DiscountFixPrice,
			&promo.MinProductPrice,
			&promo.MaxDiscountPrice,
			(*pq.StringArray)(&VariantName),
			(*pq.StringArray)(&VariantType),
		); errScan != nil {
			return nil, nil, nil, err
		}
		mapVariant := make(map[string]string, 0)
		n := len(VariantName)
		for i := 0; i < n; i++ {
			mapVariant[VariantName[i]] = VariantType[i]
		}

		productData.Variant = mapVariant
		cartItem.Shop = &shop

		cartItems = append(cartItems, &cartItem)
		products = append(products, &productData)
		promos = append(promos, &promo)
	}

	if res.Err() != nil {
		return nil, nil, nil, err
	}

	return cartItems, products, promos, err
}
