package delivery

import (
	"errors"
	"fmt"
	"murakali/config"
	"murakali/internal/constant"
	"murakali/internal/module/user"
	"murakali/internal/module/user/delivery/body"
	"murakali/internal/util"
	"murakali/pkg/httperror"
	"murakali/pkg/jwt"
	"murakali/pkg/logger"
	"murakali/pkg/pagination"
	"murakali/pkg/response"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type userHandlers struct {
	cfg    *config.Config
	userUC user.UseCase
	logger logger.Logger
}

func NewUserHandlers(cfg *config.Config, userUC user.UseCase, log logger.Logger) user.Handlers {
	return &userHandlers{cfg: cfg, userUC: userUC, logger: log}
}

func (h *userHandlers) RegisterMerchant(c *gin.Context) {
	userID, exist := c.Get("userID")
	if !exist {
		response.ErrorResponse(c.Writer, response.UnauthorizedMessage, http.StatusUnauthorized)
		return
	}

	var requestBody body.RegisterMerchant
	if err := c.ShouldBind(&requestBody); err != nil {
		response.ErrorResponse(c.Writer, response.BadRequestMessage, http.StatusBadRequest)
		return
	}

	invalidFields, err := requestBody.Validate()
	if err != nil {
		response.ErrorResponseData(c.Writer, invalidFields, response.UnprocessableEntityMessage, http.StatusUnprocessableEntity)
		return
	}
	if err := h.userUC.RegisterMerchant(c, userID.(string), requestBody.ShopName); err != nil {
		var e *httperror.Error
		if !errors.As(err, &e) {
			h.logger.Errorf("HandlerUser, Error: %s", err)
			response.ErrorResponse(c.Writer, response.InternalServerErrorMessage, http.StatusInternalServerError)
			return
		}

		response.ErrorResponse(c.Writer, e.Err.Error(), e.Status)
		return
	}

	response.SuccessResponse(c.Writer, nil, http.StatusOK)
}

func (h *userHandlers) GetWallet(c *gin.Context) {
	userID, exist := c.Get("userID")
	if !exist {
		response.ErrorResponse(c.Writer, response.UnauthorizedMessage, http.StatusUnauthorized)
		return
	}

	wallet, err := h.userUC.GetWallet(c, userID.(string))
	if err != nil {
		var e *httperror.Error
		if !errors.As(err, &e) {
			h.logger.Errorf("HandlerUser, Error: %s", err)
			response.ErrorResponse(c.Writer, response.InternalServerErrorMessage, http.StatusInternalServerError)
			return
		}

		response.ErrorResponse(c.Writer, e.Err.Error(), e.Status)
		return
	}

	wallet.PIN = ""
	response.SuccessResponse(c.Writer, wallet, http.StatusOK)
}

func (h *userHandlers) GetWalletHistory(c *gin.Context) {
	userID, exist := c.Get("userID")
	if !exist {
		response.ErrorResponse(c.Writer, response.UnauthorizedMessage, http.StatusUnauthorized)
		return
	}

	pgn := h.ValidateQuery(c)

	walletHistory, err := h.userUC.GetWalletHistory(c, userID.(string), pgn)
	if err != nil {
		var e *httperror.Error
		if !errors.As(err, &e) {
			h.logger.Errorf("HandlerUser, Error: %s", err)
			response.ErrorResponse(c.Writer, response.InternalServerErrorMessage, http.StatusInternalServerError)
			return
		}

		response.ErrorResponse(c.Writer, e.Err.Error(), e.Status)
		return
	}

	response.SuccessResponse(c.Writer, walletHistory, http.StatusOK)
}

func (h *userHandlers) ValidateQuery(c *gin.Context) *pagination.Pagination {
	limit := strings.TrimSpace(c.Query("limit"))
	page := strings.TrimSpace(c.Query("page"))
	sort := strings.TrimSpace(c.Query("sort"))

	var limitFilter int
	var pageFilter int
	sortFilter := "DESC"

	limitFilter, err := strconv.Atoi(limit)
	if err != nil || limitFilter < 1 {
		limitFilter = 18
	}

	pageFilter, err = strconv.Atoi(page)
	if err != nil || pageFilter < 1 {
		pageFilter = 1
	}

	if sort == "ASC" || sort == "asc" {
		sortFilter = "ASC"
	}

	pgn := &pagination.Pagination{
		Limit: limitFilter,
		Page:  pageFilter,
		Sort:  "created_at " + sortFilter,
	}

	return pgn
}

func (h *userHandlers) TopUpWallet(c *gin.Context) {
	userID, exist := c.Get("userID")
	if !exist {
		response.ErrorResponse(c.Writer, response.UnauthorizedMessage, http.StatusUnauthorized)
		return
	}

	var requestBody body.TopUpWalletRequest
	if err := c.ShouldBind(&requestBody); err != nil {
		response.ErrorResponse(c.Writer, response.BadRequestMessage, http.StatusBadRequest)
		return
	}

	invalidFields, err := requestBody.Validate()
	if err != nil {
		response.ErrorResponseData(c.Writer, invalidFields, response.UnprocessableEntityMessage, http.StatusUnprocessableEntity)
		return
	}

	transactionID, err := h.userUC.TopUpWallet(c, userID.(string), requestBody)
	if err != nil {
		var e *httperror.Error
		if !errors.As(err, &e) {
			h.logger.Errorf("HandlerUser, Error: %s", err)
			response.ErrorResponse(c.Writer, response.InternalServerErrorMessage, http.StatusInternalServerError)
			return
		}

		response.ErrorResponse(c.Writer, e.Err.Error(), e.Status)
		return
	}

	response.SuccessResponse(c.Writer, body.TopUpWalletResponse{TransactionID: transactionID}, http.StatusOK)
}

func (h *userHandlers) ActivateWallet(c *gin.Context) {
	userID, exist := c.Get("userID")
	if !exist {
		response.ErrorResponse(c.Writer, response.UnauthorizedMessage, http.StatusUnauthorized)
		return
	}

	var requestBody body.ActivateWalletRequest
	if err := c.ShouldBind(&requestBody); err != nil {
		response.ErrorResponse(c.Writer, response.BadRequestMessage, http.StatusBadRequest)
		return
	}

	invalidFields, err := requestBody.Validate()
	if err != nil {
		response.ErrorResponseData(c.Writer, invalidFields, response.UnprocessableEntityMessage, http.StatusUnprocessableEntity)
		return
	}

	if err := h.userUC.ActivateWallet(c, userID.(string), requestBody.Pin); err != nil {
		var e *httperror.Error
		if !errors.As(err, &e) {
			h.logger.Errorf("HandlerUser, Error: %s", err)
			response.ErrorResponse(c.Writer, response.InternalServerErrorMessage, http.StatusInternalServerError)
			return
		}

		response.ErrorResponse(c.Writer, e.Err.Error(), e.Status)
		return
	}

	response.SuccessResponse(c.Writer, nil, http.StatusOK)
}

func (h *userHandlers) DeleteAddressByID(c *gin.Context) {
	userID, exist := c.Get("userID")
	if !exist {
		response.ErrorResponse(c.Writer, response.UnauthorizedMessage, http.StatusUnauthorized)
		return
	}

	id := c.Param("id")
	addressID, err := uuid.Parse(id)
	if err != nil {
		response.ErrorResponse(c.Writer, response.BadRequestMessage, http.StatusBadRequest)
		return
	}

	if err := h.userUC.DeleteAddressByID(c, userID.(string), addressID.String()); err != nil {
		var e *httperror.Error
		if !errors.As(err, &e) {
			h.logger.Errorf("HandlerUser, Error: %s", err)
			response.ErrorResponse(c.Writer, response.InternalServerErrorMessage, http.StatusInternalServerError)
			return
		}

		response.ErrorResponse(c.Writer, e.Err.Error(), e.Status)
		return
	}

	response.SuccessResponse(c.Writer, nil, http.StatusOK)
}

func (h *userHandlers) GetAddressByID(c *gin.Context) {
	userID, exist := c.Get("userID")
	if !exist {
		response.ErrorResponse(c.Writer, response.UnauthorizedMessage, http.StatusUnauthorized)
		return
	}

	id := c.Param("id")
	addressID, err := uuid.Parse(id)
	if err != nil {
		response.ErrorResponse(c.Writer, response.BadRequestMessage, http.StatusBadRequest)
		return
	}

	address, err := h.userUC.GetAddressByID(c, userID.(string), addressID.String())
	if err != nil {
		var e *httperror.Error
		if !errors.As(err, &e) {
			h.logger.Errorf("HandlerUser, Error: %s", err)
			response.ErrorResponse(c.Writer, response.InternalServerErrorMessage, http.StatusInternalServerError)
			return
		}

		response.ErrorResponse(c.Writer, e.Err.Error(), e.Status)
		return
	}

	response.SuccessResponse(c.Writer, address, http.StatusOK)
}

func (h *userHandlers) CreateAddress(c *gin.Context) {
	userID, exist := c.Get("userID")
	if !exist {
		response.ErrorResponse(c.Writer, response.UnauthorizedMessage, http.StatusUnauthorized)
		return
	}

	var requestBody body.CreateAddressRequest
	if err := c.ShouldBind(&requestBody); err != nil {
		response.ErrorResponse(c.Writer, response.BadRequestMessage, http.StatusBadRequest)
		return
	}

	invalidFields, err := requestBody.Validate()
	if err != nil {
		response.ErrorResponseData(c.Writer, invalidFields, response.UnprocessableEntityMessage, http.StatusUnprocessableEntity)
		return
	}

	if err := h.userUC.CreateAddress(c, userID.(string), requestBody); err != nil {
		var e *httperror.Error
		if !errors.As(err, &e) {
			h.logger.Errorf("HandlerUser, Error: %s", err)
			response.ErrorResponse(c.Writer, response.InternalServerErrorMessage, http.StatusInternalServerError)
			return
		}

		response.ErrorResponse(c.Writer, e.Err.Error(), e.Status)
		return
	}

	response.SuccessResponse(c.Writer, nil, http.StatusOK)
}

func (h *userHandlers) UpdateAddressByID(c *gin.Context) {
	userID, exist := c.Get("userID")
	if !exist {
		response.ErrorResponse(c.Writer, response.UnauthorizedMessage, http.StatusUnauthorized)
		return
	}

	id := c.Param("id")
	addressID, err := uuid.Parse(id)
	if err != nil {
		response.ErrorResponse(c.Writer, response.BadRequestMessage, http.StatusBadRequest)
		return
	}

	var requestBody body.UpdateAddressRequest
	if c.ShouldBind(&requestBody) != nil {
		response.ErrorResponse(c.Writer, response.BadRequestMessage, http.StatusBadRequest)
		return
	}

	invalidFields, err := requestBody.Validate()
	if err != nil {
		response.ErrorResponseData(c.Writer, invalidFields, response.UnprocessableEntityMessage, http.StatusUnprocessableEntity)
		return
	}

	if err := h.userUC.UpdateAddressByID(c, userID.(string), addressID.String(), requestBody); err != nil {
		var e *httperror.Error
		if !errors.As(err, &e) {
			h.logger.Errorf("HandlerUser, Error: %s", err)
			response.ErrorResponse(c.Writer, response.InternalServerErrorMessage, http.StatusInternalServerError)
			return
		}

		response.ErrorResponse(c.Writer, e.Err.Error(), e.Status)
		return
	}

	response.SuccessResponse(c.Writer, nil, http.StatusOK)
}

func (h *userHandlers) GetAddress(c *gin.Context) {
	userID, exist := c.Get("userID")
	if !exist {
		response.ErrorResponse(c.Writer, response.UnauthorizedMessage, http.StatusUnauthorized)
		return
	}

	pgn := &pagination.Pagination{}
	queryRequest := h.ValidateQueryAddress(c, pgn)

	addresses, err := h.userUC.GetAddress(c, userID.(string), pgn, queryRequest)
	if err != nil {
		var e *httperror.Error
		if !errors.As(err, &e) {
			h.logger.Errorf("HandlerUser, Error: %s", err)
			response.ErrorResponse(c.Writer, response.InternalServerErrorMessage, http.StatusInternalServerError)
			return
		}

		response.ErrorResponse(c.Writer, e.Err.Error(), e.Status)
		return
	}

	response.SuccessResponse(c.Writer, addresses, http.StatusOK)
}

func (h *userHandlers) ValidateQueryAddress(c *gin.Context, pgn *pagination.Pagination) *body.GetAddressQueryRequest {
	queryRequest := &body.GetAddressQueryRequest{}

	name := strings.TrimSpace(c.Query("name"))
	isDefault := strings.TrimSpace(c.Query("is_default"))
	isShopDefault := strings.TrimSpace(c.Query("is_shop_default"))

	sort := strings.TrimSpace(c.Query("sort"))
	sortBy := strings.TrimSpace(c.Query("sortBy"))
	limit := strings.TrimSpace(c.Query("limit"))
	page := strings.TrimSpace(c.Query("page"))

	var sortFilter string
	var sortByFilter string
	var limitFilter int
	var pageFilter int
	var isDefaultFilter bool
	var isShopDefaultFilter bool

	switch isDefault {
	case constant.AddressDefault:
		isDefaultFilter = true
	default:
		isDefaultFilter = false
	}

	switch isShopDefault {
	case constant.AddressDefault:
		isShopDefaultFilter = true
	default:
		isShopDefaultFilter = false
	}

	switch sort {
	case "asc":
		sortFilter = sort
	default:
		sortFilter = "desc"
	}

	switch sortBy {
	case "province":
		sortByFilter = sortBy
	default:
		sortByFilter = "created_at"
	}

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
	pgn.Sort = fmt.Sprintf("%s %s", sortByFilter, sortFilter)

	queryRequest.Name = name
	queryRequest.IsDefaultBool = isDefaultFilter
	queryRequest.IsShopDefaultBool = isShopDefaultFilter

	return queryRequest
}

func (h *userHandlers) GetOrder(c *gin.Context) {
	userID, exist := c.Get("userID")
	if !exist {
		response.ErrorResponse(c.Writer, response.UnauthorizedMessage, http.StatusUnauthorized)
		return
	}

	pgn := &pagination.Pagination{}

	h.ValidateQueryOrder(c, pgn)

	orders, err := h.userUC.GetOrder(c, userID.(string), pgn)
	if err != nil {
		var e *httperror.Error
		if !errors.As(err, &e) {
			h.logger.Errorf("HandlerUser, Error: %s", err)
			response.ErrorResponse(c.Writer, response.InternalServerErrorMessage, http.StatusInternalServerError)
			return
		}

		response.ErrorResponse(c.Writer, e.Err.Error(), e.Status)
		return
	}

	response.SuccessResponse(c.Writer, orders, http.StatusOK)
}

func (h *userHandlers) ValidateQueryOrder(c *gin.Context, pgn *pagination.Pagination) {
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

func (h *userHandlers) EditUser(c *gin.Context) {
	var requestBody body.EditUserRequest
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

	_, err = h.userUC.EditUser(c, fmt.Sprintf("%v", userID), requestBody)
	if err != nil {
		var e *httperror.Error
		if !errors.As(err, &e) {
			h.logger.Errorf("HandlerUser, Error: %s", err)
			response.ErrorResponse(c.Writer, response.InternalServerErrorMessage, http.StatusInternalServerError)
			return
		}

		response.ErrorResponse(c.Writer, e.Err.Error(), e.Status)
		return
	}

	response.SuccessResponse(c.Writer, nil, http.StatusOK)
}

func (h *userHandlers) EditEmail(c *gin.Context) {
	var requestBody body.EditEmailRequest
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

	_, err = h.userUC.EditEmail(c, fmt.Sprintf("%v", userID), requestBody)
	if err != nil {
		var e *httperror.Error
		if !errors.As(err, &e) {
			h.logger.Errorf("HandlerUser, Error: %s", err)
			response.ErrorResponse(c.Writer, response.InternalServerErrorMessage, http.StatusInternalServerError)
			return
		}

		response.ErrorResponse(c.Writer, e.Err.Error(), e.Status)
		return
	}

	response.SuccessResponse(c.Writer, nil, http.StatusOK)
}

func (h *userHandlers) EditEmailUser(c *gin.Context) {
	var requestParam body.EditEmailUserRequest
	if err := c.ShouldBind(&requestParam); err != nil {
		response.ErrorResponse(c.Writer, response.BadRequestMessage, http.StatusBadRequest)
		return
	}

	invalidFields, err := requestParam.Validate()
	if err != nil {
		response.ErrorResponseData(c.Writer, invalidFields, response.UnprocessableEntityMessage, http.StatusUnprocessableEntity)
		return
	}

	userID, exist := c.Get("userID")
	if !exist {
		response.ErrorResponse(c.Writer, response.UnauthorizedMessage, http.StatusUnauthorized)
		return
	}

	_, err = h.userUC.EditEmailUser(c, fmt.Sprintf("%v", userID), requestParam)
	if err != nil {
		var e *httperror.Error
		if !errors.As(err, &e) {
			h.logger.Errorf("HandlerUser, Error: %s", err)
			response.ErrorResponse(c.Writer, response.InternalServerErrorMessage, http.StatusInternalServerError)
			return
		}

		response.ErrorResponse(c.Writer, e.Err.Error(), e.Status)
		return
	}

	response.SuccessResponse(c.Writer, nil, http.StatusOK)
}

func (h *userHandlers) GetSealabsPay(c *gin.Context) {
	userid, exist := c.Get("userID")

	if !exist {
		response.ErrorResponse(c.Writer, response.UnauthorizedMessage, http.StatusUnauthorized)
		return
	}

	result, err := h.userUC.GetSealabsPay(c, userid.(string))
	if err != nil {
		var e *httperror.Error
		if !errors.As(err, &e) {
			h.logger.Errorf("HandlerUser, Error: %s", err)
			response.ErrorResponse(c.Writer, response.InternalServerErrorMessage, http.StatusInternalServerError)
			return
		}

		response.ErrorResponse(c.Writer, e.Err.Error(), e.Status)
		return
	}

	response.SuccessResponse(c.Writer, result, http.StatusOK)
}

func (h *userHandlers) AddSealabsPay(c *gin.Context) {
	userid, exist := c.Get("userID")
	if !exist {
		response.ErrorResponse(c.Writer, response.UnauthorizedMessage, http.StatusUnauthorized)
		return
	}

	var requestBody body.AddSealabsPayRequest
	if err := c.ShouldBind(&requestBody); err != nil {
		response.ErrorResponse(c.Writer, response.BadRequestMessage, http.StatusBadRequest)
		return
	}

	invalidFields, err := requestBody.Validate()
	if err != nil {
		response.ErrorResponseData(c.Writer, invalidFields, response.UnprocessableEntityMessage, http.StatusUnprocessableEntity)
		return
	}

	if err := h.userUC.AddSealabsPay(c, requestBody, userid.(string)); err != nil {
		var e *httperror.Error
		if !errors.As(err, &e) {
			h.logger.Errorf("HandlerUser, Error: %s", err)
			response.ErrorResponse(c.Writer, response.InternalServerErrorMessage, http.StatusInternalServerError)
			return
		}

		response.ErrorResponse(c.Writer, e.Err.Error(), e.Status)
		return
	}
	response.SuccessResponse(c.Writer, nil, http.StatusOK)
}

func (h *userHandlers) PatchSealabsPay(c *gin.Context) {
	cardNumber := c.Param("cardNumber")

	var requestBody body.SlpCardRequest
	requestBody.CardNumber = cardNumber

	userid, exist := c.Get("userID")
	if !exist {
		response.ErrorResponse(c.Writer, response.UnauthorizedMessage, http.StatusUnauthorized)
		return
	}

	invalidFields, err := requestBody.Validate()
	if err != nil {
		response.ErrorResponseData(c.Writer, invalidFields, response.UnprocessableEntityMessage, http.StatusUnprocessableEntity)
		return
	}

	if err := h.userUC.PatchSealabsPay(c, requestBody.CardNumber, userid.(string)); err != nil {
		var e *httperror.Error
		if !errors.As(err, &e) {
			h.logger.Errorf("HandlerUser, Error: %s", err)
			response.ErrorResponse(c.Writer, response.InternalServerErrorMessage, http.StatusInternalServerError)
			return
		}

		response.ErrorResponse(c.Writer, e.Err.Error(), e.Status)
		return
	}

	response.SuccessResponse(c.Writer, nil, http.StatusOK)
}

func (h *userHandlers) DeleteSealabsPay(c *gin.Context) {
	cardNumber := c.Param("cardNumber")

	var requestBody body.SlpCardRequest
	requestBody.CardNumber = cardNumber

	invalidFields, err := requestBody.Validate()
	if err != nil {
		response.ErrorResponseData(c.Writer, invalidFields, response.UnprocessableEntityMessage, http.StatusUnprocessableEntity)
		return
	}

	if err := h.userUC.DeleteSealabsPay(c, requestBody.CardNumber); err != nil {
		var e *httperror.Error
		if !errors.As(err, &e) {
			h.logger.Errorf("HandlerUser, Error: %s", err)
			response.ErrorResponse(c.Writer, response.InternalServerErrorMessage, http.StatusInternalServerError)
			return
		}

		response.ErrorResponse(c.Writer, e.Err.Error(), e.Status)
		return
	}

	response.SuccessResponse(c.Writer, nil, http.StatusOK)
}

func (h *userHandlers) GetUserProfile(c *gin.Context) {
	userID, exist := c.Get("userID")
	if !exist {
		response.ErrorResponse(c.Writer, response.UnauthorizedMessage, http.StatusUnauthorized)
		return
	}

	profile, err := h.userUC.GetUserProfile(c, userID.(string))
	if err != nil {
		var e *httperror.Error
		if !errors.As(err, &e) {
			h.logger.Errorf("HandlerUser, Error: %s", err)
			response.ErrorResponse(c.Writer, response.InternalServerErrorMessage, http.StatusInternalServerError)
			return
		}

		response.ErrorResponse(c.Writer, e.Err.Error(), e.Status)
		return
	}

	response.SuccessResponse(c.Writer, profile, http.StatusOK)
}

func (h *userHandlers) UploadProfilePicture(c *gin.Context) {
	type Sizer interface {
		Size() int64
	}

	var imgURL string
	var img body.ImageRequest

	userID, exist := c.Get("userID")
	if !exist {
		response.ErrorResponse(c.Writer, response.UnauthorizedMessage, http.StatusUnauthorized)
		return
	}

	err := c.ShouldBind(&img)
	if err != nil {
		response.ErrorResponse(c.Writer, response.InternalServerErrorMessage, http.StatusInternalServerError)
		return
	}
	data, _, _ := c.Request.FormFile("Img")

	if data.(Sizer).Size() > constant.ImgMaxSize {
		response.ErrorResponse(c.Writer, response.PictureSizeTooBig, http.StatusInternalServerError)
		return
	}

	if data == nil {
		response.ErrorResponse(c.Writer, response.InternalServerErrorMessage, http.StatusInternalServerError)
		return
	}
	imgURL = util.UploadImageToCloudinary(c, h.cfg, data)

	err = h.userUC.UploadProfilePicture(c, imgURL, userID.(string))

	if err != nil {
		var e *httperror.Error
		if !errors.As(err, &e) {
			h.logger.Errorf("HandlerUser, Error: %s", err)
			response.ErrorResponse(c.Writer, response.InternalServerErrorMessage, http.StatusInternalServerError)
			return
		}

		response.ErrorResponse(c.Writer, e.Err.Error(), e.Status)
		return
	}
	response.SuccessResponse(c.Writer, nil, http.StatusOK)
}

func (h *userHandlers) VerifyPasswordChange(c *gin.Context) {
	userID, exist := c.Get("userID")
	if !exist {
		response.ErrorResponse(c.Writer, response.UnauthorizedMessage, http.StatusUnauthorized)
		return
	}

	err := h.userUC.VerifyPasswordChange(c, userID.(string))
	if err != nil {
		var e *httperror.Error
		if !errors.As(err, &e) {
			h.logger.Errorf("HandlerUser, Error: %s", err)
			response.ErrorResponse(c.Writer, response.InternalServerErrorMessage, http.StatusInternalServerError)
			return
		}

		response.ErrorResponse(c.Writer, e.Err.Error(), e.Status)
		return
	}

	response.SuccessResponse(c.Writer, nil, http.StatusOK)
}

func (h *userHandlers) VerifyOTP(c *gin.Context) {
	userID, exist := c.Get("userID")
	if !exist {
		response.ErrorResponse(c.Writer, response.UnauthorizedMessage, http.StatusUnauthorized)
		return
	}

	var requestBody body.VerifyOTPRequest
	if err := c.ShouldBind(&requestBody); err != nil {
		response.ErrorResponse(c.Writer, response.BadRequestMessage, http.StatusBadRequest)
		return
	}

	invalidFields, err := requestBody.Validate()
	if err != nil {
		response.ErrorResponseData(c.Writer, invalidFields, response.UnprocessableEntityMessage, http.StatusUnprocessableEntity)
		return
	}

	changePasswordToken, err := h.userUC.VerifyOTP(c, requestBody, userID.(string))
	if err != nil {
		var e *httperror.Error
		if !errors.As(err, &e) {
			h.logger.Errorf("HandlerUser, Error: %s", err)
			response.ErrorResponse(c.Writer, response.InternalServerErrorMessage, http.StatusInternalServerError)
			return
		}

		response.ErrorResponse(c.Writer, e.Err.Error(), e.Status)
		return
	}

	c.SetSameSite(http.SameSiteNoneMode)
	c.SetCookie(constant.ChangePasswordTokenCookie, changePasswordToken, h.cfg.JWT.RefreshExpMin*60, "/", h.cfg.Server.Domain, true, true)
	response.SuccessResponse(c.Writer, nil, http.StatusOK)
}

func (h *userHandlers) ChangePassword(c *gin.Context) {
	changePasswordToken, err := c.Cookie(constant.ChangePasswordTokenCookie)
	if err != nil {
		response.ErrorResponse(c.Writer, response.ForbiddenMessage, http.StatusForbidden)
		return
	}

	claims, err := jwt.ExtractJWT(changePasswordToken, h.cfg.JWT.JwtSecretKey)
	if err != nil {
		response.ErrorResponse(c.Writer, response.ForbiddenMessage, http.StatusForbidden)
		return
	}

	var requestBody body.ChangePasswordRequest
	if c.ShouldBind(&requestBody) != nil {
		response.ErrorResponse(c.Writer, response.BadRequestMessage, http.StatusBadRequest)
		return
	}

	invalidFields, err := requestBody.Validate()
	if err != nil {
		response.ErrorResponseData(c.Writer, invalidFields, response.UnprocessableEntityMessage, http.StatusUnprocessableEntity)
		return
	}

	if err := h.userUC.ChangePassword(c, claims["id"].(string), requestBody.NewPassword); err != nil {
		var e *httperror.Error
		if !errors.As(err, &e) {
			h.logger.Errorf("HandlerUser, Error: %s", err)
			response.ErrorResponse(c.Writer, response.InternalServerErrorMessage, http.StatusInternalServerError)
			return
		}

		response.ErrorResponse(c.Writer, e.Err.Error(), e.Status)
		return
	}

	c.SetSameSite(http.SameSiteNoneMode)
	c.SetCookie(constant.ChangePasswordTokenCookie, "", -1, "/", h.cfg.Server.Domain, true, true)
	response.SuccessResponse(c.Writer, nil, http.StatusOK)
}

func (h *userHandlers) WalletStepUp(c *gin.Context) {
	userID, exist := c.Get("userID")
	if !exist {
		response.ErrorResponse(c.Writer, response.UnauthorizedMessage, http.StatusUnauthorized)
		return
	}

	var requestBody body.WalletStepUpRequest
	if err := c.ShouldBind(&requestBody); err != nil {
		response.ErrorResponse(c.Writer, response.BadRequestMessage, http.StatusBadRequest)
		return
	}

	invalidFields, err := requestBody.Validate()
	if err != nil {
		response.ErrorResponseData(c.Writer, invalidFields, response.UnprocessableEntityMessage, http.StatusUnprocessableEntity)
		return
	}

	token, err := h.userUC.WalletStepUp(c, userID.(string), requestBody)
	if err != nil {
		var e *httperror.Error
		if !errors.As(err, &e) {
			h.logger.Errorf("HandlerUser, Error: %s", err)
			response.ErrorResponse(c.Writer, response.InternalServerErrorMessage, http.StatusInternalServerError)
			return
		}

		response.ErrorResponse(c.Writer, e.Err.Error(), e.Status)
		return
	}

	c.SetSameSite(http.SameSiteNoneMode)
	c.SetCookie(constant.WalletTokenCookie, token, h.cfg.JWT.RefreshExpMin*60, "/", h.cfg.Server.Domain, true, true)
	response.SuccessResponse(c.Writer, nil, http.StatusOK)
}

func (h *userHandlers) CreateSLPPayment(c *gin.Context) {
	var requestBody body.CreatePaymentRequest
	if err := c.ShouldBind(&requestBody); err != nil {
		response.ErrorResponse(c.Writer, response.BadRequestMessage, http.StatusBadRequest)
		return
	}

	invalidFields, err := requestBody.Validate()
	if err != nil {
		response.ErrorResponseData(c.Writer, invalidFields, response.UnprocessableEntityMessage, http.StatusUnprocessableEntity)
		return
	}

	url, err := h.userUC.CreateSLPPayment(c, requestBody.TransactionID)
	if err != nil {
		var e *httperror.Error
		if !errors.As(err, &e) {
			h.logger.Errorf("HandlerUser, Error: %s", err)
			response.ErrorResponse(c.Writer, response.InternalServerErrorMessage, http.StatusInternalServerError)
			return
		}

		response.ErrorResponse(c.Writer, e.Err.Error(), e.Status)
		return
	}

	response.SuccessResponse(c.Writer, body.CreateSLPPaymentResponse{RedirectURL: url}, http.StatusOK)
}

func (h *userHandlers) CreateWalletPayment(c *gin.Context) {
	walletToken, err := c.Cookie(constant.WalletTokenCookie)
	if err != nil {
		response.ErrorResponse(c.Writer, response.ForbiddenMessage, http.StatusForbidden)
		return
	}

	claims, err := jwt.ExtractJWT(walletToken, h.cfg.JWT.JwtSecretKey)
	if err != nil {
		response.ErrorResponse(c.Writer, response.ForbiddenMessage, http.StatusForbidden)
		return
	}

	if claims["scope"].(string) != "level1" {
		response.ErrorResponse(c.Writer, response.ForbiddenMessage, http.StatusForbidden)
		return
	}

	var requestBody body.CreatePaymentRequest
	if errBind := c.ShouldBind(&requestBody); errBind != nil {
		response.ErrorResponse(c.Writer, response.BadRequestMessage, http.StatusBadRequest)
		return
	}

	invalidFields, err := requestBody.Validate()
	if err != nil {
		response.ErrorResponseData(c.Writer, invalidFields, response.UnprocessableEntityMessage, http.StatusUnprocessableEntity)
		return
	}

	if err := h.userUC.CreateWalletPayment(c, requestBody.TransactionID); err != nil {
		var e *httperror.Error
		if !errors.As(err, &e) {
			h.logger.Errorf("HandlerUser, Error: %s", err)
			response.ErrorResponse(c.Writer, response.InternalServerErrorMessage, http.StatusInternalServerError)
			return
		}

		response.ErrorResponse(c.Writer, e.Err.Error(), e.Status)
		return
	}

	c.SetSameSite(http.SameSiteNoneMode)
	c.SetCookie(constant.WalletTokenCookie, "", -1, "/", h.cfg.Server.Domain, true, true)
	response.SuccessResponse(c.Writer, nil, http.StatusOK)
}

func (h *userHandlers) SLPPaymentCallback(c *gin.Context) {
	var requestBody body.SLPCallbackRequest
	if err := c.ShouldBindJSON(&requestBody); err != nil {
		response.ErrorResponse(c.Writer, response.BadRequestMessage, http.StatusBadRequest)
		return
	}

	invalidFields, err := requestBody.Validate(h.cfg)
	if err != nil {
		response.ErrorResponseData(c.Writer, invalidFields, response.UnprocessableEntityMessage, http.StatusUnprocessableEntity)
		return
	}

	id := c.Param("id")
	transactionID, err := uuid.Parse(id)
	if err != nil {
		response.ErrorResponse(c.Writer, response.BadRequestMessage, http.StatusBadRequest)
		return
	}

	if err := h.userUC.UpdateTransaction(c, transactionID.String(), requestBody); err != nil {
		var e *httperror.Error
		if !errors.As(err, &e) {
			h.logger.Errorf("HandlerUser, Error: %s", err)
			response.ErrorResponse(c.Writer, response.InternalServerErrorMessage, http.StatusInternalServerError)
			return
		}

		response.ErrorResponse(c.Writer, e.Err.Error(), e.Status)
		return
	}

	response.SuccessResponse(c.Writer, nil, http.StatusOK)
}

func (h *userHandlers) WalletPaymentCallback(c *gin.Context) {
	var requestBody body.SLPCallbackRequest
	if err := c.ShouldBindJSON(&requestBody); err != nil {
		response.ErrorResponse(c.Writer, response.BadRequestMessage, http.StatusBadRequest)
		return
	}

	invalidFields, err := requestBody.Validate(h.cfg)
	if err != nil {
		response.ErrorResponseData(c.Writer, invalidFields, response.UnprocessableEntityMessage, http.StatusUnprocessableEntity)
		return
	}

	id := c.Param("id")
	transactionID, err := uuid.Parse(id)
	if err != nil {
		response.ErrorResponse(c.Writer, response.BadRequestMessage, http.StatusBadRequest)
		return
	}

	if err := h.userUC.UpdateWalletTransaction(c, transactionID.String(), requestBody); err != nil {
		var e *httperror.Error
		if !errors.As(err, &e) {
			h.logger.Errorf("HandlerUser, Error: %s", err)
			response.ErrorResponse(c.Writer, response.InternalServerErrorMessage, http.StatusInternalServerError)
			return
		}

		response.ErrorResponse(c.Writer, e.Err.Error(), e.Status)
		return
	}

	response.SuccessResponse(c.Writer, nil, http.StatusOK)
}

func (h *userHandlers) CreateTransaction(c *gin.Context) {
	var requestBody body.CreateTransactionRequest
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

	transactionID, err := h.userUC.CreateTransaction(c, userID.(string), requestBody)
	if err != nil {
		var e *httperror.Error
		if !errors.As(err, &e) {
			h.logger.Errorf("HandlerUser, Error: %s", err)
			response.ErrorResponse(c.Writer, response.InternalServerErrorMessage, http.StatusInternalServerError)
			return
		}

		response.ErrorResponse(c.Writer, e.Err.Error(), e.Status)
		return
	}

	response.SuccessResponse(c.Writer, body.CreateTransactionResponse{TransactionID: transactionID}, http.StatusOK)
}
