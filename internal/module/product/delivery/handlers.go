package delivery

import (
	"errors"
	"murakali/config"
	"murakali/internal/module/product"
	"murakali/internal/module/product/delivery/body"
	"murakali/pkg/httperror"
	"murakali/pkg/logger"
	"murakali/pkg/pagination"
	"murakali/pkg/response"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type productHandlers struct {
	cfg       *config.Config
	productUC product.UseCase
	logger    logger.Logger
}

func NewProductHandlers(cfg *config.Config, productUC product.UseCase, log logger.Logger) product.Handlers {
	return &productHandlers{cfg: cfg, productUC: productUC, logger: log}
}

func (h *productHandlers) GetCategories(c *gin.Context) {
	categoriesResponse, err := h.productUC.GetCategories(c)
	if err != nil {
		var e *httperror.Error
		if !errors.As(err, &e) {
			h.logger.Errorf("HandlerProduct, Error: %s", err)
			response.ErrorResponse(c.Writer, response.InternalServerErrorMessage, http.StatusInternalServerError)
			return
		}

		response.ErrorResponse(c.Writer, e.Err.Error(), e.Status)
		return
	}

	response.SuccessResponse(c.Writer, categoriesResponse, http.StatusOK)
}

func (h *productHandlers) GetBanners(c *gin.Context) {
	banners, err := h.productUC.GetBanners(c)
	if err != nil {
		var e *httperror.Error
		if !errors.As(err, &e) {
			h.logger.Errorf("HandlerAuth, Error: %s", err)
			response.ErrorResponse(c.Writer, response.InternalServerErrorMessage, http.StatusInternalServerError)
			return
		}

		response.ErrorResponse(c.Writer, e.Err.Error(), e.Status)
		return
	}

	response.SuccessResponse(c.Writer, banners, http.StatusOK)
}

func (h *productHandlers) GetCategoriesByNameLevelOne(c *gin.Context) {
	var requestPath body.CategoryRequest
	requestPath.NameLevelOne = c.Param("name_lvl_one")

	categoriesResponse, err := h.productUC.GetCategoriesByName(c, requestPath.NameLevelOne)
	if err != nil {
		var e *httperror.Error
		if !errors.As(err, &e) {
			h.logger.Errorf("HandlerProduct, Error: %s", err)
			response.ErrorResponse(c.Writer, response.InternalServerErrorMessage, http.StatusInternalServerError)
			return
		}

		response.ErrorResponse(c.Writer, e.Err.Error(), e.Status)
		return
	}

	response.SuccessResponse(c.Writer, categoriesResponse, http.StatusOK)
}

func (h *productHandlers) GetCategoriesByNameLevelTwo(c *gin.Context) {
	var requestPath body.CategoryRequest
	requestPath.NameLevelTwo = c.Param("name_lvl_two")

	categoriesResponse, err := h.productUC.GetCategoriesByName(c, requestPath.NameLevelTwo)
	if err != nil {
		var e *httperror.Error
		if !errors.As(err, &e) {
			h.logger.Errorf("HandlerProduct, Error: %s", err)
			response.ErrorResponse(c.Writer, response.InternalServerErrorMessage, http.StatusInternalServerError)
			return
		}

		response.ErrorResponse(c.Writer, e.Err.Error(), e.Status)
		return
	}

	response.SuccessResponse(c.Writer, categoriesResponse, http.StatusOK)
}

func (h *productHandlers) GetCategoriesByNameLevelThree(c *gin.Context) {
	var requestPath body.CategoryRequest
	requestPath.NameLevelThree = c.Param("name_lvl_three")

	categoriesResponse, err := h.productUC.GetCategoriesByName(c, requestPath.NameLevelThree)
	if err != nil {
		var e *httperror.Error
		if !errors.As(err, &e) {
			h.logger.Errorf("HandlerProduct, Error: %s", err)
			response.ErrorResponse(c.Writer, response.InternalServerErrorMessage, http.StatusInternalServerError)
			return
		}

		response.ErrorResponse(c.Writer, e.Err.Error(), e.Status)
		return
	}

	response.SuccessResponse(c.Writer, categoriesResponse, http.StatusOK)
}

func (h *productHandlers) GetRecommendedProducts(c *gin.Context) {
	pgn := h.ValidateQueryRecommendProduct(c)
	RecommendedProducts, err := h.productUC.GetRecommendedProducts(c, pgn)
	if err != nil {
		var e *httperror.Error
		if !errors.As(err, &e) {
			h.logger.Errorf("HandlerProduct, Error: %s", err)
			response.ErrorResponse(c.Writer, response.InternalServerErrorMessage, http.StatusInternalServerError)
			return
		}

		response.ErrorResponse(c.Writer, e.Err.Error(), e.Status)
		return
	}

	response.SuccessResponse(c.Writer, RecommendedProducts, http.StatusOK)
}

func (h *productHandlers) ValidateQueryRecommendProduct(c *gin.Context) *pagination.Pagination {
	limit := strings.TrimSpace(c.Query("limit"))
	page := strings.TrimSpace(c.Query("page"))

	var limitFilter int
	var pageFilter int

	limitFilter, err := strconv.Atoi(limit)
	if err != nil || limitFilter < 1 {
		limitFilter = 18
	}

	pageFilter, err = strconv.Atoi(page)
	if err != nil || pageFilter < 1 {
		pageFilter = 1
	}
	pgn := &pagination.Pagination{
		Limit: limitFilter,
		Page:  pageFilter,
		Sort:  "unit_sold DESC",
	}

	return pgn
}
