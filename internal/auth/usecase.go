package auth

import (
	"context"
	"murakali/internal/auth/delivery/body"
	"murakali/internal/model"
)

type UseCase interface {
	RegisterEmail(ctx context.Context, body body.RegisterEmailRequest) (*model.User, error)
	RegisterUser(ctx context.Context, email string, body body.RegisterUserRequest) error
	VerifyOTP(ctx context.Context, body body.VerifyOTPRequest) (string, error)
	ResetPasswordVerifyOTP(ctx context.Context, body body.ResetPasswordVerifyOTPRequest) (string, error)
	Login(ctx context.Context, body body.LoginRequest) (string, string, error)
	RefreshToken(ctx context.Context, id string) (string, error)
	ResetPasswordEmail(ctx context.Context, body body.ResetPasswordEmailRequest) (*model.User, error)
	ResetPasswordUser(ctx context.Context, email string, body *body.ResetPasswordUserRequest) (*model.User, error)
}
