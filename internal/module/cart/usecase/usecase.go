package usecase

import (
	"context"
	"database/sql"
	"math"
	"murakali/config"
	"murakali/internal/module/cart"
	"murakali/internal/module/cart/delivery/body"
	"murakali/pkg/httperror"
	"murakali/pkg/pagination"
	"murakali/pkg/postgre"
	"murakali/pkg/response"
	"net/http"
	"time"
)

type cartUC struct {
	cfg      *config.Config
	txRepo   *postgre.TxRepo
	cartRepo cart.Repository
}

func NewCartUseCase(cfg *config.Config, txRepo *postgre.TxRepo, cartRepo cart.Repository) cart.UseCase {
	return &cartUC{cfg: cfg, txRepo: txRepo, cartRepo: cartRepo}
}

func (u *cartUC) GetCartHoverHome(ctx context.Context, userID string, limit int) (*body.CartHomeResponse, error) {
	cartHomes, err := u.cartRepo.GetCartHoverHome(ctx, userID, limit)
	if err != nil {
		if err != sql.ErrNoRows {
			return nil, err
		}
		return nil, err
	}

	for _, cart := range cartHomes {
		var maxDiscountPrice float64
		var minProductPrice float64
		var discountPercentage float64
		var discountFixPrice float64
		var resultDiscount float64
		if cart.MaxDiscountPrice == nil {
			continue
		}
		if cart.MinProductPrice != nil {
			minProductPrice = *cart.MinProductPrice
		}

		maxDiscountPrice = *cart.MaxDiscountPrice
		if cart.DiscountPercentage != nil {
			discountPercentage = *cart.DiscountPercentage
			if cart.Price >= minProductPrice && discountPercentage > 0 {
				resultDiscount = math.Min(maxDiscountPrice,
					cart.Price*(discountPercentage/100.00))
			}
		}

		if cart.DiscountFixPrice != nil {
			discountFixPrice = *cart.DiscountFixPrice
			if cart.Price >= minProductPrice && discountFixPrice > 0 {
				resultDiscount = math.Max(resultDiscount, discountFixPrice)
				resultDiscount = math.Min(resultDiscount, maxDiscountPrice)
			}
		}

		if resultDiscount > 0 {
			cart.ResultDiscount = resultDiscount
			cart.SubPrice = cart.Price - resultDiscount
		}
	}

	totalItem, err := u.cartRepo.GetTotalCart(ctx, userID)
	if err != nil {
		return nil, err
	}

	cartHomeResponses := &body.CartHomeResponse{
		Limit:     limit,
		TotalItem: totalItem,
		CartHomes: cartHomes,
	}

	return cartHomeResponses, nil
}

func (u *cartUC) GetCartItems(ctx context.Context, userID string, pgn *pagination.Pagination) (*pagination.Pagination, error) {
	totalRows, err := u.cartRepo.GetTotalCart(ctx, userID)
	if err != nil {
		return nil, err
	}
	totalPages := int(math.Ceil(float64(totalRows) / float64(pgn.Limit)))
	pgn.TotalRows = totalRows
	pgn.TotalPages = totalPages
	pgn.Limit = int(totalRows)

	carts, products, promos, err := u.cartRepo.GetCartItems(ctx, userID, pgn)
	if err != nil {
		if err != sql.ErrNoRows {
			return nil, err
		}
	}

	flagCart := make(map[string]int)
	CartResults := make([]*body.CartItemsResponse, 0)
	n := 0
	totalCarts := len(carts)
	for i := 0; i < totalCarts; i++ {
		if _, exists := flagCart[carts[i].Shop.ID.String()]; !exists {
			flagCart[carts[i].Shop.ID.String()] = n
			CartResults = append(CartResults, carts[i])
			n++
		}
		idx := flagCart[carts[i].Shop.ID.String()]
		p := products[i]
		promo := promos[i]

		p.Promo = promo
		p = u.CalculateDiscountProduct(p)
		CartResults[idx].ProductDetails = append(CartResults[idx].ProductDetails, p)
	}
	pgn.Rows = CartResults

	return pgn, nil
}

