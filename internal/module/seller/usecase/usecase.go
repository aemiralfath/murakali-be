package usecase

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"math"
	"murakali/config"
	"murakali/internal/constant"
	"murakali/internal/model"
	body2 "murakali/internal/module/location/delivery/body"
	"murakali/internal/module/seller"
	"murakali/internal/module/seller/delivery/body"
	"murakali/pkg/httperror"
	"murakali/pkg/pagination"
	"murakali/pkg/postgre"
	"murakali/pkg/response"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
)

type sellerUC struct {
	cfg        *config.Config
	txRepo     *postgre.TxRepo
	sellerRepo seller.Repository
}

func NewSellerUseCase(cfg *config.Config, txRepo *postgre.TxRepo, sellerRepo seller.Repository) seller.UseCase {
	return &sellerUC{cfg: cfg, txRepo: txRepo, sellerRepo: sellerRepo}
}

func (u *sellerUC) GetOrder(ctx context.Context, userID, orderStatusID, voucherShopID string,
	pgn *pagination.Pagination) (*pagination.Pagination, error) {
	shopID, err := u.sellerRepo.GetShopIDByUser(ctx, userID)
	if err != nil {
		return nil, err
	}

	totalRows, err := u.sellerRepo.GetTotalOrder(ctx, shopID, orderStatusID, voucherShopID)
	if err != nil {
		return nil, err
	}

	totalPages := int(math.Ceil(float64(totalRows) / float64(pgn.Limit)))
	pgn.TotalRows = totalRows
	pgn.TotalPages = totalPages

	orders, err := u.sellerRepo.GetOrders(ctx, shopID, orderStatusID, voucherShopID, pgn)
	if err != nil {
		return nil, err
	}

	pgn.Rows = orders
	return pgn, nil
}

func (u *sellerUC) ChangeOrderStatus(ctx context.Context, userID string, requestBody body.ChangeOrderStatusRequest) error {
	shopIDFromUser, err := u.sellerRepo.GetShopIDByUser(ctx, userID)
	if err != nil {
		return err
	}

	shopIDFromOrder, err := u.sellerRepo.GetShopIDByOrder(ctx, requestBody.OrderID)
	if err != nil {
		return err
	}
	if shopIDFromUser != shopIDFromOrder {
		return httperror.New(http.StatusUnauthorized, response.UnauthorizedMessage)
	}

	err = u.sellerRepo.ChangeOrderStatus(ctx, requestBody)
	if err != nil {
		return err
	}
	return nil
}

func (u *sellerUC) GetOrderByOrderID(ctx context.Context, orderID string) (*model.Order, error) {
	order, err := u.sellerRepo.GetOrderByOrderID(ctx, orderID)
	if err != nil {
		return nil, err
	}

	buyerID, err := u.sellerRepo.GetBuyerIDByOrderID(ctx, orderID)
	if err != nil {
		return nil, err
	}

	sellerID, err := u.sellerRepo.GetSellerIDByOrderID(ctx, orderID)
	if err != nil {
		return nil, err
	}

	buyerAddress, err := u.sellerRepo.GetAddressByBuyerID(ctx, buyerID)
	if err != nil {
		return nil, err
	}

	sellerAddress, err := u.sellerRepo.GetAddressBySellerID(ctx, sellerID)
	if err != nil {
		return nil, err
	}

	order.BuyerAddress = buyerAddress
	order.SellerAddress = sellerAddress

	totalWeight := 0
	for _, detail := range order.Detail {
		totalWeight += int(detail.ProductWeight) * detail.OrderQuantity
	}

	var costRedis *string
	key := fmt.Sprintf("%d:%d:%d:%s", sellerAddress.CityID, buyerAddress.CityID, totalWeight, order.CourierCode)
	costRedis, err = u.sellerRepo.GetCostRedis(ctx, key)
	if err != nil {
		res, err := u.GetCostRajaOngkir(sellerAddress.CityID, buyerAddress.CityID, totalWeight, order.CourierCode)
		if err != nil {
			return nil, err
		}

		redisValue, err := json.Marshal(res)
		if err != nil {
			return nil, err
		}

		if errInsert := u.sellerRepo.InsertCostRedis(ctx, key, string(redisValue)); errInsert != nil {
			return nil, errInsert
		}

		value := string(redisValue)
		costRedis = &value
	}

	var costResp body2.RajaOngkirCostResponse
	if err := json.Unmarshal([]byte(*costRedis), &costResp); err != nil {
		return nil, err
	}

	if len(costResp.Rajaongkir.Results) > 0 {
		for _, cost := range costResp.Rajaongkir.Results[0].Costs {
			if cost.Service == order.CourierService {
				order.CourierETD = cost.Cost[0].Etd
			}
		}
	}

	return order, nil
}

