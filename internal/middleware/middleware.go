package middleware

import (
	"murakali/config"
	"murakali/pkg/logger"
)

type MWManager struct {
	cfg     *config.Config
	origins []string
	logger  logger.Logger
}

func NewMiddlewareManager(cfg *config.Config, origins []string, logger logger.Logger) *MWManager {
	return &MWManager{cfg: cfg, origins: origins, logger: logger}
}
