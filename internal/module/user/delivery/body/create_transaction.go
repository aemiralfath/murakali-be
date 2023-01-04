package body

import (
	"murakali/internal/model"
	"murakali/pkg/httperror"
	"murakali/pkg/response"
	"net/http"
	"strings"
)

type CreateTransactionRequest struct {
	WalletID             string     `json:"wallet_id"`
	CardNumber           string     `json:"card_number"`
	VoucherMarketplaceID string     `json:"voucher_marketplace_id"`
	CartItems            []CartItem `json:"cart_items"`
}

type CreateTransactionResponse struct {
	TransactionID string `json:"transaction_id"`
}

type CartItem struct {
	ShopID         string          `json:"shop_id"`
	VoucherShopID  string          `json:"voucher_shop_id"`
	CourierID      string          `json:"courier_id"`
	ProductDetails []ProductDetail `json:"product_details"`
}

type ProductDetail struct {
	ID       string  `json:"id"`
	Quantity int     `json:"quantity"`
	SubPrice float64 `json:"sub_price"`
}
type TransactionResponse struct {
	TransactionData *model.Transaction
	OrderResponses  []*OrderResponse
}

type OrderResponse struct {
	OrderData *model.OrderModel
	Items     []*OrderItemResponse
}

type OrderItemResponse struct {
	Item              *model.OrderItem
	ProductDetailData *model.ProductDetail
	CartItemData      *model.CartItem
}

func (r *CreateTransactionRequest) Validate() (UnprocessableEntity, error) {
	unprocessableEntity := false
	entity := UnprocessableEntity{
		Fields: map[string]string{
			"payment_method": "",
			"cart_items":     "",
		},
	}

	r.WalletID = strings.TrimSpace(r.WalletID)
	r.CardNumber = strings.TrimSpace(r.CardNumber)
	if (r.WalletID == "" && r.CardNumber == "") || (r.WalletID != "" && r.CardNumber != "") {
		unprocessableEntity = true
		entity.Fields["payment_method"] = InvalidPaymentMethod
	}

	if r.CartItems == nil {
		unprocessableEntity = true
		entity.Fields["cart_items"] = FieldCannotBeEmptyMessage
	}

	if unprocessableEntity {
		return entity, httperror.New(
			http.StatusUnprocessableEntity,
			response.UnprocessableEntityMessage,
		)
	}

	return entity, nil
}