func (u *sellerUC) GetCourierSeller(ctx context.Context, userID string) (*body.CourierSellerResponse, error) {
	courier, err := u.sellerRepo.GetAllCourier(ctx)
	if err != nil {
		return nil, err
	}

	courierSeller, err := u.sellerRepo.GetCourierSeller(ctx, userID)
	if err != nil {
		return nil, err
	}

	resultCourierSeller := make([]*body.CourierSellerInfo, 0)
	totalData := len(courier)
	totalDataCourierSeller := len(courierSeller)
	for i := 0; i < totalData; i++ {
		var shopCourierIDTemp string
		var deletedAtTemp string
		for j := 0; j < totalDataCourierSeller; j++ {
			if courier[i].CourierID == courierSeller[j].CourierID {
				shopCourierIDTemp = courierSeller[j].ShopCourierID.String()
				if !courierSeller[j].DeletedAt.Time.IsZero() {
					deletedAtTemp = courierSeller[j].DeletedAt.Time.String()
				}
			}
		}
		p := &body.CourierSellerInfo{
			ShopCourierID: shopCourierIDTemp,
			CourierID:     courier[i].CourierID,
			Name:          courier[i].Name,
			Code:          courier[i].Code,
			Service:       courier[i].Service,
			Description:   courier[i].Description,
			DeletedAt:     deletedAtTemp,
		}

		resultCourierSeller = append(resultCourierSeller, p)
	}

	csr := &body.CourierSellerResponse{}
	csr.Rows = resultCourierSeller
	return csr, nil
}
func (u *sellerUC) GetSellerBySellerID(ctx context.Context, sellerID string) (*body.SellerResponse, error) {
	sellerData, err := u.sellerRepo.GetSellerBySellerID(ctx, sellerID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, httperror.New(http.StatusNotFound, body.SellerNotFoundMessage)
		}
		return nil, err
	}

	return sellerData, nil
}

func (u *sellerUC) GetSellerByUserID(ctx context.Context, userID string) (*body.SellerResponse, error) {
	sellerData, err := u.sellerRepo.GetSellerByUserID(ctx, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, httperror.New(http.StatusNotFound, body.SellerNotFoundMessage)
		}
		return nil, err
	}

	return sellerData, nil
}

func (u *sellerUC) CreateCourierSeller(ctx context.Context, userID, courierID string) error {
	_, err := u.sellerRepo.GetCourierByID(ctx, courierID)
	if err != nil {
		if err == sql.ErrNoRows {
			return httperror.New(http.StatusBadRequest, body.CourierNotFoundMessage)
		}
		return err
	}

	shopID, err := u.sellerRepo.GetShopIDByUserID(ctx, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return httperror.New(http.StatusBadRequest, response.ShopAddressNotFound)
		}
		return err
	}

	sellerCourierID, _ := u.sellerRepo.GetCourierSellerNotNullByShopAndCourierID(ctx, shopID, courierID)

	if sellerCourierID != "" {
		if err != nil {
			if err == sql.ErrNoRows {
				return httperror.New(http.StatusBadRequest, body.CourierSellerNotFoundMessage)
			}
			return err
		}

		err = u.sellerRepo.UpdateCourierSellerByID(ctx, shopID, courierID)
		if err == sql.ErrNoRows {
			return httperror.New(http.StatusBadRequest, body.CourierSellerAlreadyExistMessage)
		}
		return nil
	}

	err = u.sellerRepo.CreateCourierSeller(ctx, shopID, courierID)
	if err != nil {
		return err
	}
	return nil
}

