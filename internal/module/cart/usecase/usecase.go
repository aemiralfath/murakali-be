package usecase

import (
	"murakali/config"
	"murakali/internal/module/cart"
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
