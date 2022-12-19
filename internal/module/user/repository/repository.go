package repository

import (
	"database/sql"
	"github.com/go-redis/redis/v8"
	"murakali/internal/module/user"
)

type userRepo struct {
	PSQL        *sql.DB
	RedisClient *redis.Client
}

func NewUserRepository(psql *sql.DB, client *redis.Client) user.Repository {
	return &userRepo{
		PSQL:        psql,
		RedisClient: client,
	}
}
