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
}