func (u *sellerUC) DeleteCourierSellerByID(ctx context.Context, shopCourierID string) error {
	_, err := u.sellerRepo.GetCourierSellerByID(ctx, shopCourierID)
	if err != nil {
		if err == sql.ErrNoRows {
			return httperror.New(http.StatusBadRequest, body.CourierSellerNotFoundMessage)
		}
		return err
	}

	if err := u.sellerRepo.DeleteCourierSellerByID(ctx, shopCourierID); err != nil {
		if err == sql.ErrNoRows {
			return httperror.New(http.StatusNotFound, body.CourierSellerNotFoundMessage)
		}

		return err
	}
	return nil
}

func (u *sellerUC) GetCategoryBySellerID(ctx context.Context, shopID string) ([]*body.CategoryResponse, error) {
	categories, err := u.sellerRepo.GetCategoryBySellerID(ctx, shopID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, httperror.New(http.StatusNotFound, body.CategoryNotFoundMessage)
		}
		return nil, err
	}

	return categories, nil
}

func (u *sellerUC) UpdateResiNumberInOrderSeller(ctx context.Context, userID, orderID string, requestBody body.UpdateNoResiOrderSellerRequest) error {
	shopID, err := u.sellerRepo.GetShopIDByUserID(ctx, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return httperror.New(http.StatusBadRequest, response.ShopAddressNotFound)
		}
		return err
	}

	err = u.sellerRepo.UpdateResiNumberInOrderSeller(ctx, requestBody.NoResi, orderID, shopID, requestBody.EstimateArriveAtTime)
	if err != nil {
		return err
	}

	return nil
}

func (u *sellerUC) UpdateOnDeliveryOrder(ctx context.Context) error {
	orders, err := u.sellerRepo.GetOrdersOnDelivery(ctx)
	if err != nil {
		return nil
	}

	for _, order := range orders {
		if order.OrderStatusID == constant.OrderStatusOnDelivery && order.ArrivedAt.Valid && time.Until(order.ArrivedAt.Time) <= 0 {
			if err := u.sellerRepo.ChangeOrderStatus(ctx, body.ChangeOrderStatusRequest{OrderID: order.ID.String(),
				OrderStatusID: strconv.Itoa(constant.OrderStatusDelivered)}); err != nil {
				return err
			}
		}
	}

	return nil
}

func (u *sellerUC) UpdateExpiredAtOrder(ctx context.Context) error {
	// TODO: Add voucher promotion stock & validation
	transactions, err := u.sellerRepo.GetTransactionsExpired(ctx)
	if err != nil {
		return nil
	}

	for _, transaction := range transactions {
		err := u.txRepo.WithTransaction(func(tx postgre.Transaction) error {
			transaction.CanceledAt.Valid = true
			transaction.CanceledAt.Time = time.Now()
			if err := u.sellerRepo.UpdateTransaction(ctx, tx, transaction); err != nil {
				return err
			}

			orders, err := u.sellerRepo.GetOrderByTransactionID(ctx, tx, transaction.ID.String())
			if err != nil {
				return err
			}

			for _, order := range orders {
				order.OrderStatusID = constant.OrderStatusCanceled
				if err := u.sellerRepo.UpdateOrder(ctx, tx, order); err != nil {
					return err
				}

				orderItems, err := u.sellerRepo.GetOrderItemsByOrderID(ctx, tx, order.ID.String())
				if err != nil {
					return err
				}

				for _, item := range orderItems {
					productDetailData, err := u.sellerRepo.GetProductDetailByID(ctx, tx, item.ProductDetailID.String())
					if err != nil {
						return err
					}

					productDetailData.Stock += float64(item.Quantity)
					errProduct := u.sellerRepo.UpdateProductDetailStock(ctx, tx, productDetailData)
					if errProduct != nil {
						return errProduct
					}
				}
			}

			return nil
		})

		if err != nil {
			return err
		}
	}

	return nil
}

