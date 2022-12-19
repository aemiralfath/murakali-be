package user

import (
	"context"
	"murakali/internal/model"
	body2 "murakali/internal/module/user/delivery/body"
)

type UseCase interface {
	EditUser(ctx context.Context, userID string, body body2.EditUserRequest) (*model.User, error)
}
