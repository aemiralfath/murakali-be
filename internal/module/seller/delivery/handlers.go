package delivery

import (
	"errors"
	"fmt"
	"murakali/config"
	"murakali/internal/module/seller"
	"murakali/internal/module/seller/delivery/body"
	"murakali/pkg/httperror"
	"murakali/pkg/logger"
	"murakali/pkg/pagination"
	"murakali/pkg/response"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type sellerHandlers struct {
	cfg      *config.Config
	sellerUC seller.UseCase
	logger   logger.Logger
}

func NewSellerHandlers(cfg *config.Config, sellerUC seller.UseCase, log logger.Logger) seller.Handlers {
	return &sellerHandlers{cfg: cfg, sellerUC: sellerUC, logger: log}
}

func (h *sellerHandlers) GetOrder(c *gin.Context) {
	userID, exist := c.Get("userID")
	if !exist {
		response.ErrorResponse(c.Writer, response.UnauthorizedMessage, http.StatusUnauthorized)
		return
	}

	userIDString := fmt.Sprintf("%v", userID)

	_, err := uuid.Parse(userIDString)
	if err != nil {
		response.ErrorResponse(c.Writer, response.BadRequestMessage, http.StatusBadRequest)
		return
	}

	pgn := &pagination.Pagination{}
	orderStatusID := c.DefaultQuery("order_status", "")
	voucherShopID := c.DefaultQuery("voucher_shop", "")
	h.ValidateQueryOrder(c, pgn)

	orders, err := h.sellerUC.GetOrder(c, userID.(string), orderStatusID, voucherShopID, pgn)
	if err != nil {
		var e *httperror.Error
		if !errors.As(err, &e) {
			h.logger.Errorf("HandlerSeller, Error: %s", err)
			response.ErrorResponse(c.Writer, response.InternalServerErrorMessage, http.StatusInternalServerError)
			return
		}

		response.ErrorResponse(c.Writer, e.Err.Error(), e.Status)
		return
	}

	response.SuccessResponse(c.Writer, orders, http.StatusOK)
}

func (h *sellerHandlers) ValidateQueryOrder(c *gin.Context, pgn *pagination.Pagination) {
	limit := strings.TrimSpace(c.Query("limit"))
	page := strings.TrimSpace(c.Query("page"))

	var limitFilter int
	var pageFilter int

	limitFilter, err := strconv.Atoi(limit)
	if err != nil || limitFilter < 1 {
		limitFilter = 10
	}

	pageFilter, err = strconv.Atoi(page)
	if err != nil || pageFilter < 1 {
		pageFilter = 1
	}

	pgn.Limit = limitFilter
	pgn.Page = pageFilter
}

func (h *sellerHandlers) ChangeOrderStatus(c *gin.Context) {
	var requestBody body.ChangeOrderStatusRequest

	if err := c.ShouldBind(&requestBody); err != nil {
		response.ErrorResponse(c.Writer, response.BadRequestMessage, http.StatusBadRequest)
		return
	}

	invalidFields, err := requestBody.Validate()
	if err != nil {
		response.ErrorResponseData(c.Writer, invalidFields, response.UnprocessableEntityMessage, http.StatusUnprocessableEntity)
		return
	}

	userID, exist := c.Get("userID")
	if !exist {
		response.ErrorResponse(c.Writer, response.UnauthorizedMessage, http.StatusUnauthorized)
		return
	}

	userIDString := fmt.Sprintf("%v", userID)

	_, err = uuid.Parse(userIDString)
	if err != nil {
		response.ErrorResponse(c.Writer, response.BadRequestMessage, http.StatusBadRequest)
		return
	}

	err = h.sellerUC.ChangeOrderStatus(c, fmt.Sprintf("%v", userID), requestBody)
	if err != nil {
		var e *httperror.Error
		if !errors.As(err, &e) {
			h.logger.Errorf("HandlerSeller, Error: %s", err)
			response.ErrorResponse(c.Writer, response.InternalServerErrorMessage, http.StatusInternalServerError)
			return
		}

		response.ErrorResponse(c.Writer, e.Err.Error(), e.Status)
		return
	}

	response.SuccessResponse(c.Writer, nil, http.StatusOK)
}

func (h *sellerHandlers) GetOrderByOrderID(c *gin.Context) {
	id := c.Param("order_id")
	orderID, err := uuid.Parse(id)
	if err != nil {
		response.ErrorResponse(c.Writer, response.BadRequestMessage, http.StatusBadRequest)
		return
	}

	data, err := h.sellerUC.GetOrderByOrderID(c, orderID.String())
	if err != nil {
		var e *httperror.Error
		if !errors.As(err, &e) {
			h.logger.Errorf("HandlerSeller, Error: %s", err)
			response.ErrorResponse(c.Writer, response.InternalServerErrorMessage, http.StatusInternalServerError)
			return
		}

		response.ErrorResponse(c.Writer, e.Err.Error(), e.Status)
		return
	}

	response.SuccessResponse(c.Writer, data, http.StatusOK)
}

func (h *sellerHandlers) GetSellerBySellerID(c *gin.Context) {
	id := c.Param("seller_id")
	sellerID, err := uuid.Parse(id)
	if err != nil {
		response.ErrorResponse(c.Writer, response.BadRequestMessage, http.StatusBadRequest)
		return
	}

	data, err := h.sellerUC.GetSellerBySellerID(c, sellerID.String())
	if err != nil {
		var e *httperror.Error
		if !errors.As(err, &e) {
			h.logger.Errorf("HandlerSeller, Error: %s", err)
			response.ErrorResponse(c.Writer, response.InternalServerErrorMessage, http.StatusInternalServerError)
			return
		}

		response.ErrorResponse(c.Writer, e.Err.Error(), e.Status)
		return
	}

	response.SuccessResponse(c.Writer, data, http.StatusOK)
}

func (h *sellerHandlers) GetSellerByUserID(c *gin.Context) {
	id := c.Param("user_id")
	userID, err := uuid.Parse(id)
	if err != nil {
		response.ErrorResponse(c.Writer, response.BadRequestMessage, http.StatusBadRequest)
		return
	}

	data, err := h.sellerUC.GetSellerByUserID(c, userID.String())
	if err != nil {
		var e *httperror.Error
		if !errors.As(err, &e) {
			h.logger.Errorf("HandlerSeller, Error: %s", err)
			response.ErrorResponse(c.Writer, response.InternalServerErrorMessage, http.StatusInternalServerError)
			return
		}

		response.ErrorResponse(c.Writer, e.Err.Error(), e.Status)
		return
	}

	response.SuccessResponse(c.Writer, data, http.StatusOK)
}

func (h *sellerHandlers) GetSellerDetailInformation(c *gin.Context) {
	userID, exist := c.Get("userID")
	if !exist {
		response.ErrorResponse(c.Writer, response.UnauthorizedMessage, http.StatusUnauthorized)
		return
	}

	data, err := h.sellerUC.GetSellerByUserID(c, userID.(string))
	if err != nil {
		var e *httperror.Error
		if !errors.As(err, &e) {
			h.logger.Errorf("HandlerSeller, Error: %s", err)
			response.ErrorResponse(c.Writer, response.InternalServerErrorMessage, http.StatusInternalServerError)
			return
		}

		response.ErrorResponse(c.Writer, e.Err.Error(), e.Status)
		return
	}

	var newData = body.SellerInformationResponse{
		ID:           data.ID,
		Name:         data.Name,
		TotalProduct: data.TotalProduct,
		TotalRating:  data.TotalRating,
		RatingAVG:    data.RatingAVG,
		PhotoURL:     data.PhotoURL,
		CreatedAt:    data.CreatedAt,
	}

	response.SuccessResponse(c.Writer, newData, http.StatusOK)
}

func (h *sellerHandlers) GetCategoryBySellerID(c *gin.Context) {
	id := c.Param("seller_id")
	sellerID, err := uuid.Parse(id)
	if err != nil {
		response.ErrorResponse(c.Writer, response.BadRequestMessage, http.StatusBadRequest)
		return
	}

	data, err := h.sellerUC.GetCategoryBySellerID(c, sellerID.String())
	if err != nil {
		var e *httperror.Error
		if !errors.As(err, &e) {
			h.logger.Errorf("HandlerSeller, Error: %s", err)
			response.ErrorResponse(c.Writer, response.InternalServerErrorMessage, http.StatusInternalServerError)
			return
		}

		response.ErrorResponse(c.Writer, e.Err.Error(), e.Status)
		return
	}

	response.SuccessResponse(c.Writer, data, http.StatusOK)
}

func (h *sellerHandlers) GetCourierSeller(c *gin.Context) {
	userID, exist := c.Get("userID")
	if !exist {
		response.ErrorResponse(c.Writer, response.UnauthorizedMessage, http.StatusUnauthorized)
		return
	}

	courierSeller, err := h.sellerUC.GetCourierSeller(c, userID.(string))
	if err != nil {
		var e *httperror.Error
		if !errors.As(err, &e) {
			h.logger.Errorf("HandlerSeller, Error: %s", err)
			response.ErrorResponse(c.Writer, response.InternalServerErrorMessage, http.StatusInternalServerError)
			return
		}

		response.ErrorResponse(c.Writer, e.Err.Error(), e.Status)
		return
	}

	response.SuccessResponse(c.Writer, courierSeller, http.StatusOK)
}

func (h *sellerHandlers) CreateCourierSeller(c *gin.Context) {
	userID, exist := c.Get("userID")
	if !exist {
		response.ErrorResponse(c.Writer, response.UnauthorizedMessage, http.StatusUnauthorized)
		return
	}

	var requestBody body.CourierSellerRequest
	if err := c.ShouldBind(&requestBody); err != nil {
		response.ErrorResponse(c.Writer, response.BadRequestMessage, http.StatusBadRequest)
		return
	}

	invalidFields, err := requestBody.Validate()
	if err != nil {
		response.ErrorResponseData(c.Writer, invalidFields, response.UnprocessableEntityMessage, http.StatusUnprocessableEntity)
		return
	}

	if err := h.sellerUC.CreateCourierSeller(c, userID.(string), requestBody.CourierID); err != nil {
		var e *httperror.Error
		if !errors.As(err, &e) {
			h.logger.Errorf("HandlerSeller, Error: %s", err)
			response.ErrorResponse(c.Writer, response.InternalServerErrorMessage, http.StatusInternalServerError)
			return
		}

		response.ErrorResponse(c.Writer, e.Err.Error(), e.Status)
		return
	}

	response.SuccessResponse(c.Writer, nil, http.StatusOK)
}

func (h *sellerHandlers) DeleteCourierSellerByID(c *gin.Context) {
	id := c.Param("id")
	sellerCourierID, err := uuid.Parse(id)
	if err != nil {
		response.ErrorResponse(c.Writer, response.BadRequestMessage, http.StatusBadRequest)
		return
	}

	if err := h.sellerUC.DeleteCourierSellerByID(c, sellerCourierID.String()); err != nil {
		var e *httperror.Error
		if !errors.As(err, &e) {
			h.logger.Errorf("HandlerSeller, Error: %s", err)
			response.ErrorResponse(c.Writer, response.InternalServerErrorMessage, http.StatusInternalServerError)
			return
		}
		response.ErrorResponse(c.Writer, e.Err.Error(), e.Status)
		return
	}

	response.SuccessResponse(c.Writer, nil, http.StatusOK)
}

func (h *sellerHandlers) UpdateResiNumberInOrderSeller(c *gin.Context) {
	userID, exist := c.Get("userID")
	if !exist {
		response.ErrorResponse(c.Writer, response.UnauthorizedMessage, http.StatusUnauthorized)
		return
	}

	id := c.Param("id")
	orderID, err := uuid.Parse(id)
	if err != nil {
		response.ErrorResponse(c.Writer, response.BadRequestMessage, http.StatusBadRequest)
		return
	}

	var requestBody body.UpdateNoResiOrderSellerRequest
	if err = c.ShouldBind(&requestBody); err != nil {
		response.ErrorResponse(c.Writer, response.BadRequestMessage, http.StatusBadRequest)
		return
	}

	invalidFields, err := requestBody.ValidateUpdateNoResi()
	if err != nil {
		response.ErrorResponseData(c.Writer, invalidFields, response.UnprocessableEntityMessage, http.StatusUnprocessableEntity)
		return
	}

	if err := h.sellerUC.UpdateResiNumberInOrderSeller(c, userID.(string), orderID.String(), requestBody); err != nil {
		var e *httperror.Error
		if !errors.As(err, &e) {
			h.logger.Errorf("HandlerSeller, Error: %s", err)
			response.ErrorResponse(c.Writer, response.InternalServerErrorMessage, http.StatusInternalServerError)
			return
		}

		response.ErrorResponse(c.Writer, e.Err.Error(), e.Status)
		return
	}

	response.SuccessResponse(c.Writer, nil, http.StatusOK)
}

func (h *sellerHandlers) GetAllVoucherSeller(c *gin.Context) {
	userID, exist := c.Get("userID")
	if !exist {
		response.ErrorResponse(c.Writer, response.UnauthorizedMessage, http.StatusUnauthorized)
		return
	}

	pgn := &pagination.Pagination{}
	h.ValidateQueryPagination(c, pgn)

	voucherStatusID := c.DefaultQuery("voucher_status", "")
	switch voucherStatusID {
	case "1":
		voucherStatusID = "1"
	case "2":
		voucherStatusID = "2"
	case "3":
		voucherStatusID = "3"
	case "4":
		voucherStatusID = "4"
	default:
		voucherStatusID = "1"
	}
	h.ValidateQueryPagination(c, pgn)

	shopVouchers, err := h.sellerUC.GetAllVoucherSeller(c, userID.(string), voucherStatusID, pgn)
	if err != nil {
		var e *httperror.Error
		if !errors.As(err, &e) {
			h.logger.Errorf("HandlerSeller, Error: %s", err)
			response.ErrorResponse(c.Writer, response.InternalServerErrorMessage, http.StatusInternalServerError)
			return
		}

		response.ErrorResponse(c.Writer, e.Err.Error(), e.Status)
		return
	}

	response.SuccessResponse(c.Writer, shopVouchers, http.StatusOK)
}

func (h *sellerHandlers) CreateVoucherSeller(c *gin.Context) {
	userID, exist := c.Get("userID")
	if !exist {
		response.ErrorResponse(c.Writer, response.UnauthorizedMessage, http.StatusUnauthorized)
		return
	}

	var requestBody body.CreateVoucherRequest
	if err := c.ShouldBind(&requestBody); err != nil {
		response.ErrorResponse(c.Writer, response.BadRequestMessage, http.StatusBadRequest)
		return
	}

	invalidFields, err := requestBody.Validate()
	if err != nil {
		response.ErrorResponseData(c.Writer, invalidFields, response.UnprocessableEntityMessage, http.StatusUnprocessableEntity)
		return
	}

	if err := h.sellerUC.CreateVoucherSeller(c, userID.(string), requestBody); err != nil {
		var e *httperror.Error
		if !errors.As(err, &e) {
			h.logger.Errorf("HandlerSeller, Error: %s", err)
			response.ErrorResponse(c.Writer, response.InternalServerErrorMessage, http.StatusInternalServerError)
			return
		}

		response.ErrorResponse(c.Writer, e.Err.Error(), e.Status)
		return
	}

	response.SuccessResponse(c.Writer, nil, http.StatusOK)
}

func (h *sellerHandlers) DeleteVoucherSeller(c *gin.Context) {
	userID, exist := c.Get("userID")
	if !exist {
		response.ErrorResponse(c.Writer, response.UnauthorizedMessage, http.StatusUnauthorized)
		return
	}

	id := c.Param("id")
	voucherShopID, err := uuid.Parse(id)
	if err != nil {
		response.ErrorResponse(c.Writer, response.BadRequestMessage, http.StatusBadRequest)
		return
	}
	voucherIDShopID := &body.VoucherIDShopID{
		UserID:    userID.(string),
		ShopID:    "",
		VoucherID: voucherShopID.String(),
	}
	if err := h.sellerUC.DeleteVoucherSeller(c, voucherIDShopID); err != nil {
		var e *httperror.Error
		if !errors.As(err, &e) {
			h.logger.Errorf("HandlerSeller, Error: %s", err)
			response.ErrorResponse(c.Writer, response.InternalServerErrorMessage, http.StatusInternalServerError)
			return
		}

		response.ErrorResponse(c.Writer, e.Err.Error(), e.Status)
		return
	}

	response.SuccessResponse(c.Writer, nil, http.StatusOK)
}

func (h *sellerHandlers) UpdateVoucherSeller(c *gin.Context) {
	userID, exist := c.Get("userID")
	if !exist {
		response.ErrorResponse(c.Writer, response.UnauthorizedMessage, http.StatusUnauthorized)
		return
	}

	var requestBody body.UpdateVoucherRequest
	if err := c.ShouldBind(&requestBody); err != nil {
		response.ErrorResponse(c.Writer, response.BadRequestMessage, http.StatusBadRequest)
		return
	}

	invalidFields, err := requestBody.Validate()
	if err != nil {
		response.ErrorResponseData(c.Writer, invalidFields, response.UnprocessableEntityMessage, http.StatusUnprocessableEntity)
		return
	}

	if err := h.sellerUC.UpdateVoucherSeller(c, userID.(string), requestBody); err != nil {
		var e *httperror.Error
		if !errors.As(err, &e) {
			h.logger.Errorf("HandlerSeller, Error: %s", err)
			response.ErrorResponse(c.Writer, response.InternalServerErrorMessage, http.StatusInternalServerError)
			return
		}

		response.ErrorResponse(c.Writer, e.Err.Error(), e.Status)
		return
	}

	response.SuccessResponse(c.Writer, nil, http.StatusOK)
}

func (h *sellerHandlers) UpdateOnDeliveryOrder(c *gin.Context) {
	if err := h.sellerUC.UpdateOnDeliveryOrder(c); err != nil {
		var e *httperror.Error
		if !errors.As(err, &e) {
			h.logger.Errorf("HandlerSeller, Error: %s", err)
			response.ErrorResponse(c.Writer, response.InternalServerErrorMessage, http.StatusInternalServerError)
			return
		}

		response.ErrorResponse(c.Writer, e.Err.Error(), e.Status)
		return
	}

	response.SuccessResponse(c.Writer, nil, http.StatusOK)
}

func (h *sellerHandlers) UpdateExpiredAtOrder(c *gin.Context) {
	if err := h.sellerUC.UpdateExpiredAtOrder(c); err != nil {
		var e *httperror.Error
		if !errors.As(err, &e) {
			h.logger.Errorf("HandlerSeller, Error: %s", err)
			response.ErrorResponse(c.Writer, response.InternalServerErrorMessage, http.StatusInternalServerError)
			return
		}

		response.ErrorResponse(c.Writer, e.Err.Error(), e.Status)
		return
	}

	response.SuccessResponse(c.Writer, nil, http.StatusOK)
}

func (h *sellerHandlers) DetailVoucherSeller(c *gin.Context) {
	userID, exist := c.Get("userID")
	if !exist {
		response.ErrorResponse(c.Writer, response.UnauthorizedMessage, http.StatusUnauthorized)
		return
	}

	id := c.Param("id")
	voucherShopID, err := uuid.Parse(id)
	if err != nil {
		response.ErrorResponse(c.Writer, response.BadRequestMessage, http.StatusBadRequest)
		return
	}
	voucherIDShopID := &body.VoucherIDShopID{
		UserID:    userID.(string),
		ShopID:    "",
		VoucherID: voucherShopID.String(),
	}
	voucherShop, err := h.sellerUC.GetDetailVoucherSeller(c, voucherIDShopID)
	if err != nil {
		var e *httperror.Error
		if !errors.As(err, &e) {
			h.logger.Errorf("HandlerSeller, Error: %s", err)
			response.ErrorResponse(c.Writer, response.InternalServerErrorMessage, http.StatusInternalServerError)
			return
		}

		response.ErrorResponse(c.Writer, e.Err.Error(), e.Status)
		return
	}

	response.SuccessResponse(c.Writer, voucherShop, http.StatusOK)
}

func (h *sellerHandlers) ValidateQueryPagination(c *gin.Context, pgn *pagination.Pagination) {
	limit := strings.TrimSpace(c.Query("limit"))
	page := strings.TrimSpace(c.Query("page"))

	var limitFilter int
	var pageFilter int

	limitFilter, err := strconv.Atoi(limit)
	if err != nil || limitFilter < 1 {
		limitFilter = 10
	}

	pageFilter, err = strconv.Atoi(page)
	if err != nil || pageFilter < 1 {
		pageFilter = 1
	}

	pgn.Limit = limitFilter
	pgn.Page = pageFilter
}

func (h *sellerHandlers) GetAllPromotionSeller(c *gin.Context) {
	userID, exist := c.Get("userID")
	if !exist {
		response.ErrorResponse(c.Writer, response.UnauthorizedMessage, http.StatusUnauthorized)
		return
	}

	pgn := &pagination.Pagination{}
	promoStatusID := c.DefaultQuery("promo_status", "")
	switch promoStatusID {
	case "1":
		promoStatusID = "1"
	case "2":
		promoStatusID = "2"
	case "3":
		promoStatusID = "3"
	case "4":
		promoStatusID = "4"

	default:
		promoStatusID = "1"
	}
	h.ValidateQueryPagination(c, pgn)

	promotionSeller, err := h.sellerUC.GetAllPromotionSeller(c, userID.(string), promoStatusID, pgn)
	if err != nil {
		var e *httperror.Error
		if !errors.As(err, &e) {
			h.logger.Errorf("HandlerSeller, Error: %s", err)
			response.ErrorResponse(c.Writer, response.InternalServerErrorMessage, http.StatusInternalServerError)
			return
		}

		response.ErrorResponse(c.Writer, e.Err.Error(), e.Status)
		return
	}

	response.SuccessResponse(c.Writer, promotionSeller, http.StatusOK)
}

func (h *sellerHandlers) CreatePromotionSeller(c *gin.Context) {
	userID, exist := c.Get("userID")
	if !exist {
		response.ErrorResponse(c.Writer, response.UnauthorizedMessage, http.StatusUnauthorized)
		return
	}

	var requestBody body.CreatePromotionRequest
	if err := c.ShouldBind(&requestBody); err != nil {
		response.ErrorResponse(c.Writer, response.BadRequestMessage, http.StatusBadRequest)
		return
	}

	invalidFields, err := requestBody.Validate()
	if err != nil {
		response.ErrorResponseData(c.Writer, invalidFields, response.UnprocessableEntityMessage, http.StatusUnprocessableEntity)
		return
	}

	rowEffected, err := h.sellerUC.CreatePromotionSeller(c, userID.(string), requestBody)
	if err != nil {
		var e *httperror.Error
		if !errors.As(err, &e) {
			h.logger.Errorf("HandlerSeller, Error: %s", err)
			response.ErrorResponse(c.Writer, response.InternalServerErrorMessage, http.StatusInternalServerError)
			return
		}

		response.ErrorResponse(c.Writer, e.Err.Error(), e.Status)
		return
	}

	response.SuccessResponse(c.Writer, rowEffected, http.StatusOK)
}
func (h *sellerHandlers) UpdatePromotionSeller(c *gin.Context) {
	userID, exist := c.Get("userID")
	if !exist {
		response.ErrorResponse(c.Writer, response.UnauthorizedMessage, http.StatusUnauthorized)
		return
	}

	var requestBody body.UpdatePromotionRequest
	if err := c.ShouldBind(&requestBody); err != nil {
		response.ErrorResponse(c.Writer, response.BadRequestMessage, http.StatusBadRequest)
		return
	}

	invalidFields, err := requestBody.Validate()
	if err != nil {
		response.ErrorResponseData(c.Writer, invalidFields, response.UnprocessableEntityMessage, http.StatusUnprocessableEntity)
		return
	}

	if err := h.sellerUC.UpdatePromotionSeller(c, userID.(string), requestBody); err != nil {
		var e *httperror.Error
		if !errors.As(err, &e) {
			h.logger.Errorf("HandlerSeller, Error: %s", err)
			response.ErrorResponse(c.Writer, response.InternalServerErrorMessage, http.StatusInternalServerError)
			return
		}

		response.ErrorResponse(c.Writer, e.Err.Error(), e.Status)
		return
	}

	response.SuccessResponse(c.Writer, nil, http.StatusOK)
}

func (h *sellerHandlers) GetDetailPromotionSellerByID(c *gin.Context) {
	userID, exist := c.Get("userID")
	if !exist {
		response.ErrorResponse(c.Writer, response.UnauthorizedMessage, http.StatusUnauthorized)
		return
	}

	id := c.Param("id")
	promotionShopID, err := uuid.Parse(id)
	if err != nil {
		response.ErrorResponse(c.Writer, response.BadRequestMessage, http.StatusBadRequest)
		return
	}
	shopProductPromo := &body.ShopProductPromo{
		UserID:      userID.(string),
		ShopID:      "",
		PromotionID: promotionShopID.String(),
	}
	promotionShop, err := h.sellerUC.GetDetailPromotionSellerByID(c, shopProductPromo)
	if err != nil {
		var e *httperror.Error
		if !errors.As(err, &e) {
			h.logger.Errorf("HandlerSeller, Error: %s", err)
			response.ErrorResponse(c.Writer, response.InternalServerErrorMessage, http.StatusInternalServerError)
			return
		}

		response.ErrorResponse(c.Writer, e.Err.Error(), e.Status)
		return
	}

	response.SuccessResponse(c.Writer, promotionShop, http.StatusOK)
}

func (h *sellerHandlers) GetProductWithoutPromotionSeller(c *gin.Context) {
	userID, exist := c.Get("userID")
	if !exist {
		response.ErrorResponse(c.Writer, response.UnauthorizedMessage, http.StatusUnauthorized)
		return
	}

	pgn := &pagination.Pagination{}
	h.ValidateQueryPagination(c, pgn)

	productWithoutPromotion, err := h.sellerUC.GetProductWithoutPromotionSeller(c, userID.(string), pgn)
	if err != nil {
		var e *httperror.Error
		if !errors.As(err, &e) {
			h.logger.Errorf("HandlerSeller, Error: %s", err)
			response.ErrorResponse(c.Writer, response.InternalServerErrorMessage, http.StatusInternalServerError)
			return
		}

		response.ErrorResponse(c.Writer, e.Err.Error(), e.Status)
		return
	}

	response.SuccessResponse(c.Writer, productWithoutPromotion, http.StatusOK)
}