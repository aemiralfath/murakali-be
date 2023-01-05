package user

import (
	"context"
	"murakali/internal/model"
	"murakali/internal/module/user/delivery/body"
	"murakali/pkg/pagination"
	"murakali/pkg/postgre"

	"github.com/google/uuid"
)

type Repository interface {
	GetTransactionByID(ctx context.Context, transactionID string) (*model.Transaction, error)
	GetUserByID(ctx context.Context, id string) (*model.User, error)
	GetTotalAddress(ctx context.Context, userID, name string) (int64, error)
	GetTotalAddressDefault(ctx context.Context, userID, name string, isDefault, isShopDefault bool) (int64, error)
	GetAddresses(ctx context.Context, userID, name string, isDefault, isShopDefault bool, pagination *pagination.Pagination) ([]*model.Address, error)
	GetAllAddresses(ctx context.Context, userID, name string, pagination *pagination.Pagination) ([]*model.Address, error)
	GetAddressByID(ctx context.Context, userID, addressID string) (*model.Address, error)
	GetDefaultUserAddress(ctx context.Context, userID string) (*model.Address, error)
	GetDefaultShopAddress(ctx context.Context, userID string) (*model.Address, error)
	CreateAddress(ctx context.Context, tx postgre.Transaction, userID string, requestBody body.CreateAddressRequest) error
	UpdateAddress(ctx context.Context, tx postgre.Transaction, address *model.Address) error
	UpdateDefaultAddress(ctx context.Context, tx postgre.Transaction, status bool, address *model.Address) error
	UpdateDefaultShopAddress(ctx context.Context, tx postgre.Transaction, status bool, address *model.Address) error
	DeleteAddress(ctx context.Context, addressID string) error
	GetSealabsPay(ctx context.Context, userid string) ([]*model.SealabsPay, error)
	AddSealabsPay(ctx context.Context, tx postgre.Transaction, request body.AddSealabsPayRequest, userid string) error
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
	GetWalletUser(ctx context.Context, userID, walletID string) (*model.Wallet, error)
	GetSealabsPayUser(ctx context.Context, userID, CardNumber string) (*model.SealabsPay, error)
	GetVoucherMarketplaceByID(ctx context.Context, voucherMarketplaceID string) (*model.Voucher, error)
	GetShopByID(ctx context.Context, shopID string) (*model.Shop, error)
	GetVoucherShopByID(ctx context.Context, VoucherShopID, shopID string) (*model.Voucher, error)
	GetCourierShopByID(ctx context.Context, CourierID, shopID string) (*model.Courier, error)
	GetProductDetailByID(ctx context.Context, productDetailID string) (*model.ProductDetail, error)
	GetCartItemUser(ctx context.Context, userID, productDetailID string) (*model.CartItem, error)
	CreateTransaction(ctx context.Context, tx postgre.Transaction, transactionData *model.Transaction) (*uuid.UUID, error)
	UpdateTransaction(ctx context.Context, tx postgre.Transaction, transactionData *model.Transaction) error
	CreateOrder(ctx context.Context, tx postgre.Transaction, orderData *model.OrderModel) (*uuid.UUID, error)
	CreateOrderItem(ctx context.Context, tx postgre.Transaction, item *model.OrderItem) (*uuid.UUID, error)
	UpdateProductDetailStock(ctx context.Context, tx postgre.Transaction, productDetailData *model.ProductDetail) error
	DeleteCartItemByID(ctx context.Context, tx postgre.Transaction, cartItemData *model.CartItem) error
	GetOrderByTransactionID(ctx context.Context, transactionID string) ([]*model.OrderModel, error)
	UpdateOrder(ctx context.Context, tx postgre.Transaction, orderData *model.OrderModel) error
	CreateWallet(ctx context.Context, walletData *model.Wallet) error
	GetWalletByUserID(ctx context.Context, userID string) (*model.Wallet, error)
}
