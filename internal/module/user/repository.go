package user

import (
	"context"
	"murakali/internal/model"
	"murakali/internal/module/user/delivery/body"
	"murakali/pkg/postgre"
)

type Repository interface {
	GetSealabsPay(ctx context.Context, userid string) ([]*model.SealabsPay, error)
	AddSealabsPay(ctx context.Context, tx postgre.Transaction, request body.AddSealabsPayRequest) error
	CheckDefaultSealabsPay(ctx context.Context, userid string) (*string, error)
	SetDefaultSealabsPay(ctx context.Context, tx postgre.Transaction, card_number *string) error
}
