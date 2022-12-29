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

	h.ValidateQueryOrder(c, pgn)

	orders, err := h.sellerUC.GetOrder(c, userID.(string), pgn)
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
			h.logger.Errorf("HandlerUser, Error: %s", err)
			response.ErrorResponse(c.Writer, response.InternalServerErrorMessage, http.StatusInternalServerError)
			return
		}

		response.ErrorResponse(c.Writer, e.Err.Error(), e.Status)
		return
	}

	response.SuccessResponse(c.Writer, nil, http.StatusOK)
}