func (u *sellerUC) GetCostRajaOngkir(origin, destination, weight int, code string) (*body2.RajaOngkirCostResponse, error) {
	var responseCost body2.RajaOngkirCostResponse
	url := fmt.Sprintf("%s/cost", u.cfg.External.OngkirAPIURL)
	payload := fmt.Sprintf(
		"origin=%d&destination=%d&weight=%d&courier=%s", origin, destination, weight, code)

	req, err := http.NewRequest("POST", url, strings.NewReader(payload))
	if err != nil {
		return nil, err
	}

	req.Header.Add("key", u.cfg.External.OngkirAPIKey)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	readErr := json.NewDecoder(res.Body).Decode(&responseCost)
	if readErr != nil {
		return nil, err
	}

	return &responseCost, nil
}

func (u *sellerUC) GetAllVoucherSeller(ctx context.Context, userID, voucherStatusID string,
	pgn *pagination.Pagination) (*pagination.Pagination, error) {
	shopID, err := u.sellerRepo.GetShopIDByUserID(ctx, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, httperror.New(http.StatusBadRequest, response.UserNotHaveShop)
		}
		return nil, err
	}

	totalRows, err := u.sellerRepo.GetTotalVoucherSeller(ctx, shopID, voucherStatusID)
	if err != nil {
		return nil, err
	}

	totalPages := int(math.Ceil(float64(totalRows) / float64(pgn.Limit)))
	pgn.TotalRows = totalRows
	pgn.TotalPages = totalPages

	ShopVouchers, err := u.sellerRepo.GetAllVoucherSeller(ctx, shopID, voucherStatusID, pgn)
	if err != nil {
		return nil, err
	}

	pgn.Rows = ShopVouchers

	return pgn, nil
}

func (u *sellerUC) CreateVoucherSeller(ctx context.Context, userID string, requestBody body.CreateVoucherRequest) error {
	id, err := u.sellerRepo.GetShopIDByUserID(ctx, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return httperror.New(http.StatusBadRequest, response.UserNotHaveShop)
		}
		return err
	}

	shopID, err := uuid.Parse(id)
	if err != nil {
		return err
	}

	voucherShop := &model.Voucher{
		ShopID:             shopID,
		Code:               requestBody.Code,
		Quota:              requestBody.Quota,
		ActivedDate:        requestBody.ActiveDateTime,
		ExpiredDate:        requestBody.ExpiredDateTime,
		DiscountPercentage: &requestBody.DiscountPercentage,
		DiscountFixPrice:   &requestBody.DiscountFixPrice,
		MinProductPrice:    &requestBody.MinProductPrice,
		MaxDiscountPrice:   &requestBody.MaxDiscountPrice,
	}

	err = u.sellerRepo.CreateVoucherSeller(ctx, voucherShop)
	if err != nil {
		return err
	}

	return nil
}

func (u *sellerUC) UpdateVoucherSeller(ctx context.Context, userID string, requestBody body.UpdateVoucherRequest) error {
	shopID, err := u.sellerRepo.GetShopIDByUserID(ctx, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return httperror.New(http.StatusBadRequest, response.UserNotHaveShop)
		}
		return err
	}

	voucherIDShopID := &body.VoucherIDShopID{
		ShopID:    shopID,
		VoucherID: requestBody.VoucherID,
	}

	voucherShop, errVoucher := u.sellerRepo.GetAllVoucherSellerByIDAndShopID(ctx, voucherIDShopID)
	if errVoucher != nil {
		if errVoucher == sql.ErrNoRows {
			return httperror.New(http.StatusBadRequest, body.VoucherSellerNotFoundMessage)
		}

		return errVoucher
	}

	voucherShop.Quota = requestBody.Quota
	voucherShop.ActivedDate = requestBody.ActiveDateTime
	voucherShop.ExpiredDate = requestBody.ExpiredDateTime
	voucherShop.DiscountPercentage = &requestBody.DiscountPercentage
	voucherShop.DiscountFixPrice = &requestBody.DiscountFixPrice
	voucherShop.MinProductPrice = &requestBody.MinProductPrice
	voucherShop.MaxDiscountPrice = &requestBody.MaxDiscountPrice

	err = u.sellerRepo.UpdateVoucherSeller(ctx, voucherShop)
	if err != nil {
		return err
	}

	return nil
}

