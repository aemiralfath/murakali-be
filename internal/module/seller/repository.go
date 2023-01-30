package seller

import (
	"context"
	"murakali/internal/model"
	"murakali/internal/module/seller/delivery/body"
	"murakali/pkg/pagination"
	"murakali/pkg/postgre"
	"time"
)

type Repository interface {
	GetPerformance(ctx context.Context, shopID string) (*body.SellerPerformance, error)
	GetPerformaceRedis(ctx context.Context, key string) (*body.SellerPerformance, error)
	InsertPerformaceRedis(ctx context.Context, key string, value *body.SellerPerformance) error
	GetTotalAllSeller(ctx context.Context, shopName string) (int64, error)
	GetAllSeller(ctx context.Context, pgn *pagination.Pagination, shopName string) ([]*body.SellerResponse, error)
	GetTotalOrder(ctx context.Context, userID, orderStatusID, voucherShopID string) (int64, error)
	GetOrders(ctx context.Context, userID, orderStatusID, voucherShopID string, pgn *pagination.Pagination) ([]*model.Order, error)
	GetShopIDByUser(ctx context.Context, userID string) (string, error)
	GetShopIDByOrder(ctx context.Context, OrderID string) (string, error)
	ChangeOrderStatus(ctx context.Context, requestBody body.ChangeOrderStatusRequest) error
	CancelOrderStatus(ctx context.Context, tx postgre.Transaction, requestBody body.CancelOrderStatus) error
	CreateRefundSeller(ctx context.Context, tx postgre.Transaction, requestBody body.CancelOrderStatus) error
	GetOrderByOrderID(ctx context.Context, OrderID string) (*model.Order, error)
	GetSellerBySellerID(ctx context.Context, sellerID string) (*body.SellerResponse, error)
	GetSellerByUserID(ctx context.Context, userID string) (*body.SellerResponse, error)
	UpdateSellerInformationByUserID(ctx context.Context, shopName, userID string) error
	GetCourierByID(ctx context.Context, courierID string) (string, error)
	GetCourierSellerNotNullByShopAndCourierID(ctx context.Context, shopID, courierID string) (string, error)
	GetShopIDByUserID(ctx context.Context, userID string) (string, error)
	CreateCourierSeller(ctx context.Context, shopID string, courierID string) error
	GetCourierSellerByID(ctx context.Context, shopCourierID string) (string, error)
	UpdateCourierSellerByID(ctx context.Context, shopID, courierID string) error
	DeleteCourierSellerByID(ctx context.Context, shopCourierID string) error
	GetCategoryBySellerID(ctx context.Context, shopID string) ([]*body.CategoryResponse, error)
	GetAllCourier(ctx context.Context) ([]*body.CourierInfo, error)
	GetCourierSeller(ctx context.Context, userID string) ([]*body.CourierSellerRelationInfo, error)
	GetBuyerIDByOrderID(ctx context.Context, orderID string) (string, error)
	GetSellerIDByOrderID(ctx context.Context, orderID string) (string, error)
	GetAddressByBuyerID(ctx context.Context, userID string) (*model.Address, error)
	GetAddressBySellerID(ctx context.Context, userID string) (*model.Address, error)
	UpdateResiNumberInOrderSeller(ctx context.Context, noResi, orderID, shopID string, arriveAt time.Time) error
	GetCostRedis(ctx context.Context, key string) (*string, error)
	GetOrdersOnDelivery(ctx context.Context) ([]*model.OrderModel, error)
	InsertCostRedis(ctx context.Context, key string, value string) error
	CountCodeVoucher(ctx context.Context, code string) (int64, error)
	GetAllVoucherSeller(ctx context.Context, shopID, voucherStatusID, sortFilter string, pgn *pagination.Pagination) ([]*model.Voucher, error)
	GetTotalVoucherSeller(ctx context.Context, shopID, voucherStatusID string) (int64, error)
	CreateVoucherSeller(ctx context.Context, voucherShop *model.Voucher) error
	UpdateVoucherSeller(ctx context.Context, voucherShop *model.Voucher) error
	DeleteVoucherSeller(ctx context.Context, voucherIDShopID *body.VoucherIDShopID) error
	GetAllVoucherSellerByIDAndShopID(ctx context.Context, voucherIDShopID *body.VoucherIDShopID) (*model.Voucher, error)
	UpdateOrder(ctx context.Context, tx postgre.Transaction, orderData *model.OrderModel) error
	UpdateTransaction(ctx context.Context, tx postgre.Transaction, transactionData *model.Transaction) error
	GetOrderByTransactionID(ctx context.Context, tx postgre.Transaction, transactionID string) ([]*model.OrderModel, error)
	GetTransactionsExpired(ctx context.Context) ([]*model.Transaction, error)
	GetOrderItemsByOrderID(ctx context.Context, tx postgre.Transaction, orderID string) ([]*model.OrderItem, error)
	GetProductDetailByID(ctx context.Context, tx postgre.Transaction, productDetailID string) (*model.ProductDetail, error)
	UpdateProductDetailStock(ctx context.Context, tx postgre.Transaction, productDetailData *model.ProductDetail) error
	GetAllPromotionSeller(ctx context.Context, shopID string, promoStatusID string) ([]*body.PromotionSellerResponse, error)
	GetTotalPromotionSeller(ctx context.Context, shopID string, promoStatusID string) (int64, error)
	GetProductPromotion(ctx context.Context, shopProduct *body.ShopProduct) (*body.ProductPromotion, error)
	CreatePromotionSeller(ctx context.Context, tx postgre.Transaction, promotionShop *model.Promotion) error
	GetPromotionSellerDetailByID(ctx context.Context, shopProductPromo *body.ShopProductPromo) (*body.PromotionSellerResponse, error)
	UpdatePromotionSeller(ctx context.Context, promotion *model.Promotion) error
	GetDetailPromotionSellerByID(ctx context.Context, shopProductPromo *body.ShopProductPromo) (*body.PromotionDetailSeller, error)
	GetTotalProductWithoutPromotionSeller(ctx context.Context, shopID string) (int64, error)
	GetProductWithoutPromotionSeller(ctx context.Context, shopID string, pgn *pagination.Pagination) ([]*body.GetProductWithoutPromotion, error)
	GetWalletByUserID(ctx context.Context, tx postgre.Transaction, userID string) (*model.Wallet, error)
	UpdateWalletBalance(ctx context.Context, tx postgre.Transaction, wallet *model.Wallet) error
	InsertWalletHistory(ctx context.Context, tx postgre.Transaction, walletHistory *model.WalletHistory) error
	GetOrderModelByID(ctx context.Context, OrderID string) (*model.OrderModel, error)
	GetRefundOrderByOrderID(ctx context.Context, orderID string) (*model.Refund, error)
	GetRefundOrderByID(ctx context.Context, refundID string) (*model.Refund, error)
	GetRefundThreadByRefundID(ctx context.Context, refundID string) ([]*model.RefundThread, error)
	CreateRefundThreadSeller(ctx context.Context, refundThreadData *model.RefundThread) error
	UpdateRefundAccept(ctx context.Context, refundDataID string) error
	UpdateRefundReject(ctx context.Context, refundDataID string) error
}
