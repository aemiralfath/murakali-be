package user

import (
	"context"
	"murakali/internal/model"
	"murakali/internal/module/user/delivery/body"
	"murakali/pkg/pagination"
)

type UseCase interface {
	GetAddress(ctx context.Context, userID string, pagination *pagination.Pagination,
		queryRequest *body.GetAddressQueryRequest) (*pagination.Pagination, error)
	CreateAddress(ctx context.Context, userID string, requestBody body.CreateAddressRequest) error
	GetAddressByID(ctx context.Context, userID, addressID string) (*model.Address, error)
	UpdateAddressByID(ctx context.Context, userID, addressID string, requestBody body.UpdateAddressRequest) error
	DeleteAddressByID(ctx context.Context, userID, addressID string) error
	EditUser(ctx context.Context, userID string, requestBody body.EditUserRequest) (*model.User, error)
	EditEmail(ctx context.Context, userID string, requestBody body.EditEmailRequest) (*model.User, error)
	EditEmailUser(ctx context.Context, userID string, requestBody body.EditEmailUserRequest) (*model.User, error)
	GetSealabsPay(ctx context.Context, userid string) ([]*model.SealabsPay, error)
	AddSealabsPay(ctx context.Context, request body.AddSealabsPayRequest, userid string) error
	PatchSealabsPay(ctx context.Context, cardNumber string, userid string) error
	DeleteSealabsPay(ctx context.Context, userID, cardNumber string) error
	RegisterMerchant(ctx context.Context, userID string, shopName string) error
	GetUserProfile(ctx context.Context, userID string) (*body.ProfileResponse, error)
	UploadProfilePicture(ctx context.Context, imgURL, userID string) error
	VerifyPasswordChange(ctx context.Context, userID string) error
	SendOTPEmail(ctx context.Context, email string) error
	VerifyOTP(ctx context.Context, requestBody body.VerifyOTPRequest, userID string) (string, error)
	ChangePassword(ctx context.Context, userID string, newPassword string) error
	GetTransactionByID(ctx context.Context, transactionID string) (*body.GetTransactionByIDResponse, error)
	GetTransactionByUserID(ctx context.Context, userID string, status int, pgn *pagination.Pagination) (*pagination.Pagination, error)
	CreateTransaction(ctx context.Context, userID string, requestBody body.CreateTransactionRequest) (string, error)
	UpdateTransaction(ctx context.Context, transactionID string, requestBody body.SLPCallbackRequest) error
	UpdateTransactionPaymentMethod(ctx context.Context, transactionID, cardNumber string) error
	UpdateWalletTransaction(ctx context.Context, transactionID string, requestBody body.SLPCallbackRequest) error
	CreateSLPPayment(ctx context.Context, transactionID string) (string, error)
	GetTransactionDetailByID(ctx context.Context, transactionID, userID string) (*body.TransactionDetailResponse, error)
	GetOrder(ctx context.Context, userID, orderStatusID string, pgn *pagination.Pagination) (*pagination.Pagination, error)
	GetOrderByOrderID(ctx context.Context, orderID string) (*model.Order, error)
	ChangeOrderStatus(ctx context.Context, userID string, requestBody body.ChangeOrderStatusRequest) error
	ActivateWallet(ctx context.Context, userID, pin string) error
	GetWallet(ctx context.Context, userID string) (*model.Wallet, error)
	GetDetailWalletHistory(ctx context.Context, walletHistoryID, userID string) (*body.DetailHistoryWalletResponse, error)
	GetWalletHistory(ctx context.Context, userID string, pgn *pagination.Pagination) (*pagination.Pagination, error)
	TopUpWallet(ctx context.Context, userID string, requestBody body.TopUpWalletRequest) (string, error)
	WalletStepUp(ctx context.Context, userID string, requestBody body.WalletStepUpRequest) (string, error)
	CreateWalletPayment(ctx context.Context, transactionID string) error
	ChangeWalletPinStepUp(ctx context.Context, userID string, requestBody body.ChangeWalletPinStepUpRequest) (string, error)
	ChangeWalletPinStepUpEmail(ctx context.Context, userID string) error
	ChangeWalletPin(ctx context.Context, userID, pin string) error
	CreateRefundUser(ctx context.Context, userID string, requestBody body.CreateRefundUserRequest) error
	GetRefundOrder(ctx context.Context, userID string, refundID string) (*body.GetRefundThreadResponse, error)
	CreateRefundThreadUser(ctx context.Context, userID string, requestBody *body.CreateRefundThreadRequest) error
}
