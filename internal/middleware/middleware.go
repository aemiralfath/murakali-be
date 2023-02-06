package middleware

import (
	"github.com/go-redis/redis/v8"
	"murakali/config"
	"murakali/pkg/logger"
)

type MWManager struct {
	cfg         *config.Config
	origins     []string
	log         logger.Logger
	RedisClient *redis.Client
}

func NewMiddlewareManager(cfg *config.Config, origins []string, log logger.Logger, redisClient *redis.Client) *MWManager {
	return &MWManager{cfg: cfg, origins: origins, log: log, RedisClient: redisClient}
}
