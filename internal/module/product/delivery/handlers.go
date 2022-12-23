package delivery

import (
	"errors"
	"murakali/config"
	"murakali/internal/module/product"
	"murakali/internal/module/product/delivery/body"
	"murakali/pkg/httperror"
	"murakali/pkg/logger"
	"murakali/pkg/response"
	"net/http"

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

	RecommendedProducts, err := h.productUC.GetRecommendedProducts(c)
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
