package auth

import (
	"context"
	"murakali/internal/model"
	"murakali/pkg/postgre"
)

type Repository interface {
	CheckEmailHistory(ctx context.Context, email string) (*model.EmailHistory, error)
	GetUserByEmail(ctx context.Context, email string) (*model.User, error)
	CreateUser(ctx context.Context, email string) (*model.User, error)
	InsertNewOTPKey(ctx context.Context, email, otp string) error
	GetOTPValue(ctx context.Context, email string) (string, error)
	CreateEmailHistory(ctx context.Context, tx postgre.Transaction, email string) error
	VerifyUser(ctx context.Context, tx postgre.Transaction, email string) error
}
