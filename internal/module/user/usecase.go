package user

import (
	"context"
	"murakali/internal/model"
	"murakali/internal/module/user/delivery/body"
)

type UseCase interface {
	GetSealabsPay(ctx context.Context, userid string) ([]*model.SealabsPay, error)
	AddSealabsPay(ctx context.Context, request body.AddSealabsPayRequest, userid string) error
	PatchSealabsPay(ctx context.Context, card_number string, userid string) error
	DeleteSealabsPay(ctx context.Context, card_number string) error
	EditUser(ctx context.Context, userID string, requestBody body.EditUserRequest) (*model.User, error)
	EditEmail(ctx context.Context, userID string, requestBody body.EditEmailRequest) (*model.User, error)
	EditEmailUser(ctx context.Context, userID string, requestBody body.EditEmailUserRequest) (*model.User, error)
}
