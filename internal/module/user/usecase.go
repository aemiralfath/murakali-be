package user

import (
	"context"
	"murakali/internal/model"
)

type UseCase interface {
	GetSealabsPay(ctx context.Context, userid string) ([]*model.SealabsPay, error)
}
