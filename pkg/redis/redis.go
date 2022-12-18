package redis

import (
	"context"
	"murakali/config"

	"github.com/go-redis/redis/v8"
)

func NewRedis(cfg *config.Config) (*redis.Client, error) {
	connectionCfg := &redis.Options{
		Addr:     cfg.Redis.Address,
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	}

	client := redis.NewClient(connectionCfg)
	if err := client.Ping(context.Background()).Err(); err != nil {
		return nil, err
	}

	return client, nil
}
