package middleware

import (
	"murakali/config"
	"murakali/pkg/logger"
)

type MWManager struct {
	cfg     *config.Config
	origins []string
	log     logger.Logger
}

func NewMiddlewareManager(cfg *config.Config, origins []string, log logger.Logger) *MWManager {
	return &MWManager{cfg: cfg, origins: origins, log: log}
}
