package usecase

import (
	"murakali/config"
	"murakali/internal/module/product"
	"murakali/pkg/postgre"
)

type productUC struct {
	cfg         *config.Config
	txRepo      *postgre.TxRepo
	productRepo product.Repository
}

func NewProductUseCase(cfg *config.Config, txRepo *postgre.TxRepo, productRepo product.Repository) product.UseCase {
	return &productUC{cfg: cfg, txRepo: txRepo, productRepo: productRepo}
}
