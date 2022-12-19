package repository

import (
	"context"
	"database/sql"
	"murakali/internal/model"
	"murakali/internal/module/user"
	"murakali/internal/module/user/delivery/body"
	"murakali/pkg/postgre"

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

func (r *userRepo) CheckDefaultSealabsPay(ctx context.Context, userid string) (*string, error) {
	var temp *string
	if err := r.PSQL.QueryRowContext(ctx, CheckDefaultSealabsPayQuery, userid).
		Scan(&temp); err != nil {
		return nil, err
	}

	return temp, nil
}

func (r *userRepo) SetDefaultSealabsPayTrans(ctx context.Context, tx postgre.Transaction, card_number *string) error {
	if _, err := tx.ExecContext(ctx, SetDefaultSealabsPayTransQuery, card_number); err != nil {
		return err
	}

	return nil
}

func (r *userRepo) AddSealabsPay(ctx context.Context, tx postgre.Transaction, request body.AddSealabsPayRequest) error {

	if _, err := tx.ExecContext(ctx, CreateSealabsPayQuery, request.CardNumber, request.UserID, request.Name, request.IsDefault, request.ActiveDateTime); err != nil {
		return err
	}

	return nil
}

func (r *userRepo) SetDefaultSealabsPay(ctx context.Context, card_number string, userid string) error {
	if _, err := r.PSQL.ExecContext(ctx, SetDefaultSealabsPayQuery, card_number, userid); err != nil {
		return err
	}

	return nil
}

func (r *userRepo) PatchSealabsPay(ctx context.Context, card_number string) error {

	if _, err := r.PSQL.ExecContext(ctx, PatchSealabsPayQuery, card_number); err != nil {
		return err
	}
	return nil
}
