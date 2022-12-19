package repository

import (
	"context"
	"database/sql"
	"murakali/internal/model"
	"murakali/internal/module/user"

	"github.com/go-redis/redis/v8"
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

func (r *userRepo) GetSealabsPay(ctx context.Context, userid string) ([]*model.SealabsPay, error) {

	responses := make([]*model.SealabsPay, 0)
	res, err := r.PSQL.QueryContext(ctx, GetSealabsPayByIdQuery, userid)

	if err != nil {
		return nil, err
	}
	defer res.Close()

	for res.Next() {
		var response model.SealabsPay
		if err := res.Scan(
			&response.CardNumber, &response.UserID, &response.Name, &response.IsDefault, &response.ActiveDate, &response.CreatedAt, &response.UpdatedAt, &response.DeletedAt); err != nil {
			return nil, err
		}
		responses = append(responses, &response)
	}
	return responses, nil
}