func (u *sellerUC) GetDetailVoucherSeller(ctx context.Context, voucherIDShopID *body.VoucherIDShopID) (*model.Voucher, error) {
	shopID, err := u.sellerRepo.GetShopIDByUserID(ctx, voucherIDShopID.UserID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, httperror.New(http.StatusBadRequest, response.UserNotHaveShop)
		}
		return nil, err
	}
	voucherIDShopID.ShopID = shopID

	voucherShop, errVoucher := u.sellerRepo.GetAllVoucherSellerByIDAndShopID(ctx, voucherIDShopID)
	if errVoucher != nil {
		if errVoucher == sql.ErrNoRows {
			return nil, httperror.New(http.StatusBadRequest, body.VoucherSellerNotFoundMessage)
		}

		return nil, errVoucher
	}

	return voucherShop, nil
}

func (u *sellerUC) DeleteVoucherSeller(ctx context.Context, voucherIDShopID *body.VoucherIDShopID) error {
	shopID, err := u.sellerRepo.GetShopIDByUserID(ctx, voucherIDShopID.UserID)
	if err != nil {
		if err == sql.ErrNoRows {
			return httperror.New(http.StatusBadRequest, response.UserNotHaveShop)
		}
		return err
	}
	voucherIDShopID.ShopID = shopID

	_, errVoucher := u.sellerRepo.GetAllVoucherSellerByIDAndShopID(ctx, voucherIDShopID)
	if errVoucher != nil {
		if errVoucher == sql.ErrNoRows {
			return httperror.New(http.StatusBadRequest, body.VoucherSellerNotFoundMessage)
		}

		return errVoucher
	}

	if err := u.sellerRepo.DeleteVoucherSeller(ctx, voucherIDShopID); err != nil {
		return err
	}

	return nil
}

func (u *sellerUC) CancelOrderStatus(ctx context.Context, userID string, requestBody body.CancelOrderStatus) error {
	shopIDFromUser, err := u.sellerRepo.GetShopIDByUser(ctx, userID)
	if err != nil {
		return err
	}

	order, err := u.sellerRepo.GetOrderByOrderID(ctx, requestBody.OrderID)
	if err != nil {
		return err
	}

	if shopIDFromUser != order.ShopID {
		return httperror.New(http.StatusUnauthorized, response.UnauthorizedMessage)
	}

	if order.OrderStatus != constant.OrderStatusWaitingForSeller {
		return httperror.New(http.StatusBadRequest, response.OrderNotWaitingForSeller)
	}

	if err := u.sellerRepo.CancelOrderStatus(ctx, requestBody); err != nil {
		return err
	}

	return nil
}

func (u *sellerUC) GetAllPromotionSeller(ctx context.Context, userID, promoStatusID string,
	pgn *pagination.Pagination) (*pagination.Pagination, error) {
	shopID, err := u.sellerRepo.GetShopIDByUserID(ctx, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, httperror.New(http.StatusBadRequest, response.UserNotHaveShop)
		}
		return nil, err
	}

	totalRows, err := u.sellerRepo.GetTotalPromotionSeller(ctx, shopID, promoStatusID)
	if err != nil {
		return nil, err
	}

	totalPages := int(math.Ceil(float64(totalRows) / float64(pgn.Limit)))
	pgn.TotalRows = totalRows
	pgn.TotalPages = totalPages

	ShopVouchers, err := u.sellerRepo.GetAllPromotionSeller(ctx, shopID, promoStatusID)
	if err != nil {
		return nil, err
	}

	pgn.Rows = ShopVouchers

	return pgn, nil
}

