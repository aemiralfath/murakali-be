package repository

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/go-redis/redis/v8"
	"murakali/internal/auth"
	"murakali/internal/constant"
	"murakali/internal/model"
	"murakali/pkg/postgre"
	"time"
)

type authRepo struct {
	PSQL        *sql.DB
	RedisClient *redis.Client
}

func NewAuthRepository(psql *sql.DB, client *redis.Client) (auth.Repository, error) {
	return &authRepo{
		PSQL:        psql,
		RedisClient: client,
	}, nil
}

func (r *authRepo) CheckEmailHistory(ctx context.Context, email string) (*model.EmailHistory, error) {
	var emailHistory model.EmailHistory
	if err := r.PSQL.QueryRowContext(ctx, CheckEmailHistoryQuery, email).
		Scan(&emailHistory.ID, &emailHistory.Email); err != nil {
		return nil, err
	}

	return &emailHistory, nil
}

func (r *authRepo) GetUserByID(ctx context.Context, id string) (*model.User, error) {
	var user model.User
	if err := r.PSQL.QueryRowContext(ctx, GetUserByIDQuery, id).
		Scan(&user.ID, &user.RoleID, &user.Email); err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *authRepo) GetUserByEmail(ctx context.Context, email string) (*model.User, error) {
	var user model.User
	if err := r.PSQL.QueryRowContext(ctx, GetUserByEmailQuery, email).
		Scan(&user.ID, &user.RoleID, &user.Email, &user.Password, &user.IsVerify); err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *authRepo) GetUserByUsername(ctx context.Context, username string) (*model.User, error) {
	var user model.User
	if err := r.PSQL.QueryRowContext(ctx, GetUserByUsernameQuery, username).
		Scan(&user.ID, &user.Email, &user.IsVerify); err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *authRepo) GetUserByPhoneNo(ctx context.Context, phoneNo string) (*model.User, error) {
	var user model.User
	if err := r.PSQL.QueryRowContext(ctx, GetUserByPhoneNoQuery, phoneNo).
		Scan(&user.ID, &user.Email, &user.IsVerify); err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *authRepo) CreateUser(ctx context.Context, email string) (*model.User, error) {
	var user model.User
	if err := r.PSQL.QueryRowContext(ctx, CreateUserQuery, constant.RoleUser, email, false, false).
		Scan(&user.ID, &user.Email); err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *authRepo) CreateEmailHistory(ctx context.Context, tx postgre.Transaction, email string) error {
	_, err := tx.ExecContext(ctx, CreateEmailHistoryQuery, email)
	if err != nil {
		return err
	}

	return nil
}

func (r *authRepo) UpdateUser(ctx context.Context, tx postgre.Transaction, user *model.User) error {
	_, err := tx.ExecContext(ctx, VerifyUserQuery, user.PhoneNo, user.FullName, user.Username, user.Password, user.IsVerify, user.UpdatedAt, user.Email)
	if err != nil {
		return err
	}

	return nil
}

func (r *authRepo) InsertNewOTPKey(ctx context.Context, email, otp string) error {
	key := fmt.Sprintf("%s:%s", constant.OtpKey, email)
	value := fmt.Sprintf("%s", otp)

	duration, err := time.ParseDuration(constant.OtpDuration)
	if err != nil {
		return err
	}

	if err := r.RedisClient.Set(ctx, key, value, duration); err.Err() != nil {
		return err.Err()
	}

	return nil
}

func (r *authRepo) GetOTPValue(ctx context.Context, email string) (string, error) {
	key := fmt.Sprintf("%s:%s", constant.OtpKey, email)

	res := r.RedisClient.GetDel(ctx, key)
	if res.Err() != nil {
		return "", res.Err()
	}

	value, err := res.Result()
	if err != nil {
		return "", err
	}

	return value, nil
}
