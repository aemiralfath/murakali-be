package repository

import (
	"database/sql"
	"murakali/internal/module/product"

	"github.com/go-redis/redis/v8"
)

type productRepo struct {
	PSQL        *sql.DB
	RedisClient *redis.Client
}

func NewProductRepository(psql *sql.DB, client *redis.Client) product.Repository {
	return &productRepo{
		PSQL:        psql,
		RedisClient: client,
	}
}
