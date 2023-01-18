package auth

import (
	"context"
	"murakali/internal/model"
	"murakali/pkg/postgre"
)

type Repository interface {
	CheckEmailHistory(ctx context.Context, email string) (*model.EmailHistory, error)
	GetUserByID(ctx context.Context, id string) (*model.User, error)
	GetUserByEmail(ctx context.Context, email string) (*model.User, error)
	GetUserByUsername(ctx context.Context, username string) (*model.User, error)
	GetUserByPhoneNo(ctx context.Context, phoneNo string) (*model.User, error)
	CreateUser(ctx context.Context, email string) (*model.User, error)
	UpdatePassword(ctx context.Context, user *model.User, password string) (*model.User, error)
	InsertNewOTPKey(ctx context.Context, email, otp string) error
	GetOTPValue(ctx context.Context, email string) (string, error)
	CreateEmailHistory(ctx context.Context, tx postgre.Transaction, email string) error
	UpdateUser(ctx context.Context, tx postgre.Transaction, user *model.User) error
	DeleteOTPValue(ctx context.Context, email string) (int64, error)
	CreateUserGoogle(ctx context.Context, tx postgre.Transaction, user *model.User) (*model.User, error)
}
