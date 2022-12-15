package repository

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/go-redis/redis/v8"
	"murakali/internal/auth"
	"murakali/internal/constant"
	"murakali/internal/model"
	"time"
)

type statement struct {
	CheckEMailHistory *sql.Stmt
	GetUserByEmail    *sql.Stmt
	CreateUser        *sql.Stmt
}

type authRepo struct {
	PSQL        *sql.DB // postgres
	RedisClient *redis.Client
	statement   statement
}

func initStatement(postgres *sql.DB) (statement, error) {
	var dbStatement statement
	var err error

	dbStatement.CheckEMailHistory, err = postgres.Prepare(CheckEmailHistoryQuery)
	if err != nil {
		return dbStatement, err
	}

	dbStatement.CreateUser, err = postgres.Prepare(CreateUserQuery)
	if err != nil {
		return dbStatement, err
	}

	dbStatement.GetUserByEmail, err = postgres.Prepare(GetUserByEmailQuery)
	if err != nil {
		return dbStatement, err
	}

	return dbStatement, nil
}

func NewAuthRepository(psql *sql.DB, client *redis.Client) (auth.Repository, error) {
	dbStatement, err := initStatement(psql)
	if err != nil {
		return nil, err
	}

	return &authRepo{
		PSQL:        psql,
		RedisClient: client,
		statement:   dbStatement,
	}, nil
}

func (r *authRepo) CheckEmailHistory(ctx context.Context, email string) (*model.EmailHistory, error) {
	var emailHistory model.EmailHistory
	if err := r.statement.CheckEMailHistory.QueryRowContext(ctx, email).
		Scan(&emailHistory.ID, &emailHistory.Email); err != nil {
		return nil, err
	}

	return &emailHistory, nil
}

func (r *authRepo) GetUserByEmail(ctx context.Context, email string) (*model.User, error) {
	var user model.User
	if err := r.statement.GetUserByEmail.QueryRowContext(ctx, email).
		Scan(&user.ID, &user.Email, &user.IsVerify); err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *authRepo) CreateUser(ctx context.Context, email string) (*model.User, error) {
	var user model.User
	if err := r.statement.CreateUser.QueryRowContext(ctx, constant.RoleUser, email, false, false).
		Scan(&user.ID, &user.Email); err != nil {
		return nil, err
	}

	return &user, nil
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
