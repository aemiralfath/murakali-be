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
	PatchSealabsPay(ctx context.Context, card_number string) error
	CheckDefaultSealabsPay(ctx context.Context, userid string) (*string, error)
	SetDefaultSealabsPayTrans(ctx context.Context, tx postgre.Transaction, card_number *string) error
	SetDefaultSealabsPay(ctx context.Context, card_number string, userid string) error
	DeleteSealabsPay(ctx context.Context, card_number string) error
	GetUserByID(ctx context.Context, id string) (*model.User, error)
	CheckEmailHistory(ctx context.Context, email string) (*model.EmailHistory, error)
	InsertNewOTPKey(ctx context.Context, email, otp string) error
	GetOTPValue(ctx context.Context, email string) (string, error)
	DeleteOTPValue(ctx context.Context, email string) (int64, error)
	GetUserByUsername(ctx context.Context, username string) (*model.User, error)
	GetUserByPhoneNo(ctx context.Context, phoneNo string) (*model.User, error)
	UpdateUserField(ctx context.Context, user *model.User) error
	UpdateUserEmail(ctx context.Context, tx postgre.Transaction, user *model.User) error
	CreateEmailHistory(ctx context.Context, tx postgre.Transaction, email string) error
}
