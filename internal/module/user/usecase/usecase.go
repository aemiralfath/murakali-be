package usecase

import (
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
