package delivery

import (
	"errors"
	"github.com/gin-gonic/gin"
	"murakali/config"
	"murakali/internal/module/location"
	"murakali/pkg/httperror"
	"murakali/pkg/logger"
	"murakali/pkg/response"
	"net/http"
)

type locationHandlers struct {
	cfg        *config.Config
	locationUC location.UseCase
	logger     logger.Logger
}

func NewLocationHandlers(cfg *config.Config, locationUC location.UseCase, log logger.Logger) location.Handlers {
	return &locationHandlers{cfg: cfg, locationUC: locationUC, logger: log}
}

func (h *locationHandlers) GetProvince(c *gin.Context) {
	province, err := h.locationUC.GetProvince(c)
	if err != nil {
		var e *httperror.Error
		if !errors.As(err, &e) {
			h.logger.Errorf("HandlerLocation, Error: %s", err)
			response.ErrorResponse(c.Writer, response.InternalServerErrorMessage, http.StatusInternalServerError)
			return
		}

		response.ErrorResponse(c.Writer, e.Err.Error(), e.Status)
		return
	}

	response.SuccessResponse(c.Writer, province, http.StatusOK)
}
