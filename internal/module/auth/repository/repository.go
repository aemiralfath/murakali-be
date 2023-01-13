package repository

import (
	"context"
	"database/sql"
	"fmt"
	"murakali/internal/constant"
	"murakali/internal/model"
	"murakali/internal/module/auth"
	"murakali/pkg/postgre"
	"time"

	"github.com/go-redis/redis/v8"
)

type authRepo struct {
	PSQL        *sql.DB
	RedisClient *redis.Client
}

func NewAuthRepository(psql *sql.DB, client *redis.Client) auth.Repository {
	return &authRepo{
		PSQL:        psql,
		RedisClient: client,
	}
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
		Scan(&user.ID, &user.RoleID, &user.Email, &user.Password, &user.Username, &user.IsVerify, &user.IsSSO); err != nil {
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

func (r *authRepo) UpdatePassword(ctx context.Context, user *model.User, password string) (*model.User, error) {
	_, err := r.PSQL.ExecContext(ctx, UpdatePasswordQuery, password, user.Email)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *authRepo) CreateEmailHistory(ctx context.Context, tx postgre.Transaction, email string) error {
	_, err := tx.ExecContext(ctx, CreateEmailHistoryQuery, email)
	if err != nil {
		return err
	}

	return nil
}

func (r *authRepo) UpdateUser(ctx context.Context, tx postgre.Transaction, user *model.User) error {
	_, err := tx.ExecContext(
		ctx, VerifyUserQuery, user.PhoneNo, user.FullName,
		user.Username, user.Password, user.IsVerify, user.UpdatedAt, user.Email)
	if err != nil {
		return err
	}

	return nil
}

func (r *authRepo) InsertNewOTPKey(ctx context.Context, email, otp string) error {
	key := fmt.Sprintf("%s:%s", constant.OtpKey, email)

	duration, err := time.ParseDuration(constant.OtpDuration)
	if err != nil {
		return err
	}

	if err := r.RedisClient.Set(ctx, key, otp, duration); err.Err() != nil {
		return err.Err()
	}

	return nil
}

func (r *authRepo) GetOTPValue(ctx context.Context, email string) (string, error) {
	key := fmt.Sprintf("%s:%s", constant.OtpKey, email)

	res := r.RedisClient.Get(ctx, key)
	if res.Err() != nil {
		return "", res.Err()
	}

	value, err := res.Result()
	if err != nil {
		return "", err
	}

	return value, nil
}

func (r *authRepo) DeleteOTPValue(ctx context.Context, email string) (int64, error) {
	key := fmt.Sprintf("%s:%s", constant.OtpKey, email)

	res := r.RedisClient.Del(ctx, key)
	if res.Err() != nil {
		return -1, res.Err()
	}

	value, err := res.Result()
	if err != nil {
		return -1, err
	}

	return value, nil
}

func (r *authRepo) CreateUserGoogle(ctx context.Context, tx postgre.Transaction, user *model.User) (*model.User, error) {
	if err := tx.QueryRowContext(ctx, CreateUserGoogleQuery, constant.RoleUser, user.Username, user.Email, user.FullName, user.PhotoURL, user.IsSSO, user.IsVerify).Scan(&user.ID, &user.RoleID); err != nil {
		return nil, err
	}

	return user, nil
}
