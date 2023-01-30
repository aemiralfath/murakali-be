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
	GetTransactionByUserID(ctx context.Context, userID string, status int, pgn *pagination.Pagination) ([]*model.Transaction, error)
	GetTotalTransactionByUserID(ctx context.Context, userID string) (int64, error)
	GetUserByID(ctx context.Context, id string) (*model.User, error)
	GetUserPasswordByID(ctx context.Context, id string) (*model.User, error)
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
	AddSealabsPayTrans(ctx context.Context, tx postgre.Transaction, request body.AddSealabsPayRequest, userid string) error
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
	ChangeOrderStatus(ctx context.Context, requestBody body.ChangeOrderStatusRequest) error
	GetOrdersByTransactionID(ctx context.Context, transactionID, userID string) ([]*model.Order, error)
	GetTotalOrder(ctx context.Context, userID, orderStatusID string) (int64, error)
	GetProductUnitSoldByOrderID(ctx context.Context, tx postgre.Transaction, orderID string) ([]*body.ProductUnitSoldOrderQty, error)
	UpdateProductUnitSold(ctx context.Context, tx postgre.Transaction, productID string, newQty int64) error
	GetOrders(ctx context.Context, userID, orderStatusID string, pgn *pagination.Pagination) ([]*model.Order, error)
	GetOrderByOrderID(ctx context.Context, OrderID string) (*model.Order, error)
	GetSellerIDByOrderID(ctx context.Context, orderID string) (string, error)
	GetAddressByBuyerID(ctx context.Context, userID string) (*model.Address, error)
	GetAddressBySellerID(ctx context.Context, userID string) (*model.Address, error)
	GetBuyerIDByOrderID(ctx context.Context, orderID string) (string, error)
	GetCostRedis(ctx context.Context, key string) (*string, error)
	InsertCostRedis(ctx context.Context, key string, value string) error
	GetWalletUser(ctx context.Context, walletID string) (*model.Wallet, error)
	GetWalletHistoryByID(ctx context.Context, id string) (*model.WalletHistory, error)
	GetWalletHistoryByWalletID(ctx context.Context, pgn *pagination.Pagination, walletID string) ([]*body.HistoryWalletResponse, error)
	GetTotalWalletHistoryByWalletID(ctx context.Context, walletID string) (int64, error)
	GetSealabsPayUser(ctx context.Context, userID, CardNumber string) (*model.SealabsPay, error)
	GetVoucherMarketplaceByID(ctx context.Context, voucherMarketplaceID string) (*model.Voucher, error)
	GetShopByID(ctx context.Context, shopID string) (*model.Shop, error)
	GetVoucherShopByID(ctx context.Context, VoucherShopID, shopID string) (*model.Voucher, error)
	GetCourierShopByID(ctx context.Context, CourierID, shopID string) (*model.Courier, error)
	GetProductDetailByID(ctx context.Context, tx postgre.Transaction, productDetailID string) (*model.ProductDetail, error)
	GetCartItemUser(ctx context.Context, userID, productDetailID string) (*model.CartItem, error)
	CreateTransaction(ctx context.Context, tx postgre.Transaction, transactionData *model.Transaction) (*uuid.UUID, error)
	UpdateTransaction(ctx context.Context, tx postgre.Transaction, transactionData *model.Transaction) error
	CreateOrder(ctx context.Context, tx postgre.Transaction, orderData *model.OrderModel) (*uuid.UUID, error)
	CreateOrderItem(ctx context.Context, tx postgre.Transaction, item *model.OrderItem) (*uuid.UUID, error)
	UpdateProductDetailStock(ctx context.Context, tx postgre.Transaction, productDetailData *model.ProductDetail) error
	DeleteCartItemByID(ctx context.Context, tx postgre.Transaction, cartItemData *model.CartItem) error
	GetOrderByTransactionID(ctx context.Context, transactionID string) ([]*model.OrderModel, error)
	GetOrderDetailByTransactionID(ctx context.Context, TransactionID string) ([]*model.Order, error)
	UpdateOrder(ctx context.Context, tx postgre.Transaction, orderData *model.OrderModel) error
	CreateWallet(ctx context.Context, walletData *model.Wallet) error
	GetWalletByUserID(ctx context.Context, userID string) (*model.Wallet, error)
	InsertWalletHistory(ctx context.Context, tx postgre.Transaction, walletHistory *model.WalletHistory) error
	UpdateWalletBalance(ctx context.Context, tx postgre.Transaction, wallet *model.Wallet) error
	UpdateWallet(ctx context.Context, wallet *model.Wallet) error
	CheckUserSealabsPay(ctx context.Context, userid string) (int, error)
	CheckDeletedSealabsPay(ctx context.Context, cardNumber string) (int, error)
	AddSealabsPay(ctx context.Context, request body.AddSealabsPayRequest, userid string) error
	UpdateUserSealabsPay(ctx context.Context, request body.AddSealabsPayRequest, userid string) error
	UpdateUserSealabsPayTrans(ctx context.Context, tx postgre.Transaction, request body.AddSealabsPayRequest, userid string) error
	UpdateWalletPin(ctx context.Context, wallet *model.Wallet) error
	GetProductPromotionByProductID(ctx context.Context, productID string) (*model.Promotion, error)
	UpdateVoucherQuota(ctx context.Context, tx postgre.Transaction, upVoucher *model.Voucher) error
	UpdatePromotionQuota(ctx context.Context, tx postgre.Transaction, promo *model.Promotion) error
	GetOrderModelByID(ctx context.Context, OrderID string) (*model.OrderModel, error)
	GetRefundOrderByOrderID(ctx context.Context, orderID string) (*model.Refund, error)
	CreateRefundUser(ctx context.Context, tx postgre.Transaction, refundData *model.Refund) error
	UpdateOrderRefund(ctx context.Context, tx postgre.Transaction, orderID string, isRefund bool) error
	GetRefundOrderByID(ctx context.Context, refundID string) (*model.Refund, error)
	GetRefundThreadByRefundID(ctx context.Context, refundID string) ([]*body.RThread, error)
	CreateRefundThreadUser(ctx context.Context, refundThreadData *model.RefundThread) error
}
