package usecase

import (
	"context"
	"murakali/internal/model"

	"murakali/config"
	"murakali/internal/module/user"
	"murakali/pkg/postgre"
)

type userUC struct {
	cfg      *config.Config
	txRepo   *postgre.TxRepo
	userRepo user.Repository
}

func NewUserUseCase(cfg *config.Config, txRepo *postgre.TxRepo, userRepo user.Repository) user.UseCase {
	return &userUC{cfg: cfg, txRepo: txRepo, userRepo: userRepo}
}

func (u *userUC) GetSealabsPay(ctx context.Context, userid string) ([]*model.SealabsPay, error) {

	response, err := u.userRepo.GetSealabsPay(ctx, userid)
	if err != nil {
		return nil, err
	}

	return response, nil
}
