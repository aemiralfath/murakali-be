package usecase

import (
	"context"
	"database/sql"
	"murakali/config"
	"murakali/internal/model"
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

func (u *productUC) GetCategories(ctx context.Context) ([]*model.Category, error) {

	categories, err := u.productRepo.GetCategories(ctx)
	if err != nil {
		if err != sql.ErrNoRows {
			return nil, err
		}
	}

	return categories, nil
}
