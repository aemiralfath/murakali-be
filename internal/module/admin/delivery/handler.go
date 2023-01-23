package delivery

import (
	"errors"
	"murakali/config"
	"murakali/internal/constant"
	"murakali/internal/module/admin"
	"murakali/internal/module/admin/delivery/body"
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

type adminHandlers struct {
	cfg     *config.Config
	adminUC admin.UseCase
	logger  logger.Logger
}

func NewAdminHandlers(cfg *config.Config, adminUC admin.UseCase, log logger.Logger) admin.Handlers {
	return &adminHandlers{cfg: cfg, adminUC: adminUC, logger: log}
}

func (h *adminHandlers) GetAllVoucher(c *gin.Context) {
	pgn := &pagination.Pagination{}
	h.ValidateQueryPagination(c, pgn)

	sort := c.DefaultQuery("sort", "")
	sort = strings.ToLower(sort)
	var sortFilter string
	switch sort {
	case constant.ASC:
		sortFilter = sort
	default:
		sortFilter = constant.DESC
	}

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

	sortFilter = "v." + "created_at " + sortFilter
	shopVouchers, err := h.adminUC.GetAllVoucher(c, voucherStatusID, sortFilter, pgn)
	if err != nil {
		var e *httperror.Error
		if !errors.As(err, &e) {
			h.logger.Errorf("HandlerAdmin, Error: %s", err)
			response.ErrorResponse(c.Writer, response.InternalServerErrorMessage, http.StatusInternalServerError)
			return
		}

		response.ErrorResponse(c.Writer, e.Err.Error(), e.Status)
		return
	}

	response.SuccessResponse(c.Writer, shopVouchers, http.StatusOK)
}

func (h *adminHandlers) CreateVoucher(c *gin.Context) {
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

	if err := h.adminUC.CreateVoucher(c, requestBody); err != nil {
		var e *httperror.Error
		if !errors.As(err, &e) {
			h.logger.Errorf("HandlerAdmin, Error: %s", err)
			response.ErrorResponse(c.Writer, response.InternalServerErrorMessage, http.StatusInternalServerError)
			return
		}

		response.ErrorResponse(c.Writer, e.Err.Error(), e.Status)
		return
	}

	response.SuccessResponse(c.Writer, nil, http.StatusOK)
}

func (h *adminHandlers) DeleteVoucher(c *gin.Context) {
	id := c.Param("id")
	voucherID, err := uuid.Parse(id)
	if err != nil {
		response.ErrorResponse(c.Writer, response.BadRequestMessage, http.StatusBadRequest)
		return
	}

	if err := h.adminUC.DeleteVoucher(c, voucherID.String()); err != nil {
		var e *httperror.Error
		if !errors.As(err, &e) {
			h.logger.Errorf("HandlerAdmin, Error: %s", err)
			response.ErrorResponse(c.Writer, response.InternalServerErrorMessage, http.StatusInternalServerError)
			return
		}

		response.ErrorResponse(c.Writer, e.Err.Error(), e.Status)
		return
	}

	response.SuccessResponse(c.Writer, nil, http.StatusOK)
}

func (h *adminHandlers) UpdateVoucher(c *gin.Context) {
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

	if err := h.adminUC.UpdateVoucher(c, requestBody); err != nil {
		var e *httperror.Error
		if !errors.As(err, &e) {
			h.logger.Errorf("HandlerAdmin, Error: %s", err)
			response.ErrorResponse(c.Writer, response.InternalServerErrorMessage, http.StatusInternalServerError)
			return
		}

		response.ErrorResponse(c.Writer, e.Err.Error(), e.Status)
		return
	}

	response.SuccessResponse(c.Writer, nil, http.StatusOK)
}

func (h *adminHandlers) GetDetailVoucher(c *gin.Context) {
	id := c.Param("id")
	voucherID, err := uuid.Parse(id)
	if err != nil {
		response.ErrorResponse(c.Writer, response.BadRequestMessage, http.StatusBadRequest)
		return
	}

	voucherShop, err := h.adminUC.GetDetailVoucher(c, voucherID.String())
	if err != nil {
		var e *httperror.Error
		if !errors.As(err, &e) {
			h.logger.Errorf("HandlerAdmin, Error: %s", err)
			response.ErrorResponse(c.Writer, response.InternalServerErrorMessage, http.StatusInternalServerError)
			return
		}

		response.ErrorResponse(c.Writer, e.Err.Error(), e.Status)
		return
	}

	response.SuccessResponse(c.Writer, voucherShop, http.StatusOK)
}

func (h *adminHandlers) ValidateQueryPagination(c *gin.Context, pgn *pagination.Pagination) {
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
