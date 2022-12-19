package user

import (
	"context"
	"murakali/internal/model"
)

type Repository interface {
	GetSealabsPay(ctx context.Context, userid string) ([]*model.SealabsPay, error)
}
