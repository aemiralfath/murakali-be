package auth

import (
	"context"
	"murakali/internal/auth/delivery/body"
	"murakali/internal/model"
)

type UseCase interface {
	RegisterEmail(ctx context.Context, body body.RegisterEmailRequest) (*model.User, error)
	RegisterUser(ctx context.Context, body body.RegisterUserRequest) error
	VerifyOTP(ctx context.Context, body body.VerifyOTPRequest) error
}
