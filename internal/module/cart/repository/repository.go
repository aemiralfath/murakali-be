package repository

import (
	"database/sql"
	"murakali/internal/module/cart"

	"github.com/go-redis/redis/v8"
)

type cartRepo struct {
	PSQL        *sql.DB
	RedisClient *redis.Client
}

func NewCartRepository(psql *sql.DB, client *redis.Client) cart.Repository {
	return &cartRepo{
		PSQL:        psql,
		RedisClient: client,
	}
}
