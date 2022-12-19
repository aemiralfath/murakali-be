package delivery

import (
	"murakali/config"
	"murakali/internal/module/user"
	"murakali/pkg/logger"
	"murakali/pkg/response"
	"net/http"

	"murakali/pkg/httperror"

	"errors"
	"github.com/gin-gonic/gin"

)

type userHandlers struct {
	cfg    *config.Config
	userUC user.UseCase
	logger logger.Logger
}

func NewUserHandlers(cfg *config.Config, userUC user.UseCase, log logger.Logger) user.Handlers {
	return &userHandlers{cfg: cfg, userUC: userUC, logger: log}
}

func (h *userHandlers) GetSealabsPay(c *gin.Context) {
	userid, exist := c.Get("userID")

	if !exist {
	   response.ErrorResponse(c.Writer, response.UnauthorizedMessage, http.StatusUnauthorized)
	   return
	}
	
	result, err := h.userUC.GetSealabsPay(c,userid.(string))
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
	response.SuccessResponse(c.Writer, result)

}
