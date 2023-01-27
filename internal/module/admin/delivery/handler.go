package delivery

import (
	"errors"
	"murakali/config"
	"murakali/internal/constant"
	"murakali/internal/module/admin"
	"murakali/internal/module/admin/delivery/body"
	"murakali/internal/util"
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

func (h *adminHandlers) GetRefunds(c *gin.Context) {
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

	sortFilter = "accepted_at " + sortFilter
	refunds, err := h.adminUC.GetRefunds(c, sortFilter, pgn)
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

	response.SuccessResponse(c.Writer, refunds, http.StatusOK)
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

func (h *adminHandlers) RefundOrder(c *gin.Context) {
	id := c.Param("id")
	refundID, err := uuid.Parse(id)
	if err != nil {
		response.ErrorResponse(c.Writer, response.BadRequestMessage, http.StatusBadRequest)
		return
	}

	if err := h.adminUC.RefundOrder(c, refundID.String()); err != nil {
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

func (h *adminHandlers) GetCategories(c *gin.Context) {
	Categories, err := h.adminUC.GetCategories(c)
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

	response.SuccessResponse(c.Writer, Categories, http.StatusOK)
}

func (h *adminHandlers) UploadProductPicture(c *gin.Context) {
	type Sizer interface {
		Size() int64
	}

	var imgURL string
	var img body.ImageRequest

	err := c.ShouldBind(&img)
	if err != nil {
		response.ErrorResponse(c.Writer, response.InternalServerErrorMessage, http.StatusInternalServerError)
		return
	}
	data, _, err := c.Request.FormFile("Img")
	if err != nil {
		response.ErrorResponse(c.Writer, body.ImageIsEmpty, http.StatusInternalServerError)
		return
	}

	if data.(Sizer).Size() > constant.ImgMaxSize {
		response.ErrorResponse(c.Writer, response.PictureSizeTooBig, http.StatusInternalServerError)
		return
	}

	if data == nil {
		response.ErrorResponse(c.Writer, response.InternalServerErrorMessage, http.StatusInternalServerError)
		return
	}
	imgURL = util.UploadImageToCloudinary(c, h.cfg, data)

	response.SuccessResponse(c.Writer, imgURL, http.StatusOK)
}

func (h *adminHandlers) AddCategory(c *gin.Context) {
	var requestBody body.CategoryRequest
	if err := c.ShouldBind(&requestBody); err != nil {
		response.ErrorResponse(c.Writer, response.BadRequestMessage, http.StatusBadRequest)
		return
	}

	invalidFields, err := requestBody.Validate()
	if err != nil {
		response.ErrorResponseData(c.Writer, invalidFields, response.UnprocessableEntityMessage, http.StatusUnprocessableEntity)
		return
	}

	if err := h.adminUC.AddCategory(c, requestBody); err != nil {
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

func (h *adminHandlers) DeleteCategory(c *gin.Context) {
	id := c.Param("id")
	categoryID, err := uuid.Parse(id)
	if err != nil {
		response.ErrorResponse(c.Writer, response.BadRequestMessage, http.StatusBadRequest)
		return
	}

	if err := h.adminUC.DeleteCategory(c, categoryID.String()); err != nil {
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

func (h *adminHandlers) EditCategory(c *gin.Context) {
	var requestBody body.CategoryRequest
	if err := c.ShouldBind(&requestBody); err != nil {
		response.ErrorResponse(c.Writer, response.BadRequestMessage, http.StatusBadRequest)
		return
	}

	invalidFields, err := requestBody.Validate()
	if err != nil {
		response.ErrorResponseData(c.Writer, invalidFields, response.UnprocessableEntityMessage, http.StatusUnprocessableEntity)
		return
	}

	if err := h.adminUC.EditCategory(c, requestBody); err != nil {
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

func (h *adminHandlers) GetBanner(c *gin.Context) {
	banner, err := h.adminUC.GetBanner(c)
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

	response.SuccessResponse(c.Writer, banner, http.StatusOK)
}

func (h *adminHandlers) AddBanner(c *gin.Context) {
	var requestBody body.BannerRequest
	if err := c.ShouldBind(&requestBody); err != nil {
		response.ErrorResponse(c.Writer, response.BadRequestMessage, http.StatusBadRequest)
		return
	}

	invalidFields, err := requestBody.Validate()
	if err != nil {
		response.ErrorResponseData(c.Writer, invalidFields, response.UnprocessableEntityMessage, http.StatusUnprocessableEntity)
		return
	}

	if err := h.adminUC.AddBanner(c, requestBody); err != nil {
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

func (h *adminHandlers) DeleteBanner(c *gin.Context) {
	id := c.Param("id")
	bannerID, err := uuid.Parse(id)
	if err != nil {
		response.ErrorResponse(c.Writer, response.BadRequestMessage, http.StatusBadRequest)
		return
	}

	if err := h.adminUC.DeleteBanner(c, bannerID.String()); err != nil {
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

func (h *adminHandlers) EditBanner(c *gin.Context) {
	var requestBody body.BannerIDRequest
	if err := c.ShouldBind(&requestBody); err != nil {
		response.ErrorResponse(c.Writer, response.BadRequestMessage, http.StatusBadRequest)
		return
	}

	invalidFields, err := requestBody.IDValidate()
	if err != nil {
		response.ErrorResponseData(c.Writer, invalidFields, response.UnprocessableEntityMessage, http.StatusUnprocessableEntity)
		return
	}

	if err := h.adminUC.EditBanner(c, requestBody); err != nil {
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
