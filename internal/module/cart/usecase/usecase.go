package usecase

import (
	"context"
	"database/sql"
	"murakali/config"
	"murakali/internal/module/cart"
	"murakali/internal/module/cart/delivery/body"
	"murakali/pkg/httperror"
	"murakali/pkg/postgre"
	"murakali/pkg/response"
	"net/http"
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
		if err == sql.ErrNoRows {
			return nil, httperror.New(http.StatusBadRequest, response.AddressNotExistMessage)
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
