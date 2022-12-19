package user

import (
	"context"
	"murakali/internal/model"
	"murakali/pkg/postgre"
)

type Repository interface {
	GetUserByID(ctx context.Context, id string) (*model.User, error)
	CheckEmailHistory(ctx context.Context, email string) (*model.EmailHistory, error)
	InsertNewOTPKey(ctx context.Context, email, otp string) error
	GetOTPValue(ctx context.Context, email string) (string, error)
	DeleteOTPValue(ctx context.Context, email string) (int64, error)
	GetUserByUsername(ctx context.Context, username string) (*model.User, error)
	GetUserByPhoneNo(ctx context.Context, phoneNo string) (*model.User, error)
	UpdateUserField(ctx context.Context, user *model.User) error
	UpdateUserEmail(ctx context.Context, tx postgre.Transaction, user *model.User) error
	CreateEmailHistory(ctx context.Context, tx postgre.Transaction, email string) error
}
