package usecase

import (
	"context"
	"database/sql"
	"math"
	"murakali/config"
	"murakali/internal/module/cart"
	"murakali/internal/module/cart/delivery/body"
	"murakali/pkg/pagination"
	"murakali/pkg/postgre"
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
		CartResults[idx].Product = append(CartResults[idx].Product, p)
	}
	pgn.Rows = CartResults

	return pgn, nil
}

func (u *cartUC) CalculateDiscountProduct(p *body.ProductResponse) *body.ProductResponse {
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
