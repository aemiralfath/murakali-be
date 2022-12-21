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

	categoryList, err := h.productUC.GetCategories(c)
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

	categoryResponse := make([]*body.CategoryResponse, 0)
	for _, category := range categoryList {
		res := &body.CategoryResponse{
			ID:       category.ID,
			ParentID: category.ParentID,
			Name:     category.Name,
			PhotoURL: category.PhotoURL,
		}
		categoryResponse = append(categoryResponse, res)
	}
	response.SuccessResponse(c.Writer, categoryResponse, http.StatusOK)
}
