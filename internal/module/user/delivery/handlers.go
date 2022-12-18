package delivery

import (
	"murakali/config"
	"murakali/internal/module/user"
	"murakali/pkg/logger"
)

type userHandlers struct {
	cfg    *config.Config
	userUC user.UseCase
	logger logger.Logger
}

func NewUserHandlers(cfg *config.Config, userUC user.UseCase, log logger.Logger) user.Handlers {
	return &userHandlers{cfg: cfg, userUC: userUC, logger: log}
}
