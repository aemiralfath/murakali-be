package delivery

import (
	"murakali/config"
	"murakali/internal/module/product"
	"murakali/pkg/logger"
)

type prouctHandlers struct {
	cfg       *config.Config
	productUC product.UseCase
	logger    logger.Logger
}

func NewProductHandlers(cfg *config.Config, productUC product.UseCase, log logger.Logger) product.Handlers {
	return &prouctHandlers{cfg: cfg, productUC: productUC, logger: log}
}
