package usecase

import (
	"context"
	"murakali/internal/model"

	"murakali/config"
	"murakali/internal/module/user"
	"murakali/internal/module/user/delivery/body"
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

func (u *userUC) AddSealabsPay(ctx context.Context, request body.AddSealabsPayRequest, userid string) error {

	card_number, err := u.userRepo.CheckDefaultSealabsPay(ctx, userid)
	if err != nil {
		return err
	}

	err = u.txRepo.WithTransaction(func(tx postgre.Transaction) error {
		if u.userRepo.SetDefaultSealabsPayTrans(ctx, tx, card_number) != nil {
			return err
		}

		err = u.userRepo.AddSealabsPay(ctx, tx, request)
		if err != nil {
			return err
		}
		return nil
	})
	return nil
}

func (u *userUC) PatchSealabsPay(ctx context.Context, card_number string, userid string) error {
	err := u.userRepo.PatchSealabsPay(ctx, card_number)
	if err != nil {
		return err
	}

	if u.userRepo.SetDefaultSealabsPay(ctx, card_number, userid) != nil {
		return err
	}
	return nil
}