func (u *sellerUC) CreatePromotionSeller(ctx context.Context, userID string, requestBody body.CreatePromotionRequest) (int, error) {
	shopID, err := u.sellerRepo.GetShopIDByUserID(ctx, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return -1, httperror.New(http.StatusBadRequest, response.UserNotHaveShop)
		}
		return -1, err
	}

	data, errTx := u.txRepo.WithTransactionReturnData(func(tx postgre.Transaction) (interface{}, error) {
		countProduct := 0
		for _, productID := range requestBody.ProductIDs {
			shopProduct := &body.ShopProduct{ShopID: shopID, ProductID: productID}

			productPromo, errProductPromo := u.sellerRepo.GetProductPromotion(ctx, shopProduct)
			if errProductPromo != nil {
				if errProductPromo == sql.ErrNoRows {
					return nil, httperror.New(http.StatusBadRequest, response.ProductNotExistMessage)
				}
				return nil, errProductPromo
			}

			if productPromo.PromotionID != nil {
				return nil, httperror.New(http.StatusBadRequest, response.ProductAlreadyHasPromoMessage)
			}

			PID, err := uuid.Parse(productID)
			if err != nil {
				return nil, err
			}

			promotionShop := &model.Promotion{
				Name:               requestBody.Name,
				ProductID:          PID,
				DiscountPercentage: &requestBody.DiscountPercentage,
				DiscountFixPrice:   &requestBody.DiscountFixPrice,
				MinProductPrice:    &requestBody.MinProductPrice,
				MaxDiscountPrice:   &requestBody.MaxDiscountPrice,
				Quota:              requestBody.Quota,
				MaxQuantity:        requestBody.MaxQuantity,
				ActivedDate:        requestBody.ActiveDateTime,
				ExpiredDate:        requestBody.ExpiredDateTime,
			}

			err = u.sellerRepo.CreatePromotionSeller(ctx, tx, promotionShop)
			if err != nil {
				return nil, err
			}
			countProduct++
		}
		return countProduct, nil
	})

	if errTx != nil {
		return -1, errTx
	}

	return data.(int), nil
}

func (u *sellerUC) UpdatePromotionSeller(ctx context.Context, userID string, requestBody body.UpdatePromotionRequest) error {
	shopID, err := u.sellerRepo.GetShopIDByUserID(ctx, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return httperror.New(http.StatusBadRequest, response.UserNotHaveShop)
		}
		return err
	}

	shopProductPromo := &body.ShopProductPromo{
		ShopID:      shopID,
		ProductID:   requestBody.ProductID,
		PromotionID: requestBody.PromotionID,
	}

	promotionShop, errPromotion := u.sellerRepo.GetPromotionSellerDetailByID(ctx, shopProductPromo)
	if errPromotion != nil {
		if errPromotion == sql.ErrNoRows {
			return httperror.New(http.StatusBadRequest, body.PromotionSellerNotFoundMessage)
		}
		return errPromotion
	}
	timeNow := time.Now()

	if timeNow.After(promotionShop.ActivedDate) && timeNow.After(promotionShop.ExpiredDate) {
		return httperror.New(http.StatusBadRequest, body.PromotionExpairedMessage)
	}

	if timeNow.After(promotionShop.ActivedDate) && timeNow.Before(promotionShop.ExpiredDate) {
		if !promotionShop.ActivedDate.Equal(requestBody.ActiveDateTime) {
			return httperror.New(http.StatusBadRequest, body.PromotionCannotUpdateActivedMessage)
		}
	}

	if timeNow.Before(promotionShop.ActivedDate) && timeNow.Before(promotionShop.ExpiredDate) {
		if requestBody.ActiveDateTime.After(promotionShop.ActivedDate) {
			return httperror.New(http.StatusBadRequest, body.PromotionCannotUpdateActivedBeforeMessage)
		}
	}

	if requestBody.ExpiredDateTime.After(promotionShop.ExpiredDate) {
		return httperror.New(http.StatusBadRequest, body.PromotionCannotUpdateExpiredAfterMessage)
	}
	if requestBody.ExpiredDateTime.Before(timeNow) {
		return httperror.New(http.StatusBadRequest, body.PromotionCannotUpdateExpiredBefoferMessage)
	}

	promotion := &model.Promotion{
		ID:                 promotionShop.ID,
		Name:               requestBody.Name,
		MaxQuantity:        requestBody.MaxQuantity,
		DiscountPercentage: &requestBody.DiscountPercentage,
		DiscountFixPrice:   &requestBody.DiscountFixPrice,
		MinProductPrice:    &requestBody.MinProductPrice,
		MaxDiscountPrice:   &requestBody.MaxDiscountPrice,
		ActivedDate:        requestBody.ActiveDateTime,
		ExpiredDate:        requestBody.ExpiredDateTime,
	}

	err = u.sellerRepo.UpdatePromotionSeller(ctx, promotion)
	if err != nil {
		return err
	}

	return nil
}

