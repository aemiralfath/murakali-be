package auth

import (
	"context"
	"murakali/internal/auth/delivery/body"
	"murakali/internal/model"
)

type UseCase interface {
	Register(ctx context.Context, body body.RegisterRequest) (*model.User, error)
}
