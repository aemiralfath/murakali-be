package user

import (
	"context"
	"murakali/internal/model"
	"murakali/internal/module/user/delivery/body"
	"murakali/pkg/pagination"
	"murakali/pkg/postgre"
)

type Repository interface {
	GetUserByID(ctx context.Context, id string) (*model.User, error)
	GetTotalAddress(ctx context.Context, userID, name string, isDefault, isShopDefault bool) (int64, error)
	GetAddresses(ctx context.Context, userID, name string, isDefault, isShopDefault bool, pagination *pagination.Pagination) ([]*model.Address, error)
	GetAddressByID(ctx context.Context, userID, addressID string) (*model.Address, error)
	GetDefaultUserAddress(ctx context.Context, userID string) (*model.Address, error)
	GetDefaultShopAddress(ctx context.Context, userID string) (*model.Address, error)
	CreateAddress(ctx context.Context, tx postgre.Transaction, userID string, requestBody body.CreateAddressRequest) error
	UpdateAddress(ctx context.Context, tx postgre.Transaction, address *model.Address) error
	UpdateDefaultAddress(ctx context.Context, tx postgre.Transaction, status bool, address *model.Address) error
	UpdateDefaultShopAddress(ctx context.Context, tx postgre.Transaction, status bool, address *model.Address) error
	DeleteAddress(ctx context.Context, addressID string) error
	GetSealabsPay(ctx context.Context, userid string) ([]*model.SealabsPay, error)
	AddSealabsPay(ctx context.Context, tx postgre.Transaction, request body.AddSealabsPayRequest) error
	PatchSealabsPay(ctx context.Context, cardNumber string) error
	CheckDefaultSealabsPay(ctx context.Context, userid string) (*string, error)
	SetDefaultSealabsPayTrans(ctx context.Context, tx postgre.Transaction, cardNumber *string) error
	SetDefaultSealabsPay(ctx context.Context, cardNumber string, userid string) error
	DeleteSealabsPay(ctx context.Context, cardNmber string) error
	CheckEmailHistory(ctx context.Context, email string) (*model.EmailHistory, error)
	InsertNewOTPKey(ctx context.Context, email, otp string) error
	GetOTPValue(ctx context.Context, email string) (string, error)
	DeleteOTPValue(ctx context.Context, email string) (int64, error)
	GetUserByUsername(ctx context.Context, username string) (*model.User, error)
	GetUserByPhoneNo(ctx context.Context, phoneNo string) (*model.User, error)
	UpdateUserField(ctx context.Context, user *model.User) error
	UpdateUserEmail(ctx context.Context, tx postgre.Transaction, user *model.User) error
	CreateEmailHistory(ctx context.Context, tx postgre.Transaction, email string) error
	CheckShopByID(ctx context.Context, userID string) (int64, error)
	CheckShopUnique(ctx context.Context, shopName string) (int64, error)
	AddShop(ctx context.Context, userID string, shopName string) error
	UpdateRole(ctx context.Context, userID string) error
	UpdateProfileImage(ctx context.Context, imgURL, userID string) error
	UpdatePasswordByID(ctx context.Context, userID, newPassword string) error
	GetPasswordByID(ctx context.Context, id string) (string, error)
	GetTotalOrder(ctx context.Context, userID string) (int64, error)
	GetOrders(ctx context.Context, userID string, pgn *pagination.Pagination) ([]*model.Order, error)
}