func (u *sellerUC) GetDetailPromotionSellerByID(ctx context.Context,
	shopProductPromo *body.ShopProductPromo) (*body.PromotionDetailSeller, error) {
	shopID, err := u.sellerRepo.GetShopIDByUserID(ctx, shopProductPromo.UserID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, httperror.New(http.StatusBadRequest, response.UserNotHaveShop)
		}
		return nil, err
	}
	shopProductPromo.ShopID = shopID

	promotionShop, errPromotion := u.sellerRepo.GetDetailPromotionSellerByID(ctx, shopProductPromo)
	if errPromotion != nil {
		if errPromotion == sql.ErrNoRows {
			return nil, httperror.New(http.StatusBadRequest, body.PromotionSellerNotFoundMessage)
		}
		return nil, errPromotion
	}

	promotionShop = u.CalculateDiscountPromotionProduct(ctx, promotionShop)

	return promotionShop, nil
}

func (u *sellerUC) CalculateDiscountPromotionProduct(ctx context.Context, p *body.PromotionDetailSeller) *body.PromotionDetailSeller {
	var maxDiscountPrice float64
	var minProductPrice float64
	var discountPercentage float64
	var discountFixPrice float64
	var resultMinPriceDiscount float64
	var resultMaxPriceDiscount float64

	if p.MinProductPrice != nil {
		minProductPrice = *p.MinProductPrice
	}

	maxDiscountPrice = *p.MaxDiscountPrice
	if p.DiscountPercentage != nil {
		discountPercentage = *p.DiscountPercentage
		if p.MinPrice >= minProductPrice && discountPercentage > 0 {
			resultMinPriceDiscount = math.Min(maxDiscountPrice,
				p.MinPrice*(discountPercentage/100.00))
		}
		if p.MaxPrice >= minProductPrice && discountPercentage > 0 {
			resultMaxPriceDiscount = math.Min(maxDiscountPrice,
				p.MaxPrice*(discountPercentage/100.00))
		}
	}

	if p.DiscountFixPrice != nil {
		discountFixPrice = *p.DiscountFixPrice
		if p.MinPrice >= minProductPrice && discountFixPrice > 0 {
			resultMinPriceDiscount = math.Max(resultMinPriceDiscount, discountFixPrice)
			resultMinPriceDiscount = math.Min(resultMinPriceDiscount, maxDiscountPrice)
		}
		if p.MaxPrice >= minProductPrice && discountFixPrice > 0 {
			resultMaxPriceDiscount = math.Max(resultMaxPriceDiscount, discountFixPrice)
			resultMaxPriceDiscount = math.Min(resultMaxPriceDiscount, maxDiscountPrice)
		}
	}

	if resultMinPriceDiscount > 0 {
		p.ProductMinDiscountPrice = resultMinPriceDiscount
		p.ProductSubMinPrice = p.MinPrice - p.ProductMinDiscountPrice
	}

	if resultMaxPriceDiscount > 0 {
		p.ProductMaxDiscountPrice = resultMaxPriceDiscount
		p.ProductSubMaxPrice = p.MaxPrice - p.ProductMaxDiscountPrice
	}

	return p
}
