package redis

import (
	"context"
	"murakali/config"

	"github.com/go-redis/redis/v8"
)

func NewRedis(config *config.Config) (*redis.Client, error) {
	connectionCfg := &redis.Options{
		Addr:     config.Redis.Address,
		Password: config.Redis.Password,
		DB:       config.Redis.DB,
	}

	client := redis.NewClient(connectionCfg)
	if err := client.Ping(context.Background()).Err(); err != nil {
		return nil, err
	}

	return client, nil
}
