package user

import (
	"context"
	"murakali/internal/model"
)

type Repository interface {
	GetUserByID(ctx context.Context, id string) (*model.User, error)
	GetUserByUsername(ctx context.Context, username string) (*model.User, error)
	GetUserByPhoneNo(ctx context.Context, phoneNo string) (*model.User, error)
	UpdateUserField(ctx context.Context, user *model.User) error
}
