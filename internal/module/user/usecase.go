package user

import (
	"context"
	"murakali/internal/model"
	"murakali/internal/module/user/delivery/body"
	body2 "murakali/internal/module/user/delivery/body"
)

type UseCase interface {
	EditUser(ctx context.Context, userID string, body body2.EditUserRequest) (*model.User, error)
	EditEmail(ctx context.Context, userID string, requestBody body.EditEmailRequest) (*model.User, error)
	EditEmailUser(ctx context.Context, userID string, requestBody body.EditEmailUserRequest) (*model.User, error)
}