func (u *cartUC) AddCartItems(ctx context.Context, userID string, requestBody body.AddCartItemRequest) error {
	userModel, err := u.cartRepo.GetUserByID(ctx, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return httperror.New(http.StatusBadRequest, response.UserNotExistMessage)
		}
		return err
	}

	productDetail, err := u.cartRepo.GetProductDetailByID(ctx, requestBody.ProductDetailID)
	if err != nil {
		if err == sql.ErrNoRows {
			return httperror.New(http.StatusBadRequest, response.ProductDetailNotExistMessage)
		}
		return err
	}

	cartProductDetail, err := u.cartRepo.GetCartProductDetail(ctx, userModel.ID.String(), productDetail.ID.String())
	if err != nil {
		if err != sql.ErrNoRows {
			return err
		}
		if requestBody.Quantity > productDetail.Stock {
			return httperror.New(http.StatusBadRequest, response.QuantityReachedMaximum)
		}
		_, err = u.cartRepo.CreateCart(ctx, userModel.ID.String(), productDetail.ID.String(), requestBody.Quantity)
		if err != nil {
			return err
		}
		return nil
	}

	cartProductDetail.Quantity += requestBody.Quantity
	if cartProductDetail.Quantity > productDetail.Stock {
		return httperror.New(http.StatusBadRequest, response.QuantityReachedMaximum)
	}

	cartProductDetail.UpdatedAt.Time = time.Now()
	cartProductDetail.UpdatedAt.Valid = true
	err = u.cartRepo.UpdateCartByID(ctx, cartProductDetail)
	if err != nil {
		return err
	}

	return nil
}

func (u *cartUC) UpdateCartItems(ctx context.Context, userID string, requestBody body.CartItemRequest) error {
	userModel, err := u.cartRepo.GetUserByID(ctx, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return httperror.New(http.StatusBadRequest, response.UserNotExistMessage)
		}
		return err
	}

	productDetail, err := u.cartRepo.GetProductDetailByID(ctx, requestBody.ProductDetailID)
	if err != nil {
		if err == sql.ErrNoRows {
			return httperror.New(http.StatusBadRequest, response.ProductDetailNotExistMessage)
		}
		return err
	}

	cartProductDetail, err := u.cartRepo.GetCartProductDetail(ctx, userModel.ID.String(), productDetail.ID.String())
	if err != nil {
		return err
	}

	cartProductDetail.Quantity = requestBody.Quantity
	if cartProductDetail.Quantity > productDetail.Stock {
		return httperror.New(http.StatusBadRequest, response.QuantityReachedMaximum)
	}

	cartProductDetail.UpdatedAt.Time = time.Now()
	cartProductDetail.UpdatedAt.Valid = true
	err = u.cartRepo.UpdateCartByID(ctx, cartProductDetail)
	if err != nil {
		return err
	}

	return nil
}

func (u *cartUC) DeleteCartItems(ctx context.Context, userID, productDetailID string) error {
	userModel, err := u.cartRepo.GetUserByID(ctx, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return httperror.New(http.StatusBadRequest, response.UserNotExistMessage)
		}
		return err
	}

	productDetail, err := u.cartRepo.GetProductDetailByID(ctx, productDetailID)
	if err != nil {
		if err == sql.ErrNoRows {
			return httperror.New(http.StatusBadRequest, response.ProductDetailNotExistMessage)
		}
		return err
	}

	cartProductDetail, err := u.cartRepo.GetCartProductDetail(ctx, userModel.ID.String(), productDetail.ID.String())
	if err != nil {
		return err
	}

	err = u.cartRepo.DeleteCartByID(ctx, cartProductDetail)
	if err != nil {
		return err
	}

	return nil
}

func (u *cartUC) CalculateDiscountProduct(p *body.ProductDetailResponse) *body.ProductDetailResponse {
	if p.Promo.MaxDiscountPrice == nil {
		return p
	}
	var maxDiscountPrice float64
	var minProductPrice float64
	var discountPercentage float64
	var discountFixPrice float64
	var resultDiscount float64
	if p.Promo.MinProductPrice != nil {
		minProductPrice = *p.Promo.MinProductPrice
	}

	maxDiscountPrice = *p.Promo.MaxDiscountPrice
	if p.Promo.DiscountPercentage != nil {
		discountPercentage = *p.Promo.DiscountPercentage
		if p.ProductPrice >= minProductPrice && *p.Promo.DiscountPercentage > 0 {
			resultDiscount = math.Min(maxDiscountPrice,
				p.ProductPrice*(discountPercentage/100.00))
		}
	}

	if p.Promo.DiscountFixPrice != nil {
		discountFixPrice = *p.Promo.DiscountFixPrice
		if p.ProductPrice >= minProductPrice && discountFixPrice > 0 {
			resultDiscount = math.Max(resultDiscount, discountFixPrice)
			resultDiscount = math.Min(resultDiscount, maxDiscountPrice)
		}
	}

	if resultDiscount > 0 {
		p.Promo.ResultDiscount = resultDiscount
		p.Promo.SubPrice = p.ProductPrice - resultDiscount
	}

	return p
}
