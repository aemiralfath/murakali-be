package repository

import (
	"context"
	"database/sql"
	"fmt"
	"murakali/internal/constant"
	"murakali/internal/model"
	"murakali/internal/module/user"
	"murakali/pkg/postgre"
	"time"

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

func (r *userRepo) GetUserByID(ctx context.Context, id string) (*model.User, error) {
	var userModel model.User
	if err := r.PSQL.QueryRowContext(ctx, GetUserByIDQuery, id).
		Scan(&userModel.ID, &userModel.RoleID, &userModel.Email, &userModel.Username, &userModel.PhoneNo,
			&userModel.FullName, &userModel.Gender, &userModel.BirthDate, &userModel.IsVerify); err != nil {
		return nil, err
	}

	return &userModel, nil
}

func (r *userRepo) CheckEmailHistory(ctx context.Context, email string) (*model.EmailHistory, error) {
	var emailHistory model.EmailHistory
	if err := r.PSQL.QueryRowContext(ctx, CheckEmailHistoryQuery, email).
		Scan(&emailHistory.ID, &emailHistory.Email); err != nil {
		return nil, err
	}

	return &emailHistory, nil
}
func (r *userRepo) GetUserByUsername(ctx context.Context, username string) (*model.User, error) {
	var userModel model.User
	if err := r.PSQL.QueryRowContext(ctx, GetUserByUsernameQuery, username).
		Scan(&userModel.ID, &userModel.Email, &userModel.Username, &userModel.IsVerify); err != nil {
		return nil, err
	}

	return &userModel, nil
}

func (r *userRepo) GetUserByPhoneNo(ctx context.Context, phoneNo string) (*model.User, error) {
	var userModel model.User
	if err := r.PSQL.QueryRowContext(ctx, GetUserByPhoneNoQuery, phoneNo).
		Scan(&userModel.ID, &userModel.Email, &userModel.PhoneNo, &userModel.IsVerify); err != nil {
		return nil, err
	}

	return &userModel, nil
}
func (r *userRepo) UpdateUserField(ctx context.Context, userModel *model.User) error {
	_, err := r.PSQL.ExecContext(
		ctx, UpdateUserFieldQuery, userModel.Username, userModel.FullName,
		userModel.PhoneNo, userModel.BirthDate, userModel.Gender, userModel.UpdatedAt, userModel.Email)
	if err != nil {
		return err
	}

	return nil
}

func (r *userRepo) UpdateUserEmail(ctx context.Context, tx postgre.Transaction, user *model.User) error {
	_, err := tx.ExecContext(
		ctx, UpdateUserEmailQuery, user.Email, user.UpdatedAt, user.ID)
	if err != nil {
		return err
	}

	return nil
}

func (r *userRepo) CreateEmailHistory(ctx context.Context, tx postgre.Transaction, email string) error {
	_, err := tx.ExecContext(ctx, CreateEmailHistoryQuery, email)
	if err != nil {
		return err
	}

	return nil
}

func (r *userRepo) InsertNewOTPKey(ctx context.Context, email, otp string) error {
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

func (r *userRepo) GetOTPValue(ctx context.Context, email string) (string, error) {
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

func (r *userRepo) DeleteOTPValue(ctx context.Context, email string) (int64, error) {
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
