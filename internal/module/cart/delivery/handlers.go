package delivery

import (
	"murakali/config"
	"murakali/internal/module/cart"
	"murakali/pkg/logger"
)

type cartHandlers struct {
	cfg    *config.Config
	cartUC cart.UseCase
	logger logger.Logger
}

func NewCartHandlers(cfg *config.Config, cartUC cart.UseCase, log logger.Logger) cart.Handlers {
	return &cartHandlers{cfg: cfg, cartUC: cartUC, logger: log}
}
