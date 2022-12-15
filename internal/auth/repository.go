package auth

import (
	"context"
	"murakali/internal/model"
)

type Repository interface {
	CheckEmailHistory(ctx context.Context, email string) (*model.EmailHistory, error)
	GetUserByEmail(ctx context.Context, email string) (*model.User, error)
	CreateUser(ctx context.Context, email string) (*model.User, error)
	InsertNewOTPKey(ctx context.Context, email, otp string) error
}
