package auth

import (
	"context"
	"murakali/internal/model"
	body2 "murakali/internal/module/auth/delivery/body"
)

type UseCase interface {
	RegisterEmail(ctx context.Context, body body2.RegisterEmailRequest) (*model.User, error)
	RegisterUser(ctx context.Context, email string, body body2.RegisterUserRequest) error
	VerifyOTP(ctx context.Context, body body2.VerifyOTPRequest) (string, error)
	ResetPasswordVerifyOTP(ctx context.Context, body body2.ResetPasswordVerifyOTPRequest) (string, error)
	Login(ctx context.Context, body body2.LoginRequest) (*model.Token, error)
	RefreshToken(ctx context.Context, id string) (*model.AccessToken, error)
	ResetPasswordEmail(ctx context.Context, body body2.ResetPasswordEmailRequest) (*model.User, error)
	ResetPasswordUser(ctx context.Context, email string, body *body2.ResetPasswordUserRequest) (*model.User, error)
	CheckUniqueUsername(ctx context.Context, username string) (bool, error)
	CheckUniquePhoneNo(ctx context.Context, phoneNo string) (bool, error)
}
